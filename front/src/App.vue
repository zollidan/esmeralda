<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, watch } from "vue";
import TaskCreateForm from "./components/TaskCreateForm.vue";
import TaskTable from "./components/TaskTable.vue";
import type { ProgressBar, Task } from "./types/task";

const new_task = ref<boolean>(false);
const date = ref<string>("");
const tasks = ref<Task[]>([]);
const loading = ref(false);
const error = ref<string | null>(null);
const progressByTaskId = ref<Record<string, ProgressBar>>({});
const sockets = new Map<string, WebSocket>();

function wsUrl(taskId: string): string {
  const protocol = window.location.protocol === "https:" ? "wss" : "ws";
  return `${protocol}://${window.location.host}/api/progress/ws?task_id=${encodeURIComponent(taskId)}`;
}

function clearProgress(taskId: string) {
  const next = { ...progressByTaskId.value };
  delete next[taskId];
  progressByTaskId.value = next;
}

function closeTaskSocket(taskId: string) {
  const ws = sockets.get(taskId);
  if (!ws) return;
  ws.onopen = null;
  ws.onmessage = null;
  ws.onerror = null;
  ws.onclose = null;
  if (
    ws.readyState === WebSocket.OPEN ||
    ws.readyState === WebSocket.CONNECTING
  ) {
    ws.close();
  }
  sockets.delete(taskId);
}

function upsertTaskStatus(taskID: string, status: string) {
  tasks.value = tasks.value.map((task) =>
    task.id === taskID ? { ...task, status } : task,
  );
}

function connectTaskProgress(taskId: string) {
  if (sockets.has(taskId)) return;

  const ws = new WebSocket(wsUrl(taskId));
  sockets.set(taskId, ws);

  ws.onmessage = (event) => {
    try {
      const progress = JSON.parse(event.data) as ProgressBar;

      progressByTaskId.value = {
        ...progressByTaskId.value,
        [progress.task_id]: progress,
      };

      upsertTaskStatus(progress.task_id, progress.status);

      if (["done", "error", "failed", "cancelled"].includes(progress.status)) {
        closeTaskSocket(progress.task_id);
        setTimeout(() => {
          clearProgress(progress.task_id);
          void fetchTasks();
        }, 600);
      }
    } catch {
      // ignore malformed payloads
    }
  };

  ws.onclose = () => {
    sockets.delete(taskId);
  };

  ws.onerror = () => {
    closeTaskSocket(taskId);
  };
}

watch(tasks, (nextTasks) => {
  const activeIDs = new Set(
    nextTasks
      .filter(
        (task) => task.status === "pending" || task.status === "processing",
      )
      .map((task) => task.id),
  );

  for (const taskId of activeIDs) {
    connectTaskProgress(taskId);
  }

  for (const taskId of sockets.keys()) {
    if (!activeIDs.has(taskId) && !nextTasks.find((t) => t.id === taskId)) {
      closeTaskSocket(taskId);
      clearProgress(taskId);
    }
  }
});

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

async function exportToExcel(taskDate: string) {
  try {
    const res = await fetch(
      `/api/export?date_start=${taskDate}&date_end=${taskDate}`,
    );
    if (!res.ok) throw new Error(`export failed: ${res.status}`);
    const blob = await res.blob();
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `matches-${taskDate}.xlsx`;
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
    document.body.removeChild(a);
  } catch (err: any) {
    error.value = err.message || String(err);
  }
}

onMounted(fetchTasks);

onBeforeUnmount(() => {
  for (const taskId of sockets.keys()) {
    closeTaskSocket(taskId);
  }
});
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
      />
    </div>
  </main>
</template>
