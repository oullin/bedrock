package graph

// NodeType enumerates every graph node kind brain understands.
//
// The list mirrors upstream-brain's frontend/src/types/graph.ts so the Vue SPA
// can reuse the same discriminator. Filament-specific values are dropped;
// Inertia-specific values are added in their place.
type NodeType string

// Node is a single graph vertex. Field order and tags match upstream-brain's
// Node.php JSON serialisation exactly so the Vue SPA reads either output.
type Node struct {
	ID    string         `json:"id"`
	Type  NodeType       `json:"type"`
	Label string         `json:"label"`
	Data  map[string]any `json:"data"`
}

const (
	NodeTypeRoute             NodeType = "route"
	NodeTypeMiddleware        NodeType = "middleware"
	NodeTypeController        NodeType = "controller"
	NodeTypeAction            NodeType = "action"
	NodeTypeService           NodeType = "service"
	NodeTypeValidationRequest NodeType = "validation_request"
	NodeTypeModel             NodeType = "model"
	NodeTypeEvent             NodeType = "event"
	NodeTypeJob               NodeType = "job"
	NodeTypeCommand           NodeType = "command"
	NodeTypeChannel           NodeType = "channel"
	NodeTypeSchedule          NodeType = "schedule"
	NodeTypeView              NodeType = "view"
	NodeTypeMail              NodeType = "mail"
	NodeTypeNotification      NodeType = "notification"
	NodeTypeEnum              NodeType = "enum"
	NodeTypeInterface         NodeType = "interface"
	NodeTypeTrait             NodeType = "trait"
	NodeTypeAbstractClass     NodeType = "abstract_class"
	NodeTypeServiceProvider   NodeType = "service_provider"
	NodeTypeFacade            NodeType = "facade"
	NodeTypeLivewireComponent NodeType = "livewire_component"

	NodeTypeInertiaPage   NodeType = "inertia_page"
	NodeTypeInertiaLayout NodeType = "inertia_layout"
	NodeTypeInertiaProp   NodeType = "inertia_prop"
)

// NewNode constructs a node with an empty data bag pre-allocated.
func NewNode(id string, t NodeType, label string) *Node {
	return &Node{ID: id, Type: t, Label: label, Data: map[string]any{}}
}

// Set stores a value on the node's data bag and returns the node for chaining.
func (n *Node) Set(key string, value any) *Node {
	if n.Data == nil {
		n.Data = map[string]any{}
	}

	n.Data[key] = value

	return n
}
