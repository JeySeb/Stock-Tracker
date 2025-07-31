<template>
  <div class="relative border rounded-lg shadow-sm divide-y divide-gray-200 bg-white" :class="borderClass">
    <!-- Popular badge -->
    <div v-if="isPopular" class="absolute top-0 right-6 transform -translate-y-1/2">
      <span class="inline-flex rounded-full bg-primary-100 px-4 py-1 text-sm font-semibold text-primary-600">
        {{ popularText }}
      </span>
    </div>
    
    <!-- Header -->
    <div class="p-6">
      <h3 class="text-lg leading-6 font-medium text-gray-900">{{ title }}</h3>
      <p class="mt-4 text-sm text-gray-500">{{ description }}</p>
      
      <!-- Price -->
      <div class="mt-8">
        <span class="text-4xl font-extrabold text-gray-900">${{ price }}</span>
        <span v-if="period" class="text-base font-medium text-gray-500">{{ period }}</span>
      </div>
      
      <!-- Savings text -->
      <p v-if="savingsText" class="text-sm text-green-600 mt-2">{{ savingsText }}</p>
      
      <!-- CTA Button -->
      <component
        :is="buttonComponent"
        :to="buttonTo"
        @click="$emit('click')"
        :disabled="isDisabled"
        :class="buttonClass"
        class="mt-8 block w-full border border-transparent rounded-md py-2 text-sm font-semibold text-center"
      >
        <span v-if="isLoading">{{ loadingText }}</span>
        <span v-else>{{ buttonText }}</span>
      </component>
    </div>
    
    <!-- Features List -->
    <div class="pt-6 pb-8 px-6">
      <h4 class="text-xs font-medium text-gray-900 tracking-wide uppercase">{{ featuresTitle }}</h4>
      <ul class="mt-6 space-y-4">
        <li v-for="feature in features" :key="feature" class="flex space-x-3">
          <svg class="flex-shrink-0 h-5 w-5 text-green-500" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
          </svg>
          <span class="text-sm text-gray-500">{{ feature }}</span>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  title: string
  description: string
  price: string | number
  period?: string
  savingsText?: string
  buttonText: string
  buttonTo?: string
  features: string[]
  featuresTitle?: string
  isPopular?: boolean
  popularText?: string
  isLoading?: boolean
  loadingText?: string
  isDisabled?: boolean
  variant?: 'primary' | 'secondary' | 'outline'
}

const props = withDefaults(defineProps<Props>(), {
  featuresTitle: "What's included",
  popularText: 'Most Popular',
  loadingText: 'Processing...',
  variant: 'primary'
})

defineEmits<{
  click: []
}>()

const borderClass = computed(() => {
  if (props.isPopular) {
    return 'border-primary-200'
  }
  return 'border-gray-200'
})

const buttonComponent = computed(() => {
  return props.buttonTo ? 'router-link' : 'button'
})

const buttonClass = computed(() => {
  const baseClasses = 'hover:opacity-90 focus:outline-none focus:ring-2 focus:ring-offset-2'
  
  if (props.isDisabled) {
    return `${baseClasses} opacity-50 cursor-not-allowed bg-gray-400 text-white`
  }
  
  switch (props.variant) {
    case 'primary':
      return `${baseClasses} bg-primary-600 text-white hover:bg-primary-700 focus:ring-primary-500`
    case 'secondary':
      return `${baseClasses} bg-gray-600 text-white hover:bg-gray-700 focus:ring-gray-500`
    case 'outline':
      return `${baseClasses} border border-gray-300 bg-white text-gray-700 hover:bg-gray-50 focus:ring-primary-500`
    default:
      return `${baseClasses} bg-primary-600 text-white hover:bg-primary-700 focus:ring-primary-500`
  }
})
</script> 