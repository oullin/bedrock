<script setup lang="ts">
import type { InertiaForm } from "@inertiajs/vue3";
import { computed } from "vue";
import InputError from "@/js/components/app/InputError.vue";
import { Button } from "@/js/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/js/components/ui/card";
import { Input } from "@/js/components/ui/input";
import { Label } from "@/js/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/js/components/ui/select";
import AppLayout from "@/js/layouts/AppLayout.vue";
import { contactRoutes } from "@/js/lib/routes";
import type { Contact, ContactFormData, SelectOption } from "@/js/types";

const props = withDefaults(
  defineProps<{
    mode: "create" | "edit";
    form: InertiaForm<ContactFormData>;
    contact?: Contact | null;
    organizations?: SelectOption[];
  }>(),
  {
    contact: null,
    organizations: () => [],
  },
);

const isEdit = computed(() => props.mode === "edit");
const breadcrumbs = computed(() =>
  isEdit.value
    ? [
        { title: "CRM" },
        { title: "Contacts", href: contactRoutes.index().url },
        {
          title: `${props.contact?.first_name} ${props.contact?.last_name}`,
          href: contactRoutes.show(props.contact?.id).url,
        },
        { title: "Edit" },
      ]
    : [
        { title: "CRM" },
        { title: "Contacts", href: contactRoutes.index().url },
        { title: "Create" },
      ],
);

function submit() {
  if (isEdit.value) {
    props.form.put(contactRoutes.update(props.contact?.id).url);
    return;
  }

  props.form.post(contactRoutes.store().url);
}
</script>

<template>
  <AppLayout :title="isEdit ? 'Edit Contact' : 'Create Contact'" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4">
      <Card class="max-w-3xl">
        <CardHeader>
          <CardTitle>{{ isEdit ? "Edit Contact" : "Create Contact" }}</CardTitle>
          <CardDescription
            >Keep the CRM data aligned with the upstream demo structure.</CardDescription
          >
        </CardHeader>
        <CardContent>
          <form class="grid gap-6 md:grid-cols-2" @submit.prevent="submit">
            <div class="grid gap-2">
              <Label for="organization">Organization</Label>
              <Select v-model="form.organization_id">
                <SelectTrigger id="organization">
                  <SelectValue placeholder="Select an organization" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem
                    v-for="option in organizations"
                    :key="option.value || 'none'"
                    :value="option.value"
                  >
                    {{ option.label }}
                  </SelectItem>
                </SelectContent>
              </Select>
              <InputError :message="form.errors.organization_id" />
            </div>

            <div class="grid gap-2">
              <Label for="email">Email</Label>
              <Input id="email" v-model="form.email" type="email" />
              <InputError :message="form.errors.email" />
            </div>

            <div class="grid gap-2">
              <Label for="first_name">First name</Label>
              <Input id="first_name" v-model="form.first_name" />
              <InputError :message="form.errors.first_name" />
            </div>

            <div class="grid gap-2">
              <Label for="last_name">Last name</Label>
              <Input id="last_name" v-model="form.last_name" />
              <InputError :message="form.errors.last_name" />
            </div>

            <div class="grid gap-2">
              <Label for="phone">Phone</Label>
              <Input id="phone" v-model="form.phone" />
              <InputError :message="form.errors.phone" />
            </div>

            <div class="md:col-span-2 flex items-center gap-3">
              <Button type="submit" :disabled="form.processing">
                {{ form.processing ? "Saving..." : isEdit ? "Update contact" : "Create contact" }}
              </Button>
              <Button v-if="isEdit" type="button" variant="outline" as-child>
                <a :href="contactRoutes.show(contact?.id).url">Cancel</a>
              </Button>
              <span v-if="isEdit && form.isDirty" class="text-muted-foreground text-sm">
                Unsaved changes
              </span>
              <span v-if="isEdit && form.recentlySuccessful" class="text-sm text-green-600">
                Saved!
              </span>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  </AppLayout>
</template>
