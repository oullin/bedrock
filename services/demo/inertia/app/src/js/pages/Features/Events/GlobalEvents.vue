<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import { router, usePage } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Button } from "@/js/components/ui/button";
import { Badge } from "@/js/components/ui/badge";
import AppLayout from "@/js/layouts/AppLayout.vue";
import type { SharedPageProps } from "@/js/types";

interface LogEntry {
  id: number;
  name: string;
  detail: string;
  time: string;
}

const page = usePage<SharedPageProps>();

const breadcrumbs = [{ title: "Features" }, { title: "Events" }, { title: "Global Events" }];

let idCounter = 0;
const eventLog = ref<LogEntry[]>([]);
const removers: Array<() => void> = [];

function logEvent(name: string, detail: string | object) {
  eventLog.value.unshift({
    id: ++idCounter,
    name,
    detail:
      detail !== null && typeof detail === "object" ? JSON.stringify(detail) : String(detail ?? ""),
    time: new Date().toLocaleTimeString(),
  });
  if (eventLog.value.length > 20) {
    eventLog.value.pop();
  }
}

onMounted(() => {
  removers.push(
    router.on("before", (event) => {
      logEvent("before", { url: event.detail?.visit?.url });
    }),
    router.on("start", () => {
      logEvent("start", "Visit started");
    }),
    router.on("progress", (event) => {
      logEvent("progress", { percentage: event.detail?.progress?.percentage });
    }),
    router.on("success", () => {
      logEvent("success", "Visit completed successfully");
    }),
    router.on("finish", () => {
      logEvent("finish", "Visit finished");
    }),
  );
});

onUnmounted(() => {
  removers.forEach((remove) => remove());
});

function triggerPost() {
  router.post(page.url, {}, { preserveScroll: true });
}

function clearLog() {
  eventLog.value = [];
}
</script>

<template>
  <AppLayout title="Global Events" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Global Events"
        description="Listen to Inertia router events globally using router.on(). Events fire during the lifecycle of every Inertia visit."
      />

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard title="Trigger Events" description="Make a request to see events fire.">
          <div class="space-y-4">
            <div class="flex gap-2">
              <Button @click="triggerPost">Send POST Request</Button>
              <Button variant="outline" @click="clearLog">Clear Log</Button>
            </div>
            <div class="space-y-2 text-sm">
              <p class="font-medium">Registered listeners:</p>
              <div class="flex flex-wrap gap-1">
                <Badge variant="secondary">before</Badge>
                <Badge variant="secondary">start</Badge>
                <Badge variant="secondary">progress</Badge>
                <Badge variant="secondary">success</Badge>
                <Badge variant="secondary">finish</Badge>
              </div>
            </div>
          </div>
        </FeatureCard>

        <FeatureCard title="Event Log" :description="`${eventLog.length} event(s) captured`">
          <div class="max-h-80 space-y-1.5 overflow-auto">
            <div
              v-for="entry in eventLog"
              :key="entry.id"
              class="flex items-start gap-2 rounded-md border p-2 text-xs"
            >
              <Badge
                variant="outline"
                class="shrink-0"
                :class="{
                  'border-blue-200 text-blue-700': entry.name === 'before',
                  'border-yellow-200 text-yellow-700': entry.name === 'start',
                  'border-purple-200 text-purple-700': entry.name === 'progress',
                  'border-green-200 text-green-700': entry.name === 'success',
                  'border-gray-200 text-gray-700': entry.name === 'finish',
                }"
              >
                {{ entry.name }}
              </Badge>
              <span class="text-muted-foreground flex-1 truncate">{{ entry.detail }}</span>
              <span class="text-muted-foreground shrink-0">{{ entry.time }}</span>
            </div>
            <p v-if="eventLog.length === 0" class="text-muted-foreground py-4 text-center text-sm">
              No events captured yet. Click "Send POST Request" to trigger some.
            </p>
          </div>
        </FeatureCard>
      </div>

      <FeatureCard title="Event Lifecycle" description="The order in which events fire.">
        <div class="flex flex-wrap items-center gap-2 text-sm">
          <Badge>before</Badge>
          <span class="text-muted-foreground">-></span>
          <Badge variant="secondary">start</Badge>
          <span class="text-muted-foreground">-></span>
          <Badge variant="secondary">progress</Badge>
          <span class="text-muted-foreground">-></span>
          <Badge variant="secondary">success / error</Badge>
          <span class="text-muted-foreground">-></span>
          <Badge variant="outline">finish</Badge>
        </div>
      </FeatureCard>
    </div>
  </AppLayout>
</template>
