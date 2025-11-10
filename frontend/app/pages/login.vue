<script setup lang="ts">
definePageMeta({ layout: "auth" });

const { loggedIn, user, fetch: refreshSession } = useUserSession();

const credentials = reactive({
  username: "",
  password: "",
});

const error = ref("");

async function login() {
  error.value = "";
  try {
    await $fetch("/api/login", {
      method: "POST",
      body: credentials,
    });

    await refreshSession();
    await navigateTo("/");
  } catch {
    error.value = "Invalid username or password";
  }
}
</script>
<template>
  <div class="flex justify-center items-center min-h-screen bg-gray-50">
    <form
      @submit.prevent="login"
      class="w-72 border rounded-xl p-6 flex flex-col justify-center gap-3 shadow-sm bg-white"
    >
      <!-- Название -->
      <h1 class="text-center text-xl font-semibold mb-2 tracking-wide">
        Esmeralda
      </h1>

      <!-- Поля -->
      <input
        v-model="credentials.username"
        type="text"
        placeholder="Username"
        class="border rounded p-2 focus:ring-2 focus:ring-blue-500 outline-none"
      />
      <input
        v-model="credentials.password"
        type="password"
        placeholder="Password"
        class="border rounded p-2 focus:ring-2 focus:ring-blue-500 outline-none"
      />

      <!-- Ошибка -->
      <p v-if="error" class="text-red-600 text-sm text-center">
        {{ error }}
      </p>

      <!-- Кнопка -->
      <button
        type="submit"
        class="bg-blue-600 text-white rounded p-2 mt-2 hover:bg-blue-700 transition"
      >
        Login
      </button>
    </form>
  </div>
</template>
