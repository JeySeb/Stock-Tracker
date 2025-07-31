<template>
    <div v-if="hasError" class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div class="max-w-md w-full space-y-8">
        <div class="text-center">
          <div class="mx-auto h-24 w-24 text-red-500">
            <ExclamationTriangleIcon class="h-full w-full" />
          </div>
          <h2 class="mt-6 text-3xl font-extrabold text-gray-900">
            Something went wrong
          </h2>
          <p class="mt-2 text-sm text-gray-600">
            {{ userMessage }}
          </p>
        </div>
        
        <div class="space-y-4">
          <button
            @click="handleRetry"
            class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
          >
            Try Again
          </button>
          
          <button
            @click="handleGoHome"
            class="w-full flex justify-center py-2 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
          >
            Go to Dashboard
          </button>
          
          <details class="mt-4">
            <summary class="cursor-pointer text-sm text-gray-500 hover:text-gray-700">
              Technical Details
            </summary>
            <div class="mt-2 p-3 bg-gray-100 rounded text-xs font-mono text-gray-700">
              <div><strong>Error ID:</strong> {{ error?.id }}</div>
              <div><strong>Type:</strong> {{ error?.type }}</div>
              <div><strong>Time:</strong> {{ error?.timestamp }}</div>
              <div class="mt-2"><strong>Message:</strong> {{ error?.message }}</div>
            </div>
          </details>
        </div>
      </div>
    </div>
    
    <slot v-else />
  </template>
  
  <script setup lang="ts">
  import { ref, onErrorCaptured } from 'vue'
  import { useRouter } from 'vue-router'
  import { ExclamationTriangleIcon } from '@heroicons/vue/24/outline'
  import { errorHandler, type AppError } from '@/utils/errorHandler'
  
  const router = useRouter()
  const hasError = ref(false)
  const error = ref<AppError | null>(null)
  const userMessage = ref('')
  
  onErrorCaptured((err: Error) => {
    const appError = errorHandler.logError(err)
    
    hasError.value = true
    error.value = appError
    userMessage.value = errorHandler.getUserMessage(appError)
    
    return false // Prevent error from bubbling up
  })
  
  const handleRetry = () => {
    hasError.value = false
    error.value = null
    userMessage.value = ''
    
    // Force component re-render
    window.location.reload()
  }
  
  const handleGoHome = () => {
    hasError.value = false
    error.value = null
    userMessage.value = ''
    
    router.push('/dashboard')
  }
  </script>