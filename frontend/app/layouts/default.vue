<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
import { LayoutDashboardIcon, LineChartIcon, ArrowRightLeftIcon, ActivityIcon, LogOutIcon } from 'lucide-vue-next'

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()

const navItems = [
  { name: 'Overview', to: '/', icon: LayoutDashboardIcon },
  { name: 'Portfolio', to: '/portfolio', icon: LineChartIcon },
  { name: 'Cashflows', to: '/cashflows', icon: ArrowRightLeftIcon },
  { name: 'Servers Monitoring', to: '/monitoring', icon: ActivityIcon }
]

const currentRouteName = computed(() => {
  if (route.path === '/') return 'Overview'
  if (route.path.startsWith('/portfolio')) return 'Portfolio'
  if (route.path.startsWith('/cashflows')) return 'Cashflows'
  if (route.path.startsWith('/monitoring')) return 'Infrastructure Monitor'
  return ''
})

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>

<template>
  <div class="flex h-screen overflow-hidden bg-background text-foreground font-sans">
    <!-- Sidebar -->
    <aside class="hidden md:flex md:w-64 md:flex-col border-r border-border bg-card">
      <div class="flex h-16 items-center px-6 border-b border-border">
        <h1 class="text-xl font-bold tracking-tight text-primary flex items-center gap-2">
          <ActivityIcon class="h-6 w-6 text-primary" />
          PortInves
        </h1>
      </div>
      <div class="flex-1 overflow-y-auto py-4 px-4 space-y-1">
        <NuxtLink
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          class="flex items-center gap-3 px-3 py-2 text-sm font-medium rounded-md transition-colors"
          :class="[route.path === item.to || (item.to !== '/' && route.path.startsWith(item.to)) ? 'bg-primary/10 text-primary' : 'text-muted-foreground hover:bg-muted hover:text-foreground']"
        >
          <component :is="item.icon" class="h-5 w-5" />
          {{ item.name }}
        </NuxtLink>
      </div>
      <div class="p-4 border-t border-border flex items-center gap-3 bg-muted/10">
        <Avatar>
          <AvatarFallback class="bg-primary text-primary-foreground">{{ (authStore.user?.name || 'G').charAt(0).toUpperCase() }}</AvatarFallback>
        </Avatar>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold truncate">{{ authStore.user?.name || 'Guest' }}</p>
          <p class="text-xs text-muted-foreground truncate">{{ authStore.user?.email }}</p>
        </div>
        <Button variant="ghost" size="icon" @click="handleLogout">
          <LogOutIcon class="h-4 w-4" />
        </Button>
      </div>
    </aside>

    <!-- Main Content Area -->
    <div class="flex flex-1 flex-col overflow-hidden">
      <!-- Header -->
      <header class="flex h-16 items-center justify-between px-6 border-b border-border bg-card">
        <div class="flex items-center gap-4">
          <h2 class="text-lg font-semibold">{{ currentRouteName }}</h2>
        </div>
        <div class="flex items-center gap-4">
          <NotificationBell />
        </div>
      </header>

      <!-- Content -->
      <main class="flex-1 overflow-y-auto p-6 bg-muted/10">
        <slot />
      </main>
    </div>
  </div>
</template>
