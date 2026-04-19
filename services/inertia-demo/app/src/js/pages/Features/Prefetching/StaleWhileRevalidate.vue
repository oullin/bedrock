<script setup lang="ts">
import { Link } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Badge } from "@/js/components/ui/badge";
import AppLayout from "@/js/layouts/AppLayout.vue";
import { featureRoute } from "@/js/lib/routes";

const breadcrumbs = [
  { title: "Features" },
  { title: "Prefetching" },
  { title: "Stale While Revalidate" },
];

const linkPrefetchRoute = featureRoute("features.prefetching.link-prefetch");
const manualPrefetchRoute = featureRoute("features.prefetching.manual-prefetch");
const cacheManagementRoute = featureRoute("features.prefetching.cache-management");
</script>

<template>
  <AppLayout title="Stale While Revalidate" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Stale While Revalidate"
        description="Show cached (stale) data immediately while fetching fresh data in the background. When the fresh data arrives, the page updates seamlessly."
      />

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard title="How SWR Works" description="The Stale-While-Revalidate pattern.">
          <div class="space-y-3 text-sm">
            <div class="flex items-start gap-3 rounded-md border p-3">
              <Badge class="shrink-0">1</Badge>
              <p>User hovers a link with prefetch. The response is cached.</p>
            </div>
            <div class="flex items-start gap-3 rounded-md border p-3">
              <Badge class="shrink-0" variant="secondary">2</Badge>
              <p>User clicks the link. The cached response is shown immediately (stale data).</p>
            </div>
            <div class="flex items-start gap-3 rounded-md border p-3">
              <Badge class="shrink-0" variant="outline">3</Badge>
              <p>
                In the background, a fresh request is made. When it completes, the page updates with
                new data.
              </p>
            </div>
          </div>
        </FeatureCard>

        <FeatureCard
          title="SWR Cache Configuration"
          description="Configure cache freshness and stale windows."
        >
          <div class="space-y-3 text-sm">
            <p>
              Use the <code class="bg-muted rounded px-1.5 py-0.5 text-xs">cacheFor</code> prop with
              an array to set both fresh and stale durations:
            </p>
            <pre class="bg-muted overflow-auto rounded p-3 text-xs">
&lt;Link
  :href="url"
  prefetch="hover"
  :cache-for="[5000, 30000]"
&gt;</pre
            >
            <p class="text-muted-foreground">
              The first value (5s) is how long the data is considered fresh. The second value (30s)
              is how long stale data can be shown while revalidating.
            </p>
          </div>
        </FeatureCard>

        <FeatureCard
          title="Try It"
          description="Navigate between prefetching pages to see SWR in action."
        >
          <div class="flex flex-wrap gap-2">
            <Link
              v-if="linkPrefetchRoute"
              :href="linkPrefetchRoute"
              prefetch="hover"
              :cache-for="[5000, 30000]"
              class="inline-flex items-center gap-2 rounded-md border px-4 py-2 text-sm font-medium hover:bg-accent"
            >
              Link Prefetch
              <Badge variant="secondary">SWR</Badge>
            </Link>
            <Link
              v-if="manualPrefetchRoute"
              :href="manualPrefetchRoute"
              prefetch="hover"
              :cache-for="[5000, 30000]"
              class="inline-flex items-center gap-2 rounded-md border px-4 py-2 text-sm font-medium hover:bg-accent"
            >
              Manual Prefetch
              <Badge variant="secondary">SWR</Badge>
            </Link>
            <Link
              v-if="cacheManagementRoute"
              :href="cacheManagementRoute"
              prefetch="hover"
              :cache-for="[5000, 30000]"
              class="inline-flex items-center gap-2 rounded-md border px-4 py-2 text-sm font-medium hover:bg-accent"
            >
              Cache Management
              <Badge variant="secondary">SWR</Badge>
            </Link>
          </div>
        </FeatureCard>

        <FeatureCard title="Benefits" description="Why use Stale-While-Revalidate?">
          <ul class="text-muted-foreground space-y-2 text-sm">
            <li class="flex items-start gap-2">
              <span class="text-foreground font-medium">Instant navigation</span> - Users see
              content immediately.
            </li>
            <li class="flex items-start gap-2">
              <span class="text-foreground font-medium">Fresh data</span> - Background revalidation
              keeps data current.
            </li>
            <li class="flex items-start gap-2">
              <span class="text-foreground font-medium">Reduced load</span> - Cached responses avoid
              redundant requests during the fresh window.
            </li>
          </ul>
        </FeatureCard>
      </div>
    </div>
  </AppLayout>
</template>
