<script setup lang="ts">
import { Link } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import AppLayout from "@/js/layouts/AppLayout.vue";
import { useDemoRoute } from "@/js/lib/routes";

const breadcrumbs = [{ title: "Features" }, { title: "Navigation" }, { title: "View Transitions" }];

const linksRoute = useDemoRoute("features.navigation.links");
</script>

<template>
  <AppLayout title="View Transitions" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="View Transitions"
        description="Inertia supports the CSS View Transitions API. Add the view-transition attribute to Link components to animate page transitions using the browser's native capabilities."
      />

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard
          title="How It Works"
          description="The View Transitions API enables smooth animated transitions between DOM states."
        >
          <div class="space-y-3 text-sm">
            <p>
              When you add the
              <code class="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">view-transition</code>
              attribute to a Link component, Inertia wraps the page update inside
              <code class="rounded bg-muted px-1.5 py-0.5 font-mono text-xs"
                >document.startViewTransition()</code
              >.
            </p>
            <p>
              The browser captures a screenshot of the current page, applies the DOM update, then
              cross-fades between the old and new states. You can customize the animation with CSS.
            </p>
            <p class="text-muted-foreground">
              Note: View Transitions require a supporting browser (Chrome 111+, Edge 111+).
              Unsupported browsers skip the animation gracefully.
            </p>
          </div>
        </FeatureCard>

        <FeatureCard title="Usage" description="Add view-transition to any Link component.">
          <pre
            class="overflow-auto rounded-md bg-muted p-3 text-xs leading-relaxed"
          ><code>&lt;Link href="/page" view-transition&gt;
  Navigate with transition
&lt;/Link&gt;</code></pre>
          <p class="text-muted-foreground mt-3 text-xs">
            You can also pass it programmatically via
            <code class="rounded bg-muted px-1.5 py-0.5 font-mono text-xs"
              >router.visit(url, { viewTransition: true })</code
            >.
          </p>
        </FeatureCard>

        <FeatureCard
          title="Try It"
          description="Click the links below to navigate with and without view transitions."
        >
          <div class="flex flex-wrap gap-3">
            <Link
              :href="linksRoute.url"
              view-transition
              class="inline-flex h-9 items-center rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground hover:bg-primary/90"
            >
              Go to Links (with transition)
            </Link>
            <Link
              :href="linksRoute.url"
              class="inline-flex h-9 items-center rounded-md border bg-background px-4 text-sm font-medium shadow-xs hover:bg-accent hover:text-accent-foreground"
            >
              Go to Links (no transition)
            </Link>
          </div>
        </FeatureCard>

        <FeatureCard
          title="CSS Customization"
          description="Override the default cross-fade with custom CSS animations."
        >
          <pre
            class="overflow-auto rounded-md bg-muted p-3 text-xs leading-relaxed"
          ><code>::view-transition-old(root) {
  animation: slide-out 0.3s ease-in;
}
::view-transition-new(root) {
  animation: slide-in 0.3s ease-out;
}</code></pre>
        </FeatureCard>
      </div>
    </div>
  </AppLayout>
</template>
