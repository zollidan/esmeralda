import tailwindcss from "@tailwindcss/vite";

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: "2025-07-15",
  css: ["./app/assets/css/main.css"],

  vite: {
    plugins: [tailwindcss()],
  },

  devtools: { enabled: true },
  modules: ["@nuxt/eslint"],
  eslint: {
    config: {
      stylistic: true,
    },
  },
});
