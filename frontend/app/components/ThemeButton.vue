<script setup>
const colorMode = useColorMode()

const toggleTheme = () => {
  colorMode.preference = colorMode.value === 'dark' ? 'light' : 'dark'
}
</script>

<template>
  <button
    class="fixed bottom-5 right-5 z-50 w-12 h-12 flex items-center justify-center rounded-full shadow-lg bg-slate-200 dark:bg-slate-700 hover:scale-110 transition"
    @click="toggleTheme"
  >
    <span v-if="colorMode.value === 'dark'">🌙</span>
    <span v-else>☀️</span>
  </button>
</template>
