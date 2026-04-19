<script setup lang="ts">
import { computed } from "vue";
import { usePage } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Badge } from "@/js/components/ui/badge";
import AppLayout from "@/js/layouts/AppLayout.vue";
import type { SharedPageProps } from "@/js/types";

const page = usePage<SharedPageProps>();

const breadcrumbs = [{ title: "Features" }, { title: "State" }, { title: "Shared Props" }];

const propEntries = computed(() => {
  return Object.entries(page.props).map(([key, value]) => ({
    key,
    value,
    type: value === null ? "null" : Array.isArray(value) ? "array" : typeof value,
  }));
});
</script>

<template>
  <AppLayout title="Shared Props" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Shared Props"
        description="Shared props are data that is available on every page without explicitly passing it. Access them via usePage().props."
      />

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard
          title="All Shared Props"
          description="Every prop available via usePage().props on this page."
        >
          <div class="space-y-2">
            <div v-for="entry in propEntries" :key="entry.key" class="rounded-md border p-3">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium">{{ entry.key }}</span>
                <Badge variant="outline" class="text-xs">{{ entry.type }}</Badge>
              </div>
              <pre class="bg-muted mt-2 overflow-auto rounded p-2 text-xs">{{
                JSON.stringify(entry.value, null, 2)
              }}</pre>
            </div>
          </div>
        </FeatureCard>

        <div class="space-y-6">
          <FeatureCard title="Auth User" description="The authenticated user from shared props.">
            <template v-if="page.props.auth?.user">
              <div class="space-y-2">
                <div class="flex items-center gap-3">
                  <div
                    class="bg-primary/10 text-primary flex size-10 items-center justify-center rounded-full font-medium"
                  >
                    {{ page.props.auth.user.name?.[0] ?? "?" }}
                  </div>
                  <div>
                    <p class="font-medium">{{ page.props.auth.user.name }}</p>
                    <p class="text-muted-foreground text-sm">{{ page.props.auth.user.email }}</p>
                  </div>
                </div>
              </div>
            </template>
            <p v-else class="text-muted-foreground text-sm">No authenticated user.</p>
          </FeatureCard>

          <FeatureCard title="Routes" description="Shared route definitions.">
            <div v-if="page.props.routes" class="max-h-64 space-y-1 overflow-auto">
              <div
                v-for="(url, name) in page.props.routes"
                :key="name"
                class="flex items-center justify-between gap-2 rounded border px-2 py-1.5 text-xs"
              >
                <span class="font-mono font-medium">{{ name }}</span>
                <span class="text-muted-foreground truncate">{{ url }}</span>
              </div>
            </div>
            <p v-else class="text-muted-foreground text-sm">No routes shared.</p>
          </FeatureCard>

          <FeatureCard title="How It Works">
            <p class="text-sm">
              Shared props are defined on the server and automatically merged into every Inertia
              response. Common uses include the authenticated user, flash messages, app
              configuration, and route definitions. Access them anywhere with
              <code class="bg-muted rounded px-1.5 py-0.5 text-xs">usePage().props</code>.
            </p>
          </FeatureCard>
        </div>
      </div>
    </div>
  </AppLayout>
</template>
