<script setup lang="ts">
definePageMeta({
  middleware: ["authenticated"],
});

const { user, clear: clearSession } = useUserSession();

async function logout() {
  await clearSession();
  await navigateTo("/login");
}
</script>

<template>
  <header class="w-full bg-white shadow-sm border-b sticky top-0 z-50">
    <div class="max-w-4xl mx-auto px-4 h-14 flex justify-between items-center">
      <!-- Лого и имя -->
      <div class="flex items-center gap-2">
        <span class="font-semibold tracking-wide">Esmeralda</span>
      </div>

      <!-- Навигация -->
      <nav>
        <ul class="flex gap-6 items-center text-sm font-medium">
          <li>
            <NuxtLink
              to="/"
              class="hover:text-blue-600 transition-colors"
              active-class="text-blue-600"
            >
              Home
            </NuxtLink>
          </li>
          <li>
            <NuxtLink
              to="/files"
              class="hover:text-blue-600 transition-colors"
              active-class="text-blue-600"
            >
              Files
            </NuxtLink>
          </li>
          <li>
            <button
              @click="logout"
              class="text-red-600 hover:text-red-700 transition-colors"
            >
              Logout
            </button>
          </li>
        </ul>
      </nav>
    </div>
  </header>
</template>
