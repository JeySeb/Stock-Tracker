# Frontend Development - Phase 2: Main Features & Data Visualization

## Phase Overview

**Duration:** 3-4 weeks  
**Focus:** Implement core business features including stocks explorer, recommendations system, dashboard analytics, and subscription management  
**Goal:** Create a fully functional financial analytics platform with tier-aware features and compelling data visualization

## Justification for Phase 2 Structure

### Why This Approach?
1. **Data-Driven Features:** Build the core value proposition of the platform
2. **Tier Differentiation:** Implement the subscription model's value through feature gating
3. **User Experience:** Create compelling visualizations that demonstrate platform capabilities
4. **Revenue Foundation:** Enable subscription upgrades through preview and upgrade flows

### Critical Success Factors
- Rich, filterable data tables with excellent UX
- Tier-specific feature visibility and upgrade prompts
- Charts and visualizations that convey financial insights
- Smooth subscription management flow

---

## Step-by-Step Implementation

### Step 1: Stocks Data Management & API Integration

#### 1.1 Stocks Store (`src/stores/stocks.ts`)
```typescript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { stocksAPI } from '@/api/stocks'
import type { StockEvent, PaginatedResponse } from '@/types'

export interface StockFilters {
  ticker?: string
  company?: string
  brokerage?: string
  action?: string
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}

export const useStocksStore = defineStore('stocks', () => {
  // State
  const stocks = ref<StockEvent[]>([])
  const totalItems = ref(0)
  const currentPage = ref(1)
  const itemsPerPage = ref(50)
  const isLoading = ref(false)
  const filters = ref<StockFilters>({
    sort_by: 'event_time',
    sort_order: 'desc'
  })
  const stats = ref<{ total_stocks: number; last_updated: string } | null>(null)

  // Getters
  const totalPages = computed(() => Math.ceil(totalItems.value / itemsPerPage.value))
  const hasNextPage = computed(() => currentPage.value < totalPages.value)
  const hasPrevPage = computed(() => currentPage.value > 1)
  const offset = computed(() => (currentPage.value - 1) * itemsPerPage.value)

  // Actions
  async function fetchStocks() {
    isLoading.value = true
    try {
      const params = {
        ...filters.value,
        limit: itemsPerPage.value,
        offset: offset.value
      }
      
      const response: PaginatedResponse<StockEvent> = await stocksAPI.getStocks(params)
      stocks.value = response.data
      totalItems.value = response.pagination.total_items
      
      return response
    } finally {
      isLoading.value = false
    }
  }

  async function fetchStocksByTicker(ticker: string) {
    isLoading.value = true
    try {
      const response = await stocksAPI.getStocksByTicker(ticker)
      return response.data
    } finally {
      isLoading.value = false
    }
  }

  async function fetchStats() {
    try {
      const response = await stocksAPI.getStats()
      stats.value = response.data
      return response
    } catch (error) {
      console.error('Failed to fetch stats:', error)
    }
  }

  function updateFilters(newFilters: Partial<StockFilters>) {
    filters.value = { ...filters.value, ...newFilters }
    currentPage.value = 1 // Reset to first page when filters change
  }

  function setPage(page: number) {
    if (page >= 1 && page <= totalPages.value) {
      currentPage.value = page
    }
  }

  function clearFilters() {
    filters.value = {
      sort_by: 'event_time',
      sort_order: 'desc'
    }
    currentPage.value = 1
  }

  return {
    // State
    stocks: readonly(stocks),
    totalItems: readonly(totalItems),
    currentPage: readonly(currentPage),
    itemsPerPage: readonly(itemsPerPage),
    isLoading: readonly(isLoading),
    filters: readonly(filters),
    stats: readonly(stats),
    
    // Getters
    totalPages,
    hasNextPage,
    hasPrevPage,
    
    // Actions
    fetchStocks,
    fetchStocksByTicker,
    fetchStats,
    updateFilters,
    setPage,
    clearFilters
  }
})
```

#### 1.2 Stocks API Service (`src/api/stocks.ts`)
```typescript
import { apiClient } from './client'
import type { StockEvent, PaginatedResponse } from '@/types'

export interface StockQueryParams {
  ticker?: string
  company?: string
  brokerage?: string
  action?: string
  limit?: number
  offset?: number
  sort_by?: string
  sort_order?: string
}

export const stocksAPI = {
  async getStocks(params?: StockQueryParams): Promise<PaginatedResponse<StockEvent>> {
    const searchParams = new URLSearchParams()
    
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== '') {
          searchParams.append(key, value.toString())
        }
      })
    }
    
    const queryString = searchParams.toString()
    const url = queryString ? `/stocks?${queryString}` : '/stocks'
    
    return apiClient.get(url)
  },

  async getStocksByTicker(ticker: string): Promise<{ data: StockEvent[] }> {
    return apiClient.get(`/stocks/${ticker.toUpperCase()}`)
  },

  async getStats(): Promise<{ data: { total_stocks: number; last_updated: string } }> {
    return apiClient.get('/stocks/stats')
  }
}
```

### Step 2: Recommendations Management

#### 2.1 Recommendations Store (`src/stores/recommendations.ts`)
```typescript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { recommendationsAPI } from '@/api/recommendations'
import { useAuthStore } from './auth'
import type { Recommendation, RecommendationResponse } from '@/types'

export interface RecommendationFilters {
  limit?: number
  min_score?: number
  type?: string
  exclude?: string
}

export const useRecommendationsStore = defineStore('recommendations', () => {
  // State
  const recommendations = ref<Recommendation[]>([])
  const isLoading = ref(false)
  const meta = ref<RecommendationResponse['meta'] | null>(null)
  const filters = ref<RecommendationFilters>({
    limit: 10
  })
  const selectedRecommendation = ref<Recommendation | null>(null)
  const previewData = ref<any>(null)

  // Getters
  const authStore = useAuthStore()
  const maxRecommendations = computed(() => {
    const limits = { guest: 10, basic: 25, premium: 100 }
    return limits[authStore.userTier]
  })

  const availableFeatures = computed(() => meta.value?.features || [])
  const rateLimitRemaining = computed(() => meta.value?.rate_limit_remaining)

  // Actions
  async function fetchRecommendations() {
    isLoading.value = true
    try {
      // Ensure limit doesn't exceed tier maximum
      const adjustedFilters = {
        ...filters.value,
        limit: Math.min(filters.value.limit || 10, maxRecommendations.value)
      }
      
      const response = await recommendationsAPI.getRecommendations(adjustedFilters)
      recommendations.value = response.data
      meta.value = response.meta
      
      return response
    } finally {
      isLoading.value = false
    }
  }

  async function fetchRecommendationByTicker(ticker: string) {
    isLoading.value = true
    try {
      const response = await recommendationsAPI.getRecommendationByTicker(ticker)
      selectedRecommendation.value = response.data
      return response
    } finally {
      isLoading.value = false
    }
  }

  async function fetchPreviewForTicker(ticker: string) {
    if (!authStore.isAuthenticated) return
    
    try {
      const response = await recommendationsAPI.getPreview(ticker)
      previewData.value = response
      return response
    } catch (error) {
      console.error('Failed to fetch preview:', error)
    }
  }

  function updateFilters(newFilters: Partial<RecommendationFilters>) {
    filters.value = { ...filters.value, ...newFilters }
  }

  function clearSelectedRecommendation() {
    selectedRecommendation.value = null
    previewData.value = null
  }

  return {
    // State
    recommendations: readonly(recommendations),
    isLoading: readonly(isLoading),
    meta: readonly(meta),
    filters: readonly(filters),
    selectedRecommendation: readonly(selectedRecommendation),
    previewData: readonly(previewData),
    
    // Getters
    maxRecommendations,
    availableFeatures,
    rateLimitRemaining,
    
    // Actions
    fetchRecommendations,
    fetchRecommendationByTicker,
    fetchPreviewForTicker,
    updateFilters,
    clearSelectedRecommendation
  }
})
```

#### 2.2 Recommendations API Service (`src/api/recommendations.ts`)
```typescript
import { apiClient } from './client'
import type { RecommendationResponse } from '@/types'

export interface RecommendationQueryParams {
  limit?: number
  min_score?: number
  type?: string
  exclude?: string
}

export const recommendationsAPI = {
  async getRecommendations(params?: RecommendationQueryParams): Promise<RecommendationResponse> {
    const searchParams = new URLSearchParams()
    
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== '') {
          searchParams.append(key, value.toString())
        }
      })
    }
    
    const queryString = searchParams.toString()
    const url = queryString ? `/recommendations?${queryString}` : '/recommendations'
    
    return apiClient.get(url)
  },

  async getRecommendationByTicker(ticker: string): Promise<{ data: any; meta: any }> {
    return apiClient.get(`/recommendations/${ticker.toUpperCase()}`)
  },

  async getPreview(ticker: string): Promise<any> {
    return apiClient.get(`/recommendations/preview/${ticker.toUpperCase()}`)
  }
}
```

### Step 3: Core UI Components

#### 3.1 Data Table Component (`src/components/ui/DataTable.vue`)
```vue
<template>
  <div class="bg-white shadow-sm rounded-lg border border-gray-200">
    <!-- Table Header -->
    <div class="px-6 py-4 border-b border-gray-200">
      <div class="flex items-center justify-between">
        <h3 class="text-lg font-medium text-gray-900">{{ title }}</h3>
        <div class="flex items-center space-x-2">
          <slot name="actions" />
        </div>
      </div>
    </div>

    <!-- Table Content -->
    <div class="overflow-x-auto">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th
              v-for="column in columns"
              :key="column.key"
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
              @click="handleSort(column.key)"
            >
              <div class="flex items-center space-x-1">
                <span>{{ column.label }}</span>
                <component
                  :is="getSortIcon(column.key)"
                  v-if="column.sortable"
                  class="h-4 w-4"
                />
              </div>
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-if="isLoading" v-for="n in 5" :key="n" class="animate-pulse">
            <td v-for="column in columns" :key="column.key" class="px-6 py-4">
              <div class="h-4 bg-gray-200 rounded"></div>
            </td>
          </tr>
          <tr v-else-if="data.length === 0">
            <td :colspan="columns.length" class="px-6 py-8 text-center text-gray-500">
              <slot name="empty">
                <div class="text-center">
                  <p class="text-sm text-gray-500">{{ emptyMessage }}</p>
                </div>
              </slot>
            </td>
          </tr>
          <tr
            v-else
            v-for="(item, index) in data"
            :key="item.id || index"
            class="hover:bg-gray-50 cursor-pointer"
            @click="$emit('rowClick', item)"
          >
            <td
              v-for="column in columns"
              :key="column.key"
              class="px-6 py-4 whitespace-nowrap text-sm"
            >
              <slot :name="`cell-${column.key}`" :item="item" :value="item[column.key]">
                <span :class="getCellClass(column.key, item[column.key])">
                  {{ formatCellValue(column, item[column.key]) }}
                </span>
              </slot>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <div v-if="pagination" class="px-6 py-4 border-t border-gray-200">
      <TablePagination
        :current-page="pagination.page"
        :total-pages="pagination.total_pages"
        :total-items="pagination.total_items"
        :has-next="pagination.has_next"
        :has-prev="pagination.has_prev"
        @page-change="$emit('pageChange', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { ChevronUpIcon, ChevronDownIcon, ChevronUpDownIcon } from '@heroicons/vue/24/outline'
import TablePagination from './TablePagination.vue'

interface Column {
  key: string
  label: string
  sortable?: boolean
  type?: 'text' | 'number' | 'date' | 'currency' | 'percentage'
  format?: (value: any) => string
}

interface Props {
  title: string
  columns: Column[]
  data: any[]
  isLoading?: boolean
  emptyMessage?: string
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
  pagination?: {
    page: number
    total_pages: number
    total_items: number
    has_next: boolean
    has_prev: boolean
  }
}

const props = withDefaults(defineProps<Props>(), {
  isLoading: false,
  emptyMessage: 'No data available',
  sortOrder: 'desc'
})

const emit = defineEmits<{
  sort: [{ column: string; order: 'asc' | 'desc' }]
  rowClick: [item: any]
  pageChange: [page: number]
}>()

function handleSort(column: string) {
  const currentOrder = props.sortBy === column ? props.sortOrder : 'desc'
  const newOrder = currentOrder === 'asc' ? 'desc' : 'asc'
  emit('sort', { column, order: newOrder })
}

function getSortIcon(column: string) {
  if (props.sortBy !== column) return ChevronUpDownIcon
  return props.sortOrder === 'asc' ? ChevronUpIcon : ChevronDownIcon
}

function formatCellValue(column: Column, value: any): string {
  if (value === null || value === undefined) return '-'
  
  if (column.format) {
    return column.format(value)
  }
  
  switch (column.type) {
    case 'date':
      return new Date(value).toLocaleDateString()
    case 'currency':
      return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'USD'
      }).format(value)
    case 'percentage':
      return `${(value * 100).toFixed(2)}%`
    case 'number':
      return typeof value === 'number' ? value.toLocaleString() : value
    default:
      return value.toString()
  }
}

function getCellClass(column: string, value: any): string {
  const baseClass = 'text-gray-900'
  
  // Add specific styling for financial data
  if (column.includes('target') || column.includes('price')) {
    if (value > 0) return `${baseClass} text-financial-buy`
    if (value < 0) return `${baseClass} text-financial-sell`
  }
  
  if (column.includes('rating')) {
    if (value?.toLowerCase().includes('buy')) return `${baseClass} text-financial-buy font-medium`
    if (value?.toLowerCase().includes('sell')) return `${baseClass} text-financial-sell font-medium`
    if (value?.toLowerCase().includes('hold')) return `${baseClass} text-financial-hold font-medium`
  }
  
  return baseClass
}
</script>
```

#### 3.2 Table Pagination (`src/components/ui/TablePagination.vue`)
```vue
<template>
  <div class="flex items-center justify-between">
    <div class="flex items-center text-sm text-gray-700">
      <span>
        Showing {{ startItem }} to {{ endItem }} of {{ totalItems }} results
      </span>
    </div>
    
    <div class="flex items-center space-x-2">
      <button
        :disabled="!hasPrev"
        @click="$emit('pageChange', currentPage - 1)"
        class="inline-flex items-center px-3 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-500 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        Previous
      </button>
      
      <div class="flex items-center space-x-1">
        <template v-for="page in visiblePages" :key="page">
          <button
            v-if="page !== '...'"
            @click="$emit('pageChange', page)"
            :class="[
              'inline-flex items-center px-3 py-2 border text-sm font-medium rounded-md',
              page === currentPage
                ? 'border-primary-500 bg-primary-50 text-primary-600'
                : 'border-gray-300 bg-white text-gray-500 hover:bg-gray-50'
            ]"
          >
            {{ page }}
          </button>
          <span v-else class="px-2 text-gray-500">...</span>
        </template>
      </div>
      
      <button
        :disabled="!hasNext"
        @click="$emit('pageChange', currentPage + 1)"
        class="inline-flex items-center px-3 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-500 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        Next
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  currentPage: number
  totalPages: number
  totalItems: number
  hasNext: boolean
  hasPrev: boolean
}

const props = defineProps<Props>()

defineEmits<{
  pageChange: [page: number]
}>()

const startItem = computed(() => {
  return ((props.currentPage - 1) * 50) + 1
})

const endItem = computed(() => {
  return Math.min(props.currentPage * 50, props.totalItems)
})

const visiblePages = computed(() => {
  const pages: (number | string)[] = []
  const current = props.currentPage
  const total = props.totalPages
  
  if (total <= 7) {
    for (let i = 1; i <= total; i++) {
      pages.push(i)
    }
  } else {
    if (current <= 4) {
      for (let i = 1; i <= 5; i++) {
        pages.push(i)
      }
      pages.push('...')
      pages.push(total)
    } else if (current >= total - 3) {
      pages.push(1)
      pages.push('...')
      for (let i = total - 4; i <= total; i++) {
        pages.push(i)
      }
    } else {
      pages.push(1)
      pages.push('...')
      for (let i = current - 1; i <= current + 1; i++) {
        pages.push(i)
      }
      pages.push('...')
      pages.push(total)
    }
  }
  
  return pages
})
</script>
```

### Step 4: Stocks Explorer Module

#### 4.1 Stocks View (`src/views/StocksView.vue`)
```vue
<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div class="sm:flex sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Stock Events</h1>
        <p class="mt-2 text-sm text-gray-700">
          Track broker actions, rating changes, and price targets
        </p>
      </div>
      <div class="mt-4 sm:mt-0">
        <StockFilters
          :filters="stocksStore.filters"
          :is-loading="stocksStore.isLoading"
          @update-filters="handleFiltersUpdate"
          @clear-filters="stocksStore.clearFilters"
        />
      </div>
    </div>

    <!-- Stats Cards -->
    <div v-if="stocksStore.stats" class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div class="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <div class="w-8 h-8 bg-primary-100 rounded-full flex items-center justify-center">
              <span class="text-primary-600 font-semibold">{{ formatNumber(stocksStore.stats.total_stocks) }}</span>
            </div>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-500">Total Stock Events</p>
            <p class="text-2xl font-semibold text-gray-900">{{ stocksStore.stats.total_stocks.toLocaleString() }}</p>
          </div>
        </div>
      </div>
      
      <div class="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <div class="w-8 h-8 bg-green-100 rounded-full flex items-center justify-center">
              <span class="text-green-600 font-semibold">🕒</span>
            </div>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-500">Last Updated</p>
            <p class="text-2xl font-semibold text-gray-900">{{ formatDate(stocksStore.stats.last_updated) }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Data Table -->
    <DataTable
      title="Stock Events"
      :columns="stockColumns"
      :data="stocksStore.stocks"
      :is-loading="stocksStore.isLoading"
      :sort-by="stocksStore.filters.sort_by"
      :sort-order="stocksStore.filters.sort_order"
      :pagination="{
        page: stocksStore.currentPage,
        total_pages: stocksStore.totalPages,
        total_items: stocksStore.totalItems,
        has_next: stocksStore.hasNextPage,
        has_prev: stocksStore.hasPrevPage
      }"
      empty-message="No stock events found with current filters"
      @sort="handleSort"
      @page-change="stocksStore.setPage"
      @row-click="handleRowClick"
    >
      <template #cell-target_change="{ item }">
        <span :class="getTargetChangeClass(item)">
          {{ formatTargetChange(item) }}
        </span>
      </template>

      <template #cell-rating_change="{ item }">
        <div class="flex flex-col">
          <span class="text-xs text-gray-500">{{ item.rating_from }}</span>
          <span class="text-xs">→</span>
          <span class="text-xs font-medium">{{ item.rating_to }}</span>
        </div>
      </template>

      <template #cell-event_time="{ value }">
        <div class="text-sm">
          <div class="font-medium">{{ formatDate(value) }}</div>
          <div class="text-gray-500">{{ formatTime(value) }}</div>
        </div>
      </template>
    </DataTable>

    <!-- Stock Detail Modal -->
    <StockDetailModal
      :is-open="showDetailModal"
      :ticker="selectedTicker"
      @close="showDetailModal = false"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useStocksStore } from '@/stores/stocks'
import DataTable from '@/components/ui/DataTable.vue'
import StockFilters from '@/components/features/stocks/StockFilters.vue'
import StockDetailModal from '@/components/features/stocks/StockDetailModal.vue'
import { format } from 'date-fns'
import type { StockEvent } from '@/types'

const stocksStore = useStocksStore()
const showDetailModal = ref(false)
const selectedTicker = ref('')

const stockColumns = [
  { key: 'ticker', label: 'Ticker', sortable: true },
  { key: 'company', label: 'Company', sortable: true },
  { key: 'brokerage', label: 'Brokerage', sortable: true },
  { key: 'action', label: 'Action', sortable: true },
  { key: 'rating_change', label: 'Rating Change' },
  { key: 'target_change', label: 'Target Change' },
  { key: 'event_time', label: 'Event Time', sortable: true }
]

onMounted(async () => {
  await Promise.all([
    stocksStore.fetchStocks(),
    stocksStore.fetchStats()
  ])
})

function handleFiltersUpdate() {
  stocksStore.fetchStocks()
}

function handleSort({ column, order }: { column: string; order: 'asc' | 'desc' }) {
  stocksStore.updateFilters({ sort_by: column, sort_order: order })
  stocksStore.fetchStocks()
}

function handleRowClick(stock: StockEvent) {
  selectedTicker.value = stock.ticker
  showDetailModal.value = true
}

function formatTargetChange(item: StockEvent): string {
  const from = item.target_from
  const to = item.target_to
  if (!from || !to) return '-'
  
  const change = to - from
  const percentage = ((change / from) * 100).toFixed(1)
  return `$${from} → $${to} (${change > 0 ? '+' : ''}${percentage}%)`
}

function getTargetChangeClass(item: StockEvent): string {
  const from = item.target_from
  const to = item.target_to
  if (!from || !to) return 'text-gray-500'
  
  const change = to - from
  if (change > 0) return 'text-financial-buy font-medium'
  if (change < 0) return 'text-financial-sell font-medium'
  return 'text-gray-500'
}

function formatDate(date: string): string {
  return format(new Date(date), 'MMM dd, yyyy')
}

function formatTime(date: string): string {
  return format(new Date(date), 'HH:mm')
}

function formatNumber(num: number): string {
  if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M'
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
  return num.toString()
}
</script>
```

#### 4.2 Stock Filters Component (`src/components/features/stocks/StockFilters.vue`)
```vue
<template>
  <div class="bg-white p-4 rounded-lg shadow-sm border border-gray-200">
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      <!-- Ticker Filter -->
      <div>
        <label for="ticker" class="block text-sm font-medium text-gray-700 mb-1">
          Ticker
        </label>
        <input
          id="ticker"
          v-model="localFilters.ticker"
          type="text"
          placeholder="e.g., AAPL"
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500"
          @input="debouncedUpdate"
        />
      </div>

      <!-- Company Filter -->
      <div>
        <label for="company" class="block text-sm font-medium text-gray-700 mb-1">
          Company
        </label>
        <input
          id="company"
          v-model="localFilters.company"
          type="text"
          placeholder="e.g., Apple Inc."
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500"
          @input="debouncedUpdate"
        />
      </div>

      <!-- Brokerage Filter -->
      <div>
        <label for="brokerage" class="block text-sm font-medium text-gray-700 mb-1">
          Brokerage
        </label>
        <select
          id="brokerage"
          v-model="localFilters.brokerage"
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500"
          @change="handleUpdate"
        >
          <option value="">All Brokerages</option>
          <option value="Goldman Sachs">Goldman Sachs</option>
          <option value="Morgan Stanley">Morgan Stanley</option>
          <option value="JP Morgan">JP Morgan</option>
          <option value="Bank of America">Bank of America</option>
          <option value="Citigroup">Citigroup</option>
        </select>
      </div>

      <!-- Action Filter -->
      <div>
        <label for="action" class="block text-sm font-medium text-gray-700 mb-1">
          Action
        </label>
        <select
          id="action"
          v-model="localFilters.action"
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500"
          @change="handleUpdate"
        >
          <option value="">All Actions</option>
          <option value="upgraded by">Upgraded</option>
          <option value="downgraded by">Downgraded</option>
          <option value="initiated by">Initiated</option>
          <option value="maintained by">Maintained</option>
        </select>
      </div>
    </div>

    <!-- Filter Actions -->
    <div class="mt-4 flex items-center justify-between">
      <div class="text-sm text-gray-500">
        {{ activeFiltersCount }} active filter{{ activeFiltersCount !== 1 ? 's' : '' }}
      </div>
      <div class="space-x-2">
        <button
          @click="handleClear"
          class="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
        >
          Clear Filters
        </button>
        <button
          :disabled="isLoading"
          @click="handleUpdate"
          class="inline-flex items-center px-3 py-2 border border-transparent text-sm leading-4 font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50"
        >
          <span v-if="!isLoading">Apply Filters</span>
          <span v-else>Applying...</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { debounce } from 'lodash-es'
import type { StockFilters } from '@/stores/stocks'

interface Props {
  filters: StockFilters
  isLoading: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  updateFilters: []
  clearFilters: []
}>()

const localFilters = ref<StockFilters>({ ...props.filters })

// Watch for external filter changes
watch(() => props.filters, (newFilters) => {
  localFilters.value = { ...newFilters }
}, { deep: true })

const activeFiltersCount = computed(() => {
  let count = 0
  if (localFilters.value.ticker) count++
  if (localFilters.value.company) count++
  if (localFilters.value.brokerage) count++
  if (localFilters.value.action) count++
  return count
})

const debouncedUpdate = debounce(() => {
  handleUpdate()
}, 500)

function handleUpdate() {
  // Update the store filters and emit update event
  Object.assign(props.filters, localFilters.value)
  emit('updateFilters')
}

function handleClear() {
  localFilters.value = {
    sort_by: 'event_time',
    sort_order: 'desc'
  }
  emit('clearFilters')
}
</script>
```

### Step 5: Recommendations Module

#### 5.1 Recommendations View (`src/views/RecommendationsView.vue`)
```vue
<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div class="sm:flex sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Stock Recommendations</h1>
        <p class="mt-2 text-sm text-gray-700">
          AI-powered insights based on broker actions and market data
        </p>
      </div>
      <TierBadge :tier="authStore.userTier" />
    </div>

    <!-- Tier Info Card -->
    <RecommendationTierInfo
      :user-tier="authStore.userTier"
      :available-features="recommendationsStore.availableFeatures"
      :rate-limit-remaining="recommendationsStore.rateLimitRemaining"
      :max-recommendations="recommendationsStore.maxRecommendations"
    />

    <!-- Filters -->
    <RecommendationFilters
      :filters="recommendationsStore.filters"
      :max-recommendations="recommendationsStore.maxRecommendations"
      @update-filters="handleFiltersUpdate"
    />

    <!-- Recommendations Grid -->
    <div v-if="recommendationsStore.isLoading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="n in 6" :key="n" class="animate-pulse">
        <div class="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
          <div class="h-4 bg-gray-200 rounded w-1/4 mb-4"></div>
          <div class="h-6 bg-gray-200 rounded w-3/4 mb-2"></div>
          <div class="h-4 bg-gray-200 rounded w-1/2 mb-4"></div>
          <div class="space-y-2">
            <div class="h-3 bg-gray-200 rounded"></div>
            <div class="h-3 bg-gray-200 rounded w-5/6"></div>
          </div>
        </div>
      </div>
    </div>

    <div v-else-if="recommendationsStore.recommendations.length === 0" class="text-center py-12">
      <div class="text-gray-500">
        <p class="text-lg font-medium">No recommendations found</p>
        <p class="text-sm mt-2">Try adjusting your filters or check back later</p>
      </div>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <RecommendationCard
        v-for="recommendation in recommendationsStore.recommendations"
        :key="recommendation.id"
        :recommendation="recommendation"
        :user-tier="authStore.userTier"
        @view-details="handleViewDetails"
        @upgrade="handleUpgrade"
      />
    </div>

    <!-- Recommendation Detail Modal -->
    <RecommendationDetailModal
      :is-open="showDetailModal"
      :recommendation="selectedRecommendation"
      :preview-data="recommendationsStore.previewData"
      :user-tier="authStore.userTier"
      @close="handleCloseDetail"
      @upgrade="handleUpgrade"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useRecommendationsStore } from '@/stores/recommendations'
import TierBadge from '@/components/ui/TierBadge.vue'
import RecommendationTierInfo from '@/components/features/recommendations/RecommendationTierInfo.vue'
import RecommendationFilters from '@/components/features/recommendations/RecommendationFilters.vue'
import RecommendationCard from '@/components/features/recommendations/RecommendationCard.vue'
import RecommendationDetailModal from '@/components/features/recommendations/RecommendationDetailModal.vue'
import type { Recommendation } from '@/types'

const router = useRouter()
const authStore = useAuthStore()
const recommendationsStore = useRecommendationsStore()

const showDetailModal = ref(false)
const selectedRecommendation = ref<Recommendation | null>(null)

onMounted(() => {
  recommendationsStore.fetchRecommendations()
})

function handleFiltersUpdate() {
  recommendationsStore.fetchRecommendations()
}

async function handleViewDetails(recommendation: Recommendation) {
  selectedRecommendation.value = recommendation
  
  // Fetch preview data if user is BASIC tier
  if (authStore.userTier === 'basic') {
    await recommendationsStore.fetchPreviewForTicker(recommendation.ticker)
  }
  
  showDetailModal.value = true
}

function handleCloseDetail() {
  showDetailModal.value = false
  selectedRecommendation.value = null
  recommendationsStore.clearSelectedRecommendation()
}

function handleUpgrade() {
  router.push('/subscription')
}
</script>
```

#### 5.2 Recommendation Card (`src/components/features/recommendations/RecommendationCard.vue`)
```vue
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
import { computed } from 'vue'
import RecommendationBadge from './RecommendationBadge.vue'
import type { Recommendation, UserTier } from '@/types'

interface Props {
  recommendation: Recommendation
  userTier: UserTier
}

const props = defineProps<Props>()

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
```

### Step 6: Dashboard Implementation

#### 6.1 Dashboard View (`src/views/DashboardView.vue`)
```vue
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
        :value="stocksStore.stats?.total_stocks"
        :loading="stocksStore.isLoading"
        icon="📊"
        color="blue"
      />
      <MetricCard
        title="Active Recommendations"
        :value="recommendationsStore.recommendations.length"
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
        :value="recommendationsStore.rateLimitRemaining"
        :loading="recommendationsStore.isLoading"
        icon="⚡"
        color="amber"
        :format="(v) => v ? `${v}/hr` : '-'"
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

const router = useRouter()
const authStore = useAuthStore()
const stocksStore = useStocksStore()
const recommendationsStore = useRecommendationsStore()

const strongBuyCount = computed(() => {
  return recommendationsStore.recommendations.filter(
    r => r.recommendation_type === 'Strong Buy'
  ).length
})

const topRecommendations = computed(() => {
  return recommendationsStore.recommendations
    .filter(r => ['Strong Buy', 'Buy'].includes(r.recommendation_type))
    .sort((a, b) => b.basic_score - a.basic_score)
    .slice(0, 5)
})

const recentStockEvents = computed(() => {
  return stocksStore.stocks.slice(0, 10)
})

const recommendationChartData = computed(() => {
  const distribution = recommendationsStore.recommendations.reduce((acc, rec) => {
    acc[rec.recommendation_type] = (acc[rec.recommendation_type] || 0) + 1
    return acc
  }, {} as Record<string, number>)

  return Object.entries(distribution).map(([type, count]) => ({
    name: type,
    value: count
  }))
})

onMounted(async () => {
  await Promise.all([
    stocksStore.fetchStocks(),
    stocksStore.fetchStats(),
    recommendationsStore.fetchRecommendations()
  ])
})

function handleViewRecommendation(ticker: string) {
  router.push({ name: 'recommendations', query: { ticker } })
}

function handleViewStock(ticker: string) {
  router.push({ name: 'stocks', query: { ticker } })
}
</script>
```

### Step 7: Subscription Management

#### 7.1 Subscription Store (`src/stores/subscription.ts`)
```typescript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { subscriptionAPI } from '@/api/subscriptions'
import { useAuthStore } from './auth'
import type { Subscription } from '@/types'

export const useSubscriptionStore = defineStore('subscription', () => {
  // State
  const currentSubscription = ref<Subscription | null>(null)
  const isLoading = ref(false)
  const isProcessingPayment = ref(false)

  // Getters
  const hasActiveSubscription = computed(() => 
    currentSubscription.value?.status === 'active'
  )

  const subscriptionPlan = computed(() => currentSubscription.value?.plan)

  // Actions
  async function fetchCurrentSubscription() {
    const authStore = useAuthStore()
    if (!authStore.isAuthenticated) return

    isLoading.value = true
    try {
      const response = await subscriptionAPI.getCurrentSubscription()
      currentSubscription.value = response
    } catch (error) {
      // No active subscription is not an error
      currentSubscription.value = null
    } finally {
      isLoading.value = false
    }
  }

  async function createSubscription(plan: 'monthly' | 'yearly') {
    isLoading.value = true
    try {
      const response = await subscriptionAPI.createSubscription({ plan })
      currentSubscription.value = response
      return response
    } finally {
      isLoading.value = false
    }
  }

  async function processPayment(subscriptionId: string) {
    isProcessingPayment.value = true
    try {
      await subscriptionAPI.processPayment(subscriptionId)
      
      // Update user tier to premium
      const authStore = useAuthStore()
      if (authStore.user) {
        authStore.user.tier = 'premium'
      }
      
      // Refresh subscription data
      await fetchCurrentSubscription()
    } finally {
      isProcessingPayment.value = false
    }
  }

  return {
    // State
    currentSubscription: readonly(currentSubscription),
    isLoading: readonly(isLoading),
    isProcessingPayment: readonly(isProcessingPayment),
    
    // Getters
    hasActiveSubscription,
    subscriptionPlan,
    
    // Actions
    fetchCurrentSubscription,
    createSubscription,
    processPayment
  }
})
```

#### 7.2 Subscription API (`src/api/subscriptions.ts`)
```typescript
import { apiClient } from './client'
import type { Subscription } from '@/types'

export interface CreateSubscriptionRequest {
  plan: 'monthly' | 'yearly'
}

export const subscriptionAPI = {
  async getCurrentSubscription(): Promise<Subscription> {
    return apiClient.get('/subscriptions/current')
  },

  async createSubscription(data: CreateSubscriptionRequest): Promise<Subscription> {
    return apiClient.post('/subscriptions', data)
  },

  async processPayment(subscriptionId: string): Promise<{ message: string }> {
    return apiClient.post(`/subscriptions/${subscriptionId}/payment`)
  }
}
```

---

## Testing Strategy for Phase 2

### Component Testing
- Data table sorting and pagination
- Filter components with debounced updates
- Recommendation cards with tier-specific content
- Chart components with mock data

### Integration Testing
- Store actions with API calls
- Router navigation with authentication
- Tier-based feature access
- Subscription upgrade flow

### E2E Testing
- Complete user journey from stocks to recommendations
- Filter and search functionality
- Subscription creation and payment simulation
- Responsive behavior across devices

---

## Phase 2 Deliverables

✅ **Data Management:**
- Comprehensive stores for stocks, recommendations, and subscriptions
- Robust API integration with error handling
- Efficient caching and state management

✅ **Core Features:**
- Stocks explorer with advanced filtering and sorting
- Tier-aware recommendations system
- Interactive dashboard with charts and metrics
- Subscription management flow

✅ **UI/UX Excellence:**
- Responsive data tables with professional styling
- Compelling data visualizations
- Clear tier differentiation and upgrade prompts
- Consistent design system

✅ **Business Logic:**
- Feature gating based on user tiers
- Rate limiting awareness
- Preview functionality for upgrades
- Financial data formatting and color coding

---

## Next Steps to Phase 3

Phase 2 establishes the core platform functionality. Phase 3 will focus on:
- Advanced search and filtering capabilities
- Performance optimizations and caching
- Comprehensive AI placeholder integration
- Advanced data visualizations and charts
- Final polish, testing, and deployment preparation

The robust foundation built in Phases 1 and 2 enables Phase 3 to focus on enhancement and optimization rather than core functionality implementation. 