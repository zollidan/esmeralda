// server/api/auth/refresh.post.ts
export default defineEventHandler(async (event) => {
  const { apiBase } = useRuntimeConfig()
  const { user } = await getUserSession(event)

  if (!user?.refreshToken) {
    throw createError({ statusCode: 401, message: 'No refresh token' })
  }

  try {
    const data = await $fetch<{ access_token: string, refresh_token: string }>(
      `${apiBase}/api/auth/refresh`,
      {
        method: 'POST',
        body: { refresh_token: user.refreshToken },
      },
    )

    // Ротация — сохраняем оба новых токена
    await setUserSession(event, {
      user: {
        accessToken: data.access_token,
        refreshToken: data.refresh_token,
      },
    })

    return { ok: true }
  }
  catch {
    await clearUserSession(event)
    throw createError({ statusCode: 401, message: 'Refresh failed' })
  }
})
