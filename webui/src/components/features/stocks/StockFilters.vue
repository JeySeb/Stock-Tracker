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