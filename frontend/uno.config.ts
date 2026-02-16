import { defineConfig, presetWind3, transformerDirectives, transformerVariantGroup } from 'unocss'

export default defineConfig({
  shortcuts: {
    'media-control-button': 'h-10 w-10 flex cursor-pointer items-center justify-center border-none bg-white/0 font-semibold outline-none lg:h-12 lg:w-12 sm:h-10 sm:w-10',
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
    'footer-icon': 'scale-100 text-lg text-muted transition-all duration-100 hover:scale-130 group-hover/button:scale-130 lg:text-xl sm:text-lg active:text-primary1',
    'footer-icon-disabled': 'scale-100 text-lg text-muted lg:text-xl sm:text-lg opacity-50',
    'footer-icon-on': 'scale-110 text-lg text-primary1 transition-all duration-100 hover:scale-130 group-hover/button:scale-130 lg:text-xl sm:text-lg',
    'auto-grid-6': '[grid-template-columns:repeat(auto-fit,minmax(min(150px,100%),1fr))] grid gap-6',
  },
  theme: {
    colors: {
      primary1: 'hsl(32 100% 50%)',
      primary2: 'hsl(32 75% 65%)',
      secondary1: 'hsl(332.13 56.95% 56.27%)',
      secondary2: 'hsl(332.13 56.95% 66.27%)',
      zshade: {
        100: 'hsl(37 7% 90%)',
        200: 'hsl(37 7% 80%)',
        300: 'hsl(37 7% 70%)',
        500: 'hsl(37 7% 50%)',
        600: 'hsl(37 7% 40%)',
        700: 'hsl(37 7% 30%)',
        800: 'hsl(37 7% 20%)',
        900: 'hsl(37 7% 10%)',
      },
      zgreen: 'hsl(105.63 100% 27.84%) ',
    },
  },
  presets: [
    presetWind3(),
  ],
  transformers: [
    transformerDirectives(),
    transformerVariantGroup(),
  ],
})
