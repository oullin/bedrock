<script setup lang="ts">
import { ref } from "vue";
import { useRemember, router, Link } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Button } from "@/js/components/ui/button";
import { Input } from "@/js/components/ui/input";
import { Badge } from "@/js/components/ui/badge";
import AppLayout from "@/js/layouts/AppLayout.vue";
import { featureRoute } from "@/js/lib/routes";

const breadcrumbs = [{ title: "Features" }, { title: "State" }, { title: "Remember" }];

const flashDataRoute = featureRoute("features.state.flash-data");

const formState = useRemember({
  name: "",
  email: "",
  message: "",
});

const counter = ref<number>((router.restore("counter") as number) ?? 0);

function incrementCounter() {
  counter.value++;
  router.remember(counter.value, "counter");
}

function resetCounter() {
  counter.value = 0;
  router.remember(0, "counter");
}
</script>

<template>
  <AppLayout title="Remember" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Remember"
        description="Persist component state across navigations using useRemember and router.remember/restore. Navigate away and come back to see your data preserved."
      />

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard
          title="useRemember Form"
          description="This form state persists across navigation."
        >
          <div class="space-y-4">
            <div class="grid gap-3">
              <div class="grid gap-1.5">
                <label class="text-sm font-medium" for="name">Name</label>
                <Input id="name" v-model="formState.name" placeholder="Enter your name" />
              </div>
              <div class="grid gap-1.5">
                <label class="text-sm font-medium" for="email">Email</label>
                <Input
                  id="email"
                  v-model="formState.email"
                  type="email"
                  placeholder="you@example.com"
                />
              </div>
              <div class="grid gap-1.5">
                <label class="text-sm font-medium" for="message">Message</label>
                <Input id="message" v-model="formState.message" placeholder="Your message..." />
              </div>
            </div>

            <div class="flex items-center gap-2">
              <Link
                v-if="flashDataRoute"
                :href="flashDataRoute"
                class="inline-flex items-center rounded-md border px-3 py-2 text-sm font-medium hover:bg-accent"
              >
                Navigate Away
              </Link>
              <span class="text-muted-foreground text-xs"
                >then come back to see data preserved</span
              >
            </div>
          </div>
        </FeatureCard>

        <FeatureCard
          title="Manual Remember/Restore"
          description="Using router.remember() and router.restore() directly."
        >
          <div class="space-y-4">
            <div class="flex items-center gap-4">
              <span class="text-3xl font-semibold">{{ counter }}</span>
              <Badge variant="secondary">Persisted</Badge>
            </div>

            <div class="flex gap-2">
              <Button @click="incrementCounter">Increment</Button>
              <Button variant="outline" @click="resetCounter">Reset</Button>
            </div>

            <p class="text-muted-foreground text-sm">
              This counter is saved with
              <code class="bg-muted rounded px-1.5 py-0.5 text-xs">router.remember()</code> and
              restored with
              <code class="bg-muted rounded px-1.5 py-0.5 text-xs">router.restore()</code>. Navigate
              away and back to verify it persists.
            </p>
          </div>
        </FeatureCard>
      </div>
    </div>
  </AppLayout>
</template>
