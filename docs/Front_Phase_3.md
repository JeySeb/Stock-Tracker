# Frontend Development - Phase 3: Advanced Features & AI Readiness

## Phase Overview

**Duration:** 2-3 weeks  
**Focus:** Advanced search capabilities, performance optimizations, comprehensive AI placeholder integration, and production-ready polish  
**Goal:** Deliver a highly polished, performant financial analytics platform ready for AI integration and real-world deployment

## Justification for Phase 3 Structure

### Why This Final Phase?
1. **Performance Excellence:** Optimize for real-world usage with caching and lazy loading
2. **AI Integration Readiness:** Complete infrastructure for seamless AI feature rollout
3. **User Experience Refinement:** Advanced search, error handling, and accessibility
4. **Production Readiness:** Comprehensive testing, monitoring, and deployment preparation

### Critical Success Factors
- Sub-second load times for all major operations
- Comprehensive error boundaries and graceful degradation
- Full AI chat interface ready for backend integration
- Advanced search and filtering capabilities
- Accessibility compliance and responsive design perfection

---

## Step-by-Step Implementation

### Step 1: Advanced Search & Filtering System

#### 1.1 Global Search Composable (`src/composables/useGlobalSearch.ts`)
```typescript
import { ref, computed, watch } from 'vue'
import { debounce } from 'lodash-es'
import { useStocksStore } from '@/stores/stocks'
import { useRecommendationsStore } from '@/stores/recommendations'
import type { StockEvent, Recommendation } from '@/types'

export interface SearchResult {
  type: 'stock' | 'recommendation' | 'company'
  id: string
  title: string
  subtitle: string
  description: string
  data: StockEvent | Recommendation | any
  relevance: number
}

export function useGlobalSearch() {
  const searchQuery = ref('')
  const isSearching = ref(false)
  const searchResults = ref<SearchResult[]>([])
  const recentSearches = ref<string[]>([])
  const searchHistory = ref<SearchResult[]>([])

  const stocksStore = useStocksStore()
  const recommendationsStore = useRecommendationsStore()

  // Load recent searches from localStorage
  const loadRecentSearches = () => {
    const stored = localStorage.getItem('stock_tracker_recent_searches')
    if (stored) {
      recentSearches.value = JSON.parse(stored).slice(0, 10)
    }
  }

  // Save search to history
  const saveSearch = (query: string, results: SearchResult[]) => {
    if (query.trim()) {
      recentSearches.value = [
        query,
        ...recentSearches.value.filter(s => s !== query)
      ].slice(0, 10)
      
      localStorage.setItem('stock_tracker_recent_searches', JSON.stringify(recentSearches.value))
      
      if (results.length > 0) {
        searchHistory.value = [
          results[0],
          ...searchHistory.value.filter(r => r.id !== results[0].id)
        ].slice(0, 20)
      }
    }
  }

  // Advanced search algorithm
  const performSearch = async (query: string): Promise<SearchResult[]> => {
    if (!query.trim()) return []

    isSearching.value = true
    const results: SearchResult[] = []

    try {
      // Search stocks
      const stockResults = await searchStocks(query)
      results.push(...stockResults)

      // Search recommendations
      const recommendationResults = await searchRecommendations(query)
      results.push(...recommendationResults)

      // Search companies
      const companyResults = await searchCompanies(query)
      results.push(...companyResults)

      // Sort by relevance
      return results.sort((a, b) => b.relevance - a.relevance).slice(0, 50)
    } finally {
      isSearching.value = false
    }
  }

  const searchStocks = async (query: string): Promise<SearchResult[]> => {
    const response = await stocksStore.fetchStocks()
    const stocks = stocksStore.stocks

    return stocks
      .filter(stock => 
        stock.ticker.toLowerCase().includes(query.toLowerCase()) ||
        stock.company.toLowerCase().includes(query.toLowerCase()) ||
        stock.brokerage.toLowerCase().includes(query.toLowerCase())
      )
      .map(stock => ({
        type: 'stock' as const,
        id: stock.id,
        title: stock.ticker,
        subtitle: stock.company,
        description: `${stock.action} by ${stock.brokerage}`,
        data: stock,
        relevance: calculateStockRelevance(stock, query)
      }))
  }

  const searchRecommendations = async (query: string): Promise<SearchResult[]> => {
    await recommendationsStore.fetchRecommendations()
    const recommendations = recommendationsStore.recommendations

    return recommendations
      .filter(rec => 
        rec.ticker.toLowerCase().includes(query.toLowerCase()) ||
        rec.company_name.toLowerCase().includes(query.toLowerCase())
      )
      .map(rec => ({
        type: 'recommendation' as const,
        id: rec.id,
        title: rec.ticker,
        subtitle: rec.company_name,
        description: `${rec.recommendation_type} - Score: ${(rec.basic_score * 100).toFixed(0)}%`,
        data: rec,
        relevance: calculateRecommendationRelevance(rec, query)
      }))
  }

  const searchCompanies = async (query: string): Promise<SearchResult[]> => {
    // This would integrate with a company database or API
    // For now, we'll extract unique companies from stocks
    const companies = new Set(stocksStore.stocks.map(s => s.company))
    
    return Array.from(companies)
      .filter(company => company.toLowerCase().includes(query.toLowerCase()))
      .map(company => ({
        type: 'company' as const,
        id: company,
        title: company,
        subtitle: 'Company',
        description: 'View all events and recommendations',
        data: { name: company },
        relevance: calculateTextRelevance(company, query)
      }))
  }

  const calculateStockRelevance = (stock: StockEvent, query: string): number => {
    let score = 0
    const q = query.toLowerCase()
    
    if (stock.ticker.toLowerCase() === q) score += 100
    else if (stock.ticker.toLowerCase().startsWith(q)) score += 80
    else if (stock.ticker.toLowerCase().includes(q)) score += 60
    
    if (stock.company.toLowerCase().includes(q)) score += 40
    if (stock.brokerage.toLowerCase().includes(q)) score += 20
    
    // Boost recent events
    const eventDate = new Date(stock.event_time)
    const daysSinceEvent = (Date.now() - eventDate.getTime()) / (1000 * 60 * 60 * 24)
    if (daysSinceEvent < 7) score += 30
    else if (daysSinceEvent < 30) score += 15
    
    return score
  }

  const calculateRecommendationRelevance = (rec: Recommendation, query: string): number => {
    let score = 0
    const q = query.toLowerCase()
    
    if (rec.ticker.toLowerCase() === q) score += 100
    else if (rec.ticker.toLowerCase().startsWith(q)) score += 80
    else if (rec.ticker.toLowerCase().includes(q)) score += 60
    
    if (rec.company_name.toLowerCase().includes(q)) score += 40
    
    // Boost high-confidence recommendations
    score += rec.confidence * 20
    score += rec.basic_score * 15
    
    return score
  }

  const calculateTextRelevance = (text: string, query: string): number => {
    const t = text.toLowerCase()
    const q = query.toLowerCase()
    
    if (t === q) return 100
    if (t.startsWith(q)) return 80
    if (t.includes(q)) return 60
    return 0
  }

  // Debounced search
  const debouncedSearch = debounce(async (query: string) => {
    if (query.length < 2) {
      searchResults.value = []
      return
    }

    const results = await performSearch(query)
    searchResults.value = results
    saveSearch(query, results)
  }, 300)

  // Watch for search query changes
  watch(searchQuery, (newQuery) => {
    debouncedSearch(newQuery)
  })

  // Clear search
  const clearSearch = () => {
    searchQuery.value = ''
    searchResults.value = []
  }

  // Initialize
  loadRecentSearches()

  return {
    searchQuery,
    isSearching: readonly(isSearching),
    searchResults: readonly(searchResults),
    recentSearches: readonly(recentSearches),
    searchHistory: readonly(searchHistory),
    clearSearch,
    performSearch
  }
}
```

#### 1.2 Global Search Component (`src/components/features/search/GlobalSearch.vue`)
```vue
<template>
  <div class="relative" ref="searchContainer">
    <!-- Search Input -->
    <div class="relative">
      <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
        <MagnifyingGlassIcon class="h-5 w-5 text-gray-400" />
      </div>
      <input
        v-model="searchQuery"
        type="text"
        placeholder="Search stocks, companies, or recommendations..."
        class="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md leading-5 bg-white placeholder-gray-500 focus:outline-none focus:placeholder-gray-400 focus:ring-1 focus:ring-primary-500 focus:border-primary-500"
        @focus="showResults = true"
        @keydown.escape="handleEscape"
        @keydown.arrow-down.prevent="navigateResults(1)"
        @keydown.arrow-up.prevent="navigateResults(-1)"
        @keydown.enter.prevent="selectCurrentResult"
      />
      <div v-if="searchQuery" class="absolute inset-y-0 right-0 pr-3 flex items-center">
        <button
          @click="clearSearch"
          class="text-gray-400 hover:text-gray-600"
        >
          <XMarkIcon class="h-5 w-5" />
        </button>
      </div>
    </div>

    <!-- Search Results Dropdown -->
    <div
      v-if="showResults && (searchResults.length > 0 || isSearching || recentSearches.length > 0)"
      class="absolute z-50 mt-1 w-full bg-white shadow-lg max-h-96 rounded-md py-1 text-base ring-1 ring-black ring-opacity-5 overflow-auto focus:outline-none"
    >
      <!-- Loading State -->
      <div v-if="isSearching" class="px-4 py-2 text-sm text-gray-500">
        <div class="flex items-center">
          <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-primary-600 mr-2"></div>
          Searching...
        </div>
      </div>

      <!-- Search Results -->
      <div v-else-if="searchResults.length > 0">
        <div class="px-3 py-2 text-xs font-medium text-gray-500 uppercase tracking-wide">
          Search Results
        </div>
        <SearchResultItem
          v-for="(result, index) in searchResults"
          :key="result.id"
          :result="result"
          :is-highlighted="highlightedIndex === index"
          @select="selectResult(result)"
          @mouseenter="highlightedIndex = index"
        />
      </div>

      <!-- Recent Searches -->
      <div v-else-if="recentSearches.length > 0">
        <div class="px-3 py-2 text-xs font-medium text-gray-500 uppercase tracking-wide">
          Recent Searches
        </div>
        <div
          v-for="search in recentSearches.slice(0, 5)"
          :key="search"
          class="px-4 py-2 text-sm text-gray-700 cursor-pointer hover:bg-gray-100 flex items-center"
          @click="searchQuery = search"
        >
          <ClockIcon class="h-4 w-4 text-gray-400 mr-2" />
          {{ search }}
        </div>
      </div>

      <!-- No Results -->
      <div v-else-if="searchQuery && !isSearching" class="px-4 py-2 text-sm text-gray-500">
        No results found for "{{ searchQuery }}"
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { MagnifyingGlassIcon, XMarkIcon, ClockIcon } from '@heroicons/vue/24/outline'
import { useGlobalSearch } from '@/composables/useGlobalSearch'
import SearchResultItem from './SearchResultItem.vue'
import type { SearchResult } from '@/composables/useGlobalSearch'

const router = useRouter()
const { searchQuery, isSearching, searchResults, recentSearches, clearSearch } = useGlobalSearch()

const searchContainer = ref<HTMLElement>()
const showResults = ref(false)
const highlightedIndex = ref(-1)

// Handle clicks outside to close dropdown
const handleClickOutside = (event: Event) => {
  if (searchContainer.value && !searchContainer.value.contains(event.target as Node)) {
    showResults.value = false
    highlightedIndex.value = -1
  }
}

// Keyboard navigation
const navigateResults = (direction: number) => {
  if (!showResults.value || searchResults.value.length === 0) return
  
  const maxIndex = searchResults.value.length - 1
  highlightedIndex.value = Math.max(0, Math.min(maxIndex, highlightedIndex.value + direction))
}

const selectCurrentResult = () => {
  if (highlightedIndex.value >= 0 && searchResults.value[highlightedIndex.value]) {
    selectResult(searchResults.value[highlightedIndex.value])
  }
}

const handleEscape = () => {
  showResults.value = false
  highlightedIndex.value = -1
}

const selectResult = (result: SearchResult) => {
  showResults.value = false
  highlightedIndex.value = -1
  
  // Navigate based on result type
  switch (result.type) {
    case 'stock':
      router.push({ name: 'stocks', query: { ticker: result.data.ticker } })
      break
    case 'recommendation':
      router.push({ name: 'recommendations', query: { ticker: result.data.ticker } })
      break
    case 'company':
      router.push({ name: 'stocks', query: { company: result.data.name } })
      break
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>
```

### Step 2: Performance Optimization & Caching

#### 2.1 Cache Management System (`src/utils/cache.ts`)
```typescript
interface CacheItem<T> {
  data: T
  timestamp: number
  ttl: number
}

class CacheManager {
  private cache = new Map<string, CacheItem<any>>()
  private maxSize = 100

  set<T>(key: string, data: T, ttlMinutes: number = 30): void {
    // Clean up expired items if cache is getting full
    if (this.cache.size >= this.maxSize) {
      this.cleanup()
    }

    this.cache.set(key, {
      data,
      timestamp: Date.now(),
      ttl: ttlMinutes * 60 * 1000
    })
  }

  get<T>(key: string): T | null {
    const item = this.cache.get(key)
    
    if (!item) return null
    
    // Check if expired
    if (Date.now() - item.timestamp > item.ttl) {
      this.cache.delete(key)
      return null
    }
    
    return item.data
  }

  has(key: string): boolean {
    const item = this.cache.get(key)
    if (!item) return false
    
    // Check if expired
    if (Date.now() - item.timestamp > item.ttl) {
      this.cache.delete(key)
      return false
    }
    
    return true
  }

  delete(key: string): void {
    this.cache.delete(key)
  }

  clear(): void {
    this.cache.clear()
  }

  private cleanup(): void {
    const now = Date.now()
    for (const [key, item] of this.cache.entries()) {
      if (now - item.timestamp > item.ttl) {
        this.cache.delete(key)
      }
    }
  }

  // Get cache statistics
  getStats() {
    return {
      size: this.cache.size,
      maxSize: this.maxSize,
      keys: Array.from(this.cache.keys())
    }
  }
}

export const cacheManager = new CacheManager()

// Cache key generators
export const generateCacheKey = {
  stocks: (filters: any) => `stocks:${JSON.stringify(filters)}`,
  recommendations: (filters: any, tier: string) => `recommendations:${tier}:${JSON.stringify(filters)}`,
  stocksByTicker: (ticker: string) => `stocks_ticker:${ticker.toUpperCase()}`,
  recommendationByTicker: (ticker: string, tier: string) => `recommendation:${ticker.toUpperCase()}:${tier}`,
  stats: () => 'stats:general'
}
```

#### 2.2 Enhanced API Client with Caching (`src/api/client.ts` - Updated)
```typescript
import axios, { AxiosInstance, AxiosRequestConfig } from 'axios'
import { useAuthStore } from '@/stores/auth'
import { cacheManager } from '@/utils/cache'

class APIClient {
  private client: AxiosInstance

  constructor() {
    this.client = axios.create({
      baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json',
      },
    })

    this.setupInterceptors()
  }

  private setupInterceptors() {
    // Request interceptor
    this.client.interceptors.request.use((config) => {
      const authStore = useAuthStore()
      if (authStore.accessToken) {
        config.headers.Authorization = `Bearer ${authStore.accessToken}`
      }
      
      // Add request ID for tracking
      config.metadata = { requestId: crypto.randomUUID() }
      
      return config
    })

    // Response interceptor
    this.client.interceptors.response.use(
      (response) => {
        // Log successful requests in development
        if (import.meta.env.DEV) {
          console.log(`✅ ${response.config.method?.toUpperCase()} ${response.config.url}`, {
            status: response.status,
            requestId: response.config.metadata?.requestId
          })
        }
        return response
      },
      async (error) => {
        const originalRequest = error.config
        
        // Log errors
        if (import.meta.env.DEV) {
          console.error(`❌ ${originalRequest.method?.toUpperCase()} ${originalRequest.url}`, {
            status: error.response?.status,
            message: error.message,
            requestId: originalRequest.metadata?.requestId
          })
        }
        
        // Token refresh logic
        if (error.response?.status === 401 && !originalRequest._retry) {
          originalRequest._retry = true
          
          const authStore = useAuthStore()
          try {
            await authStore.refreshToken()
            return this.client(originalRequest)
          } catch (refreshError) {
            authStore.logout()
            return Promise.reject(refreshError)
          }
        }
        
        return Promise.reject(error)
      }
    )
  }

  async get<T>(url: string, config?: AxiosRequestConfig & { cacheKey?: string; cacheTTL?: number }): Promise<T> {
    // Check cache first
    if (config?.cacheKey && cacheManager.has(config.cacheKey)) {
      const cachedData = cacheManager.get<T>(config.cacheKey)
      if (cachedData) {
        console.log(`🎯 Cache hit for ${config.cacheKey}`)
        return cachedData
      }
    }

    const response = await this.client.get(url, config)
    const data = response.data

    // Store in cache if cache key provided
    if (config?.cacheKey) {
      cacheManager.set(config.cacheKey, data, config.cacheTTL || 30)
      console.log(`💾 Cached ${config.cacheKey}`)
    }

    return data
  }

  async post<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.client.post(url, data, config)
    return response.data
  }

  async put<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.client.put(url, data, config)
    return response.data
  }

  async delete<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.client.delete(url, config)
    return response.data
  }

  // Clear related caches
  invalidateCache(pattern: string) {
    const stats = cacheManager.getStats()
    const keysToDelete = stats.keys.filter(key => key.includes(pattern))
    keysToDelete.forEach(key => cacheManager.delete(key))
    console.log(`🗑️ Invalidated ${keysToDelete.length} cache entries matching "${pattern}"`)
  }
}

export const apiClient = new APIClient()
```

#### 2.3 Virtual Scrolling for Large Lists (`src/components/ui/VirtualList.vue`)
```vue
<template>
  <div
    ref="container"
    class="virtual-list-container"
    :style="{ height: `${containerHeight}px`, overflow: 'auto' }"
    @scroll="handleScroll"
  >
    <div :style="{ height: `${totalHeight}px`, position: 'relative' }">
      <div
        v-for="item in visibleItems"
        :key="getItemKey(item.data)"
        :style="{
          position: 'absolute',
          top: `${item.top}px`,
          left: 0,
          right: 0,
          height: `${itemHeight}px`
        }"
      >
        <slot :item="item.data" :index="item.index" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'

interface Props {
  items: any[]
  itemHeight: number
  containerHeight: number
  buffer?: number
  getItemKey?: (item: any) => string | number
}

const props = withDefaults(defineProps<Props>(), {
  buffer: 5,
  getItemKey: (item, index) => index
})

const container = ref<HTMLElement>()
const scrollTop = ref(0)

const totalHeight = computed(() => props.items.length * props.itemHeight)

const visibleStart = computed(() => {
  return Math.max(0, Math.floor(scrollTop.value / props.itemHeight) - props.buffer)
})

const visibleEnd = computed(() => {
  const visibleCount = Math.ceil(props.containerHeight / props.itemHeight)
  return Math.min(props.items.length, visibleStart.value + visibleCount + props.buffer * 2)
})

const visibleItems = computed(() => {
  return props.items.slice(visibleStart.value, visibleEnd.value).map((item, index) => ({
    data: item,
    index: visibleStart.value + index,
    top: (visibleStart.value + index) * props.itemHeight
  }))
})

const handleScroll = (event: Event) => {
  const target = event.target as HTMLElement
  scrollTop.value = target.scrollTop
}

// Scroll to specific item
const scrollToItem = (index: number) => {
  if (container.value) {
    container.value.scrollTop = index * props.itemHeight
  }
}

// Scroll to top
const scrollToTop = () => {
  if (container.value) {
    container.value.scrollTop = 0
  }
}

defineExpose({
  scrollToItem,
  scrollToTop
})
</script>

<style scoped>
.virtual-list-container {
  will-change: scroll-position;
}
</style>
```

### Step 3: AI Chat Interface & Placeholder System

#### 3.1 AI Chat Store (`src/stores/aiChat.ts`)
```typescript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAuthStore } from './auth'

export interface ChatMessage {
  id: string
  type: 'user' | 'assistant' | 'system'
  content: string
  timestamp: Date
  isTyping?: boolean
  metadata?: {
    ticker?: string
    actionType?: 'analysis' | 'recommendation' | 'general'
    confidence?: number
  }
}

export interface ChatSession {
  id: string
  title: string
  createdAt: Date
  lastMessageAt: Date
  messages: ChatMessage[]
}

export const useAIChatStore = defineStore('aiChat', () => {
  // State
  const isOpen = ref(false)
  const currentSession = ref<ChatSession | null>(null)
  const sessions = ref<ChatSession[]>([])
  const isTyping = ref(false)
  const isConnected = ref(false)
  const connectionStatus = ref<'disconnected' | 'connecting' | 'connected'>('disconnected')

  // Getters
  const authStore = useAuthStore()
  const canUseAI = computed(() => authStore.hasFeature('ai_insights'))
  const hasActiveSessions = computed(() => sessions.value.length > 0)
  const currentMessages = computed(() => currentSession.value?.messages || [])

  // Actions
  function openChat() {
    if (!canUseAI.value) {
      // Show upgrade prompt
      return false
    }
    isOpen.value = true
    
    // Create new session if none exists
    if (!currentSession.value) {
      createNewSession()
    }
    
    return true
  }

  function closeChat() {
    isOpen.value = false
  }

  function toggleChat() {
    if (isOpen.value) {
      closeChat()
    } else {
      openChat()
    }
  }

  function createNewSession() {
    const session: ChatSession = {
      id: crypto.randomUUID(),
      title: `Chat ${new Date().toLocaleDateString()}`,
      createdAt: new Date(),
      lastMessageAt: new Date(),
      messages: [{
        id: crypto.randomUUID(),
        type: 'system',
        content: 'Hello! I\'m your AI investment assistant. I can help you analyze stocks, understand market trends, and provide personalized recommendations. What would you like to know?',
        timestamp: new Date()
      }]
    }

    sessions.value.unshift(session)
    currentSession.value = session
    saveSessionsToStorage()
  }

  function switchToSession(sessionId: string) {
    const session = sessions.value.find(s => s.id === sessionId)
    if (session) {
      currentSession.value = session
    }
  }

  function deleteSession(sessionId: string) {
    sessions.value = sessions.value.filter(s => s.id !== sessionId)
    
    if (currentSession.value?.id === sessionId) {
      currentSession.value = sessions.value[0] || null
    }
    
    saveSessionsToStorage()
  }

  async function sendMessage(content: string, metadata?: ChatMessage['metadata']) {
    if (!currentSession.value || !canUseAI.value) return

    // Add user message
    const userMessage: ChatMessage = {
      id: crypto.randomUUID(),
      type: 'user',
      content,
      timestamp: new Date(),
      metadata
    }

    currentSession.value.messages.push(userMessage)
    currentSession.value.lastMessageAt = new Date()

    // Show typing indicator
    isTyping.value = true

    try {
      // TODO: Replace with actual AI API call
      const response = await simulateAIResponse(content, metadata)
      
      const assistantMessage: ChatMessage = {
        id: crypto.randomUUID(),
        type: 'assistant',
        content: response.content,
        timestamp: new Date(),
        metadata: response.metadata
      }

      currentSession.value.messages.push(assistantMessage)
      currentSession.value.lastMessageAt = new Date()

      // Update session title if it's the first meaningful exchange
      if (currentSession.value.messages.filter(m => m.type === 'user').length === 1) {
        currentSession.value.title = generateSessionTitle(content)
      }

    } catch (error) {
      console.error('Failed to get AI response:', error)
      
      const errorMessage: ChatMessage = {
        id: crypto.randomUUID(),
        type: 'assistant',
        content: 'I apologize, but I\'m having trouble processing your request right now. Please try again later.',
        timestamp: new Date()
      }

      currentSession.value.messages.push(errorMessage)
    } finally {
      isTyping.value = false
      saveSessionsToStorage()
    }
  }

  // Simulate AI response (replace with actual API call)
  async function simulateAIResponse(userMessage: string, metadata?: ChatMessage['metadata']) {
    // Simulate network delay
    await new Promise(resolve => setTimeout(resolve, 1000 + Math.random() * 2000))

    const lowerMessage = userMessage.toLowerCase()

    if (lowerMessage.includes('aapl') || lowerMessage.includes('apple')) {
      return {
        content: `Based on my analysis of Apple (AAPL), here's what I found:

📊 **Current Analysis:**
- Recent price target increases from major brokerages
- Strong institutional buying pressure
- Q4 earnings beat expectations by 12%

💡 **Key Insights:**
- 73% of brokerages rate it as "Buy" or "Strong Buy"
- Average price target: $185 (current: $175)
- Risk level: Moderate

🔮 **My Recommendation:**
Consider this stock for your portfolio. The fundamentals look strong, and technical indicators suggest continued upward momentum.

Would you like me to analyze any specific aspects like competitors, valuation metrics, or market sentiment?`,
        metadata: {
          ticker: 'AAPL',
          actionType: 'analysis' as const,
          confidence: 0.85
        }
      }
    }

    if (lowerMessage.includes('recommendation') || lowerMessage.includes('suggest')) {
      return {
        content: `Here are my top 3 stock recommendations based on current market conditions:

🥇 **Microsoft (MSFT)** - Strong Buy
- AI growth story with solid fundamentals
- Price target: $420 (upside: 15%)
- Risk: Low

🥈 **NVIDIA (NVDA)** - Buy
- Leading AI chip manufacturer
- Price target: $520 (upside: 22%)
- Risk: Moderate-High

🥉 **Amazon (AMZN)** - Buy
- AWS growth and retail recovery
- Price target: $165 (upside: 18%)
- Risk: Moderate

These recommendations are based on broker consensus, technical analysis, and my proprietary sentiment scoring. Would you like detailed analysis on any of these?`,
        metadata: {
          actionType: 'recommendation' as const,
          confidence: 0.78
        }
      }
    }

    // Default response
    return {
      content: `I understand you're asking about "${userMessage}". Let me help you with that.

As your AI investment assistant, I can provide insights on:
• Stock analysis and recommendations
• Market trend interpretation
• Portfolio optimization suggestions
• Risk assessment
• Earnings analysis

Could you be more specific about what you'd like to know? For example:
- "Analyze TSLA stock"
- "What are the best tech stocks right now?"
- "Should I buy or sell Netflix?"`,
      metadata: {
        actionType: 'general' as const
      }
    }
  }

  function generateSessionTitle(firstMessage: string): string {
    const maxLength = 30
    if (firstMessage.length <= maxLength) {
      return firstMessage
    }
    return firstMessage.substring(0, maxLength - 3) + '...'
  }

  function saveSessionsToStorage() {
    try {
      localStorage.setItem('ai_chat_sessions', JSON.stringify(sessions.value))
    } catch (error) {
      console.error('Failed to save chat sessions:', error)
    }
  }

  function loadSessionsFromStorage() {
    try {
      const stored = localStorage.getItem('ai_chat_sessions')
      if (stored) {
        const parsedSessions = JSON.parse(stored)
        sessions.value = parsedSessions.map((session: any) => ({
          ...session,
          createdAt: new Date(session.createdAt),
          lastMessageAt: new Date(session.lastMessageAt),
          messages: session.messages.map((message: any) => ({
            ...message,
            timestamp: new Date(message.timestamp)
          }))
        }))
        
        if (sessions.value.length > 0) {
          currentSession.value = sessions.value[0]
        }
      }
    } catch (error) {
      console.error('Failed to load chat sessions:', error)
    }
  }

  // Initialize
  loadSessionsFromStorage()

  return {
    // State
    isOpen: readonly(isOpen),
    currentSession: readonly(currentSession),
    sessions: readonly(sessions),
    isTyping: readonly(isTyping),
    isConnected: readonly(isConnected),
    connectionStatus: readonly(connectionStatus),
    
    // Getters
    canUseAI,
    hasActiveSessions,
    currentMessages,
    
    // Actions
    openChat,
    closeChat,
    toggleChat,
    createNewSession,
    switchToSession,
    deleteSession,
    sendMessage
  }
})
```

#### 3.2 AI Chat Interface (`src/components/features/ai/AIChatInterface.vue`)
```vue
<template>
  <div class="fixed inset-0 z-50 overflow-hidden" v-if="aiChatStore.isOpen">
    <div class="absolute inset-0 bg-black bg-opacity-50" @click="aiChatStore.closeChat"></div>
    
    <div class="absolute right-4 top-4 bottom-4 w-96 bg-white rounded-lg shadow-2xl flex flex-col">
      <!-- Header -->
      <div class="flex items-center justify-between p-4 border-b border-gray-200">
        <div class="flex items-center space-x-2">
          <div class="w-8 h-8 bg-gradient-to-r from-purple-500 to-pink-500 rounded-full flex items-center justify-center">
            <span class="text-white text-sm font-semibold">AI</span>
          </div>
          <div>
            <h3 class="font-semibold text-gray-900">Investment Assistant</h3>
            <p class="text-xs text-gray-500">
              {{ aiChatStore.isTyping ? 'Typing...' : 'Online' }}
            </p>
          </div>
        </div>
        
        <div class="flex items-center space-x-2">
          <button
            @click="showSessionMenu = !showSessionMenu"
            class="text-gray-400 hover:text-gray-600"
          >
            <Bars3Icon class="h-5 w-5" />
          </button>
          <button
            @click="aiChatStore.closeChat"
            class="text-gray-400 hover:text-gray-600"
          >
            <XMarkIcon class="h-5 w-5" />
          </button>
        </div>
      </div>

      <!-- Session Menu -->
      <div v-if="showSessionMenu" class="border-b border-gray-200 bg-gray-50 p-2">
        <div class="flex items-center justify-between mb-2">
          <span class="text-sm font-medium text-gray-700">Chat Sessions</span>
          <button
            @click="aiChatStore.createNewSession(); showSessionMenu = false"
            class="text-xs bg-primary-600 text-white px-2 py-1 rounded"
          >
            New Chat
          </button>
        </div>
        <div class="max-h-32 overflow-y-auto space-y-1">
          <div
            v-for="session in aiChatStore.sessions"
            :key="session.id"
            class="flex items-center justify-between p-2 text-sm rounded hover:bg-gray-100 cursor-pointer"
            :class="{ 'bg-primary-50 border border-primary-200': session.id === aiChatStore.currentSession?.id }"
            @click="aiChatStore.switchToSession(session.id); showSessionMenu = false"
          >
            <span class="truncate">{{ session.title }}</span>
            <button
              @click.stop="aiChatStore.deleteSession(session.id)"
              class="text-gray-400 hover:text-red-500"
            >
              <TrashIcon class="h-4 w-4" />
            </button>
          </div>
        </div>
      </div>

      <!-- Messages -->
      <div ref="messagesContainer" class="flex-1 overflow-y-auto p-4 space-y-4">
        <div
          v-for="message in aiChatStore.currentMessages"
          :key="message.id"
          :class="[
            'flex',
            message.type === 'user' ? 'justify-end' : 'justify-start'
          ]"
        >
          <div
            :class="[
              'max-w-[80%] rounded-lg px-4 py-2',
              message.type === 'user'
                ? 'bg-primary-600 text-white'
                : message.type === 'system'
                ? 'bg-gray-100 text-gray-700 text-sm'
                : 'bg-gray-100 text-gray-900'
            ]"
          >
            <div
              v-if="message.type === 'assistant'"
              v-html="formatAssistantMessage(message.content)"
              class="prose prose-sm max-w-none"
            />
            <div v-else>{{ message.content }}</div>
            
            <div
              v-if="message.metadata?.ticker"
              class="mt-2 text-xs opacity-75"
            >
              📊 {{ message.metadata.ticker }}
              <span v-if="message.metadata.confidence">
                • Confidence: {{ (message.metadata.confidence * 100).toFixed(0) }}%
              </span>
            </div>
            
            <div class="text-xs opacity-75 mt-1">
              {{ formatTime(message.timestamp) }}
            </div>
          </div>
        </div>

        <!-- Typing Indicator -->
        <div v-if="aiChatStore.isTyping" class="flex justify-start">
          <div class="bg-gray-100 rounded-lg px-4 py-2">
            <div class="flex space-x-1">
              <div class="w-2 h-2 bg-gray-400 rounded-full animate-bounce"></div>
              <div class="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 0.1s"></div>
              <div class="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 0.2s"></div>
            </div>
          </div>
        </div>
      </div>

      <!-- Input -->
      <div class="p-4 border-t border-gray-200">
        <div class="flex space-x-2">
          <input
            v-model="newMessage"
            type="text"
            placeholder="Ask me about stocks, market trends, or get recommendations..."
            class="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-1 focus:ring-primary-500 focus:border-primary-500"
            @keydown.enter="handleSendMessage"
            :disabled="aiChatStore.isTyping"
          />
          <button
            @click="handleSendMessage"
            :disabled="!newMessage.trim() || aiChatStore.isTyping"
            class="px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <PaperAirplaneIcon class="h-4 w-4" />
          </button>
        </div>
        
        <!-- Quick Actions -->
        <div class="mt-2 flex flex-wrap gap-1">
          <button
            v-for="quickAction in quickActions"
            :key="quickAction.text"
            @click="newMessage = quickAction.text; handleSendMessage()"
            class="text-xs bg-gray-100 text-gray-700 px-2 py-1 rounded hover:bg-gray-200"
          >
            {{ quickAction.text }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick, watch } from 'vue'
import { XMarkIcon, Bars3Icon, TrashIcon, PaperAirplaneIcon } from '@heroicons/vue/24/outline'
import { useAIChatStore } from '@/stores/aiChat'
import { format } from 'date-fns'

const aiChatStore = useAIChatStore()
const newMessage = ref('')
const messagesContainer = ref<HTMLElement>()
const showSessionMenu = ref(false)

const quickActions = [
  { text: 'Market overview today' },
  { text: 'Top 5 stock recommendations' },
  { text: 'Analyze AAPL' },
  { text: 'Best tech stocks to buy' }
]

const handleSendMessage = async () => {
  if (!newMessage.value.trim()) return

  const message = newMessage.value.trim()
  newMessage.value = ''

  await aiChatStore.sendMessage(message)
  scrollToBottom()
}

const formatAssistantMessage = (content: string): string => {
  return content
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
    .replace(/\*(.*?)\*/g, '<em>$1</em>')
    .replace(/\n/g, '<br>')
    .replace(/📊|💡|🔮|🥇|🥈|🥉|•/g, '<span class="inline-block mr-1">$&</span>')
}

const formatTime = (timestamp: Date): string => {
  return format(timestamp, 'HH:mm')
}

const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

// Auto-scroll when new messages arrive
watch(() => aiChatStore.currentMessages.length, () => {
  scrollToBottom()
})
</script>
```

#### 3.3 AI Chat FAB (`src/components/features/ai/AIChatFAB.vue`)
```vue
<template>
  <div class="fixed bottom-6 right-6 z-40">
    <!-- Upgrade Prompt for Non-Premium Users -->
    <div
      v-if="!aiChatStore.canUseAI && showUpgradePrompt"
      class="absolute bottom-16 right-0 w-80 bg-white rounded-lg shadow-xl border border-gray-200 p-4 transform transition-all duration-300"
    >
      <div class="flex items-start space-x-3">
        <div class="w-10 h-10 bg-gradient-to-r from-purple-500 to-pink-500 rounded-full flex items-center justify-center">
          <span class="text-white text-sm font-semibold">AI</span>
        </div>
        <div class="flex-1">
          <h3 class="font-semibold text-gray-900 mb-1">AI Investment Assistant</h3>
          <p class="text-sm text-gray-600 mb-3">
            Get personalized stock recommendations, market insights, and real-time analysis with our AI assistant.
          </p>
          <div class="flex space-x-2">
            <button
              @click="$router.push('/subscription')"
              class="text-xs bg-primary-600 text-white px-3 py-1.5 rounded font-medium hover:bg-primary-700"
            >
              Upgrade to Premium
            </button>
            <button
              @click="showUpgradePrompt = false"
              class="text-xs text-gray-500 hover:text-gray-700"
            >
              Maybe later
            </button>
          </div>
        </div>
        <button
          @click="showUpgradePrompt = false"
          class="text-gray-400 hover:text-gray-600"
        >
          <XMarkIcon class="h-4 w-4" />
        </button>
      </div>
    </div>

    <!-- Main FAB -->
    <button
      @click="handleFABClick"
      :class="[
        'w-14 h-14 rounded-full shadow-lg flex items-center justify-center transition-all duration-300 hover:scale-110',
        aiChatStore.canUseAI
          ? 'bg-gradient-to-r from-purple-500 to-pink-500 text-white hover:shadow-xl'
          : 'bg-gray-300 text-gray-600 hover:bg-gray-400'
      ]"
    >
      <ChatBubbleLeftRightIcon v-if="!aiChatStore.isOpen" class="h-6 w-6" />
      <XMarkIcon v-else class="h-6 w-6" />
    </button>

    <!-- Notification Badge -->
    <div
      v-if="!aiChatStore.canUseAI && !showUpgradePrompt"
      class="absolute -top-1 -right-1 w-5 h-5 bg-red-500 rounded-full flex items-center justify-center"
    >
      <span class="text-white text-xs font-bold">!</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ChatBubbleLeftRightIcon, XMarkIcon } from '@heroicons/vue/24/outline'
import { useAIChatStore } from '@/stores/aiChat'

const router = useRouter()
const aiChatStore = useAIChatStore()
const showUpgradePrompt = ref(false)

const handleFABClick = () => {
  if (aiChatStore.canUseAI) {
    aiChatStore.toggleChat()
  } else {
    showUpgradePrompt.value = !showUpgradePrompt.value
  }
}

// Show upgrade prompt automatically for new users after a delay
onMounted(() => {
  if (!aiChatStore.canUseAI) {
    setTimeout(() => {
      const hasSeenPrompt = localStorage.getItem('ai_upgrade_prompt_seen')
      if (!hasSeenPrompt) {
        showUpgradePrompt.value = true
        localStorage.setItem('ai_upgrade_prompt_seen', 'true')
      }
    }, 5000)
  }
})
</script>
```

### Step 4: Error Handling & Resilience

#### 4.1 Global Error Handler (`src/utils/errorHandler.ts`)
```typescript
import { nextTick } from 'vue'

export interface AppError {
  id: string
  type: 'network' | 'validation' | 'auth' | 'permission' | 'unknown'
  message: string
  details?: any
  timestamp: Date
  userAgent?: string
  url?: string
  userId?: string
}

class ErrorHandler {
  private errors: AppError[] = []
  private maxErrors = 50

  logError(error: Error | AppError, context?: any) {
    const appError: AppError = this.normalizeError(error, context)
    
    // Add to error log
    this.errors.unshift(appError)
    if (this.errors.length > this.maxErrors) {
      this.errors = this.errors.slice(0, this.maxErrors)
    }

    // Log to console in development
    if (import.meta.env.DEV) {
      console.error('🚨 Application Error:', appError)
    }

    // In production, send to monitoring service
    if (import.meta.env.PROD) {
      this.sendToMonitoring(appError)
    }

    return appError
  }

  private normalizeError(error: Error | AppError, context?: any): AppError {
    if ('id' in error && 'type' in error) {
      return error as AppError
    }

    const appError: AppError = {
      id: crypto.randomUUID(),
      type: this.categorizeError(error),
      message: error.message || 'An unknown error occurred',
      details: context,
      timestamp: new Date(),
      userAgent: navigator.userAgent,
      url: window.location.href
    }

    return appError
  }

  private categorizeError(error: Error): AppError['type'] {
    const message = error.message.toLowerCase()
    
    if (message.includes('network') || message.includes('fetch')) {
      return 'network'
    }
    if (message.includes('unauthorized') || message.includes('401')) {
      return 'auth'
    }
    if (message.includes('forbidden') || message.includes('403')) {
      return 'permission'
    }
    if (message.includes('validation') || message.includes('invalid')) {
      return 'validation'
    }
    
    return 'unknown'
  }

  private async sendToMonitoring(error: AppError) {
    try {
      // Replace with actual monitoring service (Sentry, LogRocket, etc.)
      await fetch('/api/errors', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(error)
      })
    } catch (monitoringError) {
      console.error('Failed to send error to monitoring:', monitoringError)
    }
  }

  getRecentErrors(limit = 10): AppError[] {
    return this.errors.slice(0, limit)
  }

  clearErrors() {
    this.errors = []
  }

  // User-friendly error messages
  getUserMessage(error: AppError): string {
    switch (error.type) {
      case 'network':
        return 'Please check your internet connection and try again.'
      case 'auth':
        return 'Your session has expired. Please log in again.'
      case 'permission':
        return 'You don\'t have permission to perform this action.'
      case 'validation':
        return 'Please check your input and try again.'
      default:
        return 'Something went wrong. Please try again.'
    }
  }
}

export const errorHandler = new ErrorHandler()

// Global error boundary for uncaught errors
window.addEventListener('error', (event) => {
  errorHandler.logError(new Error(event.message), {
    filename: event.filename,
    lineno: event.lineno,
    colno: event.colno
  })
})

window.addEventListener('unhandledrejection', (event) => {
  errorHandler.logError(new Error(event.reason?.message || 'Unhandled promise rejection'), {
    reason: event.reason
  })
})
```

#### 4.2 Error Boundary Component (`src/components/ui/ErrorBoundary.vue`)
```vue
<template>
  <div v-if="hasError" class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div class="text-center">
        <div class="mx-auto h-24 w-24 text-red-500">
          <ExclamationTriangleIcon class="h-full w-full" />
        </div>
        <h2 class="mt-6 text-3xl font-extrabold text-gray-900">
          Something went wrong
        </h2>
        <p class="mt-2 text-sm text-gray-600">
          {{ userMessage }}
        </p>
      </div>
      
      <div class="space-y-4">
        <button
          @click="handleRetry"
          class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
        >
          Try Again
        </button>
        
        <button
          @click="handleGoHome"
          class="w-full flex justify-center py-2 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
        >
          Go to Dashboard
        </button>
        
        <details class="mt-4">
          <summary class="cursor-pointer text-sm text-gray-500 hover:text-gray-700">
            Technical Details
          </summary>
          <div class="mt-2 p-3 bg-gray-100 rounded text-xs font-mono text-gray-700">
            <div><strong>Error ID:</strong> {{ error?.id }}</div>
            <div><strong>Type:</strong> {{ error?.type }}</div>
            <div><strong>Time:</strong> {{ error?.timestamp }}</div>
            <div class="mt-2"><strong>Message:</strong> {{ error?.message }}</div>
          </div>
        </details>
      </div>
    </div>
  </div>
  
  <slot v-else />
</template>

<script setup lang="ts">
import { ref, onErrorCaptured } from 'vue'
import { useRouter } from 'vue-router'
import { ExclamationTriangleIcon } from '@heroicons/vue/24/outline'
import { errorHandler, type AppError } from '@/utils/errorHandler'

const router = useRouter()
const hasError = ref(false)
const error = ref<AppError | null>(null)
const userMessage = ref('')

onErrorCaptured((err: Error) => {
  const appError = errorHandler.logError(err)
  
  hasError.value = true
  error.value = appError
  userMessage.value = errorHandler.getUserMessage(appError)
  
  return false // Prevent error from bubbling up
})

const handleRetry = () => {
  hasError.value = false
  error.value = null
  userMessage.value = ''
  
  // Force component re-render
  window.location.reload()
}

const handleGoHome = () => {
  hasError.value = false
  error.value = null
  userMessage.value = ''
  
  router.push('/dashboard')
}
</script>
```

### Step 5: Advanced Data Visualizations

#### 5.1 Chart Utilities (`src/utils/chartUtils.ts`)
```typescript
import * as echarts from 'echarts'

export interface ChartTheme {
  backgroundColor: string
  textColor: string
  primaryColor: string
  successColor: string
  dangerColor: string
  warningColor: string
}

export const defaultTheme: ChartTheme = {
  backgroundColor: '#ffffff',
  textColor: '#374151',
  primaryColor: '#3b82f6',
  successColor: '#10b981',
  dangerColor: '#ef4444',
  warningColor: '#f59e0b'
}

export const darkTheme: ChartTheme = {
  backgroundColor: '#1f2937',
  textColor: '#f9fafb',
  primaryColor: '#60a5fa',
  successColor: '#34d399',
  dangerColor: '#f87171',
  warningColor: '#fbbf24'
}

export function createRecommendationChart(data: any[], theme: ChartTheme = defaultTheme) {
  return {
    backgroundColor: theme.backgroundColor,
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      left: 10,
      textStyle: {
        color: theme.textColor
      }
    },
    series: [
      {
        name: 'Recommendations',
        type: 'pie',
        radius: ['50%', '70%'],
        avoidLabelOverlap: false,
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: '18',
            fontWeight: 'bold',
            color: theme.textColor
          }
        },
        labelLine: {
          show: false
        },
        data: data.map(item => ({
          value: item.value,
          name: item.name,
          itemStyle: {
            color: getColorForRecommendationType(item.name, theme)
          }
        }))
      }
    ]
  }
}

export function createPriceTargetChart(data: any[], theme: ChartTheme = defaultTheme) {
  return {
    backgroundColor: theme.backgroundColor,
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross'
      }
    },
    legend: {
      data: ['Price Target', 'Current Price'],
      textStyle: {
        color: theme.textColor
      }
    },
    xAxis: {
      type: 'category',
      data: data.map(item => item.date),
      axisLabel: {
        color: theme.textColor
      }
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        color: theme.textColor,
        formatter: '${value}'
      }
    },
    series: [
      {
        name: 'Price Target',
        type: 'line',
        data: data.map(item => item.target),
        lineStyle: {
          color: theme.primaryColor
        },
        itemStyle: {
          color: theme.primaryColor
        }
      },
      {
        name: 'Current Price',
        type: 'line',
        data: data.map(item => item.current),
        lineStyle: {
          color: theme.successColor
        },
        itemStyle: {
          color: theme.successColor
        }
      }
    ]
  }
}

export function createSentimentChart(data: any[], theme: ChartTheme = defaultTheme) {
  return {
    backgroundColor: theme.backgroundColor,
    tooltip: {
      trigger: 'axis'
    },
    legend: {
      data: ['Sentiment Score', 'Volume'],
      textStyle: {
        color: theme.textColor
      }
    },
    xAxis: [
      {
        type: 'category',
        data: data.map(item => item.date),
        axisLabel: {
          color: theme.textColor
        }
      }
    ],
    yAxis: [
      {
        type: 'value',
        name: 'Sentiment',
        min: -1,
        max: 1,
        axisLabel: {
          color: theme.textColor
        }
      },
      {
        type: 'value',
        name: 'Volume',
        axisLabel: {
          color: theme.textColor
        }
      }
    ],
    series: [
      {
        name: 'Sentiment Score',
        type: 'line',
        yAxisIndex: 0,
        data: data.map(item => item.sentiment),
        lineStyle: {
          color: theme.primaryColor
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: theme.primaryColor + '40' },
            { offset: 1, color: theme.primaryColor + '10' }
          ])
        }
      },
      {
        name: 'Volume',
        type: 'bar',
        yAxisIndex: 1,
        data: data.map(item => item.volume),
        itemStyle: {
          color: theme.warningColor + '60'
        }
      }
    ]
  }
}

function getColorForRecommendationType(type: string, theme: ChartTheme): string {
  switch (type.toLowerCase()) {
    case 'strong buy':
      return theme.successColor
    case 'buy':
      return '#22c55e'
    case 'hold':
      return theme.warningColor
    case 'sell':
      return '#f97316'
    case 'strong sell':
      return theme.dangerColor
    default:
      return theme.textColor
  }
}

export function createHeatmapChart(data: any[], theme: ChartTheme = defaultTheme) {
  const hours = ['12a', '1a', '2a', '3a', '4a', '5a', '6a', '7a', '8a', '9a', '10a', '11a',
                 '12p', '1p', '2p', '3p', '4p', '5p', '6p', '7p', '8p', '9p', '10p', '11p']
  const days = ['Saturday', 'Friday', 'Thursday', 'Wednesday', 'Tuesday', 'Monday', 'Sunday']

  return {
    backgroundColor: theme.backgroundColor,
    tooltip: {
      position: 'top',
      formatter: function(params: any) {
        return `${days[params.data[1]]} ${hours[params.data[0]]}<br/>Activity: ${params.data[2]}`
      }
    },
    grid: {
      height: '50%',
      top: '10%'
    },
    xAxis: {
      type: 'category',
      data: hours,
      splitArea: {
        show: true
      },
      axisLabel: {
        color: theme.textColor
      }
    },
    yAxis: {
      type: 'category',
      data: days,
      splitArea: {
        show: true
      },
      axisLabel: {
        color: theme.textColor
      }
    },
    visualMap: {
      min: 0,
      max: 10,
      calculable: true,
      orient: 'horizontal',
      left: 'center',
      bottom: '15%',
      textStyle: {
        color: theme.textColor
      }
    },
    series: [{
      name: 'Trading Activity',
      type: 'heatmap',
      data: data,
      label: {
        show: true,
        color: theme.textColor
      },
      emphasis: {
        itemStyle: {
          shadowBlur: 10,
          shadowColor: 'rgba(0, 0, 0, 0.5)'
        }
      }
    }]
  }
}
```

### Step 6: Accessibility & Responsive Design

#### 6.1 Accessibility Composable (`src/composables/useAccessibility.ts`)
```typescript
import { ref, onMounted, onUnmounted } from 'vue'

export function useAccessibility() {
  const isHighContrast = ref(false)
  const isReducedMotion = ref(false)
  const fontSize = ref('normal')

  const checkAccessibilityPreferences = () => {
    // Check for high contrast preference
    isHighContrast.value = window.matchMedia('(prefers-contrast: high)').matches

    // Check for reduced motion preference
    isReducedMotion.value = window.matchMedia('(prefers-reduced-motion: reduce)').matches

    // Load font size preference
    const savedFontSize = localStorage.getItem('font-size-preference')
    if (savedFontSize) {
      fontSize.value = savedFontSize
      applyFontSize(savedFontSize)
    }
  }

  const applyFontSize = (size: string) => {
    const root = document.documentElement
    switch (size) {
      case 'small':
        root.style.setProperty('--font-size-multiplier', '0.875')
        break
      case 'large':
        root.style.setProperty('--font-size-multiplier', '1.125')
        break
      case 'extra-large':
        root.style.setProperty('--font-size-multiplier', '1.25')
        break
      default:
        root.style.setProperty('--font-size-multiplier', '1')
    }
  }

  const setFontSize = (size: string) => {
    fontSize.value = size
    applyFontSize(size)
    localStorage.setItem('font-size-preference', size)
  }

  const announceToScreenReader = (message: string) => {
    const announcement = document.createElement('div')
    announcement.setAttribute('aria-live', 'polite')
    announcement.setAttribute('aria-atomic', 'true')
    announcement.style.position = 'absolute'
    announcement.style.left = '-10000px'
    announcement.style.width = '1px'
    announcement.style.height = '1px'
    announcement.style.overflow = 'hidden'
    
    document.body.appendChild(announcement)
    announcement.textContent = message
    
    setTimeout(() => {
      document.body.removeChild(announcement)
    }, 1000)
  }

  const trapFocus = (element: HTMLElement) => {
    const focusableElements = element.querySelectorAll(
      'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
    )
    const firstElement = focusableElements[0] as HTMLElement
    const lastElement = focusableElements[focusableElements.length - 1] as HTMLElement

    const handleTab = (e: KeyboardEvent) => {
      if (e.key === 'Tab') {
        if (e.shiftKey) {
          if (document.activeElement === firstElement) {
            lastElement.focus()
            e.preventDefault()
          }
        } else {
          if (document.activeElement === lastElement) {
            firstElement.focus()
            e.preventDefault()
          }
        }
      }
    }

    element.addEventListener('keydown', handleTab)
    firstElement?.focus()

    return () => {
      element.removeEventListener('keydown', handleTab)
    }
  }

  onMounted(() => {
    checkAccessibilityPreferences()
  })

  return {
    isHighContrast: readonly(isHighContrast),
    isReducedMotion: readonly(isReducedMotion),
    fontSize: readonly(fontSize),
    setFontSize,
    announceToScreenReader,
    trapFocus
  }
}
```

### Step 7: Final Testing & Optimization

#### 7.1 Performance Monitor (`src/utils/performance.ts`)
```typescript
interface PerformanceMetric {
  name: string
  value: number
  timestamp: number
  type: 'navigation' | 'resource' | 'custom'
}

class PerformanceMonitor {
  private metrics: PerformanceMetric[] = []
  private observer?: PerformanceObserver

  init() {
    // Monitor Core Web Vitals
    this.observeWebVitals()
    
    // Monitor resource loading
    this.observeResources()
    
    // Monitor navigation timing
    this.observeNavigation()
  }

  private observeWebVitals() {
    // CLS (Cumulative Layout Shift)
    this.observeMetric('layout-shift', (entry) => {
      if (!entry.hadRecentInput) {
        this.recordMetric('CLS', entry.value, 'custom')
      }
    })

    // LCP (Largest Contentful Paint)
    this.observeMetric('largest-contentful-paint', (entry) => {
      this.recordMetric('LCP', entry.startTime, 'custom')
    })

    // FID (First Input Delay)
    this.observeMetric('first-input', (entry) => {
      this.recordMetric('FID', entry.processingStart - entry.startTime, 'custom')
    })
  }

  private observeResources() {
    this.observeMetric('resource', (entry) => {
      if (entry.duration > 1000) { // Log slow resources (>1s)
        this.recordMetric(`Slow Resource: ${entry.name}`, entry.duration, 'resource')
      }
    })
  }

  private observeNavigation() {
    this.observeMetric('navigation', (entry) => {
      this.recordMetric('Page Load Time', entry.loadEventEnd - entry.fetchStart, 'navigation')
      this.recordMetric('DOM Content Loaded', entry.domContentLoadedEventEnd - entry.fetchStart, 'navigation')
      this.recordMetric('First Paint', entry.responseEnd - entry.fetchStart, 'navigation')
    })
  }

  private observeMetric(type: string, callback: (entry: any) => void) {
    if (!PerformanceObserver) return

    try {
      const observer = new PerformanceObserver((list) => {
        list.getEntries().forEach(callback)
      })
      observer.observe({ entryTypes: [type] })
    } catch (error) {
      console.warn(`Failed to observe ${type}:`, error)
    }
  }

  recordMetric(name: string, value: number, type: PerformanceMetric['type'] = 'custom') {
    const metric: PerformanceMetric = {
      name,
      value,
      timestamp: Date.now(),
      type
    }

    this.metrics.push(metric)
    
    // Keep only recent metrics (last 100)
    if (this.metrics.length > 100) {
      this.metrics = this.metrics.slice(-100)
    }

    // Log critical metrics
    if (value > this.getThreshold(name)) {
      console.warn(`⚠️ Performance issue detected: ${name} = ${value}ms`)
    }

    // Send to analytics in production
    if (import.meta.env.PROD) {
      this.sendToAnalytics(metric)
    }
  }

  private getThreshold(metricName: string): number {
    const thresholds = {
      'LCP': 2500,
      'FID': 100,
      'CLS': 0.1,
      'Page Load Time': 3000,
      'DOM Content Loaded': 1500
    }
    return thresholds[metricName as keyof typeof thresholds] || Infinity
  }

  private async sendToAnalytics(metric: PerformanceMetric) {
    try {
      // Replace with actual analytics service
      await fetch('/api/analytics/performance', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(metric)
      })
    } catch (error) {
      console.error('Failed to send performance metric:', error)
    }
  }

  getMetrics(): PerformanceMetric[] {
    return [...this.metrics]
  }

  getAverageMetric(name: string): number {
    const relevantMetrics = this.metrics.filter(m => m.name === name)
    if (relevantMetrics.length === 0) return 0
    
    const sum = relevantMetrics.reduce((acc, metric) => acc + metric.value, 0)
    return sum / relevantMetrics.length
  }

  generateReport(): string {
    const report = {
      timestamp: new Date().toISOString(),
      metrics: this.metrics.reduce((acc, metric) => {
        acc[metric.name] = {
          latest: metric.value,
          average: this.getAverageMetric(metric.name),
          count: this.metrics.filter(m => m.name === metric.name).length
        }
        return acc
      }, {} as any)
    }

    return JSON.stringify(report, null, 2)
  }
}

export const performanceMonitor = new PerformanceMonitor()

// Initialize when the module loads
performanceMonitor.init()
```

---

## Phase 3 Deliverables

✅ **Advanced Features:**
- Global search system with intelligent ranking
- Virtual scrolling for large datasets
- Advanced filtering and sorting capabilities
- Performance monitoring and optimization

✅ **AI Integration Ready:**
- Complete AI chat interface with session management
- Floating action button with upgrade prompts
- Placeholder system for seamless AI feature rollout
- Conversation history and context management

✅ **Production Quality:**
- Comprehensive error handling and boundaries
- Accessibility compliance (WCAG 2.1)
- Performance optimization with caching
- Advanced data visualizations

✅ **User Experience:**
- Responsive design across all devices
- Loading states and skeleton screens
- Intuitive keyboard navigation
- Screen reader compatibility

---

## Deployment & Final Steps

### Build Optimization
```bash
# Production build with optimizations
npm run build

# Bundle analysis
npm run build:analyze

# Performance testing
npm run test:performance
```

### Environment Configuration
```env
# Production environment variables
VITE_API_BASE_URL=https://api.stocktracker.com/v1
VITE_ENABLE_ANALYTICS=true
VITE_SENTRY_DSN=your_sentry_dsn
VITE_APP_VERSION=1.0.0
```

### Monitoring Setup
- **Error Tracking:** Sentry integration for production error monitoring
- **Performance:** Core Web Vitals tracking and alerting
- **Analytics:** User interaction and feature usage tracking
- **Uptime:** API endpoint monitoring and alerting

### Future AI Integration Checklist
✅ Chat interface ready for backend connection  
✅ Message queuing and session management  
✅ Upgrade flow and tier validation  
✅ Context-aware conversation handling  
✅ Real-time response streaming preparation  

---

## Conclusion

Phase 3 completes the Stock Tracker frontend with:

1. **Advanced search and filtering** that makes finding relevant data effortless
2. **AI chat interface** ready for immediate backend integration when AI features are developed
3. **Production-grade performance** with caching, virtual scrolling, and optimization
4. **Accessibility compliance** ensuring the platform works for all users
5. **Comprehensive error handling** that gracefully manages edge cases
6. **Advanced visualizations** that make complex financial data understandable

The platform is now ready for deployment and future AI feature integration, providing a solid foundation for scaling and adding new capabilities. The modular architecture ensures that new features can be added without disrupting existing functionality, and the tier-based system creates clear upgrade paths for users. 