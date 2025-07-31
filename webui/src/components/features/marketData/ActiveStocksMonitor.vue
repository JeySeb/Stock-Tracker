<template>
  <div class="bg-white rounded-lg shadow-sm border border-gray-200">
    <div class="px-6 py-4 border-b border-gray-200">
      <h3 class="text-lg font-semibold text-gray-900 flex items-center">
        🔥 Active Stocks Monitor
        <span class="ml-2 text-sm font-normal text-gray-500">(1D)</span>
      </h3>
    </div>
    
    <div class="overflow-x-auto">
      <table class="w-full">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer"
                @click="handleSort('ticker')">
              Ticker
              <span v-if="sortBy === 'ticker'" class="ml-1">{{ sortOrder === 'asc' ? '↑' : '↓' }}</span>
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer"
                @click="handleSort('current_price')">
              Price
              <span v-if="sortBy === 'current_price'" class="ml-1">{{ sortOrder === 'asc' ? '↑' : '↓' }}</span>
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer"
                @click="handleSort('day_change_percent')">
              Change
              <span v-if="sortBy === 'day_change_percent'" class="ml-1">{{ sortOrder === 'asc' ? '↑' : '↓' }}</span>
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer"
                @click="handleSort('volume')">
              Volume
              <span v-if="sortBy === 'volume'" class="ml-1">{{ sortOrder === 'asc' ? '↑' : '↓' }}</span>
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Activity</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Action</th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="stock in sortedStocks" :key="stock.ticker" 
              class="hover:bg-gray-50 cursor-pointer"
              @click="handleStockClick(stock.ticker)">
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="flex items-center">
                <div class="text-sm font-medium text-gray-900">{{ stock.ticker }}</div>
                <div class="ml-2 text-xs text-gray-500">{{ stock.company_name }}</div>
              </div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
              ${{ formatPrice(stock.current_price) }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span class="text-sm font-medium" :class="getChangeColorClass(stock.day_change_percent)">
                {{ formatPercentage(stock.day_change_percent) }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="flex items-center">
                <div class="text-sm text-gray-900">{{ formatVolume(stock.volume) }}</div>
                <div class="ml-2">
                  <div class="w-16 bg-gray-200 rounded-full h-2">
                    <div class="bg-blue-600 h-2 rounded-full transition-all duration-300" 
                         :style="{ width: `${Math.min(stock.volume_activity * 10, 100)}%` }"></div>
                  </div>
                </div>
              </div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="flex items-center">
                <span v-if="stock.volume_activity > 5" class="text-2xl mr-2">🔥</span>
                <span v-else-if="stock.volume_activity > 2" class="text-2xl mr-2">⚡</span>
                <span v-else class="text-2xl mr-2">📊</span>
                <span class="text-xs text-gray-500">{{ getActivityLabel(stock.volume_activity) }}</span>
              </div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              <button class="text-blue-600 hover:text-blue-800 font-medium">
                View
              </button>
            </td>
          </tr>
          <tr v-if="mostActiveStocks.length === 0 && !loading">
            <td colspan="6" class="px-6 py-4 text-center text-sm text-gray-500">
              No active stocks data available
            </td>
          </tr>
          <tr v-if="loading">
            <td colspan="6" class="px-6 py-4 text-center">
              <div class="flex items-center justify-center">
                <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
                <span class="ml-2 text-sm text-gray-500">Loading active stocks...</span>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { MarketDataStock } from '@/types'

interface Props {
  mostActiveStocks: readonly MarketDataStock[]
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
})

const emit = defineEmits<{
  stockClick: [ticker: string]
}>()

// Sorting state
const sortBy = ref<'ticker' | 'current_price' | 'day_change_percent' | 'volume'>('volume')
const sortOrder = ref<'asc' | 'desc'>('desc')

// Computed
const sortedStocks = computed(() => {
  if (!props.mostActiveStocks.length) return []
  
  return [...props.mostActiveStocks].sort((a, b) => {
    let aValue: string | number
    let bValue: string | number
    
    switch (sortBy.value) {
      case 'ticker':
        aValue = a.ticker
        bValue = b.ticker
        break
      case 'current_price':
        aValue = a.current_price
        bValue = b.current_price
        break
      case 'day_change_percent':
        aValue = a.day_change_percent
        bValue = b.day_change_percent
        break
      case 'volume':
        aValue = a.volume
        bValue = b.volume
        break
      default:
        return 0
    }
    
    if (typeof aValue === 'string' && typeof bValue === 'string') {
      return sortOrder.value === 'asc' 
        ? aValue.localeCompare(bValue)
        : bValue.localeCompare(aValue)
    } else {
      return sortOrder.value === 'asc'
        ? (aValue as number) - (bValue as number)
        : (bValue as number) - (aValue as number)
    }
  })
})

// Methods
function handleSort(column: typeof sortBy.value) {
  if (sortBy.value === column) {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortBy.value = column
    sortOrder.value = 'desc'
  }
}

function handleStockClick(ticker: string) {
  emit('stockClick', ticker)
}

function formatPrice(price?: number): string {
  if (price === undefined || price === null) return 'N/A'
  return price.toFixed(2)
}

function formatPercentage(value?: number): string {
  if (value === undefined || value === null) return 'N/A'
  return `${value >= 0 ? '+' : ''}${value.toFixed(2)}%`
}

function formatVolume(volume?: number): string {
  if (volume === undefined || volume === null) return 'N/A'
  if (volume >= 1e9) return `${(volume / 1e9).toFixed(1)}B`
  if (volume >= 1e6) return `${(volume / 1e6).toFixed(1)}M`
  if (volume >= 1e3) return `${(volume / 1e3).toFixed(1)}K`
  return volume.toString()
}

function getChangeColorClass(value?: number): string {
  if (value === undefined || value === null) return 'text-gray-500'
  return value >= 0 ? 'text-green-600' : 'text-red-600'
}

function getActivityLabel(activity: number): string {
  if (activity > 5) return 'Very High'
  if (activity > 2) return 'High'
  if (activity > 1) return 'Normal'
  return 'Low'
}
</script> 