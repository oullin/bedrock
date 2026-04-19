<script setup lang="ts">
import { router } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Button } from "@/js/components/ui/button";
import { Badge } from "@/js/components/ui/badge";
import AppLayout from "@/js/layouts/AppLayout.vue";

withDefaults(
  defineProps<{
    page?: string;
    staticData?: Record<string, unknown>;
    freshData?: Record<string, unknown>;
    dynamicData?: Record<string, unknown>;
  }>(),
  {
    page: "",
    staticData: () => ({}),
    freshData: () => ({}),
    dynamicData: () => ({}),
  },
);

const breadcrumbs = [{ title: "Features" }, { title: "Data Loading" }, { title: "Once Props" }];

function reloadAll() {
  router.reload();
}

function reloadFreshOnly() {
  router.reload({ only: ["freshData", "dynamicData"] });
}
</script>

<template>
  <AppLayout title="Once Props" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Once Props"
        description="Props marked as 'once' are only sent on the initial page load. Subsequent reloads skip them, reducing server work for static data."
      />

      <div class="flex gap-2">
        <Button @click="reloadAll">Reload All Props</Button>
        <Button variant="outline" @click="reloadFreshOnly">Reload Fresh Only</Button>
      </div>

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard
          title="Static Data (Once)"
          description="Sent only on the first visit. Cached on subsequent reloads."
        >
          <div class="space-y-2">
            <Badge>Cached</Badge>
            <pre class="bg-muted mt-2 overflow-auto rounded p-3 text-xs">{{
              JSON.stringify(staticData, null, 2)
            }}</pre>
            <p class="text-muted-foreground text-xs">
              This data will not change when you click "Reload All Props" because the server skips
              once-props on partial reloads.
            </p>
          </div>
        </FeatureCard>

        <FeatureCard title="Fresh Data" description="Always re-fetched on every request.">
          <div class="space-y-2">
            <Badge variant="secondary">Fresh</Badge>
            <pre class="bg-muted mt-2 overflow-auto rounded p-3 text-xs">{{
              JSON.stringify(freshData, null, 2)
            }}</pre>
            <p class="text-muted-foreground text-xs">
              This data updates on every reload because it is a regular prop.
            </p>
          </div>
        </FeatureCard>

        <FeatureCard
          title="Dynamic Data"
          description="Changes on every reload to demonstrate the difference."
        >
          <div class="space-y-2">
            <Badge variant="outline">Dynamic</Badge>
            <pre class="bg-muted mt-2 overflow-auto rounded p-3 text-xs">{{
              JSON.stringify(dynamicData, null, 2)
            }}</pre>
          </div>
        </FeatureCard>

        <FeatureCard title="Page Identifier" description="Which page load we are on.">
          <div class="text-2xl font-semibold">{{ page }}</div>
          <p class="text-muted-foreground mt-2 text-sm">
            Compare the timestamps in static vs fresh data after reloading to see that once-props
            retain their original values.
          </p>
        </FeatureCard>
      </div>
    </div>
  </AppLayout>
</template>
