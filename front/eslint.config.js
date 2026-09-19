import antfu from '@antfu/eslint-config'

export default antfu({
  ignores: ['**/.pnpm-store/**', '**/dist/**', '**/node_modules/**'],
  unocss: true,
  rules: {
    'no-console': 'warn',
    'curly': 'off',
    '@typescript-eslint/brace-style': 'off',
    'node/prefer-global/process': 'off',
    'unused-imports/no-unused-imports': 'off',
  },
}, {
  ignores: ['**/.pnpm-store/**', '**/dist/**', '**/node_modules/**'],
})
