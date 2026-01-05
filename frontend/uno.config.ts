import { defineConfig, presetWebFonts, presetWind3, transformerDirectives, transformerVariantGroup } from 'unocss'

export default defineConfig({
  shortcuts: {
    'touch-target': 'min-h-11 min-w-11 flex items-center justify-center',
    'touch-button': 'touch-target cursor-pointer transition-colors duration-200',
    'mobile-grid': 'grid grid-cols-1 lg:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4',
    'mobile-padding': 'p-3 lg:p-6',
    'mobile-margin': 'm-3 lg:m-6',
    'mobile-gap': 'gap-3 lg:gap-6',
    'corner-cut': 'rounded-tl-12px [corner-top-left-shape:_bevel]',
    'corner-cut-large': 'rounded-tl-36px [corner-top-left-shape:_bevel]',
    'text-primary': 'dark:text-zshade-100 text-zshade-900',
    'text-muted': 'dark:text-zshade-200 text-zshade-800',
    'border-muted': 'dark:border-zshade-700 border-zshade-200 border-solid border-1',
    'background-1': 'dark:bg-zshade-900 bg-zshade-300',
    'background-2': 'dark:bg-zshade-800 bg-zshade-200',
    'background-grad-2': 'dark:from-zshade-800 from-zshade-200 bg-gradient-to-r',
    'background-3': 'dark:bg-zshade-700 bg-zshade-100',
    'z-button': 'border-1 border-solid relative font-semibold corner-cut flex cursor-pointer items-center justify-center cursor-pointer px-3 py-1 text-sm outline-none transition-all duration-200 text-muted background-2 hover:from-primary2 hover:text-primary hover:bg-gradient-to-br',
    'footer-icon': 'scale-100 text-lg text-muted transition-all duration-100 hover:scale-130 group-hover/button:scale-130 lg:text-xl sm:text-lg dark:hover:text-primary1 hover:text-primary1 dark:group-hover/button:text-primary1 group-hover/button:text-primary1',
    'auto-grid-6': '[grid-template-columns:repeat(auto-fit,minmax(min(150px,100%),1fr))] grid gap-6',
  },
  theme: {
    colors: {
      primary1: 'hsl(32,100%,50%)',
      primary2: 'hsl(32,75%,65%)',
      secondary1: 'hsl(332.13, 56.95%, 56.27%) ',
      secondary2: 'hsl(332.13, 56.95%, 66.27%) ',
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
      zgreen: 'hsl(105.63, 100%, 27.84%) ',
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
