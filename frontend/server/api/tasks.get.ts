export default defineEventHandler(async () => {
  const { apiBase } = useRuntimeConfig()

  const data = await $fetch(`${apiBase}/api/tasks/`, {
    method: 'GET',
  })

  return data
})
