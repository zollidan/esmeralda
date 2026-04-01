<script setup lang="ts">
const { login } = useAuth()
const form = reactive({ username: '', password: '' })
const error = ref('')
const loading = ref(false)

const handleLogin = async () => {
  loading.value = true
  error.value = ''
  try {
    await login(form)
    navigateTo('/')
  }
  catch {
    error.value = 'Неверный логин или пароль'
  }
  finally {
    loading.value = false
  }
}
</script>

<template>
  <div
    class="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-900 dark:to-slate-800 transition-colors duration-300"
  >
    <div
      class="w-full max-w-sm p-8 rounded-2xl shadow-md bg-white dark:bg-slate-800 transition-colors"
    >
      <h1
        class="text-2xl font-bold mb-6 text-center text-slate-800 dark:text-white"
      >
        Вход
      </h1>

      <div class="flex flex-col gap-4">
        <input
          v-model="form.username"
          placeholder="Логин"
          class="px-4 py-2 rounded-lg border border-slate-300 dark:border-slate-600 bg-white dark:bg-slate-700 text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >

        <input
          v-model="form.password"
          type="password"
          placeholder="Пароль"
          class="px-4 py-2 rounded-lg border border-slate-300 dark:border-slate-600 bg-white dark:bg-slate-700 text-slate-900 dark:text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >

        <p
          v-if="error"
          class="text-red-500 text-sm"
        >
          {{ error }}
        </p>

        <button
          :disabled="loading"
          class="px-6 py-2 rounded-lg font-medium text-white bg-blue-600 hover:bg-blue-700 disabled:bg-slate-300 dark:disabled:bg-slate-600 disabled:cursor-not-allowed transition"
          @click="handleLogin"
        >
          {{ loading ? "Вход..." : "Войти" }}
        </button>
      </div>
    </div>
  </div>
</template>
