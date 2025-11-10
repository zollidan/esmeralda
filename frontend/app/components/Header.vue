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
  <header class="w-full h-12">
    <div class="flex justify-center items-center h-full">
      <ul class="flex gap-4">
        <li><NuxtLink to="/">Home</NuxtLink></li>
        <li><NuxtLink to="/files">Files</NuxtLink></li>
        <li>
          <button @click="logout">Logout</button>
        </li>
      </ul>
    </div>
  </header>
</template>
