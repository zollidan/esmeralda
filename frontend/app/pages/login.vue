<!-- pages/login.vue -->
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
    class="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 flex items-center justify-center"
  >
    <div class="bg-white rounded-2xl shadow-md p-8 w-full max-w-sm">
      <h1 class="text-2xl font-bold text-slate-800 mb-6 text-center">
        Вход
      </h1>
      <div class="flex flex-col gap-4">
        <input
          v-model="form.username"
          placeholder="Логин"
          class="px-4 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
        <input
          v-model="form.password"
          type="password"
          placeholder="Пароль"
          class="px-4 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
        <p
          v-if="error"
          class="text-red-500 text-sm"
        >
          {{ error }}
        </p>
        <button
          :disabled="loading"
          class="px-6 py-2 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 disabled:bg-slate-300 disabled:cursor-not-allowed transition"
          @click="handleLogin"
        >
          {{ loading ? "Вход..." : "Войти" }}
        </button>
      </div>
    </div>
  </div>
</template>
