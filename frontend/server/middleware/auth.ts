export default defineEventHandler(async (event) => {
  const url = getRequestURL(event)

  // Защищаем только /api/tasks и /api/export, но не /api/auth/*
  if (!url.pathname.startsWith('/api/') || url.pathname.startsWith('/api/auth/')) {
    return
  }

  const { user } = await getUserSession(event)
  if (!user?.accessToken) {
    throw createError({ statusCode: 401, message: 'Unauthorized' })
  }
})
