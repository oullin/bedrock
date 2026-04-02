<script setup lang="ts">
import { summarizeBillingWorkspace } from "@gollin/billing";

import { Button } from "@/components/ui/button";

const billingSummary = summarizeBillingWorkspace();
</script>

<template>
  <main class="min-h-screen bg-stone-100 px-6 py-12 text-stone-950">
    <section class="mx-auto max-w-6xl space-y-8">
      <div class="space-y-4 text-center">
        <p class="text-sm uppercase tracking-[0.35em] text-stone-500">Pricing</p>
        <h1 class="text-4xl font-semibold tracking-tight sm:text-5xl">
          Simple packaging for the first monorepo release.
        </h1>
        <p class="mx-auto max-w-2xl text-lg leading-8 text-stone-600">
          The pricing app shares the same Vite and Tailwind baseline, with `shadcn-vue` components
          owned directly in the repository and an internal package wired through `workspace:*`.
        </p>
      </div>

      <div class="grid gap-6 lg:grid-cols-3">
        <article
          v-for="tier in billingSummary.tiers"
          :key="tier.name"
          class="rounded-3xl border border-stone-300 bg-white p-6 shadow-sm"
        >
          <p class="text-sm uppercase tracking-[0.3em] text-stone-500">{{ tier.name }}</p>
          <p class="mt-4 text-4xl font-semibold text-stone-950">{{ tier.price }}</p>
          <p class="mt-4 min-h-16 text-sm leading-7 text-stone-600">{{ tier.description }}</p>
          <Button class="mt-6 w-full" :variant="tier.featured ? 'default' : 'outline'">
            {{ tier.featured ? "Start featured plan" : "View plan" }}
          </Button>
        </article>
      </div>
    </section>
  </main>
</template>
