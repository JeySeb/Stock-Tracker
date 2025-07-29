<template>
  <div 
    :class="[
      'relative bg-white rounded-lg shadow-sm border overflow-hidden transition-all duration-200',
      {
        'border-gray-200 hover:shadow-md': available,
        'border-gray-100 opacity-75': !available,
        'border-primary-200 bg-primary-50/50': premiumOnly,
        'border-amber-200 bg-amber-50/50': comingSoon
      }
    ]"
  >
    <!-- Premium Badge -->
    <div 
      v-if="premiumOnly" 
      class="absolute top-3 right-3 px-2 py-1 bg-primary-100 text-primary-700 text-xs font-medium rounded"
    >
      Premium
    </div>

    <!-- Coming Soon Badge -->
    <div 
      v-if="comingSoon" 
      class="absolute top-3 right-3 px-2 py-1 bg-amber-100 text-amber-700 text-xs font-medium rounded"
    >
      Coming Soon
    </div>

    <!-- Card Content -->
    <div 
      class="p-6"
      :class="{ 'cursor-pointer': isClickable }"
      @click="handleClick"
    >
      <div class="flex items-start space-x-4">
        <!-- Icon -->
        <div 
          class="flex-shrink-0 w-10 h-10 rounded-full flex items-center justify-center"
          :class="[
            available ? 'bg-primary-100 text-primary-600' : 'bg-gray-100 text-gray-400',
            { 'bg-amber-100 text-amber-600': comingSoon }
          ]"
        >
          <span class="text-xl">{{ icon }}</span>
        </div>

        <!-- Text Content -->
        <div class="flex-1 min-w-0">
          <h3 
            class="text-lg font-medium"
            :class="[
              available ? 'text-gray-900' : 'text-gray-500',
              { 'text-primary-900': premiumOnly }
            ]"
          >
            {{ title }}
          </h3>
          <p 
            class="mt-1 text-sm"
            :class="[
              available ? 'text-gray-500' : 'text-gray-400',
              { 'text-primary-600': premiumOnly }
            ]"
          >
            {{ description }}
          </p>
        </div>
      </div>

      <!-- Action Area -->
      <div v-if="showAction" class="mt-4">
        <template v-if="available && link">
          <RouterLink
            :to="link"
            class="inline-flex items-center text-sm font-medium text-primary-600 hover:text-primary-700"
          >
            Explore Feature
            <svg class="ml-1 w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </RouterLink>
        </template>
        <template v-else-if="!available && !comingSoon">
          <button
            @click="$emit('upgrade')"
            class="inline-flex items-center text-sm font-medium text-primary-600 hover:text-primary-700"
          >
            Upgrade to Access
            <svg class="ml-1 w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
            </svg>
          </button>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { RouteLocationRaw } from 'vue-router'

interface Props {
  /** Title of the feature */
  title: string
  /** Description of the feature */
  description: string
  /** Emoji or icon to display */
  icon?: string
  /** Whether the feature is available for the current user */
  available?: boolean
  /** Router link if the feature is accessible */
  link?: RouteLocationRaw
  /** Whether this is a premium-only feature */
  premiumOnly?: boolean
  /** Whether this feature is coming soon */
  comingSoon?: boolean
  /** Whether this is a placeholder for future AI features */
  placeholderForAi?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  icon: '✨',
  available: false,
  premiumOnly: false,
  comingSoon: false,
  placeholderForAi: false
})

const emit = defineEmits<{
  upgrade: []
}>()

// Computed
const isClickable = computed(() => {
  return (props.available && props.link) || (!props.available && !props.comingSoon)
})

const showAction = computed(() => {
  return !props.comingSoon && (props.link || !props.available)
})

// Methods
function handleClick() {
  if (!props.available && !props.comingSoon) {
    emit('upgrade')
  }
}
</script>