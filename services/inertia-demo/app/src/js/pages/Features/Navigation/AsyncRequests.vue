<script setup lang="ts">
import { Link, usePage } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import AppLayout from "@/js/layouts/AppLayout.vue";
import { useDemoRoute } from "@/js/lib/routes";
import type { SharedPageProps } from "@/js/types";

withDefaults(defineProps<{ timestamp?: string }>(), {
  timestamp: "",
});

const page = usePage<SharedPageProps>();

const breadcrumbs = [{ title: "Features" }, { title: "Navigation" }, { title: "Async Requests" }];

const asyncSlowRoute = useDemoRoute("features.navigation.async-slow");
const asyncRoute = useDemoRoute("features.navigation.async-requests");
</script>

<template>
  <AppLayout title="Async Requests" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Async Requests"
        description="Inertia requests are asynchronous by default. The loading bar at the top of the page indicates when a request is in flight. Slow endpoints make the loading state more visible."
      />

      <p class="text-muted-foreground text-sm">
        Server timestamp:
        <code class="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">{{ timestamp }}</code>
      </p>

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard
          title="Normal Request"
          description="This link reloads the current page instantly. The loading bar is barely visible."
        >
          <div class="flex flex-wrap gap-3">
            <Link
              :href="asyncRoute.url"
              class="inline-flex h-9 items-center rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground hover:bg-primary/90"
            >
              Fast reload
            </Link>
          </div>
        </FeatureCard>

        <FeatureCard
          title="Slow Request (2s delay)"
          description="This link hits a slow endpoint that sleeps for 2 seconds on the server. Watch the loading bar."
        >
          <div class="flex flex-wrap gap-3">
            <Link
              :href="asyncSlowRoute.url"
              class="inline-flex h-9 items-center rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground hover:bg-primary/90"
            >
              Slow reload (2s)
            </Link>
          </div>
          <p class="text-muted-foreground mt-2 text-xs">
            The loading bar at the top of the page fills while the server processes the request.
            This is the default Inertia progress indicator behavior.
          </p>
        </FeatureCard>

        <FeatureCard title="How It Works" description="Inertia uses XHR requests under the hood.">
          <div class="space-y-2 text-sm">
            <p>
              Every Link click and
              <code class="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">router.visit()</code>
              call fires an asynchronous XHR request. The browser stays responsive during the
              request.
            </p>
            <p>
              If a new request starts while another is in progress, Inertia cancels the previous one
              automatically to prevent stale responses from arriving out of order.
            </p>
          </div>
        </FeatureCard>

        <FeatureCard
          title="Progress Indicator"
          description="The built-in loading bar appears after a configurable delay."
        >
          <div class="space-y-2 text-sm">
            <p>
              By default, the progress bar appears after 250ms. For fast endpoints, users never see
              it. For slow endpoints, it provides visual feedback that something is happening.
            </p>
            <p>
              You can customize the delay, color, and behavior via the Inertia progress
              configuration.
            </p>
          </div>
        </FeatureCard>
      </div>
    </div>
  </AppLayout>
</template>
