<template>
  <div class="min-h-screen bg-gray-50 py-12">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- Header -->
      <div class="text-center">
        <h2 class="text-3xl font-extrabold text-gray-900 sm:text-4xl">
          Choose Your Plan
        </h2>
        <p class="mt-3 max-w-2xl mx-auto text-xl text-gray-500 sm:mt-4">
          Unlock advanced features with our premium subscription plans
        </p>
      </div>

      <!-- Current Subscription Status -->
      <div v-if="authStore.isAuthenticated && subscriptionStore.currentSubscription" class="mt-8">
        <div class="bg-white rounded-lg shadow p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Current Subscription</h3>
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-600">Plan: {{ subscriptionStore.currentSubscription.plan }}</p>
              <p class="text-sm text-gray-600">Status: 
                <span :class="getStatusColor(subscriptionStore.currentSubscription.status)">
                  {{ subscriptionStore.currentSubscription.status }}
                </span>
              </p>
              <p class="text-sm text-gray-600">End Date: {{ formatDate(subscriptionStore.currentSubscription.end_date) }}</p>
            </div>
            <div class="text-right">
              <p class="text-lg font-semibold">${{ subscriptionStore.currentSubscription.price }}</p>
              <p class="text-sm text-gray-500">{{ subscriptionStore.currentSubscription.currency }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Pricing Plans -->
      <div class="mt-12 space-y-4 sm:mt-16 sm:space-y-0 sm:grid sm:grid-cols-2 sm:gap-6 lg:max-w-4xl lg:mx-auto xl:max-w-none xl:mx-0 xl:grid-cols-2">
        <!-- Monthly Plan -->
        <div class="border border-gray-200 rounded-lg shadow-sm divide-y divide-gray-200">
          <div class="p-6">
            <h3 class="text-lg leading-6 font-medium text-gray-900">Monthly</h3>
            <p class="mt-4 text-sm text-gray-500">
              Perfect for trying out premium features
            </p>
            <p class="mt-8">
              <span class="text-4xl font-extrabold text-gray-900">$29.99</span>
              <span class="text-base font-medium text-gray-500">/month</span>
            </p>
            <button
              @click="handleSubscription('monthly')"
              :disabled="subscriptionStore.isLoading"
              class="mt-8 block w-full bg-primary-600 border border-transparent rounded-md py-2 text-sm font-semibold text-white text-center hover:bg-primary-700 disabled:opacity-50"
            >
              <span v-if="subscriptionStore.isLoading">Processing...</span>
              <span v-else>{{ getButtonText('monthly') }}</span>
            </button>
          </div>
          <div class="pt-6 pb-8 px-6">
            <h4 class="text-xs font-medium text-gray-900 tracking-wide uppercase">What's included</h4>
            <ul class="mt-6 space-y-4">
              <li v-for="feature in monthlyFeatures" :key="feature" class="flex space-x-3">
                <svg class="flex-shrink-0 h-5 w-5 text-green-500" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
                <span class="text-sm text-gray-500">{{ feature }}</span>
              </li>
            </ul>
          </div>
        </div>

        <!-- Yearly Plan -->
        <div class="border border-gray-200 rounded-lg shadow-sm divide-y divide-gray-200 relative">
          <!-- Popular badge -->
          <div class="absolute top-0 right-6 transform -translate-y-1/2">
            <span class="inline-flex rounded-full bg-primary-100 px-4 py-1 text-sm font-semibold text-primary-600">
              Most Popular
            </span>
          </div>
          <div class="p-6">
            <h3 class="text-lg leading-6 font-medium text-gray-900">Yearly</h3>
            <p class="mt-4 text-sm text-gray-500">
              Best value for committed users
            </p>
            <p class="mt-8">
              <span class="text-4xl font-extrabold text-gray-900">$299.99</span>
              <span class="text-base font-medium text-gray-500">/year</span>
            </p>
            <p class="text-sm text-green-600 mt-2">Save $59.89 compared to monthly!</p>
            <button
              @click="handleSubscription('yearly')"
              :disabled="subscriptionStore.isLoading"
              class="mt-8 block w-full bg-primary-600 border border-transparent rounded-md py-2 text-sm font-semibold text-white text-center hover:bg-primary-700 disabled:opacity-50"
            >
              <span v-if="subscriptionStore.isLoading">Processing...</span>
              <span v-else>{{ getButtonText('yearly') }}</span>
            </button>
          </div>
          <div class="pt-6 pb-8 px-6">
            <h4 class="text-xs font-medium text-gray-900 tracking-wide uppercase">What's included</h4>
            <ul class="mt-6 space-y-4">
              <li v-for="feature in yearlyFeatures" :key="feature" class="flex space-x-3">
                <svg class="flex-shrink-0 h-5 w-5 text-green-500" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
                <span class="text-sm text-gray-500">{{ feature }}</span>
              </li>
            </ul>
          </div>
        </div>
      </div>

      <!-- Tier Comparison -->
      <div class="mt-16">
        <h3 class="text-2xl font-bold text-gray-900 text-center mb-8">Feature Comparison</h3>
        <div class="overflow-hidden shadow ring-1 ring-black ring-opacity-5 md:rounded-lg">
          <table class="min-w-full divide-y divide-gray-300">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Feature</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Guest</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Registered</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Premium</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="comparison in featureComparison" :key="comparison.feature">
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{{ comparison.feature }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  <span v-if="comparison.guest" class="text-green-600">✓</span>
                  <span v-else class="text-red-600">✗</span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  <span v-if="comparison.registered" class="text-green-600">✓</span>
                  <span v-else class="text-red-600">✗</span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  <span v-if="comparison.premium" class="text-green-600">✓</span>
                  <span v-else class="text-red-600">✗</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Payment Processing Modal -->
      <div v-if="showPaymentModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" @click="showPaymentModal = false">
        <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white" @click.stop>
          <div class="mt-3 text-center">
            <div class="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-green-100">
              <svg v-if="!subscriptionStore.isProcessingPayment" class="h-6 w-6 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              <svg v-else class="animate-spin h-6 w-6 text-green-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
            </div>
            <h3 class="text-lg leading-6 font-medium text-gray-900 mt-2">
              {{ subscriptionStore.isProcessingPayment ? 'Processing Payment' : 'Payment Simulation' }}
            </h3>
            <div class="mt-2 px-7 py-3">
              <p class="text-sm text-gray-500">
                {{ subscriptionStore.isProcessingPayment ? 'Please wait while we process your payment...' : 'Click the button below to simulate payment processing.' }}
              </p>
            </div>
            <div class="items-center px-4 py-3">
              <button
                v-if="!subscriptionStore.isProcessingPayment"
                @click="processPayment"
                class="px-4 py-2 bg-primary-600 text-white text-base font-medium rounded-md w-full shadow-sm hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500"
              >
                Simulate Payment
              </button>
              <div v-else class="text-center">
                <p class="text-sm text-gray-600">Processing payment...</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useSubscriptionStore } from '@/stores/subscription'

const router = useRouter()
const authStore = useAuthStore()
const subscriptionStore = useSubscriptionStore()

const showPaymentModal = ref(false)
const pendingSubscription = ref<{ id: string; plan: string } | null>(null)

const monthlyFeatures = [
  'Real-time market data',
  'External API integration',
  'Advanced analytics',
  'Priority support',
  'All registered features'
]

const yearlyFeatures = [
  'Everything in Monthly',
  'AI-powered insights',
  'Sentiment analysis',
  'Advanced AI chatbot',
  'Custom alerts',
  'Priority customer support'
]

const featureComparison = [
  { feature: 'Basic recommendations', guest: true, registered: true, premium: true },
  { feature: 'Market analytics', guest: true, registered: true, premium: true },
  { feature: 'Real-time data', guest: false, registered: true, premium: true },
  { feature: 'External APIs', guest: false, registered: true, premium: true },
  { feature: 'AI insights', guest: false, registered: false, premium: true },
  { feature: 'Sentiment analysis', guest: false, registered: false, premium: true },
  { feature: 'AI chatbot', guest: false, registered: false, premium: true }
]

onMounted(() => {
  if (authStore.isAuthenticated) {
    subscriptionStore.fetchCurrentSubscription()
  }
})

const getButtonText = (plan: string) => {
  if (!authStore.isAuthenticated) {
    return 'Sign up to Subscribe'
  }
  
  const currentPlan = subscriptionStore.currentSubscription?.plan
  const currentStatus = subscriptionStore.currentSubscription?.status
  
  if (currentPlan === plan && currentStatus === 'active') {
    return 'Current Plan'
  }
  
  if (currentPlan === plan && currentStatus === 'pending') {
    return 'Complete Payment'
  }
  
  return `Subscribe ${plan === 'monthly' ? 'Monthly' : 'Yearly'}`
}

const getStatusColor = (status: string) => {
  switch (status) {
    case 'active': return 'text-green-600'
    case 'pending': return 'text-yellow-600'
    case 'cancelled': return 'text-red-600'
    case 'expired': return 'text-gray-600'
    default: return 'text-gray-600'
  }
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString()
}

async function handleSubscription(plan: 'monthly' | 'yearly') {
  if (!authStore.isAuthenticated) {
    router.push('/register')
    return
  }

  const currentPlan = subscriptionStore.currentSubscription?.plan
  const currentStatus = subscriptionStore.currentSubscription?.status

  // If user has the same plan and it's active, do nothing
  if (currentPlan === plan && currentStatus === 'active') {
    return
  }

  // If user has pending subscription with same plan, show payment modal
  if (currentPlan === plan && currentStatus === 'pending') {
    pendingSubscription.value = {
      id: subscriptionStore.currentSubscription!.id,
      plan
    }
    showPaymentModal.value = true
    return
  }

  try {
    const subscription = await subscriptionStore.createSubscription(plan)
    pendingSubscription.value = {
      id: subscription.id,
      plan
    }
    showPaymentModal.value = true
  } catch (error) {
    console.error('Failed to create subscription:', error)
    // Handle error (show notification, etc.)
  }
}

async function processPayment() {
  if (!pendingSubscription.value) return

  try {
    await subscriptionStore.processPayment(pendingSubscription.value.id)
    showPaymentModal.value = false
    pendingSubscription.value = null
    
    // Show success message and redirect to dashboard
    setTimeout(() => {
      router.push('/dashboard')
    }, 1000)
  } catch (error) {
    console.error('Payment processing failed:', error)
    // Handle error
  }
}
</script> 