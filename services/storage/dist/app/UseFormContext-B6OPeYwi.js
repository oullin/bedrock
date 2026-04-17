import{_ as n,a as d}from"./FeatureHeader.vue_vue_type_script_setup_true_lang-4CZFCx4c.js";import{_ as a}from"./AppLayout.vue_vue_type_script_setup_true_lang-Ct7coe5q.js";import{d as l,j as m,w as o,g as e,a as r,b as t,o as p}from"./app.js";import"./CardTitle.vue_vue_type_script_setup_true_lang-ucxQeZZT.js";import"./routes-S6IyX_wi.js";import"./CardDescription.vue_vue_type_script_setup_true_lang-dHSpRoiL.js";const c={class:"flex flex-1 flex-col gap-6 p-4 lg:p-6"},u={class:"grid gap-6 lg:grid-cols-2"},f={class:"grid gap-6 lg:grid-cols-2"},w=l({__name:"UseFormContext",setup(x){const i=[{title:"Features"},{title:"Forms"},{title:"Form Context"}];return(g,s)=>(p(),m(a,{title:"Form Context",breadcrumbs:i},{default:o(()=>[e("div",c,[r(n,{title:"useFormContext",description:"The useFormContext composable allows child components to access the parent form instance without prop drilling."}),e("div",u,[r(d,{title:"Concept",description:"How form context sharing works across components."},{default:o(()=>[...s[0]||(s[0]=[e("div",{class:"grid gap-3 text-sm"},[e("div",{class:"rounded-md border p-3"},[e("p",{class:"mb-1 font-medium"},"The Problem"),e("p",{class:"text-muted-foreground"}," When building complex forms with nested components, you need to pass the form instance through multiple layers of props. This creates tight coupling and verbose code. ")]),e("div",{class:"rounded-md border p-3"},[e("p",{class:"mb-1 font-medium"},"The Solution"),e("p",{class:"text-muted-foreground"},[e("code",{class:"bg-muted rounded px-1"},"useFormContext"),t(" uses Vue's provide/inject to make the form instance available to any descendant component without explicit prop passing. ")])])],-1)])]),_:1}),r(d,{title:"Component Structure",description:"A typical nested form component tree."},{default:o(()=>[...s[1]||(s[1]=[e("div",{class:"grid gap-2 text-sm"},[e("div",{class:"rounded-md border p-3"},[e("pre",{class:"overflow-auto text-xs leading-relaxed"},[e("code",null,`<!-- ParentForm.vue -->
<script setup>
import { useForm } from '@inertiajs/vue3'

const form = useForm({
  name: '',
  email: '',
  address: {
    street: '',
    city: '',
  },
})
<\/script>

<template>
  <form @submit.prevent="form.post('/submit')">
    <PersonalFields />
    <AddressFields />
    <button type="submit">Save</button>
  </form>
</template>`)])]),e("div",{class:"rounded-md border p-3"},[e("pre",{class:"overflow-auto text-xs leading-relaxed"},[e("code",null,`<!-- PersonalFields.vue -->
<script setup>
import { useFormContext } from '@inertiajs/vue3'

const form = useFormContext()
<\/script>

<template>
  <div>
    <input v-model="form.name" />
    <input v-model="form.email" />
  </div>
</template>`)])]),e("div",{class:"rounded-md border p-3"},[e("pre",{class:"overflow-auto text-xs leading-relaxed"},[e("code",null,`<!-- AddressFields.vue -->
<script setup>
import { useFormContext } from '@inertiajs/vue3'

const form = useFormContext()
<\/script>

<template>
  <div>
    <input v-model="form.address.street" />
    <input v-model="form.address.city" />
  </div>
</template>`)])])],-1)])]),_:1})]),e("div",f,[r(d,{title:"Benefits",description:"Why use form context."},{default:o(()=>[...s[2]||(s[2]=[e("div",{class:"grid gap-2 text-sm"},[e("div",{class:"flex items-start gap-3 rounded-md border p-3"},[e("span",{class:"bg-primary text-primary-foreground flex h-6 w-6 shrink-0 items-center justify-center rounded-full text-xs font-medium"},"1"),e("div",null,[e("p",{class:"font-medium"},"No Prop Drilling"),e("p",{class:"text-muted-foreground"}," Child components access the form directly without passing it through intermediate components. ")])]),e("div",{class:"flex items-start gap-3 rounded-md border p-3"},[e("span",{class:"bg-primary text-primary-foreground flex h-6 w-6 shrink-0 items-center justify-center rounded-full text-xs font-medium"},"2"),e("div",null,[e("p",{class:"font-medium"},"Reusable Field Components"),e("p",{class:"text-muted-foreground"}," Build form field components that work with any form instance. ")])]),e("div",{class:"flex items-start gap-3 rounded-md border p-3"},[e("span",{class:"bg-primary text-primary-foreground flex h-6 w-6 shrink-0 items-center justify-center rounded-full text-xs font-medium"},"3"),e("div",null,[e("p",{class:"font-medium"},"Clean Architecture"),e("p",{class:"text-muted-foreground"}," Keep form logic centralized while distributing UI across components. ")])])],-1)])]),_:1}),r(d,{title:"Usage Notes",description:"Important things to keep in mind."},{default:o(()=>[...s[3]||(s[3]=[e("div",{class:"grid gap-2 text-sm"},[e("div",{class:"rounded-md border p-3"},[e("p",{class:"text-muted-foreground"},[t(" The form context relies on Vue's "),e("code",{class:"bg-muted rounded px-1"},"provide"),t(" / "),e("code",{class:"bg-muted rounded px-1"},"inject"),t(" mechanism, so the child component must be a descendant of the component that created the form. ")])]),e("div",{class:"rounded-md border p-3"},[e("p",{class:"text-muted-foreground"},[t(" All reactive properties like "),e("code",{class:"bg-muted rounded px-1"},"processing"),t(", "),e("code",{class:"bg-muted rounded px-1"},"isDirty"),t(", and "),e("code",{class:"bg-muted rounded px-1"},"errors"),t(" remain reactive through the context. ")])]),e("div",{class:"rounded-md border p-3"},[e("p",{class:"text-muted-foreground"},[t(" Methods like "),e("code",{class:"bg-muted rounded px-1"},"reset()"),t(", "),e("code",{class:"bg-muted rounded px-1"},"clearErrors()"),t(", and "),e("code",{class:"bg-muted rounded px-1"},"post()"),t(" are also available via the context. ")])])],-1)])]),_:1})])])]),_:1}))}});export{w as default};
