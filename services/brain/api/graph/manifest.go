package graph

import "time"

// TabEntry describes one subgraph chunk emitted by the splitter.
// Field names match upstream-brain's manifest format.
type TabEntry struct {
	ID         string `json:"id"`
	Label      string `json:"label"`
	RouteCount int    `json:"routeCount"`
	NodeCount  int    `json:"nodeCount"`
	EdgeCount  int    `json:"edgeCount"`
	File       string `json:"file"`
	RouteFile  string `json:"routeFile,omitempty"`
	Category   string `json:"category,omitempty"`
	PanelID    string `json:"panelId,omitempty"`
}

// Manifest is the top-level index emitted to .graph-manifest.json. The viewer
// loads this first to know which subgraph files to fetch on demand.
type Manifest struct {
	Project      string     `json:"project"`
	AnalyzedAt   time.Time  `json:"analyzedAt"`
	TotalRoutes  int        `json:"totalRoutes"`
	TotalNodes   int        `json:"totalNodes"`
	TotalEdges   int        `json:"totalEdges"`
	Tabs         []TabEntry `json:"tabs"`
}
