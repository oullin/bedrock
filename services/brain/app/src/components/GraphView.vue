<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { select } from "d3-selection";
import { zoom, zoomIdentity, type D3ZoomEvent } from "d3-zoom";
import type { GraphData, GraphNode } from "@/types/graph";
import { layoutGraph, type LayoutResult } from "@/lib/layout";
import { nodeColor } from "@/lib/colors";

const props = defineProps<{ graph: GraphData }>();
const emit = defineEmits<{ (e: "select", node: GraphNode): void }>();

const layout = computed<LayoutResult>(() => layoutGraph(props.graph));

// d3 zoom for pan/zoom around the svg
const svgEl = ref<SVGSVGElement | null>(null);
const innerEl = ref<SVGGElement | null>(null);
const transform = ref("translate(0,0) scale(1)");

function path(points: { x: number; y: number }[]): string {
  if (points.length === 0) return "";
  return points
    .map((p, i) => `${i === 0 ? "M" : "L"} ${p.x.toFixed(1)} ${p.y.toFixed(1)}`)
    .join(" ");
}

function onZoom(ev: D3ZoomEvent<SVGSVGElement, unknown>) {
  transform.value = `translate(${ev.transform.x},${ev.transform.y}) scale(${ev.transform.k})`;
}

function attachZoom() {
  if (!svgEl.value) return;
  const z = zoom<SVGSVGElement, unknown>().scaleExtent([0.2, 4]).on("zoom", onZoom);
  select(svgEl.value).call(z).call(z.transform, zoomIdentity);
}

onMounted(attachZoom);
watch(() => props.graph, () => {
  transform.value = "translate(0,0) scale(1)";
});
</script>

<template>
  <svg
    ref="svgEl"
    class="graph-svg"
    :viewBox="`0 0 ${Math.max(layout.width, 600)} ${Math.max(layout.height, 400)}`"
    preserveAspectRatio="xMidYMid meet"
  >
    <g ref="innerEl" :transform="transform">
      <g class="edges">
        <path
          v-for="e in layout.edges"
          :key="e.id"
          :d="path(e.points)"
          fill="none"
          stroke="rgba(125, 133, 144, 0.6)"
          stroke-width="1.4"
        />
      </g>
      <g class="nodes">
        <g
          v-for="n in layout.nodes"
          :key="n.id"
          :transform="`translate(${n.x - n.width / 2}, ${n.y - n.height / 2})`"
          class="node"
          @click="emit('select', n)"
        >
          <rect
            :width="n.width"
            :height="n.height"
            rx="6"
            ry="6"
            :fill="nodeColor(n.type)"
            stroke="rgba(0,0,0,0.4)"
            stroke-width="1"
          />
          <text :x="n.width / 2" :y="n.height / 2 + 4" text-anchor="middle" fill="#0e1116" font-size="12" font-weight="600">
            {{ n.label }}
          </text>
        </g>
      </g>
    </g>
  </svg>
</template>

<style scoped>
.graph-svg {
  width: 100%;
  height: 100%;
  background: #0e1116;
  cursor: grab;
}
.graph-svg:active { cursor: grabbing; }
.node { cursor: pointer; }
.node:hover rect { stroke: #ffffff; stroke-width: 1.5; }
</style>
