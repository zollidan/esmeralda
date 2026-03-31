export default defineEventHandler(async (event) => {
  const taskId = event.context.params?.task_id
  if (!taskId) {
    throw createError({ statusCode: 400, statusMessage: 'task_id required' })
  }

  return backendFetch(event, `/api/tasks/${taskId}`, {
    method: 'DELETE',
  })
})
