<script setup lang="ts">
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/js/components/ui/card";

withDefaults(defineProps<{ title: string; description?: string }>(), {
  description: "",
});
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>{{ title }}</CardTitle>
      <CardDescription v-if="description">{{ description }}</CardDescription>
    </CardHeader>
    <CardContent>
      <slot />
    </CardContent>
  </Card>
</template>
