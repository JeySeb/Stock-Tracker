<template>
  <div class="relative">
    <!-- User Menu Button -->
    <button
      @click="isOpen = !isOpen"
      class="flex items-center space-x-2 text-sm font-medium text-gray-700 hover:text-gray-900 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 rounded-md px-3 py-2 transition-colors"
    >
      <div class="w-8 h-8 bg-primary-600 rounded-full flex items-center justify-center">
        <span class="text-white text-sm font-medium">
          {{ userInitials }}
        </span>
      </div>
      <span class="hidden md:block">{{ authStore.user?.first_name || 'User' }}</span>
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
      </svg>
    </button>

    <!-- Dropdown Menu -->
    <div
      v-if="isOpen"
      class="absolute right-0 mt-2 w-64 bg-white rounded-md shadow-lg py-1 z-50 border border-gray-200"
    >
      <!-- User Info Section -->
      <div class="px-4 py-3 border-b border-gray-100">
        <div class="flex items-center space-x-3">
          <div class="w-10 h-10 bg-primary-600 rounded-full flex items-center justify-center">
            <span class="text-white text-sm font-medium">
              {{ userInitials }}
            </span>
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-gray-900 truncate">
              {{ authStore.user?.first_name }} {{ authStore.user?.last_name }}
            </p>
            <p class="text-sm text-gray-500 truncate">{{ authStore.user?.email }}</p>
            <div class="flex items-center mt-1">
              <TierBadge :tier="authStore.userTier" />
            </div>
          </div>
        </div>
      </div>

      <!-- Menu Items -->
      <div class="py-1">
        <!-- Dashboard -->
        <router-link
          to="/dashboard"
          class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 transition-colors"
          @click="isOpen = false"
        >
          <svg class="w-4 h-4 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5a2 2 0 012-2h4a2 2 0 012 2v6H8V5z" />
          </svg>
          Dashboard
        </router-link>

        <!-- Subscription Management -->
        <router-link
          to="/subscription"
          class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 transition-colors"
          @click="isOpen = false"
        >
          <svg class="w-4 h-4 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1" />
          </svg>
          Subscription
        </router-link>

        <!-- Profile (disabled until implemented) -->
        <!--
        <router-link
          to="/profile"
          class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 transition-colors"
          @click="isOpen = false"
        >
          <svg class="w-4 h-4 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
          </svg>
          Profile
        </router-link>
        -->

        <!-- Divider -->
        <div class="border-t border-gray-100 my-1"></div>

        <!-- Logout -->
        <button
          @click="handleLogout"
          class="flex items-center w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 transition-colors"
        >
          <svg class="w-4 h-4 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
          </svg>
          Sign out
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import TierBadge from './TierBadge.vue'

const authStore = useAuthStore()
const isOpen = ref(false)

const userInitials = computed(() => {
  const user = authStore.user
  if (!user) return 'U'
  
  const first = user.first_name?.charAt(0) || ''
  const last = user.last_name?.charAt(0) || ''
  return (first + last).toUpperCase() || 'U'
})

function handleLogout() {
  isOpen.value = false
  authStore.logout()
}

// Close dropdown when clicking outside
function handleClickOutside(event: Event) {
  const target = event.target as Element
  if (!target.closest('.relative')) {
    isOpen.value = false
  }
}

// Add click outside listener
onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script> 