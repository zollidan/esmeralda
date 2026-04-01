<script setup lang="ts">
const date = defineModel<string>('date', { required: true })

defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  submit: []
}>()
</script>

<template>
  <section
    v-if="visible"
    class="mb-8 p-6 rounded-lg shadow-md transition-colors bg-white dark:bg-slate-800"
  >
    <h2 class="text-xl font-semibold mb-4 text-slate-700 dark:text-white">
      Создать новую задачу
    </h2>

    <div class="flex gap-4 items-end">
      <div class="flex-1">
        <label
          for="date"
          class="block text-sm font-medium mb-2 text-slate-700 dark:text-slate-300"
        >
          выбор даты
        </label>

        <input
          id="date"
          v-model="date"
          type="date"
          class="w-full px-4 py-2 rounded-lg border transition border-slate-300 dark:border-slate-600 bg-white dark:bg-slate-700 text-slate-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
        >
      </div>

      <button
        :disabled="!date"
        class="px-6 py-2 rounded-lg font-medium text-white transition bg-blue-600 hover:bg-blue-700 disabled:bg-slate-300 dark:disabled:bg-slate-600 disabled:cursor-not-allowed"
        @click="emit('submit')"
      >
        Создать
      </button>
    </div>
  </section>
</template>
