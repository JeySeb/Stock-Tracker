<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 overflow-y-auto">
    <div class="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
      <!-- Backdrop -->
      <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" @click="$emit('close')"></div>
      
      <!-- Modal panel -->
      <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-2xl sm:w-full">
        <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
          <div class="flex items-start justify-between">
            <div>
              <h3 class="text-lg leading-6 font-medium text-gray-900">
                {{ recommendation?.ticker || 'N/A' }} - {{ recommendation?.company_name || 'Stock Details' }}
              </h3>
              <div class="mt-2">
                <span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full"
                      :class="getRecommendationClass(recommendation?.recommendation_type)">
                  {{ recommendation?.recommendation_type || 'N/A' }}
                </span>
              </div>
            </div>
            <button @click="$emit('close')" class="text-gray-400 hover:text-gray-600">
              <span class="sr-only">Close</span>
              <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          
          <div v-if="recommendation" class="mt-4 space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div>
                <span class="text-sm font-medium text-gray-500">Basic Score</span>
                <p class="text-lg font-semibold">{{ (recommendation.basic_score*100).toFixed(2) }}/100</p>
              </div>
              <div>
                <span class="text-sm font-medium text-gray-500">Confidence</span>
                <p class="text-lg font-semibold">{{ recommendation.confidence }}</p>
              </div>
            </div>
            
            <div>
              <span class="text-sm font-medium text-gray-500">Latest Target Price</span>
              <p class="text-lg font-semibold">${{ recommendation.latest_target_price?.toFixed(2) || 'N/A' }}</p>
            </div>
            
            <div v-if="userTier !== 'guest'">
              <span class="text-sm font-medium text-gray-500">Scoring Factors</span>
              <div class="mt-2 space-y-2">
                <div v-for="factor in recommendation.scoring_factors?.slice(0, 3)" :key="factor.name"
                     class="flex justify-between items-center p-2 bg-gray-50 rounded">
                  <span class="text-sm">{{ factor.name }}</span>
                  <span class="text-sm font-medium">{{ factor.score }}/100</span>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
          <button @click="$emit('close')" 
                  class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-primary-600 text-base font-medium text-white hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 sm:ml-3 sm:w-auto sm:text-sm">
            Close
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Recommendation, UserTier } from '@/types'

interface Props {
  isOpen: boolean
  recommendation?: Recommendation | null
  previewData?: unknown
  userTier: UserTier
}

defineProps<Props>()
defineEmits<{
  close: []
  upgrade: []
}>()

function getRecommendationClass(type?: string) {
  switch (type) {
    case 'Strong Buy':
      return 'bg-green-100 text-green-800'
    case 'Buy':
      return 'bg-green-100 text-green-700'
    case 'Hold':
      return 'bg-yellow-100 text-yellow-800'
    case 'Sell':
      return 'bg-red-100 text-red-700'
    case 'Strong Sell':
      return 'bg-red-100 text-red-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}
</script>