<script setup>
import { ref } from "vue";

const emit = defineEmits(["parserAdded"]);

const payload = ref({
  name: "",
  description: "",
});

const submitForm = async () => {
  await $fetch("/api/parsers", {
    method: "POST",
    body: payload.value,
  });

  emit("parserAdded");
};
</script>

<template>
  <input type="text" v-model="payload.name" placeholder="name" />
  <input type="text" v-model="payload.description" placeholder="description" />
  <button @click="submitForm">Submit</button>
</template>
