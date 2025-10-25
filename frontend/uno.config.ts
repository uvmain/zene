import presetWebFonts from '@unocss/preset-web-fonts'
import { presetWind3 } from '@unocss/preset-wind3'

import {
  defineConfig,
  transformerDirectives,
  transformerVariantGroup,
} from 'unocss'

export default defineConfig({
  shortcuts: {
    'touch-target': 'min-h-11 min-w-11 flex items-center justify-center',
    'touch-button': 'touch-target cursor-pointer transition-colors duration-200',
    'mobile-grid': 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4',
    'mobile-padding': 'p-3 md:p-6',
    'mobile-margin': 'm-3 md:m-6',
    'mobile-gap': 'gap-3 md:gap-6',
    'corner-cut': '[clip-path:polygon(10px_0,100%_0,100%_100%,0_100%,0_15px)]',
    'corner-cut-large': '[clip-path:polygon(30px_0,100%_0,100%_100%,0_100%,0_45px)]',
    'text-primary': 'dark:text-zshade-100 text-zshade-900',
    'text-muted': 'dark:text-zshade-200 text-zshade-800',
    'background-1': 'dark:bg-zshade-900 bg-zshade-300',
    'background-2': 'dark:bg-zshade-800 bg-zshade-200',
    'background-grad-2': 'dark:from-zshade-800 from-zshade-200 bg-gradient-to-r',
    'background-3': 'dark:bg-zshade-700 bg-zshade-100',
    'z-button': 'border-none relative font-semibold corner-cut flex cursor-pointer items-center justify-center cursor-pointer px-3 py-1 text-sm outline-none transition-all duration-200 text-muted background-2 hover:from-primary2 hover:text-primary hover:bg-gradient-to-br',
    'footer-icon': 'scale-100 text-lg text-muted transition-all duration-100 hover:scale-130 group-hover/button:scale-130 md:text-xl sm:text-lg dark:hover:text-primary1 hover:text-primary1 dark:group-hover/button:text-primary1 group-hover/button:text-primary1',
  },
  theme: {
    colors: {
      primary1: 'hsl(32,100%,50%)',
      primary2: 'hsl(32,75%,65%)',
      zshade: {
        100: 'hsl(37,10%,90%)',
        200: 'hsl(37,10%,80%)',
        300: 'hsl(37,10%,70%)',
        400: 'hsl(37,10%,60%)',
        500: 'hsl(37,10%,50%)',
        600: 'hsl(37,10%,40%)',
        700: 'hsl(37,10%,30%)',
        800: 'hsl(37,10%,20%)',
        900: 'hsl(37,10%,10%)',
      },
    },
  },
  presets: [
    presetWind3(),
    presetWebFonts({
      provider: 'google',
      fonts: {
        jura: ['Jura:300,400,500,600,700'],
      },
    }),
  ],
  transformers: [
    transformerDirectives(),
    transformerVariantGroup(),
  ],
})
