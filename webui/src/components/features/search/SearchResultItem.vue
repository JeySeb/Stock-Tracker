<template>
  <div
    class="px-4 py-3 cursor-pointer hover:bg-gray-50 transition-colors duration-150"
    :class="{ 'bg-primary-50 border-l-4 border-primary-500': isHighlighted }"
    @click="$emit('select', result)"
  >
    <div class="flex items-start space-x-3">
      <!-- Icon based on result type -->
      <div class="flex-shrink-0 mt-1">
        <div
          class="w-8 h-8 rounded-full flex items-center justify-center text-white text-sm font-medium"
          :class="{
            'bg-blue-500': result.type === 'stock',
            'bg-green-500': result.type === 'recommendation',
            'bg-purple-500': result.type === 'company'
          }"
        >
          <span v-if="result.type === 'stock'">📈</span>
          <span v-else-if="result.type === 'recommendation'">🎯</span>
          <span v-else-if="result.type === 'company'">🏢</span>
        </div>
      </div>

      <!-- Content -->
      <div class="flex-1 min-w-0">
        <div class="flex items-center space-x-2">
          <h4 class="text-sm font-medium text-gray-900 truncate">
            {{ result.title }}
          </h4>
          <span
            class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
            :class="{
              'bg-blue-100 text-blue-800': result.type === 'stock',
              'bg-green-100 text-green-800': result.type === 'recommendation',
              'bg-purple-100 text-purple-800': result.type === 'company'
            }"
          >
            {{ result.type }}
          </span>
        </div>
        
        <p class="text-sm text-gray-600 truncate">
          {{ result.subtitle }}
        </p>
        
        <p class="text-xs text-gray-500 mt-1">
          {{ result.description }}
        </p>

        <!-- Additional info for recommendations -->
        <div v-if="result.type === 'recommendation'" class="mt-2 flex items-center space-x-4 text-xs text-gray-500">
          <span class="flex items-center">
            <span class="w-2 h-2 rounded-full mr-1" :class="getScoreColor((result.data as Recommendation).basic_score)"></span>
            Score: {{ ((result.data as Recommendation).basic_score * 100).toFixed(0) }}%
          </span>
          <span v-if="(result.data as Recommendation).confidence">
            Confidence: {{ ((result.data as Recommendation).confidence * 100).toFixed(0) }}%
          </span>
        </div>

        <!-- Additional info for stocks -->
        <div v-if="result.type === 'stock'" class="mt-2 flex items-center space-x-4 text-xs text-gray-500">
          <span>{{ (result.data as StockEvent).brokerage }}</span>
          <span>{{ (result.data as StockEvent).action }}</span>
          <span>{{ formatDate((result.data as StockEvent).event_time) }}</span>
        </div>
      </div>

      <!-- Relevance indicator -->
      <div class="flex-shrink-0">
        <div
          class="w-2 h-2 rounded-full"
          :class="getRelevanceColor(result.relevance)"
          :title="`Relevance: ${result.relevance}`"
        ></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { SearchResult } from '@/composables/useGlobalSearch'
import type { StockEvent, Recommendation } from '@/types'

interface Props {
  result: SearchResult
  isHighlighted: boolean
}

defineProps<Props>()
defineEmits<{
  select: [result: SearchResult]
}>()

const getScoreColor = (score: number) => {
  if (score >= 0.8) return 'bg-green-500'
  if (score >= 0.6) return 'bg-yellow-500'
  return 'bg-red-500'
}

const getRelevanceColor = (relevance: number) => {
  if (relevance >= 80) return 'bg-green-500'
  if (relevance >= 60) return 'bg-yellow-500'
  return 'bg-gray-400'
}

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffTime = Math.abs(now.getTime() - date.getTime())
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
  
  if (diffDays === 1) return 'Today'
  if (diffDays === 2) return 'Yesterday'
  if (diffDays <= 7) return `${diffDays - 1} days ago`
  return date.toLocaleDateString()
}
</script>