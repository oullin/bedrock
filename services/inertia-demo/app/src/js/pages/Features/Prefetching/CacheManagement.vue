<script setup lang="ts">
import { router } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Button } from "@/js/components/ui/button";
import { Badge } from "@/js/components/ui/badge";
import AppLayout from "@/js/layouts/AppLayout.vue";

const breadcrumbs = [
  { title: "Features" },
  { title: "Prefetching" },
  { title: "Cache Management" },
];

function clearCache() {
  router.clearHistory();
}

function flushAll() {
  router.flushAll();
}
</script>

<template>
  <AppLayout title="Cache Management" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Cache Management"
        description="Understand the cache lifecycle for prefetched pages and how to manage cache invalidation."
      />

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard title="Cache Lifecycle" description="How prefetch cache entries are managed.">
          <div class="space-y-3 text-sm">
            <div class="flex items-start gap-3 rounded-md border p-3">
              <Badge class="shrink-0">Fresh</Badge>
              <p>
                The cached response is within the fresh window. It will be served directly without
                any network request.
              </p>
            </div>
            <div class="flex items-start gap-3 rounded-md border p-3">
              <Badge class="shrink-0" variant="secondary">Stale</Badge>
              <p>
                The response is past its fresh window but still within the stale window. It is
                served immediately while a background revalidation runs.
              </p>
            </div>
            <div class="flex items-start gap-3 rounded-md border p-3">
              <Badge class="shrink-0" variant="outline">Expired</Badge>
              <p>
                The cache entry has exceeded both windows and is discarded. A full request is made
                on the next navigation.
              </p>
            </div>
          </div>
        </FeatureCard>

        <FeatureCard
          title="Cache Invalidation"
          description="Ways to clear or invalidate cached data."
        >
          <div class="space-y-4">
            <div class="space-y-3 text-sm">
              <p>
                Cached responses are automatically invalidated when a non-GET request (POST, PUT,
                DELETE) is made, since those requests typically change server state.
              </p>
              <p>You can also manually clear the cache using the methods below:</p>
            </div>
            <div class="flex gap-2">
              <Button variant="outline" @click="clearCache">Clear History</Button>
              <Button variant="outline" @click="flushAll">Flush All</Button>
            </div>
          </div>
        </FeatureCard>

        <FeatureCard
          title="cacheFor Options"
          description="Different ways to configure cache duration."
        >
          <div class="space-y-3">
            <div class="rounded-md border p-3">
              <p class="text-sm font-medium">Number (milliseconds)</p>
              <code class="text-muted-foreground text-xs">:cache-for="30000"</code>
              <p class="text-muted-foreground mt-1 text-xs">
                Cache for 30 seconds (fresh only, no stale window).
              </p>
            </div>
            <div class="rounded-md border p-3">
              <p class="text-sm font-medium">Array [fresh, stale]</p>
              <code class="text-muted-foreground text-xs">:cache-for="[5000, 60000]"</code>
              <p class="text-muted-foreground mt-1 text-xs">
                Fresh for 5s, then stale-while-revalidate up to 60s.
              </p>
            </div>
            <div class="rounded-md border p-3">
              <p class="text-sm font-medium">String shorthand</p>
              <code class="text-muted-foreground text-xs">cache-for="1m"</code>
              <p class="text-muted-foreground mt-1 text-xs">Human-readable duration strings.</p>
            </div>
          </div>
        </FeatureCard>

        <FeatureCard title="Best Practices" description="Tips for effective cache management.">
          <ul class="space-y-2 text-sm">
            <li class="flex items-start gap-2">
              <span class="text-primary font-medium">1.</span>
              Keep cache durations short for frequently changing data.
            </li>
            <li class="flex items-start gap-2">
              <span class="text-primary font-medium">2.</span>
              Use SWR (stale-while-revalidate) for data that can tolerate brief staleness.
            </li>
            <li class="flex items-start gap-2">
              <span class="text-primary font-medium">3.</span>
              POST/PUT/DELETE requests automatically bust the cache.
            </li>
            <li class="flex items-start gap-2">
              <span class="text-primary font-medium">4.</span>
              Consider using
              <code class="bg-muted rounded px-1 py-0.5 text-xs">prefetch="mount"</code> for links
              the user is very likely to click.
            </li>
          </ul>
        </FeatureCard>
      </div>
    </div>
  </AppLayout>
</template>
