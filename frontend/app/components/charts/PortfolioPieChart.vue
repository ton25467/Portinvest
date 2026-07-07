<script setup lang="ts">
const props = defineProps<{
  data: Array<{ name: string; value: number }>
}>()

const chartOption = computed(() => {
  return {
    tooltip: {
      trigger: 'item',
      formatter: '{b}: ${c} ({d}%)'
    },
    legend: {
      orient: 'horizontal',
      bottom: '0%',
      textStyle: {
        color: '#888'
      }
    },
    series: [
      {
        name: 'Asset Allocation',
        type: 'pie',
        radius: ['45%', '70%'],
        avoidLabelOverlap: false,
        padAngle: 3,
        itemStyle: {
          borderRadius: 8
        },
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 16,
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: props.data
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
          Loading Allocation Chart...
        </div>
      </template>
    </client-only>
  </div>
</template>
