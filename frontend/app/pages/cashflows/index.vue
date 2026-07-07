<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
import { toast } from 'vue-sonner'
import { SettingsIcon, RefreshCwIcon, PlusIcon, Loader2Icon, ArrowDownLeftIcon, ArrowUpRightIcon, Trash2Icon, XIcon } from 'lucide-vue-next'

const authStore = useAuthStore()
const config = useRuntimeConfig()

interface Cashflow {
  id: string
  portfolio_id?: string
  type: 'income' | 'expense' | 'deposit' | 'withdrawal'
  amount: number
  currency: string
  description: string
  executed_at: string
}

const cashflows = ref<Cashflow[]>([])
const loading = ref(false)
const syncLoading = ref(false)

const isAddModalOpen = ref(false)
const isSettingsModalOpen = ref(false)

const newCashflow = ref({
  type: 'income',
  amount: 0,
  currency: 'THB',
  description: ''
})

onMounted(async () => {
  await fetchCashflows()
})

async function fetchCashflows() {
  loading.value = true
  try {
    const res = await $fetch<Cashflow[]>(`${config.public.apiBase}/api/v1/cashflows`, {
      headers: { Authorization: `Bearer ${authStore.token}` }
    })
    if (res) {
      cashflows.value = res
    }
  } catch (error) {
    toast.error('Failed to load cashflows')
  } finally {
    loading.value = false
  }
}

async function handleAddCashflow() {
  try {
    await $fetch(`${config.public.apiBase}/api/v1/cashflows`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${authStore.token}` },
      body: newCashflow.value
    })
    toast.success('Cashflow added')
    isAddModalOpen.value = false
    await fetchCashflows()
  } catch (error) {
    toast.error('Failed to add cashflow')
  }
}

async function deleteCashflow(id: string) {
  try {
    await $fetch(`${config.public.apiBase}/api/v1/cashflows/${id}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${authStore.token}` }
    })
    toast.success('Cashflow deleted')
    await fetchCashflows()
  } catch (error) {
    toast.error('Failed to delete cashflow')
  }
}

async function triggerEmailSync() {
  syncLoading.value = true
  try {
    await $fetch(`${config.public.apiBase}/api/v1/cashflows/sync`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${authStore.token}` }
    })
    toast.success('Email sync completed successfully')
    await fetchCashflows()
  } catch (error) {
    toast.error('Failed to sync emails')
  } finally {
    syncLoading.value = false
  }
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleString()
}

const bankCredentials = ref({
  bank_name: 'KBANK',
  username: '',
  password: ''
})

async function saveBankCredentials() {
  try {
    await $fetch(`${config.public.apiBase}/api/v1/settings/banks`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${authStore.token}` },
      body: bankCredentials.value
    })
    toast.success('Bank credentials saved')
    isSettingsModalOpen.value = false
  } catch (error) {
    toast.error('Failed to save credentials')
  }
}

const typeColors: Record<string, string> = {
  income: 'text-green-500',
  deposit: 'text-green-500',
  expense: 'text-red-500',
  withdrawal: 'text-red-500'
}
</script>

<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <div>
        <h2 class="text-2xl font-bold tracking-tight text-foreground">Cashflows</h2>
        <p class="text-muted-foreground text-sm mt-1">Track your income, expenses, and portfolio deposits</p>
      </div>
      <div class="flex gap-2">
        <Button variant="ghost" size="icon" @click="isSettingsModalOpen = true">
          <SettingsIcon class="h-4 w-4" />
        </Button>
        <Button variant="default" :disabled="syncLoading" @click="triggerEmailSync">
          <RefreshCwIcon v-if="!syncLoading" class="mr-2 h-4 w-4" />
          <span v-else class="mr-2 animate-spin">⟳</span>
          Sync Bank Data
        </Button>
        <Button variant="outline" @click="isAddModalOpen = true">
          <PlusIcon class="mr-2 h-4 w-4" /> Add Record
        </Button>
      </div>
    </div>

    <!-- Data Table -->
    <Card class="ring-border/50 shadow-sm border-border">
      <div v-if="loading" class="p-8 flex justify-center">
        <Loader2Icon class="w-8 h-8 animate-spin text-primary" />
      </div>
      
      <div v-else-if="cashflows.length === 0" class="p-8 text-center text-muted-foreground">
        No cashflow records found. Add one or sync from your email.
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-sm text-left">
          <thead class="text-xs uppercase bg-muted/50 text-muted-foreground">
            <tr>
              <th class="px-6 py-3 font-semibold">Type</th>
              <th class="px-6 py-3 font-semibold">Description</th>
              <th class="px-6 py-3 font-semibold text-right">Amount</th>
              <th class="px-6 py-3 font-semibold">Date</th>
              <th class="px-6 py-3 font-semibold">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-border">
            <tr v-for="cf in cashflows" :key="cf.id" class="hover:bg-muted/20 transition-colors">
              <td class="px-6 py-4 flex items-center gap-2 capitalize font-medium text-foreground">
                <ArrowDownLeftIcon v-if="cf.type === 'income' || cf.type === 'deposit'" :class="typeColors[cf.type]" class="w-4 h-4" />
                <ArrowUpRightIcon v-if="cf.type === 'expense' || cf.type === 'withdrawal'" :class="typeColors[cf.type]" class="w-4 h-4" />
                {{ cf.type }}
              </td>
              <td class="px-6 py-4 text-foreground">{{ cf.description }}</td>
              <td class="px-6 py-4 text-right font-medium" :class="typeColors[cf.type]">
                {{ cf.type === 'expense' || cf.type === 'withdrawal' ? '-' : '+' }}
                {{ cf.amount.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }} {{ cf.currency }}
              </td>
              <td class="px-6 py-4 text-muted-foreground">{{ formatDate(cf.executed_at) }}</td>
              <td class="px-6 py-4">
                <Button variant="ghost" size="icon" @click="deleteCashflow(cf.id)">
                  <Trash2Icon class="h-4 w-4 text-red-500" />
                </Button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </Card>

    <!-- Add Modal -->
    <Dialog :open="isAddModalOpen" @update:open="isAddModalOpen = $event">
      <DialogContent class="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Add Cashflow Record</DialogTitle>
        </DialogHeader>
        <form @submit.prevent="handleAddCashflow" class="space-y-4 py-4">
          <div class="space-y-2">
            <Label for="type">Type</Label>
            <Select v-model="newCashflow.type">
              <SelectTrigger>
                <SelectValue placeholder="Select type" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="income">Income</SelectItem>
                <SelectItem value="expense">Expense</SelectItem>
                <SelectItem value="deposit">Deposit</SelectItem>
                <SelectItem value="withdrawal">Withdrawal</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div class="space-y-2">
            <Label for="amount">Amount</Label>
            <Input id="amount" v-model.number="newCashflow.amount" type="number" step="0.01" min="0" required />
          </div>
          <div class="space-y-2">
            <Label for="currency">Currency</Label>
            <Input id="currency" v-model="newCashflow.currency" required />
          </div>
          <div class="space-y-2">
            <Label for="description">Description</Label>
            <Input id="description" v-model="newCashflow.description" placeholder="e.g. Salary, Groceries" required />
          </div>
          <div class="flex justify-end gap-2 pt-4">
            <Button type="button" variant="outline" @click="isAddModalOpen = false">Cancel</Button>
            <Button type="submit">Save</Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>

    <!-- Bank Settings Modal -->
    <Dialog :open="isSettingsModalOpen" @update:open="isSettingsModalOpen = $event">
      <DialogContent class="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Bank Connection Settings</DialogTitle>
        </DialogHeader>
        <form @submit.prevent="saveBankCredentials" class="space-y-4 py-4">
          <div class="space-y-2">
            <Label for="bank">Bank</Label>
            <Select v-model="bankCredentials.bank_name">
              <SelectTrigger>
                <SelectValue placeholder="Select bank" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="KBANK">KBANK</SelectItem>
                <SelectItem value="KTB">KTB</SelectItem>
                <SelectItem value="SCBAM">SCBAM</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div class="space-y-2">
            <Label for="username">Username</Label>
            <Input id="username" v-model="bankCredentials.username" required />
          </div>
          <div class="space-y-2">
            <Label for="bank-password">Password</Label>
            <Input id="bank-password" v-model="bankCredentials.password" type="password" required />
          </div>
          <p class="text-xs text-muted-foreground">Credentials are stored securely using encryption.</p>
          <div class="flex justify-end gap-2 pt-4">
            <Button type="button" variant="outline" @click="isSettingsModalOpen = false">Cancel</Button>
            <Button type="submit">Save Credentials</Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  </div>
</template>
