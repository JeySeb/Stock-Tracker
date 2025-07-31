import { defineStore } from 'pinia'
import { ref, computed, readonly } from 'vue'
import { subscriptionAPI } from '@/api/subscription'
import { useAuthStore } from './auth'
import type { Subscription } from '@/types'

export const useSubscriptionStore = defineStore('subscription', () => {
  // State
  const currentSubscription = ref<Subscription | null>(null)
  const isLoading = ref(false)
  const isProcessingPayment = ref(false)

  // Getters
  const hasActiveSubscription = computed(() => 
    currentSubscription.value?.status === 'active'
  )

  const subscriptionPlan = computed(() => currentSubscription.value?.plan)

  // Actions
  async function fetchCurrentSubscription() {
    const authStore = useAuthStore()
    if (!authStore.isAuthenticated) return

    isLoading.value = true
    try {
      const response = await subscriptionAPI.getCurrentSubscription()
      currentSubscription.value = response
    } catch {
      // No active subscription is not an error
      currentSubscription.value = null
    } finally {
      isLoading.value = false
    }
  }

  async function createSubscription(plan: 'monthly' | 'yearly') {
    isLoading.value = true
    try {
      const response = await subscriptionAPI.createSubscription({ plan })
      currentSubscription.value = response
      return response
    } finally {
      isLoading.value = false
    }
  }

  async function processPayment(subscriptionId: string) {
    isProcessingPayment.value = true
    try {
      await subscriptionAPI.processPayment(subscriptionId)
      
      // Update user tier to premium
      const authStore = useAuthStore()
      if (authStore.user) {
        // Use a method to update user tier instead of direct assignment
        authStore.updateUserTier('premium')
      }
      
      // Refresh subscription data
      await fetchCurrentSubscription()
    } finally {
      isProcessingPayment.value = false
    }
  }

  return {
    // State
    currentSubscription: readonly(currentSubscription),
    isLoading: readonly(isLoading),
    isProcessingPayment: readonly(isProcessingPayment),
    
    // Getters
    hasActiveSubscription,
    subscriptionPlan,
    
    // Actions
    fetchCurrentSubscription,
    createSubscription,
    processPayment
  }
})