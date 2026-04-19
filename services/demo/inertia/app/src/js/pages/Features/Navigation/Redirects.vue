<script setup lang="ts">
import { router, usePage } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Button } from "@/js/components/ui/button";
import AppLayout from "@/js/layouts/AppLayout.vue";
import type { SharedPageProps } from "@/js/types";

withDefaults(defineProps<{ timestamp?: string }>(), {
  timestamp: "",
});

const page = usePage<SharedPageProps>();

const breadcrumbs = [{ title: "Features" }, { title: "Navigation" }, { title: "Redirects" }];

function actionUrl(suffix: string): string {
  const url = new URL(page.url, window.location.origin);
  url.pathname = url.pathname.replace(/\/?$/, "/") + suffix;
  return url.pathname;
}

function redirectBack() {
  router.post(actionUrl("back"));
}

function redirectToRoute() {
  router.post(actionUrl("to-route"));
}

function externalRedirect() {
  router.post(actionUrl("external"));
}
</script>

<template>
  <AppLayout title="Redirects" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Redirects"
        description="Server-side redirects are the primary way to navigate after form submissions or actions. Inertia follows redirects transparently and supports external redirects via window.location."
      />

      <p class="text-muted-foreground text-sm">
        Server timestamp:
        <code class="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">{{ timestamp }}</code>
      </p>

      <div class="grid gap-6 lg:grid-cols-3">
        <FeatureCard
          title="Redirect Back"
          description="The server redirects back to this page with a flash message."
        >
          <div class="flex flex-wrap gap-3">
            <Button @click="redirectBack">Redirect back</Button>
          </div>
          <p class="text-muted-foreground mt-2 text-xs">
            Simulates a "redirect back" pattern commonly used after form submissions.
          </p>
        </FeatureCard>

        <FeatureCard
          title="Named Route"
          description="The server redirects to a specific named route with a flash message."
        >
          <div class="flex flex-wrap gap-3">
            <Button variant="outline" @click="redirectToRoute">Redirect to route</Button>
          </div>
          <p class="text-muted-foreground mt-2 text-xs">
            The server resolves a named route and sends a redirect response.
          </p>
        </FeatureCard>

        <FeatureCard
          title="External Redirect"
          description="The server responds with an external location, causing a full page load."
        >
          <div class="flex flex-wrap gap-3">
            <Button variant="outline" @click="externalRedirect">External redirect</Button>
          </div>
          <p class="text-muted-foreground mt-2 text-xs">
            This triggers
            <code class="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">window.location</code> to
            navigate to an external URL (inertiajs.com).
          </p>
        </FeatureCard>
      </div>
    </div>
  </AppLayout>
</template>
