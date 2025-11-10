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
  <div>
    <h1>Welcome {{ user.name }}</h1>
    <button @click="logout">Logout</button>
  </div>
</template>
