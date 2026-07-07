import { useAuthStore } from '~/stores/auth'

export default defineNuxtRouteMiddleware((to) => {
  const authStore = useAuthStore()

  // Skip middleware for login and register pages
  if (to.path === '/login' || to.path === '/register') {
    if (authStore.isAuthenticated) {
      return navigateTo('/')
    }
    return
  }

  // Redirect to login if not authenticated
  if (!authStore.isAuthenticated) {
    return navigateTo('/login')
  }
})
