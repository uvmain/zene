/// <reference types="vite-ssg" />
/// <reference types="vitest" />

import path from 'node:path'
import vue from '@vitejs/plugin-vue'
import UnoCSS from 'unocss/vite'
import AutoImport from 'unplugin-auto-import/vite'
import IconsResolver from 'unplugin-icons/resolver'
import Icons from 'unplugin-icons/vite'
import Components from 'unplugin-vue-components/vite'
import { defineConfig } from 'vite'

export default defineConfig({
  plugins: [
    AutoImport({
      imports: [
        'vue',
        'vue-router',
        'vitest',
      ],
      dts: true,
      viteOptimizeDeps: true,
    }),
    vue({
      template: {
        compilerOptions: {
          isCustomElement: tag => tag === 'google-cast-launcher',
        },
      },
    }),
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
  ],
  resolve: {
    alias: {
      '~/': `${path.resolve(__dirname, 'src')}/`,
    },
  },
  assetsInclude: [
    '**/*.md',
  ],
  ssgOptions: {
    script: 'async',
    formatting: 'minify',
  },
  optimizeDeps: {
    include: [
      'vue',
      'vue-router',
      '@vueuse/core',
    ],
    exclude: [],
  },
  // Vitest configuration
  test: {
    environment: 'happy-dom',
    globals: true,
    include: ['src/**/*.{test,spec}.{js,mjs,cjs,ts,mts,cts,jsx,tsx}'],
    exclude: ['node_modules', 'dist', '.git', '.cache'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      exclude: [
        'node_modules/',
        'test/',
        '**/*.d.ts',
        '**/*.config.*',
        'dist/',
      ],
    },
  },
})
