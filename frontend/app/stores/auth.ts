import { defineStore } from 'pinia'

interface User {
  id: string
  name: string
  email: string
  role: string
}

export const useAuthStore = defineStore('auth', () => {
  const tokenCookie = useCookie<string | null>('token', { maxAge: 60 * 60 * 24 })
  const userState = ref<User | null>(null)

  const token = computed(() => tokenCookie.value)
  const user = computed(() => userState.value)
  const isAuthenticated = computed(() => !!token.value)

  async function login(email: string, password: string) {
    const config = useRuntimeConfig()
    const baseURL = config.public.apiBase || '/api/v1'

    try {
      const data = await $fetch<{ token: string; user: User }>('/auth/login', {
        method: 'POST',
        baseURL,
        body: { email, password },
      })

      tokenCookie.value = data.token
      userState.value = data.user
      return data
    } catch (error) {
      throw error
    }
  }

  async function register(name: string, email: string, password: string) {
    const config = useRuntimeConfig()
    const baseURL = config.public.apiBase || '/api/v1'

    try {
      return await $fetch('/auth/register', {
        method: 'POST',
        baseURL,
        body: { name, email, password },
      })
    } catch (error) {
      throw error
    }
  }

  function logout() {
    tokenCookie.value = null
    userState.value = null
  }

  // Restore user info if token is present
  async function fetchCurrentUser() {
    if (!token.value) return

    // Since we don't have a specific /auth/me, we can extract details if needed,
    // or when we make a request and get a 401, it will logout.
  }

  return {
    token,
    user,
    isAuthenticated,
    login,
    register,
    logout,
    fetchCurrentUser,
  }
})
