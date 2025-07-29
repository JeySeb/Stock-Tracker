<template>
  <div class="divide-y divide-gray-200">
    <!-- Loading State -->
    <div v-if="loading" class="space-y-3 p-6">
      <div v-for="n in 3" :key="n" class="animate-pulse flex space-x-4">
        <div class="flex-1 space-y-2">
          <div class="h-4 bg-gray-200 rounded w-1/4"></div>
          <div class="h-3 bg-gray-200 rounded w-3/4"></div>
        </div>
        <div class="h-8 w-24 bg-gray-200 rounded"></div>
      </div>
    </div>

    <!-- Empty State -->
    <div 
      v-else-if="!events || events.length === 0" 
      class="p-6 text-center text-gray-500"
    >
      No recent stock events
    </div>

    <!-- Events List -->
    <div v-else>
      <div
        v-for="event in events"
        :key="event.id"
        class="p-4 hover:bg-gray-50 transition-colors cursor-pointer"
        @click="$emit('viewDetails', event.ticker)"
      >
        <div class="flex items-start justify-between">
          <!-- Event Info -->
          <div class="flex-1 min-w-0">
            <div class="flex items-center space-x-2">
              <span class="font-medium text-gray-900">{{ event.ticker }}</span>
              <span class="text-gray-500">·</span>
              <span class="text-sm text-gray-500 truncate">{{ event.company }}</span>
            </div>
            
            <div class="mt-1 flex items-center space-x-2 text-sm">
              <!-- Action Badge -->
              <span 
                class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium"
                :class="getActionClass(event.action)"
              >
                {{ formatAction(event.action) }}
              </span>

              <!-- Rating Change -->
              <span class="text-gray-500">by</span>
              <span class="font-medium text-gray-900">{{ event.brokerage }}</span>
              
              <template v-if="event.rating_from && event.rating_to">
                <span class="text-gray-500">from</span>
                <span class="text-gray-600">{{ event.rating_from }}</span>
                <span class="text-gray-500">to</span>
                <span 
                  class="font-medium"
                  :class="getRatingClass(event.rating_to)"
                >
                  {{ event.rating_to }}
                </span>
              </template>
            </div>

            <!-- Target Price Change -->
            <div 
              v-if="event.target_from && event.target_to"
              class="mt-1 text-sm"
            >
              <span class="text-gray-500">Target Price:</span>
              <span class="text-gray-900">${{ event.target_from }}</span>
              <span class="text-gray-500">→</span>
              <span 
                :class="getTargetChangeClass(event.target_from, event.target_to)"
              >
                ${{ event.target_to }}
                ({{ formatTargetChange(event.target_from, event.target_to) }})
              </span>
            </div>
          </div>

          <!-- Timestamp -->
          <div class="ml-4 flex-shrink-0 text-right">
            <div class="text-sm text-gray-900">{{ formatTime(event.event_time) }}</div>
            <div class="text-xs text-gray-500">{{ formatDate(event.event_time) }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { format } from 'date-fns'
import type { StockEvent } from '@/types'

interface Props {
  events?: StockEvent[]
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
})

defineEmits<{
  viewDetails: [ticker: string]
}>()

// Helper Functions
function formatAction(action: string): string {
  return action.split(' ')[0].toUpperCase()
}

function getActionClass(action: string): string {
  const baseClasses = 'bg-opacity-10'
  const actionType = action.toLowerCase()
  
  if (actionType.includes('upgrade')) {
    return `${baseClasses} bg-emerald-500 text-emerald-700`
  }
  if (actionType.includes('downgrade')) {
    return `${baseClasses} bg-red-500 text-red-700`
  }
  if (actionType.includes('initiat')) {
    return `${baseClasses} bg-blue-500 text-blue-700`
  }
  return `${baseClasses} bg-gray-500 text-gray-700`
}

function getRatingClass(rating: string): string {
  const ratingLower = rating.toLowerCase()
  if (ratingLower.includes('buy')) return 'text-emerald-600'
  if (ratingLower.includes('sell')) return 'text-red-600'
  return 'text-gray-600'
}

function getTargetChangeClass(from: number, to: number): string {
  const change = to - from
  if (change > 0) return 'text-emerald-600 font-medium'
  if (change < 0) return 'text-red-600 font-medium'
  return 'text-gray-600'
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
  return format(new Date(date), 'MMM d, yyyy')
}
</script>