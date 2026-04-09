<script setup lang="ts">
import { useForm } from "@inertiajs/vue3";
import AppLayout from "@/Layouts/AppLayout.vue";

const form = useForm({ name: "" });

function submit() {
  form.post("/teams");
}
</script>

<template>
  <AppLayout title="Create Team">
    <template #header>
      <h2 class="text-xl font-semibold leading-tight text-gray-800">Create Team</h2>
    </template>

    <div class="rounded-lg bg-white p-6 shadow">
      <h3 class="text-lg font-medium text-gray-900">Team Details</h3>
      <p class="mt-1 text-sm text-gray-600">Create a new team to collaborate with others.</p>

      <form @submit.prevent="submit" class="mt-6 space-y-4">
        <div>
          <label for="name" class="block text-sm font-medium text-gray-700">Team Name</label>
          <input id="name" v-model="form.name" type="text" required autofocus
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" />
          <p v-if="form.errors.name" class="mt-1 text-sm text-red-600">{{ form.errors.name }}</p>
        </div>

        <button type="submit" :disabled="form.processing"
          class="rounded-md bg-gray-800 px-4 py-2 text-xs font-semibold uppercase text-white hover:bg-gray-700 disabled:opacity-50">
          Create
        </button>
      </form>
    </div>
  </AppLayout>
</template>
