<script setup lang="ts">
import { Head } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Badge } from "@/js/components/ui/badge";
import AppLayout from "@/js/layouts/AppLayout.vue";

const breadcrumbs = [{ title: "Features" }, { title: "Layouts" }, { title: "Head" }];
</script>

<template>
  <AppLayout title="Head Component" :breadcrumbs="breadcrumbs">
    <Head>
      <title>Custom Page Title - Inertia Demo</title>
      <meta
        name="description"
        content="Demonstrating the Inertia Head component for managing document head."
      />
      <meta property="og:title" content="Head Component Demo" />
      <meta
        property="og:description"
        content="Learn how to manage the document head with Inertia.js"
      />
    </Head>

    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Head Component"
        description="Use the Head component to set the page title, meta tags, and other elements in the document head."
      />

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard
          title="Current Head Tags"
          description="Tags set by this page's Head component."
        >
          <div class="space-y-2">
            <div class="rounded-md border p-3">
              <div class="flex items-center gap-2">
                <Badge>title</Badge>
                <span class="text-sm">Custom Page Title - Inertia Demo</span>
              </div>
            </div>
            <div class="rounded-md border p-3">
              <div class="flex items-center gap-2">
                <Badge variant="secondary">meta</Badge>
                <span class="text-sm">description: Demonstrating the Inertia Head component</span>
              </div>
            </div>
            <div class="rounded-md border p-3">
              <div class="flex items-center gap-2">
                <Badge variant="secondary">meta</Badge>
                <span class="text-sm">og:title: Head Component Demo</span>
              </div>
            </div>
            <div class="rounded-md border p-3">
              <div class="flex items-center gap-2">
                <Badge variant="secondary">meta</Badge>
                <span class="text-sm">og:description: Learn how to manage the document head</span>
              </div>
            </div>
          </div>
          <p class="text-muted-foreground mt-3 text-xs">
            Check your browser tab to see the custom title. Open DevTools and inspect the
            &lt;head&gt; element to see the meta tags.
          </p>
        </FeatureCard>

        <FeatureCard title="Usage" description="How to use the Head component.">
          <div class="space-y-4">
            <div class="space-y-2">
              <p class="text-sm font-medium">Title shorthand</p>
              <pre class="bg-muted overflow-auto rounded p-3 text-xs">
&lt;Head title="My Page" /&gt;</pre
              >
            </div>
            <div class="space-y-2">
              <p class="text-sm font-medium">Full control</p>
              <pre class="bg-muted overflow-auto rounded p-3 text-xs leading-relaxed">
&lt;Head&gt;
  &lt;title&gt;My Page&lt;/title&gt;
  &lt;meta name="description"
    content="Page description" /&gt;
  &lt;meta property="og:title"
    content="Social title" /&gt;
&lt;/Head&gt;</pre
              >
            </div>
            <p class="text-muted-foreground text-sm">
              The Head component manages document head tags across page navigations. Tags from the
              previous page are automatically cleaned up.
            </p>
          </div>
        </FeatureCard>

        <FeatureCard title="Title Callback" description="Dynamic title formatting.">
          <div class="space-y-3 text-sm">
            <p>
              You can configure a global title callback in your app setup to format all page titles
              consistently:
            </p>
            <pre class="bg-muted overflow-auto rounded p-3 text-xs leading-relaxed">
createInertiaApp({
  title: (title) =>
    `${title} - My App`,
})</pre
            >
            <p class="text-muted-foreground">
              This appends " - My App" to every page title automatically.
            </p>
          </div>
        </FeatureCard>

        <FeatureCard title="Head Priorities" description="How conflicting head tags are resolved.">
          <ul class="space-y-2 text-sm">
            <li class="flex items-start gap-2">
              <span class="text-primary font-medium">Page wins</span> - Page-level Head tags
              override layout-level ones.
            </li>
            <li class="flex items-start gap-2">
              <span class="text-primary font-medium">Last wins</span> - If multiple Head components
              exist, the last one rendered takes priority.
            </li>
            <li class="flex items-start gap-2">
              <span class="text-primary font-medium">Auto cleanup</span> - Tags are removed when
              navigating away from the page.
            </li>
          </ul>
        </FeatureCard>
      </div>
    </div>
  </AppLayout>
</template>
