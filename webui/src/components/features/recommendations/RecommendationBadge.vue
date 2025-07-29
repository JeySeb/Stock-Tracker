<template>
  <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium" :class="badgeClasses">
    {{ type }}
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'

type RecommendationType = 'Strong Buy' | 'Buy' | 'Hold' | 'Sell' | 'Strong Sell'

interface Props {
  type: RecommendationType
}

const props = defineProps<Props>()

const badgeClasses = computed(() => {
  switch (props.type) {
    case 'Strong Buy':
      return 'bg-green-100 text-green-800'
    case 'Buy':
      return 'bg-green-50 text-green-700'
    case 'Hold':
      return 'bg-yellow-100 text-yellow-800'
    case 'Sell':
      return 'bg-red-50 text-red-700'
    case 'Strong Sell':
      return 'bg-red-100 text-red-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
})
</script>