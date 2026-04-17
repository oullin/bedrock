<template>
  <div class="flex h-full flex-1 flex-col">
    <div class="border-b bg-muted/30 px-4 py-2 lg:px-6">
      <p class="text-muted-foreground text-xs font-medium uppercase tracking-wide">
        Section Layout
      </p>
    </div>

    <div class="flex flex-1 flex-col">
      <slot />
    </div>
  </div>
</template>
