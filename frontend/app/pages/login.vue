<script setup lang="ts">
definePageMeta({
  layout: "auth",
});
const { loggedIn, user, fetch: refreshSession } = useUserSession();
const credentials = reactive({
  username: "",
  password: "",
});
async function login() {
  try {
    await $fetch("/api/login", {
      method: "POST",
      body: credentials,
    });

    await refreshSession();
    await navigateTo("/");
  } catch {
    alert("Bad credentials");
  }
}
</script>

<template>
  <form @submit.prevent="login">
    <input v-model="credentials.username" type="text" placeholder="Username" />
    <input
      v-model="credentials.password"
      type="password"
      placeholder="Password"
    />
    <button type="submit">Login</button>
  </form>
</template>
