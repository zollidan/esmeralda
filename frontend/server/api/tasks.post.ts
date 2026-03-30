export default defineEventHandler(async (event) => {
  const body = await readBody(event)

  return backendFetch(event, '/api/tasks/', {
    method: 'POST',
    body,
  })
})
