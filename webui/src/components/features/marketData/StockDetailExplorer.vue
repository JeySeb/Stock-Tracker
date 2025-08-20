<template>
  <div class="bg-white rounded-lg shadow-sm border border-gray-200">
    <div class="px-6 py-4 border-b border-gray-200">
      <h3 class="text-lg font-semibold text-gray-900">Stock Detail Explorer</h3>
    </div>
    
    <div class="p-6">
      <!-- Search Section -->
      <div class="mb-6">
        <div class="flex gap-4">
          <div class="flex-1">
            <input
              v-model="searchTicker"
              type="text"
              placeholder="Enter ticker symbol (e.g., AAPL)"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              @keyup.enter="handleSearch"
            />
          </div>
          <button
            @click="handleSearch"
            :disabled="!searchTicker.trim() || isLoadingAnalysis"
            class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ isLoadingAnalysis ? 'Loading...' : 'Search' }}
          </button>
        </div>
      </div>

      <!-- Stock Analysis Results -->
      <div v-if="selectedStock" class="space-y-6">
        <!-- Price Highlight Panel -->
        <div class="bg-gradient-to-r from-blue-50 to-indigo-50 p-6 rounded-lg border border-blue-200">
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-2xl font-bold text-gray-900">{{ selectedStock.ticker }}</h4>
              <p class="text-gray-600">{{ selectedStock.company_name }}</p>
            </div>
            <div class="text-right">
              <div class="text-3xl font-bold text-gray-900">${{ formatPrice(selectedStock.current_price) }}</div>
              <div class="text-lg font-medium" :class="getChangeColorClass(selectedStock.day_change_percent)">
                {{ formatPercentage(selectedStock.day_change_percent) }}
              </div>
            </div>
          </div>
        </div>

        <!-- Risk and Volatility -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div class="bg-gray-50 p-4 rounded-lg">
            <div class="text-sm font-medium text-gray-600">Risk Level</div>
            <div class="text-lg font-semibold" :class="getRiskColorClass(selectedStock.risk_level)">
              {{ selectedStock.risk_level }}
            </div>
          </div>
          <div class="bg-gray-50 p-4 rounded-lg">
            <div class="text-sm font-medium text-gray-600">Volatility Score</div>
            <div class="text-lg font-semibold text-gray-900">{{ formatVolatilityScore(selectedStock.volatility_score) }}</div>
          </div>
          <div class="bg-gray-50 p-4 rounded-lg">
            <div class="text-sm font-medium text-gray-600">Volume</div>
            <div class="text-lg font-semibold text-gray-900">{{ formatVolume(selectedStock.volume) }}</div>
          </div>
        </div>

        <!-- Trend period controls removed -->

        <!-- Trend Chart Placeholder -->
        <div v-if="selectedStockTrend" class="bg-gray-50 p-6 rounded-lg border border-gray-200">
          <div class="flex items-center justify-between mb-4">
            <h5 class="text-lg font-medium text-gray-900">Price Trend</h5>
            <span class="text-sm text-gray-500">1W period</span>
          </div>
          <div class="h-64 bg-white rounded border border-gray-200 flex items-center justify-center">
            <div class="text-center">
              <div class="text-4xl mb-2">📈</div>
              <div class="text-gray-600">Chart data available</div>
              <div class="text-sm text-gray-500">{{ selectedStockTrend.data.length }} data points</div>
            </div>
          </div>
        </div>

        <!-- Recommendation Section -->
        <div v-if="selectedStockRecommendation" class="bg-gradient-to-r from-green-50 to-blue-50 p-6 rounded-lg border border-green-200">
          <div class="flex items-center justify-between">
            <div>
              <h5 class="text-lg font-medium text-gray-900">Recommendation</h5>
              <div class="mt-2">
                <span class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium" 
                      :class="getRecommendationColorClass(selectedStockRecommendation.recommendation_type)">
                  {{ selectedStockRecommendation.recommendation_type }}
                </span>
              </div>
            </div>
            <div class="text-right">
              <div class="text-2xl font-bold text-gray-900">{{ selectedStockRecommendation.basic_score.toFixed(1) }}</div>
              <div class="text-sm text-gray-600">Score</div>
            </div>
          </div>
          <div class="mt-4 text-sm text-gray-600">
            <div>Broker Consensus: {{ selectedStockRecommendation.broker_consensus }}</div>
            <div>Confidence: {{ selectedStockRecommendation.confidence.toFixed(1) }}%</div>
          </div>
        </div>

        <!-- Premium Features Placeholder -->
        <div v-if="!hasFeature('ai_insights')" class="bg-gradient-to-r from-purple-50 to-pink-50 p-6 rounded-lg border border-purple-200">
          <div class="flex items-center justify-between">
            <div>
              <h5 class="text-lg font-medium text-gray-900">AI Insights</h5>
              <p class="text-sm text-gray-600 mt-1">Get advanced sentiment analysis and predictions</p>
            </div>
            <button class="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700">
              Upgrade to Premium
            </button>
          </div>
        </div>
      </div>

      <!-- No Stock Selected -->
      <div v-else-if="!isLoadingAnalysis" class="text-center py-12">
        <div class="text-4xl mb-4">🔍</div>
        <h4 class="text-lg font-medium text-gray-900 mb-2">Search for a Stock</h4>
        <p class="text-gray-600">Enter a ticker symbol to view detailed analysis</p>
      </div>

      <!-- Error State -->
      <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4">
        <div class="flex items-center">
          <div class="text-red-400 mr-3">⚠️</div>
          <div class="text-red-800">{{ error }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { StockAnalysis, MarketDataTrendResponse, MarketDataRecommendation } from '@/types'

interface Props {
  selectedStock: StockAnalysis | null
  selectedStockTrend: MarketDataTrendResponse | null
  selectedStockRecommendation: MarketDataRecommendation | null
  isLoadingAnalysis?: boolean
  isLoadingTrend?: boolean
  isLoadingRecommendation?: boolean
  error?: string | null
  hasFeature: (feature: string) => boolean
  initialTicker?: string
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const props = withDefaults(defineProps<Props>(), {
  isLoadingAnalysis: false,
  isLoadingTrend: false,
  isLoadingRecommendation: false,
  error: null
})

const emit = defineEmits<{
  search: [ticker: string]
}>()

// Local state
const searchTicker = ref(props.initialTicker ? props.initialTicker.toUpperCase() : '')

// Methods
function handleSearch() {
  if (searchTicker.value.trim()) {
    emit('search', searchTicker.value.trim().toUpperCase())
  }
}

// Period controls removed

function formatPrice(price?: number): string {
  if (price === undefined || price === null) return 'N/A'
  return price.toFixed(2)
}

function formatPercentage(value?: number): string {
  if (value === undefined || value === null) return 'N/A'
  return `${value >= 0 ? '+' : ''}${value.toFixed(2)}%`
}

function formatVolume(volume: number): string {
  if (volume >= 1e9) return `${(volume / 1e9).toFixed(1)}B`
  if (volume >= 1e6) return `${(volume / 1e6).toFixed(1)}M`
  if (volume >= 1e3) return `${(volume / 1e3).toFixed(1)}K`
  return volume.toString()
}

function getChangeColorClass(value: number): string {
  return value >= 0 ? 'text-green-600' : 'text-red-600'
}

function getRiskColorClass(risk: string): string {
  switch (risk) {
    case 'High': return 'text-red-600'
    case 'Medium': return 'text-yellow-600'
    case 'Low': return 'text-green-600'
    default: return 'text-gray-600'
  }
}

function formatVolatilityScore(score?: number): string {
  if (score === undefined || score === null) return 'N/A'
  return score.toFixed(1)
}

function getRecommendationColorClass(recommendation: string): string {
  switch (recommendation) {
    case 'Strong Buy': return 'bg-green-100 text-green-800'
    case 'Buy': return 'bg-blue-100 text-blue-800'
    case 'Hold': return 'bg-yellow-100 text-yellow-800'
    case 'Sell': return 'bg-orange-100 text-orange-800'
    case 'Strong Sell': return 'bg-red-100 text-red-800'
    default: return 'bg-gray-100 text-gray-800'
  }
}
</script> 