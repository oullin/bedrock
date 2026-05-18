// Ported verbatim from upstream-brain/frontend/src/types/graph.ts so the
// Go-side graph JSON renders unchanged. Filament types are dropped;
// Inertia types are appended.

export type NodeType =
  | "route"
  | "middleware"
  | "controller"
  | "action"
  | "service"
  | "validation_request"
  | "model"
  | "event"
  | "job"
  | "command"
  | "channel"
  | "schedule"
  | "view"
  | "mail"
  | "notification"
  | "enum"
  | "interface"
  | "trait"
  | "abstract_class"
  | "service_provider"
  | "facade"
  | "livewire_component"
  | "inertia_page"
  | "inertia_layout"
  | "inertia_prop";

export type EdgeType =
  | "uses"
  | "handles_by"
  | "middleware"
  | "dispatches"
  | "listens_to"
  | "queues"
  | "renders"
  | "validates"
  | "returns"
  | "calls"
  | "queries"
  | "notifies"
  | "broadcasts"
  | "extends"
  | "implements"
  | "binds_to_impl";

export interface GraphNode {
  id: string;
  type: NodeType;
  label: string;
  data: Record<string, unknown>;
}

export interface GraphEdge {
  id: string;
  source: string;
  target: string;
  label: string;
  type: EdgeType;
}

export interface GraphMeta {
  project: string;
  analyzedAt: string;
  nodeCount: number;
  edgeCount: number;
}

export interface GraphData {
  meta: GraphMeta;
  nodes: GraphNode[];
  edges: GraphEdge[];
}

export interface TabEntry {
  id: string;
  label: string;
  routeCount: number;
  nodeCount: number;
  edgeCount: number;
  file: string;
  routeFile?: string;
  category?: string;
  panelId?: string;
}

export interface Manifest {
  project: string;
  analyzedAt: string;
  totalRoutes: number;
  totalNodes: number;
  totalEdges: number;
  tabs: TabEntry[];
}
