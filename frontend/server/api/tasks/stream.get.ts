import { proxyRequest } from 'h3'

export default defineEventHandler((event) => {
  const { apiBase } = useRuntimeConfig()
  return proxyRequest(event, `${apiBase}/api/tasks/stream`)
})
