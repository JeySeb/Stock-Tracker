<template>
    <div class="space-y-8">
      <!-- Welcome Header -->
      <div class="bg-gradient-to-r from-primary-600 to-primary-700 rounded-lg p-6 text-white">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold">
              Welcome{{ authStore.user ? `, ${authStore.user.first_name}` : ' to Stock Tracker' }}
            </h1>
            <p class="mt-2 text-primary-100">
              Your {{ authStore.userTier }} dashboard with market insights and recommendations
            </p>
          </div>
          <TierBadge :tier="authStore.userTier" class="bg-white/20 text-white border-white/30" />
        </div>
      </div>
  
      <!-- Key Metrics -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <MetricCard
          title="Total Stock Events"
          :value="totalStocksCount"
          :loading="stocksStore.isLoading"
          icon="📊"
          color="blue"
        />
        <MetricCard
          title="Active Recommendations"
          :value="activeRecommendationsCount"
          :loading="recommendationsStore.isLoading"
          icon="⭐"
          color="green"
        />
        <MetricCard
          title="Strong Buy Signals"
          :value="strongBuyCount"
          :loading="recommendationsStore.isLoading"
          icon="📈"
          color="emerald"
        />
        <MetricCard
          title="Rate Limit Remaining"
          :value="rateLimitRemaining"
          :loading="recommendationsStore.isLoading"
          icon="⚡"
          color="amber"
          :format="(v: number) => v ? `${v}/hr` : '-'"
        />
      </div>
  
      <!-- Charts Section -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <!-- Recommendations Overview Chart -->
        <div class="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Recommendation Distribution</h3>
          <RecommendationDistributionChart
            :data="recommendationChartData"
            :loading="recommendationsStore.isLoading"
          />
        </div>
  
        <!-- Top Performers -->
        <div class="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Top Recommendations</h3>
          <TopRecommendationsList
            :recommendations="topRecommendations"
            :loading="recommendationsStore.isLoading"
            :user-tier="authStore.userTier"
            @view-details="handleViewRecommendation"
          />
        </div>
      </div>
  
      <!-- Feature Cards Based on Tier -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <!-- Always Available -->
        <FeatureCard
          title="Stock Explorer"
          description="Browse and filter stock events from major brokerages"
          icon="🔍"
          :link="{ name: 'stocks' }"
          available
        />
        
        <FeatureCard
          title="Real-Time Data"
          description="Access live market data and external analytics"
          icon="📊"
          :link="{ name: 'real-time-data' }"
          available
        />
  
        <FeatureCard
          title="Recommendations"
          description="Get AI-powered stock recommendations based on broker actions"
          icon="⭐"
          :link="{ name: 'recommendations' }"
          available
        />
  
        <!-- Basic/Premium Features -->
        <FeatureCard
          title="Real-time Data"
          description="Access live market data and external analytics"
          icon="📊"
          :available="authStore.hasFeature('real_time_data')"
          @upgrade="$router.push('/subscription')"
        />
  
        <!-- Premium Only -->
        <FeatureCard
          title="AI Insights"
          description="Advanced sentiment analysis and market predictions"
          icon="🤖"
          :available="authStore.hasFeature('ai_insights')"
          premium-only
          @upgrade="$router.push('/subscription')"
        />
  
        <FeatureCard
          title="Sentiment Analysis"
          description="Track news and social media sentiment for stocks"
          icon="💭"
          :available="authStore.hasFeature('sentiment_analysis')"
          premium-only
          @upgrade="$router.push('/subscription')"
        />
  
        <!-- AI Chat Placeholder -->
        <FeatureCard
          title="AI Assistant"
          description="Chat with our AI for personalized investment insights"
          icon="💬"
          :available="false"
          coming-soon
          placeholder-for-ai
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

      <!-- Recent Activity -->
      <div class="bg-white rounded-lg shadow-sm border border-gray-200">
        <div class="px-6 py-4 border-b border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900">Recent Stock Events</h3>
        </div>
        <RecentStockEvents
          :events="recentStockEvents"
          :loading="stocksStore.isLoading"
          @view-details="handleViewStock"
        />
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { onMounted, computed } from 'vue'
  import { useRouter } from 'vue-router'
  import { useAuthStore } from '@/stores/auth'
  import { useStocksStore } from '@/stores/stocks'
  import { useRecommendationsStore } from '@/stores/recommendations'
  import TierBadge from '@/components/ui/TierBadge.vue'
  import MetricCard from '@/components/features/dashboard/MetricCard.vue'
  import FeatureCard from '@/components/features/dashboard/FeatureCard.vue'
  import RecommendationDistributionChart from '@/components/features/dashboard/RecommendationDistributionChart.vue'
  import TopRecommendationsList from '@/components/features/dashboard/TopRecommendationsList.vue'
  import RecentStockEvents from '@/components/features/dashboard/RecentStockEvents.vue'
  import type { Recommendation } from '@/types'

  const router = useRouter()
  const authStore = useAuthStore()
  const stocksStore = useStocksStore()
  const recommendationsStore = useRecommendationsStore()

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
    return recommendationsStore.recommendations
      .filter((r: Recommendation) => ['Strong Buy', 'Buy'].includes(r.recommendation_type))
      .sort((a: Recommendation, b: Recommendation) => b.basic_score - a.basic_score)
      .slice(0, 5)
      .map(rec => ({
        ...rec,
        scoring_factors: [...rec.scoring_factors]
      }))
  })
  
  const recentStockEvents = computed(() => {
    if (!stocksStore.stocks || stocksStore.stocks.length === 0) return []
    return stocksStore.stocks.slice(0, 10)
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
        recommendationsStore.fetchRecommendations()
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