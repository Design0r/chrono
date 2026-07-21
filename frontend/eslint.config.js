import js from '@eslint/js'
import globals from 'globals'
import reactHooks from 'eslint-plugin-react-hooks'
import reactRefresh from 'eslint-plugin-react-refresh'
import tseslint from 'typescript-eslint'
import { defineConfig, globalIgnores } from 'eslint/config'

export default defineConfig([
  globalIgnores(['dist']),
  {
    files: ['**/*.{ts,tsx}'],
    extends: [
      js.configs.recommended,
      tseslint.configs.recommended,
      reactHooks.configs.flat['recommended-latest'],
      reactRefresh.configs.vite,
    ],
    languageOptions: {
      ecmaVersion: 2020,
      globals: globals.browser,
    },
  },
  {
    // Route-Dateien exportieren per Konvention `Route` neben der Komponente,
    // Provider-Dateien ihren Hook. Fast Refresh verliert dort etwas Komfort,
    // vermeiden lässt es sich bei diesen Frameworks aber nicht.
    files: [
      'src/routes/**/*.tsx',
      'src/main.tsx',
      'src/auth.tsx',
      'src/components/Toast.tsx',
      'src/integrations/**/*.tsx',
    ],
    rules: {
      'react-refresh/only-export-components': 'off',
    },
  },
])
