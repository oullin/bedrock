<script setup lang="ts">
import { Head, useForm } from "@inertiajs/vue3";
import InputError from "@/js/components/app/InputError.vue";
import { Button } from "@/js/components/ui/button";
import { Checkbox } from "@/js/components/ui/checkbox";
import { Input } from "@/js/components/ui/input";
import { Label } from "@/js/components/ui/label";
import AuthLayout from "@/js/layouts/AuthLayout.vue";
import { appRoutes } from "@/js/lib/routes";

withDefaults(defineProps<{ status?: string }>(), {
  status: "",
});

const form = useForm({
  email: "test@example.com",
  password: "password",
  remember: true,
});

function submit() {
  form.post(appRoutes.login().url);
}
</script>

<template>
  <Head title="Log in" />

  <AuthLayout
    title="Inertia.js Kitchen Sink"
    description="Explore the Go port of the official demo with the seeded account below."
  >
    <div v-if="status" class="mb-4 text-center text-sm font-medium text-green-600">
      {{ status }}
    </div>

    <form class="flex flex-col gap-6" @submit.prevent="submit">
      <div class="grid gap-2">
        <Label for="email">Email address</Label>
        <Input id="email" v-model="form.email" type="email" autocomplete="email" />
        <InputError :message="form.errors.email" />
      </div>

      <div class="grid gap-2">
        <Label for="password">Password</Label>
        <Input
          id="password"
          v-model="form.password"
          type="password"
          autocomplete="current-password"
        />
        <InputError :message="form.errors.password" />
      </div>

      <Label for="remember" class="flex items-center gap-3">
        <Checkbox
          id="remember"
          :checked="form.remember"
          @update:checked="form.remember = Boolean($event)"
        />
        <span>Remember me</span>
      </Label>

      <Button type="submit" class="w-full" :disabled="form.processing">
        {{ form.processing ? "Signing in..." : "Log in" }}
      </Button>
    </form>
  </AuthLayout>
</template>
