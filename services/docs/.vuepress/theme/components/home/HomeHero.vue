<script setup lang="ts">
import { computed, ref } from 'vue'
import { RouteLink } from 'vuepress/client'
import { HERO_TABS } from './data'
import { highlightGo } from './highlightGo'
import { IconArrow, IconGh, IconStar, IconCopy, IconCheck, IconDot } from './icons'

const tab = ref(HERO_TABS[0].id)
const copied = ref(false)

const current = computed(() => HERO_TABS.find(t => t.id === tab.value)!)
const highlighted = computed(() => highlightGo(current.value.code))
const lineCount = computed(() => current.value.code.split('\n').length)

function copy() {
  navigator.clipboard?.writeText(current.value.code)
  copied.value = true
  setTimeout(() => { copied.value = false }, 1400)
}
</script>

<template>
  <section class="hero">
    <div class="hero__inner">
      <div>
        <div class="hero__tag">
          <span class="hero__tag-dot" />
          v1.4.0 — now with <code style="font-size: 11px; opacity: 0.9">bedrock/bus</code> and typed events
        </div>

        <h1 class="hero__h1">
          The Upstream <em class="serif">feel</em>.<br />
          The Go performance.
        </h1>

        <p class="hero__sub">
          A collection of 21 well-tested, independently versioned packages for building
          production-grade Go web applications. Pick one, or use them together.
        </p>

        <div class="hero__ctas">
          <RouteLink to="/getting-started" class="btn btn--primary">
            Get started <IconArrow :size="15" />
          </RouteLink>
          <a
            class="btn btn--ghost"
            href="https://github.com/gocanto/bedrock"
            target="_blank"
            rel="noopener noreferrer"
          >
            <IconGh :size="15" /> View on GitHub <IconStar :size="12" /> 4.8k
          </a>
        </div>

        <div class="hero__stats">
          <div>
            <div class="hero__stat-n">21</div>
            <div class="hero__stat-l">packages</div>
          </div>
          <div>
            <div class="hero__stat-n">Go 1.24+</div>
            <div class="hero__stat-l">required</div>
          </div>
          <div>
            <div class="hero__stat-n">9.2M</div>
            <div class="hero__stat-l">downloads</div>
          </div>
          <div>
            <div class="hero__stat-n">MIT</div>
            <div class="hero__stat-l">license</div>
          </div>
        </div>
      </div>

      <div>
        <div class="code">
          <div class="code__bar">
            <div class="code__lights"><span /><span /><span /></div>
            <div class="code__tabs">
              <button
                v-for="t in HERO_TABS"
                :key="t.id"
                :class="['code__tab', tab === t.id ? 'code__tab--active' : '']"
                @click="tab = t.id"
              >
                {{ t.label }}
              </button>
            </div>
            <button class="code__copy" @click="copy">
              <template v-if="copied"><IconCheck :size="11" /> copied</template>
              <template v-else><IconCopy :size="11" /> copy</template>
            </button>
          </div>
          <div class="code__body">
            <div class="code__gutter">
              <div v-for="n in lineCount" :key="n">{{ n }}</div>
            </div>
            <pre class="code__pre" v-html="highlighted" />
          </div>
        </div>

        <div class="term">
          <div class="term__head">
            <IconDot :size="10" style="color:#7ee07e" />
            <span>go run ./…</span>
          </div>
          <div><span class="term__prompt">$</span> go run ./cmd/api</div>
          <div class="term__info">→ bedrock/container  wiring 14 bindings (&lt;1ms)</div>
          <div class="term__info">→ bedrock/routing    34 routes registered</div>
          <div class="term__info">→ bedrock/queue      worker pool ready (4 workers)</div>
          <div class="term__ok">✓ listening on :8080 — ready in 12ms<span class="term__cursor" /></div>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.hero {
  position: relative;
  padding: 72px 32px 96px;
  overflow: hidden;
  isolation: isolate;
}
.hero::before {
  content: '';
  position: absolute; inset: -20% -10% auto -10%; height: 680px; z-index: -1;
  background:
    radial-gradient(60% 55% at 15% 30%, color-mix(in oklab, var(--accent) 18%, transparent) 0%, transparent 55%),
    radial-gradient(45% 50% at 85% 10%, color-mix(in oklab, #7b6bff 20%, transparent) 0%, transparent 55%);
  filter: blur(20px);
  pointer-events: none;
}
.hero::after {
  content: '';
  position: absolute; inset: 0; z-index: -1;
  background-image:
    linear-gradient(to right, color-mix(in oklab, var(--line) 60%, transparent) 1px, transparent 1px),
    linear-gradient(to bottom, color-mix(in oklab, var(--line) 60%, transparent) 1px, transparent 1px);
  background-size: 56px 56px;
  background-position: -1px -1px;
  mask-image: radial-gradient(70% 60% at 50% 30%, black 40%, transparent 90%);
  -webkit-mask-image: radial-gradient(70% 60% at 50% 30%, black 40%, transparent 90%);
  opacity: .35;
  pointer-events: none;
}
.hero__inner {
  max-width: var(--max); margin: 0 auto;
  display: grid; grid-template-columns: minmax(0, 1fr) minmax(0, 1.05fr);
  gap: 56px; align-items: start;
}
.hero__tag {
  display: inline-flex; align-items: center; gap: 8px;
  padding: 6px 12px 6px 8px;
  border-radius: 999px;
  background: color-mix(in oklab, var(--accent) 9%, transparent);
  border: 1px solid color-mix(in oklab, var(--accent) 22%, transparent);
  color: var(--accent);
  font-size: 12px; font-weight: 500;
  margin-bottom: 28px;
}
.hero__tag-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--accent); box-shadow: 0 0 0 3px color-mix(in oklab, var(--accent) 30%, transparent); }
.hero__h1 {
  font-size: clamp(56px, 7vw, 104px);
  line-height: 0.98; letter-spacing: -0.045em;
  margin: 0 0 24px; font-weight: 800;
}
.hero__h1 em { font-family: 'Instrument Serif', serif; font-style: italic; font-weight: 400; color: var(--accent); letter-spacing: -0.01em; }
.hero__sub {
  font-size: 18px; line-height: 1.55; color: var(--ink-2);
  max-width: 520px; margin: 0 0 36px;
}
.hero__ctas { display: flex; gap: 10px; flex-wrap: wrap; }
.btn {
  display: inline-flex; align-items: center; gap: 8px;
  padding: 12px 18px;
  font-size: 14px; font-weight: 600;
  border-radius: 10px; cursor: pointer; border: 0;
  transition: transform .1s ease, background .15s;
}
.btn:active { transform: translateY(1px); }
.btn--primary {
  background: var(--accent); color: var(--accent-ink);
  box-shadow: 0 0 0 1px color-mix(in oklab, var(--accent) 50%, transparent), 0 12px 40px -12px color-mix(in oklab, var(--accent) 70%, transparent);
}
.btn--primary:hover { background: color-mix(in oklab, var(--accent) 90%, white); }
.btn--ghost {
  background: var(--panel); color: var(--ink); border: 1px solid var(--line);
}
.btn--ghost:hover { background: var(--panel-2); border-color: var(--line-2); }
.hero__stats {
  display: grid; grid-template-columns: repeat(4, auto); gap: 36px;
  margin-top: 48px;
  padding-top: 28px;
  border-top: 1px solid var(--line);
  max-width: 560px;
}
.hero__stat-n {
  font-size: 26px; font-weight: 700; letter-spacing: -0.02em;
  color: var(--ink);
}
.hero__stat-l { font-size: 12px; color: var(--ink-3); margin-top: 2px; text-transform: uppercase; letter-spacing: 0.08em; }

.code {
  border-radius: 14px;
  background: linear-gradient(180deg, #0b1220 0%, #080d17 100%);
  border: 1px solid var(--line);
  box-shadow:
    0 1px 0 color-mix(in oklab, white 5%, transparent) inset,
    0 30px 80px -30px rgba(0,0,0,.8),
    0 0 0 1px color-mix(in oklab, var(--accent) 8%, transparent);
  overflow: hidden;
  position: relative;
}
.code__bar {
  display: flex; align-items: center; gap: 10px;
  padding: 10px 14px;
  border-bottom: 1px solid var(--line);
  background: color-mix(in oklab, black 35%, var(--panel));
}
.code__lights { display: flex; gap: 6px; }
.code__lights span { width: 11px; height: 11px; border-radius: 50%; display: block; }
.code__lights span:nth-child(1) { background:#ff5f57; }
.code__lights span:nth-child(2) { background:#febc2e; }
.code__lights span:nth-child(3) { background:#28c840; }
.code__tabs { display: flex; gap: 2px; margin-left: 10px; flex: 1; }
.code__tab {
  padding: 6px 12px; font-family: 'JetBrains Mono', monospace; font-size: 12px;
  color: var(--ink-3); background: transparent; border: 0; border-radius: 6px;
  cursor: pointer; transition: color .12s, background .12s;
}
.code__tab:hover { color: var(--ink-2); }
.code__tab--active { color: var(--accent); background: color-mix(in oklab, var(--accent) 10%, transparent); }
.code__copy {
  margin-left: auto; display: inline-flex; align-items: center; gap: 6px;
  padding: 5px 10px; font-size: 11px; font-family: 'JetBrains Mono', monospace;
  color: var(--ink-3); background: transparent; border: 1px solid var(--line); border-radius: 6px;
  cursor: pointer; transition: color .12s, border-color .12s;
}
.code__copy:hover { color: var(--ink); border-color: var(--line-2); }
.code__body {
  display: grid; grid-template-columns: 44px 1fr; gap: 0;
  font-family: 'JetBrains Mono', monospace; font-size: 12.5px; line-height: 1.72;
}
.code__gutter {
  padding: 16px 10px 20px; text-align: right; color: var(--ink-3); opacity: .5;
  border-right: 1px solid var(--line);
  user-select: none;
}
.code__pre {
  padding: 16px 20px 20px; margin: 0; overflow: auto;
  color: var(--ink);
}
.code__pre :deep(.k) { color: #f07178; }
.code__pre :deep(.s) { color: #c3e88d; }
.code__pre :deep(.c) { color: var(--ink-3); font-style: italic; }
.code__pre :deep(.fn) { color: #82aaff; }
.code__pre :deep(.t) { color: #ffcb6b; }
.code__pre :deep(.p) { color: #89ddff; }

.term {
  margin-top: 14px;
  border-radius: 12px;
  background: #05080d;
  border: 1px solid var(--line);
  padding: 14px 16px;
  font-family: 'JetBrains Mono', monospace; font-size: 12px; line-height: 1.85;
  color: var(--ink-2);
  position: relative;
}
.term__head {
  display: flex; align-items: center; gap: 8px;
  color: var(--ink-3); margin-bottom: 6px;
  font-size: 11px; text-transform: uppercase; letter-spacing: 0.12em;
}
.term__prompt { color: var(--accent); }
.term__ok { color: #7ee07e; }
.term__info { color: #89ddff; }
.term__cursor {
  display: inline-block; width: 7px; height: 13px; vertical-align: -2px;
  background: var(--accent); animation: blink 1s steps(1) infinite;
  margin-left: 2px;
}
@keyframes blink { 50% { opacity: 0; } }

@media (max-width: 960px) {
  .hero__inner { grid-template-columns: 1fr; gap: 40px; }
  .hero__stats { grid-template-columns: repeat(2, auto); gap: 24px; }
}
</style>
