<script setup lang="ts">
import { ref } from "vue";
import { usePage } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Button } from "@/js/components/ui/button";
import { Input } from "@/js/components/ui/input";
import { Badge } from "@/js/components/ui/badge";
import AppLayout from "@/js/layouts/AppLayout.vue";
import type { SharedPageProps } from "@/js/types";

const page = usePage<SharedPageProps>();

const breadcrumbs = [{ title: "Features" }, { title: "HTTP" }, { title: "useHttp" }];

const name = ref("");
const email = ref("");
const response = ref<Record<string, unknown> | null>(null);
const loading = ref(false);
const error = ref<string | Record<string, string | string[]> | null>(null);

async function submitForm() {
  loading.value = true;
  error.value = null;
  response.value = null;

  try {
    const res = await fetch("/features/http/use-http/api", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
        "X-XSRF-TOKEN": getCookie("XSRF-TOKEN"),
      },
      body: JSON.stringify({
        name: name.value,
        email: email.value,
      }),
    });

    const contentType = res.headers.get("content-type") ?? "";

    if (!res.ok) {
      if (contentType.includes("application/json")) {
        const data = await res.json();
        error.value = data.errors ?? data.message ?? "Request failed";
      } else {
        error.value = (await res.text()) || "Request failed";
      }
    } else {
      const data = await res.json();
      response.value = data;
    }
  } catch (err) {
    error.value = (err as Error).message;
  } finally {
    loading.value = false;
  }
}

function getCookie(cookieName: string): string {
  const escaped = cookieName.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
  const match = document.cookie.match(new RegExp("(^| )" + escaped + "=([^;]+)"));
  return match ? decodeURIComponent(match[2]) : "";
}

function clearAll() {
  name.value = "";
  email.value = "";
  response.value = null;
  error.value = null;
}
</script>

<template>
  <AppLayout title="useHttp" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="useHttp"
        description="Make API requests that return JSON without triggering Inertia page navigations. Useful for inline data fetching, autocomplete, and AJAX-style interactions."
      />

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard
          title="API Form"
          description="Submit to a JSON endpoint without page navigation."
        >
          <form class="space-y-4" @submit.prevent="submitForm">
            <div class="grid gap-3">
              <div class="grid gap-1.5">
                <label class="text-sm font-medium" for="name">Name</label>
                <Input id="name" v-model="name" placeholder="John Doe" />
              </div>
              <div class="grid gap-1.5">
                <label class="text-sm font-medium" for="email">Email</label>
                <Input id="email" v-model="email" type="email" placeholder="john@example.com" />
              </div>
            </div>

            <div class="flex gap-2">
              <Button type="submit" :disabled="loading">
                {{ loading ? "Sending..." : "Send POST" }}
              </Button>
              <Button type="button" variant="outline" @click="clearAll">Clear</Button>
            </div>
          </form>
        </FeatureCard>

        <FeatureCard title="Response" description="JSON response from the API endpoint.">
          <div v-if="loading" class="text-muted-foreground text-sm">Loading...</div>

          <div v-else-if="error" class="space-y-2">
            <Badge variant="destructive">Error</Badge>
            <pre
              class="overflow-auto rounded border border-red-200 bg-red-50 p-3 text-xs text-red-700"
              >{{ typeof error === "string" ? error : JSON.stringify(error, null, 2) }}</pre
            >
          </div>

          <div v-else-if="response" class="space-y-2">
            <Badge>Success</Badge>
            <pre class="bg-muted overflow-auto rounded p-3 text-xs">{{
              JSON.stringify(response, null, 2)
            }}</pre>
          </div>

          <p v-else class="text-muted-foreground text-sm">
            Submit the form to see the JSON response here. The page will not navigate or reload.
          </p>
        </FeatureCard>
      </div>

      <FeatureCard title="When to Use" description="JSON APIs vs Inertia visits.">
        <div class="space-y-3 text-sm">
          <p>
            Most Inertia interactions use
            <code class="bg-muted rounded px-1.5 py-0.5 text-xs">router.visit()</code> or
            <code class="bg-muted rounded px-1.5 py-0.5 text-xs">useForm</code>, which trigger full
            page responses. However, some use cases benefit from plain JSON endpoints:
          </p>
          <ul class="space-y-2">
            <li class="flex items-start gap-2">
              <Badge variant="secondary" class="shrink-0 text-xs">Autocomplete</Badge>
              <span>Search suggestions that update inline without page navigation.</span>
            </li>
            <li class="flex items-start gap-2">
              <Badge variant="secondary" class="shrink-0 text-xs">Inline edits</Badge>
              <span>Update a single field without refreshing the whole page.</span>
            </li>
            <li class="flex items-start gap-2">
              <Badge variant="secondary" class="shrink-0 text-xs">Third-party APIs</Badge>
              <span>Proxy calls to external services that return JSON.</span>
            </li>
            <li class="flex items-start gap-2">
              <Badge variant="secondary" class="shrink-0 text-xs">Real-time data</Badge>
              <span>Fetch live data without replacing the current page state.</span>
            </li>
          </ul>
        </div>
      </FeatureCard>
    </div>
  </AppLayout>
</template>
