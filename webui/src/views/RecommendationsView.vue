<template>
    <div class="space-y-6">
      <!-- Page Header -->
      <div class="sm:flex sm:items-center sm:justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Stock Recommendations</h1>
          <p class="mt-2 text-sm text-gray-700">
            AI-powered insights based on broker actions and market data
          </p>
        </div>
        <TierBadge :tier="authStore.userTier" />
      </div>
  
      <!-- Tier Info Card -->
      <RecommendationTierInfo
        :user-tier="authStore.userTier"
        :available-features="recommendationsStore.availableFeatures"
        :rate-limit-remaining="recommendationsStore.rateLimitRemaining"
        :max-recommendations="recommendationsStore.maxRecommendations"
      />
  
      <!-- Filters -->
      <RecommendationFilters
        :filters="recommendationsStore.filters"
        :max-recommendations="recommendationsStore.maxRecommendations"
        @update-filters="handleFiltersUpdate"
      />
  
      <!-- Recommendations Grid -->
      <div v-if="recommendationsStore.isLoading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div v-for="n in 6" :key="n" class="animate-pulse">
          <div class="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
            <div class="h-4 bg-gray-200 rounded w-1/4 mb-4"></div>
            <div class="h-6 bg-gray-200 rounded w-3/4 mb-2"></div>
            <div class="h-4 bg-gray-200 rounded w-1/2 mb-4"></div>
            <div class="space-y-2">
              <div class="h-3 bg-gray-200 rounded"></div>
              <div class="h-3 bg-gray-200 rounded w-5/6"></div>
            </div>
          </div>
        </div>
      </div>
  
      <div v-else-if="recommendationsStore.recommendations.length === 0" class="text-center py-12">
        <div class="text-gray-500">
          <p class="text-lg font-medium">No recommendations found</p>
          <p class="text-sm mt-2">Try adjusting your filters or check back later</p>
        </div>
      </div>
  
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <RecommendationCard
          v-for="recommendation in recommendationsData"
          :key="recommendation.id"
          :recommendation="recommendation"
          :user-tier="authStore.userTier"
          @view-details="handleViewDetails"
          @upgrade="handleUpgrade"
        />
      </div>
  
      <!-- Recommendation Detail Modal -->
      <RecommendationDetailModal
        :is-open="showDetailModal"
        :recommendation="selectedRecommendation"
        :preview-data="recommendationsStore.previewData"
        :user-tier="authStore.userTier"
        @close="handleCloseDetail"
        @upgrade="handleUpgrade"
      />
    </div>
  </template>
  
  <script setup lang="ts">
  import { onMounted, ref, computed } from 'vue'
  import { useRouter } from 'vue-router'
  import { useAuthStore } from '@/stores/auth'
  import { useRecommendationsStore } from '@/stores/recommendations'
  import TierBadge from '@/components/ui/TierBadge.vue'
  import RecommendationTierInfo from '@/components/features/recommendations/RecommendationTierInfo.vue'
  import RecommendationFilters from '@/components/features/recommendations/RecommendationFilters.vue'
  import RecommendationCard from '@/components/features/recommendations/RecommendationCard.vue'
  import RecommendationDetailModal from '@/components/features/recommendations/RecommendationDetailModal.vue'
  import type { Recommendation } from '@/types'
  
  const router = useRouter()
  const authStore = useAuthStore()
  const recommendationsStore = useRecommendationsStore()
  
  // Create a mutable copy of recommendations for the components
  const recommendationsData = computed(() => 
    recommendationsStore.recommendations.map(rec => ({
      ...rec,
      scoring_factors: [...rec.scoring_factors]
    }))
  )
  
  const showDetailModal = ref(false)
  const selectedRecommendation = ref<Recommendation | null>(null)
  
  onMounted(() => {
    recommendationsStore.fetchRecommendations()
  })
  
  function handleFiltersUpdate() {
    recommendationsStore.fetchRecommendations()
  }
  
  async function handleViewDetails(recommendation: Recommendation) {
    selectedRecommendation.value = recommendation
    
    // Fetch preview data if user is BASIC tier
    if (authStore.userTier === 'basic') {
      await recommendationsStore.fetchPreviewForTicker(recommendation.ticker)
    }
    
    showDetailModal.value = true
  }
  
  function handleCloseDetail() {
    showDetailModal.value = false
    selectedRecommendation.value = null
    recommendationsStore.clearSelectedRecommendation()
  }
  
  function handleUpgrade() {
    router.push('/subscription')
  }
  </script>