export default defineEventHandler((event) => {
  const { apiBase } = useRuntimeConfig()

  const taskId = event.context.params?.task_id
  if (!taskId) {
    throw createError({ statusCode: 400, statusMessage: 'task_id required' })
  }

  return proxyRequest(event, `${apiBase}/api/tasks/${taskId}`)
})
