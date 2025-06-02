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
import { VitePWA } from 'vite-plugin-pwa'

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
    VitePWA(
      {
        manifest: {
          name: 'Zene',
          short_name: 'Zene',
          icons: [
            {
              src: '/zene.png',
              sizes: '1024x1024',
              type: 'image/png',
              purpose: 'any maskable',
            },
          ],
        },
        includeAssets: [
          '/default-image.jpg',
          '/default-square.png',
          '/favicon.ico',
          '/logo.png',
          '/zene.png',
        ],
        workbox: {
          runtimeCaching: [
            {
              urlPattern: ({ url }) => {
                return url.pathname.startsWith('/api') && !['/api/check-session'].includes(url.pathname)
              },
              handler: 'CacheFirst' as const,
              options: {
                cacheName: 'api-cache',
                cacheableResponse: {
                  statuses: [0, 200, 206],
                },
              },
            },
          ],
        },
      },
    ),
    Icons(),
    // https://github.com/antfu/unplugin-vue-components
    Components({
      extensions: ['vue'],
      dts: true,
      include: [/\.vue$/, /\.vue\?vue/],
      resolvers: [
        IconsResolver({
          prefix: 'icon',
        }),
      ],
    }),
    UnoCSS(),
    Unfonts({
      google: {
        families: [
          'Montserrat',
          'Poppins',
          'Quicksand',
          'Leto',
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
  optimizeDeps: {
    include: [
      'vue',
      'vue-router',
      '@vueuse/core',
    ],
    exclude: [],
  },
})
