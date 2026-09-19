import antfu from '@antfu/eslint-config'

export default antfu({
  ignores: ['**/.pnpm-store/**', '**/dist/**', '**/node_modules/**'],
  unocss: true,
  rules: {
    'no-console': 'warn',
    'curly': 'off',
    '@typescript-eslint/brace-style': 'off',
    'unused-imports/no-unused-imports': 'off',
    'node/prefer-global/process': 'off',
  },
}, {
  ignores: ['**/.pnpm-store/**', '**/dist/**', '**/node_modules/**'],
})
