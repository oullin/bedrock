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

const page2Route = featureRoute("features.layouts.persistent-layouts-page-2");
</script>

<template>
  <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
    <FeatureHeader
      title="Persistent Layouts - Page 1"
      description="This page uses a persistent layout. The stopwatch and counter in the bar above will not reset when navigating to page 2."
    />

    <div class="grid gap-6 lg:grid-cols-2">
      <FeatureCard title="How It Works" description="Persistent layouts survive page navigations.">
        <div class="space-y-3 text-sm">
          <p>
            Normally, navigating between pages causes Vue to unmount and remount the entire
            component tree, including the layout. Persistent layouts prevent this by using
            <code class="bg-muted rounded px-1.5 py-0.5 text-xs"
              >defineOptions({ layout: MyLayout })</code
            >.
          </p>
          <p>
            Notice the stopwatch and counter above -- they maintain their state even when you
            navigate between pages that share the same layout.
          </p>
        </div>
      </FeatureCard>

      <FeatureCard title="Navigate" description="Go to page 2 without resetting layout state.">
        <div class="space-y-4">
          <p class="text-sm">
            Click the button below to navigate to page 2. The stopwatch and counter in the
            persistent layout bar above will keep running and maintain their values.
          </p>
          <Button v-if="page2Route" as-child>
            <Link :href="page2Route"> Go to Page 2 </Link>
          </Button>
          <div class="flex items-center gap-2">
            <Badge variant="secondary">Current Page</Badge>
            <span class="text-sm">Page 1</span>
          </div>
        </div>
      </FeatureCard>
    </div>
  </div>
</template>
