<script setup lang="ts">
import { router } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Button } from "@/js/components/ui/button";
import { Badge } from "@/js/components/ui/badge";
import AppLayout from "@/js/layouts/AppLayout.vue";

withDefaults(
  defineProps<{
    users?: Array<{ id?: number; name?: string; email?: string }>;
    stats?: Record<string, unknown>;
    timestamp?: string;
    randomNumber?: number;
  }>(),
  {
    users: () => [],
    stats: () => ({}),
    timestamp: "",
    randomNumber: 0,
  },
);

const breadcrumbs = [
  { title: "Features" },
  { title: "Data Loading" },
  { title: "Partial Reloads" },
];

function reloadUsers() {
  router.reload({ only: ["users"] });
}

function reloadStats() {
  router.reload({ only: ["stats"] });
}

function reloadTimestamp() {
  router.reload({ only: ["timestamp", "randomNumber"] });
}

function reloadAll() {
  router.reload();
}
</script>

<template>
  <AppLayout title="Partial Reloads" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Partial Reloads"
        description="Reload only specific props without re-fetching the entire page data. This reduces server load and speeds up updates."
      />

      <div class="flex flex-wrap gap-2">
        <Button @click="reloadUsers">Reload Users Only</Button>
        <Button variant="secondary" @click="reloadStats">Reload Stats Only</Button>
        <Button variant="outline" @click="reloadTimestamp">Reload Timestamp</Button>
        <Button variant="ghost" @click="reloadAll">Reload Everything</Button>
      </div>

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard title="Users" description="Reloaded independently from other props.">
          <div class="space-y-2">
            <div
              v-for="user in users"
              :key="user.id ?? user.email"
              class="flex items-center gap-3 rounded-md border p-3"
            >
              <div
                class="bg-primary/10 text-primary flex size-8 items-center justify-center rounded-full text-xs font-medium"
              >
                {{ (user.name ?? "U")[0] }}
              </div>
              <div class="min-w-0 flex-1">
                <p class="text-sm font-medium">{{ user.name }}</p>
                <p class="text-muted-foreground truncate text-xs">{{ user.email }}</p>
              </div>
            </div>
            <p v-if="users.length === 0" class="text-muted-foreground text-sm">No users loaded.</p>
          </div>
        </FeatureCard>

        <div class="space-y-6">
          <FeatureCard
            title="Stats"
            description="Separate prop group that can be refreshed on its own."
          >
            <div class="grid grid-cols-2 gap-3">
              <div
                v-for="(value, key) in stats"
                :key="key"
                class="rounded-md border p-3 text-center"
              >
                <div class="text-lg font-semibold">{{ value }}</div>
                <div class="text-muted-foreground text-xs capitalize">{{ key }}</div>
              </div>
            </div>
          </FeatureCard>

          <FeatureCard title="Metadata" description="Timestamp and random number reload together.">
            <div class="space-y-3">
              <div class="flex items-center justify-between rounded-md border p-3">
                <span class="text-sm font-medium">Timestamp</span>
                <Badge variant="secondary">{{ timestamp }}</Badge>
              </div>
              <div class="flex items-center justify-between rounded-md border p-3">
                <span class="text-sm font-medium">Random Number</span>
                <Badge variant="outline">{{ randomNumber }}</Badge>
              </div>
            </div>
          </FeatureCard>
        </div>
      </div>
    </div>
  </AppLayout>
</template>
