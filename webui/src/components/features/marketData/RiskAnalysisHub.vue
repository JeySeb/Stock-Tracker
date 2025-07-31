<template>
  <div class="bg-white rounded-lg shadow-sm border border-gray-200">
    <div class="px-6 py-4 border-b border-gray-200">
      <h3 class="text-lg font-semibold text-gray-900">Risk Analysis Hub</h3>
    </div>
    
    <div class="p-6">
      <!-- Tab Navigation -->
      <div class="border-b border-gray-200">
        <nav class="-mb-px flex space-x-8">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            @click="activeTab = tab.id"
            :class="[
              activeTab === tab.id
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300',
              'whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm'
            ]"
          >
            {{ tab.name }}
          </button>
        </nav>
      </div>

      <!-- Tab Content -->
      <div class="mt-6">
        <!-- High Risk Tab -->
        <div v-if="activeTab === 'high-risk'" class="space-y-4">
          <div class="flex items-center mb-4">
            <span class="text-2xl mr-2">⚠️</span>
            <h4 class="text-lg font-medium text-gray-900">High Risk Stocks</h4>
          </div>
          <div class="space-y-3">
            <div v-for="stock in highRiskStocks" :key="stock.ticker" 
                 class="flex items-center justify-between p-4 bg-red-50 rounded-lg border border-red-200 cursor-pointer hover:bg-red-100"
                 @click="handleStockClick(stock.ticker)">
              <div class="flex items-center">
                <div class="w-3 h-3 bg-red-500 rounded-full mr-3"></div>
                <div>
                  <div class="font-medium text-gray-900">{{ stock.ticker }}</div>
                  <div class="text-sm text-gray-500">{{ stock.company_name }}</div>
                </div>
              </div>
              <div class="text-right">
                <div class="font-medium text-gray-900">${{ formatPrice(stock.current_price) }}</div>
                <div class="text-sm" :class="getChangeColorClass(stock.day_change_percent)">
                  {{ formatPercentage(stock.day_change_percent) }}
                </div>
              </div>
            </div>
            <div v-if="highRiskStocks.length === 0 && !loading" class="text-center text-gray-500 py-4">
              No high risk stocks available
            </div>
          </div>
        </div>

        <!-- Low Risk Tab -->
        <div v-if="activeTab === 'low-risk'" class="space-y-4">
          <div class="flex items-center mb-4">
            <span class="text-2xl mr-2">🛡️</span>
            <h4 class="text-lg font-medium text-gray-900">Low Risk Stocks</h4>
          </div>
          <div class="space-y-3">
            <div v-for="stock in lowRiskStocks" :key="stock.ticker" 
                 class="flex items-center justify-between p-4 bg-green-50 rounded-lg border border-green-200 cursor-pointer hover:bg-green-100"
                 @click="handleStockClick(stock.ticker)">
              <div class="flex items-center">
                <div class="w-3 h-3 bg-green-500 rounded-full mr-3"></div>
                <div>
                  <div class="font-medium text-gray-900">{{ stock.ticker }}</div>
                  <div class="text-sm text-gray-500">{{ stock.company_name }}</div>
                </div>
              </div>
              <div class="text-right">
                <div class="font-medium text-gray-900">${{ formatPrice(stock.current_price) }}</div>
                <div class="text-sm" :class="getChangeColorClass(stock.day_change_percent)">
                  {{ formatPercentage(stock.day_change_percent) }}
                </div>
              </div>
            </div>
            <div v-if="lowRiskStocks.length === 0 && !loading" class="text-center text-gray-500 py-4">
              No low risk stocks available
            </div>
          </div>
        </div>

        <!-- Volatile Tab -->
        <div v-if="activeTab === 'volatile'" class="space-y-4">
          <div class="flex items-center mb-4">
            <span class="text-2xl mr-2">⚡</span>
            <h4 class="text-lg font-medium text-gray-900">Most Volatile Stocks</h4>
          </div>
          <div class="space-y-3">
            <div v-for="stock in mostVolatileStocks" :key="stock.ticker" 
                 class="flex items-center justify-between p-4 bg-yellow-50 rounded-lg border border-yellow-200 cursor-pointer hover:bg-yellow-100"
                 @click="handleStockClick(stock.ticker)">
              <div class="flex items-center">
                <div class="w-3 h-3 bg-yellow-500 rounded-full mr-3"></div>
                <div>
                  <div class="font-medium text-gray-900">{{ stock.ticker }}</div>
                  <div class="text-sm text-gray-500">{{ stock.company_name }}</div>
                </div>
              </div>
              <div class="text-right">
                <div class="font-medium text-gray-900">${{ formatPrice(stock.current_price) }}</div>
                <div class="text-sm" :class="getChangeColorClass(stock.day_change_percent)">
                  {{ formatPercentage(stock.day_change_percent) }}
                </div>
                <div class="text-xs text-gray-500">Vol: {{ formatVolatilityScore(stock.volatility_score) }}</div>
              </div>
            </div>
            <div v-if="mostVolatileStocks.length === 0 && !loading" class="text-center text-gray-500 py-4">
              No volatile stocks available
            </div>
          </div>
        </div>

        <!-- Loading State -->
        <div v-if="loading" class="flex items-center justify-center py-8">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          <span class="ml-3 text-gray-500">Loading risk analysis...</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { MarketDataStock } from '@/types'

interface Props {
  highRiskStocks: readonly MarketDataStock[]
  lowRiskStocks: readonly MarketDataStock[]
  mostVolatileStocks: readonly MarketDataStock[]
  loading?: boolean
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const props = withDefaults(defineProps<Props>(), {
  loading: false
})

const emit = defineEmits<{
  stockClick: [ticker: string]
}>()

// Tab management
const activeTab = ref('high-risk')

const tabs = [
  { id: 'high-risk', name: 'High Risk' },
  { id: 'low-risk', name: 'Low Risk' },
  { id: 'volatile', name: 'Volatile' }
]

// Methods
function formatPrice(price: number): string {
  return price.toFixed(2)
}

function formatPercentage(value: number): string {
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