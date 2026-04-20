<script setup lang="ts">
import { ref } from 'vue'
import { IconCopy, IconCheck } from './icons'

const copied = ref<string | null>(null)

function copy(txt: string, key: string) {
  navigator.clipboard?.writeText(txt)
  copied.value = key
  setTimeout(() => { copied.value = null }, 1200)
}
</script>

<template>
  <section id="quickstart" class="qs">
    <div class="qs__inner">
      <div>
        <div class="qs__label">Quickstart</div>
        <h2 class="qs__h">From zero to <em class="serif">running</em> in under a minute.</h2>
        <p class="qs__sub">
          Start with a plain <code class="mono" style="color: var(--accent); font-size: 0.95em">go mod</code>
          and pull in only the packages you need. No generators, no required CLI.
        </p>

        <div class="qs__steps">
          <div class="qs__step">
            <div class="qs__num">1</div>
            <div>
              <div class="qs__step-t">Start a new Go module</div>
              <p class="qs__step-b">No scaffolding required — bedrock integrates into any vanilla Go project.</p>
            </div>
          </div>
          <div class="qs__step">
            <div class="qs__num">2</div>
            <div>
              <div class="qs__step-t">Add just the packages you need</div>
              <p class="qs__step-b">Every package is a separate Go module. Add queue, cache or auth as your app grows.</p>
            </div>
          </div>
          <div class="qs__step">
            <div class="qs__num">3</div>
            <div>
              <div class="qs__step-t">Run it.</div>
              <p class="qs__step-b">Plain <code class="mono">go run</code> in dev, single static binary in prod. No runtime, no magic.</p>
            </div>
          </div>
        </div>
      </div>

      <div class="qs__term">
        <div class="qs__term-head">zsh · ~/my-app</div>

        <div class="qs__line" @click="copy('mkdir my-app && cd my-app && go mod init my-app', 'a')">
          <span class="qs__prompt">$</span>
          <span class="qs__cmd">mkdir my-app &amp;&amp; cd my-app &amp;&amp; go mod init my-app</span>
          <span class="qs__cp">
            <template v-if="copied === 'a'"><IconCheck :size="11" /> copied</template>
            <template v-else><IconCopy :size="11" /> copy</template>
          </span>
        </div>
        <div class="qs__line" @click="copy('go get github.com/bedrock/packages/container github.com/bedrock/packages/routing', 'b')">
          <span class="qs__prompt">$</span>
          <span class="qs__cmd">go get github.com/bedrock/packages/container github.com/bedrock/packages/routing</span>
          <span class="qs__cp">
            <template v-if="copied === 'b'"><IconCheck :size="11" /> copied</template>
            <template v-else><IconCopy :size="11" /> copy</template>
          </span>
        </div>
        <div class="qs__out">→ go: downloading github.com/bedrock/packages/container</div>
        <div class="qs__out">→ go: downloading github.com/bedrock/packages/routing</div>

        <div class="qs__line" @click="copy('go get github.com/bedrock/packages/queue github.com/bedrock/packages/cache', 'c')">
          <span class="qs__prompt">$</span>
          <span class="qs__cmd">go get github.com/bedrock/packages/queue github.com/bedrock/packages/cache</span>
          <span class="qs__cp">
            <template v-if="copied === 'c'"><IconCheck :size="11" /> copied</template>
            <template v-else><IconCopy :size="11" /> copy</template>
          </span>
        </div>
        <div class="qs__out qs__ok">✓ added 4 modules to go.mod</div>

        <div class="qs__line" @click="copy('go run ./cmd/api', 'd')">
          <span class="qs__prompt">$</span>
          <span class="qs__cmd">go run ./cmd/api</span>
          <span class="qs__cp">
            <template v-if="copied === 'd'"><IconCheck :size="11" /> copied</template>
            <template v-else><IconCopy :size="11" /> copy</template>
          </span>
        </div>
        <div class="qs__out qs__ok">✓ listening on :8080 — ready in 12ms</div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.qs {
  padding: 72px 32px 64px;
  max-width: var(--max); margin: 0 auto;
}
.qs__inner { display: grid; grid-template-columns: 1fr 1fr; gap: 56px; align-items: start; }
.qs__label {
  display: inline-flex; align-items: center; gap: 10px;
  font-family: 'JetBrains Mono', monospace; font-size: 11px;
  color: var(--ink-3); text-transform: uppercase; letter-spacing: 0.2em;
  margin-bottom: 18px;
}
.qs__label::before { content: ''; width: 24px; height: 1px; background: var(--line-2); }
.qs__h {
  font-size: clamp(28px, 3.2vw, 40px); line-height: 1.05; letter-spacing: -0.025em;
  margin: 0 0 16px; font-weight: 600;
}
.qs__h em { font-family: 'Instrument Serif', serif; font-style: italic; font-weight: 400; color: var(--accent); }
.qs__sub { color: var(--ink-2); font-size: 15px; line-height: 1.6; margin: 0 0 28px; max-width: 440px; }
.qs__steps { display: flex; flex-direction: column; gap: 18px; }
.qs__step { display: grid; grid-template-columns: 32px 1fr; gap: 14px; align-items: start; }
.qs__num {
  width: 28px; height: 28px; border-radius: 50%;
  display: grid; place-items: center;
  font-family: 'JetBrains Mono', monospace; font-size: 12px;
  background: color-mix(in oklab, var(--accent) 12%, transparent);
  color: var(--accent);
  border: 1px solid color-mix(in oklab, var(--accent) 30%, transparent);
  margin-top: 2px;
}
.qs__step-t { font-weight: 600; font-size: 15px; margin: 2px 0 4px; }
.qs__step-b { font-size: 13.5px; color: var(--ink-2); margin: 0; line-height: 1.55; }

.qs__term {
  background: #05080d;
  border: 1px solid var(--line);
  border-radius: 12px;
  padding: 18px 20px;
  font-family: 'JetBrains Mono', monospace;
  font-size: 13px;
  color: var(--ink);
  position: sticky; top: 90px;
}
.qs__term-head {
  display: flex; align-items: center; gap: 8px;
  font-size: 11px; color: var(--ink-3); text-transform: uppercase; letter-spacing: 0.12em;
  padding-bottom: 10px; margin-bottom: 10px;
  border-bottom: 1px solid var(--line);
}
.qs__term-head::before {
  content: ''; width: 8px; height: 8px; border-radius: 50%; background: #7ee07e;
  box-shadow: 0 0 12px #7ee07e;
}
.qs__line {
  display: flex; align-items: center; gap: 10px;
  padding: 6px 8px; margin: 0 -8px;
  border-radius: 6px;
  cursor: pointer;
  transition: background .15s;
  line-height: 1.7;
}
.qs__line:hover { background: color-mix(in oklab, var(--accent) 8%, transparent); }
.qs__prompt { color: var(--accent); user-select: none; }
.qs__cmd { flex: 1; overflow-x: auto; white-space: nowrap; }
.qs__cp {
  opacity: 0; transition: opacity .15s;
  font-size: 10px; color: var(--ink-3);
  display: inline-flex; align-items: center; gap: 4px;
  padding: 3px 8px; border: 1px solid var(--line); border-radius: 5px;
}
.qs__line:hover .qs__cp { opacity: 1; }
.qs__out {
  color: var(--ink-3); font-size: 12px;
  padding: 4px 0; line-height: 1.7;
}
.qs__ok { color: #7ee07e; }

@media (max-width: 900px) {
  .qs__inner { grid-template-columns: 1fr; gap: 36px; }
  .qs__term { position: static; }
}
</style>
