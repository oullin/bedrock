<script setup lang="ts">
import { Deferred } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Skeleton } from "@/js/components/ui/skeleton";
import { Badge } from "@/js/components/ui/badge";
import AppLayout from "@/js/layouts/AppLayout.vue";

withDefaults(
  defineProps<{
    quickStat?: number;
    slowStats?: Record<string, unknown>;
    heavyData?: Array<{ name?: string; status?: string }>;
  }>(),
  {
    quickStat: 0,
    slowStats: undefined,
    heavyData: undefined,
  },
);

const breadcrumbs = [{ title: "Features" }, { title: "Data Loading" }, { title: "Deferred Props" }];
</script>

<template>
  <AppLayout title="Deferred Props" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Deferred Props"
        description="Deferred props allow you to load slow data after the initial page render, showing skeleton fallbacks until the data arrives."
      />

      <div class="grid gap-6 lg:grid-cols-3">
        <FeatureCard title="Quick Stat" description="This prop loads immediately with the page.">
          <div class="text-3xl font-semibold tracking-tight">{{ quickStat }}</div>
          <p class="text-muted-foreground mt-1 text-sm">Available instantly</p>
        </FeatureCard>

        <FeatureCard
          title="Slow Stats"
          description="These stats are deferred and load after the initial render."
        >
          <Deferred data="slowStats">
            <template #fallback>
              <div class="space-y-3">
                <Skeleton class="h-8 w-24" />
                <Skeleton class="h-4 w-32" />
                <Skeleton class="h-4 w-28" />
              </div>
            </template>
            <div class="space-y-2">
              <div class="text-2xl font-semibold">{{ slowStats?.total }}</div>
              <div class="flex items-center gap-2">
                <Badge variant="secondary">Active: {{ slowStats?.active }}</Badge>
                <Badge variant="outline">Inactive: {{ slowStats?.inactive }}</Badge>
              </div>
            </div>
          </Deferred>
        </FeatureCard>

        <FeatureCard title="Heavy Data" description="A large dataset loaded in the background.">
          <Deferred data="heavyData">
            <template #fallback>
              <div class="space-y-2">
                <Skeleton v-for="i in 5" :key="i" class="h-6 w-full" />
              </div>
            </template>
            <div class="space-y-2">
              <div
                v-for="(item, index) in heavyData"
                :key="index"
                class="flex items-center justify-between rounded-md border p-2 text-sm"
              >
                <span>{{ item.name ?? item }}</span>
                <Badge v-if="item.status" variant="secondary">{{ item.status }}</Badge>
              </div>
              <p v-if="!heavyData?.length" class="text-muted-foreground text-sm">
                No data available.
              </p>
            </div>
          </Deferred>
        </FeatureCard>
      </div>
    </div>
  </AppLayout>
</template>
