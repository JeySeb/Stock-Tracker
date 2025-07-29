<template>
  <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
    <div class="flex items-center">
      <div class="flex-shrink-0">
        <div class="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
          <span class="text-blue-600 text-sm font-semibold">{{ userTier.charAt(0).toUpperCase() }}</span>
        </div>
      </div>
      <div class="ml-4">
        <h3 class="text-sm font-medium text-blue-900">{{ userTier.charAt(0).toUpperCase() + userTier.slice(1) }} Tier</h3>
        <p class="text-sm text-blue-700">
          {{ rateLimitRemaining }} requests remaining | 
          Up to {{ maxRecommendations }} recommendations
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { UserTier } from '@/types'

interface Props {
  userTier: UserTier
  availableFeatures?: string[]
  rateLimitRemaining?: number
  maxRecommendations?: number
}

withDefaults(defineProps<Props>(), {
  availableFeatures: () => [],
  rateLimitRemaining: 0,
  maxRecommendations: 10
})
</script>