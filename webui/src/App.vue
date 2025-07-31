<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import { watch } from 'vue'
import { useAuthStore } from '@/stores/auth'
import AIChatInterface from '@/components/features/ai/AIChatInterface.vue'
import AIChatFAB from '@/components/features/ai/AIChatFAB.vue'
import GlobalSearch from '@/components/features/search/GlobalSearch.vue'

const authStore = useAuthStore()

// Debug authentication state changes
watch(() => authStore.isAuthenticated, (newValue, oldValue) => {
  console.log('🔄 Auth state changed:', { from: oldValue, to: newValue })
  console.log('👤 User:', authStore.user)
  console.log('🔑 Has token:', !!authStore.accessToken)
}, { immediate: true })

watch(() => authStore.user, (newUser) => {
  console.log('👤 User changed:', newUser)
}, { immediate: true })
</script>

<template>
  <!-- Main Layout -->
  <div class="min-h-screen bg-gray-50">
    <!-- Navigation Header -->
    <nav v-if="authStore.isAuthenticated" class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <!-- Left side -->
          <div class="flex">
            <div class="flex-shrink-0 flex items-center">
              <h1 class="text-xl font-bold text-primary-600">Stock Tracker</h1>
            </div>
            <div class="hidden sm:ml-6 sm:flex sm:space-x-8">
              <RouterLink
                to="/dashboard"
                class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-900 border-b-2"
                :class="[$route.name === 'dashboard' ? 'border-primary-500' : 'border-transparent hover:border-gray-300']"
              >
                Dashboard
              </RouterLink>
              <RouterLink
                to="/stocks"
                class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-900 border-b-2"
                :class="[$route.name === 'stocks' ? 'border-primary-500' : 'border-transparent hover:border-gray-300']"
              >
                Stocks
              </RouterLink>
              <RouterLink
                to="/recommendations"
                class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-900 border-b-2"
                :class="[$route.name === 'recommendations' ? 'border-primary-500' : 'border-transparent hover:border-gray-300']"
              >
                Recommendations
              </RouterLink>
              <RouterLink
                to="/real-time-data"
                class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-900 border-b-2"
                :class="[$route.name === 'real-time-data' ? 'border-primary-500' : 'border-transparent hover:border-gray-300']"
              >
                Real-Time Data
              </RouterLink>
            </div>
          </div>
          <!-- Right side -->
          <div class="flex items-center space-x-4">
            <!-- Global Search -->
            <div class="w-64">
              <GlobalSearch />
            </div>
            <!-- Subscription Link -->
            <RouterLink
              to="/subscription"
              class="inline-flex items-center px-3 py-2 text-sm font-medium text-gray-500 hover:text-gray-700"
              :class="[$route.name === 'subscription' ? 'text-primary-600' : '']"
            >
              Subscription
            </RouterLink>
            <button
              @click="authStore.logout"
              class="ml-3 inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
            >
              Logout
            </button>
          </div>
        </div>
      </div>
    </nav>

    <!-- Main Content -->
    <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <RouterView />
    </main>

    <!-- AI Chat Components -->
    <AIChatInterface />
    <AIChatFAB />
  </div>
</template>

<style>
/* Global styles */
:root {
  --color-primary-50: #f0f9ff;
  --color-primary-100: #e0f2fe;
  --color-primary-500: #0ea5e9;
  --color-primary-600: #0284c7;
  --color-primary-700: #0369a1;
}

.bg-primary-600 {
  background-color: var(--color-primary-600);
}

.hover\:bg-primary-700:hover {
  background-color: var(--color-primary-700);
}

.border-primary-500 {
  border-color: var(--color-primary-500);
}

.focus\:ring-primary-500:focus {
  --tw-ring-color: var(--color-primary-500);
}
</style>
