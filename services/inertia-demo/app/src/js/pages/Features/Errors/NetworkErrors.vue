<script setup lang="ts">
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Badge } from "@/js/components/ui/badge";
import AppLayout from "@/js/layouts/AppLayout.vue";

const breadcrumbs = [{ title: "Features" }, { title: "Errors" }, { title: "Network Errors" }];
</script>

<template>
  <AppLayout title="Network Errors" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Network Errors"
        description="Handle network failures (timeouts, offline, DNS errors) that prevent the request from reaching the server."
      />

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard
          title="What Are Network Errors?"
          description="When the request never reaches the server."
        >
          <div class="space-y-3 text-sm">
            <p>
              Network errors occur when the browser cannot complete the HTTP request. Unlike HTTP
              exceptions (which get a response with a status code), network errors produce no
              response at all.
            </p>
            <div class="space-y-2">
              <div class="flex items-center gap-2 rounded-md border p-2">
                <Badge variant="outline">Offline</Badge>
                <span class="text-muted-foreground text-xs">Device has no internet connection</span>
              </div>
              <div class="flex items-center gap-2 rounded-md border p-2">
                <Badge variant="outline">Timeout</Badge>
                <span class="text-muted-foreground text-xs">Server took too long to respond</span>
              </div>
              <div class="flex items-center gap-2 rounded-md border p-2">
                <Badge variant="outline">DNS</Badge>
                <span class="text-muted-foreground text-xs">Domain could not be resolved</span>
              </div>
              <div class="flex items-center gap-2 rounded-md border p-2">
                <Badge variant="outline">CORS</Badge>
                <span class="text-muted-foreground text-xs">Cross-origin request blocked</span>
              </div>
            </div>
          </div>
        </FeatureCard>

        <FeatureCard
          title="Handling Network Errors"
          description="How to detect and respond to failures."
        >
          <div class="space-y-4 text-sm">
            <div class="space-y-2">
              <p class="font-medium">Global listener</p>
              <pre class="bg-muted overflow-auto rounded p-3 text-xs leading-relaxed">
router.on('error', (event) => {
  // Show a notification
  alert('A network error occurred.')
})</pre
              >
            </div>
            <div class="space-y-2">
              <p class="font-medium">Per-visit callback</p>
              <pre class="bg-muted overflow-auto rounded p-3 text-xs leading-relaxed">
router.visit(url, {
  onError(errors) {
    // Handle the error
  },
})</pre
              >
            </div>
          </div>
        </FeatureCard>

        <FeatureCard title="Testing Network Errors" description="How to simulate network failures.">
          <div class="space-y-3 text-sm">
            <p>To test network error handling, you can:</p>
            <ul class="space-y-2">
              <li class="flex items-start gap-2">
                <Badge variant="secondary" class="shrink-0 text-xs">1</Badge>
                <span>Open DevTools, go to Network tab, and select "Offline" mode.</span>
              </li>
              <li class="flex items-start gap-2">
                <Badge variant="secondary" class="shrink-0 text-xs">2</Badge>
                <span
                  >Use DevTools to throttle the connection to "Slow 3G" and set a short
                  timeout.</span
                >
              </li>
              <li class="flex items-start gap-2">
                <Badge variant="secondary" class="shrink-0 text-xs">3</Badge>
                <span>Stop your development server while the app is loaded.</span>
              </li>
            </ul>
            <p class="text-muted-foreground">
              Then try navigating or submitting a form to see the error handling in action.
            </p>
          </div>
        </FeatureCard>

        <FeatureCard title="Best Practices" description="Resilient error handling patterns.">
          <ul class="space-y-2 text-sm">
            <li class="flex items-start gap-2">
              <span class="text-primary font-medium">Show a toast</span> - Let users know something
              went wrong without disrupting the page.
            </li>
            <li class="flex items-start gap-2">
              <span class="text-primary font-medium">Offer retry</span> - Provide a button to retry
              the failed request.
            </li>
            <li class="flex items-start gap-2">
              <span class="text-primary font-medium">Preserve state</span> - Do not clear form data
              on network errors.
            </li>
            <li class="flex items-start gap-2">
              <span class="text-primary font-medium">Queue offline</span> - For critical actions,
              queue them and retry when connectivity returns.
            </li>
          </ul>
        </FeatureCard>
      </div>
    </div>
  </AppLayout>
</template>
