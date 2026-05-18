package workflow

import "sort"

// Marking tracks active places and token counts.
type Marking struct {
	Places map[string]int
}

func NewMarking(places ...string) Marking {
	m := Marking{Places: make(map[string]int, len(places))}

	for _, place := range places {
		if place == "" {
			continue
		}

		m.Places[place]++
	}

	return m
}

func (m Marking) Clone() Marking {
	cloned := Marking{Places: make(map[string]int, len(m.Places))}

	for place, count := range m.Places {
		cloned.Places[place] = count
	}

	return cloned
}

func (m Marking) Has(place string) bool {
	return m.Places[place] > 0
}

func (m Marking) Tokens(place string) int {
	return m.Places[place]
}

func (m Marking) Add(place string, count int) {
	if count <= 0 {
		return
	}

	if m.Places == nil {
		m.Places = map[string]int{}
	}

	m.Places[place] += count
}

func (m Marking) Remove(place string, count int) {
	if count <= 0 || m.Places == nil {
		return
	}

	next := m.Places[place] - count

	if next <= 0 {
		delete(m.Places, place)

		return
	}

	m.Places[place] = next
}

func (m Marking) ActivePlaces() []string {
	places := make([]string, 0, len(m.Places))

	for place, count := range m.Places {
		if count > 0 {
			places = append(places, place)
		}
	}

	sort.Strings(places)

	return places
}
