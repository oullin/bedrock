<script setup lang="ts">
import { Link, Deferred, router, useForm } from "@inertiajs/vue3";
import InputError from "@/js/components/app/InputError.vue";
import { Badge } from "@/js/components/ui/badge";
import { Button } from "@/js/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/js/components/ui/card";
import { Skeleton } from "@/js/components/ui/skeleton";
import { Label } from "@/js/components/ui/label";
import { Textarea } from "@/js/components/ui/textarea";
import AppLayout from "@/js/layouts/AppLayout.vue";
import { contactRoutes } from "@/js/lib/routes";
import type { Contact, Note } from "@/js/types";

const props = withDefaults(defineProps<{ contact: Contact; notes?: Note[] }>(), {
  notes: undefined,
});

const breadcrumbs = [
  { title: "CRM" },
  { title: "Contacts", href: contactRoutes.index().url },
  { title: `${props.contact.first_name} ${props.contact.last_name}` },
];

const noteForm = useForm({
  body: "",
});

function submitNote() {
  noteForm.post(contactRoutes.storeNote(props.contact.id).url, {
    preserveScroll: true,
    onSuccess: () => noteForm.reset(),
  });
}

function toggleFavorite() {
  router.post(contactRoutes.favorite(props.contact.id).url, {}, { preserveScroll: true });
}

function deleteContact() {
  if (!confirm("Are you sure you want to delete this contact?")) return;

  router.delete(contactRoutes.destroy(props.contact.id).url);
}
</script>

<template>
  <AppLayout :title="`${contact.first_name} ${contact.last_name}`" :breadcrumbs="breadcrumbs">
    <div class="flex h-full flex-1 flex-col gap-6 p-4">
      <Card>
        <CardHeader class="flex flex-row items-start justify-between gap-4">
          <div class="space-y-2">
            <div class="flex items-center gap-2">
              <CardTitle class="text-2xl tracking-tight">
                {{ contact.first_name }} {{ contact.last_name }}
              </CardTitle>
              <Badge v-if="contact.organization" variant="outline">
                {{ contact.organization.name }}
              </Badge>
            </div>
            <div class="text-muted-foreground grid gap-1 text-sm">
              <p>{{ contact.email }}</p>
              <p>{{ contact.phone || "No phone number" }}</p>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <Button variant="outline" @click="toggleFavorite">
              {{ contact.is_favorite ? "Favorited" : "Favorite" }}
            </Button>
            <Button as-child>
              <Link :href="contactRoutes.edit(contact.id).url">Edit</Link>
            </Button>
            <Button variant="destructive" @click="deleteContact">Delete</Button>
          </div>
        </CardHeader>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Add note</CardTitle>
        </CardHeader>
        <CardContent class="space-y-3">
          <form class="space-y-3" @submit.prevent="submitNote">
            <Label for="note-body" class="sr-only">Note</Label>
            <Textarea
              id="note-body"
              v-model="noteForm.body"
              name="body"
              placeholder="Write a note for the activity feed..."
            />
            <InputError :message="noteForm.errors.body" />
            <Button type="submit" :disabled="noteForm.processing">
              {{ noteForm.processing ? "Saving..." : "Save note" }}
            </Button>
          </form>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Activity</CardTitle>
        </CardHeader>
        <CardContent>
          <Deferred data="notes">
            <template #fallback>
              <div class="space-y-3">
                <Skeleton class="h-20 w-full" />
                <Skeleton class="h-20 w-full" />
                <Skeleton class="h-20 w-full" />
              </div>
            </template>

            <div v-if="!notes || notes.length === 0" class="text-muted-foreground text-sm">
              No notes yet.
            </div>
            <div v-else class="space-y-3">
              <div v-for="note in notes" :key="note.id" class="rounded-lg bg-muted/30 p-4">
                <div class="mb-2 flex items-center justify-between gap-3">
                  <p class="text-sm font-medium">{{ note.user?.name ?? "Unknown" }}</p>
                  <time class="text-muted-foreground text-xs">{{
                    new Date(note.created_at).toLocaleString()
                  }}</time>
                </div>
                <p class="text-sm leading-6">{{ note.body }}</p>
              </div>
            </div>
          </Deferred>
        </CardContent>
      </Card>
    </div>
  </AppLayout>
</template>
