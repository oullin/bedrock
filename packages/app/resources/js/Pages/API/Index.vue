<script setup lang="ts">
import { useForm, router } from "@inertiajs/vue3";
import { ref } from "vue";
import AppLayout from "@/Layouts/AppLayout.vue";

defineProps<{
  tokens: Array<{ id: string; name: string; permissions: string[]; last_used_at?: string }>;
  availablePermissions: string[];
}>();

const createForm = useForm({ name: "", permissions: [] as string[] });
const newToken = ref<string | null>(null);

function createToken() {
  createForm.post("/user/api-tokens", {
    onSuccess: (page: any) => {
      if (page.props?.flash?.token) {
        newToken.value = page.props.flash.token;
      }
      createForm.reset();
    },
  });
}

function deleteToken(tokenId: string) {
  if (confirm("Are you sure?")) {
    router.delete(`/user/api-tokens/${tokenId}`);
  }
}
</script>

<template>
  <AppLayout title="API Tokens">
    <template #header>
      <h2 class="text-xl font-semibold leading-tight text-gray-800">API Tokens</h2>
    </template>

    <div class="space-y-6">
      <!-- Create Token -->
      <div class="rounded-lg bg-white p-6 shadow">
        <h3 class="text-lg font-medium text-gray-900">Create API Token</h3>
        <p class="mt-1 text-sm text-gray-600">API tokens allow third-party services to authenticate with our application on your behalf.</p>

        <form @submit.prevent="createToken" class="mt-6 space-y-4">
          <div>
            <label for="token-name" class="block text-sm font-medium text-gray-700">Name</label>
            <input id="token-name" v-model="createForm.name" type="text" required
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" />
            <p v-if="createForm.errors.name" class="mt-1 text-sm text-red-600">{{ createForm.errors.name }}</p>
          </div>

          <div v-if="availablePermissions.length > 0">
            <label class="block text-sm font-medium text-gray-700">Permissions</label>
            <div class="mt-2 grid grid-cols-2 gap-2">
              <label v-for="perm in availablePermissions" :key="perm" class="flex items-center gap-2">
                <input v-model="createForm.permissions" type="checkbox" :value="perm" class="rounded border-gray-300 text-indigo-600" />
                <span class="text-sm text-gray-700">{{ perm }}</span>
              </label>
            </div>
          </div>

          <button type="submit" :disabled="createForm.processing"
            class="rounded-md bg-gray-800 px-4 py-2 text-xs font-semibold uppercase text-white hover:bg-gray-700 disabled:opacity-50">
            Create
          </button>
        </form>

        <!-- New Token Display -->
        <div v-if="newToken" class="mt-4 rounded-md bg-gray-100 p-4">
          <p class="text-sm font-medium text-gray-700">Please copy your new API token. For your security, it won't be shown again.</p>
          <code class="mt-2 block break-all text-sm text-gray-900">{{ newToken }}</code>
        </div>
      </div>

      <!-- Existing Tokens -->
      <div v-if="tokens.length > 0" class="rounded-lg bg-white p-6 shadow">
        <h3 class="text-lg font-medium text-gray-900">Manage API Tokens</h3>

        <div class="mt-4 space-y-3">
          <div v-for="token in tokens" :key="token.id" class="flex items-center justify-between">
            <div>
              <span class="text-sm font-medium text-gray-700">{{ token.name }}</span>
              <span v-if="token.last_used_at" class="ml-2 text-xs text-gray-400">Last used {{ token.last_used_at }}</span>
              <div class="text-xs text-gray-500">{{ token.permissions.join(', ') }}</div>
            </div>
            <button @click="deleteToken(token.id)" class="text-sm text-red-600 hover:text-red-800">Delete</button>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>
