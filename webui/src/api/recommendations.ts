import { apiClient } from './client'
import type { RecommendationResponse } from '@/types'

export interface RecommendationQueryParams {
  limit?: number
  min_score?: number
  type?: string
  exclude?: string
}

export const recommendationsAPI = {
  async getRecommendations(params?: RecommendationQueryParams): Promise<RecommendationResponse> {
    const searchParams = new URLSearchParams()
    
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== '') {
          searchParams.append(key, value.toString())
        }
      })
    }
    
    const queryString = searchParams.toString()
    const url = queryString ? `/recommendations?${queryString}` : '/recommendations'
    
    return apiClient.get(url)
  },

  async getRecommendationByTicker(ticker: string): Promise<{ data: any; meta: any }> {
    return apiClient.get(`/recommendations/${ticker.toUpperCase()}`)
  },

  async getPreview(ticker: string): Promise<any> {
    return apiClient.get(`/recommendations/preview/${ticker.toUpperCase()}`)
  }
}