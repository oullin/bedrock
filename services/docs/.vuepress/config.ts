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
    ['link', { rel: 'stylesheet', href: 'https://fonts.googleapis.com/css2?family=Geist:wght@300;400;500;600;700;800&display=swap' }],
    ['link', { rel: 'stylesheet', href: 'https://fonts.googleapis.com/css2?family=Instrument+Serif:ital@0;1&display=swap' }],
    ['link', { rel: 'stylesheet', href: 'https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;600&display=swap' }],
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
            { text: 'Directory Structure', link: '/getting-started/directory-structure' },
          ],
        },
        {
          text: 'Architecture Concepts',
          collapsible: true,
          children: [
            { text: 'Request Lifecycle', link: '/architecture/lifecycle' },
            { text: 'Application Bootstrap', link: '/architecture/application' },
            { text: 'Service Container', link: '/architecture/service-container' },
            { text: 'Service Providers', link: '/architecture/service-providers' },
            { text: 'Facades', link: '/architecture/facades' },
            { text: 'Configuration', link: '/architecture/configuration' },
            { text: 'Drivers', link: '/architecture/drivers' },
          ],
        },
        {
          text: 'Concepts',
          collapsible: true,
          children: [
            { text: 'Testing', link: '/concepts/testing' },
          ],
        },
        {
          text: 'Architecture',
          collapsible: true,
          children: [
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
            { text: 'precognition', link: '/packages/precognition' },
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
            { text: 'scout', link: '/packages/scout' },
          ],
        },
        {
          text: 'Events & Jobs',
          collapsible: true,
          children: [
            { text: 'events', link: '/packages/events' },
            { text: 'bus', link: '/packages/bus' },
            { text: 'queue', link: '/packages/queue' },
            { text: 'jobqueue', link: '/packages/jobqueue' },
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
            { text: 'broadcasting', link: '/packages/broadcasting' },
            { text: 'echo', link: '/packages/echo' },
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
            { text: 'AI SDK', link: '/ai/sdk' },
            { text: 'MCP', link: '/ai/mcp' },
            { text: 'Boost', link: '/ai/boost' },
          ],
        },
        {
          text: 'Support & Utilities',
          collapsible: true,
          children: [
            { text: 'collection', link: '/packages/collection' },
            { text: 'support', link: '/packages/support' },
            { text: 'log', link: '/packages/log' },
            { text: 'translation', link: '/packages/translation' },
            { text: 'concurrency', link: '/packages/concurrency' },
            { text: 'conditionable', link: '/packages/conditionable' },
            { text: 'jsonx', link: '/packages/jsonx' },
            { text: 'lottery', link: '/packages/lottery' },
            { text: 'money', link: '/packages/money' },
            { text: 'process', link: '/packages/process' },
            { text: 'seo', link: '/packages/seo' },
            { text: 'str', link: '/packages/str' },
            { text: 'wayfinder', link: '/packages/wayfinder' },
          ],
        },
        {
          text: 'Developer Tools',
          collapsible: true,
          children: [
            { text: 'prompts', link: '/packages/prompts' },
            { text: 'pail', link: '/packages/pail' },
            { text: 'envoy', link: '/packages/envoy' },
            { text: 'debugbar', link: '/packages/debugbar' },
          ],
        },
        {
          text: 'Products',
          collapsible: true,
          children: [
            { text: 'inception', link: '/packages/inception' },
            { text: 'authflows', link: '/packages/authflows' },
            { text: 'inertia', link: '/packages/inertia' },
            { text: 'jetstream', link: '/packages/jetstream' },
            { text: 'spark', link: '/packages/spark' },
          ],
        },
      ],
    },
  }),
})
