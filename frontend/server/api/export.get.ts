export default defineEventHandler(async (event) => {
  const { apiBase } = useRuntimeConfig()
  const q = getQuery(event)

  const response = await fetch(
    `${apiBase}/api/export?${new URLSearchParams(q as Record<string, string>)}`,
  )

  if (!response.ok) {
    throw createError({ statusCode: response.status, statusMessage: 'Export failed' })
  }

  const buffer = await response.arrayBuffer()

  setResponseHeader(event, 'Content-Type', 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet')
  setResponseHeader(event, 'Content-Disposition', response.headers.get('Content-Disposition') ?? 'attachment; filename=export.xlsx')

  return new Uint8Array(buffer)
})
