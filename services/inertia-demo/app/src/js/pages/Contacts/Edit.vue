<script setup lang="ts">
import { useForm } from "@inertiajs/vue3";
import ContactForm from "./Partials/ContactForm.vue";
import type { Contact, ContactFormData, SelectOption } from "@/js/types";

const props = withDefaults(
  defineProps<{ contact: Contact; form: ContactFormData; organizations?: SelectOption[] }>(),
  {
    organizations: () => [],
  },
);

const form = useForm({ ...props.form });
</script>

<template>
  <ContactForm mode="edit" :contact="contact" :form="form" :organizations="organizations" />
</template>
