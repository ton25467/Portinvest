import { useAuthStore } from '~/stores/auth'

export function useApi() {
  const config = useRuntimeConfig()
  const authStore = useAuthStore()
  const router = useRouter()

  const apiFetch = async <T = any>(path: string, options: any = {}): Promise<T> => {
    // Determine the base URL
    const baseURL = config.public.apiBase || '/api/v1'

    // Add authorization header if token is present
    const headers = {
      ...options.headers,
    }

    if (authStore.token) {
      headers['Authorization'] = `Bearer ${authStore.token}`
    }

    try {
      return await $fetch<T>(path, {
        baseURL,
        ...options,
        headers,
      })
    } catch (error: any) {
      // If unauthorized, clear auth and redirect to login
      if (error.status === 401) {
        authStore.logout()
        router.push('/login')
      }
      throw error
    }
  }

  return {
    fetch: apiFetch,
  }
}
