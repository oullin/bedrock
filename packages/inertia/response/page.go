package response

import "github.com/bedrock/packages/inertia/protocol"

// Page is the Inertia.js page object sent to the client. On XHR visits
// it is the JSON body; on initial visits it is embedded in a
// <script type="application/json"> element in the root HTML template.
type Page struct {
	Component      string                     `json:"component"`
	Props          map[string]any             `json:"props"`
	URL            string                     `json:"url"`
	Version        string                     `json:"version"`
	EncryptHistory bool                       `json:"encryptHistory,omitempty"`
	ClearHistory   bool                       `json:"clearHistory,omitempty"`
	MergeProps     []string                   `json:"mergeProps,omitempty"`
	DeepMergeProps []string                   `json:"deepMergeProps,omitempty"`
	DeferredProps  map[string][]string        `json:"deferredProps,omitempty"`
	ScrollProps    map[string]protocol.Scroll `json:"scrollProps,omitempty"`
	OnceProps      map[string]protocol.Once   `json:"onceProps,omitempty"`
}
