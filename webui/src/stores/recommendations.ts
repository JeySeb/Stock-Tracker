import { defineStore } from 'pinia'
import { ref, computed, readonly } from 'vue'
import { recommendationsAPI } from '@/api/recommendations'
import { useAuthStore } from './auth'
import type { Recommendation, RecommendationResponse } from '@/types'

export interface RecommendationFilters {
  limit?: number
  min_score?: number
  type?: string
  exclude?: string
}

export const useRecommendationsStore = defineStore('recommendations', () => {
  // State
  const recommendations = ref<Recommendation[]>([])
  const isLoading = ref(false)
  const meta = ref<RecommendationResponse['meta'] | null>(null)
  const filters = ref<RecommendationFilters>({
    limit: 10
  })
  const selectedRecommendation = ref<Recommendation | null>(null)
  const previewData = ref<any>(null)

  // Getters
  const authStore = useAuthStore()
  const maxRecommendations = computed(() => {
    const limits = { guest: 10, basic: 25, premium: 100 }
    return limits[authStore.userTier]
  })

  const availableFeatures = computed(() => meta.value?.features || [])
  const rateLimitRemaining = computed(() => meta.value?.rate_limit_remaining)

  // Actions
  async function fetchRecommendations() {
    isLoading.value = true
    try {
      // Ensure limit doesn't exceed tier maximum
      const adjustedFilters = {
        ...filters.value,
        limit: Math.min(filters.value.limit || 10, maxRecommendations.value)
      }
      
      const response = await recommendationsAPI.getRecommendations(adjustedFilters)
      recommendations.value = response.data
      meta.value = response.meta
      
      return response
    } catch (error) {
      console.error('Failed to fetch recommendations:', error)
      // Fallback to empty array to prevent crashes
      recommendations.value = []
      meta.value = null
    } finally {
      isLoading.value = false
    }
  }

  async function fetchRecommendationByTicker(ticker: string) {
    isLoading.value = true
    try {
      const response = await recommendationsAPI.getRecommendationByTicker(ticker)
      selectedRecommendation.value = response.data
      return response
    } finally {
      isLoading.value = false
    }
  }

  async function fetchPreviewForTicker(ticker: string) {
    if (!authStore.isAuthenticated) return
    
    try {
      const response = await recommendationsAPI.getPreview(ticker)
      previewData.value = response
      return response
    } catch (error) {
      console.error('Failed to fetch preview:', error)
    }
  }

  function updateFilters(newFilters: Partial<RecommendationFilters>) {
    filters.value = { ...filters.value, ...newFilters }
  }

  function clearSelectedRecommendation() {
    selectedRecommendation.value = null
    previewData.value = null
  }

  return {
    // State
    recommendations: readonly(recommendations),
    isLoading: readonly(isLoading),
    meta: readonly(meta),
    filters: readonly(filters),
    selectedRecommendation: readonly(selectedRecommendation),
    previewData: readonly(previewData),
    
    // Getters
    maxRecommendations,
    availableFeatures,
    rateLimitRemaining,
    
    // Actions
    fetchRecommendations,
    fetchRecommendationByTicker,
    fetchPreviewForTicker,
    updateFilters,
    clearSelectedRecommendation
  }
})