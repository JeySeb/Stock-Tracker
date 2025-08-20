<template>
  <div>
    <!-- Message Content -->
    <div
      v-html="formatted"
      class="prose prose-sm max-w-none mb-3"
    />
    
    <!-- Interactive Buttons -->
    <div v-if="message.buttons && message.buttons.length > 0" class="space-y-2">
      <button
        v-for="button in message.buttons"
        :key="button.id"
        @click="button.action"
        class="block w-full px-4 py-2 text-sm font-medium text-primary-700 bg-primary-50 border border-primary-200 rounded-lg hover:bg-primary-100 hover:border-primary-300 transition-colors duration-200 text-left"
      >
        {{ button.title }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ChatMessage } from '@/stores/aiChat'
import { computed } from 'vue'
import { formatWhatsAppText } from '@/utils/whatsappFormat'

interface Props {
  message: ChatMessage
}

const props = defineProps<Props>()
const formatted = computed(() => formatWhatsAppText(props.message.content))
</script>