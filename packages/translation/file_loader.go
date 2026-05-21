package translation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// FileLoader loads translation messages from JSON files on disk, mirroring
// the upstream @bedrock\Translation\FileLoader.
//
// File layout:
//   - Grouped translations: {path}/{locale}/{group}.json
//   - Flat JSON translations: {path}/{locale}.json  (loaded when group == "*")
//
// Multiple paths are supported; later-added paths override earlier ones on a
// per-key basis (deep merge).  Namespaced translations are loaded from the
// hint directory registered via AddNamespace.
type FileLoader struct {
	paths      []string
	jsonPaths  []string
	namespaces map[string]string
}

var _ Loader = (*FileLoader)(nil)

// NewFileLoader returns a *FileLoader seeded with the given paths.
func NewFileLoader(paths ...string) *FileLoader {
	return &FileLoader{
		paths:      append([]string{}, paths...),
		jsonPaths:  []string{},
		namespaces: make(map[string]string),
	}
}

// Load returns the merged message map for locale/group/namespace.
// A nil namespace is treated as the global namespace ("*").
// Returns a non-nil map; missing files produce an empty map.
// Returns a wrapped ErrMalformedJSON on invalid JSON.
func (l *FileLoader) Load(locale, group string, namespace *string) map[string]any {
	ns := globalNS

	if namespace != nil && *namespace != "" {
		ns = *namespace
	}

	if ns != globalNS {
		return l.loadNamespaced(locale, group, ns)
	}

	result := map[string]any{}

	// Load grouped files from each path.
	if group != globalNS {
		for _, p := range l.paths {
			file := filepath.Join(p, locale, group+".json")
			m, err := loadJSONFile(file)

			if err != nil {
				return nil // signals ErrMalformedJSON to caller
			}

			result = mergeMaps(result, m)
		}
	}

	// Load flat JSON files from regular paths and dedicated json paths.
	if group == globalNS {
		for _, p := range l.paths {
			file := filepath.Join(p, locale+".json")
			m, err := loadJSONFile(file)

			if err != nil {
				return nil
			}

			result = mergeMaps(result, m)
		}

		for _, p := range l.jsonPaths {
			file := filepath.Join(p, locale+".json")
			m, err := loadJSONFile(file)

			if err != nil {
				return nil
			}

			result = mergeMaps(result, m)
		}
	}

	return result
}

func (l *FileLoader) loadNamespaced(locale, group, namespace string) map[string]any {
	hint, ok := l.namespaces[namespace]

	if !ok {
		return map[string]any{}
	}

	// Load base file from hint directory.
	file := filepath.Join(hint, locale, group+".json")
	base, err := loadJSONFile(file)

	if err != nil {
		return nil
	}

	// Apply overrides from each regular path.
	for _, p := range l.paths {
		override := filepath.Join(p, "vendor", namespace, locale, group+".json")
		m, err := loadJSONFile(override)

		if err != nil {
			return nil
		}

		base = mergeMaps(base, m)
	}

	return base
}

// AddPath appends a directory to the grouped-translation search paths.
func (l *FileLoader) AddPath(path string) {
	l.paths = append(l.paths, path)
}

// AddNamespace registers a namespace with its hint directory.
func (l *FileLoader) AddNamespace(namespace, hint string) {
	l.namespaces[namespace] = hint
}

// AddJsonPath appends a directory scanned for flat locale JSON files.
func (l *FileLoader) AddJsonPath(path string) {
	l.jsonPaths = append(l.jsonPaths, path)
}

// Paths returns a copy of the grouped-translation search paths.
func (l *FileLoader) Paths() []string {
	return append([]string{}, l.paths...)
}

// JsonPaths returns a copy of the flat JSON search paths.
func (l *FileLoader) JsonPaths() []string {
	return append([]string{}, l.jsonPaths...)
}

// Namespaces returns a shallow copy of the registered namespace map.
func (l *FileLoader) Namespaces() map[string]string {
	cp := make(map[string]string, len(l.namespaces))

	for k, v := range l.namespaces {
		cp[k] = v
	}

	return cp
}

// loadJSONFile reads and unmarshals a JSON file.
// Returns an empty map when the file does not exist.
// Returns (nil, ErrMalformedJSON) on parse failure.
func loadJSONFile(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		if os.IsNotExist(err) {
			return map[string]any{}, nil
		}

		return map[string]any{}, nil // unreadable → treat as missing
	}

	if len(data) == 0 {
		return map[string]any{}, nil
	}

	var m map[string]any

	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrMalformedJSON, path)
	}

	return m, nil
}

// mergeMaps deep-merges overlay onto base.  For keys present in both, if both
// values are map[string]any the merge recurses; otherwise overlay wins.
func mergeMaps(base, overlay map[string]any) map[string]any {
	result := make(map[string]any, len(base)+len(overlay))

	for k, v := range base {
		result[k] = v
	}

	for k, v := range overlay {
		if bv, ok := result[k]; ok {
			if bMap, ok := bv.(map[string]any); ok {
				if oMap, ok := v.(map[string]any); ok {
					result[k] = mergeMaps(bMap, oMap)

					continue
				}
			}
		}

		result[k] = v
	}

	return result
}
