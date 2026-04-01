import type { ProgressBar, Task } from '~/types/task'

function errorMessage(err: unknown): string {
  if (err instanceof Error) return err.message
  return String(err)
}

export function useTasks() {
  const { $api } = useNuxtApp()

  const date = ref('')
  const tasks = ref<Task[]>([])
  const loading = ref(true)
  const error = ref<string | null>(null)
  const progressByTaskId = ref<Record<string, ProgressBar>>({})

  const createTask = async () => {
    if (!date.value) return
    error.value = null
    try {
      await $api('/api/tasks', {
        // 👈
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ date: date.value }),
      })
      date.value = ''
    }
    catch (err: unknown) {
      error.value = errorMessage(err)
    }
  }

  const deleteTask = async (taskId: string) => {
    try {
      await $api(`/api/tasks/${taskId}`, { method: 'DELETE' })
      tasks.value = tasks.value.filter(t => t.id !== taskId)
    }
    catch (err: unknown) {
      error.value = errorMessage(err)
    }
  }

  const exportToExcel = async (taskDate: string, format: 'xlsx' | 'csv' = 'xlsx') => {
    try {
      const res = await fetch(
        `/api/export?date_start=${taskDate}&date_end=${taskDate}&format=${format}`,
      )
      if (!res.ok) throw new Error(`export failed: ${res.status}`)
      const blob = await res.blob()
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `matches-${taskDate}.${format}`
      document.body.appendChild(a)
      a.click()
      window.URL.revokeObjectURL(url)
      document.body.removeChild(a)
    }
    catch (err: unknown) {
      error.value = errorMessage(err)
    }
  }

  const initStream = () => {
    const { $api } = useNuxtApp()
    const { clear } = useUserSession()

    const abortController = new AbortController()

    const connect = async () => {
      try {
        const response = await fetch('/api/tasks/stream', {
          signal: abortController.signal,
        })

        if (response.status === 401) {
          try {
            await $api('/api/auth/refresh', { method: 'POST' })
            connect()
          }
          catch {
            await clear()
            navigateTo('/login')
          }
          return
        }

        const reader = response.body!.getReader()
        const decoder = new TextDecoder()
        let buffer = ''

        while (true) {
          const { done, value } = await reader.read()
          if (done) break

          buffer += decoder.decode(value, { stream: true })
          const lines = buffer.split('\n\n')
          buffer = lines.pop() ?? ''

          for (const chunk of lines) {
            const dataLine = chunk
              .split('\n')
              .find(l => l.startsWith('data:'))
            if (!dataLine) continue

            const data = JSON.parse(dataLine.slice(5).trim())

            if (Array.isArray(data)) {
              tasks.value = data
              loading.value = false
            }
            else if ('task_id' in data) {
              progressByTaskId.value = {
                ...progressByTaskId.value,
                [data.task_id]: data,
              }
            }
            else if ('id' in data) {
              const idx = tasks.value.findIndex(t => t.id === data.id)
              if (idx !== -1) tasks.value[idx] = data
              else tasks.value.unshift(data)
            }
          }
        }
      }
      catch (err) {
        if ((err as Error).name === 'AbortError') return
        loading.value = false
      }
    }

    connect()
    return {
      close: () => abortController.abort(),
    }
  }

  return {
    date,
    tasks,
    loading,
    error,
    progressByTaskId,
    initStream,
    createTask,
    deleteTask,
    exportToExcel,
  }
}
