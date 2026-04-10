<script setup lang="ts">
import { useForm } from "@inertiajs/vue3";
import GuestLayout from "@/Layouts/GuestLayout.vue";

const props = defineProps<{ email: string; token: string }>();

const form = useForm({
  token: props.token,
  email: props.email,
  password: "",
  password_confirmation: "",
});

function submit() {
  form.post("/reset-password", {
    onFinish: () => form.reset("password", "password_confirmation"),
  });
}
</script>

<template>
  <GuestLayout title="Reset Password">
    <form @submit.prevent="submit">
      <div>
        <label for="email" class="block text-sm font-medium text-gray-700">Email</label>
        <input id="email" v-model="form.email" type="email" required autofocus
          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" />
        <p v-if="form.errors.email" class="mt-1 text-sm text-red-600">{{ form.errors.email }}</p>
      </div>

      <div class="mt-4">
        <label for="password" class="block text-sm font-medium text-gray-700">Password</label>
        <input id="password" v-model="form.password" type="password" required
          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" />
        <p v-if="form.errors.password" class="mt-1 text-sm text-red-600">{{ form.errors.password }}</p>
      </div>

      <div class="mt-4">
        <label for="password_confirmation" class="block text-sm font-medium text-gray-700">Confirm Password</label>
        <input id="password_confirmation" v-model="form.password_confirmation" type="password" required
          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" />
      </div>

      <div class="mt-4 flex items-center justify-end">
        <button type="submit" :disabled="form.processing"
          class="rounded-md bg-gray-800 px-4 py-2 text-xs font-semibold uppercase text-white hover:bg-gray-700 disabled:opacity-50">
          Reset Password
        </button>
      </div>
    </form>
  </GuestLayout>
</template>
