<script setup lang="ts">
import { onMounted, onBeforeUnmount, ref } from 'vue'
import { RouteLink } from 'vuepress/client'
import BrandMark from './BrandMark.vue'
import { IconSearch, IconGh, IconDiscord } from './icons'

const scrolled = ref(false)
const onScroll = () => { scrolled.value = window.scrollY > 12 }

onMounted(() => {
  window.addEventListener('scroll', onScroll, { passive: true })
  onScroll()
})
onBeforeUnmount(() => window.removeEventListener('scroll', onScroll))
</script>

<template>
  <header :class="['nav', scrolled ? 'nav--scrolled' : '']">
    <div class="nav__brand">
      <div class="nav__mark" aria-label="bedrock">
        <BrandMark :size="28" />
      </div>
      <span class="nav__brand-name">bedrock</span>
      <span class="nav__brand-v">v1.4</span>
    </div>

    <nav class="nav__links">
      <RouteLink class="nav__link nav__link--active" to="/">Home</RouteLink>
      <RouteLink class="nav__link" to="/getting-started">Get started</RouteLink>
      <a class="nav__link" href="#packages">Packages</a>
      <RouteLink class="nav__link" to="/concepts/request-lifecycle">Guides</RouteLink>
      <RouteLink class="nav__link" to="/packages/container">Reference</RouteLink>
    </nav>

    <div class="nav__spacer" />

    <div class="nav__search" role="searchbox" tabindex="0">
      <IconSearch :size="15" />
      <span>Search docs…</span>
      <span class="nav__kbd">⌘K</span>
    </div>

    <div class="nav__icons">
      <a class="nav__icon" href="#" aria-label="Discord"><IconDiscord :size="18" /></a>
      <a class="nav__icon" href="https://github.com/gocanto/bedrock" target="_blank" rel="noopener noreferrer" aria-label="GitHub"><IconGh :size="18" /></a>
    </div>
  </header>
</template>

<style scoped>
.nav {
  position: sticky; top: 0; z-index: 40;
  display: flex; align-items: center; justify-content: space-between;
  padding: 14px 32px;
  background: color-mix(in oklab, var(--bg) 70%, transparent);
  backdrop-filter: blur(14px) saturate(140%);
  -webkit-backdrop-filter: blur(14px) saturate(140%);
  border-bottom: 1px solid transparent;
  transition: border-color .2s ease, background .2s ease;
}
.nav--scrolled { border-bottom-color: var(--line); }
.nav__brand { display: flex; align-items: center; gap: 10px; font-weight: 700; letter-spacing: -0.01em; }
.nav__mark {
  width: 28px; height: 28px;
  display: grid; place-items: center;
  position: relative;
  filter: drop-shadow(0 4px 12px color-mix(in oklab, var(--accent) 50%, transparent));
}
.nav__mark svg { display: block; }
.nav__brand-name { font-size: 15px; }
.nav__brand-v { font-size: 11px; color: var(--ink-3); font-family: 'JetBrains Mono', monospace; padding: 3px 7px; border-radius: 999px; border: 1px solid var(--line); margin-left: 4px; }
.nav__links { display: flex; gap: 4px; margin-left: 28px; }
.nav__link { padding: 8px 12px; font-size: 13.5px; color: var(--ink-2); border-radius: 8px; transition: color .15s, background .15s; }
.nav__link:hover { color: var(--ink); background: var(--panel); }
.nav__link--active { color: var(--ink); background: var(--panel); }
.nav__spacer { flex: 1; }
.nav__search {
  display: flex; align-items: center; gap: 8px;
  padding: 7px 12px; min-width: 260px;
  border: 1px solid var(--line); border-radius: 10px;
  background: var(--panel); color: var(--ink-3); font-size: 13px;
  cursor: text; transition: border-color .15s;
}
.nav__search:hover { border-color: var(--line-2); }
.nav__kbd { margin-left: auto; font-family: 'JetBrains Mono', monospace; font-size: 11px; padding: 2px 6px; border-radius: 5px; background: var(--bg); border: 1px solid var(--line); color: var(--ink-3); }
.nav__icons { display: flex; gap: 4px; margin-left: 12px; }
.nav__icon {
  width: 36px; height: 36px; display: grid; place-items: center;
  border-radius: 9px; color: var(--ink-2);
  transition: color .15s, background .15s;
}
.nav__icon:hover { color: var(--ink); background: var(--panel); }

@media (max-width: 960px) {
  .nav__links, .nav__search { display: none; }
}
</style>
