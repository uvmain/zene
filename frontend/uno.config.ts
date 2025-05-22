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
        // https://colorhunt.co/palette/22092c872341be3144f05941
        200: '#F05941',
        400: '#BE3144',
        600: '#872341',
        800: '#22092C',
        // 200: '#EEEEEE',
        // 400: '#8c8c8c',
        // 600: '#F79B72',
        // 800: '#2A4759',
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
