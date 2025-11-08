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
  <div v-if="data == undefined">No data</div>
  <div v-else>
    <ul>
      <li v-for="parser in data" :key="parser.id">
        <h3>
          {{ parser.name }} <span v-if="parser.isActive">(Active)</span
          ><span v-else>(Inactive)</span>
        </h3>
        <p>{{ parser.description }}</p>
      </li>
    </ul>
  </div>
  <div>
    <button @click="toggleAddNewParser">Add new parser</button>
  </div>
  <div v-if="showAddParserForm">
    <form action="">
      <input type="text" placeholder="Parser Name" />
      <input type="text" placeholder="Parser Description" />
      <button type="submit">Submit</button>
    </form>
  </div>
</template>
