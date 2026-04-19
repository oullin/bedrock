<script setup lang="ts">
import { useForm, usePage } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import InputError from "@/js/components/app/InputError.vue";
import { Button } from "@/js/components/ui/button";
import { Input } from "@/js/components/ui/input";
import { Label } from "@/js/components/ui/label";
import AppLayout from "@/js/layouts/AppLayout.vue";
import type { SharedPageProps } from "@/js/types";

type FieldName = "username" | "email" | "password" | "password_confirmation";

const fields: FieldName[] = ["username", "email", "password", "password_confirmation"];

const page = usePage<SharedPageProps>();

const breadcrumbs = [{ title: "Features" }, { title: "Forms" }, { title: "Precognition" }];

const form = useForm("post", page.url, {
  username: "",
  email: "",
  password: "",
  password_confirmation: "",
});

function submit() {
  form.submit();
}

function fieldClass(field: FieldName): string {
  if (!form.touched(field)) return "";
  if (form.invalid(field)) return "border-red-500";
  if (form.valid(field)) return "border-green-500";
  return "";
}
</script>

<template>
  <AppLayout title="Precognition" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Precognition"
        description="Precognitive validation allows you to validate form fields in real-time by sending partial requests to the server before the form is submitted."
      />

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard
          title="Create Account"
          description="Fields are validated as you type, providing instant feedback."
        >
          <form class="grid gap-4" @submit.prevent="submit">
            <div class="grid gap-2">
              <Label for="username">Username</Label>
              <Input
                id="username"
                v-model="form.username"
                placeholder="Choose a username"
                :class="fieldClass('username')"
                @blur="form.validate('username')"
              />
              <InputError :message="form.errors.username" />
              <p
                v-if="form.touched('username') && form.valid('username') && form.username"
                class="text-xs text-green-600"
              >
                Username looks good!
              </p>
            </div>

            <div class="grid gap-2">
              <Label for="email">Email</Label>
              <Input
                id="email"
                v-model="form.email"
                type="email"
                placeholder="you@example.com"
                :class="fieldClass('email')"
                @blur="form.validate('email')"
              />
              <InputError :message="form.errors.email" />
            </div>

            <div class="grid gap-2">
              <Label for="password">Password</Label>
              <Input
                id="password"
                v-model="form.password"
                type="password"
                placeholder="At least 8 characters"
                :class="fieldClass('password')"
                @blur="form.validate('password')"
              />
              <InputError :message="form.errors.password" />
            </div>

            <div class="grid gap-2">
              <Label for="password_confirmation">Confirm Password</Label>
              <Input
                id="password_confirmation"
                v-model="form.password_confirmation"
                type="password"
                placeholder="Repeat your password"
                :class="fieldClass('password_confirmation')"
                @blur="form.validate('password_confirmation')"
              />
              <InputError :message="form.errors.password_confirmation" />
            </div>

            <div class="flex items-center gap-3">
              <Button type="submit" :disabled="form.processing">
                {{ form.processing ? "Creating..." : "Create Account" }}
              </Button>
              <span v-if="form.recentlySuccessful" class="text-sm text-green-600">
                Account created!
              </span>
            </div>
          </form>
        </FeatureCard>

        <div class="grid gap-6">
          <FeatureCard title="Field States" description="Real-time status of each form field.">
            <div class="grid gap-2 text-sm">
              <div
                v-for="field in fields"
                :key="field"
                class="flex items-center justify-between rounded-md border p-3"
              >
                <span class="font-medium font-mono">{{ field }}</span>
                <div class="flex items-center gap-2">
                  <span
                    v-if="form.touched(field)"
                    class="rounded-full px-2 py-0.5 text-xs font-medium bg-blue-100 text-blue-800"
                  >
                    touched
                  </span>
                  <span
                    v-if="form.touched(field) && form.invalid(field)"
                    class="rounded-full px-2 py-0.5 text-xs font-medium bg-red-100 text-red-800"
                  >
                    invalid
                  </span>
                  <span
                    v-else-if="form.touched(field) && form.valid(field)"
                    class="rounded-full px-2 py-0.5 text-xs font-medium bg-green-100 text-green-800"
                  >
                    valid
                  </span>
                  <span
                    v-else
                    class="rounded-full px-2 py-0.5 text-xs font-medium bg-gray-100 text-gray-600"
                  >
                    pristine
                  </span>
                </div>
              </div>
            </div>
          </FeatureCard>

          <FeatureCard
            title="How Precognition Works"
            description="Understanding the validation flow."
          >
            <div class="grid gap-2 text-sm">
              <div class="rounded-md border p-3">
                <p class="mb-1 font-medium">1. Field Blur</p>
                <p class="text-muted-foreground">
                  When a user leaves a field, a precognitive request is sent to validate just that
                  field.
                </p>
              </div>
              <div class="rounded-md border p-3">
                <p class="mb-1 font-medium">2. Server Validation</p>
                <p class="text-muted-foreground">
                  The server runs the same validation rules but returns only errors, not a full
                  response.
                </p>
              </div>
              <div class="rounded-md border p-3">
                <p class="mb-1 font-medium">3. Instant Feedback</p>
                <p class="text-muted-foreground">
                  Errors are displayed immediately without submitting the entire form.
                </p>
              </div>
            </div>
          </FeatureCard>
        </div>
      </div>
    </div>
  </AppLayout>
</template>
