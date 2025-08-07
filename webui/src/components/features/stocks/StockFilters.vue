<template>
  <div class="bg-white rounded-lg shadow-sm border border-gray-200">
    <!-- Header with Active Filters -->
    <div class="p-4 border-b border-gray-100 bg-gradient-to-r from-gray-50 to-white">
      <div class="flex items-center justify-between mb-3">
        <div class="flex items-center space-x-3">
          <div class="flex-shrink-0">
            <div class="w-7 h-7 bg-primary-100 rounded-lg flex items-center justify-center">
              <svg class="w-4 h-4 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z" />
              </svg>
            </div>
          </div>
          <div>
            <h3 class="text-base font-semibold text-gray-900">Filters</h3>
            <p class="text-xs text-gray-500">Refine your stock event search</p>
          </div>
        </div>
        <div class="flex items-center space-x-2">
          <div v-if="hasActiveFilters" class="flex items-center space-x-2">
            <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-primary-100 text-primary-800">
              {{ activeFiltersCount }} active filter{{ activeFiltersCount !== 1 ? 's' : '' }}
            </span>
            <button
              @click="handleClear"
              class="inline-flex items-center px-2 py-1 text-xs font-medium text-red-600 hover:text-red-700 hover:bg-red-50 rounded-lg transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
            >
              <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
              Clear all
            </button>
          </div>
        </div>
      </div>
      
      <!-- Active Filters Display -->
      <div v-if="hasActiveFilters" class="flex flex-wrap gap-1.5">
        <template v-if="localFilters.tickers?.length">
          <span
            v-for="ticker in localFilters.tickers"
            :key="ticker"
            class="inline-flex items-center px-2 py-1 rounded-full text-xs bg-blue-100 text-blue-800 font-medium"
          >
            <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
            </svg>
            <span>{{ ticker }}</span>
            <button
              @click="removeFilter('tickers', ticker)"
              class="ml-1 text-blue-600 hover:text-blue-800 focus:outline-none"
            >
              <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </span>
        </template>

        <template v-if="localFilters.companies?.length">
          <span
            v-for="company in localFilters.companies"
            :key="company"
            class="inline-flex items-center px-2 py-1 rounded-full text-xs bg-indigo-100 text-indigo-800 font-medium"
          >
            <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
            </svg>
            <span>{{ company }}</span>
            <button
              @click="removeFilter('companies', company)"
              class="ml-1 text-indigo-600 hover:text-indigo-800 focus:outline-none"
            >
              <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </span>
        </template>

        <template v-if="localFilters.brokerages?.length">
          <span
            v-for="brokerage in localFilters.brokerages"
            :key="brokerage"
            class="inline-flex items-center px-2 py-1 rounded-full text-xs bg-green-100 text-green-800 font-medium"
          >
            <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
            </svg>
            <span>{{ brokerage }}</span>
            <button
              @click="removeFilter('brokerages', brokerage)"
              class="ml-1 text-green-600 hover:text-green-800 focus:outline-none"
            >
              <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </span>
        </template>

        <template v-if="localFilters.actions?.length">
          <span
            v-for="action in localFilters.actions"
            :key="action"
            class="inline-flex items-center px-2 py-1 rounded-full text-xs bg-purple-100 text-purple-800 font-medium"
          >
            <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
            </svg>
            <span>{{ formatAction(action) }}</span>
            <button
              @click="removeFilter('actions', action)"
              class="ml-1 text-purple-600 hover:text-purple-800 focus:outline-none"
            >
              <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </span>
        </template>
      </div>
    </div>

    <!-- Basic Filters -->
    <div class="p-4">
      <div class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-4 gap-4">
        <!-- Ticker Search -->
        <div class="space-y-1">
          <label class="block text-xs font-medium text-gray-700">Tickers</label>
          <MultiSelect
            v-model="tickersArray"
            :options="tickerOptions"
            placeholder="Search tickers..."
            class="w-full"
          />
        </div>

        <!-- Company Search -->
        <div class="space-y-1">
          <label class="block text-xs font-medium text-gray-700">Companies</label>
          <MultiSelect
            v-model="companiesArray"
            :options="companyOptions"
            placeholder="Search companies..."
            class="w-full"
          />
        </div>

        <!-- Brokerage Selection -->
        <div class="space-y-1">
          <label class="block text-xs font-medium text-gray-700">Brokerages</label>
          <MultiSelect
            v-model="brokeragesArray"
            :options="brokerageOptions"
            placeholder="Select brokerages..."
            class="w-full"
          />
        </div>

        <!-- Action Selection -->
        <div class="space-y-1">
          <label class="block text-xs font-medium text-gray-700">Actions</label>
          <MultiSelect
            v-model="actionsArray"
            :options="actionOptions"
            placeholder="Select actions..."
            class="w-full"
          />
        </div>
      </div>
    </div>

    <!-- Advanced Filters -->
    <div class="space-y-3">
      <!-- Advanced Filters -->
      <div v-show="showAdvancedFilters" class="border-t pt-3 space-y-3 px-4">
        <!-- Rating Filters -->
        <div>
          <h3 class="text-xs font-medium text-gray-900 mb-2">Rating Range</h3>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-2">
            <input
              v-model.number="localFilters.rating_from"
              type="number"
              min="1"
              max="5"
              step="0.1"
              placeholder="Min rating"
              class="w-full px-2 py-1.5 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500 text-xs"
              @input="debouncedUpdate"
            />
            <input
              v-model.number="localFilters.rating_to"
              type="number"
              min="1"
              max="5"
              step="0.1"
              placeholder="Max rating"
              class="w-full px-2 py-1.5 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500 text-xs"
              @input="debouncedUpdate"
            />
          </div>
        </div>

        <!-- Target Price Filters -->
        <div>
          <h3 class="text-xs font-medium text-gray-900 mb-2">Target Price Range</h3>
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-2">
            <input
              v-model.number="localFilters.target_from"
              type="number"
              min="0"
              step="0.01"
              placeholder="Min target"
              class="w-full px-2 py-1.5 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500 text-xs"
              @input="debouncedUpdate"
            />
            <input
              v-model.number="localFilters.target_to"
              type="number"
              min="0"
              step="0.01"
              placeholder="Max target"
              class="w-full px-2 py-1.5 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500 text-xs"
              @input="debouncedUpdate"
            />
            <input
              v-model.number="localFilters.min_target_change"
              type="number"
              step="0.1"
              placeholder="Min change %"
              class="w-full px-2 py-1.5 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500 text-xs"
              @input="debouncedUpdate"
            />
            <input
              v-model.number="localFilters.max_target_change"
              type="number"
              step="0.1"
              placeholder="Max change %"
              class="w-full px-2 py-1.5 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500 text-xs"
              @input="debouncedUpdate"
            />
          </div>
        </div>

        <!-- Time Filters -->
        <div>
          <h3 class="text-xs font-medium text-gray-900 mb-2">Time Range</h3>
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-2">
            <input
              v-model.number="localFilters.last_hours"
              type="number"
              min="1"
              placeholder="Hours"
              class="w-full px-2 py-1.5 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500 text-xs"
              @input="debouncedUpdate"
            />
            <input
              v-model.number="localFilters.last_days"
              type="number"
              min="1"
              placeholder="Days"
              class="w-full px-2 py-1.5 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500 text-xs"
              @input="debouncedUpdate"
            />
            <input
              v-model.number="localFilters.last_weeks"
              type="number"
              min="1"
              placeholder="Weeks"
              class="w-full px-2 py-1.5 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500 text-xs"
              @input="debouncedUpdate"
            />
            <input
              v-model.number="localFilters.last_months"
              type="number"
              min="1"
              placeholder="Months"
              class="w-full px-2 py-1.5 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500 text-xs"
              @input="debouncedUpdate"
            />
          </div>
        </div>

        <!-- Date Range Filters -->
        <div>
          <h3 class="text-xs font-medium text-gray-900 mb-2">Date Range</h3>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-2">
            <input
              v-model="localFilters.date_from"
              type="datetime-local"
              class="w-full px-2 py-1.5 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500 text-xs"
              @input="debouncedUpdate"
            />
            <input
              v-model="localFilters.date_to"
              type="datetime-local"
              class="w-full px-2 py-1.5 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500 text-xs"
              @input="debouncedUpdate"
            />
          </div>
        </div>
      </div>
    </div>
    <!-- Filter Actions -->
    <div class="p-4 border-t border-gray-100 bg-gray-50">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
        <!-- Left side: Rows per page and Advanced filters toggle -->
        <div class="flex flex-col sm:flex-row sm:items-center gap-3">
          <!-- Rows per page selector -->
          <div class="flex items-center space-x-2">
            <label class="text-xs font-medium text-gray-700">Rows per page:</label>
            <div class="relative">
              <select
                :value="itemsPerPage"
                @change="handleItemsPerPageChange"
                class="form-select block w-full px-3 py-1.5 text-xs font-medium text-gray-900 bg-white border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500"
                style="min-width: 60px"
              >
              
              <option value="10">10</option>
                <option value="25">25</option>
                <option value="50">50</option>
                <option value="100">100</option>
                <option value="200">200</option>
              </select>
            </div>
          </div>
          <!-- Divider -->
          <div class="hidden sm:block w-px h-4 bg-gray-300"></div>
          
          <!-- Advanced filters toggle -->
          <button
            @click="showAdvancedFilters = !showAdvancedFilters"
            class="inline-flex items-center px-2 py-1.5 text-xs font-medium text-primary-600 hover:text-primary-700 hover:bg-primary-50 rounded-md transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
          >
            <svg class="w-3 h-3 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z" />
            </svg>
            {{ showAdvancedFilters ? 'Hide Advanced Filters' : 'Show Advanced Filters' }}
          </button>
        </div>
        
        <!-- Right side: Action buttons -->
        <div class="flex items-center space-x-2">
          <button
            @click="handleClear"
            class="inline-flex items-center px-3 py-1.5 border border-gray-300 shadow-sm text-xs font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-colors duration-200"
          >
            <svg class="w-3 h-3 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
            Clear
          </button>
          <button
            :disabled="isLoading"
            @click="handleUpdate"
            class="inline-flex items-center px-4 py-1.5 border border-transparent text-xs font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200 shadow-sm"
          >
            <svg v-if="!isLoading" class="w-3 h-3 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z" />
            </svg>
            <svg v-else class="animate-spin w-3 h-3 mr-1.5" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span v-if="!isLoading">Apply Filters</span>
            <span v-else>Applying...</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
  
  <script setup lang="ts">
  import { ref, computed, watch, onMounted } from 'vue'
import { debounce } from 'lodash-es'
import type { StockFilters } from '@/stores/stocks'
import MultiSelect from '@/components/ui/MultiSelect.vue'
import { stocksAPI } from '@/api/stocks'

interface Props {
  filters: StockFilters
  isLoading: boolean
  itemsPerPage: number
}

const props = defineProps<Props>()

const emit = defineEmits<{
  updateFilters: [filters: StockFilters]
  clearFilters: []
  itemsPerPageChange: [items: number]
}>()

const localFilters = ref<StockFilters>({ 
  ...props.filters,
  tickers: props.filters.tickers || [],
  companies: props.filters.companies || [],
  brokerages: props.filters.brokerages || [],
  actions: props.filters.actions || []
})
const showAdvancedFilters = ref(false)

// Watch for external filter changes
watch(() => props.filters, (newFilters) => {
  localFilters.value = { 
    ...newFilters,
    tickers: newFilters.tickers || [],
    companies: newFilters.companies || [],
    brokerages: newFilters.brokerages || [],
    actions: newFilters.actions || []
  }
}, { deep: true })

// Predefined options for dropdowns
const brokerageOptions = [
  { label: 'Goldman Sachs', value: 'Goldman Sachs' },
  { label: 'Morgan Stanley', value: 'Morgan Stanley' },
  { label: 'JP Morgan', value: 'JP Morgan' },
  { label: 'Bank of America', value: 'Bank of America' },
  { label: 'Citigroup', value: 'Citigroup' },
  { label: 'Wells Fargo', value: 'Wells Fargo' },
  { label: 'UBS', value: 'UBS' },
  { label: 'Deutsche Bank', value: 'Deutsche Bank' }
]

const actionOptions = [
  { label: 'Upgraded', value: 'upgraded by' },
  { label: 'Downgraded', value: 'downgraded by' },
  { label: 'Initiated', value: 'initiated by' },
  { label: 'Maintained', value: 'maintained by' },
  { label: 'Resumed', value: 'resumed by' },
  { label: 'Suspended', value: 'suspended by' }
]

// Computed properties for arrays to ensure they're never undefined
const tickersArray = computed({
  get: () => localFilters.value.tickers || [],
  set: (value: string[]) => { localFilters.value.tickers = value }
})
const companiesArray = computed({
  get: () => localFilters.value.companies || [],
  set: (value: string[]) => { localFilters.value.companies = value }
})
const brokeragesArray = computed({
  get: () => localFilters.value.brokerages || [],
  set: (value: string[]) => { localFilters.value.brokerages = value }
})
const actionsArray = computed({
  get: () => localFilters.value.actions || [],
  set: (value: string[]) => { localFilters.value.actions = value }
})

// Real data options from API
const tickerOptions = ref<{ label: string; value: string }[]>([])
const companyOptions = ref<{ label: string; value: string }[]>([])

// Fetch options from API
async function fetchOptions() {
  try {
    const [tickersResponse, companiesResponse] = await Promise.all([
      stocksAPI.getUniqueTickers(),
      stocksAPI.getUniqueCompanies()
    ])
    
    tickerOptions.value = tickersResponse.data.map(ticker => ({ label: ticker, value: ticker }))
    companyOptions.value = companiesResponse.data.map(company => ({ label: company, value: company }))
  } catch (error) {
    console.error('Failed to fetch options:', error)
  }
}

// Fetch options on component mount
onMounted(() => {
  fetchOptions()
})

// Computed properties for UI state
const hasActiveFilters = computed(() => activeFiltersCount.value > 0)

const activeFiltersCount = computed(() => {
  let count = 0
  const filters = localFilters.value

  // Array filters
  if (filters.tickers?.length) count += filters.tickers.length
  if (filters.companies?.length) count += filters.companies.length
  if (filters.brokerages?.length) count += filters.brokerages.length
  if (filters.actions?.length) count += filters.actions.length

  // Rating filters
  if (filters.rating_from !== undefined) count++
  if (filters.rating_to !== undefined) count++

  // Target price filters
  if (filters.target_from !== undefined) count++
  if (filters.target_to !== undefined) count++
  if (filters.min_target_change !== undefined) count++
  if (filters.max_target_change !== undefined) count++
  if (filters.has_target_price !== undefined) count++

  // Brokerage score filters
  if (filters.min_broker_score !== undefined) count++
  if (filters.max_broker_score !== undefined) count++

  // Time filters
  if (filters.last_hours !== undefined) count++
  if (filters.last_days !== undefined) count++
  if (filters.last_weeks !== undefined) count++
  if (filters.last_months !== undefined) count++

  // Date filters
  if (filters.date_from) count++
  if (filters.date_to) count++
  if (filters.date_ranges) count++

  return count
})

// Helper functions
function formatAction(action: string): string {
  return action
    .replace(' by', '')
    .split('_')
    .map(word => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ')
}

function removeFilter(filterType: keyof StockFilters, value: string) {
  const filters = localFilters.value[filterType] as string[]
  if (Array.isArray(filters)) {
    const index = filters.indexOf(value)
    if (index !== -1) {
      filters.splice(index, 1)
      handleUpdate()
    }
  }
}

const debouncedUpdate = debounce(() => {
  handleUpdate()
}, 500)

function handleUpdate() {
  // Convert date inputs to RFC3339 format
  const updatedFilters = { ...localFilters.value }
  
  if (updatedFilters.date_from) {
    updatedFilters.date_from = new Date(updatedFilters.date_from).toISOString()
  }
  if (updatedFilters.date_to) {
    updatedFilters.date_to = new Date(updatedFilters.date_to).toISOString()
  }

  // Emit update event with the updated filters
  emit('updateFilters', updatedFilters)
}

function handleClear() {
  localFilters.value = {
    sort_by: 'event_time',
    sort_order: 'desc',
    tickers: [],
    companies: [],
    brokerages: [],
    actions: []
  }
  showAdvancedFilters.value = false
  emit('clearFilters')
}

function handleItemsPerPageChange(event: Event) {
  const target = event.target as HTMLSelectElement
  const newItemsPerPage = parseInt(target.value, 10)
  emit('itemsPerPageChange', newItemsPerPage)
}
  </script>