<script setup lang="ts">
import { Deferred, router } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Button } from "@/js/components/ui/button";
import { Badge } from "@/js/components/ui/badge";
import { Skeleton } from "@/js/components/ui/skeleton";
import AppLayout from "@/js/layouts/AppLayout.vue";

withDefaults(
  defineProps<{
    regularData?: Record<string, unknown>;
    optionalData?: Record<string, unknown>;
    deferredData?: Record<string, unknown>;
  }>(),
  {
    regularData: () => ({}),
    optionalData: undefined,
    deferredData: undefined,
  },
);

const breadcrumbs = [{ title: "Features" }, { title: "Data Loading" }, { title: "Optional Props" }];

function loadOptional() {
  router.reload({ only: ["optionalData"] });
}

function loadAll() {
  router.reload();
}
</script>

<template>
  <AppLayout title="Optional Props" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Optional Props"
        description="Optional props are not included in the initial page load. They are only sent when explicitly requested via partial reloads."
      />

      <div class="flex gap-2">
        <Button @click="loadOptional">Load Optional Data</Button>
        <Button variant="outline" @click="loadAll">Reload All</Button>
      </div>

      <div class="grid gap-6 lg:grid-cols-3">
        <FeatureCard title="Regular Data" description="Always included in every response.">
          <div class="space-y-2">
            <div class="flex items-center gap-2">
              <Badge>Always Loaded</Badge>
            </div>
            <pre class="bg-muted mt-2 overflow-auto rounded p-3 text-xs">{{
              JSON.stringify(regularData, null, 2)
            }}</pre>
          </div>
        </FeatureCard>

        <FeatureCard title="Optional Data" description="Only sent when explicitly requested.">
          <div class="space-y-2">
            <div class="flex items-center gap-2">
              <Badge variant="secondary">{{ optionalData ? "Loaded" : "Not Loaded" }}</Badge>
            </div>
            <template v-if="optionalData">
              <pre class="bg-muted overflow-auto rounded p-3 text-xs">{{
                JSON.stringify(optionalData, null, 2)
              }}</pre>
            </template>
            <p v-else class="text-muted-foreground text-sm">
              Click "Load Optional Data" to fetch this prop. It is skipped in the initial page load
              to reduce payload size.
            </p>
          </div>
        </FeatureCard>

        <FeatureCard title="Deferred Data" description="Loaded automatically after initial render.">
          <Deferred data="deferredData">
            <template #fallback>
              <div class="space-y-3">
                <Badge variant="outline">Loading...</Badge>
                <Skeleton class="h-6 w-full" />
                <Skeleton class="h-6 w-3/4" />
              </div>
            </template>
            <div class="space-y-2">
              <Badge variant="outline">Auto-Loaded</Badge>
              <pre class="bg-muted overflow-auto rounded p-3 text-xs">{{
                JSON.stringify(deferredData, null, 2)
              }}</pre>
            </div>
          </Deferred>
        </FeatureCard>
      </div>

      <FeatureCard title="Comparison">
        <div class="overflow-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="border-b">
                <th class="py-2 pr-4 text-left font-medium">Strategy</th>
                <th class="py-2 pr-4 text-left font-medium">Initial Load</th>
                <th class="py-2 text-left font-medium">When Available</th>
              </tr>
            </thead>
            <tbody class="text-muted-foreground">
              <tr class="border-b">
                <td class="py-2 pr-4 font-medium text-foreground">Regular</td>
                <td class="py-2 pr-4">Included</td>
                <td class="py-2">Immediately</td>
              </tr>
              <tr class="border-b">
                <td class="py-2 pr-4 font-medium text-foreground">Optional</td>
                <td class="py-2 pr-4">Excluded</td>
                <td class="py-2">On explicit request</td>
              </tr>
              <tr>
                <td class="py-2 pr-4 font-medium text-foreground">Deferred</td>
                <td class="py-2 pr-4">Excluded</td>
                <td class="py-2">Auto-fetched after render</td>
              </tr>
            </tbody>
          </table>
        </div>
      </FeatureCard>
    </div>
  </AppLayout>
</template>
