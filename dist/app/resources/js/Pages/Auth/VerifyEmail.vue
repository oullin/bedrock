<script setup lang="ts">
import { useForm } from "@inertiajs/vue3";
import AppLayout from "@/Layouts/AppLayout.vue";

defineProps<{ status?: string }>();

const form = useForm({});

function submit() {
  form.post("/email/verification-notification");
}
</script>

<template>
  <AppLayout title="Email Verification">
    <template #header>
      <h2 class="text-xl font-semibold leading-tight text-gray-800">Email Verification</h2>
    </template>

    <div class="rounded-md bg-white p-6 shadow">
      <div class="mb-4 text-sm text-gray-600">
        Before continuing, please verify your email address by clicking the link we sent you.
        If you didn't receive the email, we can send another.
      </div>

      <div v-if="status" class="mb-4 text-sm font-medium text-green-600">
        A new verification link has been sent to your email address.
      </div>

      <form @submit.prevent="submit">
        <button type="submit" :disabled="form.processing"
          class="rounded-md bg-gray-800 px-4 py-2 text-xs font-semibold uppercase text-white hover:bg-gray-700 disabled:opacity-50">
          Resend Verification Email
        </button>
      </form>
    </div>
  </AppLayout>
</template>
