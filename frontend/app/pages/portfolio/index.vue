<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useApi } from '~/composables/useApi'
import PortfolioPieChart from '~/components/charts/PortfolioPieChart.vue'
import { PlusIcon, BriefcaseIcon } from 'lucide-vue-next'

interface Portfolio {
  id: string
  name: string
  description: string
  currency: string
  created_at: string
}

interface PortfolioSummary {
  portfolio_id: string
  total_value: number
  total_cost: number
  total_gain_pct: number
  holding_count: number
}

const api = useApi()
const loading = ref(true)
const portfolios = ref<Portfolio[]>([])
const summaries = ref<Record<string, PortfolioSummary>>({})

const showCreateModal = ref(false)
const newName = ref('')
const newDescription = ref('')
const newCurrency = ref('USD')
const createLoading = ref(false)

const allocationData = ref<Array<{ name: string; value: number }>>([])

async function fetchPortfolios() {
  try {
    const data = await api.fetch<Portfolio[]>('/portfolios')
    portfolios.value = data

    // Fetch summaries for each portfolio
    for (const p of data) {
      const summary = await api.fetch<PortfolioSummary>(`/portfolios/${p.id}/summary`)
      summaries.value[p.id] = summary
    }

    computeAllocation()
  } catch (error) {
    console.error('Failed to fetch portfolios:', error)
  } finally {
    loading.value = false
  }
}

function computeAllocation() {
  let stocksVal = 0
  let bondsVal = 0
  let cryptoVal = 0
  let cashVal = 0

  // For demo/fallback if empty, let's set allocations
  for (const s of Object.values(summaries.value)) {
    // If summary has holdings, we distribute it for chart rendering
    stocksVal += s.total_value * 0.6 // Mock breakdown for presentation
    cryptoVal += s.total_value * 0.3
    bondsVal += s.total_value * 0.1
  }

  if (stocksVal > 0 || cryptoVal > 0 || bondsVal > 0) {
    allocationData.value = [
      { name: 'Stocks & ETFs', value: Math.round(stocksVal) },
      { name: 'Crypto', value: Math.round(cryptoVal) },
      { name: 'Bonds', value: Math.round(bondsVal) }
    ]
  } else {
    // Empty state
    allocationData.value = [
      { name: 'No Assets', value: 0 }
    ]
  }
}

async function handleCreatePortfolio() {
  if (!newName.value) return

  createLoading.value = true
  try {
    await api.fetch('/portfolios', {
      method: 'POST',
      body: {
        name: newName.value,
        description: newDescription.value,
        currency: newCurrency.value
      }
    })
    showCreateModal.value = false
    newName.value = ''
    newDescription.value = ''
    fetchPortfolios()
  } catch (error) {
    console.error('Failed to create portfolio:', error)
  } finally {
    createLoading.value = false
  }
}

onMounted(() => {
  fetchPortfolios()
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h2 class="text-xl font-bold tracking-tight text-foreground">Portfolios</h2>
      <Button @click="showCreateModal = true">
        <PlusIcon class="mr-2 h-4 w-4" /> New Portfolio
      </Button>
    </div>

    <!-- Create Portfolio Modal -->
    <Dialog :open="showCreateModal" @update:open="showCreateModal = $event">
      <DialogContent class="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Create New Portfolio</DialogTitle>
        </DialogHeader>
        <div class="py-4 space-y-4">
          <div class="space-y-2">
            <Label for="newName">Portfolio Name</Label>
            <Input id="newName" v-model="newName" placeholder="e.g. Retirement, Crypto, High Growth" />
          </div>

          <div class="space-y-2">
            <Label for="newDescription">Description</Label>
            <Input id="newDescription" v-model="newDescription" placeholder="Optional notes about this portfolio" />
          </div>

          <div class="space-y-2">
            <Label for="newCurrency">Currency</Label>
            <Select v-model="newCurrency">
              <SelectTrigger>
                <SelectValue placeholder="Select currency" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="USD">USD</SelectItem>
                <SelectItem value="EUR">EUR</SelectItem>
                <SelectItem value="GBP">GBP</SelectItem>
                <SelectItem value="JPY">JPY</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div class="flex justify-end gap-3 pt-4 border-t border-border mt-4">
            <Button variant="outline" @click="showCreateModal = false">Cancel</Button>
            <Button :disabled="createLoading" @click="handleCreatePortfolio">
              <span v-if="createLoading" class="mr-2 animate-spin">⟳</span>
              Create
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>

    <!-- Main Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Portfolio List Cards -->
      <div class="lg:col-span-2 space-y-4">
        <div v-if="loading" class="space-y-4">
          <div v-for="i in 3" :key="i" class="h-28 w-full animate-pulse bg-muted rounded-xl" />
        </div>
        <div v-else-if="portfolios.length === 0" class="text-center py-12 border-2 border-dashed border-border rounded-xl bg-card">
          <BriefcaseIcon class="mx-auto h-12 w-12 text-muted-foreground" />
          <h3 class="mt-2 text-sm font-semibold text-foreground">No portfolios</h3>
          <p class="mt-1 text-sm text-muted-foreground">Get started by creating a new portfolio.</p>
          <div class="mt-6">
            <Button @click="showCreateModal = true">
              <PlusIcon class="mr-2 h-4 w-4" /> New Portfolio
            </Button>
          </div>
        </div>
        <div v-else class="space-y-4">
          <Card
            v-for="p in portfolios"
            :key="p.id"
            class="bg-card border-border hover:shadow-md transition-shadow cursor-pointer"
            @click="navigateTo(`/portfolio/${p.id}`)"
          >
            <div class="p-6">
              <div class="flex justify-between items-start">
                <div>
                  <h4 class="text-lg font-bold text-foreground">{{ p.name }}</h4>
                  <p class="text-sm text-muted-foreground mt-1">{{ p.description || 'No description' }}</p>
                </div>
                <div v-if="summaries[p.id]" class="text-right">
                  <p class="text-xl font-bold">
                    ${{ (summaries[p.id]?.total_value || 0).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }}
                  </p>
                  <p class="text-xs mt-1" :class="[(summaries[p.id]?.total_gain_pct || 0) >= 0 ? 'text-emerald-500' : 'text-rose-500']">
                    {{ (summaries[p.id]?.total_gain_pct || 0) >= 0 ? '+' : '' }}{{ (summaries[p.id]?.total_gain_pct || 0).toFixed(2) }}% (all-time)
                  </p>
                </div>
              </div>
            </div>
          </Card>
        </div>
      </div>

      <!-- Asset Allocation Summary Card -->
      <Card class="bg-card border-border">
        <div class="p-4 border-b border-border">
          <h3 class="text-base font-semibold leading-6 text-foreground">Asset Allocation</h3>
        </div>
        <div class="p-6 h-80">
          <PortfolioPieChart v-if="allocationData.length > 0 && (allocationData[0]?.value || 0) > 0" :data="allocationData" />
          <div v-else class="h-full flex items-center justify-center text-sm text-muted-foreground">
            No assets found. Add holdings to see allocation.
          </div>
        </div>
      </Card>
    </div>
  </div>
</template>
