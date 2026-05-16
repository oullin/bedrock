import dagre from "dagre";
import type { GraphData, GraphEdge, GraphNode } from "@/types/graph";

export interface LaidOutNode extends GraphNode {
  x: number;
  y: number;
  width: number;
  height: number;
}

export interface LaidOutEdge extends GraphEdge {
  points: { x: number; y: number }[];
}

export interface LayoutResult {
  width: number;
  height: number;
  nodes: LaidOutNode[];
  edges: LaidOutEdge[];
}

const NODE_W = 160;
const NODE_H = 44;

// Runs Dagre to assign (x,y) to every node + a polyline path to every
// edge. Pure function — caller passes the data, gets back a layout.
export function layoutGraph(data: GraphData): LayoutResult {
  const g = new dagre.graphlib.Graph();
  g.setGraph({ rankdir: "LR", marginx: 32, marginy: 32, nodesep: 24, ranksep: 64 });
  g.setDefaultEdgeLabel(() => ({}));

  for (const n of data.nodes) {
    g.setNode(n.id, { width: NODE_W, height: NODE_H });
  }
  for (const e of data.edges) {
    if (!g.hasNode(e.source) || !g.hasNode(e.target)) continue;
    g.setEdge(e.source, e.target);
  }
  dagre.layout(g);

  const nodes: LaidOutNode[] = data.nodes.map((n) => {
    const meta = g.node(n.id) ?? { x: 0, y: 0 };
    return { ...n, x: meta.x, y: meta.y, width: NODE_W, height: NODE_H };
  });
  const edges: LaidOutEdge[] = data.edges.map((e) => {
    const meta = g.edge({ v: e.source, w: e.target }) as { points?: { x: number; y: number }[] } | undefined;
    return { ...e, points: meta?.points ?? [] };
  });
  const { width = 0, height = 0 } = g.graph();
  return { width, height, nodes, edges };
}
