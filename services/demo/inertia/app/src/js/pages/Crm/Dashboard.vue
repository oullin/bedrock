<script setup lang="ts">
import { Deferred, Link, usePoll } from "@inertiajs/vue3";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/js/components/ui/card";
import { Skeleton } from "@/js/components/ui/skeleton";
import AppLayout from "@/js/layouts/AppLayout.vue";
import { appRoutes, contactRoutes } from "@/js/lib/routes";
import type { Note } from "@/js/types";

withDefaults(
  defineProps<{
    totalContacts?: number;
    totalOrganizations?: number;
    recentNotesCount?: number;
    recentActivity?: Note[];
  }>(),
  {
    totalContacts: undefined,
    totalOrganizations: undefined,
    recentNotesCount: undefined,
    recentActivity: () => [],
  },
);

usePoll(30000, { only: ["recentActivity"] });

const breadcrumbs = [{ title: "Dashboard", href: appRoutes.dashboard().url }];
</script>

<template>
  <AppLayout title="Dashboard" :breadcrumbs="breadcrumbs">
    <div class="flex h-full flex-1 flex-col gap-6 p-4">
      <div class="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader>
            <CardTitle class="text-sm font-medium">Total Contacts</CardTitle>
          </CardHeader>
          <CardContent>
            <Deferred data="totalContacts">
              <template #fallback>
                <Skeleton class="h-9 w-20" />
              </template>
              <div class="text-3xl font-semibold tracking-tight">{{ totalContacts }}</div>
            </Deferred>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle class="text-sm font-medium">Organizations</CardTitle>
          </CardHeader>
          <CardContent>
            <Deferred data="totalOrganizations">
              <template #fallback>
                <Skeleton class="h-9 w-20" />
              </template>
              <div class="text-3xl font-semibold tracking-tight">{{ totalOrganizations }}</div>
            </Deferred>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle class="text-sm font-medium">Notes This Week</CardTitle>
          </CardHeader>
          <CardContent>
            <Deferred data="recentNotesCount">
              <template #fallback>
                <Skeleton class="h-9 w-20" />
              </template>
              <div class="text-3xl font-semibold tracking-tight">{{ recentNotesCount }}</div>
            </Deferred>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Recent Activity</CardTitle>
          <CardDescription>Latest contact notes in the CRM.</CardDescription>
        </CardHeader>
        <CardContent>
          <div v-if="recentActivity.length === 0" class="text-muted-foreground py-4 text-sm">
            No recent activity.
          </div>
          <div v-else class="space-y-2">
            <div
              v-for="note in recentActivity"
              :key="note.id"
              class="flex items-start gap-4 rounded-lg bg-muted/30 px-4 py-3"
            >
              <div class="flex-1 space-y-1">
                <div class="flex items-center gap-2">
                  <span class="text-sm font-medium">{{ note.user?.name ?? "Unknown" }}</span>
                  <span class="text-muted-foreground text-xs">added a note on</span>
                  <Link
                    v-if="note.contact"
                    :href="contactRoutes.show(note.contact?.id).url"
                    class="text-sm font-medium text-primary hover:underline"
                  >
                    {{ note.contact?.name }}
                  </Link>
                  <span v-else class="text-sm font-medium text-muted-foreground"
                    >Unknown contact</span
                  >
                </div>
                <p class="text-muted-foreground line-clamp-2 text-sm">{{ note.body }}</p>
              </div>
              <time class="text-muted-foreground text-xs whitespace-nowrap">
                {{ new Date(note.created_at).toLocaleDateString() }}
              </time>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  </AppLayout>
</template>
