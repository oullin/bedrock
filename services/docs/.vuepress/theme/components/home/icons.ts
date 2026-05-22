import { h, defineComponent, type Component } from 'vue'
import Icon from './Icon.vue'

interface IconProps {
  size?: number
  stroke?: number
}

function makeIcon(paths: string[], opts: { fill?: string; stroke?: number } = {}): Component {
  return defineComponent({
    name: 'BedrockIcon',
    props: {
      size: { type: Number, default: 18 },
      stroke: { type: Number, default: 1.6 },
    },
    setup(props: IconProps) {
      return () => h(
        Icon,
        {
          size: props.size,
          stroke: opts.stroke ?? props.stroke,
          fill: opts.fill ?? 'none',
        },
        () => paths.map(d => h('path', { d })),
      )
    },
  })
}

function makeIconRaw(children: () => any, opts: { fill?: string; stroke?: number } = {}): Component {
  return defineComponent({
    name: 'BedrockIcon',
    props: {
      size: { type: Number, default: 18 },
      stroke: { type: Number, default: 1.6 },
    },
    setup(props: IconProps) {
      return () => h(
        Icon,
        {
          size: props.size,
          stroke: opts.stroke ?? props.stroke,
          fill: opts.fill ?? 'none',
        },
        children,
      )
    },
  })
}

export const IconContainer = makeIcon([
  'M4 7l8-4 8 4-8 4-8-4Z',
  'M4 12l8 4 8-4',
  'M4 17l8 4 8-4',
])

export const IconRouting = makeIconRaw(() => [
  h('circle', { cx: 5, cy: 6, r: 2 }),
  h('circle', { cx: 19, cy: 6, r: 2 }),
  h('circle', { cx: 12, cy: 18, r: 2 }),
  h('path', { d: 'M7 6h10' }),
  h('path', { d: 'M12 8v8' }),
  h('path', { d: 'M6 8l6 8' }),
  h('path', { d: 'M18 8l-6 8' }),
])

export const IconHttp = makeIconRaw(() => [
  h('path', { d: 'M3 12h18' }),
  h('path', { d: 'M12 3a15 15 0 0 1 0 18' }),
  h('path', { d: 'M12 3a15 15 0 0 0 0 18' }),
  h('circle', { cx: 12, cy: 12, r: 9 }),
])

export const IconValidation = makeIcon(['M20 7L10 17l-5-5', 'M4 20h16'])

export const IconAuth = makeIconRaw(() => [
  h('rect', { x: 5, y: 11, width: 14, height: 10, rx: 2 }),
  h('path', { d: 'M8 11V8a4 4 0 0 1 8 0v3' }),
])

export const IconCache = makeIconRaw(() => [
  h('ellipse', { cx: 12, cy: 6, rx: 8, ry: 3 }),
  h('path', { d: 'M4 6v6c0 1.7 3.6 3 8 3s8-1.3 8-3V6' }),
  h('path', { d: 'M4 12v6c0 1.7 3.6 3 8 3s8-1.3 8-3v-6' }),
])

export const IconRedis = makeIcon([
  'M3 7l9-4 9 4-9 4-9-4Z',
  'M3 12l9 4 9-4',
  'M3 17l9 4 9-4',
])

export const IconEvents = makeIconRaw(() => [
  h('circle', { cx: 12, cy: 12, r: 2 }),
  h('path', { d: 'M12 2v4' }),
  h('path', { d: 'M12 18v4' }),
  h('path', { d: 'M2 12h4' }),
  h('path', { d: 'M18 12h4' }),
  h('path', { d: 'M5 5l2.5 2.5' }),
  h('path', { d: 'M16.5 16.5L19 19' }),
  h('path', { d: 'M5 19l2.5-2.5' }),
  h('path', { d: 'M16.5 7.5L19 5' }),
])

export const IconQueue = makeIconRaw(() => [
  h('rect', { x: 3, y: 5, width: 4, height: 14, rx: 1 }),
  h('rect', { x: 10, y: 5, width: 4, height: 14, rx: 1 }),
  h('rect', { x: 17, y: 5, width: 4, height: 14, rx: 1 }),
])

export const IconBus = makeIconRaw(() => [
  h('circle', { cx: 6, cy: 12, r: 2 }),
  h('circle', { cx: 18, cy: 6, r: 2 }),
  h('circle', { cx: 18, cy: 18, r: 2 }),
  h('path', { d: 'M8 12h6' }),
  h('path', { d: 'M14 12l2-4' }),
  h('path', { d: 'M14 12l2 4' }),
])

export const IconMail = makeIconRaw(() => [
  h('rect', { x: 3, y: 5, width: 18, height: 14, rx: 2 }),
  h('path', { d: 'M3 7l9 6 9-6' }),
])

export const IconEncryption = makeIcon([
  'M12 3l8 3v6c0 5-3.5 8-8 9-4.5-1-8-4-8-9V6l8-3Z',
  'M9 12l2 2 4-4',
])

export const IconHash = makeIcon([
  'M5 9h14',
  'M5 15h14',
  'M10 3l-2 18',
  'M16 3l-2 18',
])

export const IconFile = makeIcon([
  'M6 3h8l4 4v14a1 1 0 0 1-1 1H6a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1Z',
  'M14 3v4h4',
])

export const IconPagination = makeIconRaw(() => [
  h('rect', { x: 3, y: 6, width: 18, height: 12, rx: 2 }),
  h('path', { d: 'M9 12h6' }),
  h('path', { d: 'M12 9l-3 3 3 3' }),
  h('path', { d: 'M12 9l3 3-3 3' }),
])

export const IconTranslation = makeIcon([
  'M3 5h10',
  'M8 3v2',
  'M4 17l5-10 5 10',
  'M5 14h8',
  'M14 21l4-8 4 8',
  'M15 18h6',
])

export const IconSupport = makeIconRaw(() => [
  h('circle', { cx: 12, cy: 12, r: 9 }),
  h('path', { d: 'M8 12h.01' }),
  h('path', { d: 'M12 12h.01' }),
  h('path', { d: 'M16 12h.01' }),
])

export const IconLog = makeIconRaw(() => [
  h('rect', { x: 4, y: 3, width: 16, height: 18, rx: 2 }),
  h('path', { d: 'M8 8h8' }),
  h('path', { d: 'M8 12h8' }),
  h('path', { d: 'M8 16h5' }),
])

export const IconInception = makeIconRaw(() => [
  h('path', { d: 'M12 3v18' }),
  h('path', { d: 'M3 12h18' }),
  h('circle', { cx: 12, cy: 12, r: 9 }),
  h('circle', { cx: 12, cy: 12, r: 4 }),
])

export const IconBilling = makeIcon(['M12 3l1.8 5.2L19 10l-5.2 1.8L12 17l-1.8-5.2L5 10l5.2-1.8L12 3Z'])

export const IconArrow = makeIcon(['M5 12h14', 'M13 6l6 6-6 6'])

export const IconGh = makeIconRaw(
  () => [
    h('path', {
      fill: 'currentColor',
      d: 'M12 2a10 10 0 0 0-3.16 19.49c.5.09.68-.22.68-.48v-1.7c-2.78.6-3.37-1.34-3.37-1.34-.45-1.16-1.11-1.46-1.11-1.46-.9-.62.07-.6.07-.6 1 .07 1.52 1.03 1.52 1.03.89 1.52 2.34 1.08 2.91.82.09-.64.35-1.08.63-1.33-2.22-.25-4.55-1.11-4.55-4.94 0-1.09.39-1.99 1.03-2.69-.1-.26-.45-1.28.1-2.67 0 0 .84-.27 2.75 1.03a9.48 9.48 0 0 1 5 0c1.9-1.3 2.74-1.03 2.74-1.03.55 1.39.2 2.41.1 2.67.64.7 1.03 1.6 1.03 2.69 0 3.84-2.34 4.69-4.57 4.93.36.31.68.92.68 1.85v2.74c0 .27.18.58.69.48A10 10 0 0 0 12 2Z',
    }),
  ],
  { stroke: 0 },
)

export const IconCopy = makeIconRaw(() => [
  h('rect', { x: 9, y: 9, width: 11, height: 11, rx: 2 }),
  h('path', { d: 'M5 15V5a1 1 0 0 1 1-1h10' }),
])

export const IconCheck = makeIcon(['M5 12l5 5 9-11'])

export const IconBolt = makeIcon(['M13 3L4 14h7l-1 7 9-11h-7l1-7Z'])

export const IconShield = makeIcon(['M12 3l8 3v6c0 5-3.5 8-8 9-4.5-1-8-4-8-9V6l8-3Z'])

export const IconStack = makeIcon([
  'M4 7l8-4 8 4-8 4-8-4Z',
  'M4 12l8 4 8-4',
  'M4 17l8 4 8-4',
])

export const IconLayers = makeIcon([
  'M12 3l9 5-9 5-9-5 9-5Z',
  'M3 13l9 5 9-5',
])

export const IconStar = makeIconRaw(
  () => [h('path', { d: 'M12 2l3 6.9 7.6.7-5.7 5 1.7 7.4L12 18.3 5.4 22l1.7-7.4-5.7-5 7.6-.7L12 2Z' })],
  { fill: 'currentColor', stroke: 0 },
)

export const IconDot = makeIconRaw(
  () => [h('circle', { cx: 12, cy: 12, r: 4 })],
  { fill: 'currentColor', stroke: 0 },
)

export const IconSearch = makeIconRaw(() => [
  h('circle', { cx: 11, cy: 11, r: 7 }),
  h('path', { d: 'M20 20l-3.5-3.5' }),
])

export const IconDiscord = makeIconRaw(
  () => [
    h('path', {
      fill: 'currentColor',
      d: 'M20 5.2A17 17 0 0 0 16 4l-.2.4c1.4.3 2.6.8 3.7 1.5A13 13 0 0 0 4.5 6c1-.7 2.3-1.2 3.7-1.5L8 4a17 17 0 0 0-4 1.2C1.3 9.3.6 13.2.9 17c1.7 1.3 3.3 2 4.9 2.4l.9-1.4a10 10 0 0 1-1.7-.9l.4-.3a11.4 11.4 0 0 0 11.2 0l.4.3c-.5.4-1.1.7-1.7.9l.9 1.4c1.6-.5 3.2-1.1 4.9-2.5.3-4.3-.6-8.2-3.1-11.7ZM9 15a1.7 1.7 0 0 1-1.6-1.8c0-1 .7-1.8 1.6-1.8s1.7.8 1.6 1.8c0 1-.7 1.8-1.6 1.8Zm6 0a1.7 1.7 0 0 1-1.6-1.8c0-1 .7-1.8 1.6-1.8s1.7.8 1.6 1.8c0 1-.7 1.8-1.6 1.8Z',
    }),
  ],
  { stroke: 0 },
)

export const ICON_MAP: Record<string, Component> = {
  IconContainer,
  IconRouting,
  IconHttp,
  IconValidation,
  IconAuth,
  IconCache,
  IconRedis,
  IconEvents,
  IconQueue,
  IconBus,
  IconMail,
  IconEncryption,
  IconHash,
  IconFile,
  IconPagination,
  IconTranslation,
  IconSupport,
  IconLog,
  IconInception,
  IconBilling,
}
