export default defineNuxtPlugin(() => {
  const { clear } = useUserSession()
  let refreshPromise: Promise<void> | null = null

  const api = $fetch.create({
    onResponseError: async ({ response, request, options }) => {
      if (response.status !== 401) return

      if (!refreshPromise) {
        refreshPromise = $fetch('/api/auth/refresh', { method: 'POST' })
          .then(() => {})
          .finally(() => { refreshPromise = null })
      }

      try {
        await refreshPromise
        return $fetch(request, options)
      }
      catch {
        await clear()
        navigateTo('/login')
      }
    },
  })

  return { provide: { api } }
})
