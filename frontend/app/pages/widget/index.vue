<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '~/stores/auth'
import { LockIcon, Loader2Icon, WalletIcon } from 'lucide-vue-next'

definePageMeta({
  layout: 'widget'
})

const authStore = useAuthStore()
const config = useRuntimeConfig()

// State
const overview = ref<any>(null)
const servers = ref<any[]>([])
const loading = ref(true)

// Polling interval
let pollTimer: any = null

onMounted(async () => {
  if (!authStore.isAuthenticated) {
    // Attempt automatic login for widget if token exists
    await authStore.checkAuth()
  }
  await fetchWidgetData()
  
  // Poll every 30 seconds
  pollTimer = setInterval(fetchWidgetData, 30000)
})

onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer)
})

async function fetchWidgetData() {
  if (!authStore.isAuthenticated) return
  
  try {
    const [overviewRes, serversRes] = await Promise.all([
      $fetch(`${config.public.apiBase}/api/v1/dashboard/overview`, {
        headers: { Authorization: `Bearer ${authStore.token}` }
      }),
      $fetch(`${config.public.apiBase}/api/v1/servers`, {
        headers: { Authorization: `Bearer ${authStore.token}` }
      })
    ])
    
    overview.value = overviewRes
    servers.value = (serversRes as any[]) || []
  } catch (error) {
    console.error('Failed to fetch widget data', error)
  } finally {
    loading.value = false
  }
}

function getStatusColor(status: string) {
  if (status === 'online') return 'bg-green-500'
  if (status === 'degraded') return 'bg-yellow-500'
  return 'bg-red-500'
}
</script>

<template>
  <div class="p-4 text-slate-200">
    <!-- Unauthenticated State -->
    <div v-if="!authStore.isAuthenticated" class="flex flex-col items-center justify-center h-48 text-center">
      <LockIcon class="w-8 h-8 text-slate-500 mb-2" />
      <p class="text-sm font-medium">Please login via the main dashboard first.</p>
    </div>

    <div v-else-if="loading" class="flex justify-center p-8">
      <Loader2Icon class="w-6 h-6 animate-spin text-primary" />
    </div>

    <!-- Authenticated Widget Content -->
    <div v-else class="space-y-6">
      
      <!-- Net Worth Section -->
      <div class="bg-slate-800/60 rounded-lg p-4 border border-slate-700 shadow-inner">
        <h3 class="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-1">Total Net Worth</h3>
        <p class="text-2xl font-bold text-white flex items-center">
          <WalletIcon class="w-5 h-5 mr-2 text-primary" />
          {{ overview?.total_value.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00' }} <span class="text-sm font-normal text-slate-400 ml-1">THB</span>
        </p>
        <div class="mt-2 text-xs flex justify-between text-slate-300">
          <span>Portfolios: {{ overview?.total_portfolios || 0 }}</span>
        </div>
      </div>

      <!-- Infrastructure Status Section -->
      <div>
        <div class="flex items-center justify-between mb-3">
          <h3 class="text-xs font-semibold text-slate-400 uppercase tracking-wider">Infrastructure</h3>
          <Badge variant="secondary" class="text-[10px]">
            {{ servers.length }} Nodes
          </Badge>
        </div>
        
        <div class="space-y-2">
          <div v-if="servers.length === 0" class="text-xs text-slate-500 italic">No servers tracked.</div>
          
          <div 
            v-for="server in servers" 
            :key="server.id"
            class="flex items-center justify-between bg-slate-800/40 p-2.5 rounded-md border border-slate-700/50"
          >
            <div class="flex items-center gap-2 overflow-hidden">
              <span class="w-2 h-2 rounded-full flex-shrink-0" :class="getStatusColor(server.status)"></span>
              <span class="text-sm font-medium truncate">{{ server.name }}</span>
            </div>
            <span class="text-[10px] text-slate-400 uppercase bg-slate-800 px-1.5 py-0.5 rounded border border-slate-600">{{ server.type }}</span>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>
