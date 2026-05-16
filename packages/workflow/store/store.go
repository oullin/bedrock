// Package store provides marking store implementations: SingleState for
// state-machine subjects (one active place) and MultiState for general
// Petri-Net subjects (any number of active places).
package store

import (
	"fmt"

	"github.com/bedrock/packages/workflow"
)

// SingleState persists exactly one active place per subject — the classic
// state-machine shape. Use this for entities with a single string `State`
// field.
type SingleState[T any] struct {
	Getter func(T) string
	Setter func(T, string)
}

// MultiState persists any number of active places — used for general Petri-Net
// workflows where multiple tokens may be active concurrently.
type MultiState[T any] struct {
	Getter func(T) []string
	Setter func(T, []string)
}

func (s *SingleState[T]) GetMarking(subject T, _ *workflow.Definition) (workflow.Marking, error) {
	if s == nil || s.Getter == nil {
		return workflow.Marking{}, fmt.Errorf("single state getter is required")
	}

	place := s.Getter(subject)

	if place == "" {
		return workflow.Marking{}, nil
	}

	return workflow.NewMarking(place), nil
}

func (s *SingleState[T]) SetMarking(subject T, marking workflow.Marking, _ *workflow.Definition, _ map[string]any) error {
	if s == nil || s.Setter == nil {
		return fmt.Errorf("single state setter is required")
	}

	places := marking.ActivePlaces()

	if len(places) > 1 {
		return fmt.Errorf("single state store cannot persist %d active places", len(places))
	}

	place := ""

	if len(places) == 1 {
		place = places[0]
	}

	s.Setter(subject, place)

	return nil
}

func (s *MultiState[T]) GetMarking(subject T, _ *workflow.Definition) (workflow.Marking, error) {
	if s == nil || s.Getter == nil {
		return workflow.Marking{}, fmt.Errorf("multi state getter is required")
	}

	return workflow.NewMarking(s.Getter(subject)...), nil
}

func (s *MultiState[T]) SetMarking(subject T, marking workflow.Marking, _ *workflow.Definition, _ map[string]any) error {
	if s == nil || s.Setter == nil {
		return fmt.Errorf("multi state setter is required")
	}

	s.Setter(subject, marking.ActivePlaces())

	return nil
}
