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
            <template v-if="isLoading">
              <tr v-for="n in 5" :key="n" class="animate-pulse">
                <td v-for="column in columns" :key="column.key" class="px-6 py-4">
                  <div class="h-4 bg-gray-200 rounded"></div>
                </td>
              </tr>
            </template>
            <tr v-else-if="data.length === 0">
              <td :colspan="columns.length" class="px-6 py-8 text-center text-gray-500">
                <slot name="empty">
                  <div class="text-center">
                    <p class="text-sm text-gray-500">{{ emptyMessage }}</p>
                  </div>
                </slot>
              </td>
            </tr>
            <template v-else>
              <tr
                v-for="(item, index) in data"
                :key="(item as any)?.id || `row-${index}`"
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
            </template>
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
          :items-per-page="pagination.items_per_page"
          @page-change="$emit('pageChange', $event)"
        />
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ChevronUpIcon, ChevronDownIcon, ChevronUpDownIcon } from '@heroicons/vue/24/outline'
  import TablePagination from './TablePagination.vue'
  
  interface Column<T = Record<string, unknown>> {
    key: keyof T
    label: string
    sortable?: boolean
    type?: 'text' | 'number' | 'date' | 'currency' | 'percentage'
    format?: (value: T[keyof T]) => string
  }
  
  interface Props<T = Record<string, unknown>> {
    title: string
    columns: Column<T>[]
    data: T[]
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
      items_per_page: number
    }
  }
  
  const props = withDefaults(defineProps<Props>(), {
    isLoading: false,
    emptyMessage: 'No data available',
    sortOrder: 'desc'
  })
  
  const emit = defineEmits<{
    sort: [{ column: string; order: 'asc' | 'desc' }]
    rowClick: [item: Record<string, unknown>]
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
  
  function formatCellValue(column: Column, value: unknown): string {
    if (value === null || value === undefined) return '-'
    
    if (column.format) {
      return column.format(value as never)
    }
    
    switch (column.type) {
      case 'date':
        return new Date(value as string | number).toLocaleDateString()
      case 'currency':
        return new Intl.NumberFormat('en-US', {
          style: 'currency',
          currency: 'USD'
        }).format(value as number)
      case 'percentage':
        return `${((value as number) * 100).toFixed(2)}%`
      case 'number':
        return typeof value === 'number' ? value.toLocaleString() : String(value)
      default:
        return String(value)
    }
  }
  
  function getCellClass(column: string, value: unknown): string {
    const baseClass = 'text-gray-900'
    
    // Add specific styling for financial data
    if (column.includes('target') || column.includes('price')) {
      if (typeof value === 'number' && value > 0) return `${baseClass} text-financial-buy`
      if (typeof value === 'number' && value < 0) return `${baseClass} text-financial-sell`
    }
    
    if (column.includes('rating')) {
      const stringValue = String(value).toLowerCase()
      if (stringValue.includes('buy')) return `${baseClass} text-financial-buy font-medium`
      if (stringValue.includes('sell')) return `${baseClass} text-financial-sell font-medium`
      if (stringValue.includes('hold')) return `${baseClass} text-financial-hold font-medium`
    }
    
    return baseClass
  }
  </script>