package translation

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// Translator is the central translation component, mirroring Laravel's
// Illuminate\Translation\Translator.
//
// It resolves translation keys against a Loader, applies placeholder
// substitution atomically (longest match first), supports namespace/fallback
// locale chains, and delegates pluralisation to a MessageSelector.
type Translator struct {
	loader      Loader
	locale      string
	fallback    string
	loaded      map[string]map[string]map[string]map[string]any // [ns][locale][group]→messages
	selector    *MessageSelector
	stringables map[reflect.Type]func(any) string

	missingKeyFn     func(key, locale, namespace, group string)
	determineLocales func([]string) []string
}

// NewTranslator returns an initialised *Translator backed by loader and using
// locale as the default locale.
func NewTranslator(loader Loader, locale string) *Translator {
	return &Translator{
		loader:      loader,
		locale:      locale,
		loaded:      make(map[string]map[string]map[string]map[string]any),
		selector:    NewMessageSelector(),
		stringables: make(map[reflect.Type]func(any) string),
	}
}

// ── Locale management ─────────────────────────────────────────────────────

// GetLocale returns the currently active locale.
func (t *Translator) GetLocale() string { return t.locale }

// SetLocale changes the active locale.  Locales containing "/" or "\" are
// silently rejected to prevent directory-traversal attacks.
func (t *Translator) SetLocale(locale string) {
	if strings.ContainsAny(locale, "/\\") {
		return
	}

	t.locale = locale
}

// SetFallback sets the fallback locale used when a key is not found in the
// primary locale.
func (t *Translator) SetFallback(fallback string) { t.fallback = fallback }

// GetFallback returns the fallback locale.
func (t *Translator) GetFallback() string { return t.fallback }

// ── Loader delegation ─────────────────────────────────────────────────────

// GetLoader returns the underlying Loader.
func (t *Translator) GetLoader() Loader { return t.loader }

// AddNamespace registers a package namespace with a hint path on the loader.
func (t *Translator) AddNamespace(namespace, hint string) {
	t.loader.AddNamespace(namespace, hint)
}

// AddPath appends a translation directory on a FileLoader.  No-op for other
// loader types.
func (t *Translator) AddPath(path string) {
	if fl, ok := t.loader.(*FileLoader); ok {
		fl.AddPath(path)
	}
}

// AddJsonPath registers an additional flat-JSON locale directory.
func (t *Translator) AddJsonPath(path string) {
	t.loader.AddJsonPath(path)
}

// ── Lines ─────────────────────────────────────────────────────────────────

// AddLines injects translation lines into the loaded cache for locale.
// namespace defaults to "*" if empty.
func (t *Translator) AddLines(lines map[string]any, locale string, namespace ...string) {
	ns := globalNS

	if len(namespace) > 0 && namespace[0] != "" {
		ns = namespace[0]
	}

	t.ensureLoaded(ns, locale, globalNS)

	existing := t.loaded[ns][locale][globalNS]
	t.loaded[ns][locale][globalNS] = mergeMaps(existing, lines)
}

// ── Has ───────────────────────────────────────────────────────────────────

// Has reports whether a translation key exists (using locale, then fallback).
func (t *Translator) Has(key string, locale *string) bool {
	result := t.Get(key, nil, locale)
	s, ok := result.(string)

	return ok && s != key
}

// HasForLocale reports whether a translation key exists in locale only,
// without falling back.
func (t *Translator) HasForLocale(key string, locale *string) bool {
	loc := t.resolveLocale(locale)
	ns, group, item := t.ParseKey(key)

	return t.getLine(ns, group, item, loc) != nil
}

// ── Get ───────────────────────────────────────────────────────────────────

// Get retrieves a translation for key, applying replace substitutions.
// Returns the key itself when no translation is found.
// An optional boolean fourth parameter (default true) controls whether the
// fallback locale is tried.
func (t *Translator) Get(key string, replace map[string]any, locale *string, fallback ...bool) any {
	useFallback := true

	if len(fallback) > 0 {
		useFallback = fallback[0]
	}

	ns, group, item := t.ParseKey(key)
	loc := t.resolveLocale(locale)

	locales := t.localeArray(loc)

	if !useFallback {
		locales = []string{loc}
	}

	for _, l := range locales {
		line := t.getLine(ns, group, item, l)

		if line != nil {
			return t.makeReplacementsFor(line, replace)
		}
	}

	if t.missingKeyFn != nil {
		t.missingKeyFn(key, loc, ns, group)
	}

	// No translation found: apply replacements to the key itself.
	return t.makeReplacementsFor(key, replace)
}

// ── Choice ────────────────────────────────────────────────────────────────

// Choice retrieves a pluralised translation for key based on number.
// number may be int, float64, a slice/array, or a Countable.
func (t *Translator) Choice(key string, number any, replace map[string]any, locale *string) string {
	n := toFloat(number)
	loc := t.localeForChoice(key, locale)

	line := t.Get(key, nil, &loc)
	lineStr, _ := line.(string)

	chosen := t.getSelector().Choose(lineStr, n, loc)

	// Merge :count into replacements; user-supplied value wins.
	merged := make(map[string]any, len(replace)+1)
	merged["count"] = n

	for k, v := range replace {
		merged[k] = v
	}

	return t.MakeReplacements(chosen, merged)
}

// localeForChoice returns the locale to use for choice pluralisation:
// the requested locale if it has the key, otherwise the fallback.
func (t *Translator) localeForChoice(key string, locale *string) string {
	loc := t.resolveLocale(locale)

	if t.HasForLocale(key, &loc) {
		return loc
	}

	if t.fallback != "" {
		return t.fallback
	}

	return loc
}

// ── Key parsing ───────────────────────────────────────────────────────────

// ParseKey parses a translation key into (namespace, group, item).
//
//   - "vendor::messages.welcome" → ("vendor", "messages", "welcome")
//   - "vendor::messages"         → ("vendor", "messages", "")    ← whole group
//   - "messages.welcome"         → ("*", "messages", "welcome")
//   - "welcome"                  → ("*", "*", "welcome")         ← flat JSON
func (t *Translator) ParseKey(key string) (namespace, group, item string) {
	if idx := strings.Index(key, "::"); idx >= 0 {
		namespace = key[:idx]
		remaining := key[idx+2:]

		if dot := strings.Index(remaining, "."); dot >= 0 {
			group = remaining[:dot]
			item = remaining[dot+1:]
		} else {
			// Namespaced, no dot → return whole group.
			group = remaining
			item = ""
		}
	} else {
		namespace = globalNS

		if dot := strings.Index(key, "."); dot >= 0 {
			group = key[:dot]
			item = key[dot+1:]
		} else {
			// No namespace, no dot → flat JSON lookup.
			group = globalNS
			item = key
		}
	}

	return
}

// ── Replacements ─────────────────────────────────────────────────────────

// MakeReplacements performs atomic placeholder substitution on line.
//
// Replacement values that are func(string) string are treated as tag
// handlers: patterns <key>content</key> are replaced by handler(content).
//
// For remaining :placeholder tokens, substitution is done in a single atomic
// pass (longest key first) using temporary sentinel tokens so that a
// replacement value containing ":something" does not trigger further
// substitution.
func (t *Translator) MakeReplacements(line string, replace map[string]any) string {
	if len(replace) == 0 {
		return line
	}

	// Phase 1: tag replacements (func(string) string values).
	line = applyTagReplacements(line, replace)

	// Phase 2: collect :placeholder pairs (non-func values).
	type pair struct {
		key string
		val string
	}

	var pairs []pair

	for k, v := range replace {
		if _, isFn := v.(func(string) string); isFn {
			continue
		}

		pairs = append(pairs, pair{key: k, val: t.stringValue(v)})
	}

	// Sort longest key first to avoid partial substitution of shorter keys.
	sort.Slice(pairs, func(i, j int) bool {
		return len(pairs[i].key) > len(pairs[j].key)
	})

	// Phase 3: first pass — replace :Key/:key/:KEY with sentinels.
	// UC/UU variants are only injected when they produce a different pattern
	// from the base key (e.g. numeric keys "0" must not receive UC treatment).
	const sentinelFmt = "\x00REPL%d\x00"

	for i, p := range pairs {
		sentinel := fmt.Sprintf(sentinelFmt, i)

		ucKey := ucFirst(p.key)
		uuKey := strings.ToUpper(p.key)

		if ucKey != p.key {
			line = strings.ReplaceAll(line, ":"+ucKey, sentinel+"_UC")
		}

		if uuKey != p.key && uuKey != ucKey {
			line = strings.ReplaceAll(line, ":"+uuKey, sentinel+"_UU")
		}

		line = strings.ReplaceAll(line, ":"+p.key, sentinel)
	}

	// Phase 4: second pass — swap sentinels for final values.
	for i, p := range pairs {
		sentinel := fmt.Sprintf(sentinelFmt, i)
		line = strings.ReplaceAll(line, sentinel+"_UC", ucFirst(p.val))
		line = strings.ReplaceAll(line, sentinel+"_UU", strings.ToUpper(p.val))
		line = strings.ReplaceAll(line, sentinel, p.val)
	}

	return line
}

// makeReplacementsFor converts a raw translation value (string or []string)
// to a string and applies replacements.  Arrays have each element processed.
func (t *Translator) makeReplacementsFor(line any, replace map[string]any) any {
	switch v := line.(type) {
	case string:
		return t.MakeReplacements(v, replace)
	case map[string]any:
		result := make(map[string]any, len(v))

		for k, val := range v {
			result[k] = t.makeReplacementsFor(val, replace)
		}

		return result
	case []any:
		result := make([]any, len(v))

		for i, val := range v {
			result[i] = t.makeReplacementsFor(val, replace)
		}

		return result
	default:
		return line
	}
}

// ── Internal helpers ──────────────────────────────────────────────────────

// getLine retrieves a raw translation value for (ns, group, item, locale).
// Returns nil when the key does not exist.
func (t *Translator) getLine(ns, group, item, locale string) any {
	t.load(ns, group, locale)

	msgs := t.loaded[ns][locale][group]

	if msgs == nil || len(msgs) == 0 {
		return nil
	}

	if item == "" || item == globalNS {
		return msgs
	}

	// Navigate dot-notation.
	return dotGet(msgs, item)
}

// load ensures the (ns, group, locale) triple is present in the loaded cache.
// It records a non-nil (possibly empty) map even when the loader returns
// nothing, preventing repeated load calls for missing keys.
func (t *Translator) load(ns, group, locale string) {
	if t.loaded[ns] != nil &&
		t.loaded[ns][locale] != nil &&
		t.loaded[ns][locale][group] != nil {
		return
	}

	t.ensureLoaded(ns, locale, group)

	var nsPtr *string

	if ns != globalNS {
		nsPtr = &ns
	}

	msgs := t.loader.Load(locale, group, nsPtr)

	if msgs == nil {
		msgs = map[string]any{}
	}

	t.loaded[ns][locale][group] = msgs
}

// ensureLoaded initialises the nested maps in t.loaded up to (ns, locale, group).
func (t *Translator) ensureLoaded(ns, locale, group string) {
	if t.loaded[ns] == nil {
		t.loaded[ns] = make(map[string]map[string]map[string]any)
	}

	if t.loaded[ns][locale] == nil {
		t.loaded[ns][locale] = make(map[string]map[string]any)
	}

	if t.loaded[ns][locale][group] == nil {
		t.loaded[ns][locale][group] = map[string]any{}
	}
}

// localeArray builds the ordered list of locales to attempt for a key lookup.
func (t *Translator) localeArray(locale string) []string {
	locales := []string{locale}

	if t.fallback != "" && t.fallback != locale {
		locales = append(locales, t.fallback)
	}

	if t.determineLocales != nil {
		locales = t.determineLocales(locales)
	}

	return locales
}

// resolveLocale returns locale if non-nil, else t.locale.
func (t *Translator) resolveLocale(locale *string) string {
	if locale != nil {
		return *locale
	}

	return t.locale
}

// stringValue converts a replacement value to its string representation.
// Priority: registered Stringable handler → fmt.Stringer → fmt.Sprintf("%v").
// nil → "".
func (t *Translator) stringValue(v any) string {
	if v == nil {
		return ""
	}

	rt := reflect.TypeOf(v)

	if fn, ok := t.stringables[rt]; ok {
		return fn(v)
	}

	if s, ok := v.(fmt.Stringer); ok {
		return s.String()
	}

	switch c := v.(type) {
	case float64:
		return floatToStr(c)
	case float32:
		return floatToStr(float64(c))
	}

	return fmt.Sprintf("%v", v)
}

// floatToStr formats a float without trailing zeros when it is whole.
func floatToStr(f float64) string {
	if f == float64(int64(f)) {
		return fmt.Sprintf("%d", int64(f))
	}

	return fmt.Sprintf("%g", f)
}

// ── Customisation hooks ───────────────────────────────────────────────────

// GetSelector returns the MessageSelector.
func (t *Translator) GetSelector() *MessageSelector { return t.selector }

// SetSelector replaces the MessageSelector.
func (t *Translator) SetSelector(s *MessageSelector) { t.selector = s }

func (t *Translator) getSelector() *MessageSelector {
	if t.selector == nil {
		t.selector = NewMessageSelector()
	}

	return t.selector
}

// Stringable registers a handler for converting values of type T to string
// during placeholder substitution.
func (t *Translator) Stringable(typ reflect.Type, handler func(any) string) {
	t.stringables[typ] = handler
}

// HandleMissingKeysUsing sets a callback invoked when a key is not found
// after all locale fallbacks are exhausted.
func (t *Translator) HandleMissingKeysUsing(fn func(key, locale, namespace, group string)) {
	t.missingKeyFn = fn
}

// DetermineLocalesUsing sets a callback that receives the initial locale list
// (primary + fallback) and returns the final ordered list to try.
func (t *Translator) DetermineLocalesUsing(fn func([]string) []string) {
	t.determineLocales = fn
}

// ── Package-level helpers ─────────────────────────────────────────────────

// dotGet navigates a nested map[string]any using dot-notation key.
func dotGet(m map[string]any, key string) any {
	parts := strings.SplitN(key, ".", 2)
	val, ok := m[parts[0]]

	if !ok {
		return nil
	}

	if len(parts) == 1 {
		return val
	}

	sub, ok := val.(map[string]any)

	if !ok {
		return nil
	}

	return dotGet(sub, parts[1])
}

// ucFirst returns s with its first Unicode letter upper-cased.
func ucFirst(s string) string {
	if s == "" {
		return s
	}

	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])

	return string(r)
}

// applyTagReplacements processes <tag>content</tag> patterns in line using
// func(string) string values from replace.  All occurrences of each tag are
// replaced.
func applyTagReplacements(line string, replace map[string]any) string {
	for k, v := range replace {
		fn, ok := v.(func(string) string)

		if !ok {
			continue
		}
		// Match <k>...</k> non-greedily, allowing multiline content.
		pattern := regexp.MustCompile(`(?s)<` + regexp.QuoteMeta(k) + `>(.*?)</` + regexp.QuoteMeta(k) + `>`)
		line = pattern.ReplaceAllStringFunc(line, func(match string) string {
			sub := pattern.FindStringSubmatch(match)

			if len(sub) < 2 {
				return match
			}

			return fn(sub[1])
		})
	}

	return line
}

// toFloat converts a number-like value to float64 for pluralisation.
func toFloat(number any) float64 {
	switch v := number.(type) {
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case float32:
		return float64(v)
	case float64:
		return v
	}

	if c, ok := number.(Countable); ok {
		return float64(c.Len())
	}

	rv := reflect.ValueOf(number)

	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		return float64(rv.Len())
	}

	return 0
}
