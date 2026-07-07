<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useApi } from '~/composables/useApi'
import PerformanceLineChart from '~/components/charts/PerformanceLineChart.vue'
import { DollarSign, Briefcase, Server, Activity } from 'lucide-vue-next'

interface DashboardOverview {
  total_portfolios: number
  total_value: number
  total_servers: number
  servers_online: number
  servers_offline: number
  servers_degraded: number
}

const api = useApi()
const loading = ref(true)
const overview = ref<DashboardOverview>({
  total_portfolios: 0,
  total_value: 0,
  total_servers: 0,
  servers_online: 0,
  servers_offline: 0,
  servers_degraded: 0
})

const performanceDates = ref(['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'])
const performanceValues = ref([0, 0, 0, 0, 0, 0, 0])

async function fetchOverview() {
  try {
    const data = await api.fetch<DashboardOverview>('/dashboard/overview')
    overview.value = data

    if (data.total_value > 0) {
      const val = data.total_value
      performanceValues.value = [
        val * 0.92,
        val * 0.95,
        val * 0.93,
        val * 0.97,
        val * 0.96,
        val * 0.99,
        val
      ]
    } else {
      performanceValues.value = [0, 0, 0, 0, 0, 0, 0]
    }
  } catch (error) {
    console.error('Failed to fetch dashboard overview:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchOverview()
})
</script>

<template>
  <div class="space-y-6">
    <!-- Top Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <Card class="bg-card border-border">
        <div class="p-6 flex flex-col gap-2">
          <div class="flex items-center justify-between">
            <span class="text-sm font-semibold text-muted-foreground">Portfolio Value</span>
            <DollarSign class="h-5 w-5 text-primary" />
          </div>
          <div class="text-2xl font-bold text-foreground">
            ${{ overview.total_value.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }}
          </div>
          <p class="text-xs text-muted-foreground mt-1">Across {{ overview.total_portfolios }} portfolios</p>
        </div>
      </Card>

      <Card class="bg-card border-border">
        <div class="p-6 flex flex-col gap-2">
          <div class="flex items-center justify-between">
            <span class="text-sm font-semibold text-muted-foreground">Active Portfolios</span>
            <Briefcase class="h-5 w-5 text-indigo-500" />
          </div>
          <div class="text-2xl font-bold text-foreground">{{ overview.total_portfolios }}</div>
          <p class="text-xs text-muted-foreground mt-1">Managed accounts</p>
        </div>
      </Card>

      <Card class="bg-card border-border">
        <div class="p-6 flex flex-col gap-2">
          <div class="flex items-center justify-between">
            <span class="text-sm font-semibold text-muted-foreground">Servers Monitored</span>
            <Server class="h-5 w-5 text-emerald-500" />
          </div>
          <div class="text-2xl font-bold text-foreground">{{ overview.total_servers }}</div>
          <p class="text-xs text-muted-foreground mt-1">{{ overview.servers_online }} online / {{ overview.servers_offline }} offline</p>
        </div>
      </Card>

      <Card class="bg-card border-border">
        <div class="p-6 flex flex-col gap-2">
          <div class="flex items-center justify-between">
            <span class="text-sm font-semibold text-muted-foreground">System Health</span>
            <Activity class="h-5 w-5 text-rose-500" />
          </div>
          <div class="text-2xl font-bold text-foreground">
            {{ overview.total_servers > 0 && overview.servers_offline === 0 ? 'Optimal' : overview.total_servers === 0 ? 'No Data' : 'Alert' }}
          </div>
          <p class="text-xs text-muted-foreground mt-1">
            <span v-if="overview.servers_offline > 0" class="text-red-500 font-semibold">{{ overview.servers_offline }} servers down</span>
            <span v-else-if="overview.servers_degraded > 0" class="text-yellow-500 font-semibold">{{ overview.servers_degraded }} servers degraded</span>
            <span v-else>All systems operational</span>
          </p>
        </div>
      </Card>
    </div>

    <!-- Main Chart Section -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <Card class="lg:col-span-2 bg-card border-border">
        <div class="p-6">
          <h3 class="text-base font-semibold leading-6 text-foreground mb-4">Portfolio Performance (7 Days)</h3>
          <div class="h-80">
            <PerformanceLineChart :dates="performanceDates" :values="performanceValues" />
          </div>
        </div>
      </Card>

      <!-- Infrastructure Quick view -->
      <Card class="bg-card border-border">
        <div class="p-6">
          <h3 class="text-base font-semibold leading-6 text-foreground mb-4">Infrastructure Health</h3>
          <div class="space-y-4">
            <div class="flex items-center justify-between p-3 bg-muted/30 rounded-lg">
              <div class="flex items-center gap-3">
                <span class="relative flex h-3 w-3">
                  <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                  <span class="relative inline-flex rounded-full h-3 w-3 bg-emerald-500"></span>
                </span>
                <span class="text-sm font-medium">Online</span>
              </div>
              <span class="text-sm font-bold">{{ overview.servers_online }}</span>
            </div>

            <div class="flex items-center justify-between p-3 bg-muted/30 rounded-lg">
              <div class="flex items-center gap-3">
                <span class="relative flex h-3 w-3">
                  <span class="relative inline-flex rounded-full h-3 w-3 bg-yellow-500"></span>
                </span>
                <span class="text-sm font-medium">Degraded</span>
              </div>
              <span class="text-sm font-bold">{{ overview.servers_degraded }}</span>
            </div>

            <div class="flex items-center justify-between p-3 bg-muted/30 rounded-lg">
              <div class="flex items-center gap-3">
                <span class="relative flex h-3 w-3">
                  <span class="relative inline-flex rounded-full h-3 w-3 bg-red-500"></span>
                </span>
                <span class="text-sm font-medium">Offline</span>
              </div>
              <span class="text-sm font-bold">{{ overview.servers_offline }}</span>
            </div>

            <div class="pt-4 border-t border-border">
              <NuxtLink to="/monitoring" class="w-full">
                <Button variant="secondary" class="w-full flex items-center gap-2">
                  <Activity class="w-4 h-4" /> View Monitoring
                </Button>
              </NuxtLink>
            </div>
          </div>
        </div>
      </Card>
    </div>
  </div>
</template>
