<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { Link, router } from "@inertiajs/vue3";
import { Badge } from "@/js/components/ui/badge";
import { Button } from "@/js/components/ui/button";
import { Input } from "@/js/components/ui/input";
import AppLayout from "@/js/layouts/AppLayout.vue";
import { organizationRoutes } from "@/js/lib/routes";
import type { OffsetPaginated, Organization } from "@/js/types";

const props = withDefaults(
  defineProps<{
    organizations?: OffsetPaginated<Organization>;
    filters?: { search: string };
  }>(),
  {
    organizations: () => ({ data: [], total: 0, per_page: 20, current_page: 1, last_page: 1 }),
    filters: () => ({ search: "" }),
  },
);

const breadcrumbs = [
  { title: "CRM" },
  { title: "Organizations", href: organizationRoutes.index().url },
];

const search = ref(props.filters.search ?? "");
let timeout: ReturnType<typeof setTimeout> | undefined;

watch(search, (value) => {
  clearTimeout(timeout);
  timeout = setTimeout(() => {
    router.visit(organizationRoutes.index().url, {
      data: { search: value || undefined },
      only: ["organizations", "filters"],
      preserveState: true,
      preserveScroll: true,
      reset: ["organizations"],
    });
  }, 300);
});

const orgList = computed(() => props.organizations.data ?? []);

const pages = computed((): (number | "...")[] => {
  const last = props.organizations.last_page;
  const current = props.organizations.current_page;

  if (last <= 7) {
    return Array.from({ length: last }, (_, i) => i + 1);
  }

  const items = new Set<number>([1, last]);
  for (let i = current - 1; i <= current + 1; i++) {
    if (i >= 1 && i <= last) items.add(i);
  }

  const sorted = [...items].sort((a, b) => a - b);
  const result: (number | "...")[] = [];

  for (let i = 0; i < sorted.length; i++) {
    if (i > 0 && sorted[i] - sorted[i - 1] > 1) {
      result.push("...");
    }
    result.push(sorted[i]);
  }

  return result;
});

function goToPage(page: number) {
  router.visit(organizationRoutes.index().url, {
    data: {
      search: search.value || undefined,
      page: page > 1 ? page : undefined,
    },
    only: ["organizations", "filters"],
    preserveState: true,
    preserveScroll: true,
  });
}
</script>

<template>
  <AppLayout title="Organizations" :breadcrumbs="breadcrumbs">
    <div class="flex h-full flex-1 flex-col gap-4 p-4">
      <h2 class="text-2xl font-semibold tracking-tight">Organizations</h2>

      <div class="relative max-w-sm">
        <Input v-model="search" placeholder="Search organizations..." />
      </div>

      <div class="space-y-2">
        <Link
          v-for="organization in orgList"
          :key="organization.id"
          :href="organizationRoutes.show(organization.id).url"
          prefetch="hover"
          class="flex items-center gap-4 rounded-lg bg-muted/30 p-4 hover:bg-muted/50"
        >
          <div
            class="bg-primary/10 text-primary flex size-10 shrink-0 items-center justify-center rounded-full text-sm font-medium"
          >
            {{ organization.name?.[0] ?? "?" }}
          </div>
          <div class="flex-1">
            <span class="font-medium">{{ organization.name }}</span>
          </div>
          <Badge variant="outline">
            {{ organization.contacts_count }}
            {{ organization.contacts_count === 1 ? "contact" : "contacts" }}
          </Badge>
        </Link>

        <div v-if="orgList.length === 0" class="text-muted-foreground py-12 text-center">
          No organizations found.
        </div>
      </div>

      <div v-if="organizations.last_page > 1" class="flex items-center justify-center gap-1 py-4">
        <template v-for="(page, index) in pages" :key="index">
          <span v-if="page === '...'" class="text-muted-foreground px-2 text-sm">...</span>
          <Button
            v-else
            variant="outline"
            size="sm"
            :class="{ 'bg-accent font-bold': page === organizations.current_page }"
            @click="goToPage(page)"
          >
            {{ page }}
          </Button>
        </template>
      </div>
    </div>
  </AppLayout>
</template>
