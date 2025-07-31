<template>
    <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div class="max-w-md w-full space-y-8">
        <div>
          <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
            Sign in to your account
          </h2>
        </div>
        
        <form class="mt-8 space-y-6" @submit.prevent="handleLogin">
          <div class="space-y-4">
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
              />
              <p v-if="errors.email" class="mt-1 text-sm text-red-600">
                {{ errors.email }}
              </p>
            </div>
            
            <div>
              <label for="password" class="block text-sm font-medium text-gray-700">
                Password
              </label>
              <input
                id="password"
                v-model="form.password"
                type="password"
                required
                class="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500 focus:z-10 sm:text-sm"
                :class="{ 'border-red-300 focus:border-red-500 focus:ring-red-500': errors.password }"
                placeholder="Enter your password"
              />
              <p v-if="errors.password" class="mt-1 text-sm text-red-600">
                {{ errors.password }}
              </p>
            </div>
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
  
          <div>
            <button
              type="submit"
              :disabled="authStore.isLoading"
              class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50"
            >
              <span v-if="!authStore.isLoading">Sign in</span>
              <span v-else>Signing in...</span>
            </button>
          </div>
  
          <div class="text-center">
            <router-link to="/register" class="text-primary-600 hover:text-primary-500">
              Don't have an account? Sign up
            </router-link>
          </div>

          <!-- Demo Login Section -->
          <div class="border-t border-gray-200 pt-6">
            <div class="text-center mb-4">
              <span class="text-sm text-gray-500">For testing purposes:</span>
            </div>
            <div class="space-y-2">
              <button
                @click="handleDemoLogin"
                type="button"
                class="w-full flex justify-center py-2 px-4 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
              >
                🚀 Demo Login (Test UI Styling)
              </button>
              <button
                @click="checkAuthState"
                type="button"
                class="w-full flex justify-center py-2 px-4 border border-blue-300 text-sm font-medium rounded-md text-blue-700 bg-blue-50 hover:bg-blue-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              >
                🔍 Check Auth State
              </button>
              <button
                @click="testLogin"
                type="button"
                class="w-full flex justify-center py-2 px-4 border border-green-300 text-sm font-medium rounded-md text-green-700 bg-green-50 hover:bg-green-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
              >
                🧪 Test Login
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { reactive } from 'vue'
  import { useRouter } from 'vue-router'
  import { useAuthStore } from '@/stores/auth'
  import { useAuthValidation } from '@/composables/useAuthValidation'
  import type { AuthFormData } from '@/composables/useAuthValidation'

  const router = useRouter()
  const authStore = useAuthStore()
  const { errors, validateLoginForm, handleApiError } = useAuthValidation()

  const form = reactive<AuthFormData>({
    email: '',
    password: ''
  })
  
  async function handleLogin() {
    if (!validateLoginForm(form)) {
      return
    }
    
    try {
      console.log('🔐 Attempting login with:', { email: form.email, password: '[HIDDEN]' })
      
      const loginResponse = await authStore.login(form.email, form.password)
      console.log('✅ Login successful, response:', loginResponse)
      console.log('🔍 Auth state after login:', {
        isAuthenticated: authStore.isAuthenticated,
        user: authStore.user,
        userTier: authStore.userTier,
        hasAccessToken: !!authStore.accessToken
      })
      
      // Small delay to ensure state is updated
      await new Promise(resolve => setTimeout(resolve, 100))
      
      console.log('🔄 Redirecting to dashboard after login')
      const navigationResult = await router.push('/dashboard')
      console.log('🎯 Login navigation result:', navigationResult)
      console.log('📍 Current route after login navigation:', router.currentRoute.value.name)
    } catch (error: unknown) {
      console.error('Login failed:', error)
      handleApiError(error)
    }
  }

  function handleDemoLogin() {
    authStore.demoLogin()
    router.push('/dashboard')
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

  function testLogin() {
    // Fill form with test credentials (using the registered user)
    form.email = 'test2@email.com'
    form.password = 'TestPassword123!'
    
    // Trigger login
    handleLogin()
  }
  </script>