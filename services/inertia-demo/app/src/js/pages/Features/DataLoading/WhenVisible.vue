<script setup lang="ts">
import { WhenVisible } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Skeleton } from "@/js/components/ui/skeleton";
import { Badge } from "@/js/components/ui/badge";
import AppLayout from "@/js/layouts/AppLayout.vue";

withDefaults(
  defineProps<{
    section1?: Record<string, unknown>;
    section2?: Record<string, unknown>;
    section3?: Record<string, unknown>;
  }>(),
  {
    section1: undefined,
    section2: undefined,
    section3: undefined,
  },
);

const breadcrumbs = [{ title: "Features" }, { title: "Data Loading" }, { title: "When Visible" }];
</script>

<template>
  <AppLayout title="When Visible" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="When Visible"
        description="Load deferred data only when the component scrolls into the viewport. Scroll down to trigger each section."
      />

      <FeatureCard
        title="Introduction"
        description="The sections below will load their data lazily as you scroll down."
      >
        <p class="text-muted-foreground text-sm">
          WhenVisible uses an IntersectionObserver to detect when an element enters the viewport,
          then triggers a partial reload to fetch the associated data. This is ideal for long pages
          where not all data is needed immediately.
        </p>
      </FeatureCard>

      <div class="h-32" />

      <WhenVisible data="section1">
        <template #fallback>
          <FeatureCard title="Section 1" description="Loading...">
            <div class="space-y-3">
              <Skeleton class="h-6 w-full" />
              <Skeleton class="h-6 w-3/4" />
              <Skeleton class="h-6 w-1/2" />
            </div>
          </FeatureCard>
        </template>
        <FeatureCard title="Section 1" description="Loaded when scrolled into view.">
          <div class="space-y-2">
            <div class="flex items-center gap-2">
              <Badge>Loaded</Badge>
              <span class="text-sm">{{ section1?.title ?? "Section 1 data" }}</span>
            </div>
            <p class="text-muted-foreground text-sm">
              {{ section1?.content ?? "This content was fetched when the section became visible." }}
            </p>
          </div>
        </FeatureCard>
      </WhenVisible>

      <div class="h-32" />

      <WhenVisible data="section2">
        <template #fallback>
          <FeatureCard title="Section 2" description="Loading...">
            <div class="space-y-3">
              <Skeleton class="h-6 w-full" />
              <Skeleton class="h-6 w-2/3" />
              <Skeleton class="h-6 w-1/3" />
            </div>
          </FeatureCard>
        </template>
        <FeatureCard title="Section 2" description="Also loaded on visibility.">
          <div class="space-y-2">
            <div class="flex items-center gap-2">
              <Badge variant="secondary">Loaded</Badge>
              <span class="text-sm">{{ section2?.title ?? "Section 2 data" }}</span>
            </div>
            <p class="text-muted-foreground text-sm">
              {{ section2?.content ?? "This section loaded independently when you scrolled here." }}
            </p>
          </div>
        </FeatureCard>
      </WhenVisible>

      <div class="h-32" />

      <WhenVisible data="section3">
        <template #fallback>
          <FeatureCard title="Section 3" description="Loading...">
            <div class="space-y-3">
              <Skeleton class="h-6 w-full" />
              <Skeleton class="h-6 w-5/6" />
              <Skeleton class="h-6 w-2/3" />
            </div>
          </FeatureCard>
        </template>
        <FeatureCard title="Section 3" description="The final lazy-loaded section.">
          <div class="space-y-2">
            <div class="flex items-center gap-2">
              <Badge variant="outline">Loaded</Badge>
              <span class="text-sm">{{ section3?.title ?? "Section 3 data" }}</span>
            </div>
            <p class="text-muted-foreground text-sm">
              {{
                section3?.content ??
                "Each section makes its own partial reload request to the server."
              }}
            </p>
          </div>
        </FeatureCard>
      </WhenVisible>
    </div>
  </AppLayout>
</template>
