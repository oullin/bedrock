<script setup lang="ts">
import { useForm } from "@inertiajs/vue3";
import GuestLayout from "@/Layouts/GuestLayout.vue";

const form = useForm({
  email: "",
  password: "",
  remember: false,
});

function submit() {
  form.post("/login", {
    onFinish: () => form.reset("password"),
  });
}
</script>

<template>
  <GuestLayout title="Log in">
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

      <div class="mt-4 flex items-center">
        <input id="remember" v-model="form.remember" type="checkbox" class="rounded border-gray-300 text-indigo-600 shadow-sm focus:ring-indigo-500" />
        <label for="remember" class="ms-2 text-sm text-gray-600">Remember me</label>
      </div>

      <div class="mt-4 flex items-center justify-between">
        <a href="/forgot-password" class="text-sm text-gray-600 underline hover:text-gray-900">Forgot your password?</a>
        <button type="submit" :disabled="form.processing"
          class="rounded-md bg-gray-800 px-4 py-2 text-xs font-semibold uppercase text-white hover:bg-gray-700 disabled:opacity-50">
          Log in
        </button>
      </div>
    </form>
  </GuestLayout>
</template>
