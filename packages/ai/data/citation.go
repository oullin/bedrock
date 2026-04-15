package data

// Citation represents a source citation in a response.
// Mirrors Laravel\Ai\Responses\Data\Citation.
type Citation struct {
	Title  string `json:"title"`
	URL    string `json:"url"`
	Source string `json:"source"`
}

// ToMap returns a map representation.
func (c Citation) ToMap() map[string]any {
	return map[string]any{"title": c.Title, "url": c.URL, "source": c.Source}
}

// UrlCitation is a URL-based citation variant.
// Mirrors Laravel\Ai\Responses\Data\UrlCitation.
type UrlCitation struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	StartIndex  int    `json:"start_index"`
	EndIndex    int    `json:"end_index"`
}

// ToMap returns a map representation.
func (u UrlCitation) ToMap() map[string]any {
	return map[string]any{
		"title":       u.Title,
		"url":         u.URL,
		"start_index": u.StartIndex,
		"end_index":   u.EndIndex,
	}
}
