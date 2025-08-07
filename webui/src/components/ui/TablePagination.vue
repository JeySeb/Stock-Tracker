<template>
    <div class="flex items-center justify-between">
      <div class="flex items-center text-sm text-gray-700">
        <span>
          Showing {{ startItem }} to {{ endItem }} of {{ totalItems }} results
        </span>
      </div>
      
      <div v-if="visiblePages.length > 0" class="flex items-center space-x-2">
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
              @click="$emit('pageChange', page as number)"
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
    itemsPerPage: number
  }
  
  const props = defineProps<Props>()
  
  defineEmits<{
    pageChange: [page: number]
  }>()
  
  const startItem = computed(() => {
    return ((props.currentPage - 1) * props.itemsPerPage) + 1
  })
  
  const endItem = computed(() => {
    return Math.min(props.currentPage * props.itemsPerPage, props.totalItems)
  })
  
  const visiblePages = computed(() => {
    const pages: (number | string)[] = []
    const current = props.currentPage
    const total = props.totalPages
    
    // If there's only one page, don't show pagination
    if (total <= 1) {
      return []
    }
    
    if (total <= 7) {
      // Show all pages if there are 7 or fewer
      for (let i = 1; i <= total; i++) {
        pages.push(i)
      }
    } else {
      if (current <= 4) {
        // Show first 5 pages + ellipsis + last page
        for (let i = 1; i <= 5; i++) {
          pages.push(i)
        }
        pages.push('...')
        pages.push(total)
      } else if (current >= total - 3) {
        // Show first page + ellipsis + last 5 pages
        pages.push(1)
        pages.push('...')
        for (let i = total - 4; i <= total; i++) {
          pages.push(i)
        }
      } else {
        // Show first page + ellipsis + current-1, current, current+1 + ellipsis + last page
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