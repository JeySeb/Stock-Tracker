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
          
          <!-- Desktop Navigation Links -->
          <div class="hidden sm:ml-6 sm:flex sm:space-x-6">
            <!-- Always visible links for exploration -->
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
            
            <!-- Authenticated-only links -->
            <router-link
              v-if="authStore.isAuthenticated"
              to="/dashboard"
              class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-900 border-b-2 transition-colors"
              :class="[$route.name === 'dashboard' ? 'border-primary-500' : 'border-transparent hover:border-gray-300']"
            >
              Dashboard
            </router-link>
            <router-link
              v-if="authStore.isAuthenticated && authStore.userTier === 'premium'"
              to="/real-time-data"
              class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-900 border-b-2 transition-colors"
              :class="[$route.name === 'real-time-data' ? 'border-primary-500' : 'border-transparent hover:border-gray-300']"
            >
              Real-Time Data
            </router-link>
          </div>
        </div>
        
        <!-- Right Side Controls -->
        <div class="flex items-center space-x-3">
          <!-- Global Search (only for authenticated users) -->
          <div v-if="authStore.isAuthenticated" class="hidden md:block w-56">
            <GlobalSearch />
          </div>
          
          <!-- User Controls -->
          <div class="flex items-center space-x-3">
            <!-- Tier Badge (only for authenticated users) -->
            <TierBadge v-if="authStore.isAuthenticated" :tier="authStore.userTier" />
            
            <!-- User Menu (only for authenticated users) -->
            <UserMenu v-if="authStore.isAuthenticated" />
            
            <!-- Auth Buttons (for guest users) -->
            <div v-else class="hidden sm:flex items-center space-x-3">
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
          </div>
          
          <!-- Mobile menu button -->
          <div class="sm:hidden">
            <button
              @click="mobileMenuOpen = !mobileMenuOpen"
              class="inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-primary-500"
            >
              <span class="sr-only">Open main menu</span>
              <!-- Icon when menu is closed -->
              <svg
                v-if="!mobileMenuOpen"
                class="block h-6 w-6"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
              </svg>
              <!-- Icon when menu is open -->
              <svg
                v-else
                class="block h-6 w-6"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Mobile menu -->
    <div v-if="mobileMenuOpen" class="sm:hidden">
      <div class="pt-2 pb-3 space-y-1 bg-white border-t border-gray-200">
        <!-- Navigation Links -->
        <router-link
          to="/stocks"
          class="block px-3 py-2 text-base font-medium text-gray-700 hover:text-gray-900 hover:bg-gray-50 transition-colors"
          :class="[$route.name === 'stocks' ? 'bg-primary-50 text-primary-700' : '']"
          @click="mobileMenuOpen = false"
        >
          Stocks
        </router-link>
        <router-link
          to="/recommendations"
          class="block px-3 py-2 text-base font-medium text-gray-700 hover:text-gray-900 hover:bg-gray-50 transition-colors"
          :class="[$route.name === 'recommendations' ? 'bg-primary-50 text-primary-700' : '']"
          @click="mobileMenuOpen = false"
        >
          Recommendations
        </router-link>
        
        <!-- Authenticated-only links -->
        <router-link
          v-if="authStore.isAuthenticated"
          to="/dashboard"
          class="block px-3 py-2 text-base font-medium text-gray-700 hover:text-gray-900 hover:bg-gray-50 transition-colors"
          :class="[$route.name === 'dashboard' ? 'bg-primary-50 text-primary-700' : '']"
          @click="mobileMenuOpen = false"
        >
          Dashboard
        </router-link>
        <router-link
          v-if="authStore.isAuthenticated && authStore.userTier === 'premium'"
          to="/real-time-data"
          class="block px-3 py-2 text-base font-medium text-gray-700 hover:text-gray-900 hover:bg-gray-50 transition-colors"
          :class="[$route.name === 'real-time-data' ? 'bg-primary-50 text-primary-700' : '']"
          @click="mobileMenuOpen = false"
        >
          Real-Time Data
        </router-link>
        
        <!-- Divider for authenticated users -->
        <div v-if="authStore.isAuthenticated" class="border-t border-gray-200 my-2"></div>
        
        <!-- User-specific mobile menu items -->
        <div v-if="authStore.isAuthenticated" class="px-3 py-2">
          <!-- User info -->
          <div class="flex items-center space-x-3 mb-3">
            <div class="w-8 h-8 bg-primary-600 rounded-full flex items-center justify-center">
              <span class="text-white text-sm font-medium">
                {{ userInitials }}
              </span>
            </div>
            <div>
              <p class="text-sm font-medium text-gray-900">
                {{ authStore.user?.first_name }} {{ authStore.user?.last_name }}
              </p>
              <p class="text-xs text-gray-500">{{ authStore.user?.email }}</p>
            </div>
          </div>
          
          <!-- User menu items -->
          <router-link
            to="/subscription"
            class="block px-3 py-2 text-base font-medium text-gray-700 hover:text-gray-900 hover:bg-gray-50 transition-colors"
            @click="mobileMenuOpen = false"
          >
            Subscription
          </router-link>
          
          <button
            @click="handleMobileLogout"
            class="block w-full text-left px-3 py-2 text-base font-medium text-gray-700 hover:text-gray-900 hover:bg-gray-50 transition-colors"
          >
            Sign out
          </button>
        </div>
        
        <!-- Guest user mobile menu -->
        <div v-else class="px-3 py-2 space-y-2">
          <router-link
            to="/login"
            class="block w-full text-center px-3 py-2 text-base font-medium text-gray-700 hover:text-gray-900 hover:bg-gray-50 transition-colors"
            @click="mobileMenuOpen = false"
          >
            Sign In
          </router-link>
          <router-link
            to="/register"
            class="block w-full text-center px-3 py-2 text-base font-medium text-white bg-primary-600 hover:bg-primary-700 transition-colors rounded-md"
            @click="mobileMenuOpen = false"
          >
            Get Started
          </router-link>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRoute } from 'vue-router'
import TierBadge from '@/components/ui/TierBadge.vue'
import UserMenu from '@/components/ui/UserMenu.vue'
import GlobalSearch from '@/components/features/search/GlobalSearch.vue'

const authStore = useAuthStore()
const route = useRoute()
const mobileMenuOpen = ref(false)

const userInitials = computed(() => {
  const user = authStore.user
  if (!user) return 'U'
  
  const first = user.first_name?.charAt(0) || ''
  const last = user.last_name?.charAt(0) || ''
  return (first + last).toUpperCase() || 'U'
})

function handleMobileLogout() {
  mobileMenuOpen.value = false
  authStore.logout()
}
</script>