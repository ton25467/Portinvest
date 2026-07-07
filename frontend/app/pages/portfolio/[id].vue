<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useApi } from '~/composables/useApi'
import { PlusIcon, LineChartIcon } from 'lucide-vue-next'

interface Holding {
  id: string
  portfolio_id: string
  symbol: string
  name: string
  asset_type: string
  quantity: number
  avg_buy_price: number
  current_price: number
  created_at: string
}

interface Portfolio {
  id: string
  name: string
  description: string
  currency: string
}

const route = useRoute()
const api = useApi()

const portfolioId = route.params.id as string
const portfolio = ref<Portfolio | null>(null)
const holdings = ref<Holding[]>([])
const loading = ref(true)

const showTxModal = ref(false)
const txSymbol = ref('')
const txName = ref('')
const txAssetType = ref('stock')
const txType = ref('buy')
const txQuantity = ref<number | null>(null)
const txPrice = ref<number | null>(null)
const txFee = ref(0)
const txNotes = ref('')
const txLoading = ref(false)

async function fetchData() {
  try {
    const [pData, hData] = await Promise.all([
      api.fetch<Portfolio>(`/portfolios/${portfolioId}`),
      api.fetch<Holding[]>(`/portfolios/${portfolioId}/holdings`)
    ])
    portfolio.value = pData
    holdings.value = hData
  } catch (error) {
    console.error('Failed to fetch portfolio details:', error)
  } finally {
    loading.value = false
  }
}

const tableData = computed(() => {
  return holdings.value.map(h => {
    const total_value = h.quantity * h.current_price
    const total_cost = h.quantity * h.avg_buy_price
    const pnl = total_value - total_cost
    const pnl_pct = total_cost > 0 ? (pnl / total_cost) * 100 : 0
    return {
      ...h,
      total_value,
      pnl,
      pnl_pct
    }
  })
})

const overallTotalValue = computed(() => {
  return tableData.value.reduce((sum, h) => sum + h.total_value, 0)
})

const overallTotalCost = computed(() => {
  return tableData.value.reduce((sum, h) => sum + (h.quantity * h.avg_buy_price), 0)
})

const overallPnL = computed(() => {
  return overallTotalValue.value - overallTotalCost.value
})

const overallPnLPct = computed(() => {
  return overallTotalCost.value > 0 ? (overallPnL.value / overallTotalCost.value) * 100 : 0
})

async function handleAddTransaction() {
  if (!txSymbol.value || !txQuantity.value || !txPrice.value) return

  txLoading.value = true
  try {
    await api.fetch(`/portfolios/${portfolioId}/transactions`, {
      method: 'POST',
      body: {
        symbol: txSymbol.value.toUpperCase(),
        name: txName.value || txSymbol.value.toUpperCase(),
        asset_type: txAssetType.value,
        type: txType.value,
        quantity: txQuantity.value,
        price: txPrice.value,
        fee: txFee.value,
        notes: txNotes.value
      }
    })

    showTxModal.value = false
    // Reset form
    txSymbol.value = ''
    txName.value = ''
    txQuantity.value = null
    txPrice.value = null
    txFee.value = 0
    txNotes.value = ''

    fetchData()
  } catch (error) {
    console.error('Failed to log transaction:', error)
  } finally {
    txLoading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<template>
  <div class="space-y-6">
    <!-- Header Summary Block -->
    <div v-if="portfolio" class="flex flex-col md:flex-row justify-between items-start md:items-center gap-4 bg-card border border-border p-6 rounded-xl shadow-sm">
      <div>
        <h3 class="text-xl font-bold text-foreground">{{ portfolio.name }}</h3>
        <p class="text-sm text-muted-foreground mt-1">{{ portfolio.description || 'No description' }}</p>
      </div>

      <div class="flex items-center gap-6">
        <div>
          <span class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">Total Value</span>
          <p class="text-2xl font-bold mt-1 text-foreground">
            ${{ overallTotalValue.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }}
          </p>
        </div>
        <div>
          <span class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">Profit / Loss</span>
          <p class="text-2xl font-bold mt-1" :class="[overallPnL >= 0 ? 'text-emerald-500' : 'text-rose-500']">
            {{ overallPnL >= 0 ? '+' : '' }}${{ overallPnL.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }}
            <span class="text-sm font-semibold">({{ overallPnL >= 0 ? '+' : '' }}{{ overallPnLPct.toFixed(2) }}%)</span>
          </p>
        </div>
        <Button @click="showTxModal = true">
          <PlusIcon class="mr-2 h-4 w-4" /> Add Transaction
        </Button>
      </div>
    </div>

    <!-- Transaction Modal -->
    <Dialog :open="showTxModal" @update:open="showTxModal = $event">
      <DialogContent class="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Record Transaction</DialogTitle>
        </DialogHeader>
        <div class="py-4 space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <Label for="txSymbol">Symbol</Label>
              <Input id="txSymbol" v-model="txSymbol" placeholder="e.g. AAPL, BTC" required />
            </div>
            <div class="space-y-2">
              <Label for="txName">Name</Label>
              <Input id="txName" v-model="txName" placeholder="e.g. Apple Inc." />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <Label for="txAssetType">Asset Type</Label>
              <Select v-model="txAssetType">
                <SelectTrigger>
                  <SelectValue placeholder="Select asset type" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="stock">Stock / ETF</SelectItem>
                  <SelectItem value="crypto">Cryptocurrency</SelectItem>
                  <SelectItem value="bond">Bond</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div class="space-y-2">
              <Label for="txType">Action</Label>
              <Select v-model="txType">
                <SelectTrigger>
                  <SelectValue placeholder="Select action" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="buy">Buy</SelectItem>
                  <SelectItem value="sell">Sell</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div class="grid grid-cols-3 gap-4">
            <div class="space-y-2">
              <Label for="txQuantity">Quantity</Label>
              <Input id="txQuantity" v-model.number="txQuantity" type="number" step="any" placeholder="0.0" required />
            </div>
            <div class="space-y-2">
              <Label for="txPrice">Price</Label>
              <Input id="txPrice" v-model.number="txPrice" type="number" step="any" placeholder="0.00" required />
            </div>
            <div class="space-y-2">
              <Label for="txFee">Fee</Label>
              <Input id="txFee" v-model.number="txFee" type="number" step="any" placeholder="0.00" />
            </div>
          </div>

          <div class="space-y-2">
            <Label for="txNotes">Notes</Label>
            <Input id="txNotes" v-model="txNotes" placeholder="Optional notes (e.g. brokerage name)" />
          </div>

          <div class="flex justify-end gap-3 pt-4 border-t border-border mt-4">
            <Button variant="outline" @click="showTxModal = false">Cancel</Button>
            <Button :disabled="txLoading" @click="handleAddTransaction">
              <span v-if="txLoading" class="mr-2 animate-spin">⟳</span>
              Save Transaction
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>

    <!-- Holdings Table -->
    <Card class="bg-card border-border">
      <div class="p-4 border-b border-border">
        <h4 class="text-base font-semibold leading-6 text-foreground">Holdings</h4>
      </div>

      <div v-if="loading" class="p-6">
        <div class="h-40 w-full animate-pulse bg-muted rounded-xl" />
      </div>
      <div v-else-if="holdings.length === 0" class="text-center py-12">
        <LineChartIcon class="mx-auto h-12 w-12 text-muted-foreground" />
        <h4 class="mt-2 text-sm font-semibold text-foreground">No holdings recorded</h4>
        <p class="mt-1 text-sm text-muted-foreground">Add buy transactions to build your portfolio.</p>
        <div class="mt-6">
          <Button @click="showTxModal = true">
            <PlusIcon class="mr-2 h-4 w-4" /> Add Transaction
          </Button>
        </div>
      </div>
      <div v-else class="overflow-x-auto">
        <table class="w-full text-sm text-left">
          <thead class="text-xs uppercase bg-muted/50 text-muted-foreground">
            <tr>
              <th class="px-6 py-3 font-semibold">Symbol</th>
              <th class="px-6 py-3 font-semibold">Name</th>
              <th class="px-6 py-3 font-semibold">Type</th>
              <th class="px-6 py-3 font-semibold">Holdings</th>
              <th class="px-6 py-3 font-semibold">Avg Cost</th>
              <th class="px-6 py-3 font-semibold">Current Price</th>
              <th class="px-6 py-3 font-semibold">Market Value</th>
              <th class="px-6 py-3 font-semibold">Unrealized P&L</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-border">
            <tr v-for="row in tableData" :key="row.id" class="hover:bg-muted/20 transition-colors">
              <td class="px-6 py-4 font-bold text-foreground">{{ row.symbol }}</td>
              <td class="px-6 py-4 text-sm text-muted-foreground">{{ row.name }}</td>
              <td class="px-6 py-4">
                <Badge variant="secondary" class="capitalize">{{ row.asset_type }}</Badge>
              </td>
              <td class="px-6 py-4">{{ row.quantity }}</td>
              <td class="px-6 py-4">${{ row.avg_buy_price.toFixed(2) }}</td>
              <td class="px-6 py-4">${{ row.current_price.toFixed(2) }}</td>
              <td class="px-6 py-4 font-semibold">${{ row.total_value.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }}</td>
              <td class="px-6 py-4">
                <span class="font-semibold" :class="[row.pnl >= 0 ? 'text-emerald-500' : 'text-rose-500']">
                  {{ row.pnl >= 0 ? '+' : '' }}${{ row.pnl.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }}
                  <span class="text-xs font-normal">({{ row.pnl >= 0 ? '+' : '' }}{{ row.pnl_pct.toFixed(2) }}%)</span>
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </Card>
  </div>
</template>
