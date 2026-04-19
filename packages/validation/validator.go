package validation

import (
	"strings"

	contract "github.com/bedrock/packages/contracts/validation"
	"github.com/bedrock/packages/validation/rules"
)

// RuleFunc is a validation function.  It receives the attribute name, the
// value under validation, any rule parameters, and a RuleContext.  It returns
// true when the value is valid.
type RuleFunc = rules.RuleFunc

// RuleContext is passed to every RuleFunc and provides access to the
// surrounding validation state.
type RuleContext = rules.RuleContext

// PresenceVerifier handles database-level constraints (exists, unique).
type PresenceVerifier = rules.PresenceVerifier

// Validator evaluates a set of rules against a data map and accumulates
// failures into a MessageBag.
type Validator struct {
	data             map[string]any
	flatData         map[string]any
	rawRules         map[string]any
	parsedRules      map[string][]ParsedRule
	customMessages   map[string]string
	attrNames        map[string]string
	errs             *MessageBag
	extensions       map[string]rules.RuleFunc
	implicitExt      map[string]bool
	presenceVerifier rules.PresenceVerifier
	failedRules      map[string][]string // attribute → []ruleName
	validated        map[string]any      // populated after Passes() succeeds
	hasRun           bool

	// excludedAttrs holds attributes marked for exclusion by Exclude* rules
	excludedAttrs map[string]bool
}

// newValidator creates a Validator from already-parsed state.

// Passes runs all validation rules and returns true when every rule passes.

// Fails returns true when at least one rule fails.

// Validate runs all rules and returns a *ValidationException when any rule
// fails.

// Validated returns only the attributes that passed validation.  Returns an
// error if validation has not run or if it failed.

// Safe returns a ValidatedInput containing only the validated attributes.

// Errors returns the MessageBag of accumulated failures.

// Failed returns a map of attribute → []ruleName for every failed rule.

// SetData replaces the data under validation.

// SetRules replaces the rule set.

// AddRules merges additional rules into the existing set.

// GetData returns the data under validation.

// GetRules returns the parsed rule map.

// HasRule reports whether the attribute has any of the named rules.

// SetCustomMessages registers custom error messages.

// SetAttributeNames registers human-readable attribute names used in messages.

// AddExtension registers a custom rule function.

// AddImplicitExtension registers a custom implicit rule function.

// SetPresenceVerifier sets the database presence verifier.

// ─── internal ──────────────────────────────────────────────────────────────────

// Expand wildcards to concrete attributes

// Even if not in data, we still need to run implicit rules

// Handle special marker rules

// skip remaining rules for null with nullable

// field not present at all — skip

// Check for exclusion rules

// Skip non-implicit rules when value is absent/blank

// unknown rule — pass

// User extensions take precedence over built-ins

// 1. Custom message keyed by "attribute.rule"

// 2. Custom message keyed by just "rule"

// 3. Type-qualified default (e.g. "Max.string")

// 4. Default message

// Value of the attribute

// Params

// Common single-param rules

// Convert dot-notation and underscores to spaces

// getValue returns the value for attribute using dot-notation lookup.

// Try nested map traversal

// flatDataHas reports whether the attribute key exists in the flattened data.

// dotGet traverses a nested map using dot notation.

// ruleContext implements rules.RuleContext, used internally by rule functions.
type ruleContext struct {
	validator     *Validator
	attribute     string
	customMessage string
}

var _ contract.Validator = (*Validator)(nil)

func newValidator(
	data map[string]any,
	rawRules map[string]any,
	customMessages map[string]string,
	attrNames map[string]string,
	extensions map[string]rules.RuleFunc,
	implicitExt map[string]bool,
	pv rules.PresenceVerifier,
) *Validator {
	v := &Validator{
		data:             data,
		rawRules:         rawRules,
		customMessages:   customMessages,
		attrNames:        attrNames,
		errs:             NewMessageBag(),
		extensions:       extensions,
		implicitExt:      implicitExt,
		presenceVerifier: pv,
		failedRules:      make(map[string][]string),
		excludedAttrs:    make(map[string]bool),
	}

	if v.data == nil {
		v.data = map[string]any{}
	}

	if v.rawRules == nil {
		v.rawRules = map[string]any{}
	}

	if v.customMessages == nil {
		v.customMessages = map[string]string{}
	}

	if v.attrNames == nil {
		v.attrNames = map[string]string{}
	}

	if v.extensions == nil {
		v.extensions = map[string]rules.RuleFunc{}
	}

	if v.implicitExt == nil {
		v.implicitExt = map[string]bool{}
	}

	v.flatData = FlattenData(v.data)
	v.parsedRules = v.parseAllRules()

	return v
}

func (v *Validator) Passes() bool {
	v.errs = NewMessageBag()
	v.failedRules = make(map[string][]string)
	v.excludedAttrs = make(map[string]bool)
	v.validated = nil
	v.hasRun = true

	v.validateAll()

	if v.errs.IsEmpty() {
		v.buildValidated()
	}

	return v.errs.IsEmpty()
}

func (v *Validator) Fails() bool {
	return !v.Passes()
}

func (v *Validator) Validate() error {
	if v.Passes() {
		return nil
	}

	return &ValidationException{Bag: v.errs}
}

func (v *Validator) Validated() (map[string]any, error) {
	if !v.hasRun {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if !v.errs.IsEmpty() {
		return nil, &ValidationException{Bag: v.errs}
	}

	return v.validated, nil
}

func (v *Validator) Safe(keys ...string) (*ValidatedInput, error) {
	vd, err := v.Validated()

	if err != nil {
		return nil, err
	}

	if len(keys) > 0 {
		filtered := make(map[string]any, len(keys))

		for _, k := range keys {
			if val, ok := vd[k]; ok {
				filtered[k] = val
			}
		}

		return &ValidatedInput{data: filtered}, nil
	}

	return &ValidatedInput{data: vd}, nil
}

func (v *Validator) Errors() contract.MessageBag {
	return v.errs
}

func (v *Validator) Failed() map[string][]string {
	return v.failedRules
}

func (v *Validator) SetData(data map[string]any) contract.Validator {
	v.data = data
	v.flatData = FlattenData(data)
	v.hasRun = false

	return v
}

func (v *Validator) SetRules(ruleMap map[string]any) contract.Validator {
	v.rawRules = ruleMap
	v.parsedRules = v.parseAllRules()
	v.hasRun = false

	return v
}

func (v *Validator) AddRules(ruleMap map[string]any) contract.Validator {
	for k, r := range ruleMap {
		existing := Explode(v.rawRules[k])
		v.rawRules[k] = append(existing, Explode(r)...)
	}

	v.parsedRules = v.parseAllRules()
	v.hasRun = false

	return v
}

func (v *Validator) GetData() map[string]any {
	return v.data
}

func (v *Validator) GetRules() map[string][]ParsedRule {
	return v.parsedRules
}

func (v *Validator) HasRule(attribute string, ruleNames ...string) bool {
	pr, ok := v.parsedRules[attribute]

	if !ok {
		return false
	}

	for _, rule := range pr {
		for _, name := range ruleNames {
			if strings.EqualFold(rule.Name, name) {
				return true
			}
		}
	}

	return false
}

func (v *Validator) SetCustomMessages(messages map[string]string) contract.Validator {
	for k, msg := range messages {
		v.customMessages[k] = msg
	}

	return v
}

func (v *Validator) SetAttributeNames(names map[string]string) contract.Validator {
	for k, name := range names {
		v.attrNames[k] = name
	}

	return v
}

func (v *Validator) AddExtension(name string, fn rules.RuleFunc) *Validator {
	v.extensions[StudlyCase(name)] = fn

	return v
}

func (v *Validator) AddImplicitExtension(name string, fn rules.RuleFunc) *Validator {
	studly := StudlyCase(name)
	v.extensions[studly] = fn
	v.implicitExt[studly] = true

	return v
}

func (v *Validator) SetPresenceVerifier(pv rules.PresenceVerifier) *Validator {
	v.presenceVerifier = pv

	return v
}

func (v *Validator) validateAll() {
	for attribute, parsedRuleList := range v.parsedRules {

		attrs := ExpandWildcards(attribute, v.flatData)

		if len(attrs) == 0 {

			attrs = []string{attribute}
		}

		for _, attr := range attrs {
			v.validateAttributeRules(attr, parsedRuleList)
		}
	}
}

func (v *Validator) validateAttributeRules(attribute string, parsedRuleList []ParsedRule) {
	value := v.getValue(attribute)
	bail := false

	for _, rule := range parsedRuleList {
		if rule.IsObject() {
			v.validateRuleObject(attribute, value, rule)

			if bail && v.errs.Has(attribute) {
				return
			}

			continue
		}

		name := rule.Name

		switch name {
		case "Bail":
			bail = true

			continue
		case "Nullable":
			if value == nil {
				return
			}

			continue
		case "Sometimes":
			if !v.flatDataHas(attribute) {
				return
			}

			continue
		}

		if isExcludeRule(name) {
			ctx := v.makeContext(attribute)

			if rules.ShouldExclude(attribute, name, rule.Parameters, ctx) {
				v.excludedAttrs[attribute] = true
			}

			continue
		}

		if !v.isImplicitRule(name) && !v.flatDataHas(attribute) {
			continue
		}

		if !v.runRule(attribute, value, rule) {
			if bail {
				return
			}
		}
	}
}

func (v *Validator) validateRuleObject(attribute string, value any, rule ParsedRule) {
	obj := rule.Object
	ctx := v.makeContext(attribute)
	failed := false

	var customMsg string

	obj.Validate(attribute, value, func(message string) {
		failed = true
		customMsg = message
	})

	if failed {
		if customMsg != "" {
			ctx.SetMessage(customMsg)
		}

		v.addFailureWithContext(attribute, ruleObjectName(rule), rule.Parameters, ctx)
	}
}

func (v *Validator) runRule(attribute string, value any, rule ParsedRule) bool {
	fn := v.lookupRule(rule.Name)

	if fn == nil {
		return true
	}

	ctx := v.makeContext(attribute)
	passed := fn(attribute, value, rule.Parameters, ctx)

	if !passed {
		v.addFailureWithContext(attribute, rule.Name, rule.Parameters, ctx)
	}

	return passed
}

func (v *Validator) lookupRule(name string) rules.RuleFunc {

	if fn, ok := v.extensions[name]; ok {
		return fn
	}

	fn, _ := rules.Lookup(name)

	return fn
}

func (v *Validator) isImplicitRule(name string) bool {
	if v.implicitExt[name] {
		return true
	}

	return rules.IsImplicit(name)
}

func (v *Validator) addFailureWithContext(attribute, ruleName string, params []string, ctx *ruleContext) {
	msg := ctx.customMessage

	if msg == "" {
		msg = v.makeMessage(attribute, ruleName, params)
	}

	v.errs.Add(attribute, msg)
	v.failedRules[attribute] = append(v.failedRules[attribute], ruleName)
}

func (v *Validator) makeMessage(attribute, ruleName string, params []string) string {

	key1 := attribute + "." + strings.ToLower(ruleName)

	if msg, ok := v.customMessages[key1]; ok {
		return v.replacePlaceholders(msg, attribute, ruleName, params)
	}

	key2 := strings.ToLower(ruleName)

	if msg, ok := v.customMessages[key2]; ok {
		return v.replacePlaceholders(msg, attribute, ruleName, params)
	}

	value := v.getValue(attribute)
	typeKey := rules.MessageTypeForSize(ruleName, value)

	if typeKey != ruleName {
		if msg, ok := rules.DefaultMessages[typeKey]; ok && msg != "" {
			return v.replacePlaceholders(msg, attribute, ruleName, params)
		}
	}

	if msg, ok := rules.DefaultMessages[ruleName]; ok && msg != "" {
		return v.replacePlaceholders(msg, attribute, ruleName, params)
	}

	return "The " + v.displayAttribute(attribute) + " field is invalid."
}

func (v *Validator) replacePlaceholders(msg, attribute, ruleName string, params []string) string {
	displayAttr := v.displayAttribute(attribute)

	msg = strings.ReplaceAll(msg, ":attribute", displayAttr)
	msg = strings.ReplaceAll(msg, ":Attribute", titleCase(displayAttr))
	msg = strings.ReplaceAll(msg, ":ATTRIBUTE", strings.ToUpper(displayAttr))

	value := v.getValue(attribute)
	msg = strings.ReplaceAll(msg, ":input", stringify(value))

	paramPlaceholders := []string{
		":min", ":max", ":size", ":other", ":value", ":values",
		":date", ":format", ":digits", ":decimal",
	}
	paramNames := []string{"min", "max", "size", "other", "value", "values", "date", "format", "digits", "decimal"}

	for i, ph := range paramPlaceholders {
		if i < len(params) {
			msg = strings.ReplaceAll(msg, ph, params[i])
		} else if i < len(paramNames) {
			_ = paramNames[i]
		}
	}

	if len(params) > 0 {
		msg = strings.ReplaceAll(msg, ":min", params[0])
		msg = strings.ReplaceAll(msg, ":max", params[0])
		msg = strings.ReplaceAll(msg, ":size", params[0])
		msg = strings.ReplaceAll(msg, ":date", params[0])
		msg = strings.ReplaceAll(msg, ":format", params[0])
		msg = strings.ReplaceAll(msg, ":digits", params[0])
		msg = strings.ReplaceAll(msg, ":decimal", params[0])
		msg = strings.ReplaceAll(msg, ":value", params[0])
		msg = strings.ReplaceAll(msg, ":other", params[0])
	}

	if len(params) > 1 {
		msg = strings.ReplaceAll(msg, ":max", params[1])
	}

	if len(params) > 0 {
		msg = strings.ReplaceAll(msg, ":values", strings.Join(params, ", "))
	}

	return msg
}

func (v *Validator) displayAttribute(attribute string) string {
	if name, ok := v.attrNames[attribute]; ok {
		return name
	}

	s := strings.ReplaceAll(attribute, "_", " ")
	s = strings.ReplaceAll(s, ".", " ")

	return s
}

func (v *Validator) getValue(attribute string) any {
	if val, ok := v.flatData[attribute]; ok {
		return val
	}

	return dotGet(v.data, attribute)
}

func (v *Validator) flatDataHas(attribute string) bool {
	_, ok := v.flatData[attribute]

	return ok
}

func dotGet(data map[string]any, key string) any {
	if val, ok := data[key]; ok {
		return val
	}

	parts := strings.SplitN(key, ".", 2)

	if len(parts) == 1 {
		return nil
	}

	next, ok := data[parts[0]]

	if !ok {
		return nil
	}

	switch child := next.(type) {
	case map[string]any:
		return dotGet(child, parts[1])
	}

	return nil
}

func (v *Validator) parseAllRules() map[string][]ParsedRule {
	out := make(map[string][]ParsedRule, len(v.rawRules))

	for attr, raw := range v.rawRules {
		switch r := raw.(type) {
		case []any:
			out[attr] = Explode(r)
		case string:
			out[attr] = Explode(r)
		case []string:
			out[attr] = Explode(r)
		case []ParsedRule:
			out[attr] = r
		default:
			out[attr] = Explode(raw)
		}
	}

	return out
}

func (v *Validator) buildValidated() {
	result := make(map[string]any)

	for attr := range v.parsedRules {
		if v.excludedAttrs[attr] {
			continue
		}

		val := v.getValue(attr)

		if val != nil || v.flatDataHas(attr) {
			result[attr] = val
		}
	}

	v.validated = result
}

func (v *Validator) makeContext(attribute string) *ruleContext {
	return &ruleContext{validator: v, attribute: attribute}
}

func (c *ruleContext) GetValue(attribute string) any {
	return c.validator.getValue(attribute)
}

func (c *ruleContext) GetData() map[string]any {
	return c.validator.flatData
}

func (c *ruleContext) GetOriginalData() map[string]any {
	return c.validator.data
}

func (c *ruleContext) IsSometimes(attribute string) bool {
	pr, ok := c.validator.parsedRules[attribute]

	if !ok {
		return false
	}

	for _, r := range pr {
		if r.Name == "Sometimes" {
			return true
		}
	}

	return false
}

func (c *ruleContext) IsPresent(attribute string) bool {
	return c.validator.flatDataHas(attribute)
}

func (c *ruleContext) SetMessage(msg string) {
	c.customMessage = msg
}

func (c *ruleContext) GetPresenceVerifier() rules.PresenceVerifier {
	return c.validator.presenceVerifier
}

// Helpers

func isExcludeRule(name string) bool {
	switch name {
	case "Exclude", "ExcludeIf", "ExcludeUnless", "ExcludeWith", "ExcludeWithout":
		return true
	}

	return false
}

func ruleObjectName(rule ParsedRule) string {
	if rule.Object == nil {
		return rule.Name
	}

	return rule.String()
}

func titleCase(s string) string {
	if s == "" {
		return s
	}

	runes := []rune(s)
	capitalize := true

	for i, r := range runes {
		if capitalize && r >= 'a' && r <= 'z' {
			runes[i] = r - 32
			capitalize = false
		} else if r == ' ' {
			capitalize = true
		}
	}

	return string(runes)
}

func stringify(v any) string {
	return rules.StringifyValue(v)
}
