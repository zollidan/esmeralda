import type { Task } from "~/types/task";

function errorMessage(err: unknown): string {
  if (err instanceof Error) return err.message;
  return String(err);
}

export function useTasks() {
  const date = ref("");
  const tasks = ref<Task[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);

  const fetchTasks = async () => {
    loading.value = true;
    error.value = null;
    try {
      const res = await fetch("/api/tasks");
      if (!res.ok) throw new Error(`status ${res.status}`);
      tasks.value = await res.json();
    } catch (err: unknown) {
      error.value = errorMessage(err);
    } finally {
      loading.value = false;
    }
  };

  const createTask = async () => {
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
    } catch (err: unknown) {
      error.value = errorMessage(err);
    }
  };

  const deleteTask = async (taskId: string) => {
    try {
      const res = await fetch(`/api/tasks/${taskId}`, { method: "DELETE" });
      if (!res.ok) throw new Error(`delete failed: ${res.status}`);
      await fetchTasks();
    } catch (err: unknown) {
      error.value = errorMessage(err);
    }
  };

  const exportToExcel = async (taskDate: string) => {
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
    } catch (err: unknown) {
      error.value = errorMessage(err);
    }
  };

  return { date, tasks, loading, error, fetchTasks, createTask, deleteTask, exportToExcel };
}
