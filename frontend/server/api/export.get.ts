export default defineEventHandler(async (event) => {
  const q = getQuery(event)
  const queryString = new URLSearchParams(q as Record<string, string>).toString()

  const buffer = await backendFetch<ArrayBuffer>(event, `/api/export?${queryString}`, {
    responseType: 'arrayBuffer' as any,
  })

  setResponseHeader(event, 'Content-Type', 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet')
  setResponseHeader(event, 'Content-Disposition', `attachment; filename=export.xlsx`)

  return new Uint8Array(buffer)
})
