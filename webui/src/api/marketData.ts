import { apiClient } from './client'
import type { 
  MarketSummary, 
  MarketDataStock, 
  StockAnalysis, 
  MarketDataTrendResponse, 
  MarketDataRecommendation,
  APIResponse
} from '@/types'

export interface MarketDataQueryParams {
  period?: '1d' | '1w' | '1m' | '3m'
  limit?: number
}

export const marketDataAPI = {
  // Market Summary
  async getMarketSummary(): Promise<MarketSummary> {
    const searchParams = new URLSearchParams()
    //if (params?.period) {
    //  searchParams.append('period', params.period)
    //}
    
    const queryString = searchParams.toString()
    const url = queryString ? `/market-data/summary?${queryString}` : '/market-data/summary'
    
    const response = await apiClient.get<APIResponse<MarketSummary>>(url, {
      cacheKey: `market-summary`,
      cacheTTL: 300 // 5 minutes cache
    })
    
    return response.data
  },

  // Top Performers
  async getTopPerformers(params?: MarketDataQueryParams): Promise<MarketDataStock[]> {
    const searchParams = new URLSearchParams()
    if (params?.limit) searchParams.append('limit', params.limit.toString())
    //if (params?.period) searchParams.append('period', params.period)
    
    const queryString = searchParams.toString()
    const url = queryString ? `/market-data/top-performers?${queryString}` : '/market-data/top-performers'
    
    const response = await apiClient.get<APIResponse<MarketDataStock[]>>(url, {
      cacheKey: `top-performers-${params?.period || '1d'}-${params?.limit || 5}`,
      cacheTTL: 300
    })
    
    return response.data
  },

  // Worst Performers
  async getWorstPerformers(params?: MarketDataQueryParams): Promise<MarketDataStock[]> {
    const searchParams = new URLSearchParams()
    if (params?.limit) searchParams.append('limit', params.limit.toString())
    //if (params?.period) searchParams.append('period', params.period)
    
    const queryString = searchParams.toString()
    const url = queryString ? `/market-data/worst-performers?${queryString}` : '/market-data/worst-performers'
    
    const response = await apiClient.get<APIResponse<MarketDataStock[]>>(url, {
      cacheKey: `worst-performers-${params?.period || null}-${params?.limit || 5}`,
      cacheTTL: 300
    })
    
    return response.data
  },

  // High Risk Stocks
  async getHighRiskStocks(params?: { limit?: number }): Promise<MarketDataStock[]> {
    const searchParams = new URLSearchParams()
    if (params?.limit) searchParams.append('limit', params.limit.toString())
    
    const queryString = searchParams.toString()
    const url = queryString ? `/market-data/high-risk?${queryString}` : '/market-data/high-risk'
    
    const response = await apiClient.get<APIResponse<MarketDataStock[]>>(url, {
      cacheKey: `high-risk-${params?.limit || 5}`,
      cacheTTL: 300
    })
    
    return response.data
  },

  // Low Risk Stocks
  async getLowRiskStocks(params?: { limit?: number }): Promise<MarketDataStock[]> {
    const searchParams = new URLSearchParams()
    if (params?.limit) searchParams.append('limit', params.limit.toString())
    
    const queryString = searchParams.toString()
    const url = queryString ? `/market-data/low-risk?${queryString}` : '/market-data/low-risk'
    
    const response = await apiClient.get<APIResponse<MarketDataStock[]>>(url, {
      cacheKey: `low-risk-${params?.limit || 5}`,
      cacheTTL: 300
    })
    
    return response.data
  },

  // Most Volatile Stocks
  async getMostVolatileStocks(params?: MarketDataQueryParams): Promise<MarketDataStock[]> {
    const searchParams = new URLSearchParams()
    if (params?.limit) searchParams.append('limit', params.limit.toString())
    //if (params?.period) searchParams.append('period', params.period)
    
    const queryString = searchParams.toString()
    const url = queryString ? `/market-data/most-volatile?${queryString}` : '/market-data/most-volatile'
    
    const response = await apiClient.get<APIResponse<MarketDataStock[]>>(url, {
      cacheKey: `most-volatile-${params?.period ||null}-${params?.limit || 5}`,
      cacheTTL: 300
    })
    
    return response.data
  },

  // Most Active Stocks
  async getMostActiveStocks(params?: MarketDataQueryParams): Promise<MarketDataStock[]> {
    const searchParams = new URLSearchParams()
    if (params?.limit) searchParams.append('limit', params.limit.toString())
    //if (params?.period) searchParams.append('period', params.period)
    
    const queryString = searchParams.toString()
    const url = queryString ? `/market-data/most-active?${queryString}` : '/market-data/most-active'
    
    const response = await apiClient.get<APIResponse<MarketDataStock[]>>(url, {
      cacheKey: `most-active-${params?.period || null}-${params?.limit || 10}`,
      cacheTTL: 300
    })
    
    return response.data
  },

  // Stock Analysis
  async getStockAnalysis(ticker: string): Promise<StockAnalysis> {
    const response = await apiClient.get<APIResponse<StockAnalysis>>(`/market-data/analysis/${ticker.toUpperCase()}`, {
      cacheKey: `stock-analysis-${ticker.toUpperCase()}`,
      cacheTTL: 300
    })
    
    return response.data
  },

  // Stock Trend Data
  async getStockTrend(ticker: string, period: string): Promise<MarketDataTrendResponse> {
    const response = await apiClient.get<APIResponse<MarketDataTrendResponse>>(`/market-data/trend/${ticker.toUpperCase()}?period=${period}`, {
      cacheKey: `stock-trend-${ticker.toUpperCase()}-${period}`,
      cacheTTL: 300
    })
    
    return response.data
  },

  // Stock Recommendations
  async getStockRecommendations(ticker: string): Promise<MarketDataRecommendation> {
    const response = await apiClient.get<APIResponse<MarketDataRecommendation>>(`/recommendations/${ticker.toUpperCase()}`, {
      cacheKey: `stock-recommendations-${ticker.toUpperCase()}`,
      cacheTTL: 300
    })
    
    return response.data
  },

  // Preview Recommendations (for premium features)
  async getRecommendationsPreview(ticker: string): Promise<MarketDataRecommendation> {
    const response = await apiClient.get<APIResponse<MarketDataRecommendation>>(`/recommendations/preview/${ticker.toUpperCase()}`, {
      cacheKey: `recommendations-preview-${ticker.toUpperCase()}`,
      cacheTTL: 300
    })
    
    return response.data
  }
} 