<template>
    <div class="fixed inset-0 z-50 overflow-hidden" v-if="aiChatStore.isOpen">
      <div class="absolute inset-0 bg-black bg-opacity-50" @click="aiChatStore.closeChat"></div>
      
      <div class="absolute right-4 top-4 bottom-4 w-96 bg-white rounded-lg shadow-2xl flex flex-col">
        <!-- Header -->
        <div class="flex items-center justify-between p-4 border-b border-gray-200">
          <div class="flex items-center space-x-2">
            <div class="w-8 h-8 bg-gradient-to-r from-purple-500 to-pink-500 rounded-full flex items-center justify-center">
              <span class="text-white text-sm font-semibold">AI</span>
            </div>
            <div>
              <h3 class="font-semibold text-gray-900">Investment Assistant</h3>
              <p class="text-xs text-gray-500">
                {{ aiChatStore.isTyping ? 'Typing...' : 'Online' }}
              </p>
            </div>
          </div>
          
          <div class="flex items-center space-x-2">
            <button
              @click="showSessionMenu = !showSessionMenu"
              class="text-gray-400 hover:text-gray-600"
            >
              <Bars3Icon class="h-5 w-5" />
            </button>
            <button
              @click="aiChatStore.closeChat"
              class="text-gray-400 hover:text-gray-600"
            >
              <XMarkIcon class="h-5 w-5" />
            </button>
          </div>
        </div>
  
        <!-- Session Menu -->
        <div v-if="showSessionMenu" class="border-b border-gray-200 bg-gray-50 p-2">
          <div class="flex items-center justify-between mb-2">
            <span class="text-sm font-medium text-gray-700">Chat Sessions</span>
            <button
              @click="aiChatStore.createNewSession(); showSessionMenu = false"
              class="text-xs bg-primary-600 text-white px-2 py-1 rounded"
            >
              New Chat
            </button>
          </div>
          <div class="max-h-32 overflow-y-auto space-y-1">
            <div
              v-for="session in aiChatStore.sessions"
              :key="session.id"
              class="flex items-center justify-between p-2 text-sm rounded hover:bg-gray-100 cursor-pointer"
              :class="{ 'bg-primary-50 border border-primary-200': session.id === aiChatStore.currentSession?.id }"
              @click="aiChatStore.switchToSession(session.id); showSessionMenu = false"
            >
              <span class="truncate">{{ session.title }}</span>
              <button
                @click.stop="aiChatStore.deleteSession(session.id)"
                class="text-gray-400 hover:text-red-500"
              >
                <TrashIcon class="h-4 w-4" />
              </button>
            </div>
          </div>
        </div>
  
        <!-- Messages -->
        <div ref="messagesContainer" class="flex-1 overflow-y-auto p-4 space-y-4">
          <ChatMessage
            v-for="message in aiChatStore.currentMessages"
            :key="message.id"
            :message="message"
          />
  
          <!-- Typing Indicator -->
          <div v-if="aiChatStore.isTyping" class="flex justify-start">
            <div class="bg-gray-100 rounded-lg px-4 py-2">
              <div class="flex space-x-1">
                <div class="w-2 h-2 bg-gray-400 rounded-full animate-bounce"></div>
                <div class="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 0.1s"></div>
                <div class="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 0.2s"></div>
              </div>
            </div>
          </div>
        </div>
  
        <!-- Input -->
        <div class="p-4 border-t border-gray-200">
          <div class="flex space-x-2">
            <input
              v-model="newMessage"
              type="text"
              placeholder="Ask me about stocks, market trends, or get recommendations..."
              class="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-1 focus:ring-primary-500 focus:border-primary-500"
              @keydown.enter="handleSendMessage"
              :disabled="aiChatStore.isTyping"
            />
            <button
              @click="handleSendMessage"
              :disabled="!newMessage.trim() || aiChatStore.isTyping"
              class="px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <PaperAirplaneIcon class="h-4 w-4" />
            </button>
          </div>
          
          <!-- Quick Actions -->
          <div class="mt-2 flex flex-wrap gap-1">
            <button
              v-for="quickAction in quickActions"
              :key="quickAction.text"
              @click="newMessage = quickAction.text; handleSendMessage()"
              class="text-xs bg-gray-100 text-gray-700 px-2 py-1 rounded hover:bg-gray-200"
            >
              {{ quickAction.text }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, nextTick, watch } from 'vue'
  import { XMarkIcon, Bars3Icon, TrashIcon, PaperAirplaneIcon } from '@heroicons/vue/24/outline'
  import { useAIChatStore } from '@/stores/aiChat'
  import ChatMessage from './ChatMessage.vue'
  
  const aiChatStore = useAIChatStore()
  const newMessage = ref('')
  const messagesContainer = ref<HTMLElement>()
  const showSessionMenu = ref(false)
  
  const quickActions = [
    { text: 'Top 8 recent stocks events' },
    { text: 'Top 5 stock recommendations' },
    { text: 'Which are the top 5 most relevant brokers?' },
    { text: 'What do you recommend according to financial news?' }
  ]
  
  const handleSendMessage = async () => {
    if (!newMessage.value.trim()) return
  
    const message = newMessage.value.trim()
    newMessage.value = ''
  
    await aiChatStore.sendMessage(message)
    scrollToBottom()
  }
  

  
  const scrollToBottom = () => {
    nextTick(() => {
      if (messagesContainer.value) {
        messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
      }
    })
  }
  
  // Auto-scroll when new messages arrive
  watch(() => aiChatStore.currentMessages.length, () => {
    scrollToBottom()
  })
  </script>