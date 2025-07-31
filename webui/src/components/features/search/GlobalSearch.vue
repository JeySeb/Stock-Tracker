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
import type { StockEvent, Recommendation } from '@/types'

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
      const stockData = result.data as StockEvent
      router.push({ name: 'stocks', query: { ticker: stockData.ticker } })
      break
    case 'recommendation':
      const recData = result.data as Recommendation
      router.push({ name: 'recommendations', query: { ticker: recData.ticker } })
      break
    case 'company':
      const companyData = result.data as { name: string }
      router.push({ name: 'stocks', query: { company: companyData.name } })
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