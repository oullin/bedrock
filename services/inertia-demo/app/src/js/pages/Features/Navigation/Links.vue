<script setup lang="ts">
import { Link, usePage } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import AppLayout from "@/js/layouts/AppLayout.vue";
import type { SharedPageProps } from "@/js/types";

withDefaults(defineProps<{ timestamp?: string }>(), {
  timestamp: "",
});

const page = usePage<SharedPageProps>();
const currentUrl = page.url;

const breadcrumbs = [{ title: "Features" }, { title: "Navigation" }, { title: "Links" }];
</script>

<template>
  <AppLayout title="Links" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Links"
        description="The Link component creates anchor tags that make Inertia requests instead of full page loads. It supports different HTTP methods, custom headers, and scroll preservation."
      />

      <p class="text-muted-foreground text-sm">
        Server timestamp:
        <code class="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">{{ timestamp }}</code>
      </p>

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard
          title="GET Request"
          description="Standard navigation link that reloads this page."
        >
          <div class="flex flex-wrap gap-3">
            <Link
              :href="currentUrl"
              class="inline-flex h-9 items-center rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground hover:bg-primary/90"
            >
              GET (reload page)
            </Link>
          </div>
        </FeatureCard>

        <FeatureCard
          title="POST Request"
          description="Link that issues a POST request to the server."
        >
          <div class="flex flex-wrap gap-3">
            <Link
              :href="currentUrl"
              method="post"
              as="button"
              class="inline-flex h-9 items-center rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground hover:bg-primary/90"
            >
              POST action
            </Link>
          </div>
        </FeatureCard>

        <FeatureCard
          title="PUT / DELETE"
          description="Links using PUT and DELETE methods rendered as buttons."
        >
          <div class="flex flex-wrap gap-3">
            <Link
              :href="currentUrl"
              method="put"
              as="button"
              class="inline-flex h-9 items-center rounded-md border bg-background px-4 text-sm font-medium shadow-xs hover:bg-accent hover:text-accent-foreground"
            >
              PUT action
            </Link>
            <Link
              :href="currentUrl"
              method="delete"
              as="button"
              class="inline-flex h-9 items-center rounded-md bg-destructive px-4 text-sm font-medium text-white hover:bg-destructive/90"
            >
              DELETE action
            </Link>
          </div>
        </FeatureCard>

        <FeatureCard
          title="Custom Headers"
          description="Link that sends a custom X-Demo-Header with the request."
        >
          <div class="flex flex-wrap gap-3">
            <Link
              :href="currentUrl"
              :headers="{ 'X-Demo-Header': 'hello-from-link' }"
              class="inline-flex h-9 items-center rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground hover:bg-primary/90"
            >
              GET with custom header
            </Link>
          </div>
          <p class="text-muted-foreground mt-2 text-xs">
            Open DevTools Network tab to see the X-Demo-Header on the request.
          </p>
        </FeatureCard>

        <FeatureCard
          title="Preserve Scroll"
          description="Link with preserve-scroll keeps the current scroll position after navigation."
        >
          <div class="flex flex-wrap gap-3">
            <Link
              :href="currentUrl"
              preserve-scroll
              class="inline-flex h-9 items-center rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground hover:bg-primary/90"
            >
              Reload (preserve scroll)
            </Link>
            <Link
              :href="currentUrl"
              class="inline-flex h-9 items-center rounded-md border bg-background px-4 text-sm font-medium shadow-xs hover:bg-accent hover:text-accent-foreground"
            >
              Reload (reset scroll)
            </Link>
          </div>
        </FeatureCard>

        <FeatureCard
          title="Replace History"
          description="Link with replace removes the current URL from the browser history stack."
        >
          <div class="flex flex-wrap gap-3">
            <Link
              :href="currentUrl"
              replace
              class="inline-flex h-9 items-center rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground hover:bg-primary/90"
            >
              Reload (replace history)
            </Link>
          </div>
          <p class="text-muted-foreground mt-2 text-xs">
            After clicking, pressing the browser back button will skip this page entry.
          </p>
        </FeatureCard>
      </div>
    </div>
  </AppLayout>
</template>
