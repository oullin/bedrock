<script setup lang="ts">
import type { AsTag, PrimitiveProps } from "reka-ui";
import type { HTMLAttributes } from "vue";
import type { SidebarMenuButtonVariants } from ".";
import { Primitive } from "reka-ui";
import { cn } from "@/js/lib/utils";
import { sidebarMenuButtonVariants } from ".";

export interface SidebarMenuButtonProps extends Omit<PrimitiveProps, "as"> {
  as?: AsTag;
  variant?: SidebarMenuButtonVariants["variant"];
  size?: SidebarMenuButtonVariants["size"];
  isActive?: boolean;
  class?: HTMLAttributes["class"];
}

const props = withDefaults(defineProps<SidebarMenuButtonProps>(), {
  as: "button",
  variant: "default",
  size: "default",
});
</script>

<template>
  <Primitive
    data-slot="sidebar-menu-button"
    data-sidebar="menu-button"
    :data-size="size"
    :data-active="isActive"
    :class="cn(sidebarMenuButtonVariants({ variant, size }), props.class)"
    :as="as"
    :as-child="asChild"
    v-bind="$attrs"
  >
    <slot />
  </Primitive>
</template>
