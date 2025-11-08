<script setup>
import { ref } from "vue";

const showAddParserForm = ref(false);

function toggleAddNewParser() {
  showAddParserForm.value = !showAddParserForm.value;
}

const { data } = await useFetch("/api/parsers", {
  method: "GET",
});
</script>

<template>
  <div v-if="data == undefined || data.length === 0">No data</div>
  <div v-else>
    <ul>
      <li v-for="parser in data" :key="parser.ID">
        <h3>
          {{ parser.Name }}
          <span v-if="parser.IsActive">(Active)</span>
          <span v-else>(Inactive)</span>
        </h3>
        <p>{{ parser.Description }}</p>
      </li>
    </ul>
  </div>
  <div>
    <button @click="toggleAddNewParser">Add new parser</button>
  </div>
  <div v-if="showAddParserForm">
    <AddParserForm />
  </div>
</template>
