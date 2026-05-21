package data

// RankedDocument holds a single reranked document with its relevance score.
type RankedDocument struct {
	Index    int     `json:"index"`
	Document string  `json:"document"`
	Score    float64 `json:"score"`
}

// ToMap returns a map representation.
func (r RankedDocument) ToMap() map[string]any {
	return map[string]any{
		"index":    r.Index,
		"document": r.Document,
		"score":    r.Score,
	}
}
