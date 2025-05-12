import antfu from '@antfu/eslint-config'

export default antfu({
  vue: true,
  unocss: true,
  typescript: {
    tsconfigPath: 'tsconfig.json',
  },
  ignores: ['dist', '**/dist/**'],
}, {
  files: ['**/*.{vue,ts,js,json}'],
  rules: {
    'no-console': 'off',
    'brace-style': ['error', 'stroustrup'],
    'curly': ['off'],
    'vue/html-self-closing': 'off',
  },
})
