<script setup lang="ts">
import HomeNav from '../home/HomeNav.vue'
import HomeHero from '../home/HomeHero.vue'
import HomeWhy from '../home/HomeWhy.vue'
import HomePackages from '../home/HomePackages.vue'
import HomeQuickstart from '../home/HomeQuickstart.vue'
import HomeFooter from '../home/HomeFooter.vue'
</script>

<template>
  <div class="bedrock-home">
    <HomeNav />
    <main>
      <HomeHero />
      <HomeWhy />
      <HomePackages />
      <HomeQuickstart />
    </main>
    <HomeFooter />
  </div>
</template>
