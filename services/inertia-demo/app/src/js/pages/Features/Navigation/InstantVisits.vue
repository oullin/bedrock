<script setup lang="ts">
import { Link } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import AppLayout from "@/js/layouts/AppLayout.vue";
import { useDemoRoute } from "@/js/lib/routes";

withDefaults(defineProps<{ sourceTimestamp?: string }>(), {
  sourceTimestamp: "",
});

const breadcrumbs = [{ title: "Features" }, { title: "Navigation" }, { title: "Instant Visits" }];

const targetRoute = useDemoRoute("features.navigation.instant-visit-target");
</script>

<template>
  <AppLayout title="Instant Visits" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Instant Visits"
        description='Instant visits preload the target page as soon as the link renders, so navigation feels instantaneous. With prefetch="mount", the request fires immediately when the component mounts — the response is often cached before the user even clicks.'
      />

      <p class="text-muted-foreground text-sm">
        Source timestamp:
        <code class="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">{{ sourceTimestamp }}</code>
      </p>

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard
          title="Instant Visit Link"
          description="This link prefetches the request on mount so it's ready before you click."
        >
          <div class="flex flex-wrap gap-3">
            <Link
              :href="targetRoute.url"
              prefetch="mount"
              class="inline-flex h-9 items-center rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground hover:bg-primary/90"
            >
              Visit target page (instant)
            </Link>
          </div>
          <p class="text-muted-foreground mt-2 text-xs">
            The request starts on mount/hover, so by the time you click the response is often
            already cached.
          </p>
        </FeatureCard>

        <FeatureCard
          title="Normal Visit Link"
          description="A standard link for comparison. The request starts after the click event."
        >
          <div class="flex flex-wrap gap-3">
            <Link
              :href="targetRoute.url"
              class="inline-flex h-9 items-center rounded-md border bg-background px-4 text-sm font-medium shadow-xs hover:bg-accent hover:text-accent-foreground"
            >
              Visit target page (normal)
            </Link>
          </div>
          <p class="text-muted-foreground mt-2 text-xs">
            This is a regular Inertia link with no prefetching.
          </p>
        </FeatureCard>

        <FeatureCard
          title="How It Works"
          description="Instant visits leverage Inertia's prefetching to preload page data."
        >
          <div class="space-y-2 text-sm">
            <p>
              When
              <code class="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">prefetch="mount"</code>
              is set, Inertia starts loading the target page as soon as the link component mounts.
            </p>
            <p>
              When the user clicks, the cached response is used immediately if available, making the
              navigation appear instantaneous.
            </p>
          </div>
        </FeatureCard>

        <FeatureCard title="Usage" description="Add the prefetch attribute to Link components.">
          <pre
            class="overflow-auto rounded-md bg-muted p-3 text-xs leading-relaxed"
          ><code>&lt;Link href="/target" prefetch="mount"&gt;
  Instant visit
&lt;/Link&gt;

&lt;!-- Other prefetch strategies --&gt;
&lt;Link href="/target" prefetch="hover"&gt;
  Prefetch on hover
&lt;/Link&gt;</code></pre>
        </FeatureCard>
      </div>
    </div>
  </AppLayout>
</template>
