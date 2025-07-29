<template>
    <div class="min-h-screen bg-gray-50">
      <!-- Header -->
      <AppHeader />
      
      <!-- Main Content -->
      <div class="flex">
        <!-- Sidebar -->
        <AppSidebar v-if="showSidebar" />
        
        <!-- Content Area -->
        <main class="flex-1 p-6">
          <router-view />
        </main>
      </div>
      
      <!-- AI Chat Placeholder (Future Phase) -->
      <AIChatPlaceholder />
    </div>
  </template>
  
  <script setup lang="ts">
  import { computed } from 'vue'
  import { useRoute } from 'vue-router'
  import AppHeader from './AppHeader.vue'
  import AppSidebar from './AppSidebar.vue'
  import AIChatPlaceholder from '../placeholder/AIChatPlaceholder.vue'
  
  const route = useRoute()
  
  const showSidebar = computed(() => {
    const hiddenSidebarRoutes = ['login', 'register']
    return !hiddenSidebarRoutes.includes(route.name as string)
  })
  </script>