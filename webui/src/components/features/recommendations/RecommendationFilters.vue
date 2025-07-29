<template>
  <div class="bg-white p-4 rounded-lg shadow-sm border border-gray-200">
    <div class="flex flex-wrap gap-4 items-center">
      <div class="flex-1 min-w-48">
        <label class="block text-sm font-medium text-gray-700 mb-1">Filter by Type</label>
        <select 
          v-model="localFilters.recommendation_type"
          @change="emitUpdate"
          class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
        >
          <option value="">All Types</option>
          <option value="Strong Buy">Strong Buy</option>
          <option value="Buy">Buy</option>
          <option value="Hold">Hold</option>
          <option value="Sell">Sell</option>
          <option value="Strong Sell">Strong Sell</option>
        </select>
      </div>
      
      <div class="flex-1 min-w-48">
        <label class="block text-sm font-medium text-gray-700 mb-1">Minimum Score</label>
        <input 
          v-model.number="localFilters.min_score"
          @input="emitUpdate"
          type="number"
          min="0"
          max="100"
          class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          placeholder="0"
        />
      </div>
      
      <div class="flex-1 min-w-48">
        <label class="block text-sm font-medium text-gray-700 mb-1">Limit</label>
        <select 
          v-model.number="localFilters.limit"
          @change="emitUpdate"
          class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
        >
          <option :value="10">10</option>
          <option :value="25">25</option>
          <option :value="50">50</option>
          <option :value="maxRecommendations">{{ maxRecommendations }}</option>
        </select>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

interface Props {
  filters: Record<string, any>
  maxRecommendations: number
}

const props = defineProps<Props>()
const emit = defineEmits<{
  updateFilters: [filters: Record<string, any>]
}>()

const localFilters = ref({
  recommendation_type: props.filters.recommendation_type || '',
  min_score: props.filters.min_score || 0,
  limit: props.filters.limit || 10
})

watch(() => props.filters, (newFilters) => {
  localFilters.value = { ...newFilters }
}, { deep: true })

function emitUpdate() {
  emit('updateFilters', localFilters.value)
}
</script>