<script setup lang="ts">
import type { ProgressBar, Task } from '~/types/task'

const openRows = ref<Set<string>>(new Set())

function toggleRow(id: string) {
  if (openRows.value.has(id)) {
    openRows.value.delete(id)
  }
  else {
    openRows.value.add(id)
  }
  openRows.value = new Set(openRows.value)
}

const { progressByTaskId } = defineProps<{
  tasks: Task[]
  loading: boolean
  error: string | null
  progressByTaskId: Record<string, ProgressBar>
}>()

function getProgress(taskId: string) {
  return progressByTaskId[taskId]
}

const emit = defineEmits<{
  export: [taskDate: string]
  delete: [taskId: string]
}>()

function statusClass(status: string): string {
  if (status === 'pending') return 'bg-yellow-100 text-yellow-800'
  if (status === 'processing')
    return 'bg-blue-100 text-blue-800 flex items-center gap-1'
  if (status === 'done') return 'bg-green-100 text-green-800'
  if (status === 'error' || status === 'failed')
    return 'bg-red-100 text-red-800'
  if (status === 'cancelled') return 'bg-gray-100 text-gray-800'
  return 'bg-slate-100 text-slate-700'
}
</script>

<template>
  <section
    class="p-6 rounded-lg shadow-md transition-colors bg-white dark:bg-slate-800"
  >
    <h2 class="text-2xl font-semibold mb-4 text-slate-800 dark:text-white">
      Все задачи
    </h2>

    <!-- Loading -->
    <div
      v-if="loading"
      class="text-center py-8"
    >
      <div
        class="inline-block animate-spin rounded-full h-8 w-8 border-4 border-slate-300 dark:border-slate-600 border-t-blue-600"
      />
      <p class="mt-2 text-slate-600 dark:text-slate-400">
        Loading...
      </p>
    </div>

    <!-- Error -->
    <div
      v-if="error"
      class="mb-4 px-4 py-3 rounded-lg border bg-red-50 border-red-200 text-red-700 dark:bg-red-900/30 dark:border-red-800 dark:text-red-400"
    >
      <strong>Error:</strong> {{ error }}
    </div>

    <!-- Empty -->
    <div
      v-if="!loading && !tasks.length"
      class="text-center py-8 text-slate-500 dark:text-slate-400"
    >
      No tasks found
    </div>

    <!-- Table -->
    <div
      v-if="tasks.length"
      class="overflow-x-auto"
    >
      <table class="w-full">
        <thead
          class="border-b bg-slate-50 border-slate-200 dark:bg-slate-700 dark:border-slate-600"
        >
          <tr>
            <th
              class="px-4 py-3 text-left text-sm font-semibold text-slate-700 dark:text-slate-200"
            >
              День парсинга
            </th>
            <th
              class="px-4 py-3 text-left text-sm font-semibold text-slate-700 dark:text-slate-200"
            >
              Статус
            </th>
            <th
              class="px-4 py-3 text-left text-sm font-semibold text-slate-700 dark:text-slate-200"
            >
              Дата создания
            </th>
            <th class="px-4 py-3 w-10" />
          </tr>
        </thead>

        <tbody>
          <template
            v-for="task in tasks"
            :key="task.id"
          >
            <!-- Основная строка -->
            <tr
              class="border-b transition border-slate-100 hover:bg-slate-50 dark:border-slate-700 dark:hover:bg-slate-700"
            >
              <td
                class="px-4 py-3 text-sm font-medium text-slate-900 dark:text-white"
              >
                {{ new Date(task.date).toLocaleDateString("ru-RU") }}
              </td>

              <td class="px-4 py-3 text-sm">
                <span
                  :class="statusClass(task.status)"
                  class="px-2 py-1 rounded-full text-xs font-medium inline-flex items-center gap-1"
                >
                  <svg
                    v-if="task.status === 'processing'"
                    class="animate-spin h-3 w-3"
                    viewBox="0 0 24 24"
                  >
                    <circle
                      class="opacity-25"
                      cx="12"
                      cy="12"
                      r="10"
                      stroke="currentColor"
                      stroke-width="4"
                      fill="none"
                    />
                    <path
                      class="opacity-75"
                      fill="currentColor"
                      d="M4 12a8 8 0 018-8v4l3-3-3-3v4a10 10 0 00-10 10h2z"
                    />
                  </svg>

                  <span>
                    {{ task.status }}
                    <template
                      v-if="
                        task.status === 'processing' && getProgress(task.id)
                      "
                    >
                      ({{ getProgress(task.id)?.percent }}%)
                    </template>
                  </span>
                </span>
              </td>

              <td class="px-4 py-3 text-sm text-slate-600 dark:text-slate-400">
                {{ new Date(task.created_at).toLocaleString("ru-RU") }}
              </td>

              <td class="px-4 py-3">
                <button
                  class="w-7 h-7 flex items-center justify-center rounded-md border transition-all duration-200 border-slate-200 text-slate-400 hover:bg-slate-100 hover:text-slate-600 dark:border-slate-600 dark:text-slate-300 dark:hover:bg-slate-600 dark:hover:text-white"
                  :class="{
                    'rotate-180 bg-slate-100 text-slate-600 dark:bg-slate-600 dark:text-white':
                      openRows.has(task.id),
                  }"
                  @click="toggleRow(task.id)"
                >
                  <svg
                    width="12"
                    height="12"
                    viewBox="0 0 12 12"
                  >
                    <path
                      d="M2 4l4 4 4-4"
                      stroke="currentColor"
                      stroke-width="1.5"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    />
                  </svg>
                </button>
              </td>
            </tr>

            <!-- Expand -->
            <tr
              v-if="openRows.has(task.id)"
              class="border-b bg-slate-50 border-slate-200 dark:bg-slate-700 dark:border-slate-600"
            >
              <td
                colspan="4"
                class="px-4 py-4 pl-14"
              >
                <div class="flex flex-col gap-4">
                  <!-- ID -->
                  <div>
                    <p
                      class="text-xs font-semibold uppercase tracking-wide mb-1 text-slate-400"
                    >
                      ID задачи
                    </p>
                    <p
                      class="text-xs font-mono px-2 py-1 rounded inline-block bg-slate-200 text-slate-600 dark:bg-slate-600 dark:text-slate-200"
                    >
                      {{ task.id }}
                    </p>
                  </div>

                  <!-- Progress -->
                  <div v-if="progressByTaskId[task.id]">
                    <p
                      class="text-xs font-semibold uppercase tracking-wide mb-2 text-slate-400"
                    >
                      Прогресс
                    </p>

                    <div
                      class="h-2 w-56 rounded-full overflow-hidden bg-slate-200 dark:bg-slate-600"
                    >
                      <div
                        class="h-full bg-blue-500 transition-all duration-500"
                        :style="{
                          width: `${progressByTaskId[task.id]?.percent}%`,
                        }"
                      />
                    </div>

                    <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">
                      {{ progressByTaskId[task.id]?.message }}
                      ({{ progressByTaskId[task.id]?.percent }}%)
                    </p>
                  </div>

                  <!-- Buttons -->
                  <div class="flex gap-2">
                    <button
                      :disabled="task.status !== 'done'"
                      class="px-4 py-2 rounded-lg text-sm font-medium text-white transition bg-green-600 hover:bg-green-700 disabled:bg-slate-300 dark:disabled:bg-slate-600"
                      @click="emit('export', task.date)"
                    >
                      Скачать Excel
                    </button>

                    <button
                      :disabled="task.status !== 'done'"
                      class="px-4 py-2 rounded-lg text-sm font-medium text-white transition bg-red-700 hover:bg-red-800 disabled:bg-slate-300 dark:disabled:bg-slate-600"
                      @click="emit('delete', task.id)"
                    >
                      Удалить задачу
                    </button>
                  </div>
                </div>
              </td>
            </tr>
          </template>
        </tbody>
      </table>
    </div>
  </section>
</template>
