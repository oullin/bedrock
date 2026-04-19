<script setup lang="ts">
import { router, usePage } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Button } from "@/js/components/ui/button";
import AppLayout from "@/js/layouts/AppLayout.vue";
import type { SharedPageProps } from "@/js/types";

withDefaults(defineProps<{ timestamp?: string }>(), {
  timestamp: "",
});

const page = usePage<SharedPageProps>();

const breadcrumbs = [{ title: "Features" }, { title: "Navigation" }, { title: "URL Fragments" }];

function actionUrl(suffix: string): string {
  const url = new URL(page.url, window.location.origin);
  url.pathname = url.pathname.replace(/\/?$/, "/") + suffix;
  return url.pathname;
}

function redirectWithHash() {
  router.post(actionUrl("redirect-with-hash"));
}

function preserveFragment() {
  router.post(actionUrl("preserve-fragment"));
}
</script>

<template>
  <AppLayout title="URL Fragments" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="URL Fragments"
        description="Inertia supports URL fragments (hash anchors). You can navigate to specific sections, redirect with a hash, or preserve the current fragment across navigations."
      />

      <p class="text-muted-foreground text-sm">
        Server timestamp:
        <code class="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">{{ timestamp }}</code>
      </p>

      <FeatureCard
        title="Fragment Actions"
        description="Test server-side redirect behaviors with URL fragments."
      >
        <div class="flex flex-wrap gap-3">
          <Button @click="redirectWithHash">Redirect with #section-2</Button>
          <Button variant="outline" @click="preserveFragment">Redirect (preserve fragment)</Button>
        </div>
        <p class="text-muted-foreground mt-2 text-xs">
          "Redirect with hash" sends a POST and the server redirects to this page with #section-2 in
          the URL. "Preserve fragment" redirects without a hash.
        </p>
      </FeatureCard>

      <FeatureCard
        title="Hash Links"
        description="Click these links to jump to different sections on this page."
      >
        <div class="flex flex-wrap gap-3">
          <a
            href="#section-1"
            class="inline-flex h-9 items-center rounded-md border bg-background px-4 text-sm font-medium shadow-xs hover:bg-accent hover:text-accent-foreground"
          >
            Jump to Section 1
          </a>
          <a
            href="#section-2"
            class="inline-flex h-9 items-center rounded-md border bg-background px-4 text-sm font-medium shadow-xs hover:bg-accent hover:text-accent-foreground"
          >
            Jump to Section 2
          </a>
          <a
            href="#section-3"
            class="inline-flex h-9 items-center rounded-md border bg-background px-4 text-sm font-medium shadow-xs hover:bg-accent hover:text-accent-foreground"
          >
            Jump to Section 3
          </a>
        </div>
      </FeatureCard>

      <div id="section-1" class="scroll-mt-20">
        <FeatureCard title="Section 1" description="This section has the anchor id 'section-1'.">
          <div class="space-y-2 text-sm">
            <p>
              URL fragments let you link directly to specific parts of a page. When a user visits a
              URL with a hash, the browser scrolls to the element with the matching id.
            </p>
            <p>
              Inertia preserves this behavior during client-side navigations so hash-based scrolling
              works the same as traditional server-rendered pages.
            </p>
          </div>
        </FeatureCard>
      </div>

      <div id="section-2" class="scroll-mt-20">
        <FeatureCard title="Section 2" description="This section has the anchor id 'section-2'.">
          <div class="space-y-2 text-sm">
            <p>
              When the server redirects to a URL that includes a hash (e.g.
              <code class="rounded bg-muted px-1.5 py-0.5 font-mono text-xs"
                >/url-fragments#section-2</code
              >), Inertia will scroll to this element after the page renders.
            </p>
            <p>Try the "Redirect with #section-2" button above to see this in action.</p>
          </div>
        </FeatureCard>
      </div>

      <div id="section-3" class="scroll-mt-20">
        <FeatureCard title="Section 3" description="This section has the anchor id 'section-3'.">
          <div class="space-y-2 text-sm">
            <p>
              URL fragments are especially useful for long pages with multiple sections,
              documentation with subsections, or linking to specific content in lists.
            </p>
            <p>
              Unlike full page visits, fragment-only navigation stays on the same page without
              making a server request.
            </p>
          </div>
        </FeatureCard>
      </div>
    </div>
  </AppLayout>
</template>
