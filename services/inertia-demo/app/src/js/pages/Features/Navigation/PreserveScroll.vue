<script setup lang="ts">
import { Link, usePage } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import AppLayout from "@/js/layouts/AppLayout.vue";
import type { SharedPageProps } from "@/js/types";

withDefaults(defineProps<{ timestamp?: string }>(), {
  timestamp: "",
});

const page = usePage<SharedPageProps>();

const breadcrumbs = [{ title: "Features" }, { title: "Navigation" }, { title: "Preserve Scroll" }];
</script>

<template>
  <AppLayout title="Preserve Scroll" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Preserve Scroll"
        description="By default, Inertia resets the scroll position to the top after navigation. Use preserve-scroll to maintain the current position."
      />

      <p class="text-muted-foreground text-sm">
        Server timestamp:
        <code class="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">{{ timestamp }}</code>
      </p>

      <FeatureCard
        title="Scroll Actions"
        description="Use the links below to reload the page. Scroll down first to see the difference."
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
            Reload (reset scroll)
          </Link>
        </div>
      </FeatureCard>

      <div class="space-y-4">
        <FeatureCard
          v-for="n in 20"
          :key="n"
          :title="`Section ${n}`"
          :description="`Scroll content block ${n} of 20. Scroll down and click a reload link to test scroll preservation.`"
        >
          <p class="text-muted-foreground text-sm">
            This is filler content for section {{ n }}. It exists so you can scroll the page and
            observe how preserve-scroll works when navigating.
          </p>
        </FeatureCard>
      </div>
    </div>
  </AppLayout>
</template>
