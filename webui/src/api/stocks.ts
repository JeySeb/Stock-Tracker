import { apiClient } from './client'
import type { StockEvent, PaginatedResponse } from '@/types'

export interface StockQueryParams {
  // Basic filters
  ticker?: string
  company?: string
  brokerage?: string
  action?: string
  limit?: number
  offset?: number
  sort_by?: string
  sort_order?: string

  // Array filters
  tickers?: string[]
  companies?: string[]
  brokerages?: string[]
  actions?: string[]

  // Rating filters
  rating_from?: number
  rating_to?: number

  // Target price filters
  target_from?: number
  target_to?: number
  min_target_change?: number
  max_target_change?: number
  has_target_price?: boolean

  // Brokerage score filters
  min_broker_score?: number
  max_broker_score?: number

  // Time filters
  last_hours?: number
  last_days?: number
  last_weeks?: number
  last_months?: number

  // Date filters
  date_from?: string // RFC3339 date
  date_to?: string // RFC3339 date
  date_ranges?: string // Multiple ranges separated by |, formatted as from,to
}

export const stocksAPI = {
  async getStocks(params?: StockQueryParams): Promise<PaginatedResponse<StockEvent>> {
    const searchParams = new URLSearchParams()
    
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== '') {
          // Handle array parameters
          if (Array.isArray(value)) {
            value.forEach(item => {
              if (item !== undefined && item !== null && item !== '') {
                searchParams.append(`${key}[]`, item.toString())
              }
            })
          } else {
            searchParams.append(key, value.toString())
          }
        }
      })
    }
    
    const queryString = searchParams.toString()
    const url = queryString ? `/stocks/enhanced?${queryString}` : '/stocks/enhanced'
    
    return apiClient.get(url)
  },

  async getStocksByTicker(ticker: string): Promise<{ data: StockEvent[] }> {
    return apiClient.get(`/stocks/${ticker.toUpperCase()}`)
  },

  async getStats(): Promise<{ data: { total_stocks: number; last_updated: string } }> {
    return apiClient.get('/stocks/stats')
  },

  async getUniqueTickers(): Promise<{ data: string[] }> {
    return apiClient.get('/stocks/tickers')
  },

  async getUniqueCompanies(): Promise<{ data: string[] }> {
    return apiClient.get('/stocks/companies')
  }
}