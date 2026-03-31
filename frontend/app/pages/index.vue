<script setup lang="ts">
definePageMeta({
  middleware: 'auth',
})

const {
  date,
  tasks,
  loading,
  error,
  initStream,
  progressByTaskId,
  createTask,
  deleteTask,
  exportToExcel,
} = useTasks()

const { logout } = useAuth()

const visible = ref(false)
let es: { close: () => void } | null = null

onMounted(async () => {
  es = initStream()
})

onUnmounted(() => {
  es?.close()
})
</script>

<template>
  <main
    class="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-900 dark:to-slate-800 transition-colors duration-300"
  >
    <div class="max-w-6xl mx-auto px-4 py-8">
      <div class="flex items-center justify-between mb-8">
        <h1 class="text-4xl font-bold text-slate-800 dark:text-white">
          aaf-bet.ru
        </h1>

        <button
          class="px-4 py-2 rounded-lg font-medium transition bg-slate-200 text-slate-700 hover:bg-slate-300 dark:bg-slate-700 dark:text-slate-200 dark:hover:bg-slate-600"
          @click="logout"
        >
          Выйти
        </button>
      </div>

      <button
        class="px-6 py-2 my-4 rounded-lg font-medium text-white transition bg-blue-600 hover:bg-blue-700 disabled:bg-slate-300 dark:disabled:bg-slate-600 disabled:cursor-not-allowed"
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
