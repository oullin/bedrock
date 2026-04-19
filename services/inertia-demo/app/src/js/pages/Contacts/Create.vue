<script setup lang="ts">
import { useForm } from "@inertiajs/vue3";
import ContactForm from "./Partials/ContactForm.vue";
import type { ContactFormData, SelectOption } from "@/js/types";

const props = withDefaults(
  defineProps<{ form: ContactFormData; organizations?: SelectOption[] }>(),
  {
    organizations: () => [],
  },
);

const form = useForm({ ...props.form });
</script>

<template>
  <ContactForm mode="create" :form="form" :organizations="organizations" />
</template>
