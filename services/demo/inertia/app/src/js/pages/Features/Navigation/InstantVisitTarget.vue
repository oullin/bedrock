<script setup lang="ts">
import { Link } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import AppLayout from "@/js/layouts/AppLayout.vue";
import { useDemoRoute } from "@/js/lib/routes";

withDefaults(
  defineProps<{
    greeting?: string;
    serverTimestamp?: string;
    items?: Array<{ id: number; title: string }>;
  }>(),
  {
    greeting: "",
    serverTimestamp: "",
    items: () => [],
  },
);

const breadcrumbs = [
  { title: "Features" },
  { title: "Navigation" },
  { title: "Instant Visits" },
  { title: "Target" },
];

const sourceRoute = useDemoRoute("features.navigation.instant-visits");
</script>

<template>
  <AppLayout title="Instant Visit Target" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader title="Instant Visit Target" :description="greeting" />

      <p class="text-muted-foreground text-sm">
        Server timestamp:
        <code class="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">{{ serverTimestamp }}</code>
      </p>

      <FeatureCard
        title="Go Back"
        description="Navigate back to the source page to try the instant visit again."
      >
        <div class="flex flex-wrap gap-3">
          <Link
            :href="sourceRoute.url"
            class="inline-flex h-9 items-center rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground hover:bg-primary/90"
          >
            Back to Instant Visits
          </Link>
        </div>
      </FeatureCard>

      <FeatureCard
        title="Server Data"
        description="This data was loaded from the server (possibly prefetched)."
      >
        <div class="space-y-2">
          <div
            v-for="item in items"
            :key="item.id"
            class="flex items-center gap-4 rounded-lg bg-muted/30 px-4 py-3"
          >
            <span
              class="flex size-8 shrink-0 items-center justify-center rounded-full bg-primary/10 text-xs font-medium text-primary"
            >
              {{ item.id }}
            </span>
            <span class="text-sm font-medium">{{ item.title }}</span>
          </div>
          <div v-if="items.length === 0" class="text-muted-foreground py-4 text-center text-sm">
            No items loaded.
          </div>
        </div>
      </FeatureCard>
    </div>
  </AppLayout>
</template>
