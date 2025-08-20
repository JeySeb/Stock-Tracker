<template>
  <div>
    <div
      v-if="message.type === 'assistant'"
      v-html="formatted"
      class="prose prose-sm max-w-none"
    />
    <div v-else>{{ message.content }}</div>
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