export default defineEventHandler(async (event) => {
  const { apiBase } = useRuntimeConfig()

  const data = await $fetch(
    `${apiBase}/api/tasks/${event.context.params.task_id}`,
    {
      method: 'DELETE',
    },
  )

  return data
})
