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
  <section class="bg-white rounded-lg shadow-md p-6">
    <h2 class="text-2xl font-semibold text-slate-800 mb-4">
      All Tasks
    </h2>

    <div
      v-if="loading"
      class="text-center py-8"
    >
      <div
        class="inline-block animate-spin rounded-full h-8 w-8 border-4 border-slate-300 border-t-blue-600"
      />
      <p class="mt-2 text-slate-600">
        Loading...
      </p>
    </div>

    <div
      v-if="error"
      class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-4"
    >
      <strong>Error:</strong> {{ error }}
    </div>

    <div
      v-if="!loading && !tasks.length"
      class="text-center py-8 text-slate-500"
    >
      No tasks found
    </div>

    <div
      v-if="tasks.length"
      class="overflow-x-auto"
    >
      <table class="w-full">
        <thead class="bg-slate-50 border-b border-slate-200">
          <tr>
            <th
              class="px-4 py-3 text-left text-sm font-semibold text-slate-700"
            >
              День парсинга
            </th>
            <th
              class="px-4 py-3 text-left text-sm font-semibold text-slate-700"
            >
              Статус
            </th>
            <th
              class="px-4 py-3 text-left text-sm font-semibold text-slate-700"
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
            <tr class="hover:bg-slate-50 transition border-b border-slate-100">
              <td class="px-4 py-3 text-sm font-medium text-slate-900">
                {{ new Date(task.date).toLocaleDateString() }}
              </td>
              <td class="px-4 py-3 text-sm">
                <span
                  :class="statusClass(task.status)"
                  class="px-2 py-1 rounded-full text-xs font-medium inline-flex items-center gap-1"
                >
                  <!-- Спиннер -->
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

                  <!-- Текст -->
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
              <td class="px-4 py-3 text-sm text-slate-600">
                {{ new Date(task.created_at).toLocaleString() }}
              </td>
              <td class="px-4 py-3">
                <button
                  class="w-7 h-7 flex items-center justify-center rounded-md border border-slate-200 text-slate-400 hover:bg-slate-100 hover:text-slate-600 transition-all duration-200"
                  :class="{
                    'rotate-180 bg-slate-100 text-slate-600': openRows.has(
                      task.id,
                    ),
                  }"
                  @click="toggleRow(task.id)"
                >
                  <svg
                    width="12"
                    height="12"
                    viewBox="0 0 12 12"
                    fill="none"
                    xmlns="http://www.w3.org/2000/svg"
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

            <!-- Expandable строка -->
            <tr
              v-if="openRows.has(task.id)"
              class="bg-slate-50 border-b border-slate-200"
            >
              <td
                colspan="4"
                class="px-4 py-4 pl-14"
              >
                <div class="flex flex-col gap-4">
                  <!-- Task ID -->
                  <div>
                    <p
                      class="text-xs font-semibold text-slate-400 uppercase tracking-wide mb-1"
                    >
                      Task ID
                    </p>
                    <p
                      class="text-xs font-mono text-slate-600 bg-slate-100 inline-block px-2 py-1 rounded"
                    >
                      {{ task.id }}
                    </p>
                  </div>

                  <!-- Прогресс бар -->
                  <div v-if="progressByTaskId[task.id]">
                    <p
                      class="text-xs font-semibold text-slate-400 uppercase tracking-wide mb-2"
                    >
                      Прогресс
                    </p>
                    <div
                      class="h-2 w-56 bg-slate-200 rounded-full overflow-hidden"
                    >
                      <div
                        class="h-full bg-blue-500 transition-all duration-500 rounded-full"
                        :style="{ width: `${progressByTaskId[task.id]?.percent}%` }"
                      />
                    </div>
                    <p class="mt-1 text-xs text-slate-500">
                      {{ progressByTaskId[task.id]?.message }}
                      ({{ progressByTaskId[task.id]?.percent }}%)
                    </p>
                  </div>

                  <!-- Кнопки -->
                  <div class="flex gap-2">
                    <button
                      :disabled="task.status !== 'done'"
                      class="px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-lg cursor-pointer hover:bg-green-700 disabled:bg-slate-300 disabled:cursor-not-allowed transition inline-flex items-center gap-2"
                      @click="emit('export', task.date)"
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-4 w-4"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                        />
                      </svg>
                      Export
                    </button>
                    <button
                      :disabled="task.status !== 'done'"
                      class="px-4 py-2 bg-red-800 text-white text-sm font-medium rounded-lg cursor-pointer hover:bg-red-700 disabled:bg-slate-300 disabled:cursor-not-allowed transition inline-flex items-center gap-2"
                      @click="emit('delete', task.id)"
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-4 w-4"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                        />
                      </svg>
                      Delete
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
