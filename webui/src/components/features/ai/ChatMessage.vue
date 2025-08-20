<template>
  <div
    :class="[
      'flex',
      message.type === 'user' ? 'justify-end' : 'justify-start'
    ]"
  >
    <div
      :class="[
        'max-w-[80%] rounded-lg px-4 py-2',
        message.type === 'user'
          ? 'bg-primary-600 text-white'
          : message.type === 'system'
          ? 'bg-gray-100 text-gray-700 text-sm'
          : 'bg-gray-100 text-gray-900'
      ]"
    >
      <!-- Message Content -->
      <ChatTextContent 
        v-if="message.messageType === 'text'"
        :message="message"
      />
      
      <ChatButtonContent 
        v-else-if="message.messageType === 'interactive_button'"
        :message="message"
      />
      
      <ChatURLContent 
        v-else-if="message.messageType === 'interactive_url'"
        :message="message"
      />
      
      <!-- Metadata -->
      <div
        v-if="message.metadata?.ticker"
        class="mt-2 text-xs opacity-75"
      >
        📊 {{ message.metadata.ticker }}
        <span v-if="message.metadata.confidence">
          • Confidence: {{ (message.metadata.confidence * 100).toFixed(0) }}%
        </span>
      </div>
      
      <!-- Timestamp -->
      <div class="text-xs opacity-75 mt-1">
        {{ formatTime(message.timestamp) }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { format } from 'date-fns'
import type { ChatMessage } from '@/stores/aiChat'
import ChatTextContent from './ChatTextContent.vue'
import ChatButtonContent from './ChatButtonContent.vue'
import ChatURLContent from './ChatURLContent.vue'

interface Props {
  message: ChatMessage
}

defineProps<Props>()

const formatTime = (timestamp: Date): string => {
  return format(timestamp, 'HH:mm')
}
</script>