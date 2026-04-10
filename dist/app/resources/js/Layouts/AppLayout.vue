<script setup lang="ts">
import { Head, Link, router } from "@inertiajs/vue3";
import { ref, computed } from "vue";

const props = defineProps<{
  title?: string;
}>();

const page = computed(() => router.page);
const user = computed(() => page.value?.props?.auth?.user);
const currentTeam = computed(() => page.value?.props?.auth?.currentTeam);

const showNav = ref(false);

function logout() {
  router.post("/logout");
}

function switchTeam(teamId: string) {
  router.put("/current-team", { team_id: teamId });
}
</script>

<template>
  <div class="min-h-screen bg-gray-100">
    <Head :title="title" />

    <!-- Navigation -->
    <nav class="border-b border-gray-200 bg-white">
      <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div class="flex h-16 justify-between">
          <div class="flex">
            <div class="flex shrink-0 items-center">
              <Link href="/dashboard" class="text-xl font-bold text-gray-900">Bedrock</Link>
            </div>

            <div class="hidden space-x-8 sm:-my-px sm:ms-10 sm:flex">
              <Link href="/dashboard" class="inline-flex items-center border-b-2 border-transparent px-1 pt-1 text-sm font-medium text-gray-500 hover:border-gray-300 hover:text-gray-700">
                Dashboard
              </Link>
            </div>
          </div>

          <div class="hidden sm:ms-6 sm:flex sm:items-center">
            <!-- Team Switcher -->
            <div v-if="currentTeam" class="relative me-3">
              <span class="text-sm text-gray-500">{{ currentTeam.name }}</span>
            </div>

            <!-- User Menu -->
            <div class="relative">
              <button @click="showNav = !showNav" class="flex items-center text-sm font-medium text-gray-500 hover:text-gray-700">
                {{ user?.name || user?.email || 'Account' }}
                <svg class="ms-2 h-4 w-4 fill-current" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" /></svg>
              </button>

              <div v-show="showNav" class="absolute right-0 z-50 mt-2 w-48 rounded-md bg-white py-1 shadow-lg ring-1 ring-black/5">
                <Link href="/user/profile" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">Profile</Link>
                <Link v-if="page?.props?.jetstream?.hasApiTokens" href="/user/api-tokens" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">API Tokens</Link>
                <button @click="logout" class="block w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100">Log Out</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </nav>

    <!-- Page Header -->
    <header v-if="$slots.header" class="bg-white shadow">
      <div class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
        <slot name="header" />
      </div>
    </header>

    <!-- Page Content -->
    <main>
      <div class="mx-auto max-w-7xl py-6 sm:px-6 lg:px-8">
        <slot />
      </div>
    </main>
  </div>
</template>
