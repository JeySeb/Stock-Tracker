<template>
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
    <!-- Market Sentiment -->
    <div class="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
      <div class="flex items-center justify-between">
        <div>
          <p class="text-sm font-medium text-gray-600">Market Sentiment</p>
          <p class="text-2xl font-bold" :class="sentimentColorClass">
            {{ formatPercentage(marketSummary?.avg_day_change_percent) }}
          </p>
        </div>
        <div class="text-3xl" :class="sentimentIconClass">
          {{ sentimentIcon }}
        </div>
      </div>
      <div class="mt-2 text-sm text-gray-500">
        {{ sentimentDescription }}
      </div>
    </div>

    <!-- Bull/Bear Ratio -->
    <div class="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
      <div class="flex items-center justify-between">
        <div>
          <p class="text-sm font-medium text-gray-600">Bull/Bear Ratio</p>
          <p class="text-2xl font-bold text-gray-900">
            {{ bullBearRatio }}
          </p>
        </div>
        <div class="text-3xl">
          🎯
        </div>
      </div>
      <div class="mt-2 text-sm text-gray-500">
        {{ bullishCount }} bullish, {{ bearishCount }} bearish
      </div>
    </div>

    <!-- Most Active -->
    <div class="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
      <div class="flex items-center justify-between">
        <div>
          <p class="text-sm font-medium text-gray-600">Most Active</p>
          <p class="text-2xl font-bold text-blue-600 cursor-pointer hover:text-blue-700" 
             @click="handleTickerClick(marketSummary?.most_active_ticker)">
            {{ marketSummary?.most_active_ticker || 'N/A' }}
          </p>
        </div>
        <div class="text-3xl">
          🔥
        </div>
      </div>
      <div class="mt-2 text-sm text-gray-500">
        Highest volume today
      </div>
    </div>

    <!-- Best Performer -->
    <div class="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
      <div class="flex items-center justify-between">
        <div>
          <p class="text-sm font-medium text-gray-600">Best Performer</p>
          <p class="text-2xl font-bold text-green-600 cursor-pointer hover:text-green-700" 
             @click="handleTickerClick(marketSummary?.best_performer)">
            {{ marketSummary?.best_performer || 'N/A' }}
          </p>
        </div>
        <div class="text-3xl">
          📈
        </div>
      </div>
      <div class="mt-2 text-sm text-gray-500">
        Top gainer today
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { MarketSummary } from '@/types'

interface Props {
  marketSummary: MarketSummary | null
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
})

const emit = defineEmits<{
  tickerClick: [ticker: string]
}>()

// Computed properties
const sentimentColorClass = computed(() => {
  if (!props.marketSummary?.avg_day_change_percent) return 'text-gray-900'
  return props.marketSummary.avg_day_change_percent >= 0 ? 'text-green-600' : 'text-red-600'
})

const sentimentIconClass = computed(() => {
  if (!props.marketSummary?.avg_day_change_percent) return 'text-gray-400'
  return props.marketSummary.avg_day_change_percent >= 0 ? 'text-green-500' : 'text-red-500'
})

const sentimentIcon = computed(() => {
  if (!props.marketSummary?.avg_day_change_percent) return '➖'
  return props.marketSummary.avg_day_change_percent >= 0 ? '📈' : '📉'
})

const sentimentDescription = computed(() => {
  if (!props.marketSummary?.avg_day_change_percent) return 'No data available'
  return props.marketSummary.avg_day_change_percent >= 0 ? 'Bullish market sentiment' : 'Bearish market sentiment'
})

const bullBearRatio = computed(() => {
  if (!props.marketSummary?.bullish_count || !props.marketSummary?.bearish_count) return 'N/A'
  const ratio = props.marketSummary.bullish_count / props.marketSummary.bearish_count
  return ratio.toFixed(2)
})

const bullishCount = computed(() => props.marketSummary?.bullish_count || 0)
const bearishCount = computed(() => props.marketSummary?.bearish_count || 0)

// Methods
function formatPercentage(value?: number): string {
  if (value === undefined || value === null) return 'N/A'
  return `${value >= 0 ? '+' : ''}${value.toFixed(2)}%`
}

function handleTickerClick(ticker?: string) {
  if (ticker) {
    emit('tickerClick', ticker)
  }
}
</script> 