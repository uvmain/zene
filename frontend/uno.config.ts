import presetWebFonts from '@unocss/preset-web-fonts'
import { presetWind3 } from '@unocss/preset-wind3'

import {
  defineConfig,
  transformerDirectives,
  transformerVariantGroup,
} from 'unocss'

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
    'corner-cut': '[clip-path:polygon(10px_0,100%_0,100%_100%,0_100%,0_15px)]',
    'corner-cut-large': '[clip-path:polygon(30px_0,100%_0,100%_100%,0_100%,0_45px)]',
    'text-primary': 'text-zgray-200 dark:text-zgray-800',
    'background-primary': 'bg-zgray-800 dark:bg-zgray-200',
    'z-button': 'corner-cut flex cursor-pointer items-center justify-center cursor-pointer border-none bg-accent2 px-3 py-1 text-sm text-zgray-200 outline-none transition-all duration-200 hover:bg-accent1 hover:bg-accent1 hover:from-accent2 hover:bg-gradient-to-br',
  },
  theme: {
    colors: {
      accent1: '#d60072',
      accent2: '#eb7d00',
      zgray: {
        200: '#f2f2e3',
        300: '#f2f2e3',
        400: '#99998bff',
        600: '#52524cff',
        800: '#292926ff',
      },
      zene: {
        // https://colorhunt.co/palette/22092c872341be3144f05941
        200: '#F05941',
        400: '#BE3144',
        600: '#872341',
        700: '#544F57',
        800: '#22092C',
      },
    },
  },
  presets: [
    presetWind3(),
    presetWebFonts({
      provider: 'google',
      fonts: {
        jura: ['Jura:300,400,500,600,700'],
        montserrat: ['Montserrat:100,200,300,400,500,600,700,800,900'],
      },
    }),
  ],
  transformers: [
    transformerDirectives(),
    transformerVariantGroup(),
  ],
})
