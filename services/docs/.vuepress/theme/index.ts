import { getDirname, path } from 'vuepress/utils'
import { prismjsPlugin } from '@vuepress/plugin-prismjs'
import type { Theme } from 'vuepress/core'

const __dirname = getDirname(import.meta.url)

export interface NavbarItem {
  text: string
  link?: string
  children?: NavbarItem[]
}

export interface SidebarItem {
  text: string
  link?: string
  collapsible?: boolean
  children?: SidebarItem[]
}

export type SidebarConfig = Record<string, SidebarItem[]>

export interface BedrockThemeOptions {
  repo?: string
  navbar?: NavbarItem[]
  sidebar?: SidebarConfig
}

export function bedrockTheme(options: BedrockThemeOptions = {}): Theme {
  return {
    name: 'vuepress-theme-bedrock',

    clientConfigFile: path.resolve(__dirname, './client.ts'),

    // Inject theme options as a build-time constant accessible in all components
    define: {
      __BEDROCK_THEME_OPTIONS__: JSON.stringify(options),
    },

    plugins: [
      prismjsPlugin({ theme: 'one-dark-pro' }),
    ],
  }
}
