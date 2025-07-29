<template>
  <div 
    class="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium"
    :class="[tierStyles[tier].bg, tierStyles[tier].text]"
  >
    <span class="mr-1">{{ tierStyles[tier].icon }}</span>
    {{ formatTier }}
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
    bg: 'bg-gray-100',
    text: 'text-gray-700',
    icon: '👋'
  },
  basic: {
    bg: 'bg-blue-100',
    text: 'text-blue-700',
    icon: '⭐'
  },
  premium: {
    bg: 'bg-primary-100',
    text: 'text-primary-700',
    icon: '✨'
  }
} as const

// Format tier name for display
const formatTier = computed(() => {
  return props.tier.charAt(0).toUpperCase() + props.tier.slice(1)
})
</script>

<style scoped>
/* Additional tier-specific styles can be added here if needed */
.bg-primary-100 {
  @apply bg-blue-100;
}

.text-primary-700 {
  @apply text-blue-700;
}
</style>