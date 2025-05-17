/// <reference types="vite-ssg" />

import { fileURLToPath, URL } from 'node:url'
import vue from '@vitejs/plugin-vue'
import UnoCSS from 'unocss/vite'
import AutoImport from 'unplugin-auto-import/vite'
import Unfonts from 'unplugin-fonts/vite'
import IconsResolver from 'unplugin-icons/resolver'
import Icons from 'unplugin-icons/vite'
import Components from 'unplugin-vue-components/vite'
import { defineConfig } from 'vite'
import Pages from 'vite-plugin-pages'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    AutoImport({
      imports: [
        'vue',
        'vue-router',
      ],
      dts: true,
      viteOptimizeDeps: true,
    }),
    vue(),
    // https://github.com/hannoeru/vite-plugin-pagesQ
    Pages(),
    Icons(),
    // https://github.com/antfu/unplugin-vue-components
    Components({
      // allow auto load markdown components under `./src/components/`
      extensions: ['vue', 'md'],
      dts: true,
      // allow auto import and register components used in markdown
      include: [/\.vue$/, /\.vue\?vue/, /\.md$/],
      // custom resolvers
      resolvers: [
        // auto import icons
        // https://github.com/antfu/unplugin-icons
        IconsResolver({
          prefix: 'icon',
          // enabledCollections: ['carbon']
        }),
      ],
    }),
    // https://github.com/antfu/unocss
    // see uno.config.ts for config
    UnoCSS(),
    Unfonts({
      google: {
        families: [
          'Montserrat',
          'Poppins',
          'Quicksand',
        ],
      },
    }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  assetsInclude: [
    '**/*.md',
  ],
  ssgOptions: {
    script: 'async',
    format: 'cjs',
    formatting: 'minify',
    beastiesOptions: {
      reduceInlineStyles: false,
    },
  },
})
