package graph

// EdgeType is the relationship kind a directed edge represents.
type EdgeType string

const (
	EdgeTypeUses        EdgeType = "uses"
	EdgeTypeHandlesBy   EdgeType = "handles_by"
	EdgeTypeMiddleware  EdgeType = "middleware"
	EdgeTypeDispatches  EdgeType = "dispatches"
	EdgeTypeListensTo   EdgeType = "listens_to"
	EdgeTypeQueues      EdgeType = "queues"
	EdgeTypeRenders     EdgeType = "renders"
	EdgeTypeValidates   EdgeType = "validates"
	EdgeTypeReturns     EdgeType = "returns"
	EdgeTypeCalls       EdgeType = "calls"
	EdgeTypeQueries     EdgeType = "queries"
	EdgeTypeNotifies    EdgeType = "notifies"
	EdgeTypeBroadcasts  EdgeType = "broadcasts"
	EdgeTypeExtends     EdgeType = "extends"
	EdgeTypeImplements  EdgeType = "implements"
	EdgeTypeBindsToImpl EdgeType = "binds_to_impl"
)

// Edge is a directed relationship between two nodes.
type Edge struct {
	ID     string   `json:"id"`
	Source string   `json:"source"`
	Target string   `json:"target"`
	Label  string   `json:"label"`
	Type   EdgeType `json:"type"`
}
