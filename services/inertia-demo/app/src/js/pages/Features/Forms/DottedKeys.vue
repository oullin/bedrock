<script setup lang="ts">
import { useForm, usePage } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import InputError from "@/js/components/app/InputError.vue";
import { Button } from "@/js/components/ui/button";
import { Input } from "@/js/components/ui/input";
import { Label } from "@/js/components/ui/label";
import AppLayout from "@/js/layouts/AppLayout.vue";
import type { SharedPageProps } from "@/js/types";

const page = usePage<SharedPageProps>();

const breadcrumbs = [{ title: "Features" }, { title: "Forms" }, { title: "Dotted Keys" }];

const form = useForm({
  address: {
    street: "",
    city: "",
    state: "",
    zip: "",
  },
  contact: {
    emails: ["", ""],
    phone: "",
  },
});

function submit() {
  form.post(page.url);
}
</script>

<template>
  <AppLayout title="Dotted Keys" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Dotted Keys"
        description="Inertia supports dotted key notation for nested form fields. This is useful for flat form structures that map to nested server-side data."
      />

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard
          title="Nested Form Fields"
          description="Uses dotted notation to represent nested data structures."
        >
          <form class="grid gap-4" @submit.prevent="submit">
            <div>
              <p class="mb-3 text-sm font-medium">Address</p>
              <div class="grid gap-3">
                <div class="grid gap-2">
                  <Label for="street">address.street</Label>
                  <Input id="street" v-model="form.address.street" placeholder="123 Main St" />
                  <InputError :message="form.errors['address.street']" />
                </div>

                <div class="grid grid-cols-2 gap-3">
                  <div class="grid gap-2">
                    <Label for="city">address.city</Label>
                    <Input id="city" v-model="form.address.city" placeholder="Springfield" />
                    <InputError :message="form.errors['address.city']" />
                  </div>
                  <div class="grid gap-2">
                    <Label for="state">address.state</Label>
                    <Input id="state" v-model="form.address.state" placeholder="IL" />
                    <InputError :message="form.errors['address.state']" />
                  </div>
                </div>

                <div class="grid gap-2">
                  <Label for="zip">address.zip</Label>
                  <Input id="zip" v-model="form.address.zip" placeholder="62701" />
                  <InputError :message="form.errors['address.zip']" />
                </div>
              </div>
            </div>

            <div>
              <p class="mb-3 text-sm font-medium">Contact</p>
              <div class="grid gap-3">
                <div class="grid gap-2">
                  <Label for="email0">contact.emails[0]</Label>
                  <Input
                    id="email0"
                    v-model="form.contact.emails[0]"
                    type="email"
                    placeholder="primary@example.com"
                  />
                  <InputError :message="form.errors['contact.emails[0]']" />
                </div>

                <div class="grid gap-2">
                  <Label for="email1">contact.emails[1]</Label>
                  <Input
                    id="email1"
                    v-model="form.contact.emails[1]"
                    type="email"
                    placeholder="secondary@example.com"
                  />
                  <InputError :message="form.errors['contact.emails[1]']" />
                </div>

                <div class="grid gap-2">
                  <Label for="phone">contact.phone</Label>
                  <Input id="phone" v-model="form.contact.phone" placeholder="(555) 123-4567" />
                  <InputError :message="form.errors['contact.phone']" />
                </div>
              </div>
            </div>

            <div class="flex items-center gap-3">
              <Button type="submit" :disabled="form.processing">
                {{ form.processing ? "Saving..." : "Save" }}
              </Button>
              <Button type="button" variant="outline" @click="form.reset()">Reset</Button>
              <span v-if="form.recentlySuccessful" class="text-sm text-green-600">Saved!</span>
            </div>
          </form>
        </FeatureCard>

        <div class="grid gap-6">
          <FeatureCard
            title="Form Data"
            description="The current form data using dotted key notation."
          >
            <div class="text-sm">
              <pre class="bg-muted overflow-auto rounded-md p-3 text-xs">{{
                JSON.stringify(form.data(), null, 2)
              }}</pre>
            </div>
          </FeatureCard>

          <FeatureCard
            title="How Dotted Keys Work"
            description="Understanding the flat-to-nested mapping."
          >
            <div class="grid gap-2 text-sm">
              <div class="rounded-md border p-3">
                <p class="mb-1 font-medium">Flat Keys</p>
                <p class="text-muted-foreground">
                  Form fields use flat dotted notation like
                  <code class="bg-muted rounded px-1">address.street</code> and
                  <code class="bg-muted rounded px-1">contact.emails[0]</code>.
                </p>
              </div>
              <div class="rounded-md border p-3">
                <p class="mb-1 font-medium">Server Flattens to Dotted Keys</p>
                <p class="text-muted-foreground">
                  Inertia sends the nested JSON shown above. The server's
                  <code class="bg-muted rounded px-1">ParseForm</code> flattens it back to dotted
                  keys so handlers access values via
                  <code class="bg-muted rounded px-1">r.FormValue("address.street")</code>.
                </p>
                <pre class="bg-muted mt-2 overflow-auto rounded p-2 text-xs">
// Wire format (nested JSON)
{ "address": { "street": "123 Main St" } }

// After ParseForm (flat dotted keys)
r.FormValue("address.street") // "123 Main St"</pre
                >
              </div>
              <div class="rounded-md border p-3">
                <p class="mb-1 font-medium">Error Mapping</p>
                <p class="text-muted-foreground">
                  Server validation errors use the same dotted notation, so
                  <code class="bg-muted rounded px-1">form.errors['address.street']</code>
                  maps correctly.
                </p>
              </div>
            </div>
          </FeatureCard>
        </div>
      </div>
    </div>
  </AppLayout>
</template>
