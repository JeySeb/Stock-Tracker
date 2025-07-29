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
        :data="stocksData"
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
  
  // Create a mutable copy of stocks for the DataTable
  const stocksData = computed(() => [...stocksStore.stocks])
  
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