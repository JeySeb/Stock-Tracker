<template>
    <div class="space-y-8">

  
      <!-- Enhanced Stock Heat Map -->
      <StockHeatMap
        :events="recentStockEvents"
        :loading="stocksStore.isLoading"
      />
      <!-- Enhanced Analytics Section -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Enhanced Stock Events (2/3 width) -->
        <div class="lg:col-span-2">
          <!-- Header -->
          <div class="bg-gradient-to-r from-blue-600 to-indigo-700 p-6 text-white">
            <div class="flex flex-col md:flex-row md:items-center md:justify-between">
              <div>
                <h1 class="text-2xl font-bold">Recent Stock Events</h1>
                <p class="opacity-90 mt-1">Tracking of <span class="font-semibold">Stocks Market</span> activity</p>
              </div>
              <div class="flex items-center gap-2">
                <span v-if="stocksStore.isLoading" class="animate-pulse bg-white/20 h-6 w-16 rounded"></span>
                <span v-else class="text-sm font-medium text-white">
                  Total: {{ totalStocksCount?.toLocaleString() ?? '-' }}
                </span>
              </div>
            </div>
          </div>
          
          <div class="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
            <EnhancedStockEvents
              :events="recentStockEvents"
              :loading="stocksStore.isLoading"
              @view-stock="handleViewStock"
            />
          </div>
        </div>

        <!-- Broker Credibility Report (1/3 width) -->
        <div>
          <BrokerCredibilityReport
            :brokers="brokersStore.topBrokersByReportCount"
            :loading="brokersStore.isLoading"
            @view-broker="handleViewBroker"
          />
        </div>
      </div>

      <!-- Enhanced Recommendations Section -->
      <EnhancedRecommendationsCard
        :recommendations="Array.from(recommendationsStore.recommendations || [])"
        :loading="recommendationsStore.isLoading"
        :user-tier="authStore.userTier"
        :rate-limit-remaining="rateLimitRemaining"
        @view-details="handleViewRecommendation"
        @view-all="$router.push('/recommendations')"
        @refresh-recommendations="refreshRecommendations"
      />
  
      <!-- Feature Cards Based on Tier -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <!-- Always Available -->
                <FeatureCard
          title="Stock Explorer"
          description="Browse and filter stock events from major brokerages"
          :icon="MagnifyingGlassIcon"
          :link="{ name: 'stocks' }"
          available
          icon-bg-color="bg-green-100"
          icon-text-color="text-green-600"
        />
        
  
        <FeatureCard
          title="Recommendations"
          description="Get AI-powered stock recommendations based on broker actions"
          :icon="StarIcon"
          :link="{ name: 'recommendations' }"
          available
          icon-bg-color="bg-yellow-100"
          icon-text-color="text-yellow-600"
        />

        <!-- Basic/Premium Features -->
        <FeatureCard
          title="Real-time Data"
          description="Access live market data and external analytics"
          :icon="ChartBarIcon"
          :available="authStore.hasFeature('real_time_data')"
          :link="authStore.hasFeature('real_time_data') ? { name: 'real-time-data' } : undefined"
          @upgrade="$router.push('/subscription')"
          icon-bg-color="bg-blue-100"
          icon-text-color="text-blue-600"
        />
  
      </div>
  
      <!-- Debug Section (Development Only) -->
      <div v-if="isDev" class="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
        <h3 class="text-lg font-semibold text-yellow-900 mb-2">Debug Info</h3>
        <div class="space-y-2 text-sm">
          <div>Auth State: {{ authStore.isAuthenticated ? '✅ Authenticated' : '❌ Not Authenticated' }}</div>
          <div>User Tier: {{ authStore.userTier }}</div>
          <div>Has Real-time Data: {{ authStore.hasFeature('real_time_data') ? '✅ Yes' : '❌ No' }}</div>
          <button 
            @click="testRealTimeDataAccess"
            class="px-3 py-1 bg-blue-600 text-white rounded text-xs hover:bg-blue-700"
          >
            Test Real-Time Data Access
          </button>
        </div>
      </div>

    </div>
  </template>
  
  <script setup lang="ts">
  import { onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useStocksStore } from '@/stores/stocks'
import { useRecommendationsStore } from '@/stores/recommendations'
import { useBrokersStore } from '@/stores/brokers'
import TierBadge from '@/components/ui/TierBadge.vue'
import MetricCard from '@/components/features/dashboard/MetricCard.vue'
import FeatureCard from '@/components/features/dashboard/FeatureCard.vue'
import RecommendationDistributionChart from '@/components/features/dashboard/RecommendationDistributionChart.vue'
import TopRecommendationsList from '@/components/features/dashboard/TopRecommendationsList.vue'
import StockHeatMap from '@/components/features/dashboard/StockHeatMap.vue'
import EnhancedStockEvents from '@/components/features/dashboard/EnhancedStockEvents.vue'
import BrokerCredibilityReport from '@/components/features/dashboard/BrokerCredibilityReport.vue'
import EnhancedRecommendationsCard from '@/components/features/dashboard/EnhancedRecommendationsCard.vue'
import { MagnifyingGlassIcon, StarIcon, ChartBarIcon, CpuChipIcon } from '@heroicons/vue/24/outline'
import type { Recommendation } from '@/types'

  const router = useRouter()
  const authStore = useAuthStore()
  const stocksStore = useStocksStore()
  const recommendationsStore = useRecommendationsStore()
  const brokersStore = useBrokersStore()

  const strongBuyCount = computed(() => {
    if (!recommendationsStore.recommendations || recommendationsStore.recommendations.length === 0) return 0
    return recommendationsStore.recommendations.filter(
      (r: Recommendation) => r.recommendation_type === 'Strong Buy'
    ).length
  })

  const activeRecommendationsCount = computed(() => {
    return recommendationsStore.recommendations?.length || 0
  })

  const totalStocksCount = computed(() => {
    return stocksStore.stats?.total_stocks || 0
  })

  const rateLimitRemaining = computed(() => {
    return recommendationsStore.rateLimitRemaining || 0
  })
  
  const topRecommendations = computed(() => {
    if (!recommendationsStore.recommendations || recommendationsStore.recommendations.length === 0) return []
    return [...recommendationsStore.recommendations]
      .sort((a: Recommendation, b: Recommendation) => b.basic_score - a.basic_score)
      .slice(0, 5)
      .map((rec: Recommendation) => ({
        ...rec,
        scoring_factors: Array.from(rec.scoring_factors)
      }))
  })
  
  const recentStockEvents = computed(() => {
    if (!stocksStore.stocks || stocksStore.stocks.length === 0) return []
    return stocksStore.stocks.slice(0, 60) // Show up to 60 stocks for the heat map
  })
  
  const recommendationChartData = computed(() => {
    if (!recommendationsStore.recommendations || recommendationsStore.recommendations.length === 0) return []
    const distribution = recommendationsStore.recommendations.reduce((acc: Record<string, number>, rec: Recommendation) => {
      acc[rec.recommendation_type] = (acc[rec.recommendation_type] || 0) + 1
      return acc
    }, {} as Record<string, number>)

    return Object.entries(distribution).map(([type, count]) => ({
      name: type,
      value: count
    }))
  })
  
  const isDev = computed(() => import.meta.env.DEV)
  
  onMounted(async () => {
    try {
      await Promise.allSettled([
        stocksStore.fetchStocks(),
        stocksStore.fetchStats(),
        recommendationsStore.fetchRecommendations(),
        brokersStore.fetchBrokerScores()
      ])
    } catch (error) {
      console.error('Error loading dashboard data:', error)
    }
  })
  
  function handleViewRecommendation(ticker: string) {
    router.push({ name: 'recommendations', query: { ticker } })
  }
  
  function handleViewStock(ticker: string) {
    router.push({ name: 'stocks', query: { ticker } })
  }
  
  function handleViewBroker(brokerName: string) {
    router.push({ name: 'stocks', query: { brokerage: brokerName } })
  }

  async function refreshRecommendations() {
    try {
      await recommendationsStore.fetchRecommendations()
    } catch (error) {
      console.error('Error refreshing recommendations:', error)
    }
  }
  
  function testRealTimeDataAccess() {
    console.log('🧪 Testing Real-Time Data access...')
    console.log('Auth state:', {
      isAuthenticated: authStore.isAuthenticated,
      userTier: authStore.userTier,
      hasFeature: authStore.hasFeature('real_time_data')
    })
    
    // Try to navigate to real-time data
    router.push('/real-time-data').then(() => {
      console.log('✅ Successfully navigated to real-time data')
    }).catch((error) => {
      console.error('❌ Failed to navigate to real-time data:', error)
    })
  }
  </script>