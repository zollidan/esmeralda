<script setup lang="ts">
import { ref, onMounted } from "vue";

type Task = {
  id: string;
  date: string;
  status: string;
  created_at: string;
};

const date = ref<string>("");
const tasks = ref<Task[]>([]);
const loading = ref(false);
const error = ref<string | null>(null);

async function fetchTasks() {
  loading.value = true;
  error.value = null;
  try {
    const res = await fetch("/api/tasks");
    if (!res.ok) throw new Error(`status ${res.status}`);
    tasks.value = await res.json();
  } catch (err: any) {
    error.value = err.message || String(err);
  } finally {
    loading.value = false;
  }
}

async function createTask() {
  if (!date.value) return;
  error.value = null;
  try {
    const res = await fetch("/api/tasks", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ date: date.value }),
    });
    if (!res.ok) throw new Error(`create failed: ${res.status}`);
    date.value = "";
    await fetchTasks();
  } catch (err: any) {
    error.value = err.message || String(err);
  }
}

onMounted(fetchTasks);
</script>

<template>
  <main class="container">
    <h1>Tasks</h1>

    <section class="create">
      <label for="date">Choose date</label>
      <input id="date" type="date" v-model="date" />
      <button @click="createTask">Create Task</button>
    </section>

    <section class="list">
      <h2>All tasks</h2>
      <div v-if="loading">Loading...</div>
      <div v-if="error" class="error">Error: {{ error }}</div>
      <ul v-if="tasks.length">
        <li v-for="t in tasks" :key="t.id">
          <strong>{{ t.date }}</strong> — {{ t.status }} —
          <small>{{ t.created_at }}</small>
        </li>
      </ul>
      <div v-else-if="!loading">No tasks</div>
    </section>
  </main>
</template>

<style scoped>
.container {
  max-width: 720px;
  margin: 2rem auto;
  font-family:
    system-ui,
    -apple-system,
    "Segoe UI",
    Roboto,
    "Helvetica Neue",
    Arial;
}
.create {
  display: flex;
  gap: 0.5rem;
  align-items: center;
  margin-bottom: 1rem;
}
input[type="date"] {
  padding: 0.4rem;
}
button {
  padding: 0.45rem 0.8rem;
}
.list ul {
  list-style: none;
  padding: 0;
}
.list li {
  padding: 0.5rem 0;
  border-bottom: 1px solid #eee;
}
.error {
  color: red;
}
</style>
