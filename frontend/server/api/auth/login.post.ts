export default defineEventHandler(async (event) => {
  const { apiBase } = useRuntimeConfig()
  const body = await readBody(event)

  const data = await $fetch<{ access_token: string, refresh_token: string }>(
    `${apiBase}/api/auth/login`,
    { method: 'POST', body },
  )

  await setUserSession(event, {
    user: {
      accessToken: data.access_token,
      refreshToken: data.refresh_token,
    },
  })

  return { ok: true }
})
