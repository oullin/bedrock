<script setup lang="ts">
import { Link } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Button } from "@/js/components/ui/button";
import { Badge } from "@/js/components/ui/badge";
import PersistentDemoLayout from "@/js/layouts/PersistentDemoLayout.vue";
import { featureRoute } from "@/js/lib/routes";

defineOptions({
  layout: PersistentDemoLayout,
});

const page1Route = featureRoute("features.layouts.persistent-layouts");
</script>

<template>
  <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
    <FeatureHeader
      title="Persistent Layouts - Page 2"
      description="You navigated here but the layout state (stopwatch and counter) was preserved!"
    />

    <div class="grid gap-6 lg:grid-cols-2">
      <FeatureCard title="State Preserved" description="The layout was not remounted.">
        <div class="space-y-3 text-sm">
          <p>
            Notice that the stopwatch above kept ticking and the counter retained its value. This
            happened because both this page and Page 1 use the same persistent layout component.
          </p>
          <p>
            Inertia detects that the layout class is the same and reuses the existing instance
            instead of creating a new one.
          </p>
        </div>
      </FeatureCard>

      <FeatureCard
        title="Navigate Back"
        description="Return to page 1 to confirm state persistence."
      >
        <div class="space-y-4">
          <p class="text-sm">Go back to page 1. The layout state will continue to persist.</p>
          <Button v-if="page1Route" as-child>
            <Link :href="page1Route"> Back to Page 1 </Link>
          </Button>
          <div class="flex items-center gap-2">
            <Badge variant="secondary">Current Page</Badge>
            <span class="text-sm">Page 2</span>
          </div>
        </div>
      </FeatureCard>
    </div>
  </div>
</template>
