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
                {{ stockData?.ticker || ticker || 'Stock Details' }}
              </h3>
              <p class="mt-2 text-sm text-gray-500">
                {{ stockData?.company || 'Stock event details and history' }}
              </p>
            </div>
            <button @click="$emit('close')" class="text-gray-400 hover:text-gray-600">
              <span class="sr-only">Close</span>
              <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          
          <div v-if="stockData" class="mt-4 space-y-4">
            <!-- Stock Information -->
            <div class="bg-gray-50 p-4 rounded-lg">
              <h4 class="text-sm font-medium text-gray-900 mb-3">Stock Information</h4>
              <div class="grid grid-cols-2 gap-4 text-sm">
                <div>
                  <span class="text-gray-500">Ticker:</span>
                  <span class="ml-2 font-medium text-gray-900">{{ stockData.ticker }}</span>
                </div>
                <div>
                  <span class="text-gray-500">Company:</span>
                  <span class="ml-2 font-medium text-gray-900">{{ stockData.company }}</span>
                </div>
                <div>
                  <span class="text-gray-500">Brokerage:</span>
                  <span class="ml-2 font-medium text-gray-900">{{ stockData.brokerage }}</span>
                </div>
                <div>
                  <span class="text-gray-500">Action:</span>
                  <span class="ml-2 font-medium text-gray-900">{{ stockData.action }}</span>
                </div>
              </div>
            </div>
            
            <!-- Rating Changes -->
            <div v-if="stockData.rating_from || stockData.rating_to" class="bg-blue-50 p-4 rounded-lg">
              <h4 class="text-sm font-medium text-blue-900 mb-3">Rating Changes</h4>
              <div class="flex items-center space-x-2 text-sm">
                <span class="text-blue-700">Rating:</span>
                <span v-if="stockData.rating_from" class="font-medium text-gray-900">{{ stockData.rating_from }}</span>
                <span v-if="stockData.rating_from && stockData.rating_to" class="text-blue-700">→</span>
                <span v-if="stockData.rating_to" 
                      class="font-medium"
                      :class="getRatingClass(stockData.rating_to)">
                  {{ stockData.rating_to }}
                </span>
              </div>
            </div>
            
            <!-- Price Target Analysis -->
            <div v-if="stockData.target_from || stockData.target_to" class="bg-green-50 p-4 rounded-lg">
              <h4 class="text-sm font-medium text-green-900 mb-3">Price Target Analysis</h4>
              <div class="space-y-2 text-sm">
                <div class="flex items-center space-x-2">
                  <span class="text-green-700">Target Price:</span>
                  <span v-if="stockData.target_from" class="font-medium text-gray-900">${{ stockData.target_from }}</span>
                  <span v-if="stockData.target_from && stockData.target_to" class="text-green-700">→</span>
                  <span v-if="stockData.target_to" 
                        class="font-medium"
                        :class="getTargetChangeClass(stockData.target_from, stockData.target_to)">
                    ${{ stockData.target_to }}
                  </span>
                </div>
                <div v-if="stockData.target_from && stockData.target_to" class="text-sm">
                  <span class="text-green-700">Change:</span>
                  <span class="ml-2 font-medium"
                        :class="getTargetChangeClass(stockData.target_from, stockData.target_to)">
                    {{ formatTargetChange(stockData.target_from, stockData.target_to) }}
                  </span>
                </div>
              </div>
            </div>

            <!-- Event Details -->
            <div class="bg-purple-50 p-4 rounded-lg">
              <h4 class="text-sm font-medium text-purple-900 mb-3">Event Details</h4>
              <div class="space-y-2 text-sm">
                <div>
                  <span class="text-purple-700">Event Time:</span>
                  <span class="ml-2 font-medium text-gray-900">{{ formatDateTime(stockData.event_time) }}</span>
                </div>
                <div>
                  <span class="text-purple-700">Created:</span>
                  <span class="ml-2 font-medium text-gray-900">{{ formatDateTime(stockData.created_at) }}</span>
                </div>
              </div>
            </div>

            <!-- Stock Recommendation -->
            <div v-if="recommendation" class="bg-amber-50 p-4 rounded-lg">
              <h4 class="text-sm font-medium text-amber-900 mb-3">Stock Recommendation</h4>
              <div class="space-y-3 text-sm">
                <div class="flex items-center justify-between">
                  <span class="text-amber-700">Recommendation:</span>
                  <span class="font-medium" :class="getRecommendationClass(recommendation.recommendation_type)">
                    {{ recommendation.recommendation_type }}
                  </span>
                </div>
                
                <div class="flex items-center justify-between">
                  <span class="text-amber-700">Basic Score:</span>
                  <span class="font-medium text-gray-900">{{ (recommendation.basic_score * 100).toFixed(1) }}%</span>
                </div>

                <div class="flex items-center justify-between">
                  <span class="text-amber-700">Confidence:</span>
                  <span class="font-medium text-gray-900">{{ (recommendation.confidence * 100).toFixed(1) }}%</span>
                </div>

                <div class="flex items-center justify-between">
                  <span class="text-amber-700">Latest Target Price:</span>
                  <span class="font-medium text-gray-900">${{ recommendation.latest_target_price }}</span>
                </div>

                <div class="mt-4">
                  <h5 class="text-sm font-medium text-amber-900 mb-2">Scoring Factors</h5>
                  <div class="space-y-2">
                    <div v-for="factor in recommendation.scoring_factors" :key="factor.name" 
                         class="flex items-center justify-between text-xs">
                      <span class="text-amber-700">{{ factor.name }}:</span>
                      <div class="flex items-center space-x-2">
                        <span class="font-medium text-gray-900">
                          {{ (factor.score * 100).toFixed(1) }}%
                        </span>
                        <span class="text-gray-500">({{ factor.explanation }})</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div v-else-if="isLoadingRecommendation" class="bg-amber-50 p-4 rounded-lg">
              <div class="flex items-center justify-center space-x-2">
                <span class="text-amber-700">Loading recommendation...</span>
              </div>
            </div>
          </div>

          <!-- Placeholder when no stock data -->
          <div v-else class="mt-4 space-y-4">
            <div class="bg-gray-50 p-4 rounded-lg">
              <h4 class="text-sm font-medium text-gray-900 mb-2">Stock Information</h4>
              <p class="text-sm text-gray-600">
                Detailed stock event information for {{ ticker }} would be displayed here.
                This includes broker actions, rating changes, price targets, and historical data.
              </p>
            </div>
            
            <div class="bg-blue-50 p-4 rounded-lg">
              <h4 class="text-sm font-medium text-blue-900 mb-2">Recent Events</h4>
              <p class="text-sm text-blue-700">
                Timeline of recent broker actions and recommendations would appear in this section.
              </p>
            </div>
            
            <div class="bg-green-50 p-4 rounded-lg">
              <h4 class="text-sm font-medium text-green-900 mb-2">Price Target Analysis</h4>
              <p class="text-sm text-green-700">
                Price target trends and analyst consensus information would be shown here.
              </p>
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
import { format } from 'date-fns'
import type { StockEvent, Recommendation } from '@/types'
import { recommendationsAPI } from '@/api/recommendations'
import { ref, watch } from 'vue'

interface Props {
  isOpen: boolean
  ticker?: string
  stockData?: StockEvent | null
}

const props = defineProps<Props>()
defineEmits<{
  close: []
}>()

const recommendation = ref<Recommendation | null>(null)
const isLoadingRecommendation = ref(false)

watch(() => props.isOpen, async (isOpen) => {
  if (isOpen && props.ticker) {
    isLoadingRecommendation.value = true
    try {
      const response = await recommendationsAPI.getRecommendationByTicker(props.ticker)
      recommendation.value = response.data
    } catch (error) {
      console.error('Error fetching recommendation:', error)
      recommendation.value = null
    } finally {
      isLoadingRecommendation.value = false
    }
  } else {
    recommendation.value = null
  }
})

// Helper Functions
function getRatingClass(rating: string): string {
  const ratingLower = rating.toLowerCase()
  if (ratingLower.includes('buy')) return 'text-emerald-600'
  if (ratingLower.includes('sell')) return 'text-red-600'
  return 'text-gray-600'
}

function getTargetChangeClass(from: number, to: number): string {
  const change = to - from
  if (change > 0) return 'text-emerald-600'
  if (change < 0) return 'text-red-600'
  return 'text-gray-600'
}

function formatTargetChange(from: number, to: number): string {
  const change = to - from
  const percentage = ((change / from) * 100).toFixed(1)
  return change > 0 ? `+${percentage}%` : `${percentage}%`
}

function formatDateTime(date: string): string {
  return format(new Date(date), 'MMM d, yyyy HH:mm')
}

function getRecommendationClass(type: string): string {
  const typeLower = type.toLowerCase()
  if (typeLower.includes('buy') || typeLower.includes('strong buy')) return 'text-emerald-600'
  if (typeLower.includes('sell') || typeLower.includes('strong sell')) return 'text-red-600'
  if (typeLower.includes('hold') || typeLower.includes('neutral')) return 'text-amber-600'
  return 'text-gray-600'
}
</script>