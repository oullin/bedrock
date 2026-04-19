<script setup lang="ts">
import { ref } from "vue";
import { router, usePage } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Button } from "@/js/components/ui/button";
import { Badge } from "@/js/components/ui/badge";
import AppLayout from "@/js/layouts/AppLayout.vue";
import type { SharedPageProps } from "@/js/types";

interface CallbackLogEntry {
  id: number;
  name: string;
  time: string;
}

const page = usePage<SharedPageProps>();

const breadcrumbs = [{ title: "Features" }, { title: "Events" }, { title: "Visit Callbacks" }];

const callbackLog = ref<CallbackLogEntry[]>([]);
let nextId = 0;

function logCallback(name: string) {
  callbackLog.value.unshift({
    id: nextId++,
    name,
    time: new Date().toLocaleTimeString(),
  });
  if (callbackLog.value.length > 20) {
    callbackLog.value.pop();
  }
}

function sendWithCallbacks() {
  router.post(
    page.url,
    {},
    {
      preserveScroll: true,
      onBefore() {
        logCallback("onBefore");
      },
      onStart() {
        logCallback("onStart");
      },
      onProgress(progress) {
        logCallback("onProgress");
      },
      onSuccess() {
        logCallback("onSuccess");
      },
      onError(errors) {
        logCallback("onError");
      },
      onFinish() {
        logCallback("onFinish");
      },
    },
  );
}

function sendWithError() {
  router.post(
    page.url,
    { trigger_error: true },
    {
      preserveScroll: true,
      onBefore() {
        logCallback("onBefore");
      },
      onStart() {
        logCallback("onStart");
      },
      onSuccess() {
        logCallback("onSuccess");
      },
      onError(errors) {
        logCallback("onError");
      },
      onFinish() {
        logCallback("onFinish");
      },
    },
  );
}

function clearLog() {
  callbackLog.value = [];
}
</script>

<template>
  <AppLayout title="Visit Callbacks" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Visit Callbacks"
        description="Attach per-visit callbacks (onBefore, onStart, onSuccess, onError, onFinish) to individual router calls for fine-grained control."
      />

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard title="Trigger Visits" description="Each visit has its own set of callbacks.">
          <div class="space-y-4">
            <div class="flex flex-wrap gap-2">
              <Button @click="sendWithCallbacks">Successful POST</Button>
              <Button variant="destructive" @click="sendWithError">POST with Error</Button>
              <Button variant="outline" @click="clearLog">Clear Log</Button>
            </div>
            <p class="text-muted-foreground text-sm">
              Unlike global events, visit callbacks are scoped to a single request. They are passed
              as options to
              <code class="bg-muted rounded px-1.5 py-0.5 text-xs">router.post()</code>,
              <code class="bg-muted rounded px-1.5 py-0.5 text-xs">router.visit()</code>, etc.
            </p>
          </div>
        </FeatureCard>

        <FeatureCard title="Callback Log" :description="`${callbackLog.length} callback(s) fired`">
          <div class="max-h-80 space-y-1.5 overflow-auto">
            <div
              v-for="entry in callbackLog"
              :key="entry.id"
              class="flex items-center gap-2 rounded-md border p-2 text-xs"
            >
              <Badge
                variant="outline"
                class="shrink-0"
                :class="{
                  'border-blue-200 text-blue-700': entry.name === 'onBefore',
                  'border-yellow-200 text-yellow-700': entry.name === 'onStart',
                  'border-purple-200 text-purple-700': entry.name === 'onProgress',
                  'border-green-200 text-green-700': entry.name === 'onSuccess',
                  'border-red-200 text-red-700': entry.name === 'onError',
                  'border-gray-200 text-gray-700': entry.name === 'onFinish',
                }"
              >
                {{ entry.name }}
              </Badge>
              <span class="text-muted-foreground flex-1">Visit callback fired</span>
              <span class="text-muted-foreground shrink-0">{{ entry.time }}</span>
            </div>
            <p
              v-if="callbackLog.length === 0"
              class="text-muted-foreground py-4 text-center text-sm"
            >
              No callbacks fired yet. Click a button to trigger a visit.
            </p>
          </div>
        </FeatureCard>
      </div>

      <FeatureCard
        title="Callback Reference"
        description="Available per-visit callbacks and when they fire."
      >
        <div class="overflow-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="border-b">
                <th class="py-2 pr-4 text-left font-medium">Callback</th>
                <th class="py-2 pr-4 text-left font-medium">When</th>
                <th class="py-2 text-left font-medium">Use Case</th>
              </tr>
            </thead>
            <tbody class="text-muted-foreground">
              <tr class="border-b">
                <td class="py-2 pr-4"><Badge variant="outline" class="text-xs">onBefore</Badge></td>
                <td class="py-2 pr-4">Before the request is made</td>
                <td class="py-2">Show confirmation, cancel visit</td>
              </tr>
              <tr class="border-b">
                <td class="py-2 pr-4"><Badge variant="outline" class="text-xs">onStart</Badge></td>
                <td class="py-2 pr-4">Request begins</td>
                <td class="py-2">Show loading state</td>
              </tr>
              <tr class="border-b">
                <td class="py-2 pr-4">
                  <Badge variant="outline" class="text-xs">onProgress</Badge>
                </td>
                <td class="py-2 pr-4">Upload progress updates</td>
                <td class="py-2">Update progress bar</td>
              </tr>
              <tr class="border-b">
                <td class="py-2 pr-4">
                  <Badge variant="outline" class="text-xs">onSuccess</Badge>
                </td>
                <td class="py-2 pr-4">Successful response</td>
                <td class="py-2">Show toast, redirect</td>
              </tr>
              <tr class="border-b">
                <td class="py-2 pr-4"><Badge variant="outline" class="text-xs">onError</Badge></td>
                <td class="py-2 pr-4">Validation errors returned</td>
                <td class="py-2">Focus first error field</td>
              </tr>
              <tr>
                <td class="py-2 pr-4"><Badge variant="outline" class="text-xs">onFinish</Badge></td>
                <td class="py-2 pr-4">After success or error</td>
                <td class="py-2">Hide loading, cleanup</td>
              </tr>
            </tbody>
          </table>
        </div>
      </FeatureCard>
    </div>
  </AppLayout>
</template>
