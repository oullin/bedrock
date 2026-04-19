<script setup lang="ts">
import { Link } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Badge } from "@/js/components/ui/badge";
import AppLayout from "@/js/layouts/AppLayout.vue";
import { featureRoute } from "@/js/lib/routes";

const breadcrumbs = [{ title: "Features" }, { title: "Prefetching" }, { title: "Link Prefetch" }];

const swrRoute = featureRoute("features.prefetching.stale-while-revalidate");
const manualPrefetchRoute = featureRoute("features.prefetching.manual-prefetch");
const cacheManagementRoute = featureRoute("features.prefetching.cache-management");
</script>

<template>
  <AppLayout title="Link Prefetch" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Link Prefetch"
        description="Prefetch page data before the user clicks a link, making navigation feel instant."
      />

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard
          title="Prefetch on Hover"
          description="Data is fetched when the user hovers over the link."
        >
          <div class="space-y-3">
            <Link
              v-if="swrRoute"
              :href="swrRoute"
              prefetch="hover"
              class="inline-flex items-center gap-2 rounded-md border px-4 py-2 text-sm font-medium hover:bg-accent"
            >
              Hover me to prefetch
              <Badge variant="secondary">hover</Badge>
            </Link>
            <p class="text-muted-foreground text-xs">
              The page data starts loading as soon as your mouse enters this link.
            </p>
          </div>
        </FeatureCard>

        <FeatureCard
          title="Prefetch on Click"
          description="Data is fetched on mouse down, before the click event fires."
        >
          <div class="space-y-3">
            <Link
              v-if="manualPrefetchRoute"
              :href="manualPrefetchRoute"
              prefetch="click"
              class="inline-flex items-center gap-2 rounded-md border px-4 py-2 text-sm font-medium hover:bg-accent"
            >
              Click prefetch
              <Badge variant="secondary">click</Badge>
            </Link>
            <p class="text-muted-foreground text-xs">
              Prefetching begins on mousedown, saving a few hundred milliseconds.
            </p>
          </div>
        </FeatureCard>

        <FeatureCard
          title="Prefetch on Mount"
          description="Data is prefetched immediately when the component mounts."
        >
          <div class="space-y-3">
            <Link
              v-if="cacheManagementRoute"
              :href="cacheManagementRoute"
              prefetch="mount"
              class="inline-flex items-center gap-2 rounded-md border px-4 py-2 text-sm font-medium hover:bg-accent"
            >
              Already prefetched
              <Badge variant="secondary">mount</Badge>
            </Link>
            <p class="text-muted-foreground text-xs">
              This page's data was fetched as soon as this component rendered.
            </p>
          </div>
        </FeatureCard>

        <FeatureCard
          title="Cache Duration"
          description="Control how long prefetched data stays fresh."
        >
          <div class="space-y-3">
            <Link
              v-if="swrRoute"
              :href="swrRoute"
              prefetch="hover"
              :cache-for="30000"
              class="inline-flex items-center gap-2 rounded-md border px-4 py-2 text-sm font-medium hover:bg-accent"
            >
              Cached for 30s
              <Badge variant="outline">cacheFor=30s</Badge>
            </Link>
            <p class="text-muted-foreground text-xs">
              Use <code class="bg-muted rounded px-1 py-0.5 text-xs">cacheFor</code> to specify how
              long the prefetched response should be cached.
            </p>
          </div>
        </FeatureCard>
      </div>
    </div>
  </AppLayout>
</template>
