export const useAuth = () => {
  const { loggedIn, user, clear, fetch: refreshSession } = useUserSession()

  const login = async (credentials: { username: string, password: string }) => {
    await $fetch('/api/auth/login', {
      method: 'POST',
      body: credentials,
    })
    await refreshSession()
  }

  const logout = async () => {
    await $fetch('/api/auth/logout', { method: 'POST' })
    await clear()
    navigateTo('/login')
  }

  return { loggedIn, user, login, logout }
}
