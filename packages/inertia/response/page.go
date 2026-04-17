package response

// Page is the Inertia.js page object sent to the client. On XHR visits
// it is the JSON body; on initial visits it is embedded in a
// <script type="application/json"> element in the root HTML template.
type Page struct {
	Component      string              `json:"component"`
	Props          map[string]any      `json:"props"`
	URL            string              `json:"url"`
	Version        string              `json:"version"`
	EncryptHistory bool                `json:"encryptHistory,omitempty"`
	ClearHistory   bool                `json:"clearHistory,omitempty"`
	MergeProps     []string            `json:"mergeProps,omitempty"`
	DeepMergeProps []string            `json:"deepMergeProps,omitempty"`
	DeferredProps  map[string][]string `json:"deferredProps,omitempty"`
	ScrollProps    map[string]Scroll   `json:"scrollProps,omitempty"`
	OnceProps      map[string]Once     `json:"onceProps,omitempty"`
}

// Scroll describes infinite-scroll pagination metadata for a prop.
type Scroll struct {
	PageName     string `json:"pageName"`
	PreviousPage any    `json:"previousPage"`
	NextPage     any    `json:"nextPage"`
	CurrentPage  any    `json:"currentPage"`
	Reset        bool   `json:"reset"`
}

// Once describes client-side reuse metadata for a once prop.
type Once struct {
	Prop      string `json:"prop"`
	ExpiresAt *int64 `json:"expiresAt,omitempty"`
}
