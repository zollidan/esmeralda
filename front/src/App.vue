<script setup lang="ts">
import { onMounted, ref } from "vue";
import TaskCreateForm from "./components/TaskCreateForm.vue";
import TaskTable from "./components/TaskTable.vue";
import type { ProgressBar, Task } from "./types/task";
import { createTaskActions } from "../utils/utils";

const new_task = ref<boolean>(false);
const date = ref<string>("");
const tasks = ref<Task[]>([]);
const loading = ref(false);
const error = ref<string | null>(null);
const progressByTaskId = ref<Record<string, ProgressBar>>({});

const { fetchTasks, createTask, deleteTask, exportToExcel } = createTaskActions(
  {
    date,
    tasks,
    loading,
    error,
  },
);

onMounted(fetchTasks);
</script>

<template>
  <main class="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100">
    <div class="max-w-6xl mx-auto px-4 py-8">
      <h1 class="text-4xl font-bold text-slate-800 mb-8">aaf-bet.ru</h1>

      <button
        class="px-6 py-2 my-4 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 disabled:bg-slate-300 disabled:cursor-not-allowed transition"
        @click="new_task = !new_task"
      >
        {{ new_task ? "Скрыть" : "Создать новую задачу" }}
      </button>

      <TaskCreateForm
        v-model:date="date"
        :visible="new_task"
        @submit="createTask"
      />

      <TaskTable
        :tasks="tasks"
        :loading="loading"
        :error="error"
        :progress-by-task-id="progressByTaskId"
        @export="exportToExcel"
        @delete="deleteTask"
      />
    </div>
  </main>
</template>
