<template>
  <div class="space-y-8">
    <!-- Header -->
    <div class="bg-gradient-to-r from-blue-600 to-indigo-700 rounded-lg p-6 text-white">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold">Real-Time Market Data</h1>
          <p class="mt-2 text-blue-100">
            Live market insights and performance tracking
          </p>
        </div>
        <div class="flex items-center space-x-4">
          <div class="text-right">
            <div class="text-sm text-blue-200">Last Updated</div>
            <div class="text-sm font-medium">{{ lastUpdatedTime }}</div>
          </div>
          <button
            @click="refreshData"
            :disabled="isAnyLoading"
            class="px-4 py-2 bg-white/20 text-white rounded-lg hover:bg-white/30 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="isAnyLoading" class="flex items-center">
              <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
              Refreshing...
            </span>
            <span v-else>🔄 Refresh</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Error Alert -->
    <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4">
      <div class="flex items-center">
        <div class="text-red-400 mr-3">⚠️</div>
        <div class="text-red-800">{{ error }}</div>
        <button @click="clearError" class="ml-auto text-red-600 hover:text-red-800">
          ✕
        </button>
      </div>
    </div>

    <!-- Market Summary Cards -->
    <MarketSummaryCards
      :market-summary="marketSummary"
      :loading="isLoadingSummary"
      @ticker-click="handleTickerClick"
    />

    <!-- Main Content Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <!-- Left Sidebar: Performance Leaders -->
      <div class="lg:col-span-1">
        <PerformanceLeaders
          :top-performers="topPerformers"
          :worst-performers="worstPerformers"
          :loading="isLoadingPerformers"
          @stock-click="handleStockClick"
        />
      </div>

      <!-- Center: Stock Detail Explorer -->
      <div class="lg:col-span-1">
        <StockDetailExplorer
          :selected-stock="selectedStock"
          :selected-stock-trend="selectedStockTrend"
          :selected-stock-recommendation="selectedStockRecommendation"
          :is-loading-analysis="isLoadingStockAnalysis"
          :is-loading-trend="isLoadingStockTrend"
          :is-loading-recommendation="isLoadingStockRecommendation"
          :error="error"
          :has-feature="hasFeature"
          initial-ticker="WELL"
          @search="handleStockSearch"
        />
      </div>
    </div>

    <!-- Active Stocks Monitor -->
    <ActiveStocksMonitor
      :most-active-stocks="mostActiveStocks"
      :loading="isLoadingActiveStocks"
      @stock-click="handleStockClick"
    />

    <!-- Premium Features Promo -->
    <div v-if="!hasFeature('real_time_data')" class="bg-gradient-to-r from-purple-50 to-pink-50 p-6 rounded-lg border border-purple-200">
      <div class="flex items-center justify-between">
        <div>
          <h3 class="text-lg font-medium text-gray-900">Upgrade to Premium</h3>
          <p class="text-sm text-gray-600 mt-1">Get access to advanced market data, AI insights, and real-time analytics</p>
        </div>
        <button class="px-6 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700">
          Upgrade Now
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, computed, ref } from 'vue'
import { useMarketDataStore } from '@/stores/marketData'
import {
  MarketSummaryCards,
  PerformanceLeaders,
  StockDetailExplorer,
  ActiveStocksMonitor
} from '@/components/features/marketData'

const marketDataStore = useMarketDataStore()

// Local state
const selectedTicker = ref<string | null>(null)
// Deprecated local period handling removed with simplified trend controls

// Computed properties
const marketSummary = computed(() => marketDataStore.marketSummary)
const topPerformers = computed(() => marketDataStore.topPerformers)
const worstPerformers = computed(() => marketDataStore.worstPerformers)
const mostActiveStocks = computed(() => marketDataStore.mostActiveStocks)
const selectedStock = computed(() => marketDataStore.selectedStock)
const selectedStockTrend = computed(() => marketDataStore.selectedStockTrend)
const selectedStockRecommendation = computed(() => marketDataStore.selectedStockRecommendation)

// Loading states
const isLoadingSummary = computed(() => marketDataStore.isLoadingSummary)
const isLoadingPerformers = computed(() => marketDataStore.isLoadingPerformers)
const isLoadingActiveStocks = computed(() => marketDataStore.isLoadingActiveStocks)
const isLoadingStockAnalysis = computed(() => marketDataStore.isLoadingStockAnalysis)
const isLoadingStockTrend = computed(() => marketDataStore.isLoadingStockTrend)
const isLoadingStockRecommendation = computed(() => marketDataStore.isLoadingStockRecommendation)

const isAnyLoading = computed(() => 
  isLoadingSummary.value || 
  isLoadingPerformers.value || 
  isLoadingActiveStocks.value ||
  isLoadingStockAnalysis.value ||
  isLoadingStockTrend.value ||
  isLoadingStockRecommendation.value
)

// Error and feature states
const error = computed(() => marketDataStore.error)
const hasFeature = computed(() => marketDataStore.hasFeature)

// Last updated time
const lastUpdatedTime = computed(() => {
  if (marketSummary.value?.last_updated) {
    return new Date(marketSummary.value.last_updated).toLocaleTimeString()
  }
  return 'Never'
})

// Methods
async function refreshData() {
  await marketDataStore.fetchAllMarketData('1d')
}

function clearError() {
  marketDataStore.clearError()
}

function handleTickerClick(ticker: string) {
  selectedTicker.value = ticker
  handleStockSearch(ticker)
}

function handleStockClick(ticker: string) {
  selectedTicker.value = ticker
  handleStockSearch(ticker)
}

async function handleStockSearch(ticker: string) {
  if (!ticker.trim()) return
  
  try {
    await Promise.allSettled([
      marketDataStore.fetchStockAnalysis(ticker),
      // Use a consistent default period for trend
      marketDataStore.fetchStockTrend(ticker, '1w'),
      marketDataStore.fetchStockRecommendation(ticker)
    ])
  } catch (err) {
    console.error('Failed to fetch stock data:', err)
  }
}

// Lifecycle
onMounted(async () => {
  console.log('🔧 RealTimeDataView mounted, auth state:', {
    isAuthenticated: marketDataStore.isAuthenticated,
    userTier: marketDataStore.userTier,
    hasFeature: marketDataStore.hasFeature('real_time_data')
  })
  
  try {
    await marketDataStore.fetchAllMarketData('1d')
    // Preload default explorer content
    selectedTicker.value = 'WELL'
    await handleStockSearch('WELL')
  } catch (err) {
    console.error('Failed to load market data:', err)
  }
})
</script> 