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
    ['link', { rel: 'preconnect', href: 'https://fonts.googleapis.com' }],
    ['link', { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' }],
    ['link', { rel: 'stylesheet', href: 'https://fonts.googleapis.com/css2?family=Lexend:wght@400;500;600;700&display=swap' }],
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
          text: 'Concepts',
          collapsible: true,
          children: [
            { text: 'Request Lifecycle', link: '/concepts/request-lifecycle' },
            { text: 'Testing', link: '/concepts/testing' },
          ],
        },
        {
          text: 'Architecture',
          collapsible: true,
          children: [
            { text: 'bedrock', link: '/packages/bedrock' },
            { text: 'bootstrap', link: '/packages/bootstrap' },
            { text: 'container', link: '/packages/container' },
            { text: 'config', link: '/packages/config' },
            { text: 'contracts', link: '/packages/contracts' },
            { text: 'facades', link: '/packages/facades' },
          ],
        },
        {
          text: 'The Basics',
          collapsible: true,
          children: [
            { text: 'routing', link: '/packages/routing' },
            { text: 'controllers', link: '/basics/controllers' },
            { text: 'middleware', link: '/basics/middleware' },
            { text: 'url generation', link: '/basics/url-generation' },
            { text: 'httpx', link: '/packages/httpx' },
            { text: 'inertia', link: '/packages/inertia' },
            { text: 'session', link: '/packages/session' },
            { text: 'cookie', link: '/packages/cookie' },
            { text: 'csrf protection', link: '/basics/csrf' },
            { text: 'validation', link: '/packages/validation' },
            { text: 'httppreview', link: '/packages/httppreview' },
          ],
        },
        {
          text: 'Security',
          collapsible: true,
          children: [
            { text: 'auth', link: '/packages/auth' },
            { text: 'encryption', link: '/packages/encryption' },
            { text: 'hashing', link: '/packages/hashing' },
            { text: 'oauthserver', link: '/packages/oauthserver' },
            { text: 'socialauth', link: '/packages/socialauth' },
          ],
        },
        {
          text: 'Database',
          collapsible: true,
          children: [
            { text: 'database', link: '/packages/database' },
          ],
        },
        {
          text: 'Data & Storage',
          collapsible: true,
          children: [
            { text: 'cache', link: '/packages/cache' },
            { text: 'redis', link: '/packages/redis' },
            { text: 'filesystem', link: '/packages/filesystem' },
            { text: 'pagination', link: '/packages/pagination' },
            { text: 'search', link: '/packages/search' },
          ],
        },
        {
          text: 'Events & Jobs',
          collapsible: true,
          children: [
            { text: 'events', link: '/packages/events' },
            { text: 'bus', link: '/packages/bus' },
            { text: 'queue', link: '/packages/queue' },
            { text: 'pipeline', link: '/packages/pipeline' },
          ],
        },
        {
          text: 'Communication',
          collapsible: true,
          children: [
            { text: 'mailx', link: '/packages/mailx' },
            { text: 'notifications', link: '/packages/notifications' },
          ],
        },
        {
          text: 'Real-time',
          collapsible: true,
          children: [
            { text: 'broadcastclient', link: '/packages/broadcastclient' },
            { text: 'websockets', link: '/packages/websockets' },
          ],
        },
        {
          text: 'Feature Flags',
          collapsible: true,
          children: [
            { text: 'featureflags', link: '/packages/featureflags' },
          ],
        },
        {
          text: 'AI & Integrations',
          collapsible: true,
          children: [
            { text: 'ai', link: '/packages/ai' },
            { text: 'mcp', link: '/packages/mcp' },
            { text: 'boost', link: '/packages/boost' },
          ],
        },
        {
          text: 'Support & Utilities',
          collapsible: true,
          children: [
            { text: 'support', link: '/packages/support' },
            { text: 'log', link: '/packages/log' },
            { text: 'translation', link: '/packages/translation' },
            { text: 'concurrency', link: '/packages/concurrency' },
            { text: 'conditionable', link: '/packages/conditionable' },
            { text: 'helpers', link: '/packages/helpers' },
            { text: 'jsonx', link: '/packages/jsonx' },
            { text: 'seo', link: '/packages/seo' },
            { text: 'routegen', link: '/packages/routegen' },
          ],
        },
        {
          text: 'Developer Tools',
          collapsible: true,
          children: [
            { text: 'prompts', link: '/packages/prompts' },
            { text: 'debugbar', link: '/packages/debugbar' },
          ],
        },
        {
          text: 'Products',
          collapsible: true,
          children: [
            { text: 'inception', link: '/packages/inception' },
            { text: 'billing', link: '/packages/billing' },
          ],
        },
      ],
    },
  }),
})
