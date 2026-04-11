<script setup lang="ts">
import { useForm } from "@inertiajs/vue3";
import GuestLayout from "@/Layouts/GuestLayout.vue";

defineProps<{ status?: string }>();

const form = useForm({ email: "" });

function submit() {
  form.post("/forgot-password");
}
</script>

<template>
  <GuestLayout title="Forgot Password">
    <div class="mb-4 text-sm text-gray-600">
      Forgot your password? Enter your email and we'll send you a password reset link.
    </div>

    <div v-if="status" class="mb-4 text-sm font-medium text-green-600">{{ status }}</div>

    <form @submit.prevent="submit">
      <div>
        <label for="email" class="block text-sm font-medium text-gray-700">Email</label>
        <input id="email" v-model="form.email" type="email" required autofocus
          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" />
        <p v-if="form.errors.email" class="mt-1 text-sm text-red-600">{{ form.errors.email }}</p>
      </div>

      <div class="mt-4 flex items-center justify-end">
        <button type="submit" :disabled="form.processing"
          class="rounded-md bg-gray-800 px-4 py-2 text-xs font-semibold uppercase text-white hover:bg-gray-700 disabled:opacity-50">
          Email Password Reset Link
        </button>
      </div>
    </form>
  </GuestLayout>
</template>
