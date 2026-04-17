<script setup lang="ts">
import { useForm, usePage } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import InputError from "@/js/components/app/InputError.vue";
import { Button } from "@/js/components/ui/button";
import { Input } from "@/js/components/ui/input";
import { Textarea } from "@/js/components/ui/textarea";
import { Label } from "@/js/components/ui/label";
import { Checkbox } from "@/js/components/ui/checkbox";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/js/components/ui/select";
import AppLayout from "@/js/layouts/AppLayout.vue";
import type { SharedPageProps } from "@/js/types";

const page = usePage<SharedPageProps>();

const breadcrumbs = [{ title: "Features" }, { title: "Forms" }, { title: "useForm" }];

const form = useForm({
  name: "",
  email: "",
  bio: "",
  role: "",
  subscribe: false,
});

function submit() {
  form.post(page.url);
}
</script>

<template>
  <AppLayout title="useForm" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="useForm"
        description="The useForm helper provides a convenient way to manage form state, validation errors, and submission in Inertia.js applications."
      />

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard title="Form" description="Submit this form to see useForm in action.">
          <form class="grid gap-4" @submit.prevent="submit">
            <div class="grid gap-2">
              <Label for="name">Name</Label>
              <Input id="name" v-model="form.name" placeholder="John Doe" />
              <InputError :message="form.errors.name" />
            </div>

            <div class="grid gap-2">
              <Label for="email">Email</Label>
              <Input id="email" v-model="form.email" type="email" placeholder="john@example.com" />
              <InputError :message="form.errors.email" />
            </div>

            <div class="grid gap-2">
              <Label for="bio">Bio</Label>
              <Textarea
                id="bio"
                v-model="form.bio"
                rows="3"
                placeholder="Tell us about yourself..."
              />
              <InputError :message="form.errors.bio" />
            </div>

            <div class="grid gap-2">
              <Label for="role">Role</Label>
              <Select v-model="form.role">
                <SelectTrigger id="role">
                  <SelectValue placeholder="Select a role" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="admin">Admin</SelectItem>
                  <SelectItem value="editor">Editor</SelectItem>
                  <SelectItem value="viewer">Viewer</SelectItem>
                </SelectContent>
              </Select>
              <InputError :message="form.errors.role" />
            </div>

            <div class="flex items-center gap-2">
              <Checkbox
                id="subscribe"
                :checked="form.subscribe"
                @update:checked="form.subscribe = Boolean($event)"
              />
              <Label for="subscribe" class="cursor-pointer">Subscribe to newsletter</Label>
            </div>

            <div class="flex items-center gap-3">
              <Button type="submit" :disabled="form.processing">
                {{ form.processing ? "Submitting..." : "Submit" }}
              </Button>
              <Button type="button" variant="outline" @click="form.reset()">Reset</Button>
              <Button type="button" variant="ghost" @click="form.clearErrors()">
                Clear Errors
              </Button>
            </div>
          </form>
        </FeatureCard>

        <FeatureCard
          title="Reactive Properties"
          description="These properties update in real-time as you interact with the form."
        >
          <div class="grid gap-3 text-sm">
            <div class="flex items-center justify-between rounded-md border p-3">
              <span class="font-medium">processing</span>
              <span
                class="rounded-full px-2 py-0.5 text-xs font-medium"
                :class="
                  form.processing ? 'bg-yellow-100 text-yellow-800' : 'bg-gray-100 text-gray-600'
                "
              >
                {{ form.processing }}
              </span>
            </div>

            <div class="flex items-center justify-between rounded-md border p-3">
              <span class="font-medium">isDirty</span>
              <span
                class="rounded-full px-2 py-0.5 text-xs font-medium"
                :class="
                  form.isDirty ? 'bg-orange-100 text-orange-800' : 'bg-gray-100 text-gray-600'
                "
              >
                {{ form.isDirty }}
              </span>
            </div>

            <div class="flex items-center justify-between rounded-md border p-3">
              <span class="font-medium">hasErrors</span>
              <span
                class="rounded-full px-2 py-0.5 text-xs font-medium"
                :class="form.hasErrors ? 'bg-red-100 text-red-800' : 'bg-gray-100 text-gray-600'"
              >
                {{ form.hasErrors }}
              </span>
            </div>

            <div class="flex items-center justify-between rounded-md border p-3">
              <span class="font-medium">wasSuccessful</span>
              <span
                class="rounded-full px-2 py-0.5 text-xs font-medium"
                :class="
                  form.wasSuccessful ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-600'
                "
              >
                {{ form.wasSuccessful }}
              </span>
            </div>

            <div class="flex items-center justify-between rounded-md border p-3">
              <span class="font-medium">recentlySuccessful</span>
              <span
                class="rounded-full px-2 py-0.5 text-xs font-medium"
                :class="
                  form.recentlySuccessful
                    ? 'bg-green-100 text-green-800'
                    : 'bg-gray-100 text-gray-600'
                "
              >
                {{ form.recentlySuccessful }}
              </span>
            </div>

            <div class="mt-2 rounded-md border p-3">
              <p class="mb-1 font-medium">Current Data</p>
              <pre class="bg-muted overflow-auto rounded p-2 text-xs">{{
                JSON.stringify(form.data(), null, 2)
              }}</pre>
            </div>

            <div v-if="form.hasErrors" class="rounded-md border border-red-200 bg-red-50 p-3">
              <p class="mb-1 font-medium text-red-800">Errors</p>
              <pre class="overflow-auto text-xs text-red-700">{{
                JSON.stringify(form.errors, null, 2)
              }}</pre>
            </div>
          </div>
        </FeatureCard>
      </div>
    </div>
  </AppLayout>
</template>
