<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '~/stores/auth'
import { toast } from 'vue-sonner'
import { BellIcon } from 'lucide-vue-next'

const authStore = useAuthStore()
const config = useRuntimeConfig()

interface Notification {
  id: string
  title: string
  message: string
  is_read: boolean
  created_at: string
}

const notifications = ref<Notification[]>([])
let ws: WebSocket | null = null

const unreadCount = computed(() => notifications.value.filter(n => !n.is_read).length)

onMounted(async () => {
  if (authStore.isAuthenticated) {
    await fetchNotifications()
    connectWebSocket()
  }
})

onUnmounted(() => {
  if (ws) {
    ws.close()
  }
})

async function fetchNotifications() {
  try {
    const res = await $fetch<Notification[]>(`${config.public.apiBase}/api/v1/notifications?unread_only=false`, {
      headers: {
        Authorization: `Bearer ${authStore.token}`
      }
    })
    if (res) {
      notifications.value = res
    }
  } catch (error) {
    console.error('Failed to fetch notifications:', error)
  }
}

function connectWebSocket() {
  if (!authStore.token) return
  
  // Construct WS URL from API Base
  const apiBase = config.public.apiBase as string || 'http://localhost:8080'
  const wsUrl = apiBase.replace(/^http/, 'ws') + `/ws?token=${authStore.token}`
  
  ws = new WebSocket(wsUrl)
  
  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data)
      if (data.type === 'notification' && data.notification) {
        notifications.value.unshift(data.notification)
        toast.info(data.notification.title, {
          description: data.notification.message,
        })
      }
    } catch (e) {
      // Ignore parse errors or ping messages
    }
  }

  ws.onclose = () => {
    // Retry connection after 5s
    setTimeout(() => {
      if (authStore.isAuthenticated) {
        connectWebSocket()
      }
    }, 5000)
  }
}

async function markAsRead(id: string) {
  try {
    await $fetch(`${config.public.apiBase}/api/v1/notifications/${id}/read`, {
      method: 'PUT',
      headers: {
        Authorization: `Bearer ${authStore.token}`
      }
    })
    const notif = notifications.value.find(n => n.id === id)
    if (notif) notif.is_read = true
  } catch (error) {
    console.error('Failed to mark notification as read:', error)
  }
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleString()
}
</script>

<template>
  <Popover>
    <PopoverTrigger asChild>
      <Button variant="ghost" size="icon" class="relative">
        <BellIcon class="w-5 h-5" />
        <span v-if="unreadCount > 0" class="absolute top-1 right-1 w-2.5 h-2.5 bg-red-500 rounded-full border border-background"></span>
      </Button>
    </PopoverTrigger>

    <PopoverContent class="w-80 p-0" align="end">
      <div class="flex flex-col">
        <div class="p-3 border-b border-border bg-muted/30 font-semibold flex justify-between items-center">
          Notifications
          <Badge v-if="unreadCount > 0" variant="default" class="text-xs">{{ unreadCount }} New</Badge>
        </div>
        
        <div v-if="notifications.length === 0" class="p-6 text-center text-muted-foreground text-sm">
          No notifications yet.
        </div>
        
        <div v-else class="max-h-96 overflow-y-auto divide-y divide-border">
          <div 
            v-for="notif in notifications" 
            :key="notif.id"
            class="p-4 transition-colors hover:bg-muted/50 cursor-pointer"
            :class="{ 'bg-primary/5': !notif.is_read }"
            @click="!notif.is_read && markAsRead(notif.id)"
          >
            <div class="flex justify-between items-start mb-1">
              <h4 class="text-sm font-semibold text-foreground flex items-center gap-2">
                <span v-if="!notif.is_read" class="w-2 h-2 rounded-full bg-primary flex-shrink-0"></span>
                {{ notif.title }}
              </h4>
              <span class="text-[10px] text-muted-foreground whitespace-nowrap ml-2">{{ formatDate(notif.created_at) }}</span>
            </div>
            <p class="text-sm text-muted-foreground line-clamp-2 pl-4" :class="{ 'pl-0': notif.is_read }">{{ notif.message }}</p>
          </div>
        </div>
      </div>
    </PopoverContent>
  </Popover>
</template>
