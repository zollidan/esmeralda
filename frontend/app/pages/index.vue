<script setup lang="ts">
definePageMeta({
  middleware: ["authenticated"],
});

const { user, clear: clearSession } = useUserSession();
</script>

<template>
  <div>
    <h1>Welcome {{ user.name }}</h1>
  </div>
</template>
