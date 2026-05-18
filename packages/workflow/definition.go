package workflow

import (
	"fmt"
	"slices"
)

// Definition describes the static structure of a Petri-Net workflow.
type Definition struct {
	Places             []string
	Transitions        []Transition
	InitialMarking     Marking
	Metadata           map[string]any
	PlaceMetadata      map[string]map[string]any
	TransitionMetadata map[string]map[string]any

	placesMap map[string]struct{}
}

func (d *Definition) Clone() *Definition {
	if d == nil {
		return nil
	}

	out := &Definition{
		Places:             append([]string(nil), d.Places...),
		Transitions:        make([]Transition, len(d.Transitions)),
		InitialMarking:     d.InitialMarking.Clone(),
		Metadata:           cloneMap(d.Metadata),
		PlaceMetadata:      cloneNestedMap(d.PlaceMetadata),
		TransitionMetadata: cloneNestedMap(d.TransitionMetadata),
		placesMap:          make(map[string]struct{}, len(d.Places)),
	}

	for _, place := range d.Places {
		out.placesMap[place] = struct{}{}
	}

	for i, transition := range d.Transitions {
		out.Transitions[i] = Transition{
			Name:     transition.Name,
			From:     append([]string(nil), transition.From...),
			To:       append([]string(nil), transition.To...),
			Metadata: cloneMap(transition.Metadata),
		}
	}

	return out
}

func (d *Definition) Transition(name string) (Transition, bool) {
	for _, transition := range d.Transitions {
		if transition.Name == name {
			if metadata, ok := d.TransitionMetadata[name]; ok && len(transition.Metadata) == 0 {
				transition.Metadata = cloneMap(metadata)
			}

			return transition, true
		}
	}

	return Transition{}, false
}

func (d *Definition) HasPlace(place string) bool {
	if d.placesMap == nil {
		d.placesMap = make(map[string]struct{}, len(d.Places))

		for _, p := range d.Places {
			d.placesMap[p] = struct{}{}
		}
	}

	_, ok := d.placesMap[place]

	return ok
}

func (d *Definition) MetadataValue(key string) (any, bool) {
	value, ok := d.Metadata[key]

	return value, ok
}

func (d *Definition) PlaceMetadataValue(place string, key string) (any, bool) {
	values, ok := d.PlaceMetadata[place]

	if !ok {
		return nil, false
	}

	value, ok := values[key]

	return value, ok
}

func (d *Definition) TransitionMetadataValue(transition string, key string) (any, bool) {
	values, ok := d.TransitionMetadata[transition]

	if !ok {
		return nil, false
	}

	value, ok := values[key]

	return value, ok
}

func (d *Definition) Validate() error {
	if len(d.Places) == 0 {
		return fmt.Errorf("definition requires at least one place")
	}

	for _, transition := range d.Transitions {
		if transition.Name == "" {
			return fmt.Errorf("transition name cannot be empty")
		}

		if len(transition.From) == 0 {
			return fmt.Errorf("transition %q requires at least one from place", transition.Name)
		}

		if len(transition.To) == 0 {
			return fmt.Errorf("transition %q requires at least one to place", transition.Name)
		}

		for _, place := range append(append([]string(nil), transition.From...), transition.To...) {
			if !d.HasPlace(place) {
				return fmt.Errorf("transition %q references unknown place %q", transition.Name, place)
			}
		}
	}

	if len(d.InitialMarking.Places) == 0 {
		return fmt.Errorf("definition requires an initial marking")
	}

	for place := range d.InitialMarking.Places {
		if !d.HasPlace(place) {
			return fmt.Errorf("initial marking references unknown place %q", place)
		}
	}

	return d.validateReachability()
}

func (d *Definition) validateReachability() error {
	reachable := make(map[string]struct{})
	queue := []string{}

	for place := range d.InitialMarking.Places {
		reachable[place] = struct{}{}
		queue = append(queue, place)
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, transition := range d.Transitions {
			if slices.Contains(transition.From, current) {
				for _, to := range transition.To {
					if _, ok := reachable[to]; !ok {
						reachable[to] = struct{}{}
						queue = append(queue, to)
					}
				}
			}
		}
	}

	for _, place := range d.Places {
		if _, ok := reachable[place]; !ok {
			return fmt.Errorf("place %q is unreachable from the initial marking", place)
		}
	}

	return nil
}

func cloneMap(src map[string]any) map[string]any {
	if len(src) == 0 {
		return nil
	}

	out := make(map[string]any, len(src))

	for key, value := range src {
		out[key] = value
	}

	return out
}

func cloneNestedMap(src map[string]map[string]any) map[string]map[string]any {
	if len(src) == 0 {
		return nil
	}

	out := make(map[string]map[string]any, len(src))

	for key, value := range src {
		out[key] = cloneMap(value)
	}

	return out
}
