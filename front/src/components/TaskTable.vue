<script setup lang="ts">
import type { ProgressBar, Task } from "../types/task";

defineProps<{
  tasks: Task[];
  loading: boolean;
  error: string | null;
  progressByTaskId: Record<string, ProgressBar>;
}>();

const emit = defineEmits<{
  export: [taskDate: string];
}>();

function statusClass(status: string): string {
  if (status === "pending") return "bg-yellow-100 text-yellow-800";
  if (status === "processing") return "bg-blue-100 text-blue-800";
  if (status === "done") return "bg-green-100 text-green-800";
  if (status === "error" || status === "failed")
    return "bg-red-100 text-red-800";
  if (status === "cancelled") return "bg-gray-100 text-gray-800";
  return "bg-slate-100 text-slate-700";
}
</script>

<template>
  <section class="bg-white rounded-lg shadow-md p-6">
    <h2 class="text-2xl font-semibold text-slate-800 mb-4">All Tasks</h2>

    <div v-if="loading" class="text-center py-8">
      <div
        class="inline-block animate-spin rounded-full h-8 w-8 border-4 border-slate-300 border-t-blue-600"
      ></div>
      <p class="mt-2 text-slate-600">Loading...</p>
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

    <div v-if="tasks.length" class="overflow-x-auto">
      <table class="w-full">
        <thead class="bg-slate-50 border-b border-slate-200">
          <tr>
            <th
              class="px-4 py-3 text-left text-sm font-semibold text-slate-700"
            >
              Дата парсинга
            </th>
            <th
              class="px-4 py-3 text-left text-sm font-semibold text-slate-700"
            >
              Статус
            </th>
            <th
              class="px-4 py-3 text-left text-sm font-semibold text-slate-700"
            >
              Запущено в
            </th>
            <th
              class="px-4 py-3 text-left text-sm font-semibold text-slate-700"
            >
              Экспорт
            </th>
          </tr>
        </thead>

        <tbody class="divide-y divide-slate-200">
          <tr
            v-for="task in tasks"
            :key="task.id"
            class="hover:bg-slate-50 transition"
          >
            <td class="px-4 py-3 text-sm font-medium text-slate-900">
              {{ new Date(task.date).toLocaleDateString() }}
            </td>
            <td class="px-4 py-3 text-sm">
              <span
                :class="statusClass(task.status)"
                class="px-2 py-1 rounded-full text-xs font-medium"
              >
                {{ task.status }}
              </span>

              <!-- progress bar temporarily hidden -->
              <!-- <div v-if="progressByTaskId[task.id]" class="mt-2">
                <div class="h-2 w-48 bg-slate-200 rounded-full overflow-hidden">
                  <div
                    class="h-full bg-blue-500 transition-all"
                    :style="{
                      width: `${progressPercent(progressByTaskId[task.id])}%`,
                    }"
                  ></div>
                </div>
                <p class="mt-1 text-xs text-slate-500">
                  {{ progressByTaskId[task.id]?.current_match ?? 0 }} /
                  {{ progressByTaskId[task.id]?.total_matches ?? 0 }}
                </p>
              </div> -->
            </td>
            <td class="px-4 py-3 text-sm text-slate-600">
              {{ new Date(task.created_at).toLocaleString() }}
            </td>
            <td class="px-4 py-3 text-sm">
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
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>
