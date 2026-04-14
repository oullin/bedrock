import { defineUserConfig } from 'vuepress'
import { viteBundler } from '@vuepress/bundler-vite'
import tailwindcss from '@tailwindcss/vite'
import { bedrockTheme } from './theme/index'

export default defineUserConfig({
  lang: 'en-US',
  title: 'Bedrock',
  description: 'Foundational Go packages for building web applications.',

  // Prevent dark mode flash: apply saved preference before first paint
  head: [
    [
      'script',
      {},
      `(function(){var s=localStorage.getItem('bedrock-color-scheme');if(s==='dark'||(!s&&window.matchMedia('(prefers-color-scheme:dark)').matches)){document.documentElement.classList.add('dark');}})();`,
    ],
  ],

  bundler: viteBundler({
    viteOptions: {
      plugins: [tailwindcss()],
      server: {
        port: parseInt(process.env.PORT ?? '8080'),
      },
    },
  }),

  theme: bedrockTheme({
    repo: 'https://github.com/gocanto/bedrock',

    navbar: [
      { text: 'Home', link: '/' },
      { text: 'Getting Started', link: '/getting-started' },
      { text: 'Packages', link: '/packages/auth' },
    ],

    sidebar: {
      '/': [
        {
          text: 'Introduction',
          children: [
            { text: 'Home', link: '/' },
            { text: 'Getting Started', link: '/getting-started' },
          ],
        },
        {
          text: 'Packages',
          collapsible: true,
          children: [
            { text: 'auth', link: '/packages/auth' },
            { text: 'bus', link: '/packages/bus' },
            { text: 'cache', link: '/packages/cache' },
            { text: 'cookie', link: '/packages/cookie' },
            { text: 'fortify', link: '/packages/fortify' },
            { text: 'httpx', link: '/packages/httpx' },
            { text: 'jetstream', link: '/packages/jetstream' },
            { text: 'queue', link: '/packages/queue' },
            { text: 'routing', link: '/packages/routing' },
            { text: 'session', link: '/packages/session' },
            { text: 'spark', link: '/packages/spark' },
          ],
        },
      ],
    },
  }),
})
