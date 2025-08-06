<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <h4 class="text-md font-medium text-gray-900">Most Reliable Brokers</h4>
      <span class="text-xs text-gray-500">Based on report volume & performance</span>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="space-y-3">
      <div v-for="n in 5" :key="n" class="animate-pulse flex items-center justify-between p-3 border border-gray-200 rounded-lg">
        <div class="flex items-center space-x-3">
          <div class="h-8 w-8 bg-gray-200 rounded"></div>
          <div class="space-y-1">
            <div class="h-4 bg-gray-200 rounded w-24"></div>
            <div class="h-3 bg-gray-200 rounded w-16"></div>
          </div>
        </div>
        <div class="h-6 w-16 bg-gray-200 rounded"></div>
      </div>
    </div>

    <!-- Empty State -->
    <div 
      v-else-if="!topBrokers || topBrokers.length === 0" 
      class="text-center text-gray-500 py-6"
    >
      No broker data available
    </div>

    <!-- Brokers List -->
    <div v-else class="space-y-3">
      <div
        v-for="(broker, index) in topBrokers"
        :key="broker.id"
        class="flex items-center justify-between p-3 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors cursor-pointer"
        @click="$emit('viewBroker', broker.name)"
      >
        <!-- Left Side: Broker Info -->
        <div class="flex items-center space-x-3">
          <!-- Ranking Badge -->
          <div 
            class="flex items-center justify-center w-8 h-8 rounded-full text-sm font-semibold"
            :class="getRankingBadgeClass(index + 1)"
          >
            {{ index + 1 }}
          </div>
          
          <!-- Broker Details -->
          <div>
            <div class="font-medium text-gray-900">{{ broker.name }}</div>
            <div class="text-sm text-gray-500">
              {{ broker.report_count.toLocaleString() }} reports
            </div>
          </div>
        </div>

        <!-- Right Side: Score and Visual -->
        <div class="flex items-center space-x-3">
          <!-- Score Visualization -->
          <div class="text-right">
            <div class="text-sm font-medium text-gray-900">
              {{ formatScore(broker.calculated_score) }}
            </div>
            <div class="text-xs text-gray-500">Score</div>
          </div>
          
          <!-- Score Bar -->
          <div class="w-16 bg-gray-200 rounded-full h-2">
            <div 
              class="h-2 rounded-full transition-all duration-300"
              :class="getScoreBarClass(broker.calculated_score)"
              :style="{ width: `${(broker.calculated_score * 100)}%` }"
            ></div>
          </div>

          <!-- Trend Arrow -->
          <ChevronRightIcon class="h-4 w-4 text-gray-400" />
        </div>
      </div>
    </div>

    <!-- Summary Stats -->
    <div v-if="!loading && topBrokers.length > 0" class="mt-4 pt-4 border-t border-gray-200">
      <div class="grid grid-cols-3 gap-4 text-center">
        <div>
          <div class="text-lg font-semibold text-gray-900">{{ totalReports.toLocaleString() }}</div>
          <div class="text-xs text-gray-500">Total Reports</div>
        </div>
        <div>
          <div class="text-lg font-semibold text-gray-900">{{ averageScore.toFixed(2) }}</div>
          <div class="text-xs text-gray-500">Avg Score</div>
        </div>
        <div>
          <div class="text-lg font-semibold text-gray-900">{{ topBrokers.length }}</div>
          <div class="text-xs text-gray-500">Active Brokers</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { ChevronRightIcon } from '@heroicons/vue/24/outline'
import type { BrokerScore } from '@/api/brokers'

interface Props {
  brokers?: BrokerScore[]
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
})

defineEmits<{
  viewBroker: [brokerName: string]
}>()

// Computed properties
const topBrokers = computed(() => {
  if (!props.brokers || props.brokers.length === 0) return []
  return props.brokers.slice(0, 10) // Show top 5 brokers
})

const totalReports = computed(() => {
  if (!props.brokers || props.brokers.length === 0) return 0
  return props.brokers.reduce((sum, broker) => sum + broker.report_count, 0)
})

const averageScore = computed(() => {
  if (!props.brokers || props.brokers.length === 0) return 0
  const sum = props.brokers.reduce((sum, broker) => sum + broker.calculated_score, 0)
  return sum / props.brokers.length
})

// Helper Functions
function getRankingBadgeClass(rank: number): string {
  switch (rank) {
    case 1:
      return 'bg-yellow-100 text-yellow-800'
    case 2:
      return 'bg-gray-100 text-gray-800'
    case 3:
      return 'bg-amber-100 text-amber-800'
    default:
      return 'bg-blue-100 text-blue-800'
  }
}

function formatScore(score: number): string {
  return (score * 100).toFixed(1) + '%'
}

function getScoreBarClass(score: number): string {
  if (score >= 0.7) return 'bg-emerald-500'
  if (score >= 0.5) return 'bg-yellow-500'
  return 'bg-red-500'
}
</script>