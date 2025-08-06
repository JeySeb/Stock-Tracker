<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import { watch } from 'vue'
import { useAuthStore } from '@/stores/auth'
import AIChatInterface from '@/components/features/ai/AIChatInterface.vue'
import AIChatFAB from '@/components/features/ai/AIChatFAB.vue'
import AppHeader from '@/components/layout/AppHeader.vue'

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
    <!-- Unified Navigation Header -->
    <AppHeader />

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
