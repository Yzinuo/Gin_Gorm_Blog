import {
  defineConfig,
  presetIcons,
  presetTypography,
  presetUno,
  transformerDirectives,
  transformerVariantGroup,
} from 'unocss'

export default defineConfig({
  theme: {
    colors: {
      'surface': 'var(--bg-surface)',
      'raised': 'var(--bg-raised)',
      'foreground': 'var(--text-primary)',
      'muted': 'var(--text-muted)',
      'brand': { DEFAULT: 'var(--brand)', hover: 'var(--brand-hover)', soft: 'var(--brand-soft)' },
      'on-brand': 'var(--on-brand)',
      'line': 'var(--border-color)',
    },
  },
  shortcuts: [
    ['f-c-c', 'flex justify-center items-center'],
  ],
  presets: [
    presetUno(),
    presetIcons({ warn: true }),
    presetTypography(),
  ],
  transformers: [
    transformerDirectives(),
    transformerVariantGroup(),
  ],
})
