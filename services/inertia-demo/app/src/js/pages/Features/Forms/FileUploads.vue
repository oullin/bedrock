<script setup lang="ts">
import { computed, ref } from "vue";
import { useForm, usePage } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import InputError from "@/js/components/app/InputError.vue";
import { Button } from "@/js/components/ui/button";
import { Label } from "@/js/components/ui/label";
import AppLayout from "@/js/layouts/AppLayout.vue";
import type { SharedPageProps } from "@/js/types";

const page = usePage<SharedPageProps>();

const breadcrumbs = [{ title: "Features" }, { title: "Forms" }, { title: "File Uploads" }];

const photoInput = ref<HTMLInputElement | null>(null);
const documentsInput = ref<HTMLInputElement | null>(null);

const form = useForm({
  photo: null as File | null,
  documents: [] as File[],
});

function formatSize(bytes: number): string {
  if (!bytes) return "0 B";
  const units = ["B", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(1024));
  return `${(bytes / Math.pow(1024, i)).toFixed(1)} ${units[i]}`;
}

const photoName = computed(() => (form.photo ? form.photo.name : null));
const photoSize = computed(() => (form.photo ? formatSize(form.photo.size) : null));

const documentList = computed(() =>
  Array.from(form.documents).map((f: File) => ({
    name: f.name,
    size: formatSize(f.size),
  })),
);

function onPhotoChange(e: Event) {
  form.photo = (e.target as HTMLInputElement).files?.[0] || null;
}

function onDocumentsChange(e: Event) {
  form.documents = Array.from((e.target as HTMLInputElement).files || []);
}

function reset() {
  form.reset();
  if (photoInput.value) photoInput.value.value = "";
  if (documentsInput.value) documentsInput.value.value = "";
}

function submit() {
  form.post(page.url, {
    forceFormData: true,
  });
}
</script>

<template>
  <AppLayout title="File Uploads" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="File Uploads"
        description="Upload files with Inertia.js using multipart form data. The useForm helper automatically detects file inputs and sets the correct encoding."
      />

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard title="Upload Form" description="Select files and submit to upload.">
          <form class="grid gap-4" @submit.prevent="submit">
            <div class="grid gap-2">
              <Label for="photo">Photo (single file)</Label>
              <input
                id="photo"
                ref="photoInput"
                type="file"
                accept="image/*"
                class="border-input bg-background file:bg-primary file:text-primary-foreground hover:file:bg-primary/90 flex w-full rounded-md border px-3 py-2 text-sm file:mr-3 file:rounded file:border-0 file:px-3 file:py-1 file:text-sm file:font-medium"
                @change="onPhotoChange"
              />
              <InputError :message="form.errors.photo" />
            </div>

            <div class="grid gap-2">
              <Label for="documents">Documents (multiple files)</Label>
              <input
                id="documents"
                ref="documentsInput"
                type="file"
                multiple
                class="border-input bg-background file:bg-primary file:text-primary-foreground hover:file:bg-primary/90 flex w-full rounded-md border px-3 py-2 text-sm file:mr-3 file:rounded file:border-0 file:px-3 file:py-1 file:text-sm file:font-medium"
                @change="onDocumentsChange"
              />
              <InputError :message="form.errors.documents" />
            </div>

            <div v-if="form.progress" class="grid gap-1">
              <div class="flex items-center justify-between text-sm">
                <span>Uploading...</span>
                <span>{{ form.progress.percentage }}%</span>
              </div>
              <div class="bg-muted h-2 overflow-hidden rounded-full">
                <div
                  class="bg-primary h-full rounded-full transition-all duration-300"
                  :style="{ width: `${form.progress.percentage}%` }"
                />
              </div>
            </div>

            <div class="flex items-center gap-3">
              <Button type="submit" :disabled="form.processing">
                {{ form.processing ? "Uploading..." : "Upload Files" }}
              </Button>
              <Button type="button" variant="outline" @click="reset()">Clear</Button>
              <span v-if="form.recentlySuccessful" class="text-sm text-green-600"> Uploaded! </span>
            </div>
          </form>
        </FeatureCard>

        <FeatureCard title="Selected Files" description="Preview of files selected for upload.">
          <div class="grid gap-3 text-sm">
            <div class="rounded-md border p-3">
              <p class="mb-2 font-medium">Photo</p>
              <div v-if="photoName" class="flex items-center justify-between">
                <span class="truncate">{{ photoName }}</span>
                <span class="text-muted-foreground ml-2 shrink-0">{{ photoSize }}</span>
              </div>
              <p v-else class="text-muted-foreground italic">No photo selected</p>
            </div>

            <div class="rounded-md border p-3">
              <p class="mb-2 font-medium">Documents ({{ documentList.length }})</p>
              <div v-if="documentList.length" class="grid gap-1">
                <div
                  v-for="(doc, i) in documentList"
                  :key="i"
                  class="flex items-center justify-between rounded bg-gray-50 px-2 py-1"
                >
                  <span class="truncate">{{ doc.name }}</span>
                  <span class="text-muted-foreground ml-2 shrink-0">{{ doc.size }}</span>
                </div>
              </div>
              <p v-else class="text-muted-foreground italic">No documents selected</p>
            </div>

            <div class="rounded-md border p-3">
              <p class="mb-1 font-medium">Upload Progress</p>
              <pre class="bg-muted overflow-auto rounded p-2 text-xs">{{
                JSON.stringify(form.progress, null, 2)
              }}</pre>
            </div>
          </div>
        </FeatureCard>
      </div>
    </div>
  </AppLayout>
</template>
