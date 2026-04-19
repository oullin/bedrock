<script setup lang="ts">
import { router, usePage } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Button } from "@/js/components/ui/button";
import { Badge } from "@/js/components/ui/badge";
import AppLayout from "@/js/layouts/AppLayout.vue";
import type { SharedPageProps } from "@/js/types";

const page = usePage<SharedPageProps>();

const breadcrumbs = [{ title: "Features" }, { title: "State" }, { title: "Flash Data" }];

function triggerFlash(kind: string) {
  router.post(page.url, { kind }, { preserveScroll: true });
}
</script>

<template>
  <AppLayout title="Flash Data" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Flash Data"
        description="Flash messages are one-time data passed from the server after a redirect. They are available in page.props.flash and cleared on the next request."
      />

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard
          title="Trigger Flash Messages"
          description="Click a button to send a POST and receive a flash message."
        >
          <div class="flex flex-wrap gap-2">
            <Button @click="triggerFlash('success')"> Success Flash </Button>
            <Button variant="destructive" @click="triggerFlash('error')"> Error Flash </Button>
            <Button variant="outline" @click="triggerFlash('warning')"> Warning Flash </Button>
          </div>
          <p class="text-muted-foreground mt-4 text-sm">
            After clicking, the server sets a flash message and redirects back. The flash message
            appears in the header bar above and in the panel to the right.
          </p>
        </FeatureCard>

        <FeatureCard title="Current Flash" description="What is currently in page.props.flash.">
          <div v-if="page.props.flash?.message" class="space-y-3">
            <div class="flex items-center gap-2">
              <Badge
                :variant="
                  { error: 'destructive', success: 'success', warning: 'warning' }[
                    page.props.flash.kind
                  ] ?? 'default'
                "
              >
                {{ page.props.flash.kind ?? "info" }}
              </Badge>
            </div>
            <div
              class="rounded-md border p-3"
              :class="{
                'border-green-200 bg-green-50': page.props.flash.kind === 'success',
                'border-red-200 bg-red-50': page.props.flash.kind === 'error',
                'border-yellow-200 bg-yellow-50': page.props.flash.kind === 'warning',
              }"
            >
              <p class="text-sm font-medium">{{ page.props.flash.title ?? "Notice" }}</p>
              <p class="text-sm">{{ page.props.flash.message }}</p>
            </div>
            <pre class="bg-muted overflow-auto rounded p-3 text-xs">{{
              JSON.stringify(page.props.flash, null, 2)
            }}</pre>
          </div>
          <p v-else class="text-muted-foreground text-sm">
            No flash data. Click one of the buttons to trigger a flash message.
          </p>
        </FeatureCard>
      </div>

      <FeatureCard title="How It Works">
        <div class="space-y-2 text-sm">
          <p>
            Flash data is set on the server during a request (typically before a redirect). Inertia
            includes it in
            <code class="bg-muted rounded px-1.5 py-0.5 text-xs">page.props.flash</code>
            as a shared prop.
          </p>
          <p>
            Flash data is automatically cleared after the next Inertia request, making it ideal for
            success/error notifications after form submissions.
          </p>
        </div>
      </FeatureCard>
    </div>
  </AppLayout>
</template>
