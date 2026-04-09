<script setup lang="ts">
import { ref } from "vue";
import { useForm } from "@inertiajs/vue3";
import GuestLayout from "@/Layouts/GuestLayout.vue";

const useRecovery = ref(false);

const form = useForm({
  code: "",
  recovery_code: "",
});

function submit() {
  form.post("/two-factor-challenge", {
    onFinish: () => form.reset("code", "recovery_code"),
  });
}

function toggleRecovery() {
  useRecovery.value = !useRecovery.value;
  form.reset();
}
</script>

<template>
  <GuestLayout title="Two Factor Challenge">
    <div class="mb-4 text-sm text-gray-600">
      <template v-if="!useRecovery">
        Please confirm access to your account by entering the authentication code from your authenticator app.
      </template>
      <template v-else>
        Please confirm access to your account by entering one of your emergency recovery codes.
      </template>
    </div>

    <form @submit.prevent="submit">
      <div v-if="!useRecovery">
        <label for="code" class="block text-sm font-medium text-gray-700">Code</label>
        <input id="code" v-model="form.code" type="text" inputmode="numeric" autofocus autocomplete="one-time-code"
          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" />
        <p v-if="form.errors.code" class="mt-1 text-sm text-red-600">{{ form.errors.code }}</p>
      </div>

      <div v-else>
        <label for="recovery_code" class="block text-sm font-medium text-gray-700">Recovery Code</label>
        <input id="recovery_code" v-model="form.recovery_code" type="text" autofocus autocomplete="one-time-code"
          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" />
        <p v-if="form.errors.recovery_code" class="mt-1 text-sm text-red-600">{{ form.errors.recovery_code }}</p>
      </div>

      <div class="mt-4 flex items-center justify-between">
        <button type="button" @click="toggleRecovery" class="text-sm text-gray-600 underline hover:text-gray-900">
          {{ useRecovery ? "Use an authentication code" : "Use a recovery code" }}
        </button>

        <button type="submit" :disabled="form.processing"
          class="rounded-md bg-gray-800 px-4 py-2 text-xs font-semibold uppercase text-white hover:bg-gray-700 disabled:opacity-50">
          Log in
        </button>
      </div>
    </form>
  </GuestLayout>
</template>
