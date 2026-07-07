<script setup lang="ts">
import { createChart } from 'lightweight-charts'

const props = defineProps<{
  data: Array<{ time: string; open: number; high: number; low: number; close: number }>
}>()

const container = ref<HTMLElement | null>(null)
let chart: any = null
let candlestickSeries: any = null

const initChart = () => {
  if (!container.value) return

  chart = createChart(container.value, {
    width: container.value.clientWidth,
    height: 320,
    layout: {
      background: { color: 'transparent' },
      textColor: '#888'
    },
    grid: {
      vertLines: { color: 'rgba(136, 136, 136, 0.1)' },
      horzLines: { color: 'rgba(136, 136, 136, 0.1)' }
    }
  })

  candlestickSeries = chart.addCandlestickSeries({
    upColor: '#10b981',
    downColor: '#ef4444',
    borderDownColor: '#ef4444',
    borderUpColor: '#10b981',
    wickDownColor: '#ef4444',
    wickUpColor: '#10b981'
  })

  candlestickSeries.setData(props.data)

  const handleResize = () => {
    if (chart && container.value) {
      chart.applyOptions({ width: container.value.clientWidth })
    }
  }

  window.addEventListener('resize', handleResize)

  onUnmounted(() => {
    window.removeEventListener('resize', handleResize)
    if (chart) {
      chart.remove()
      chart = null
    }
  })
}

onMounted(() => {
  // Ensure DOM is fully ready
  setTimeout(initChart, 50)
})

watch(() => props.data, (newData) => {
  if (candlestickSeries) {
    candlestickSeries.setData(newData)
  }
}, { deep: true })
</script>

<template>
  <div class="relative w-full">
    <div ref="container" class="w-full h-80"></div>
  </div>
</template>
