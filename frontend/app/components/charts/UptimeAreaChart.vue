<script setup lang="ts">
const props = defineProps<{
  timestamps: string[]
  responseTimes: number[]
}>()

const chartOption = computed(() => {
  return {
    tooltip: {
      trigger: 'axis',
      valueFormatter: (value: any) => `${value} ms`
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
      data: props.timestamps,
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
        name: 'Response Time',
        type: 'line',
        smooth: true,
        showSymbol: false,
        data: props.responseTimes,
        itemStyle: { color: '#10b981' },
        lineStyle: { width: 2 },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(16, 185, 129, 0.25)' },
              { offset: 1, color: 'rgba(16, 185, 129, 0.0)' }
            ]
          }
        }
      }
    ]
  }
})
</script>

<template>
  <div class="h-64 w-full">
    <client-only>
      <VChart :option="chartOption" :autoresize="true" />
      <template #placeholder>
        <div class="h-full w-full flex items-center justify-center text-sm text-muted-foreground">
          Loading Response Time Chart...
        </div>
      </template>
    </client-only>
  </div>
</template>
