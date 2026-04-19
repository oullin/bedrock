<script setup lang="ts">
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import AppLayout from "@/js/layouts/AppLayout.vue";

const breadcrumbs = [{ title: "Features" }, { title: "Forms" }, { title: "Form Context" }];
</script>

<template>
  <AppLayout title="Form Context" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="useFormContext"
        description="The useFormContext composable allows child components to access the parent form instance without prop drilling."
      />

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard
          title="Concept"
          description="How form context sharing works across components."
        >
          <div class="grid gap-3 text-sm">
            <div class="rounded-md border p-3">
              <p class="mb-1 font-medium">The Problem</p>
              <p class="text-muted-foreground">
                When building complex forms with nested components, you need to pass the form
                instance through multiple layers of props. This creates tight coupling and verbose
                code.
              </p>
            </div>
            <div class="rounded-md border p-3">
              <p class="mb-1 font-medium">The Solution</p>
              <p class="text-muted-foreground">
                <code class="bg-muted rounded px-1">useFormContext</code> uses Vue's provide/inject
                to make the form instance available to any descendant component without explicit
                prop passing.
              </p>
            </div>
          </div>
        </FeatureCard>

        <FeatureCard
          title="Component Structure"
          description="A typical nested form component tree."
        >
          <div class="grid gap-2 text-sm">
            <div class="rounded-md border p-3">
              <pre class="overflow-auto text-xs leading-relaxed"><code>&lt;!-- ParentForm.vue --&gt;
&lt;script setup&gt;
import { useForm } from '@inertiajs/vue3'

const form = useForm({
  name: '',
  email: '',
  address: {
    street: '',
    city: '',
  },
})
&lt;/script&gt;

&lt;template&gt;
  &lt;form @submit.prevent="form.post('/submit')"&gt;
    &lt;PersonalFields /&gt;
    &lt;AddressFields /&gt;
    &lt;button type="submit"&gt;Save&lt;/button&gt;
  &lt;/form&gt;
&lt;/template&gt;</code></pre>
            </div>

            <div class="rounded-md border p-3">
              <pre
                class="overflow-auto text-xs leading-relaxed"
              ><code>&lt;!-- PersonalFields.vue --&gt;
&lt;script setup&gt;
import { useFormContext } from '@inertiajs/vue3'

const form = useFormContext()
&lt;/script&gt;

&lt;template&gt;
  &lt;div&gt;
    &lt;input v-model="form.name" /&gt;
    &lt;input v-model="form.email" /&gt;
  &lt;/div&gt;
&lt;/template&gt;</code></pre>
            </div>

            <div class="rounded-md border p-3">
              <pre
                class="overflow-auto text-xs leading-relaxed"
              ><code>&lt;!-- AddressFields.vue --&gt;
&lt;script setup&gt;
import { useFormContext } from '@inertiajs/vue3'

const form = useFormContext()
&lt;/script&gt;

&lt;template&gt;
  &lt;div&gt;
    &lt;input v-model="form.address.street" /&gt;
    &lt;input v-model="form.address.city" /&gt;
  &lt;/div&gt;
&lt;/template&gt;</code></pre>
            </div>
          </div>
        </FeatureCard>
      </div>

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard title="Benefits" description="Why use form context.">
          <div class="grid gap-2 text-sm">
            <div class="flex items-start gap-3 rounded-md border p-3">
              <span
                class="bg-primary text-primary-foreground flex h-6 w-6 shrink-0 items-center justify-center rounded-full text-xs font-medium"
                >1</span
              >
              <div>
                <p class="font-medium">No Prop Drilling</p>
                <p class="text-muted-foreground">
                  Child components access the form directly without passing it through intermediate
                  components.
                </p>
              </div>
            </div>
            <div class="flex items-start gap-3 rounded-md border p-3">
              <span
                class="bg-primary text-primary-foreground flex h-6 w-6 shrink-0 items-center justify-center rounded-full text-xs font-medium"
                >2</span
              >
              <div>
                <p class="font-medium">Reusable Field Components</p>
                <p class="text-muted-foreground">
                  Build form field components that work with any form instance.
                </p>
              </div>
            </div>
            <div class="flex items-start gap-3 rounded-md border p-3">
              <span
                class="bg-primary text-primary-foreground flex h-6 w-6 shrink-0 items-center justify-center rounded-full text-xs font-medium"
                >3</span
              >
              <div>
                <p class="font-medium">Clean Architecture</p>
                <p class="text-muted-foreground">
                  Keep form logic centralized while distributing UI across components.
                </p>
              </div>
            </div>
          </div>
        </FeatureCard>

        <FeatureCard title="Usage Notes" description="Important things to keep in mind.">
          <div class="grid gap-2 text-sm">
            <div class="rounded-md border p-3">
              <p class="text-muted-foreground">
                The form context relies on Vue's
                <code class="bg-muted rounded px-1">provide</code> /
                <code class="bg-muted rounded px-1">inject</code> mechanism, so the child component
                must be a descendant of the component that created the form.
              </p>
            </div>
            <div class="rounded-md border p-3">
              <p class="text-muted-foreground">
                All reactive properties like <code class="bg-muted rounded px-1">processing</code>,
                <code class="bg-muted rounded px-1">isDirty</code>, and
                <code class="bg-muted rounded px-1">errors</code> remain reactive through the
                context.
              </p>
            </div>
            <div class="rounded-md border p-3">
              <p class="text-muted-foreground">
                Methods like <code class="bg-muted rounded px-1">reset()</code>,
                <code class="bg-muted rounded px-1">clearErrors()</code>, and
                <code class="bg-muted rounded px-1">post()</code> are also available via the
                context.
              </p>
            </div>
          </div>
        </FeatureCard>
      </div>
    </div>
  </AppLayout>
</template>
