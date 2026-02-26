import path from 'node:path'
import vue from '@vitejs/plugin-vue'
import UnoCSS from 'unocss/vite'
import AutoImport from 'unplugin-auto-import/vite'
import Unfonts from 'unplugin-fonts/vite'
import IconsResolver from 'unplugin-icons/resolver'
import Icons from 'unplugin-icons/vite'
import Components from 'unplugin-vue-components/vite'
import { defineConfig } from 'vite'
import VueRouter from 'vue-router/vite'

export default defineConfig({
  plugins: [
    VueRouter({
      dts: 'route-map.d.ts',
    }),
    AutoImport({
      imports: [
        'vue',
        'vue-router',
      ],
      dts: true,
      viteOptimizeDeps: true,
    }),
    vue({
      // template: {
      //   compilerOptions: {
      //     isCustomElement: tag => tag === 'google-cast-launcher',
      //   },
      // },
    }),
    Unfonts({
      google: {
        families: [
          {
            name: 'Jura',
            styles: 'wght@300..700',
            defer: true,
          },
        ],
      },
    }),
    Icons(),
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
  ],
  resolve: {
    alias: {
      '~/': `${path.resolve(__dirname, 'src')}/`,
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
