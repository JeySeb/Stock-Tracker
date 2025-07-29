<template>
    <header class="bg-white shadow-sm border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <!-- Logo -->
          <div class="flex items-center">
            <h1 class="text-2xl font-bold text-gray-900">
              Stock Tracker
            </h1>
          </div>
          
          <!-- User Controls -->
          <div class="flex items-center space-x-4">
            <!-- Tier Badge -->
            <TierBadge v-if="authStore.isAuthenticated" :tier="authStore.userTier" />
            
            <!-- User Menu -->
            <UserMenu v-if="authStore.isAuthenticated" />
            
            <!-- Auth Buttons -->
            <div v-else class="space-x-2">
              <router-link
                to="/login"
                class="text-gray-700 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium"
              >
                Login
              </router-link>
              <router-link
                to="/register"
                class="bg-primary-600 text-white hover:bg-primary-700 px-3 py-2 rounded-md text-sm font-medium"
              >
                Sign Up
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </header>
  </template>
  
  <script setup lang="ts">
  import { useAuthStore } from '@/stores/auth'
  import TierBadge from '@/components/ui/TierBadge.vue'
  import UserMenu from '@/components/ui/UserMenu.vue'
  
  const authStore = useAuthStore()
  </script>