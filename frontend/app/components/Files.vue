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
</script>

<template>
  <div class="p-4 max-w-4xl mx-auto">
    <div v-if="loading" class="text-center text-gray-500">Загрузка...</div>
    <div v-else-if="error" class="text-center text-red-600">{{ error }}</div>

    <div v-else class="mt-4 overflow-x-auto rounded-lg border border-gray-200">
      <table class="min-w-full divide-y divide-gray-200 bg-white">
        <thead class="bg-gray-50">
          <tr>
            <th
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
            >
              ID
            </th>
            <th
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
            >
              Имя файла
            </th>
            <th
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
            >
              Размер
            </th>
            <th
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
            >
              Создан
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="file in files" :key="file.ID" class="hover:bg-gray-50">
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
              {{ file.ID }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
              {{ file.filename }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
              {{ toMB(file.file_size) }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              {{ toDate(file.CreatedAt) }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
