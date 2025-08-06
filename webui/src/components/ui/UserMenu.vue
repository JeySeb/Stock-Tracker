<template>
  <div class="relative">
    <!-- User Menu Button -->
    <button
      @click="isOpen = !isOpen"
      class="flex items-center space-x-2 text-sm font-medium text-gray-700 hover:text-gray-900 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 rounded-md px-3 py-2"
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
      class="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg py-1 z-50 border border-gray-200"
    >
      <!-- User Info -->
      <div class="px-4 py-2 border-b border-gray-100">
        <p class="text-sm font-medium text-gray-900">
          {{ authStore.user?.first_name }} {{ authStore.user?.last_name }}
        </p>
        <p class="text-sm text-gray-500">{{ authStore.user?.email }}</p>
      </div>

      <!-- Menu Items -->
      <router-link
        to="/subscription"
        class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 transition-colors"
        @click="isOpen = false"
      >
        Subscription
      </router-link>
      
      <button
        @click="handleLogout"
        class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 transition-colors"
      >
        Sign out
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'

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
if (typeof window !== 'undefined') {
  document.addEventListener('click', handleClickOutside)
}
</script> 