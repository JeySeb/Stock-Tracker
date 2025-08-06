import { defineStore } from 'pinia'
import { ref, computed, readonly } from 'vue'
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
  const itemsPerPage = ref(100) // Increased to get more stocks for heat map
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
    } catch (error) {
      console.error('Failed to fetch stocks:', error)
      // Fallback to empty array to prevent crashes
      stocks.value = []
      totalItems.value = 0
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
      // Fallback to mock data
      stats.value = {
        total_stocks: 0,
        last_updated: new Date().toISOString()
      }
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