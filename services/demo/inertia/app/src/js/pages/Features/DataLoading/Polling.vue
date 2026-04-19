<script setup lang="ts">
import { usePoll } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Button } from "@/js/components/ui/button";
import { Badge } from "@/js/components/ui/badge";
import AppLayout from "@/js/layouts/AppLayout.vue";

withDefaults(
  defineProps<{ currentTime?: string; randomNumber?: number; contactCount?: number }>(),
  {
    currentTime: "",
    randomNumber: 0,
    contactCount: 0,
  },
);

const breadcrumbs = [{ title: "Features" }, { title: "Data Loading" }, { title: "Polling" }];

const { stop, start } = usePoll(5000);
</script>

<template>
  <AppLayout title="Polling" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Polling"
        description="Automatically refresh page data at regular intervals using usePoll. This page polls every 5 seconds."
      />

      <div class="flex gap-2">
        <Button variant="outline" @click="stop">Stop Polling</Button>
        <Button @click="start">Start Polling</Button>
      </div>

      <div class="grid gap-6 lg:grid-cols-3">
        <FeatureCard title="Current Time" description="Server time, refreshed every 5 seconds.">
          <div class="flex items-center gap-3">
            <div class="text-2xl font-mono font-semibold">{{ currentTime }}</div>
          </div>
          <p class="text-muted-foreground mt-2 text-xs">
            This value comes from the server and updates automatically.
          </p>
        </FeatureCard>

        <FeatureCard title="Random Number" description="A new random number on each poll.">
          <div class="text-3xl font-semibold tracking-tight">{{ randomNumber }}</div>
          <Badge class="mt-2" variant="secondary">Changes every 5s</Badge>
        </FeatureCard>

        <FeatureCard title="Contact Count" description="Live count of contacts in the database.">
          <div class="text-3xl font-semibold tracking-tight">{{ contactCount }}</div>
          <p class="text-muted-foreground mt-2 text-xs">Reflects real-time database state.</p>
        </FeatureCard>
      </div>

      <FeatureCard title="How It Works" description="usePoll under the hood.">
        <div class="space-y-3 text-sm">
          <p>
            <code class="bg-muted rounded px-1.5 py-0.5 text-xs">usePoll(5000)</code> sets up an
            interval that calls
            <code class="bg-muted rounded px-1.5 py-0.5 text-xs">router.reload()</code>
            every 5 seconds. The poll automatically pauses when the tab is not visible and resumes
            when you return.
          </p>
          <p>
            You can control polling with the returned
            <code class="bg-muted rounded px-1.5 py-0.5 text-xs">stop()</code> and
            <code class="bg-muted rounded px-1.5 py-0.5 text-xs">start()</code> functions.
          </p>
        </div>
      </FeatureCard>
    </div>
  </AppLayout>
</template>
