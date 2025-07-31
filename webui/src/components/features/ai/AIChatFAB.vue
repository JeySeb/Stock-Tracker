<template>
    <div class="fixed bottom-6 right-6 z-40">
      <!-- Upgrade Prompt for Non-Premium Users -->
      <div
        v-if="!aiChatStore.canUseAI && showUpgradePrompt"
        class="absolute bottom-16 right-0 w-80 bg-white rounded-lg shadow-xl border border-gray-200 p-4 transform transition-all duration-300"
      >
        <div class="flex items-start space-x-3">
          <div class="w-10 h-10 bg-gradient-to-r from-purple-500 to-pink-500 rounded-full flex items-center justify-center">
            <span class="text-white text-sm font-semibold">AI</span>
          </div>
          <div class="flex-1">
            <h3 class="font-semibold text-gray-900 mb-1">AI Investment Assistant</h3>
            <p class="text-sm text-gray-600 mb-3">
              Get personalized stock recommendations, market insights, and real-time analysis with our AI assistant.
            </p>
            <div class="flex space-x-2">
              <button
                @click="$router.push('/subscription')"
                class="text-xs bg-primary-600 text-white px-3 py-1.5 rounded font-medium hover:bg-primary-700"
              >
                Upgrade to Premium
              </button>
              <button
                @click="showUpgradePrompt = false"
                class="text-xs text-gray-500 hover:text-gray-700"
              >
                Maybe later
              </button>
            </div>
          </div>
          <button
            @click="showUpgradePrompt = false"
            class="text-gray-400 hover:text-gray-600"
          >
            <XMarkIcon class="h-4 w-4" />
          </button>
        </div>
      </div>
  
      <!-- Main FAB -->
      <button
        @click="handleFABClick"
        :class="[
          'w-14 h-14 rounded-full shadow-lg flex items-center justify-center transition-all duration-300 hover:scale-110',
          aiChatStore.canUseAI
            ? 'bg-gradient-to-r from-purple-500 to-pink-500 text-white hover:shadow-xl'
            : 'bg-gray-300 text-gray-600 hover:bg-gray-400'
        ]"
      >
        <ChatBubbleLeftRightIcon v-if="!aiChatStore.isOpen" class="h-6 w-6" />
        <XMarkIcon v-else class="h-6 w-6" />
      </button>
  
      <!-- Notification Badge -->
      <div
        v-if="!aiChatStore.canUseAI && !showUpgradePrompt"
        class="absolute -top-1 -right-1 w-5 h-5 bg-red-500 rounded-full flex items-center justify-center"
      >
        <span class="text-white text-xs font-bold">!</span>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, onMounted } from 'vue'
  import { ChatBubbleLeftRightIcon, XMarkIcon } from '@heroicons/vue/24/outline'
  import { useAIChatStore } from '@/stores/aiChat'
  
  const aiChatStore = useAIChatStore()
  const showUpgradePrompt = ref(false)
  
  const handleFABClick = () => {
    if (aiChatStore.canUseAI) {
      aiChatStore.toggleChat()
    } else {
      showUpgradePrompt.value = !showUpgradePrompt.value
    }
  }
  
  // Show upgrade prompt automatically for new users after a delay
  onMounted(() => {
    if (!aiChatStore.canUseAI) {
      setTimeout(() => {
        const hasSeenPrompt = localStorage.getItem('ai_upgrade_prompt_seen')
        if (!hasSeenPrompt) {
          showUpgradePrompt.value = true
          localStorage.setItem('ai_upgrade_prompt_seen', 'true')
        }
      }, 5000)
    }
  })
  </script>