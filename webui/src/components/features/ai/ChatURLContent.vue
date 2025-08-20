<template>
  <div>
    <!-- Header Text -->
    <div 
      v-if="message.urlAction?.headerText"
      class="font-semibold text-gray-900 mb-2"
    >
      {{ message.urlAction.headerText }}
    </div>
    
    <!-- Message Content -->
    <div
      v-html="formatted"
      class="prose prose-sm max-w-none mb-3"
    />
    
    <!-- URL Action Button -->
    <div v-if="message.urlAction">
      <a
        :href="message.urlAction.url"
        target="_blank"
        rel="noopener noreferrer"
        class="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-primary-600 border border-transparent rounded-lg hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-colors duration-200"
      >
        <span>{{ message.urlAction.buttonText }}</span>
        <ArrowTopRightOnSquareIcon class="ml-2 h-4 w-4" />
      </a>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ArrowTopRightOnSquareIcon } from '@heroicons/vue/24/outline'
import type { ChatMessage } from '@/stores/aiChat'
import { computed } from 'vue'
import { formatWhatsAppText } from '@/utils/whatsappFormat'

interface Props {
  message: ChatMessage
}

const props = defineProps<Props>()
const formatted = computed(() => formatWhatsAppText(props.message.content))
</script>