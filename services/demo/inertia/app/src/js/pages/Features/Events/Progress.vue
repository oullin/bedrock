<script setup lang="ts">
import { Link } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Button } from "@/js/components/ui/button";
import { Badge } from "@/js/components/ui/badge";
import AppLayout from "@/js/layouts/AppLayout.vue";
import { featureRoute } from "@/js/lib/routes";

const breadcrumbs = [{ title: "Features" }, { title: "Events" }, { title: "Progress" }];

const deferredPropsRoute = featureRoute("features.data-loading.deferred-props");
const pollingRoute = featureRoute("features.data-loading.polling");
</script>

<template>
  <AppLayout title="Progress" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Progress Indicator"
        description="Inertia shows a progress bar during page navigations. The loading bar at the top of this app demonstrates this feature."
      />

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard title="Try It" description="Navigate to see the progress bar in action.">
          <div class="space-y-4">
            <p class="text-sm">
              Click a link below and watch the loading bar at the top of the page. Slower endpoints
              make the progress bar more visible.
            </p>
            <div class="flex flex-wrap gap-2">
              <Button v-if="deferredPropsRoute" as-child variant="outline">
                <Link :href="deferredPropsRoute"> Deferred Props (slow) </Link>
              </Button>
              <Button v-if="pollingRoute" as-child variant="outline">
                <Link :href="pollingRoute"> Polling Page </Link>
              </Button>
            </div>
          </div>
        </FeatureCard>

        <FeatureCard title="Configuration" description="How the progress bar is configured.">
          <div class="space-y-3 text-sm">
            <pre class="bg-muted overflow-auto rounded p-4 text-xs leading-relaxed">
createInertiaApp({
  progress: {
    // Delay before showing (ms)
    delay: 250,

    // Color of the bar
    color: '#4B5563',

    // Show the spinner
    includeCSS: true,

    // Show spinner icon
    showSpinner: false,
  },
})</pre
            >
            <p class="text-muted-foreground">
              The progress bar only appears after a short delay to avoid flickering on fast
              navigations.
            </p>
          </div>
        </FeatureCard>

        <FeatureCard
          title="Custom Progress Bar"
          description="This app uses a custom LoadingBar component."
        >
          <div class="space-y-3 text-sm">
            <p>
              Instead of using the built-in NProgress bar, this demo app uses a custom
              <code class="bg-muted rounded px-1.5 py-0.5 text-xs">LoadingBar</code> component that
              listens to router events.
            </p>
            <div class="space-y-2">
              <div class="flex items-center gap-2 rounded-md border p-2">
                <Badge variant="secondary">start</Badge>
                <span class="text-muted-foreground text-xs">Show the progress bar</span>
              </div>
              <div class="flex items-center gap-2 rounded-md border p-2">
                <Badge variant="secondary">progress</Badge>
                <span class="text-muted-foreground text-xs">Update the bar width</span>
              </div>
              <div class="flex items-center gap-2 rounded-md border p-2">
                <Badge variant="secondary">finish</Badge>
                <span class="text-muted-foreground text-xs">Complete and hide the bar</span>
              </div>
            </div>
          </div>
        </FeatureCard>

        <FeatureCard title="Disabling Progress" description="Opt out of the progress indicator.">
          <div class="space-y-3 text-sm">
            <p>You can disable the progress bar globally or per-visit:</p>
            <pre class="bg-muted overflow-auto rounded p-3 text-xs leading-relaxed">
// Disable globally
createInertiaApp({
  progress: false,
})

// Disable per-visit
router.visit(url, {
  showProgress: false,
})</pre
            >
          </div>
        </FeatureCard>
      </div>
    </div>
  </AppLayout>
</template>
