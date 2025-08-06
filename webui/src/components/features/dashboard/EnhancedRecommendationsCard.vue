<template>
  <div class="bg-white p-6 rounded-lg shadow-sm border border-gray-200 space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <h3 class="text-lg font-semibold text-gray-900">Your Active Recommendations</h3>
      <div class="flex items-center gap-3">
        <!-- Rate Limit Indicator -->
        <div v-if="rateLimitRemaining !== undefined" class="flex items-center space-x-1">
          <div class="w-2 h-2 rounded-full" :class="rateLimitRemaining > 10 ? 'bg-emerald-500' : rateLimitRemaining > 5 ? 'bg-yellow-500' : 'bg-red-500'"></div>
          <span class="text-xs text-gray-500">{{ rateLimitRemaining }} left</span>
        </div>
        
        <!-- Total Count -->
        <span v-if="!loading" class="text-sm font-medium text-gray-600">
          {{ activeCount?.toLocaleString() ?? '-' }} Active
        </span>
        <span v-else class="animate-pulse bg-gray-200 h-6 w-16 rounded"></span>
      </div>
    </div>

    <!-- Main Content Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Left: Distribution Chart -->
      <div class="lg:col-span-2">
        <div class="flex items-center justify-between mb-4">
          <h4 class="text-sm font-medium text-gray-700">Recommendation Distribution</h4>
          <button 
            @click="$emit('refreshRecommendations')"
            class="text-xs text-blue-600 hover:text-blue-800 transition-colors"
            :disabled="loading"
          >
            {{ loading ? 'Updating...' : 'Refresh' }}
          </button>
        </div>
        
        <RecommendationDistributionChart
          :data="chartData"
          :loading="loading"
        />
        
        <!-- Quick Stats -->
        <div class="mt-4 grid grid-cols-3 gap-4">
          <div class="text-center">
            <div class="text-lg font-semibold text-emerald-600">{{ strongBuyCount }}</div>
            <div class="text-xs text-gray-500">Strong Buy</div>
          </div>
          <div class="text-center">
            <div class="text-lg font-semibold text-blue-600">{{ buyCount }}</div>
            <div class="text-xs text-gray-500">Buy</div>
          </div>
          <div class="text-center">
            <div class="text-lg font-semibold text-gray-600">{{ holdCount }}</div>
            <div class="text-xs text-gray-500">Hold</div>
          </div>
        </div>
      </div>

      <!-- Right: Key Metrics -->
      <div class="space-y-4">
        <!-- Active Recommendations Metric -->
        <MetricCard
          title="Active Recommendations"
          :value="activeCount"
          :loading="loading"
          icon="🧠"
          color="emerald"
          :trend="recommendationTrend"
          trend-label="vs last week"
          :size-factor="0.9"
        />
        
        <!-- Average Score Metric -->
        <MetricCard
          title="Average Confidence"
          :value="averageConfidence"
          :loading="loading"
          icon="📊"
          color="blue"
          :format="(v: number) => `${(v * 100).toFixed(1)}%`"
          :size-factor="0.9"
        />
        
        <!-- Top Performer Indicator -->
        <div v-if="topPerformer && !loading" class="bg-gradient-to-r from-emerald-50 to-blue-50 p-4 rounded-lg border border-emerald-200">
          <div class="text-sm font-medium text-emerald-800 mb-1">🎯 Top Pick</div>
          <div class="font-semibold text-gray-900">{{ topPerformer.ticker }}</div>
          <div class="text-sm text-gray-600">{{ topPerformer.company_name }}</div>
          <div class="text-xs text-emerald-700 mt-1">
            Score: {{ (topPerformer.basic_score * 100).toFixed(1) }}%
          </div>
        </div>
      </div>
    </div>

    <!-- Top Recommendations List -->
    <div class="border-t border-gray-200 pt-6">
      <div class="flex items-center justify-between mb-4">
        <h4 class="text-sm font-medium text-gray-700">Top Recommendations</h4>
        <button 
          @click="$emit('viewAll')"
          class="text-xs text-blue-600 hover:text-blue-800 transition-colors"
        >
          View All →
        </button>
      </div>
      
      <TopRecommendationsList
        :recommendations="topRecommendations"
        :loading="loading"
        :user-tier="userTier"
        @view-details="$emit('viewDetails', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import MetricCard from './MetricCard.vue'
import RecommendationDistributionChart from './RecommendationDistributionChart.vue'
import TopRecommendationsList from './TopRecommendationsList.vue'
import type { Recommendation, UserTier } from '@/types'

interface Props {
  recommendations?: Recommendation[]
  loading?: boolean
  userTier: UserTier
  rateLimitRemaining?: number
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
})

defineEmits<{
  viewDetails: [ticker: string]
  viewAll: []
  refreshRecommendations: []
}>()

// Computed properties
const activeCount = computed(() => {
  return props.recommendations?.length || 0
})

const averageConfidence = computed(() => {
  if (!props.recommendations || props.recommendations.length === 0) return 0
  const sum = props.recommendations.reduce((acc, rec) => acc + rec.confidence, 0)
  return sum / props.recommendations.length
})

const strongBuyCount = computed(() => {
  if (!props.recommendations || props.recommendations.length === 0) return 0
  return props.recommendations.filter(r => r.recommendation_type === 'Strong Buy').length
})

const buyCount = computed(() => {
  if (!props.recommendations || props.recommendations.length === 0) return 0
  return props.recommendations.filter(r => r.recommendation_type === 'Buy').length
})

const holdCount = computed(() => {
  if (!props.recommendations || props.recommendations.length === 0) return 0
  return props.recommendations.filter(r => r.recommendation_type === 'Hold').length
})

const topPerformer = computed(() => {
  if (!props.recommendations || props.recommendations.length === 0) return null
  return [...props.recommendations]
    .sort((a, b) => b.basic_score - a.basic_score)[0]
})

const topRecommendations = computed(() => {
  if (!props.recommendations || props.recommendations.length === 0) return []
  return [...props.recommendations]
    .sort((a, b) => b.basic_score - a.basic_score)
    .slice(0, 3)
    .map(rec => ({
      ...rec,
      scoring_factors: Array.from(rec.scoring_factors)
    }))
})

const chartData = computed(() => {
  if (!props.recommendations || props.recommendations.length === 0) return []
  const distribution = props.recommendations.reduce((acc: Record<string, number>, rec) => {
    acc[rec.recommendation_type] = (acc[rec.recommendation_type] || 0) + 1
    return acc
  }, {} as Record<string, number>)

  return Object.entries(distribution).map(([type, count]) => ({
    name: type,
    value: count
  }))
})

// Mock trend data - in real app this would come from historical data
const recommendationTrend = computed(() => {
  // This would typically be calculated from historical data
  return Math.random() > 0.5 ? Math.floor(Math.random() * 20) : -Math.floor(Math.random() * 10)
})
</script>