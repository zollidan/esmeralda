import tailwindcss from '@tailwindcss/vite'

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  modules: ['@nuxt/eslint', 'nuxt-auth-utils', '@nuxtjs/color-mode'],

  devtools: { enabled: true },
  css: ['./app/assets/css/main.css'],
  colorMode: {
    classSuffix: '',
  },
  runtimeConfig: {
    apiBase: 'http://localhost:8080',
  },
  compatibilityDate: '2025-07-15',

  vite: {
    plugins: [tailwindcss()],
    optimizeDeps: {
      include: ['@vue/devtools-core', '@vue/devtools-kit'],
    },
  },
  eslint: {
    config: {
      stylistic: true,
    },
  },
})
