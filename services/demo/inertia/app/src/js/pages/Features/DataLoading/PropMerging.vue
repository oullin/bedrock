<script setup lang="ts">
import { router } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Button } from "@/js/components/ui/button";
import { Badge } from "@/js/components/ui/badge";
import AppLayout from "@/js/layouts/AppLayout.vue";

withDefaults(
  defineProps<{
    notifications?: Array<{ id?: number; title?: string; body?: string }>;
    activities?: Array<{ id?: number; type?: string; description?: string }>;
    timestamp?: string;
  }>(),
  {
    notifications: () => [],
    activities: () => [],
    timestamp: "",
  },
);

const breadcrumbs = [{ title: "Features" }, { title: "Data Loading" }, { title: "Prop Merging" }];

function loadMore() {
  router.reload({ only: ["notifications", "activities", "timestamp"] });
}

function resetAll() {
  router.reload({ reset: ["notifications", "activities"] });
}
</script>

<template>
  <AppLayout title="Prop Merging" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Prop Merging"
        description="When props are marked as mergeable on the server, subsequent reloads append new data rather than replacing. This is useful for infinite scroll and activity feeds."
      />

      <div class="flex items-center gap-3">
        <Button @click="loadMore">Load More</Button>
        <Button variant="outline" @click="resetAll">Reset</Button>
        <Badge variant="secondary">Last updated: {{ timestamp }}</Badge>
      </div>

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard
          title="Notifications"
          :description="`${notifications.length} notification(s) accumulated`"
        >
          <div class="space-y-2">
            <div
              v-for="(notification, index) in notifications"
              :key="notification.id ?? index"
              class="flex items-start gap-3 rounded-md border p-3"
            >
              <div
                class="bg-primary/10 flex size-7 shrink-0 items-center justify-center rounded-full text-xs font-medium"
              >
                {{ index + 1 }}
              </div>
              <div class="min-w-0 flex-1">
                <p class="text-sm font-medium">{{ notification.title ?? notification }}</p>
                <p v-if="notification.body" class="text-muted-foreground text-xs">
                  {{ notification.body }}
                </p>
              </div>
            </div>
            <p v-if="notifications.length === 0" class="text-muted-foreground text-sm">
              No notifications yet. Click "Load More" to fetch some.
            </p>
          </div>
        </FeatureCard>

        <FeatureCard
          title="Activities"
          :description="`${activities.length} ${activities.length === 1 ? 'activity' : 'activities'} accumulated`"
        >
          <div class="space-y-2">
            <div
              v-for="(activity, index) in activities"
              :key="activity.id ?? index"
              class="flex items-center gap-3 rounded-md border p-3"
            >
              <Badge variant="outline" class="shrink-0">{{ activity.type ?? "event" }}</Badge>
              <p class="text-sm">{{ activity.description ?? activity }}</p>
            </div>
            <p v-if="activities.length === 0" class="text-muted-foreground text-sm">
              No activities yet. Click "Load More" to fetch some.
            </p>
          </div>
        </FeatureCard>
      </div>

      <FeatureCard title="How It Works">
        <p class="text-sm">
          When the server marks a prop with
          <code class="bg-muted rounded px-1.5 py-0.5 text-xs">merge</code>, Inertia appends new
          items to the existing array rather than replacing it. Use
          <code class="bg-muted rounded px-1.5 py-0.5 text-xs">reset</code> in the reload options to
          clear the accumulated data and start fresh.
        </p>
      </FeatureCard>
    </div>
  </AppLayout>
</template>
