import type { NodeType } from "@/types/graph";

// Colours match upstream-brain's frontend palette so users moving between
// the two tools see consistent visual cues.
export const NODE_COLORS: Record<NodeType, string> = {
  route: "#2dd4bf",
  middleware: "#fb923c",
  controller: "#60a5fa",
  action: "#7dd3fc",
  service: "#a78bfa",
  validation_request: "#f472b6",
  model: "#f87171",
  event: "#fbbf24",
  job: "#94a3b8",
  command: "#34d399",
  channel: "#22d3ee",
  schedule: "#fde047",
  view: "#a3e635",
  mail: "#f59e0b",
  notification: "#e879f9",
  enum: "#cbd5e1",
  interface: "#bef264",
  trait: "#fcd34d",
  abstract_class: "#fb7185",
  service_provider: "#c084fc",
  facade: "#67e8f9",
  livewire_component: "#f0abfc",
  inertia_page: "#818cf8",
  inertia_layout: "#a5b4fc",
  inertia_prop: "#c7d2fe",
};

export function nodeColor(t: NodeType): string {
  return NODE_COLORS[t] ?? "#94a3b8";
}
