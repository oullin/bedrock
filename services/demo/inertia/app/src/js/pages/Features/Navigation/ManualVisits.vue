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

const breadcrumbs = [{ title: "Features" }, { title: "Navigation" }, { title: "Manual Visits" }];

function visitGet() {
  router.get(page.url);
}

function visitPost() {
  router.post(page.url);
}

function visitPut() {
  router.put(page.url);
}

function visitPatch() {
  router.patch(page.url);
}

function visitDelete() {
  router.delete(page.url);
}

function visitWithOptions() {
  router.get(
    page.url,
    {},
    {
      preserveState: true,
      preserveScroll: true,
      replace: true,
    },
  );
}
</script>

<template>
  <AppLayout title="Manual Visits" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Manual Visits"
        description="In addition to the Link component, you can make programmatic Inertia visits using the router object with methods like get, post, put, patch, and delete."
      />

      <p class="text-muted-foreground text-sm">
        Server timestamp:
        <code class="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">{{ timestamp }}</code>
      </p>

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard
          title="HTTP Methods"
          description="Each button triggers a different HTTP method via router."
        >
          <div class="flex flex-wrap gap-3">
            <Button @click="visitGet">router.get</Button>
            <Button variant="outline" @click="visitPost">router.post</Button>
            <Button variant="outline" @click="visitPut">router.put</Button>
            <Button variant="outline" @click="visitPatch">router.patch</Button>
            <Button variant="destructive" @click="visitDelete">router.delete</Button>
          </div>
          <p class="text-muted-foreground mt-2 text-xs">
            Open DevTools Network tab to see the HTTP method used for each request.
          </p>
        </FeatureCard>

        <FeatureCard
          title="With Options"
          description="Manual visits support the same options as Link: preserveState, preserveScroll, replace, and more."
        >
          <div class="flex flex-wrap gap-3">
            <Button @click="visitWithOptions">GET (preserve state + scroll, replace)</Button>
          </div>
          <pre
            class="mt-3 overflow-auto rounded-md bg-muted p-3 text-xs leading-relaxed"
          ><code>router.get(url, {}, {
  preserveState: true,
  preserveScroll: true,
  replace: true,
})</code></pre>
        </FeatureCard>

        <FeatureCard
          title="Usage"
          description="Import router from @inertiajs/vue3 and call visit methods."
        >
          <pre
            class="overflow-auto rounded-md bg-muted p-3 text-xs leading-relaxed"
          ><code>import { router } from '@inertiajs/vue3'

router.get('/url')
router.post('/url', data)
router.put('/url', data)
router.patch('/url', data)
router.delete('/url')</code></pre>
        </FeatureCard>

        <FeatureCard
          title="Generic Visit"
          description="The router.visit() method accepts a method option for full control."
        >
          <pre
            class="overflow-auto rounded-md bg-muted p-3 text-xs leading-relaxed"
          ><code>router.visit('/url', {
  method: 'post',
  data: { name: 'John' },
  preserveState: true,
  onSuccess: (page) => { ... },
})</code></pre>
        </FeatureCard>
      </div>
    </div>
  </AppLayout>
</template>
