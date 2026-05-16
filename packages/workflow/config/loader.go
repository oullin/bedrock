// Package config builds workflow Definitions from a bedrock config.Repository.
//
// It reuses packages/config rather than wiring its own viper instance, so the
// same dot-notation, env binding, and merge rules apply.
//
// Expected schema (under a configurable root key, default "workflow"):
//
//	workflow:
//	  name: subscription
//	  places: [trial, active, cancelled]
//	  initial: [trial]
//	  transitions:
//	    - name: activate
//	      from: [trial]
//	      to: [active]
//	    - name: cancel
//	      from: [active]
//	      to: [cancelled]
//	  metadata:                       # optional workflow-level metadata
//	    purpose: billing
//	  places_metadata:                # optional per-place metadata
//	    trial:
//	      ttl: 30d
//	  transitions_metadata:           # optional per-transition metadata
//	    activate:
//	      audit_level: high
package config

import (
	"fmt"

	bconfig "github.com/bedrock/packages/config"
	"github.com/bedrock/packages/workflow"
)

// Load builds a workflow Definition from the repository under root key "workflow".

// LoadAt builds a workflow Definition from the repository under `rootKey`.

type transitionDecl struct {
	Name string
	From []string
	To   []string
}

func Load(r *bconfig.Repository) (*workflow.Definition, error) {
	return LoadAt(r, "workflow")
}

func LoadAt(r *bconfig.Repository, rootKey string) (*workflow.Definition, error) {
	if r == nil {
		return nil, fmt.Errorf("config: repository is required")
	}

	if rootKey == "" {
		return nil, fmt.Errorf("config: root key is required")
	}

	prefix := rootKey + "."

	if !r.Has(rootKey) && !r.Has(prefix+"places") {
		return nil, fmt.Errorf("config: no workflow definition at %q", rootKey)
	}

	builder := workflow.NewDefinitionBuilder()

	places := coerceStringSlice(r.Get(prefix + "places"))

	if len(places) == 0 {
		return nil, fmt.Errorf("config: workflow %q requires at least one place", rootKey)
	}

	for _, place := range places {
		builder.AddPlace(place)
	}

	initial := coerceStringSlice(r.Get(prefix + "initial"))

	if len(initial) == 0 {
		return nil, fmt.Errorf("config: workflow %q requires an initial place list", rootKey)
	}

	builder.SetInitialPlaces(initial...)

	transitions, err := coerceTransitions(r.Get(prefix + "transitions"))

	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	for _, t := range transitions {
		builder.AddTransition(t.Name, t.From, t.To)
	}

	for key, value := range coerceStringMap(r.Get(prefix + "metadata")) {
		builder.SetMetadata(key, value)
	}

	for place, raw := range coerceStringMap(r.Get(prefix + "places_metadata")) {
		for key, value := range coerceStringMap(raw) {
			builder.SetPlaceMetadata(place, key, value)
		}
	}

	for transition, raw := range coerceStringMap(r.Get(prefix + "transitions_metadata")) {
		for key, value := range coerceStringMap(raw) {
			builder.SetTransitionMetadata(transition, key, value)
		}
	}

	return builder.Build()
}

// coerceTransitions accepts the most common shapes a YAML/JSON decoder produces
// for a list of objects.
func coerceTransitions(raw any) ([]transitionDecl, error) {
	if raw == nil {
		return nil, nil
	}

	list, ok := raw.([]any)

	if !ok {
		return nil, fmt.Errorf("transitions must be a list, got %T", raw)
	}

	out := make([]transitionDecl, 0, len(list))

	for i, item := range list {
		entry := coerceStringMap(item)

		if entry == nil {
			return nil, fmt.Errorf("transition[%d] must be an object, got %T", i, item)
		}

		decl := transitionDecl{}

		if name, ok := entry["name"].(string); ok {
			decl.Name = name
		}

		decl.From = coerceStringSlice(entry["from"])
		decl.To = coerceStringSlice(entry["to"])

		out = append(out, decl)
	}

	return out, nil
}

func coerceStringSlice(raw any) []string {
	if raw == nil {
		return nil
	}

	switch v := raw.(type) {
	case []string:
		return append([]string(nil), v...)
	case []any:
		out := make([]string, 0, len(v))

		for _, item := range v {
			if s, ok := item.(string); ok {
				out = append(out, s)
			}
		}

		return out
	case string:
		return []string{v}
	}

	return nil
}

func coerceStringMap(raw any) map[string]any {
	switch v := raw.(type) {
	case map[string]any:
		return v
	case map[any]any:
		out := make(map[string]any, len(v))

		for key, value := range v {
			if s, ok := key.(string); ok {
				out[s] = value
			}
		}

		return out
	}

	return nil
}
