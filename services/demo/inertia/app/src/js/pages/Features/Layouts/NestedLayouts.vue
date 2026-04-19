<script setup lang="ts">
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Badge } from "@/js/components/ui/badge";
import AppLayout from "@/js/layouts/AppLayout.vue";
import SectionLayout from "@/js/layouts/SectionLayout.vue";

defineOptions({
  layout: [AppLayout, SectionLayout],
});
</script>

<template>
  <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
    <FeatureHeader
      title="Nested Layouts"
      description="Compose multiple layouts in a hierarchy. Inner layouts wrap page content while outer layouts provide the application shell."
    />

    <div class="grid gap-6 lg:grid-cols-2">
      <FeatureCard title="Layout Hierarchy" description="How nested layouts compose together.">
        <div class="space-y-3">
          <div class="rounded-md border-2 border-dashed p-4">
            <div class="flex items-center gap-2">
              <Badge>Level 1</Badge>
              <span class="text-sm font-medium">AppLayout</span>
            </div>
            <p class="text-muted-foreground mt-1 text-xs">
              Provides sidebar, header, and flash messages.
            </p>

            <div class="mt-3 rounded-md border-2 border-dashed p-4">
              <div class="flex items-center gap-2">
                <Badge variant="secondary">Level 2</Badge>
                <span class="text-sm font-medium">Section Layout</span>
              </div>
              <p class="text-muted-foreground mt-1 text-xs">
                Adds section-specific navigation or context.
              </p>

              <div class="mt-3 rounded-md border-2 border-dashed p-4">
                <div class="flex items-center gap-2">
                  <Badge variant="outline">Level 3</Badge>
                  <span class="text-sm font-medium">Page Content</span>
                </div>
                <p class="text-muted-foreground mt-1 text-xs">The actual page component content.</p>
              </div>
            </div>
          </div>
        </div>
      </FeatureCard>

      <FeatureCard
        title="Implementation"
        description="Using defineOptions with an array of layouts."
      >
        <div class="space-y-3 text-sm">
          <pre class="bg-muted overflow-auto rounded p-4 text-xs leading-relaxed">
import AppLayout from '@/layouts/AppLayout.vue'
import SectionLayout from '@/layouts/SectionLayout.vue'

defineOptions({
  layout: [AppLayout, SectionLayout],
})</pre
          >
          <p class="text-muted-foreground">
            When you pass an array of layouts to
            <code class="bg-muted rounded px-1.5 py-0.5 text-xs">defineOptions</code>, Inertia nests
            them from left to right. The first layout is the outermost wrapper.
          </p>
        </div>
      </FeatureCard>

      <FeatureCard title="Current Layout Stack" description="This page's layout composition.">
        <div class="space-y-2">
          <div class="flex items-center gap-3 rounded-md border p-3">
            <Badge>Active</Badge>
            <span class="text-sm">AppLayout (provides sidebar + header)</span>
          </div>
          <div class="flex items-center gap-3 rounded-md border p-3">
            <Badge variant="secondary">Active</Badge>
            <span class="text-sm">SectionLayout (section-level wrapper)</span>
          </div>
          <div class="flex items-center gap-3 rounded-md border p-3">
            <Badge variant="outline">Content</Badge>
            <span class="text-sm">NestedLayouts page (this content)</span>
          </div>
        </div>
      </FeatureCard>

      <FeatureCard title="Use Cases" description="When to use nested layouts.">
        <ul class="space-y-2 text-sm">
          <li class="flex items-start gap-2">
            <span class="text-primary font-medium">Admin sections</span> - Add an admin sub-nav
            inside the main app layout.
          </li>
          <li class="flex items-start gap-2">
            <span class="text-primary font-medium">Settings pages</span> - Wrap settings with a
            tabbed sidebar.
          </li>
          <li class="flex items-start gap-2">
            <span class="text-primary font-medium">Multi-step forms</span> - Shared progress bar
            across form steps.
          </li>
          <li class="flex items-start gap-2">
            <span class="text-primary font-medium">Dashboard areas</span> - Different dashboard
            sections with shared context.
          </li>
        </ul>
      </FeatureCard>
    </div>
  </div>
</template>
