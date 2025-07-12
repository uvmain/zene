import {
  defineConfig,
  presetUno,
  presetWind,
  transformerDirectives,
  transformerVariantGroup,
} from 'unocss'

function getSafelist(): string[] {
  const base = 'prose prose-sm m-auto text-left'.split(' ')
  const unusedSafelist: string[] = []
  return [...unusedSafelist, ...base]
}

export default defineConfig({
  shortcuts: {
    // Mobile-friendly touch targets
    'touch-target': 'min-h-11 min-w-11 flex items-center justify-center',
    'touch-button': 'touch-target cursor-pointer transition-colors duration-200',
    // Mobile-first responsive grid
    'mobile-grid': 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4',
    // Standard mobile spacing
    'mobile-padding': 'p-3 md:p-6',
    'mobile-margin': 'm-3 md:m-6',
    'mobile-gap': 'gap-3 md:gap-6',
  },
  theme: {
    colors: {
      zene: {
        // https://colorhunt.co/palette/22092c872341be3144f05941
        200: '#F05941',
        400: '#BE3144',
        600: '#872341',
        700: '#544F57',
        800: '#22092C',
        // 200: '#BED754',
        // 400: '#E3651D',
        // 600: '#750E21',
        // 800: '#191919',
      },
    },
  },
  presets: [
    presetUno(),
    presetWind(),
  ],
  transformers: [
    transformerDirectives(),
    transformerVariantGroup(),
  ],
  safelist: getSafelist(),
})
