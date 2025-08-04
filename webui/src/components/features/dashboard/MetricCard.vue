<template>
  <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6" :style="sizeStyles">
    <!-- Loading State -->
    <div v-if="loading" class="animate-pulse space-y-4">
      <div class="h-8 bg-gray-200 rounded w-1/3"></div>
      <div class="h-10 bg-gray-200 rounded w-2/3"></div>
    </div>

    <!-- Content -->
    <div v-else class="space-y-2">
      <!-- Header with Icon -->
      <div class="flex items-center justify-between">
        <h3 class="text-sm font-medium text-gray-500" :style="{ fontSize: `${0.875 * sizeFactor}rem` }">{{ title }}</h3>
        <div 
          class="rounded-full flex items-center justify-center"
          :class="getIconBackgroundClass"
          :style="{ width: `${2 * sizeFactor}rem`, height: `${2 * sizeFactor}rem` }"
        >
          <span :style="{ fontSize: `${1.125 * sizeFactor}rem` }">{{ icon }}</span>
        </div>
      </div>

      <!-- Value -->
      <div class="flex items-baseline">
        <div class="flex items-baseline">
          <span class="font-semibold text-gray-900" :style="{ fontSize: `${1.5 * sizeFactor}rem` }">
            {{ formattedValue }}
          </span>
          <span 
            v-if="subtitle" 
            class="ml-2 text-gray-500"
            :style="{ fontSize: `${0.875 * sizeFactor}rem` }"
          >
            {{ subtitle }}
          </span>
        </div>
      </div>

      <!-- Trend Indicator (if provided) -->
      <div 
        v-if="trend !== undefined" 
        class="flex items-center space-x-1"
      >
        <span 
          class="font-medium"
          :class="getTrendClass"
          :style="{ fontSize: `${0.875 * sizeFactor}rem` }"
        >
          {{ trend > 0 ? '+' : '' }}{{ trend }}%
        </span>
        <span 
          class="text-gray-500"
          v-if="trendLabel"
          :style="{ fontSize: `${0.75 * sizeFactor}rem` }"
        >
          {{ trendLabel }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

type ColorType = 'blue' | 'green' | 'red' | 'yellow' | 'purple' | 'emerald' | 'amber'

interface Props {
  /** The title of the metric card */
  title: string
  /** The main value to display */
  value?: number | null
  /** Optional subtitle to display next to the value */
  subtitle?: string
  /** Loading state of the card */
  loading?: boolean
  /** Emoji or icon to display */
  icon?: string
  /** Color theme of the card */
  color: ColorType
  /** Optional trend percentage */
  trend?: number
  /** Optional label for trend */
  trendLabel?: string
  /** Optional function to format the value */
  format?: (value: number) => string
  /** Size multiplier for the card (1 is default size) */
  sizeFactor?: number
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  icon: '📊',
  subtitle: '',
  trendLabel: '',
  format: (v: number) => v?.toLocaleString() ?? '-',
  sizeFactor: 1
})

// Computed Properties
const formattedValue = computed(() => {
  if (props.value === undefined || props.value === null) return '-'
  return props.format(props.value)
})

const getTrendClass = computed(() => {
  if (!props.trend) return 'text-gray-500'
  return props.trend > 0 
    ? 'text-green-600' 
    : props.trend < 0 
      ? 'text-red-600' 
      : 'text-gray-500'
})

const getIconBackgroundClass = computed(() => {
  const colorMap = {
    blue: ['bg-blue-100', 'text-blue-600'],
    green: ['bg-green-100', 'text-green-600'],
    red: ['bg-red-100', 'text-red-600'],
    yellow: ['bg-yellow-100', 'text-yellow-600'],
    purple: ['bg-purple-100', 'text-purple-600'],
    emerald: ['bg-emerald-100', 'text-emerald-600'],
    amber: ['bg-amber-100', 'text-amber-600']
  }
  return colorMap[props.color] || ['bg-gray-100', 'text-gray-600']
})

const sizeStyles = computed(() => ({
  padding: `${1.5 * props.sizeFactor}rem`
}))
</script>

<style scoped>
/* Dynamic color classes will be handled by Tailwind's JIT compiler */
</style>