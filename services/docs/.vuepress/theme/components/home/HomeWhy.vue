<script setup lang="ts">
import { IconBolt, IconLayers, IconShield, IconStack } from './icons'

const items = [
  { icon: IconLayers,
    title: 'Use one. Use all.',
    body: "Each package is independent with its own version. No framework lock-in — adopt at your own pace." },
  { icon: IconBolt,
    title: 'Built for production.',
    body: "Zero-alloc hot paths, sensible defaults, pressure-tested in services handling 10k req/s." },
  { icon: IconShield,
    title: 'Security out of the box.',
    body: "Timing-safe crypto, CSRF-aware routing, rate limiting primitives. Secure defaults, always." },
  { icon: IconStack,
    title: 'Familiar patterns, idiomatic Go.',
    body: "Conventional, batteries-included APIs in idiomatic Go — interfaces, contexts, errors." },
]
</script>

<template>
  <section class="why">
    <div class="why__label">Why bedrock</div>
    <h2 class="why__h">
      A well-loved standard library for <em class="serif">real</em> web apps — the parts Go leaves up to you.
    </h2>

    <div class="why__grid">
      <div v-for="(it, i) in items" :key="i" class="why__cell">
        <div class="why__icon"><component :is="it.icon" :size="16" /></div>
        <h3 class="why__t">{{ it.title }}</h3>
        <p class="why__b">{{ it.body }}</p>
      </div>
    </div>
  </section>
</template>

<style scoped>
.why {
  padding: 88px 32px 48px;
  max-width: var(--max); margin: 0 auto;
}
.why__label {
  display: inline-flex; align-items: center; gap: 10px;
  font-family: 'JetBrains Mono', monospace; font-size: 11px;
  color: var(--ink-3); text-transform: uppercase; letter-spacing: 0.2em;
  margin-bottom: 18px;
}
.why__label::before { content: ''; width: 24px; height: 1px; background: var(--line-2); }
.why__h {
  font-size: clamp(28px, 3.2vw, 40px); line-height: 1.05; letter-spacing: -0.025em;
  margin: 0 0 64px; font-weight: 600; max-width: 720px;
}
.why__h em { font-family: 'Instrument Serif', serif; font-style: italic; font-weight: 400; color: var(--accent); }
.why__grid {
  display: grid; grid-template-columns: repeat(4, 1fr); gap: 0;
  border-top: 1px solid var(--line);
  border-left: 1px solid var(--line);
}
.why__cell {
  padding: 28px 24px;
  border-right: 1px solid var(--line);
  border-bottom: 1px solid var(--line);
  min-height: 180px;
  display: flex; flex-direction: column; gap: 12px;
  transition: background .2s;
}
.why__cell:hover { background: color-mix(in oklab, var(--panel) 50%, transparent); }
.why__icon {
  width: 36px; height: 36px; border-radius: 9px;
  display: grid; place-items: center;
  background: color-mix(in oklab, var(--accent) 10%, transparent);
  border: 1px solid color-mix(in oklab, var(--accent) 20%, transparent);
  color: var(--accent);
}
.why__t { font-size: 16px; font-weight: 600; letter-spacing: -0.01em; margin: 4px 0 0; }
.why__b { font-size: 13.5px; line-height: 1.6; color: var(--ink-2); margin: 0; }
@media (max-width: 860px) { .why__grid { grid-template-columns: repeat(2, 1fr); } }
</style>
