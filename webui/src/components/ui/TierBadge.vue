<template>
  <div 
    class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium transition-colors"
    :class="[tierStyles[tier].bg, tierStyles[tier].text, tierStyles[tier].border]"
  >
    <span class="mr-1 text-xs">{{ tierStyles[tier].icon }}</span>
    <span class="hidden sm:inline">{{ formatTier }}</span>
    <span class="sm:hidden">{{ tierStyles[tier].short }}</span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { UserTier } from '@/types'

interface Props {
  /** The user's subscription tier */
  tier: UserTier
}

const props = defineProps<Props>()

// Tier-specific styles and configurations
const tierStyles = {
  guest: {
    bg: 'bg-gray-50',
    text: 'text-gray-600',
    border: 'border border-gray-200',
    icon: '👋',
    short: 'Free'
  },
  basic: {
    bg: 'bg-blue-50',
    text: 'text-blue-700',
    border: 'border border-blue-200',
    icon: '⭐',
    short: 'Basic'
  },
  premium: {
    bg: 'bg-gradient-to-r from-primary-50 to-blue-50',
    text: 'text-primary-700',
    border: 'border border-primary-200',
    icon: '✨',
    short: 'Pro'
  }
} as const

// Format tier name for display
const formatTier = computed(() => {
  return props.tier.charAt(0).toUpperCase() + props.tier.slice(1)
})
</script>

<style scoped>
/* Additional tier-specific styles can be added here if needed */
.bg-primary-50 {
  @apply bg-blue-50;
}

.text-primary-700 {
  @apply text-blue-700;
}

.border-primary-200 {
  @apply border-blue-200;
}
</style>