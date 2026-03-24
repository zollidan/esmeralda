<script setup lang="ts">
import type { ProgressBar } from '~/types/task'

const { date, tasks, loading, error, fetchTasks, createTask, deleteTask, exportToExcel } = useTasks()

const visible = ref(false)
const progressByTaskId = ref<Record<string, ProgressBar>>({})

onMounted(fetchTasks)
</script>

<template>
  <main class="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100">
    <div class="max-w-6xl mx-auto px-4 py-8">
      <h1 class="text-4xl font-bold text-slate-800 mb-8">
        aaf-bet.ru
      </h1>

      <button
        class="px-6 py-2 my-4 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 disabled:bg-slate-300 disabled:cursor-not-allowed transition"
        @click="visible = !visible"
      >
        {{ visible ? "Скрыть" : "Создать новую задачу" }}
      </button>

      <TaskCreateForm
        v-model:date="date"
        :visible="visible"
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
