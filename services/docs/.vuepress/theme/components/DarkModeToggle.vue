<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { Moon, Sun } from 'lucide-vue-next'
import { Button } from './ui/button'

const isDark = ref(false)

onMounted(() => {
  const saved = localStorage.getItem('bedrock-color-scheme')
  isDark.value =
    saved === 'dark' ||
    (!saved && window.matchMedia('(prefers-color-scheme: dark)').matches)
  document.documentElement.classList.toggle('dark', isDark.value)
})

function toggle() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('bedrock-color-scheme', isDark.value ? 'dark' : 'light')
}
</script>

<template>
  <Button
    variant="ghost"
    size="icon"
    :aria-label="isDark ? 'Switch to light mode' : 'Switch to dark mode'"
    class="text-zinc-500 hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-white"
    @click="toggle"
  >
    <Sun v-if="isDark" class="h-4 w-4" />
    <Moon v-else class="h-4 w-4" />
  </Button>
</template>
