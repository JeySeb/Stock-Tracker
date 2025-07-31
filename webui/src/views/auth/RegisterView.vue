<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <!-- Header -->
      <div>
        <div class="mx-auto h-12 w-12 flex items-center justify-center rounded-full bg-primary-100">
          <svg class="h-6 w-6 text-primary-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
          </svg>
        </div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
          Create your account
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          Join Stock Tracker to get AI-powered investment insights
        </p>
      </div>
      
      <!-- Registration Form -->
      <form class="mt-8 space-y-6" @submit.prevent="handleRegister">
        <div class="space-y-4">
          <!-- First Name -->
          <div>
            <label for="first_name" class="block text-sm font-medium text-gray-700">
              First Name
            </label>
            <input
              id="first_name"
              v-model="form.first_name"
              type="text"
              required
              minlength="1"
              maxlength="100"
              class="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500 focus:z-10 sm:text-sm"
              :class="{ 'border-red-300 focus:border-red-500 focus:ring-red-500': errors.first_name }"
              placeholder="Enter your first name"
              @blur="validateRegistrationForm(form)"
            />
            <p v-if="errors.first_name" class="mt-1 text-sm text-red-600">
              {{ errors.first_name }}
            </p>
          </div>

          <!-- Last Name -->
          <div>
            <label for="last_name" class="block text-sm font-medium text-gray-700">
              Last Name
            </label>
            <input
              id="last_name"
              v-model="form.last_name"
              type="text"
              required
              minlength="1"
              maxlength="100"
              class="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500 focus:z-10 sm:text-sm"
              :class="{ 'border-red-300 focus:border-red-500 focus:ring-red-500': errors.last_name }"
              placeholder="Enter your last name"
              @blur="validateRegistrationForm(form)"
            />
            <p v-if="errors.last_name" class="mt-1 text-sm text-red-600">
              {{ errors.last_name }}
            </p>
          </div>
          
          <!-- Email -->
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700">
              Email address
            </label>
            <input
              id="email"
              v-model="form.email"
              type="email"
              required
              class="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500 focus:z-10 sm:text-sm"
              :class="{ 'border-red-300 focus:border-red-500 focus:ring-red-500': errors.email }"
              placeholder="Enter your email"
              @blur="validateRegistrationForm(form)"
            />
            <p v-if="errors.email" class="mt-1 text-sm text-red-600 flex items-center">
              <svg class="h-4 w-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
              </svg>
              {{ errors.email }}
            </p>
          </div>
          
          <!-- Password -->
          <div>
            <label for="password" class="block text-sm font-medium text-gray-700">
              Password
            </label>
            <div class="relative">
              <input
                id="password"
                v-model="form.password"
                :type="showPassword ? 'text' : 'password'"
                required
                minlength="8"
                class="mt-1 appearance-none relative block w-full px-3 py-2 pr-10 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500 focus:z-10 sm:text-sm"
                :class="{ 'border-red-300 focus:border-red-500 focus:ring-red-500': errors.password }"
                placeholder="Create a secure password"
                @blur="validateRegistrationForm(form)"
              />
              <button
                type="button"
                class="absolute inset-y-0 right-0 pr-3 flex items-center"
                @click="showPassword = !showPassword"
              >
                <svg v-if="showPassword" class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21" />
                </svg>
                <svg v-else class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
              </button>
            </div>
            <p v-if="errors.password" class="mt-1 text-sm text-red-600">
              {{ errors.password }}
            </p>
            <div class="mt-2">
              <div class="text-xs text-gray-500">
                Password must be at least 8 characters long
              </div>
              <div class="mt-1 flex items-center justify-between">
                <div class="flex space-x-1">
                  <div 
                    class="h-1 flex-1 rounded-full"
                    :class="passwordStrength >= 1 ? 'bg-red-500' : 'bg-gray-200'"
                  ></div>
                  <div 
                    class="h-1 flex-1 rounded-full"
                    :class="passwordStrength >= 2 ? 'bg-yellow-500' : 'bg-gray-200'"
                  ></div>
                  <div 
                    class="h-1 flex-1 rounded-full"
                    :class="passwordStrength >= 3 ? 'bg-green-500' : 'bg-gray-200'"
                  ></div>
                </div>
                <span class="text-xs ml-2" :class="passwordStrengthColor">
                  {{ passwordStrengthText }}
                </span>
              </div>
            </div>
          </div>

          <!-- Confirm Password -->
          <div>
            <label for="confirm_password" class="block text-sm font-medium text-gray-700">
              Confirm Password
            </label>
            <input
              id="confirm_password"
              v-model="form.confirm_password"
              type="password"
              required
              class="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500 focus:z-10 sm:text-sm"
              :class="{ 'border-red-300 focus:border-red-500 focus:ring-red-500': errors.confirm_password }"
              placeholder="Confirm your password"
              @blur="validateRegistrationForm(form)"
            />
            <p v-if="errors.confirm_password" class="mt-1 text-sm text-red-600">
              {{ errors.confirm_password }}
            </p>
          </div>
        </div>

        <!-- Subscription Plan Selection -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-3">
            Choose Your Plan (Optional)
          </label>
          <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
            <!-- Free Plan -->
            <div class="relative rounded-lg border border-gray-300 bg-white px-4 py-3 cursor-pointer focus-within:ring-2 focus-within:ring-primary-500 focus-within:border-primary-500 hover:bg-gray-50" :class="{ 'border-primary-500 bg-primary-50': selectedPlan === 'free' }">
              <input
                id="plan-free"
                v-model="selectedPlan"
                value="free"
                type="radio"
                name="plan"
                class="sr-only"
              />
              <label for="plan-free" class="cursor-pointer">
                <div class="flex items-center justify-between">
                  <div class="text-sm">
                    <div class="font-medium text-gray-900">Free</div>
                    <div class="text-gray-500">Basic features</div>
                  </div>
                  <div class="text-lg font-bold text-gray-900">$0</div>
                </div>
              </label>
            </div>

            <!-- Monthly Plan -->
            <div class="relative rounded-lg border border-gray-300 bg-white px-4 py-3 cursor-pointer focus-within:ring-2 focus-within:ring-primary-500 focus-within:border-primary-500 hover:bg-gray-50" :class="{ 'border-primary-500 bg-primary-50': selectedPlan === 'monthly' }">
              <input
                id="plan-monthly"
                v-model="selectedPlan"
                value="monthly"
                type="radio"
                name="plan"
                class="sr-only"
              />
              <label for="plan-monthly" class="cursor-pointer">
                <div class="flex items-center justify-between">
                  <div class="text-sm">
                    <div class="font-medium text-gray-900">Monthly</div>
                    <div class="text-gray-500">Premium features</div>
                  </div>
                  <div class="text-lg font-bold text-gray-900">$29.99</div>
                </div>
              </label>
            </div>

            <!-- Yearly Plan -->
            <div class="relative rounded-lg border border-gray-300 bg-white px-4 py-3 cursor-pointer focus-within:ring-2 focus-within:ring-primary-500 focus-within:border-primary-500 hover:bg-gray-50" :class="{ 'border-primary-500 bg-primary-50': selectedPlan === 'yearly' }">
              <input
                id="plan-yearly"
                v-model="selectedPlan"
                value="yearly"
                type="radio"
                name="plan"
                class="sr-only"
              />
              <label for="plan-yearly" class="cursor-pointer">
                <div class="flex items-center justify-between">
                  <div class="text-sm">
                    <div class="font-medium text-gray-900">Yearly</div>
                    <div class="text-gray-500">
                      <div>Premium features</div>
                      <div class="text-green-600 font-medium">Save $59.89!</div>
                    </div>
                  </div>
                  <div class="text-lg font-bold text-gray-900">$299.99</div>
                </div>
              </label>
            </div>
          </div>
          <p class="mt-2 text-xs text-gray-500">
            You can upgrade or change your plan anytime after registration.
          </p>
        </div>

        <!-- Terms and Conditions -->
        <div class="flex items-center">
          <input
            id="terms"
            v-model="form.acceptTerms"
            type="checkbox"
            required
            class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
          />
          <label for="terms" class="ml-2 block text-sm text-gray-900">
            I agree to the
            <a href="#" class="text-primary-600 hover:text-primary-500">Terms of Service</a>
            and
            <a href="#" class="text-primary-600 hover:text-primary-500">Privacy Policy</a>
          </label>
        </div>

        <!-- General Error Display -->
        <div v-if="errors.general" class="rounded-md bg-red-50 p-4">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
              </svg>
            </div>
            <div class="ml-3">
              <p class="text-sm text-red-800">
                {{ errors.general }}
              </p>
            </div>
          </div>
        </div>

        <!-- Submit Button -->
        <div>
          <button
            type="submit"
            :disabled="authStore.isLoading || !isFormValid"
            class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="authStore.isLoading" class="flex items-center">
              <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Creating account...
            </span>
            <span v-else>Create account</span>
          </button>
        </div>

        <!-- Login Link -->
        <div class="text-center">
          <router-link to="/login" class="text-primary-600 hover:text-primary-500">
            Already have an account? Sign in
          </router-link>
        </div>

        <!-- Demo Registration Section -->
        <div class="border-t border-gray-200 pt-6">
          <div class="text-center mb-4">
            <span class="text-sm text-gray-500">For testing purposes:</span>
          </div>
          <div class="space-y-2">
            <button
              @click="handleDemoRegister"
              type="button"
              class="w-full flex justify-center py-2 px-4 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
            >
              🚀 Demo Registration (Test UI Styling)
            </button>
            <button
              @click="testErrorHandling"
              type="button"
              class="w-full flex justify-center py-2 px-4 border border-red-300 text-sm font-medium rounded-md text-red-700 bg-red-50 hover:bg-red-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
            >
              🧪 Test Error Handling
            </button>
            <button
              @click="checkAuthState"
              type="button"
              class="w-full flex justify-center py-2 px-4 border border-blue-300 text-sm font-medium rounded-md text-blue-700 bg-blue-50 hover:bg-blue-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              🔍 Check Auth State
            </button>
          </div>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, computed, watch, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useSubscriptionStore } from '@/stores/subscription'
import { useAuthValidation } from '@/composables/useAuthValidation'
import type { AuthFormData } from '@/composables/useAuthValidation'

const router = useRouter()
const authStore = useAuthStore()
const subscriptionStore = useSubscriptionStore()
const { 
  errors, 
  validateRegistrationForm, 
  handleApiError, 
  calculatePasswordStrength,
  getPasswordStrengthColor,
  getPasswordStrengthText 
} = useAuthValidation()

const form = reactive<AuthFormData>({
  first_name: '',
  last_name: '',
  email: '',
  password: '',
  confirm_password: '',
  acceptTerms: false
})

const showPassword = ref(false)
const selectedPlan = ref('free')

// Password strength calculation
const passwordStrength = computed(() => calculatePasswordStrength(form.password))
const passwordStrengthColor = computed(() => getPasswordStrengthColor(passwordStrength.value))
const passwordStrengthText = computed(() => getPasswordStrengthText(passwordStrength.value))

// Form validation
const isFormValid = computed(() => {
  return (form.first_name?.trim() || '') &&
         (form.last_name?.trim() || '') &&
         form.email.trim() &&
         form.password.length >= 8 &&
         form.password === form.confirm_password &&
         form.acceptTerms &&
         !Object.values(errors).some(error => error)
})

// Watch for password changes to validate confirm password
watch(() => form.password, () => {
  if (form.confirm_password) {
    validateRegistrationForm(form)
  }
})

// Watch for confirm password changes
watch(() => form.confirm_password, () => {
  validateRegistrationForm(form)
})

async function handleRegister() {
  if (!validateRegistrationForm(form)) {
    return
  }
  
  try {
    const userData = {
      email: form.email.trim(),
      password: form.password,
      first_name: form.first_name?.trim() || '',
      last_name: form.last_name?.trim() || ''
    }
    
    console.log('Attempting registration with:', { ...userData, password: '[HIDDEN]' })
    
    const registrationResponse = await authStore.register(userData)
    console.log('✅ Registration successful, response:', registrationResponse)
    console.log('🔍 Auth state after registration:', {
      isAuthenticated: authStore.isAuthenticated,
      user: authStore.user,
      userTier: authStore.userTier,
      hasAccessToken: !!authStore.accessToken
    })
    
    // Small delay to ensure state is updated
    await new Promise(resolve => setTimeout(resolve, 100))
    
    // Handle subscription plan selection
    if (selectedPlan.value === 'free') {
      console.log('🔄 Redirecting to dashboard (free plan)')
      const navigationResult = await router.push('/dashboard')
      console.log('🎯 Navigation result:', navigationResult)
      console.log('📍 Current route after navigation:', router.currentRoute.value.name)
    } else {
      // Create subscription and redirect to subscription page for payment
      try {
        console.log('📦 Creating subscription:', selectedPlan.value)
        await subscriptionStore.createSubscription(selectedPlan.value as 'monthly' | 'yearly')
        console.log('🔄 Redirecting to subscription page')
        const navigationResult = await router.push('/subscription')
        console.log('🎯 Navigation result:', navigationResult)
        console.log('📍 Current route after navigation:', router.currentRoute.value.name)
      } catch (subscriptionError) {
        console.error('Failed to create subscription:', subscriptionError)
        // Still redirect to dashboard if subscription creation fails
        console.log('🔄 Redirecting to dashboard (subscription creation failed)')
        const navigationResult = await router.push('/dashboard')
        console.log('🎯 Navigation result:', navigationResult)
        console.log('📍 Current route after navigation:', router.currentRoute.value.name)
      }
    }
  } catch (error: unknown) {
    console.error('Registration failed:', error)
    handleApiError(error)
    
    // Ensure the error is visible to the user
    // Scroll to the top if there's a general error
    if (errors.general) {
      window.scrollTo({ top: 0, behavior: 'smooth' })
    }
  }
}

function handleDemoRegister() {
  // Fill form with demo data
  form.first_name = 'Demo'
  form.last_name = 'User'
  form.email = 'demo@example.com'
  form.password = 'DemoPassword123!'
  form.confirm_password = 'DemoPassword123!'
  form.acceptTerms = true
  
  // Use demo login for testing
  authStore.demoLogin()
  router.push('/dashboard')
}

function testErrorHandling() {
  // Fill form with an email that likely exists
  form.first_name = 'Test'
  form.last_name = 'User'
  form.email = 'juans.pinzonr@gmail.com' // This email already exists
  form.password = 'TestPassword123!'
  form.confirm_password = 'TestPassword123!'
  form.acceptTerms = true
  selectedPlan.value = 'free'
  
  // Trigger registration to test error handling
  handleRegister()
}

function checkAuthState() {
  console.log('🔍 Current auth state:', {
    isAuthenticated: authStore.isAuthenticated,
    user: authStore.user,
    userTier: authStore.userTier,
    hasAccessToken: !!authStore.accessToken,
    localStorage: {
      access: !!localStorage.getItem('stock_tracker_access_token'),
      refresh: !!localStorage.getItem('stock_tracker_refresh_token')
    }
  })
}
</script>