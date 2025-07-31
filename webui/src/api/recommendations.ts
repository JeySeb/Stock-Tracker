import { apiClient } from './client'
import type { RecommendationResponse, Recommendation } from '@/types'

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
    
    console.log('Fetching recommendations with URL:', url)
    console.log('Parameters:', params)
    
    try {
      const response = await apiClient.get<RecommendationResponse>(url)
      console.log('Recommendations response:', response)
      return response
    } catch (error) {
      console.error('Error fetching recommendations:', {
        error,
        params,
        url
      })
      throw error
    }
  },

  async getRecommendationByTicker(ticker: string): Promise<{ data: Recommendation; meta: unknown }> {
    return apiClient.get(`/recommendations/${ticker.toUpperCase()}`)
  },

  async getPreview(ticker: string): Promise<Recommendation> {
    return apiClient.get(`/recommendations/preview/${ticker.toUpperCase()}`)
  }
}