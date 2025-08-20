import { defineStore } from 'pinia'
import { ref, computed, readonly } from 'vue'
import { useAuthStore } from './auth'
import { aiChatAPI, type EnhancedChatMessage } from '@/api/aiChat'

// Re-export for backward compatibility
export interface ChatMessage extends EnhancedChatMessage {}

// Enhanced message interface
export interface ChatMessage {
  id: string
  type: 'user' | 'assistant' | 'system'
  content: string
  timestamp: Date
  messageType: 'text' | 'interactive_button' | 'interactive_url'
  isTyping?: boolean
  
  // For interactive buttons
  buttons?: Array<{
    id: string
    title: string
    action: () => void
  }>
  
  // For interactive URLs
  urlAction?: {
    headerText: string
    buttonText: string
    url: string
  }
  
  metadata?: {
    ticker?: string
    actionType?: 'analysis' | 'recommendation' | 'general'
    confidence?: number
  }
}

export interface ChatSession {
  id: string
  title: string
  createdAt: Date
  lastMessageAt: Date
  messages: ChatMessage[]
}

export const useAIChatStore = defineStore('aiChat', () => {
  // State
  const isOpen = ref(false)
  const currentSession = ref<ChatSession | null>(null)
  const sessions = ref<ChatSession[]>([])
  const isTyping = ref(false)
  const isConnected = ref(false)
  const connectionStatus = ref<'disconnected' | 'connecting' | 'connected'>('disconnected')

  // Getters
  const authStore = useAuthStore()
  const canUseAI = computed(() => authStore.hasFeature('ai_insights'))
  const hasActiveSessions = computed(() => sessions.value.length > 0)
  const currentMessages = computed(() => currentSession.value?.messages || [])

  // Actions
  function openChat() {
    if (!canUseAI.value) {
      // Show upgrade prompt
      return false
    }
    isOpen.value = true
    
    // Create new session if none exists
    if (!currentSession.value) {
      createNewSession()
    }
    
    return true
  }

  function closeChat() {
    isOpen.value = false
  }

  function toggleChat() {
    if (isOpen.value) {
      closeChat()
    } else {
      openChat()
    }
  }

  function createNewSession() {
    const session: ChatSession = {
      id: crypto.randomUUID(),
      title: `Chat ${new Date().toLocaleDateString()}`,
      createdAt: new Date(),
      lastMessageAt: new Date(),
      messages: [{
        id: crypto.randomUUID(),
        type: 'system',
        content: '¡Hola! Soy Mia, tu asistente de inversiones. ¿En qué puedo ayudarte?',
        timestamp: new Date(),
        messageType: 'text'
      }]
    }

    sessions.value.unshift(session)
    currentSession.value = session
    saveSessionsToStorage()
  }

  function switchToSession(sessionId: string) {
    const session = sessions.value.find(s => s.id === sessionId)
    if (session) {
      currentSession.value = session
    }
  }

  function deleteSession(sessionId: string) {
    sessions.value = sessions.value.filter(s => s.id !== sessionId)
    
    if (currentSession.value?.id === sessionId) {
      currentSession.value = sessions.value[0] || null
    }
    
    saveSessionsToStorage()
  }

  async function sendMessage(content: string, metadata?: ChatMessage['metadata']) {
    if (!currentSession.value || !canUseAI.value) return

    // Add user message
    const userMessage: ChatMessage = {
      id: crypto.randomUUID(),
      type: 'user',
      content,
      timestamp: new Date(),
      messageType: 'text',
      metadata
    }

    currentSession.value.messages.push(userMessage)
    currentSession.value.lastMessageAt = new Date()

    // Show typing indicator
    isTyping.value = true

    try {
      // Use the new LangGraph-Interpreter API
      const authStore = useAuthStore()
      const userId = authStore.user?.id || 'anonymous'
      
      const response = await aiChatAPI.sendMessage(content, userId)
      
      // Process and transform API response messages
      const enhancedMessages = aiChatAPI.processAPIResponse(response, handleButtonClick)
      
      // Add all assistant messages to the session
      enhancedMessages.forEach(message => {
        currentSession.value!.messages.push(message)
      })
      
      currentSession.value.lastMessageAt = new Date()

      // Update session title if it's the first meaningful exchange
      if (currentSession.value.messages.filter(m => m.type === 'user').length === 1) {
        currentSession.value.title = generateSessionTitle(content)
      }

    } catch (error) {
      console.error('Failed to get AI response:', error)
      
      const errorMessage: ChatMessage = {
        id: crypto.randomUUID(),
        type: 'assistant',
        content: 'I apologize, but I\'m having trouble processing your request right now. Please try again later.',
        timestamp: new Date(),
        messageType: 'text'
      }

      currentSession.value.messages.push(errorMessage)
    } finally {
      isTyping.value = false
      saveSessionsToStorage()
    }
  }

  // Handle button clicks from interactive messages
  function handleButtonClick(buttonId: string, buttonTitle: string) {
    // Send the button selection as a new message
    sendMessage(buttonTitle, {
      actionType: 'general',
      buttonId,
      isButtonResponse: true
    } as any)
  }



  function generateSessionTitle(firstMessage: string): string {
    const maxLength = 30
    if (firstMessage.length <= maxLength) {
      return firstMessage
    }
    return firstMessage.substring(0, maxLength - 3) + '...'
  }

  function saveSessionsToStorage() {
    try {
      localStorage.setItem('ai_chat_sessions', JSON.stringify(sessions.value))
    } catch (error) {
      console.error('Failed to save chat sessions:', error)
    }
  }

  function loadSessionsFromStorage() {
    try {
      const stored = localStorage.getItem('ai_chat_sessions')
      if (stored) {
        const parsedSessions = JSON.parse(stored) as Array<{
          id: string
          title: string
          createdAt: string
          lastMessageAt: string
          messages: Array<{
            id: string
            type: string
            content: string
            timestamp: string
            messageType?: string
            buttons?: Array<{ id: string; title: string }>
            urlAction?: { headerText: string; buttonText: string; url: string }
            metadata?: ChatMessage['metadata']
          }>
        }>
        sessions.value = parsedSessions.map((session) => ({
          ...session,
          createdAt: new Date(session.createdAt),
          lastMessageAt: new Date(session.lastMessageAt),
          messages: session.messages.map((message) => ({
            ...message,
            type: message.type as 'user' | 'assistant' | 'system',
            messageType: message.messageType as 'text' | 'interactive_button' | 'interactive_url' || 'text',
            timestamp: new Date(message.timestamp),
            // Restore button actions if they exist
            buttons: message.buttons?.map(btn => ({
              ...btn,
              action: () => handleButtonClick(btn.id, btn.title)
            }))
          }))
        }))
        
        if (sessions.value.length > 0) {
          currentSession.value = sessions.value[0]
        }
      }
    } catch (error) {
      console.error('Failed to load chat sessions:', error)
    }
  }

  // Initialize
  loadSessionsFromStorage()

  return {
    // State
    isOpen: readonly(isOpen),
    currentSession: readonly(currentSession),
    sessions: readonly(sessions),
    isTyping: readonly(isTyping),
    isConnected: readonly(isConnected),
    connectionStatus: readonly(connectionStatus),
    
    // Getters
    canUseAI,
    hasActiveSessions,
    currentMessages,
    
    // Actions
    openChat,
    closeChat,
    toggleChat,
    createNewSession,
    switchToSession,
    deleteSession,
    sendMessage,
    handleButtonClick
  }
})