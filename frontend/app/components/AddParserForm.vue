<script setup>
const username = ref("");
const password = ref("");

const submitForm = () => {
  const payload = {
    username: username.value,
    password: password.value,
  };
};
</script>

<template>
  <input type="text" v-model="username" placeholder="username" />
  <input type="password" v-model="password" placeholder="password" />
  <button @click="submitForm" type="submit">Submit</button>
</template>
