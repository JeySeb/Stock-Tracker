<template>
  <header class="bg-white shadow-sm border-b border-gray-200">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center h-16">
        <!-- Logo and Navigation -->
        <div class="flex items-center">
          <!-- App Name as Home Link -->
          <router-link 
            to="/"  
            class="flex-shrink-0 flex items-center hover:opacity-80 transition-opacity"
          >
            <h1 class="text-xl font-bold text-primary-600">Stock Tracker</h1>
          </router-link>
          
          <!-- Navigation Links (only for authenticated users) -->
          <div v-if="authStore.isAuthenticated" class="hidden sm:ml-6 sm:flex sm:space-x-8">
            <router-link
              to="/dashboard"
              class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-900 border-b-2 transition-colors"
              :class="[$route.name === 'dashboard' ? 'border-primary-500' : 'border-transparent hover:border-gray-300']"
            >
              Dashboard
            </router-link>
            <router-link
              to="/stocks"
              class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-900 border-b-2 transition-colors"
              :class="[$route.name === 'stocks' ? 'border-primary-500' : 'border-transparent hover:border-gray-300']"
            >
              Stocks
            </router-link>
            <router-link
              to="/recommendations"
              class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-900 border-b-2 transition-colors"
              :class="[$route.name === 'recommendations' ? 'border-primary-500' : 'border-transparent hover:border-gray-300']"
            >
              Recommendations
            </router-link>
            <router-link
              to="/real-time-data"
              class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-900 border-b-2 transition-colors"
              :class="[$route.name === 'real-time-data' ? 'border-primary-500' : 'border-transparent hover:border-gray-300']"
            >
              Real-Time Data
            </router-link>
          </div>
        </div>
        
        <!-- Right Side Controls -->
        <div class="flex items-center space-x-4">
          <!-- Global Search (only for authenticated users) -->
          <div v-if="authStore.isAuthenticated" class="w-64">
            <GlobalSearch />
          </div>
          
          <!-- Subscription Link (only for authenticated users) -->
          <router-link
            v-if="authStore.isAuthenticated"
            to="/subscription"
            class="inline-flex items-center px-3 py-2 text-sm font-medium text-gray-500 hover:text-gray-700 transition-colors"
            :class="[$route.name === 'subscription' ? 'text-primary-600' : '']"
          >
            Subscription
          </router-link>
          
          <!-- User Controls -->
          <div class="flex items-center space-x-4">
            <!-- Tier Badge (only for authenticated users) -->
            <TierBadge v-if="authStore.isAuthenticated" :tier="authStore.userTier" />
            
            <!-- User Menu (only for authenticated users) -->
            <UserMenu v-if="authStore.isAuthenticated" />
            
            <!-- Auth Buttons (for guest users) -->
            <div v-else class="space-x-2">
              <router-link
                to="/login"
                class="text-gray-700 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium transition-colors"
              >
                Sign In
              </router-link>
              <router-link
                to="/register"
                class="bg-primary-600 text-white hover:bg-primary-700 px-4 py-2 rounded-md text-sm font-medium transition-colors"
              >
                Get Started
              </router-link>
            </div>
            
            <!-- Logout Button (only for authenticated users) -->
            <button
              v-if="authStore.isAuthenticated"
              @click="authStore.logout"
              class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-colors"
            >
              Logout
            </button>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { useRoute } from 'vue-router'
import TierBadge from '@/components/ui/TierBadge.vue'
import UserMenu from '@/components/ui/UserMenu.vue'
import GlobalSearch from '@/components/features/search/GlobalSearch.vue'

const authStore = useAuthStore()
const route = useRoute()
</script>