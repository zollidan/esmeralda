<script setup lang="ts">
import { useSportRuStatus } from '~/composables/useSportRuStatus'

const { status, loading, start, stop } = useSportRuStatus()

onMounted(() => start())
onUnmounted(() => stop())
</script>

<template>
  <div class="flex items-center gap-2 text-sm">
    <template v-if="loading && !status">
      <span class="w-2 h-2 rounded-full bg-gray-400" />
      <span class="text-gray-400">Checking...</span>
    </template>
    <template v-else>
      <span
        class="w-2 h-2 rounded-full"
        :class="
          status?.status === 'online'
            ? 'bg-green-500 animate-pulse'
            : 'bg-red-500'
        "
      />
      <span
        :class="status?.status === 'online' ? 'text-green-600' : 'text-red-500'"
      >
        SportRU: {{ status?.status === "online" ? "Online" : "Offline" }}
      </span>
      <span
        v-if="status?.latency_ms"
        class="text-gray-400"
      >
        {{ status.latency_ms }}ms
      </span>
    </template>
  </div>
</template>
