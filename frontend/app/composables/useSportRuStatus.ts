import type { SportRuAPIStatus } from '~/types/task'

export function useSportRuStatus() {
  const { $api } = useNuxtApp()
  const status = ref<SportRuAPIStatus | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  let intervalId: ReturnType<typeof setInterval> | null = null

  const fetch = async () => {
    loading.value = true
    error.value = null
    try {
      status.value = await $api<SportRuAPIStatus>('/api/sportru/health')
    }
    catch (err: unknown) {
      error.value = err instanceof Error ? err.message : String(err)
      status.value = { status: 'offline', latency_ms: 0 }
    }
    finally {
      loading.value = false
    }
  }

  const start = (intervalMs = 30_000) => {
    fetch()
    intervalId = setInterval(fetch, intervalMs)
  }

  const stop = () => {
    if (intervalId) clearInterval(intervalId)
  }

  return { status, loading, error, fetch, start, stop }
}
