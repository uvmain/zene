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
  shortcuts: {},
  theme: {
    colors: {
      zene: {
        200: '#F05941',
        400: '#BE3144',
        600: '#872341',
        800: '#22092C',
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
