export default defineEventHandler(async (event) => {
  const { apiBase } = useRuntimeConfig()
  const body = await readBody(event)

  const data = await $fetch(`${apiBase}/api/tasks/`, {
    method: 'POST',
    body: body,
  })

  return data
})
