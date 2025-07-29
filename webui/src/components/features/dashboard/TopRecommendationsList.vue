<template>
  <div class="space-y-4">
    <!-- Loading State -->
    <div v-if="loading" class="space-y-3">
      <div v-for="n in 5" :key="n" class="animate-pulse flex items-center space-x-4">
        <div class="h-12 w-12 bg-gray-200 rounded-full"></div>
        <div class="flex-1 space-y-2">
          <div class="h-4 bg-gray-200 rounded w-1/4"></div>
          <div class="h-3 bg-gray-200 rounded w-3/4"></div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div 
      v-else-if="!recommendations || recommendations.length === 0" 
      class="text-center py-8 text-gray-500"
    >
      No recommendations available
    </div>

    <!-- Recommendations List -->
    <div v-else class="divide-y divide-gray-200">
      <div
        v-for="recommendation in recommendations"
        :key="recommendation.id"
        class="py-3 flex items-center space-x-4 cursor-pointer hover:bg-gray-50 transition-colors"
        @click="$emit('viewDetails', recommendation.ticker)"
      >
        <!-- Recommendation Badge -->
        <div 
          class="flex-shrink-0 w-12 h-12 rounded-full flex items-center justify-center"
          :class="getBadgeClasses(recommendation.recommendation_type)"
        >
          <span class="text-lg">{{ getRecommendationIcon(recommendation.recommendation_type) }}</span>
        </div>

        <!-- Info -->
        <div class="flex-1 min-w-0">
          <div class="flex items-center justify-between">
            <h4 class="text-sm font-medium text-gray-900 truncate">
              {{ recommendation.ticker }}
              <span class="ml-2 text-xs text-gray-500">{{ recommendation.company_name }}</span>
            </h4>
            <span 
              class="text-xs font-medium"
              :class="getScoreClass(recommendation.basic_score)"
            >
              {{ (recommendation.basic_score * 100).toFixed(0) }}%
            </span>
          </div>

          <!-- Market Data (Basic/Premium only) -->
          <div 
            v-if="userTier !== 'guest' && recommendation.external_data"
            class="mt-1 text-xs text-gray-500 flex items-center space-x-4"
          >
            <span>
              ${{ recommendation.external_data.current_price.toFixed(2) }}
            </span>
            <span :class="getChangeClass(recommendation.external_data.price_change_24h)">
              {{ formatChange(recommendation.external_data.price_change_24h) }}%
            </span>
          </div>

          <!-- Basic Info (Guest) -->
          <div 
            v-else
            class="mt-1 text-xs text-gray-500"
          >
            {{ recommendation.total_events }} events · Last updated {{ formatDate(recommendation.last_event_time) }}
          </div>
        </div>

        <!-- Action Icon -->
        <div class="flex-shrink-0 text-gray-400">
          <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { format } from 'date-fns'
import type { Recommendation, UserTier } from '@/types'

interface Props {
  recommendations?: Recommendation[]
  loading?: boolean
  userTier: UserTier
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
})

defineEmits<{
  viewDetails: [ticker: string]
}>()

// Helper Functions
function getRecommendationIcon(type: string): string {
  const icons = {
    'Strong Buy': '⭐',
    'Buy': '📈',
    'Hold': '⚖️',
    'Sell': '📉',
    'Strong Sell': '⚠️'
  }
  return icons[type as keyof typeof icons] || '📊'
}

function getBadgeClasses(type: string): string {
  const baseClasses = 'bg-opacity-20'
  const typeClasses = {
    'Strong Buy': 'bg-emerald-500 text-emerald-700',
    'Buy': 'bg-emerald-400 text-emerald-600',
    'Hold': 'bg-gray-400 text-gray-600',
    'Sell': 'bg-red-400 text-red-600',
    'Strong Sell': 'bg-red-500 text-red-700'
  }
  return `${baseClasses} ${typeClasses[type as keyof typeof typeClasses] || 'bg-gray-200 text-gray-500'}`
}

function getScoreClass(score: number): string {
  if (score >= 0.7) return 'text-emerald-600'
  if (score >= 0.4) return 'text-amber-600'
  return 'text-red-600'
}

function getChangeClass(change: number): string {
  if (change > 0) return 'text-emerald-600'
  if (change < 0) return 'text-red-600'
  return 'text-gray-600'
}

function formatChange(change: number): string {
  const formatted = Math.abs(change).toFixed(1)
  return change > 0 ? `+${formatted}` : `-${formatted}`
}

function formatDate(date: string): string {
  return format(new Date(date), 'MMM d, yyyy')
}
</script>