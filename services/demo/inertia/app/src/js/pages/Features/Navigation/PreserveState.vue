<script setup lang="ts">
import { ref } from "vue";
import { router, usePage } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Button } from "@/js/components/ui/button";
import AppLayout from "@/js/layouts/AppLayout.vue";
import type { SharedPageProps } from "@/js/types";

withDefaults(defineProps<{ serverCounter?: number; timestamp?: string }>(), {
  serverCounter: 0,
  timestamp: "",
});

const page = usePage<SharedPageProps>();
const localCounter = ref(0);

const breadcrumbs = [{ title: "Features" }, { title: "Navigation" }, { title: "Preserve State" }];

function reloadWithState() {
  router.reload({ preserveState: true } as never);
}

function reloadWithoutState() {
  router.visit(page.url, { preserveState: false });
}
</script>

<template>
  <AppLayout title="Preserve State" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Preserve State"
        description="By default, page visits reset local component state. Use preserveState to keep local reactive data intact across reloads."
      />

      <p class="text-muted-foreground text-sm">
        Server timestamp:
        <code class="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">{{ timestamp }}</code>
      </p>

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard
          title="Local Counter"
          description="This counter lives in local component state (ref). It resets on navigation unless preserveState is used."
        >
          <div class="flex items-center gap-4">
            <span class="text-3xl font-semibold tabular-nums">{{ localCounter }}</span>
            <Button @click="localCounter++">Increment</Button>
          </div>
        </FeatureCard>

        <FeatureCard
          title="Server Counter"
          description="This value comes from the server. It is always fresh after a reload."
        >
          <div class="flex items-center gap-4">
            <span class="text-3xl font-semibold tabular-nums">{{ serverCounter }}</span>
          </div>
        </FeatureCard>

        <FeatureCard
          title="Reload with preserveState"
          description="The local counter keeps its value because preserveState is true."
        >
          <div class="flex flex-wrap gap-3">
            <Button @click="reloadWithState">Reload (preserve state)</Button>
          </div>
          <p class="text-muted-foreground mt-2 text-xs">
            Increment the local counter, then click this button. The counter stays the same while
            the timestamp updates.
          </p>
        </FeatureCard>

        <FeatureCard
          title="Reload without preserveState"
          description="The local counter resets to 0 because the component is re-mounted."
        >
          <div class="flex flex-wrap gap-3">
            <Button variant="outline" @click="reloadWithoutState">Reload (reset state)</Button>
          </div>
          <p class="text-muted-foreground mt-2 text-xs">
            Increment the local counter, then click this button. The counter resets to 0.
          </p>
        </FeatureCard>
      </div>
    </div>
  </AppLayout>
</template>
