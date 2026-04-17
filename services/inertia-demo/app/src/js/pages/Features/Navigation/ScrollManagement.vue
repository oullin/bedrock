<script setup lang="ts">
import { Link, usePage } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import AppLayout from "@/js/layouts/AppLayout.vue";
import type { SharedPageProps } from "@/js/types";

withDefaults(defineProps<{ timestamp?: string; items?: Array<{ id: number; title: string }> }>(), {
  timestamp: "",
  items: () => [],
});

const page = usePage<SharedPageProps>();

const breadcrumbs = [
  { title: "Features" },
  { title: "Navigation" },
  { title: "Scroll Management" },
];
</script>

<template>
  <AppLayout title="Scroll Management" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Scroll Management"
        description="Inertia resets the scroll position to the top after each visit by default. Use preserve-scroll to maintain scroll position, which is useful for lists and infinite scroll patterns."
      />

      <p class="text-muted-foreground text-sm">
        Server timestamp:
        <code class="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">{{ timestamp }}</code>
      </p>

      <FeatureCard
        title="Scroll Controls"
        description="Scroll down through the list below, then click a reload link to test scroll behavior."
      >
        <div class="flex flex-wrap gap-3">
          <Link
            :href="page.url"
            preserve-scroll
            class="inline-flex h-9 items-center rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground hover:bg-primary/90"
          >
            Reload (preserve scroll)
          </Link>
          <Link
            :href="page.url"
            class="inline-flex h-9 items-center rounded-md border bg-background px-4 text-sm font-medium shadow-xs hover:bg-accent hover:text-accent-foreground"
          >
            Reload (reset to top)
          </Link>
        </div>
      </FeatureCard>

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
      </div>

      <FeatureCard
        title="Bottom Controls"
        description="You scrolled to the bottom. Try the reload links here too."
      >
        <div class="flex flex-wrap gap-3">
          <Link
            :href="page.url"
            preserve-scroll
            class="inline-flex h-9 items-center rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground hover:bg-primary/90"
          >
            Reload (preserve scroll)
          </Link>
          <Link
            :href="page.url"
            class="inline-flex h-9 items-center rounded-md border bg-background px-4 text-sm font-medium shadow-xs hover:bg-accent hover:text-accent-foreground"
          >
            Reload (reset to top)
          </Link>
        </div>
      </FeatureCard>
    </div>
  </AppLayout>
</template>
