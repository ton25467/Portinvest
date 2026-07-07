<script setup lang="ts">
const props = defineProps<{
  dates: string[]
  values: number[]
}>()

const chartOption = computed(() => {
  return {
    tooltip: {
      trigger: 'axis',
      valueFormatter: (value: any) => `$${Number(value).toLocaleString()}`
    },
    grid: {
      left: '2%',
      right: '2%',
      bottom: '5%',
      top: '5%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: props.dates,
      axisLabel: { color: '#888' },
      axisLine: { lineStyle: { color: 'rgba(136, 136, 136, 0.2)' } }
    },
    yAxis: {
      type: 'value',
      axisLabel: { color: '#888' },
      splitLine: { lineStyle: { color: 'rgba(136, 136, 136, 0.1)' } }
    },
    series: [
      {
        name: 'Portfolio Value',
        type: 'line',
        smooth: true,
        showSymbol: false,
        data: props.values,
        itemStyle: { color: 'var(--color-primary-500, #3b82f6)' },
        lineStyle: { width: 3 },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(59, 130, 246, 0.3)' },
              { offset: 1, color: 'rgba(59, 130, 246, 0.0)' }
            ]
          }
        }
      }
    ]
  }
})
</script>

<template>
  <div class="h-80 w-full">
    <client-only>
      <VChart :option="chartOption" :autoresize="true" />
      <template #placeholder>
        <div class="h-full w-full flex items-center justify-center text-sm text-muted-foreground">
          Loading Performance Chart...
        </div>
      </template>
    </client-only>
  </div>
</template>
