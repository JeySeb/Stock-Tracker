<template>
  <div class="space-y-6">
    <!-- Recent Stock Events with Rating Visualizations -->
    <div>
      
      <!-- Loading State -->
      <div v-if="loading" class="space-y-3">
        <div v-for="n in 5" :key="n" class="animate-pulse flex space-x-4 p-3 border border-gray-200 rounded-lg">
          <div class="flex-1 space-y-2">
            <div class="h-4 bg-gray-200 rounded w-1/4"></div>
            <div class="h-3 bg-gray-200 rounded w-3/4"></div>
            <div class="h-2 bg-gray-200 rounded w-1/2"></div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div 
        v-else-if="!recentEvents || recentEvents.length === 0" 
        class="text-center text-gray-500 py-8"
      >
        No recent stock events available
      </div>

      <!-- Events List -->
      <div v-else class="space-y-3">
        <div
          v-for="event in recentEvents"
          :key="event.id"
          class="border border-gray-200 rounded-lg p-4 hover:bg-gray-50 transition-colors cursor-pointer"
          @click="$emit('viewStock', event.ticker)"
        >
          <div class="flex items-start justify-between">
            <!-- Left Side: Stock Info -->
            <div class="flex-1 min-w-0">
              <div class="flex items-center space-x-2 mb-2">
                <span class="font-semibold text-gray-900">{{ event.ticker }}</span>
                <span class="text-gray-400">·</span>
                <span class="text-sm text-gray-600 truncate">{{ event.company }}</span>
              </div>
              
              <!-- Brokerage and Action -->
              <div class="flex items-center space-x-2 mb-3">
                <span 
                  class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800"
                >
                  {{ event.brokerage }}
                </span>
                <span 
                  class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium"
                  :class="getActionClass(event.action)"
                >
                  {{ formatAction(event.action) }}
                </span>
              </div>

              <!-- Rating Change Visualization -->
              <div v-if="event.rating_from && event.rating_to" class="mb-3">
                <div class="flex items-center space-x-2 text-sm">
                  <span class="text-gray-500">Rating:</span>
                  <span class="font-medium text-gray-700">{{ event.rating_from }}</span>
                  <ArrowRightIcon class="h-4 w-4 text-gray-400" />
                  <span 
                    class="font-medium"
                    :class="getRatingClass(event.rating_to)"
                  >
                    {{ event.rating_to }}
                  </span>
                </div>
                
                <!-- Rating Impact Bar -->
                <div class="mt-2">
                  <div class="flex items-center space-x-2">
                    <div class="flex-1 bg-gray-200 rounded-full h-2">
                      <div 
                        class="h-2 rounded-full transition-all duration-300"
                        :class="getRatingBarClass(event.rating_from, event.rating_to)"
                        :style="{ width: `${getRatingImpactWidth(event.rating_from, event.rating_to)}%` }"
                      ></div>
                    </div>
                    <span 
                      class="text-xs font-medium"
                      :class="getRatingClass(event.rating_to)"
                    >
                      {{ getRatingImpactText(event.rating_from, event.rating_to) }}
                    </span>
                  </div>
                </div>
              </div>

              <!-- Target Price Change -->
              <div v-if="event.target_from && event.target_to" class="text-sm">
                <div class="flex items-center space-x-2">
                  <span class="text-gray-500">Target:</span>
                  <span class="text-gray-700">${{ event.target_from.toLocaleString() }}</span>
                  <ArrowRightIcon class="h-3 w-3 text-gray-400" />
                  <span 
                    class="font-medium"
                    :class="getTargetChangeClass(event.target_from, event.target_to)"
                  >
                    ${{ event.target_to.toLocaleString() }}
                  </span>
                  <span 
                    class="text-xs px-1 py-0.5 rounded"
                    :class="getTargetPercentageClass(event.target_from, event.target_to)"
                  >
                    {{ formatTargetChange(event.target_from, event.target_to) }}
                  </span>
                </div>
              </div>
            </div>

            <!-- Right Side: Time -->
            <div class="ml-4 flex-shrink-0 text-right">
              <div class="text-sm font-medium text-gray-900">{{ formatTime(event.event_time) }}</div>
              <div class="text-xs text-gray-500">{{ formatDate(event.event_time) }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { format } from 'date-fns'
import { ArrowRightIcon } from '@heroicons/vue/24/outline'
import type { StockEvent } from '@/types'

interface Props {
  events?: StockEvent[]
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
})

defineEmits<{
  viewStock: [ticker: string]
}>()

// Computed properties
const recentEvents = computed(() => {
  if (!props.events || props.events.length === 0) return []
  return props.events.slice(0, 5) // Show top 5 most recent events
})

// Helper Functions
function formatAction(action: string): string {
  return action.split(' ')[0].toUpperCase()
}

function getActionClass(action: string): string {
  const actionType = action.toLowerCase()
  
  if (actionType.includes('upgrade')) {
    return 'bg-emerald-100 text-emerald-800'
  }
  if (actionType.includes('downgrade')) {
    return 'bg-red-100 text-red-800'
  }
  if (actionType.includes('initiat')) {
    return 'bg-blue-100 text-blue-800'
  }
  return 'bg-gray-100 text-gray-800'
}

function getRatingClass(rating: string): string {
  const ratingLower = rating.toLowerCase()
  if (ratingLower.includes('strong buy') || ratingLower.includes('buy')) return 'text-emerald-600'
  if (ratingLower.includes('strong sell') || ratingLower.includes('sell')) return 'text-red-600'
  if (ratingLower.includes('hold')) return 'text-yellow-600'
  return 'text-gray-600'
}

function getRatingBarClass(from: string, to: string): string {
  const impact = getRatingImpact(from, to)
  if (impact > 0) return 'bg-emerald-500'
  if (impact < 0) return 'bg-red-500'
  return 'bg-gray-400'
}

function getRatingImpact(from: string, to: string): number {
  const ratingScores: Record<string, number> = {
    'strong sell': 1,
    'sell': 2,
    'hold': 3,
    'buy': 4,
    'strong buy': 5
  }
  
  const fromScore = ratingScores[from.toLowerCase()] || 3
  const toScore = ratingScores[to.toLowerCase()] || 3
  
  return toScore - fromScore
}

function getRatingImpactWidth(from: string, to: string): number {
  const impact = Math.abs(getRatingImpact(from, to))
  return Math.min(impact * 25, 100) // Scale to percentage
}

function getRatingImpactText(from: string, to: string): string {
  const impact = getRatingImpact(from, to)
  if (impact > 0) return '↑ Positive'
  if (impact < 0) return '↓ Negative'
  return '→ Neutral'
}

function getTargetChangeClass(from: number, to: number): string {
  const change = to - from
  if (change > 0) return 'text-emerald-600'
  if (change < 0) return 'text-red-600'
  return 'text-gray-600'
}

function getTargetPercentageClass(from: number, to: number): string {
  const change = to - from
  if (change > 0) return 'bg-emerald-100 text-emerald-800'
  if (change < 0) return 'bg-red-100 text-red-800'
  return 'bg-gray-100 text-gray-800'
}

function formatTargetChange(from: number, to: number): string {
  const change = to - from
  const percentage = ((change / from) * 100).toFixed(1)
  return change > 0 ? `+${percentage}%` : `${percentage}%`
}

function formatTime(date: string): string {
  return format(new Date(date), 'HH:mm')
}

function formatDate(date: string): string {
  return format(new Date(date), 'MMM d')
}
</script>