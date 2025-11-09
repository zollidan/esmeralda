<script setup>
import { ref } from "vue";
import { useCookie } from "#app";

const credentials = ref({
  username: "",
  password: "",
});

const token = useCookie("auth_token");

async function handleLogin() {
  try {
    const response = await $fetch("/api/login", {
      method: "POST",
      body: {
        username: credentials.value.username,
        password: credentials.value.password,
      },
    });

    console.log("Response:", response);

    if (response?.token) {
      console.log("Login successful");
      token.value = response.token;
    } else {
      console.log("Login failed — no token in response");
    }
  } catch (error) {
    console.error("Error during login:", error);
  }
}
</script>

<template>
  <div>
    <h1>Login Page</h1>
    <input type="text" placeholder="Username" v-model="credentials.username" />
    <input
      type="password"
      placeholder="Password"
      v-model="credentials.password"
    />
    <button type="submit" @click.prevent="handleLogin">Login</button>
  </div>
</template>
