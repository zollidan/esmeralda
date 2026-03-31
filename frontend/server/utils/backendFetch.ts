import type { H3Event } from 'h3'

export const backendFetch = async <T>(
  event: H3Event,
  path: string,
  options: Record<string, unknown> = {},
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
