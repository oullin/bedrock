<script setup lang="ts">
import { RouteLink } from 'vuepress/client'
import BrandMark from './BrandMark.vue'
import { IconArrow, IconGh, IconDiscord } from './icons'
</script>

<template>
  <footer class="ft">
    <div class="ft__cta">
      <h2 class="ft__cta-h">Prepare for a dev experience that can finally <em class="serif">keep pace</em> with the speed of your mind.</h2>
      <p class="ft__cta-sub">Install one package. Or all 21. Ship in an afternoon.</p>
      <div class="ft__cta-actions">
        <RouteLink to="/getting-started" class="btn btn--primary">Get started <IconArrow :size="15" /></RouteLink>
        <a class="btn btn--ghost" href="https://github.com/gocanto/bedrock" target="_blank" rel="noopener noreferrer"><IconGh :size="15" /> Star on GitHub</a>
      </div>
    </div>

    <div class="ft__bottom">
      <div>
        <div class="ft__brand-h">
          <div class="ft__mark"><BrandMark :size="24" /></div> bedrock
        </div>
        <p class="ft__tagline">A collection of well-tested, Upstream-inspired libraries for building production-grade Go web applications.</p>
        <div class="ft__social">
          <a href="https://github.com/gocanto/bedrock" target="_blank" rel="noopener noreferrer" aria-label="GitHub"><IconGh :size="14" /></a>
          <a href="#" aria-label="Discord"><IconDiscord :size="14" /></a>
        </div>
      </div>
      <div class="ft__col">
        <h4>Docs</h4>
        <ul>
          <li><RouteLink to="/getting-started">Get started</RouteLink></li>
          <li><RouteLink to="/concepts/request-lifecycle">Guides</RouteLink></li>
          <li><RouteLink to="/packages/container">Reference</RouteLink></li>
          <li><RouteLink to="/getting-started">Upgrading</RouteLink></li>
        </ul>
      </div>
      <div class="ft__col">
        <h4>Packages</h4>
        <ul>
          <li><RouteLink to="/packages/routing">HTTP</RouteLink></li>
          <li><RouteLink to="/packages/queue">Data &amp; messaging</RouteLink></li>
          <li><RouteLink to="/packages/encryption">Infrastructure</RouteLink></li>
          <li><a href="#packages">All 21</a></li>
        </ul>
      </div>
      <div class="ft__col">
        <h4>Community</h4>
        <ul>
          <li><a href="#">Discord</a></li>
          <li><a href="https://github.com/gocanto/bedrock/discussions" target="_blank" rel="noopener noreferrer">GitHub discussions</a></li>
          <li><a href="#">Showcase</a></li>
          <li><a href="#">Roadmap</a></li>
        </ul>
      </div>
      <div class="ft__col">
        <h4>Project</h4>
        <ul>
          <li><a href="#">Changelog</a></li>
          <li><a href="https://github.com/gocanto/bedrock/blob/main/CONTRIBUTING.md" target="_blank" rel="noopener noreferrer">Contributing</a></li>
          <li><a href="#">Security</a></li>
          <li><a href="https://github.com/gocanto/bedrock/blob/main/LICENSE" target="_blank" rel="noopener noreferrer">License · MIT</a></li>
        </ul>
      </div>
    </div>

    <div class="ft__meta">
      <div class="ft__meta-l">
        <span>© 2026 bedrock</span>
        <span>·</span>
        <span>v1.4.0</span>
      </div>
      <div>Built with care. No telemetry, ever.</div>
    </div>
  </footer>
</template>

<style scoped>
.ft {
  margin-top: 40px;
  border-top: 1px solid var(--line);
  background: linear-gradient(180deg, transparent, color-mix(in oklab, var(--accent) 5%, transparent));
  position: relative;
  overflow: hidden;
}
.ft__cta {
  padding: 96px 32px 96px;
  text-align: center;
  position: relative;
}
.ft__cta::before {
  content: ''; position: absolute; inset: 0;
  background: radial-gradient(50% 80% at 50% 100%, color-mix(in oklab, var(--accent) 22%, transparent) 0%, transparent 70%);
  pointer-events: none;
}
.ft__cta-h {
  font-size: clamp(32px, 4vw, 54px); line-height: 1.05; letter-spacing: -0.03em;
  margin: 0 auto 10px; font-weight: 600; max-width: 780px; position: relative;
}
.ft__cta-h em { font-family: 'Instrument Serif', serif; font-style: italic; font-weight: 400; color: var(--accent); }
.ft__cta-sub {
  font-size: 16px; color: var(--ink-2); margin: 0 auto 32px; max-width: 520px; position: relative;
}
.ft__cta-actions { display: inline-flex; gap: 10px; position: relative; }

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

.ft__bottom {
  display: grid; grid-template-columns: 2fr 1fr 1fr 1fr 1fr; gap: 40px;
  padding: 48px 32px;
  max-width: var(--max); margin: 0 auto;
  border-top: 1px solid var(--line);
}
.ft__brand-h { display: flex; align-items: center; gap: 10px; margin-bottom: 14px; font-weight: 700; font-size: 15px; }
.ft__mark {
  width: 24px; height: 24px;
  display: grid; place-items: center;
}
.ft__tagline { color: var(--ink-3); font-size: 13px; max-width: 280px; line-height: 1.55; margin: 0 0 16px; }
.ft__social { display: flex; gap: 8px; }
.ft__social a {
  width: 32px; height: 32px; border-radius: 8px;
  display: grid; place-items: center;
  border: 1px solid var(--line); color: var(--ink-2);
  transition: color .15s, border-color .15s, background .15s;
}
.ft__social a:hover { color: var(--ink); border-color: var(--line-2); background: var(--panel); }
.ft__col h4 {
  font-size: 11px; color: var(--ink-3); text-transform: uppercase;
  letter-spacing: 0.14em; margin: 0 0 12px; font-weight: 600;
}
.ft__col ul { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 8px; }
.ft__col a { color: var(--ink-2); font-size: 13.5px; }
.ft__col a:hover { color: var(--ink); }

.ft__meta {
  display: flex; align-items: center; justify-content: space-between;
  padding: 20px 32px; color: var(--ink-3); font-size: 12px;
  border-top: 1px solid var(--line);
}
.ft__meta-l { display: flex; gap: 16px; font-family: 'JetBrains Mono', monospace; }
@media (max-width: 860px) {
  .ft__bottom { grid-template-columns: 1fr 1fr; }
  .ft__meta { flex-direction: column; gap: 12px; }
}
</style>
