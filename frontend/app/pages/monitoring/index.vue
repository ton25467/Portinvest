<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useApi } from '~/composables/useApi'
import { PlusIcon, ServerIcon, CheckCircleIcon, XCircleIcon } from 'lucide-vue-next'

interface Server {
  id: string
  name: string
  host: string
  port: number
  type: string
  status: string
  last_checked_at: string | null
}

interface ServiceCheckLog {
  server_id: string
  service_check_id: string
  status: string
  status_code: number
  response_time_ms: number
  checked_at: string
  server_status: string
}

const api = useApi()
const servers = ref<Server[]>([])
const loading = ref(true)

const showCreateModal = ref(false)
const newName = ref('')
const newHost = ref('')
const newPort = ref(80)
const newType = ref('web')
const createLoading = ref(false)

const logsList = ref<Array<{ time: string; msg: string; status: string }>>([])

async function fetchServers() {
  try {
    const data = await api.fetch<Server[]>('/servers')
    servers.value = data
  } catch (error) {
    console.error('Failed to fetch servers:', error)
  } finally {
    loading.value = false
  }
}

async function handleCreateServer() {
  if (!newName.value || !newHost.value) return

  createLoading.value = true
  try {
    const srv = await api.fetch<Server>('/servers', {
      method: 'POST',
      body: {
        name: newName.value,
        host: newHost.value,
        port: newPort.value,
        type: newType.value
      }
    })

    // Automatically create a default HTTP service check for this server
    await api.fetch(`/servers/${srv.id}/checks`, {
      method: 'POST',
      body: {
        name: 'HTTP Port Check',
        endpoint: `http://${newHost.value}:${newPort.value}`,
        method: 'GET',
        expected_status: 200,
        interval_seconds: 15
      }
    })

    showCreateModal.value = false
    newName.value = ''
    newHost.value = ''
    newPort.value = 80
    fetchServers()
  } catch (error) {
    console.error('Failed to register server:', error)
  } finally {
    createLoading.value = false
  }
}

// Setup WebSocket for real-time monitoring results
onMounted(() => {
  fetchServers()

  // Connect to Go websocket
  const config = useRuntimeConfig()
  const wsUrl = config.public.wsBase || 'ws://localhost:8080/ws'
  const { status, data } = useWebSocket(wsUrl, {
    autoReconnect: {
      retries: 10,
      delay: 3000
    }
  })

  // Watch incoming messages
  watch(data, (newMsg) => {
    if (!newMsg) return
    try {
      const parsed = JSON.parse(newMsg)
      if (parsed.type === 'service_check_result') {
        const update = parsed as ServiceCheckLog

        // 1. Update server status in list
        const srv = servers.value.find(s => s.id === update.server_id)
        if (srv) {
          srv.status = update.server_status
          srv.last_checked_at = update.checked_at
        }

        // 2. Add to real-time log feed
        const time = new Date(update.checked_at).toLocaleTimeString()
        const logMsg = `Check finished for Server ${srv?.name || update.server_id} in ${update.response_time_ms}ms (Status: ${update.status_code})`
        logsList.value.unshift({
          time,
          msg: logMsg,
          status: update.status
        })

        // Cap logs at 50 entries
        if (logsList.value.length > 50) {
          logsList.value.pop()
        }
      }
    } catch (e) {
      // Ignore ping frames or parsing errors
    }
  })
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h2 class="text-xl font-bold tracking-tight text-foreground">Infrastructure Monitoring</h2>
      <Button @click="showCreateModal = true">
        <PlusIcon class="mr-2 h-4 w-4" /> Register Server
      </Button>
    </div>

    <!-- Register Server Modal -->
    <Dialog :open="showCreateModal" @update:open="showCreateModal = $event">
      <DialogContent class="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Register Monitored Server</DialogTitle>
        </DialogHeader>
        <div class="py-4 space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <Label for="newName">Server Name</Label>
              <Input id="newName" v-model="newName" placeholder="e.g. Production Web" />
            </div>

            <div class="space-y-2">
              <Label for="newType">Type</Label>
              <Select v-model="newType">
                <SelectTrigger>
                  <SelectValue placeholder="Select type" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="web">Web Server</SelectItem>
                  <SelectItem value="db">Database Server</SelectItem>
                  <SelectItem value="api">API Gateway</SelectItem>
                  <SelectItem value="cache">Caching Layer</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div class="grid grid-cols-3 gap-4">
            <div class="col-span-2 space-y-2">
              <Label for="newHost">Host / IP</Label>
              <Input id="newHost" v-model="newHost" placeholder="e.g. 192.168.1.100 or app.com" />
            </div>

            <div class="space-y-2">
              <Label for="newPort">Port</Label>
              <Input id="newPort" v-model="newPort" type="number" placeholder="80" />
            </div>
          </div>

          <div class="flex justify-end gap-3 pt-4 border-t border-border mt-4">
            <Button variant="outline" @click="showCreateModal = false">Cancel</Button>
            <Button :disabled="createLoading" @click="handleCreateServer">
              <span v-if="createLoading" class="mr-2 animate-spin">⟳</span>
              Register
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>

    <!-- Main Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Servers Grid -->
      <div class="lg:col-span-2">
        <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div v-for="i in 4" :key="i" class="h-32 w-full animate-pulse bg-muted rounded-xl" />
        </div>
        <div v-else-if="servers.length === 0" class="text-center py-12 bg-card border border-border rounded-xl">
          <ServerIcon class="mx-auto h-12 w-12 text-muted-foreground" />
          <h3 class="mt-2 text-sm font-semibold text-foreground">No servers monitored</h3>
          <p class="mt-1 text-sm text-muted-foreground">Register servers to begin tracking health and uptime.</p>
          <div class="mt-6">
            <Button @click="showCreateModal = true">
              <PlusIcon class="mr-2 h-4 w-4" /> Register Server
            </Button>
          </div>
        </div>
        <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <Card
            v-for="s in servers"
            :key="s.id"
            class="bg-card border-border hover:shadow-sm transition-shadow"
          >
            <div class="p-6">
              <div class="flex justify-between items-start">
                <div>
                  <h4 class="text-lg font-bold flex items-center gap-2">
                    {{ s.name }}
                    <span class="text-xs text-muted-foreground bg-muted px-2 py-0.5 rounded capitalize">
                      {{ s.type }}
                    </span>
                  </h4>
                  <p class="text-sm text-muted-foreground mt-1 font-mono">{{ s.host }}:{{ s.port }}</p>
                  <p class="text-xs text-muted-foreground mt-2">
                    Last Checked: {{ s.last_checked_at ? new Date(s.last_checked_at).toLocaleTimeString() : 'Never' }}
                  </p>
                </div>

                <!-- Pulse health status icon -->
                <span class="relative flex h-4 w-4">
                  <span
                    class="animate-ping absolute inline-flex h-full w-full rounded-full opacity-75"
                    :class="[s.status === 'online' ? 'bg-emerald-400' : s.status === 'degraded' ? 'bg-yellow-400' : 'bg-red-400']"
                  ></span>
                  <span
                    class="relative inline-flex rounded-full h-4 w-4"
                    :class="[s.status === 'online' ? 'bg-emerald-500' : s.status === 'degraded' ? 'bg-yellow-500' : 'bg-red-500']"
                  ></span>
                </span>
              </div>
            </div>
          </Card>
        </div>
      </div>

      <!-- Real-Time Activity Log Feed -->
      <Card class="bg-card border-border flex flex-col h-[400px]">
        <div class="p-4 border-b border-border flex items-center justify-between">
          <h3 class="text-base font-semibold leading-6 text-foreground">Live Activity Feed</h3>
          <span class="flex h-2 w-2 relative">
            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
            <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
          </span>
        </div>

        <div class="flex-1 overflow-y-auto space-y-3 p-4 h-[300px]">
          <div v-if="logsList.length === 0" class="text-center py-12 text-sm text-muted-foreground">
            Waiting for live check reports...
          </div>
          <div
            v-for="(log, idx) in logsList"
            :key="idx"
            class="flex items-start gap-3 p-2 rounded bg-muted/20 border border-border text-xs"
          >
            <CheckCircleIcon v-if="log.status === 'up'" class="h-4 w-4 mt-0.5 shrink-0 text-emerald-500" />
            <XCircleIcon v-else class="h-4 w-4 mt-0.5 shrink-0 text-rose-500" />
            <div class="flex-1 min-w-0">
              <p class="font-medium text-foreground leading-normal">{{ log.msg }}</p>
              <span class="text-muted-foreground mt-1 block font-mono text-[10px]">{{ log.time }}</span>
            </div>
          </div>
        </div>
      </Card>
    </div>
  </div>
</template>
