<template>
  <div class="space-y-6">
    <!-- Top Performers -->
    <div class="bg-white rounded-lg shadow-sm border border-gray-200">
      <div class="px-6 py-4 border-b border-gray-200">
        <h3 class="text-lg font-semibold text-gray-900 flex items-center">
          📈 Top Performers
          <span class="ml-2 text-sm font-normal text-gray-500">(1D)</span>
        </h3>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Ticker</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Price</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Change</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Volatility</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Action</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="stock in topPerformers" :key="stock.ticker" 
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
                  <div class="w-16 bg-gray-200 rounded-full h-2">
                    <div class="bg-blue-600 h-2 rounded-full" 
                         :style="{ width: `${Math.min((stock.volatility_score || 0) * 10, 100)}%` }"></div>
                  </div>
                  <span class="ml-2 text-xs text-gray-500">{{ formatVolatilityScore(stock.volatility_score) }}</span>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                <button class="text-blue-600 hover:text-blue-800 font-medium">
                  View
                </button>
              </td>
            </tr>
            <tr v-if="topPerformers.length === 0 && !loading">
              <td colspan="5" class="px-6 py-4 text-center text-sm text-gray-500">
                No top performers data available
              </td>
            </tr>
            <tr v-if="loading">
              <td colspan="5" class="px-6 py-4 text-center">
                <div class="flex items-center justify-center">
                  <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
                  <span class="ml-2 text-sm text-gray-500">Loading...</span>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Worst Performers -->
    <div class="bg-white rounded-lg shadow-sm border border-gray-200">
      <div class="px-6 py-4 border-b border-gray-200">
        <h3 class="text-lg font-semibold text-gray-900 flex items-center">
          📉 Worst Performers
          <span class="ml-2 text-sm font-normal text-gray-500">(1D)</span>
        </h3>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Ticker</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Price</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Change</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Volatility</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Action</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="stock in worstPerformers" :key="stock.ticker" 
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
                  <div class="w-16 bg-gray-200 rounded-full h-2">
                    <div class="bg-red-600 h-2 rounded-full" 
                         :style="{ width: `${Math.min((stock.volatility_score || 0) * 10, 100)}%` }"></div>
                  </div>
                  <span class="ml-2 text-xs text-gray-500">{{ formatVolatilityScore(stock.volatility_score) }}</span>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                <button class="text-blue-600 hover:text-blue-800 font-medium">
                  View
                </button>
              </td>
            </tr>
            <tr v-if="worstPerformers.length === 0 && !loading">
              <td colspan="5" class="px-6 py-4 text-center text-sm text-gray-500">
                No worst performers data available
              </td>
            </tr>
            <tr v-if="loading">
              <td colspan="5" class="px-6 py-4 text-center">
                <div class="flex items-center justify-center">
                  <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
                  <span class="ml-2 text-sm text-gray-500">Loading...</span>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { MarketDataStock } from '@/types'

interface Props {
  topPerformers: readonly MarketDataStock[]
  worstPerformers: readonly MarketDataStock[]
  loading?: boolean
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const props = withDefaults(defineProps<Props>(), {
  loading: false
})

const emit = defineEmits<{
  stockClick: [ticker: string]
}>()

// Methods
function formatPrice(price?: number): string {
  if (price === undefined || price === null) return 'N/A'
  return price.toFixed(2)
}

function formatPercentage(value?: number): string {
  if (value === undefined || value === null) return 'N/A'
  return `${value >= 0 ? '+' : ''}${value.toFixed(2)}%`
}

function getChangeColorClass(value?: number): string {
  if (value === undefined || value === null) return 'text-gray-500'
  return value >= 0 ? 'text-green-600' : 'text-red-600'
}

function formatVolatilityScore(score?: number): string {
  if (score === undefined || score === null) return 'N/A'
  return score.toFixed(1)
}

function handleStockClick(ticker: string) {
  emit('stockClick', ticker)
}
</script> 