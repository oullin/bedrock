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

const breadcrumbs = [
  { title: "Features" },
  { title: "Navigation" },
  { title: "History Management" },
];

function submitAction() {
  router.post(page.url);
}

function replaceVisit() {
  router.get(page.url, {}, { replace: true });
}
</script>

<template>
  <AppLayout title="History Management" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="History Management"
        description="Inertia lets you control how visits interact with the browser history stack. POST requests automatically redirect with a flash message, and you can use the replace option to avoid adding entries."
      />

      <p class="text-muted-foreground text-sm">
        Server timestamp:
        <code class="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">{{ timestamp }}</code>
      </p>

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard
          title="POST with Redirect"
          description="Submit a POST request. The server redirects back with a flash message, which is the standard POST/redirect/GET pattern."
        >
          <div class="flex flex-wrap gap-3">
            <Button @click="submitAction">Submit POST</Button>
          </div>
          <p class="text-muted-foreground mt-2 text-xs">
            The server sets a flash message and redirects back. The flash banner appears at the top
            of the page.
          </p>
        </FeatureCard>

        <FeatureCard
          title="Replace History"
          description="Visit with replace: true to swap the current history entry instead of adding a new one."
        >
          <div class="flex flex-wrap gap-3">
            <Button variant="outline" @click="replaceVisit">Reload (replace entry)</Button>
          </div>
          <p class="text-muted-foreground mt-2 text-xs">
            After clicking, pressing back in the browser will skip this page in the history stack.
          </p>
        </FeatureCard>

        <FeatureCard
          title="How It Works"
          description="Inertia uses pushState and replaceState under the hood."
        >
          <div class="space-y-2 text-sm">
            <p>
              By default, each visit calls
              <code class="rounded bg-muted px-1.5 py-0.5 font-mono text-xs"
                >history.pushState()</code
              >, adding a new entry to the browser history.
            </p>
            <p>
              When you set
              <code class="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">replace: true</code>,
              Inertia uses
              <code class="rounded bg-muted px-1.5 py-0.5 font-mono text-xs"
                >history.replaceState()</code
              >
              instead, overwriting the current entry.
            </p>
          </div>
        </FeatureCard>

        <FeatureCard
          title="POST/Redirect/GET"
          description="Inertia follows the PRG pattern to prevent form resubmission."
        >
          <div class="space-y-2 text-sm">
            <p>
              When you submit a POST request, the server processes it and responds with a redirect.
              Inertia follows the redirect as a GET request, updating the page with fresh data.
            </p>
            <p>This prevents the browser from showing "resubmit form?" dialogs on page refresh.</p>
          </div>
        </FeatureCard>
      </div>
    </div>
  </AppLayout>
</template>
