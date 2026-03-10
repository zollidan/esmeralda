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

async function exportToExcel(taskId: string) {
  try {
    const res = await fetch(`/api/tasks/${taskId}/export`, {
      method: "GET",
    });
    if (!res.ok) throw new Error(`export failed: ${res.status}`);

    const blob = await res.blob();
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `task-${taskId}.xlsx`;
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
    document.body.removeChild(a);
  } catch (err: any) {
    error.value = err.message || String(err);
  }
}

onMounted(fetchTasks);
</script>

<template>
  <main class="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100">
    <div class="max-w-6xl mx-auto px-4 py-8">
      <h1 class="text-4xl font-bold text-slate-800 mb-8">
        Football Match Statistics
      </h1>

      <section class="bg-white rounded-lg shadow-md p-6 mb-8">
        <h2 class="text-xl font-semibold text-slate-700 mb-4">
          Create New Task
        </h2>
        <div class="flex gap-4 items-end">
          <div class="flex-1">
            <label
              for="date"
              class="block text-sm font-medium text-slate-700 mb-2"
            >
              Choose date
            </label>
            <input
              id="date"
              type="date"
              v-model="date"
              class="w-full px-4 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition"
            />
          </div>
          <button
            @click="createTask"
            :disabled="!date"
            class="px-6 py-2 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 disabled:bg-slate-300 disabled:cursor-not-allowed transition"
          >
            Create Task
          </button>
        </div>
      </section>

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
                  Date
                </th>
                <th
                  class="px-4 py-3 text-left text-sm font-semibold text-slate-700"
                >
                  Status
                </th>
                <th
                  class="px-4 py-3 text-left text-sm font-semibold text-slate-700"
                >
                  Created At
                </th>
                <th
                  class="px-4 py-3 text-left text-sm font-semibold text-slate-700"
                >
                  Actions
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-200">
              <tr
                v-for="t in tasks"
                :key="t.id"
                class="hover:bg-slate-50 transition"
              >
                <td class="px-4 py-3 text-sm font-medium text-slate-900">
                  {{ t.date }}
                </td>
                <td class="px-4 py-3 text-sm">
                  <span
                    :class="{
                      'bg-yellow-100 text-yellow-800': t.status === 'pending',
                      'bg-blue-100 text-blue-800': t.status === 'processing',
                      'bg-green-100 text-green-800': t.status === 'done',
                      'bg-red-100 text-red-800':
                        t.status === 'error' || t.status === 'failed',
                      'bg-gray-100 text-gray-800': t.status === 'cancelled',
                    }"
                    class="px-2 py-1 rounded-full text-xs font-medium"
                  >
                    {{ t.status }}
                  </span>
                </td>
                <td class="px-4 py-3 text-sm text-slate-600">
                  {{ new Date(t.created_at).toLocaleString() }}
                </td>
                <td class="px-4 py-3 text-sm">
                  <button
                    @click="exportToExcel(t.id)"
                    :disabled="t.status !== 'done'"
                    class="px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-lg hover:bg-green-700 disabled:bg-slate-300 disabled:cursor-not-allowed transition inline-flex items-center gap-2"
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
    </div>
  </main>
</template>
