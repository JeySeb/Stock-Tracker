import { apiClient } from './client'
import type { StockEvent, PaginatedResponse } from '@/types'

export interface StockQueryParams {
  ticker?: string
  company?: string
  brokerage?: string
  action?: string
  limit?: number
  offset?: number
  sort_by?: string
  sort_order?: string
}

export const stocksAPI = {
  async getStocks(params?: StockQueryParams): Promise<PaginatedResponse<StockEvent>> {
    const searchParams = new URLSearchParams()
    
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== '') {
          searchParams.append(key, value.toString())
        }
      })
    }
    
    const queryString = searchParams.toString()
    const url = queryString ? `/stocks?${queryString}` : '/stocks'
    
    return apiClient.get(url)
  },

  async getStocksByTicker(ticker: string): Promise<{ data: StockEvent[] }> {
    return apiClient.get(`/stocks/${ticker.toUpperCase()}`)
  },

  async getStats(): Promise<{ data: { total_stocks: number; last_updated: string } }> {
    return apiClient.get('/stocks/stats')
  }
}