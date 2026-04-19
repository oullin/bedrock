<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import { router, usePage } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Button } from "@/js/components/ui/button";
import { Badge } from "@/js/components/ui/badge";
import AppLayout from "@/js/layouts/AppLayout.vue";
import type { SharedPageProps } from "@/js/types";

const page = usePage<SharedPageProps>();

const breadcrumbs = [{ title: "Features" }, { title: "Errors" }, { title: "HTTP Exceptions" }];

const lastException = ref<{ status: number; time: string } | null>(null);
let removeListener: (() => void) | undefined;

onMounted(() => {
  // "invalid" event exists at runtime but is not in Inertia's GlobalEventsMap types
  removeListener = (router as any).on("invalid", (event: any) => {
    lastException.value = {
      status: event.detail?.response?.status,
      time: new Date().toLocaleTimeString(),
    };
    event.preventDefault();
  });
});

onUnmounted(() => {
  if (removeListener) removeListener();
});

function trigger403() {
  router.post(page.url, { status: 403 }, { preserveScroll: true });
}

function trigger404() {
  router.post(page.url, { status: 404 }, { preserveScroll: true });
}

function trigger500() {
  router.post(page.url, { status: 500 }, { preserveScroll: true });
}
</script>

<template>
  <AppLayout title="HTTP Exceptions" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="HTTP Exceptions"
        description="Handle server-side HTTP errors (403, 404, 500, etc.) gracefully in your Inertia application."
      />

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard
          title="Trigger Errors"
          description="Click to trigger different HTTP error responses."
        >
          <div class="space-y-4">
            <div class="flex flex-wrap gap-2">
              <Button variant="outline" @click="trigger403"> 403 Forbidden </Button>
              <Button variant="outline" @click="trigger404"> 404 Not Found </Button>
              <Button variant="destructive" @click="trigger500"> 500 Server Error </Button>
            </div>
            <p class="text-muted-foreground text-sm">
              These buttons send POST requests that trigger the server to respond with the
              corresponding HTTP status code.
            </p>
          </div>
        </FeatureCard>

        <FeatureCard
          title="Last Exception"
          description="The most recent error caught by the listener."
        >
          <div v-if="lastException" class="space-y-2">
            <div class="flex items-center gap-2">
              <Badge variant="destructive">{{ lastException.status }}</Badge>
              <span class="text-muted-foreground text-sm">at {{ lastException.time }}</span>
            </div>
            <p class="text-sm">
              This error was intercepted by the
              <code class="bg-muted rounded px-1.5 py-0.5 text-xs">router.on('invalid')</code>
              listener.
            </p>
          </div>
          <p v-else class="text-muted-foreground text-sm">
            No exceptions caught yet. Trigger an error using the buttons.
          </p>
        </FeatureCard>
      </div>

      <FeatureCard
        title="Error Handling Strategies"
        description="Different ways to handle HTTP errors."
      >
        <div class="space-y-4">
          <div class="rounded-md border p-3">
            <p class="text-sm font-medium">Error Pages</p>
            <p class="text-muted-foreground mt-1 text-sm">
              By default, Inertia renders error pages for non-Inertia responses. You can customize
              these by creating error page components.
            </p>
          </div>
          <div class="rounded-md border p-3">
            <p class="text-sm font-medium">Global Event Listener</p>
            <pre class="bg-muted mt-2 overflow-auto rounded p-3 text-xs leading-relaxed">
router.on('invalid', (event) => {
  const status = event.detail.response.status
  // Handle the error
  event.preventDefault()
})</pre
            >
          </div>
          <div class="rounded-md border p-3">
            <p class="text-sm font-medium">Per-Visit Error Handling</p>
            <pre class="bg-muted mt-2 overflow-auto rounded p-3 text-xs leading-relaxed">
router.post(url, data, {
  onError(errors) {
    // Handle validation errors
  },
})</pre
            >
            <p class="text-muted-foreground mt-1 text-xs">
              Note: onError handles validation errors (422), not HTTP exceptions.
            </p>
          </div>
        </div>
      </FeatureCard>
    </div>
  </AppLayout>
</template>
