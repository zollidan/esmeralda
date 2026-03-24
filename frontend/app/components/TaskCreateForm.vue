<script setup lang="ts">
const date = defineModel<string>("date", { required: true });

defineProps<{
  visible: boolean;
}>();

const emit = defineEmits<{
  submit: [];
}>();
</script>

<template>
  <section v-if="visible" class="bg-white rounded-lg shadow-md p-6 mb-8">
    <h2 class="text-xl font-semibold text-slate-700 mb-4">
      Создать новую задачу
    </h2>

    <div class="flex gap-4 items-end">
      <div class="flex-1">
        <label for="date" class="block text-sm font-medium text-slate-700 mb-2">
          выбор даты
        </label>
        <input
          id="date"
          v-model="date"
          type="date"
          class="w-full px-4 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition"
        />
      </div>

      <button
        :disabled="!date"
        class="px-6 py-2 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 disabled:bg-slate-300 disabled:cursor-not-allowed transition"
        @click="emit('submit')"
      >
        Создать
      </button>
    </div>
  </section>
</template>
