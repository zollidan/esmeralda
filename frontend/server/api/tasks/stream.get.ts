import { proxyRequest } from 'h3'

export default defineEventHandler(async (event) => {
  const { apiBase } = useRuntimeConfig()
  const { user } = await requireUserSession(event)

  event.node.req.headers['authorization'] = `Bearer ${user.accessToken}`

  return proxyRequest(event, `${apiBase}/api/tasks/stream`)
})
