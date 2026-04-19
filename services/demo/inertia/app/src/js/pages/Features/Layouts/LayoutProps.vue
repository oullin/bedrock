<script setup lang="ts">
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Badge } from "@/js/components/ui/badge";
import AppLayout from "@/js/layouts/AppLayout.vue";

const breadcrumbs = [{ title: "Features" }, { title: "Layouts" }, { title: "Layout Props" }];

const layoutTitle = "Layout Props";
const layoutBreadcrumbs = breadcrumbs;

const currentLayoutProps = {
  title: layoutTitle,
  breadcrumbs: layoutBreadcrumbs,
};
</script>

<template>
  <AppLayout :title="layoutTitle" :breadcrumbs="layoutBreadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Layout Props"
        description="Pass data from page components to their layout. The layout receives props like title and breadcrumbs from each page."
      />

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard
          title="Props Passed to Layout"
          description="Data this page sends to AppLayout."
        >
          <div class="space-y-2">
            <div class="rounded-md border p-3">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium">title</span>
                <Badge variant="secondary">{{ currentLayoutProps.title }}</Badge>
              </div>
            </div>
            <div class="rounded-md border p-3">
              <div class="mb-2 flex items-center justify-between">
                <span class="text-sm font-medium">breadcrumbs</span>
                <Badge variant="outline">Array[{{ currentLayoutProps.breadcrumbs.length }}]</Badge>
              </div>
              <div class="flex flex-wrap gap-1">
                <Badge
                  v-for="(crumb, index) in currentLayoutProps.breadcrumbs"
                  :key="index"
                  variant="secondary"
                  class="text-xs"
                >
                  {{ crumb.title }}
                </Badge>
              </div>
            </div>
          </div>
        </FeatureCard>

        <FeatureCard title="How It Works" description="Passing props to layouts.">
          <div class="space-y-3 text-sm">
            <div class="space-y-2">
              <p class="font-medium">Template wrapper approach</p>
              <pre class="bg-muted overflow-auto rounded p-3 text-xs leading-relaxed">
&lt;template&gt;
  &lt;AppLayout
    title="My Page"
    :breadcrumbs="crumbs"
  &gt;
    &lt;!-- page content --&gt;
  &lt;/AppLayout&gt;
&lt;/template&gt;</pre
              >
            </div>
            <p class="text-muted-foreground">
              Each page wraps its content in the layout component and passes page-specific data as
              props. The layout uses these props to render the header, breadcrumbs, and other shared
              UI.
            </p>
          </div>
        </FeatureCard>

        <FeatureCard
          title="Persistent Layout Props"
          description="How props work with defineOptions layout."
        >
          <div class="space-y-3 text-sm">
            <p>
              When using persistent layouts via
              <code class="bg-muted rounded px-1.5 py-0.5 text-xs">defineOptions</code>, the layout
              receives page props through the Inertia page object rather than direct Vue prop
              binding.
            </p>
            <pre class="bg-muted overflow-auto rounded p-3 text-xs leading-relaxed">
// In the layout component
const page = usePage()
const title = computed(
  () => page.props.title
)</pre
            >
          </div>
        </FeatureCard>

        <FeatureCard title="Common Layout Props" description="Typical data passed to layouts.">
          <ul class="space-y-2 text-sm">
            <li class="flex items-start gap-2">
              <Badge variant="secondary" class="shrink-0 text-xs">title</Badge>
              <span>Page title shown in header and browser tab.</span>
            </li>
            <li class="flex items-start gap-2">
              <Badge variant="secondary" class="shrink-0 text-xs">breadcrumbs</Badge>
              <span>Navigation path for the current page.</span>
            </li>
            <li class="flex items-start gap-2">
              <Badge variant="secondary" class="shrink-0 text-xs">actions</Badge>
              <span>Page-specific action buttons for the header.</span>
            </li>
            <li class="flex items-start gap-2">
              <Badge variant="secondary" class="shrink-0 text-xs">sidebar</Badge>
              <span>Whether to show/hide sidebar sections.</span>
            </li>
          </ul>
        </FeatureCard>
      </div>
    </div>
  </AppLayout>
</template>
