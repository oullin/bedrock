<script setup lang="ts">
import { useForm, usePage } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import InputError from "@/js/components/app/InputError.vue";
import { Button } from "@/js/components/ui/button";
import { Input } from "@/js/components/ui/input";
import { Textarea } from "@/js/components/ui/textarea";
import { Label } from "@/js/components/ui/label";
import AppLayout from "@/js/layouts/AppLayout.vue";
import type { SharedPageProps } from "@/js/types";

const page = usePage<SharedPageProps>();

const breadcrumbs = [{ title: "Features" }, { title: "Forms" }, { title: "Validation" }];

const primaryForm = useForm({
  name: "",
  email: "",
  age: "",
});

const secondaryForm = useForm({
  title: "",
  body: "",
});

function submitPrimary() {
  primaryForm.post(page.url);
}

function submitSecondary() {
  secondaryForm.post(page.url, {
    errorBag: "secondary",
  });
}
</script>

<template>
  <AppLayout title="Validation" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Validation"
        description="Inertia.js handles server-side validation errors automatically. Errors are returned as props and mapped to form fields."
      />

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard
          title="Primary Form"
          description="Standard validation with inline error messages."
        >
          <form class="grid gap-4" @submit.prevent="submitPrimary">
            <div class="grid gap-2">
              <Label for="name">Name <span class="text-red-500">*</span></Label>
              <Input
                id="name"
                v-model="primaryForm.name"
                placeholder="Required field"
                :class="{ 'border-red-500': primaryForm.errors.name }"
              />
              <InputError :message="primaryForm.errors.name" />
            </div>

            <div class="grid gap-2">
              <Label for="email">Email <span class="text-red-500">*</span></Label>
              <Input
                id="email"
                v-model="primaryForm.email"
                type="email"
                placeholder="Valid email required"
                :class="{ 'border-red-500': primaryForm.errors.email }"
              />
              <InputError :message="primaryForm.errors.email" />
            </div>

            <div class="grid gap-2">
              <Label for="age"
                >Age <span class="text-muted-foreground text-xs">(18-120)</span></Label
              >
              <Input
                id="age"
                v-model="primaryForm.age"
                type="number"
                min="18"
                max="120"
                placeholder="Must be between 18 and 120"
                :class="{ 'border-red-500': primaryForm.errors.age }"
              />
              <InputError :message="primaryForm.errors.age" />
            </div>

            <div class="flex items-center gap-3">
              <Button type="submit" :disabled="primaryForm.processing">
                {{ primaryForm.processing ? "Validating..." : "Submit" }}
              </Button>
              <Button type="button" variant="outline" @click="primaryForm.clearErrors()">
                Clear Errors
              </Button>
              <span v-if="primaryForm.recentlySuccessful" class="text-sm text-green-600">
                Valid!
              </span>
            </div>

            <div
              v-if="primaryForm.hasErrors"
              class="rounded-md border border-red-200 bg-red-50 p-3"
            >
              <p class="mb-1 text-sm font-medium text-red-800">
                There were {{ Object.keys(primaryForm.errors).length }} error(s) with your
                submission
              </p>
              <ul class="list-inside list-disc text-sm text-red-700">
                <li v-for="(error, field) in primaryForm.errors" :key="field">{{ error }}</li>
              </ul>
            </div>
          </form>
        </FeatureCard>

        <div class="grid gap-6">
          <FeatureCard
            title="Secondary Form (Error Bag)"
            description="Uses errorBag to scope errors separately from the primary form."
          >
            <form class="grid gap-4" @submit.prevent="submitSecondary">
              <div class="grid gap-2">
                <Label for="title">Title</Label>
                <Input
                  id="title"
                  v-model="secondaryForm.title"
                  placeholder="Post title"
                  :class="{ 'border-red-500': secondaryForm.errors.title }"
                />
                <InputError :message="secondaryForm.errors.title" />
              </div>

              <div class="grid gap-2">
                <Label for="body">Body</Label>
                <Textarea
                  id="body"
                  v-model="secondaryForm.body"
                  rows="3"
                  :class="{ 'border-red-500': secondaryForm.errors.body }"
                  placeholder="Post body content"
                />
                <InputError :message="secondaryForm.errors.body" />
              </div>

              <div class="flex items-center gap-3">
                <Button type="submit" :disabled="secondaryForm.processing">
                  {{ secondaryForm.processing ? "Saving..." : "Save Post" }}
                </Button>
                <Button type="button" variant="outline" @click="secondaryForm.clearErrors()">
                  Clear Errors
                </Button>
              </div>
            </form>
          </FeatureCard>

          <FeatureCard
            title="How Error Bags Work"
            description="Error bags isolate validation errors between multiple forms on the same page."
          >
            <div class="grid gap-2 text-sm">
              <div class="rounded-md border p-3">
                <p class="text-muted-foreground">
                  When you have multiple forms on a single page, each form can use an
                  <code class="bg-muted rounded px-1">errorBag</code> to ensure that validation
                  errors from one form do not leak into another.
                </p>
              </div>
              <div class="rounded-md border p-3">
                <p class="mb-1 font-medium">Primary Errors</p>
                <pre class="bg-muted overflow-auto rounded p-2 text-xs">{{
                  JSON.stringify(primaryForm.errors, null, 2)
                }}</pre>
              </div>
              <div class="rounded-md border p-3">
                <p class="mb-1 font-medium">Secondary Errors (error bag)</p>
                <pre class="bg-muted overflow-auto rounded p-2 text-xs">{{
                  JSON.stringify(secondaryForm.errors, null, 2)
                }}</pre>
              </div>
            </div>
          </FeatureCard>
        </div>
      </div>
    </div>
  </AppLayout>
</template>
