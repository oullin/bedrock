<script setup lang="ts">
import { computed, ref } from 'vue'
import { RouteLink } from 'vuepress/client'
import { PACKAGES, TIERS } from './data'
import { ICON_MAP, IconSearch, IconArrow } from './icons'

const q = ref('')
const tier = ref<'all' | 'http' | 'data' | 'infra'>('all')

const filtered = computed(() => PACKAGES.filter(p => {
  if (tier.value !== 'all' && p.tier !== tier.value) return false
  if (!q.value) return true
  const s = q.value.toLowerCase()
  return p.name.toLowerCase().includes(s) || p.blurb.toLowerCase().includes(s) || p.tag.toLowerCase().includes(s)
}))

const grouped = computed(() =>
  TIERS.map(t => ({ ...t, items: filtered.value.filter(p => p.tier === t.id) }))
    .filter(g => g.items.length > 0),
)
</script>

<template>
  <section id="packages" class="pkgs">
    <div class="pkgs__inner">
      <div class="pkgs__header">
        <div>
          <div class="pkgs__label">The collection</div>
          <h2 class="pkgs__h">21 packages, <em class="serif">ready</em> to install individually.</h2>
          <p class="pkgs__sub">
            Every package has its own module path, docs and changelog. Pull just what you need.
          </p>
        </div>
        <div class="pkgs__controls">
          <div class="pkgs__search">
            <IconSearch :size="14" />
            <input v-model="q" placeholder="Filter…" />
          </div>
          <div class="pkgs__filters">
            <button
              :class="['pkgs__filter', tier === 'all' ? 'pkgs__filter--active' : '']"
              @click="tier = 'all'"
            >All</button>
            <button
              v-for="t in TIERS"
              :key="t.id"
              :class="['pkgs__filter', tier === t.id ? 'pkgs__filter--active' : '']"
              @click="tier = t.id"
            >{{ t.label }}</button>
          </div>
        </div>
      </div>

      <div v-if="grouped.length === 0" class="pkgs__empty">No packages match "{{ q }}"</div>

      <div v-for="g in grouped" :key="g.id" class="pkgs__tier">
        <div class="pkgs__tier-head">
          <span class="pkgs__tier-label">/ {{ g.id }}</span>
          <span class="pkgs__tier-name">{{ g.label }}</span>
          <span class="pkgs__tier-sub">{{ g.sub }} · {{ g.items.length }}</span>
        </div>
        <div class="pkgs__grid">
          <RouteLink v-for="p in g.items" :key="p.id" :to="p.href" class="pkg">
            <div class="pkg__top">
              <div class="pkg__icon">
                <component :is="ICON_MAP[p.icon]" :size="16" />
              </div>
              <div class="pkg__name">{{ p.name }}</div>
              <div class="pkg__v">v{{ p.v }}</div>
            </div>
            <p class="pkg__blurb">{{ p.blurb }}</p>
            <div class="pkg__meta">
              <span>↓ {{ p.downloads }}</span>
              <span>·</span>
              <span>since {{ p.since }}</span>
              <span class="pkg__cta">docs <IconArrow :size="11" /></span>
            </div>
          </RouteLink>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.pkgs {
  padding: 64px 32px 48px;
  background: linear-gradient(180deg, transparent 0%, color-mix(in oklab, var(--panel) 40%, transparent) 40%, transparent 100%);
  position: relative;
}
.pkgs__inner { max-width: var(--max); margin: 0 auto; }
.pkgs__header {
  display: flex; align-items: end; justify-content: space-between; gap: 40px;
  margin-bottom: 44px; flex-wrap: wrap;
}
.pkgs__label {
  display: inline-flex; align-items: center; gap: 10px;
  font-family: 'JetBrains Mono', monospace; font-size: 11px;
  color: var(--ink-3); text-transform: uppercase; letter-spacing: 0.2em;
  margin-bottom: 14px;
}
.pkgs__label::before { content: ''; width: 24px; height: 1px; background: var(--line-2); }
.pkgs__h {
  font-size: clamp(28px, 3.2vw, 40px); line-height: 1.05; letter-spacing: -0.025em;
  margin: 0 0 8px; font-weight: 600; max-width: 640px;
}
.pkgs__h em { font-family: 'Instrument Serif', serif; font-style: italic; font-weight: 400; color: var(--accent); }
.pkgs__sub { color: var(--ink-2); font-size: 15px; margin: 0; max-width: 520px; line-height: 1.5; }
.pkgs__controls {
  display: flex; gap: 8px; align-items: center; flex-wrap: wrap;
}
.pkgs__search {
  display: flex; align-items: center; gap: 8px;
  padding: 8px 12px;
  background: var(--panel);
  border: 1px solid var(--line);
  border-radius: 10px;
  min-width: 240px;
}
.pkgs__search:focus-within { border-color: color-mix(in oklab, var(--accent) 40%, var(--line)); }
.pkgs__search input {
  background: transparent; border: 0; outline: 0;
  color: var(--ink); font: inherit; font-size: 13px; flex: 1;
}
.pkgs__search input::placeholder { color: var(--ink-3); }
.pkgs__filters { display: flex; gap: 2px; padding: 3px; background: var(--panel); border: 1px solid var(--line); border-radius: 10px; }
.pkgs__filter {
  padding: 6px 12px; font-size: 13px; color: var(--ink-3); background: transparent;
  border: 0; border-radius: 7px; cursor: pointer; font-family: inherit;
  transition: color .12s, background .12s;
}
.pkgs__filter:hover { color: var(--ink-2); }
.pkgs__filter--active { background: var(--bg-2); color: var(--ink); }

.pkgs__tier { margin-bottom: 56px; }
.pkgs__tier-head {
  display: flex; align-items: baseline; gap: 14px;
  padding-bottom: 12px; margin-bottom: 20px;
  border-bottom: 1px solid var(--line);
}
.pkgs__tier-label {
  font-family: 'JetBrains Mono', monospace; font-size: 11px;
  color: var(--accent); letter-spacing: 0.14em; text-transform: uppercase;
}
.pkgs__tier-name { font-size: 22px; font-weight: 600; letter-spacing: -0.02em; }
.pkgs__tier-sub { color: var(--ink-3); font-size: 14px; margin-left: auto; font-family: 'JetBrains Mono', monospace; }

.pkgs__grid {
  display: grid; grid-template-columns: repeat(3, 1fr); gap: 1px;
  background: var(--line);
  border: 1px solid var(--line);
  border-radius: 14px;
  overflow: hidden;
}
.pkg {
  background: var(--panel);
  padding: 22px 22px 20px;
  display: flex; flex-direction: column; gap: 10px;
  position: relative;
  cursor: pointer;
  transition: background .15s;
  min-height: 172px;
}
.pkg:hover { background: var(--panel-2); }
.pkg::after {
  content: ''; position: absolute; inset: 0;
  border: 1px solid transparent; pointer-events: none;
  transition: border-color .15s;
}
.pkg:hover::after { border-color: color-mix(in oklab, var(--accent) 30%, transparent); }
.pkg__top { display: flex; align-items: center; gap: 10px; }
.pkg__icon {
  width: 32px; height: 32px; border-radius: 8px;
  display: grid; place-items: center;
  background: var(--bg-2);
  color: var(--ink-2);
  transition: color .15s, background .15s;
}
.pkg:hover .pkg__icon { color: var(--accent); background: color-mix(in oklab, var(--accent) 10%, var(--bg-2)); }
.pkg__name {
  font-family: 'JetBrains Mono', monospace; font-size: 14px; font-weight: 500; color: var(--ink);
  letter-spacing: -0.01em;
}
.pkg__v {
  margin-left: auto;
  font-family: 'JetBrains Mono', monospace; font-size: 11px; color: var(--ink-3);
}
.pkg__blurb {
  font-size: 13.5px; line-height: 1.55; color: var(--ink-2); margin: 0;
  flex: 1;
}
.pkg__meta {
  display: flex; align-items: center; gap: 10px;
  font-family: 'JetBrains Mono', monospace; font-size: 11px; color: var(--ink-3);
  padding-top: 10px;
  border-top: 1px dashed color-mix(in oklab, var(--line) 80%, transparent);
  opacity: 0.8;
}
.pkg__cta {
  margin-left: auto;
  display: inline-flex; align-items: center; gap: 4px;
  color: var(--ink-3); font-family: inherit; font-size: 12px;
  transition: color .15s, gap .2s;
}
.pkg:hover .pkg__cta { color: var(--accent); gap: 8px; }

.pkgs__empty { padding: 40px; text-align: center; color: var(--ink-3); font-size: 14px; }

@media (max-width: 960px) { .pkgs__grid { grid-template-columns: repeat(2, 1fr); } }
@media (max-width: 620px) { .pkgs__grid { grid-template-columns: 1fr; } }
</style>
