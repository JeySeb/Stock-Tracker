<template>
    <div class="bg-white p-6 rounded-lg shadow-sm border border-gray-200 hover:shadow-md transition-shadow">
      <!-- Header -->
      <div class="flex items-start justify-between mb-4">
        <div>
          <h3 class="text-lg font-semibold text-gray-900">{{ recommendation.ticker }}</h3>
          <p class="text-sm text-gray-600 truncate">{{ recommendation.company_name }}</p>
        </div>
        <div class="flex flex-col items-end">
          <RecommendationBadge :type="recommendation.recommendation_type" />
          <span class="text-xs text-gray-500 mt-1">
            Score: {{ (recommendation.basic_score * 100).toFixed(0) }}%
          </span>
        </div>
      </div>
  
      <!-- Key Metrics -->
      <div class="grid grid-cols-2 gap-4 mb-4">
        <div>
          <p class="text-xs text-gray-500">Target Price</p>
          <p class="text-sm font-medium">${{ recommendation.latest_target_price.toFixed(2) }}</p>
        </div>
        <div>
          <p class="text-xs text-gray-500">Avg Change</p>
          <p :class="getChangeClass(recommendation.avg_target_change)">
            {{ formatChange(recommendation.avg_target_change) }}%
          </p>
        </div>
        <div>
          <p class="text-xs text-gray-500">Total Events</p>
          <p class="text-sm font-medium">{{ recommendation.total_events }}</p>
        </div>
        <div>
          <p class="text-xs text-gray-500">Confidence</p>
          <p class="text-sm font-medium">{{ (recommendation.confidence * 100).toFixed(0) }}%</p>
        </div>
      </div>
  
      <!-- External Data (Basic/Premium only) -->
      <div v-if="recommendation.external_data && userTier !== 'guest'" class="border-t pt-4 mb-4">
        <p class="text-xs text-gray-500 mb-2">Market Data</p>
        <div class="grid grid-cols-2 gap-2 text-xs">
          <div>
            <span class="text-gray-500">Price:</span>
            <span class="font-medium ml-1">${{ recommendation.external_data.current_price.toFixed(2) }}</span>
          </div>
          <div>
            <span class="text-gray-500">Change:</span>
            <span :class="getChangeClass(recommendation.external_data.price_change_24h)" class="ml-1">
              {{ formatChange(recommendation.external_data.price_change_24h) }}%
            </span>
          </div>
        </div>
      </div>
  
      <!-- AI Insights Preview (Premium only) -->
      <div v-if="recommendation.ai_insights && userTier === 'premium'" class="border-t pt-4 mb-4">
        <p class="text-xs text-gray-500 mb-2">AI Insights</p>
        <div class="text-xs space-y-1">
          <div>
            <span class="text-gray-500">Sentiment:</span>
            <span class="font-medium ml-1 capitalize">{{ recommendation.ai_insights.news_sentiment }}</span>
          </div>
          <div>
            <span class="text-gray-500">Risk:</span>
            <span class="font-medium ml-1 capitalize">{{ recommendation.ai_insights.risk_assessment }}</span>
          </div>
        </div>
      </div>
  
      <!-- Upgrade Prompt for Guest/Basic -->
      <div v-if="userTier !== 'premium'" class="border-t pt-4 mb-4">
        <div class="text-center text-xs text-gray-500">
          <p v-if="userTier === 'guest'">
            📊 Sign up for real-time data and detailed analytics
          </p>
          <p v-else-if="userTier === 'basic'">
            🤖 Upgrade to Premium for AI insights and predictions
          </p>
        </div>
      </div>
  
      <!-- Actions -->
      <div class="flex space-x-2">
        <button
          @click="$emit('viewDetails', recommendation)"
          class="flex-1 bg-primary-600 text-white text-sm font-medium py-2 px-4 rounded-md hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500"
        >
          View Details
        </button>
        <button
          v-if="userTier !== 'premium'"
          @click="$emit('upgrade')"
          class="px-4 py-2 border border-primary-600 text-primary-600 text-sm font-medium rounded-md hover:bg-primary-50 focus:outline-none focus:ring-2 focus:ring-primary-500"
        >
          Upgrade
        </button>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import RecommendationBadge from './RecommendationBadge.vue'
  import type { Recommendation, UserTier } from '@/types'
  
  interface Props {
    recommendation: Recommendation
    userTier: UserTier
  }
  
  defineProps<Props>()
  
  defineEmits<{
    viewDetails: [recommendation: Recommendation]
    upgrade: []
  }>()
  
  function getChangeClass(change: number): string {
    if (change > 0) return 'text-financial-buy font-medium'
    if (change < 0) return 'text-financial-sell font-medium'
    return 'text-gray-600 font-medium'
  }
  
  function formatChange(change: number): string {
    const formatted = Math.abs(change).toFixed(1)
    return change > 0 ? `+${formatted}` : `-${formatted}`
  }
  </script>