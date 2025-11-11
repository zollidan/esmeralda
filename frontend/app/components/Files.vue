<script setup>
import { ref } from "vue";

const files = ref([]);
const error = ref(null);
const loading = ref(true);

try {
  const { data, error: fetchError } = await useFetch("/api/files", {
    method: "GET",
  });

  if (fetchError.value) {
    error.value = fetchError.value.message || "Ошибка при загрузке данных";
  } else if (
    !data.value ||
    !Array.isArray(data.value) ||
    data.value.length === 0
  ) {
    error.value = "Нет файлов для отображения";
  } else {
    files.value = data.value;
  }
} catch (err) {
  error.value = err.message || "Ошибка при выполнении запроса";
} finally {
  loading.value = false;
}

const toMB = (bytes) => {
  return (bytes / (1024 * 1024)).toFixed(2) + " MB";
};

const toDate = (timestamp) => {
  const date = new Date(timestamp);
  return date.toLocaleDateString("ru-RU", {
    year: "numeric",
    month: "numeric",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
};

const columns = [
  { key: "ID", label: "ID" },
  { key: "filename", label: "Имя файла" },
  {
    key: "file_size",
    label: "Размер",
    render: (val) => toMB(val),
  },
  {
    key: "CreatedAt",
    label: "Создан",
    render: (val) => toDate(val),
  },
];
</script>

<template>
  <div class="p-4 max-w-4xl mx-auto">
    <div v-if="loading" class="text-center text-gray-500">Загрузка...</div>
    <div v-else-if="error" class="text-center text-red-600">{{ error }}</div>

    <div v-else>
      <Table :columns="columns" :items="files" />
    </div>
  </div>
</template>
