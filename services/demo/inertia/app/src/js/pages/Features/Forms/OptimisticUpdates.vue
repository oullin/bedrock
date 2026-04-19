<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { router } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Button } from "@/js/components/ui/button";
import AppLayout from "@/js/layouts/AppLayout.vue";
import type { Contact } from "@/js/types";

const props = withDefaults(defineProps<{ contacts?: Contact[] }>(), {
  contacts: () => [],
});

function cloneRaw(data: Contact[]): Contact[] {
  return JSON.parse(JSON.stringify(data));
}

const localContacts = ref<Contact[]>(cloneRaw(props.contacts));

watch(
  () => props.contacts,
  (value) => {
    localContacts.value = cloneRaw(value);
  },
);

const breadcrumbs = [{ title: "Features" }, { title: "Forms" }, { title: "Optimistic Updates" }];

const favorites = computed(() => localContacts.value.filter((c) => c.is_favorite));

function toggleFavorite(contact: Contact) {
  contact.is_favorite = !contact.is_favorite;

  router.post(
    `/features/forms/optimistic-toggle/${contact.id}`,
    {},
    {
      preserveScroll: true,
      only: ["contacts"],
      onError() {
        contact.is_favorite = !contact.is_favorite;
      },
    },
  );
}
</script>

<template>
  <AppLayout title="Optimistic Updates" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Optimistic Updates"
        description="Optimistic UI updates provide instant visual feedback by assuming the server request will succeed, then reconciling when the response arrives."
      />

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard
          title="Contacts"
          description="Click the star button to toggle favorites. The UI updates optimistically."
        >
          <div v-if="localContacts.length" class="grid gap-2">
            <div
              v-for="contact in localContacts"
              :key="contact.id"
              class="flex items-center justify-between rounded-md border p-3 transition-colors"
              :class="{ 'bg-yellow-50 border-yellow-200': contact.is_favorite }"
            >
              <div class="min-w-0">
                <p class="font-medium">{{ contact.first_name }} {{ contact.last_name }}</p>
                <p class="text-muted-foreground truncate text-sm">{{ contact.email }}</p>
              </div>
              <Button variant="ghost" size="sm" class="shrink-0" @click="toggleFavorite(contact)">
                <span class="text-lg">{{ contact.is_favorite ? "★" : "☆" }}</span>
              </Button>
            </div>
          </div>
          <div v-else class="text-muted-foreground py-8 text-center text-sm">
            <p>No contacts available.</p>
            <p class="mt-1">Contacts will be passed as props from the server.</p>
          </div>
        </FeatureCard>

        <div class="grid gap-6">
          <FeatureCard
            title="How It Works"
            description="Understanding optimistic updates in Inertia."
          >
            <div class="grid gap-2 text-sm">
              <div class="rounded-md border p-3">
                <p class="mb-1 font-medium">1. User Action</p>
                <p class="text-muted-foreground">
                  User clicks the favorite toggle. The UI updates immediately using
                  <code class="bg-muted rounded px-1">router.post()</code>.
                </p>
              </div>
              <div class="rounded-md border p-3">
                <p class="mb-1 font-medium">2. Server Request</p>
                <p class="text-muted-foreground">
                  A POST request is sent to update the server state. We use
                  <code class="bg-muted rounded px-1">preserveScroll</code> and
                  <code class="bg-muted rounded px-1">only: ['contacts']</code>
                  for a smooth experience.
                </p>
              </div>
              <div class="rounded-md border p-3">
                <p class="mb-1 font-medium">3. Reconciliation</p>
                <p class="text-muted-foreground">
                  When the response arrives, Inertia reconciles the server state with the UI. If the
                  request failed, the UI reverts automatically.
                </p>
              </div>
            </div>
          </FeatureCard>

          <FeatureCard title="Favorites" description="Current favorite contacts.">
            <div class="text-sm">
              <div v-if="favorites.length" class="grid gap-1">
                <div
                  v-for="fav in favorites"
                  :key="fav.id"
                  class="flex items-center gap-2 rounded bg-yellow-50 px-3 py-2"
                >
                  <span class="text-yellow-600">★</span>
                  <span>{{ fav.first_name }} {{ fav.last_name }}</span>
                </div>
              </div>
              <p v-else class="text-muted-foreground italic">No favorites yet</p>
            </div>
          </FeatureCard>
        </div>
      </div>
    </div>
  </AppLayout>
</template>
