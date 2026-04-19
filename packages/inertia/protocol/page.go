package protocol

// Scroll describes infinite-scroll pagination metadata for a prop.
// It is produced by props.Resolve and serialized into the page's
// scrollProps field.
type Scroll struct {
	PageName     string `json:"pageName"`
	PreviousPage any    `json:"previousPage"`
	NextPage     any    `json:"nextPage"`
	CurrentPage  any    `json:"currentPage"`
	Reset        bool   `json:"reset"`
}

// Once describes client-side reuse metadata for a once prop. It is
// produced by props.Resolve and serialized into the page's onceProps
// field.
type Once struct {
	Prop      string `json:"prop"`
	ExpiresAt *int64 `json:"expiresAt,omitempty"`
}
