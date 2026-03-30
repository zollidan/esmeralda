export const backendFetch = async <T>(
  event: any,
  path: string,
  options: Record<string, any> = {},
): Promise<T> => {
  const { apiBase } = useRuntimeConfig()
  const { user } = await requireUserSession(event)

  return $fetch<T>(`${apiBase}${path}`, {
    ...options,
    headers: {
      ...options.headers,
      Authorization: `Bearer ${user.accessToken}`,
    },
  })
}
