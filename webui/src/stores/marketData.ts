import { defineStore } from 'pinia'
import { ref, computed, readonly } from 'vue'
import { marketDataAPI } from '@/api/marketData'
import { useAuthStore } from '@/stores/auth'
import type { 
  MarketSummary, 
  MarketDataStock, 
  StockAnalysis, 
  MarketDataTrendResponse, 
  MarketDataRecommendation 
} from '@/types'

export const useMarketDataStore = defineStore('marketData', () => {
  // Constants
  const defaultPeriod = '1w' as const

  // State
  const marketSummary = ref<MarketSummary | null>(null)
  const topPerformers = ref<MarketDataStock[]>([])
  const worstPerformers = ref<MarketDataStock[]>([])
  const highRiskStocks = ref<MarketDataStock[]>([])
  const lowRiskStocks = ref<MarketDataStock[]>([])
  const mostVolatileStocks = ref<MarketDataStock[]>([])
  const mostActiveStocks = ref<MarketDataStock[]>([])
  const selectedStock = ref<StockAnalysis | null>(null)
  const selectedStockTrend = ref<MarketDataTrendResponse | null>(null)
  const selectedStockRecommendation = ref<MarketDataRecommendation | null>(null)
  
  // Loading states
  const isLoadingSummary = ref(false)
  const isLoadingPerformers = ref(false)
  const isLoadingRiskAnalysis = ref(false)
  const isLoadingActiveStocks = ref(false)
  const isLoadingStockAnalysis = ref(false)
  const isLoadingStockTrend = ref(false)
  const isLoadingStockRecommendation = ref(false)

  // Error states
  const error = ref<string | null>(null)

  // Auth store
  const authStore = useAuthStore()

  // Getters
  const isAuthenticated = computed(() => authStore.isAuthenticated)
  const userTier = computed(() => authStore.userTier)
  const hasFeature = computed(() => (feature: string) => authStore.hasFeature(feature))

  // Actions
  async function fetchMarketSummary() {
    isLoadingSummary.value = true
    error.value = null
    
    try {
      const data = await marketDataAPI.getMarketSummary()
      marketSummary.value = data
    } catch (err) {
      console.error('Failed to fetch market summary:', err)
      error.value = 'Failed to load market summary'
    } finally {
      isLoadingSummary.value = false
    }
  }

  async function fetchPerformers(period: '1d' | '1w' | '1m' | '3m' = defaultPeriod) {
    isLoadingPerformers.value = true
    error.value = null
    
    try {
      const [topData, worstData] = await Promise.all([
        marketDataAPI.getTopPerformers({ period, limit: 5 }),
        marketDataAPI.getWorstPerformers({ period, limit: 5 })
      ])
      
      topPerformers.value = topData
      worstPerformers.value = worstData
    } catch (err) {
      console.error('Failed to fetch performers:', err)
      error.value = 'Failed to load performance data'
    } finally {
      isLoadingPerformers.value = false
    }
  }

  async function fetchRiskAnalysis() {
    isLoadingRiskAnalysis.value = true
    error.value = null
    
    try {
      const [highRiskData, lowRiskData, volatileData] = await Promise.all([
        marketDataAPI.getHighRiskStocks({ limit: 5 }),
        marketDataAPI.getLowRiskStocks({ limit: 5 }),
        marketDataAPI.getMostVolatileStocks({ period: defaultPeriod, limit: 5 })
      ])
      
      highRiskStocks.value = highRiskData
      lowRiskStocks.value = lowRiskData
      mostVolatileStocks.value = volatileData
    } catch (err) {
      console.error('Failed to fetch risk analysis:', err)
      error.value = 'Failed to load risk analysis'
    } finally {
      isLoadingRiskAnalysis.value = false
    }
  }

  async function fetchActiveStocks(period: '1d' | '1w' | '1m' | '3m' = defaultPeriod) {
    isLoadingActiveStocks.value = true
    error.value = null
    
    try {
      const data = await marketDataAPI.getMostActiveStocks({ period, limit: 10 })
      mostActiveStocks.value = data
    } catch (err) {
      console.error('Failed to fetch active stocks:', err)
      error.value = 'Failed to load active stocks'
    } finally {
      isLoadingActiveStocks.value = false
    }
  }

  async function fetchStockAnalysis(ticker: string) {
    isLoadingStockAnalysis.value = true
    error.value = null
    
    try {
      const data = await marketDataAPI.getStockAnalysis(ticker)
      selectedStock.value = data
    } catch (err) {
      console.error('Failed to fetch stock analysis:', err)
      error.value = 'Failed to load stock analysis'
      selectedStock.value = null
    } finally {
      isLoadingStockAnalysis.value = false
    }
  }

  async function fetchStockTrend(ticker: string, period: string) {
    isLoadingStockTrend.value = true
    error.value = null
    
    try {
      const data = await marketDataAPI.getStockTrend(ticker, period)
      selectedStockTrend.value = data
    } catch (err) {
      console.error('Failed to fetch stock trend:', err)
      error.value = 'Failed to load trend data'
      selectedStockTrend.value = null
    } finally {
      isLoadingStockTrend.value = false
    }
  }

  async function fetchStockRecommendation(ticker: string) {
    isLoadingStockRecommendation.value = true
    error.value = null
    
    try {
      const data = await marketDataAPI.getStockRecommendations(ticker)
      selectedStockRecommendation.value = data
    } catch (err) {
      console.error('Failed to fetch stock recommendation:', err)
      error.value = 'Failed to load recommendation'
      selectedStockRecommendation.value = null
    } finally {
      isLoadingStockRecommendation.value = false
    }
  }

  async function fetchAllMarketData(period: '1d' | '1w' | '1m' | '3m' = defaultPeriod) {
    try {
      await Promise.allSettled([
        fetchMarketSummary(),
        fetchPerformers(period),
        fetchRiskAnalysis(),
        fetchActiveStocks(period)
      ])
    } catch (err) {
      console.error('Failed to fetch market data:', err)
      error.value = 'Failed to load market data'
    }
  }

  function clearSelectedStock() {
    selectedStock.value = null
    selectedStockTrend.value = null
    selectedStockRecommendation.value = null
  }

  function clearError() {
    error.value = null
  }

  return {
    // State
    marketSummary: readonly(marketSummary),
    topPerformers: readonly(topPerformers),
    worstPerformers: readonly(worstPerformers),
    highRiskStocks: readonly(highRiskStocks),
    lowRiskStocks: readonly(lowRiskStocks),
    mostVolatileStocks: readonly(mostVolatileStocks),
    mostActiveStocks: readonly(mostActiveStocks),
    selectedStock: readonly(selectedStock),
    selectedStockTrend: readonly(selectedStockTrend),
    selectedStockRecommendation: readonly(selectedStockRecommendation),
    
    // Loading states
    isLoadingSummary: readonly(isLoadingSummary),
    isLoadingPerformers: readonly(isLoadingPerformers),
    isLoadingRiskAnalysis: readonly(isLoadingRiskAnalysis),
    isLoadingActiveStocks: readonly(isLoadingActiveStocks),
    isLoadingStockAnalysis: readonly(isLoadingStockAnalysis),
    isLoadingStockTrend: readonly(isLoadingStockTrend),
    isLoadingStockRecommendation: readonly(isLoadingStockRecommendation),
    
    // Error state
    error: readonly(error),
    
    // Getters
    isAuthenticated,
    userTier,
    hasFeature,
    
    // Actions
    fetchMarketSummary,
    fetchPerformers,
    fetchRiskAnalysis,
    fetchActiveStocks,
    fetchStockAnalysis,
    fetchStockTrend,
    fetchStockRecommendation,
    fetchAllMarketData,
    clearSelectedStock,
    clearError
  }
}) 