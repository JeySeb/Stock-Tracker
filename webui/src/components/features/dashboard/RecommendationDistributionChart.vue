<template>
  <div class="h-64">
    <!-- Loading State -->
    <div v-if="loading" class="h-full flex items-center justify-center">
      <div class="animate-pulse space-y-4 w-full">
        <div class="h-4 bg-gray-200 rounded w-1/4 mx-auto"></div>
        <div class="h-40 bg-gray-200 rounded mx-8"></div>
      </div>
    </div>

    <!-- Empty State -->
    <div 
      v-else-if="!data || data.length === 0" 
      class="h-full flex items-center justify-center text-gray-500"
    >
      No data available
    </div>

    <!-- Chart -->
    <v-chart
      v-else
      class="h-full w-full"
      :option="chartOption"
      :autoresize="true"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import type { PieSeriesOption } from 'echarts'

interface ChartData {
  name: string
  value: number
}

interface Props {
  data?: ChartData[]
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
})

// Color mapping for different recommendation types
const colorMap = {
  'Strong Buy': '#10B981', // emerald-500
  'Buy': '#34D399',       // emerald-400
  'Hold': '#9CA3AF',      // gray-400
  'Sell': '#F87171',      // red-400
  'Strong Sell': '#EF4444' // red-500
}

// Chart configuration
const chartOption = computed(() => ({
  tooltip: {
    trigger: 'item',
    formatter: '{b}: {c} ({d}%)'
  },
  legend: {
    orient: 'vertical',
    right: '10%',
    top: 'center',
    formatter: (name: string) => {
      const item = props.data?.find(d => d.name === name)
      return `${name}: ${item?.value || 0}`
    }
  },
  series: [
    {
      type: 'pie',
      radius: ['40%', '70%'],
      center: ['40%', '50%'],
      avoidLabelOverlap: true,
      itemStyle: {
        borderRadius: 4,
        borderColor: '#fff',
        borderWidth: 2
      },
      label: {
        show: false
      },
      emphasis: {
        label: {
          show: true,
          fontSize: '14',
          fontWeight: 'bold'
        }
      },
      labelLine: {
        show: false
      },
      data: props.data?.map(item => ({
        ...item,
        itemStyle: {
          color: colorMap[item.name as keyof typeof colorMap] || '#CBD5E1'
        }
      }))
    } as PieSeriesOption
  ]
}))
</script>

<style scoped>
:deep(.echarts) {
  width: 100%;
  height: 100%;
}
</style>