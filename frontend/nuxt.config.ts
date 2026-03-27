import tailwindcss from '@tailwindcss/vite'

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({

  modules: ['@nuxt/eslint'],

  devtools: { enabled: true },
  css: ['./app/assets/css/main.css'], runtimeConfig: {
    apiBase: 'http://localhost:8080',
  }, compatibilityDate: '2025-07-15',

  vite: {
    plugins: [tailwindcss()],
  },
  eslint: {
    config: {
      stylistic: true,
    },
  },
})
