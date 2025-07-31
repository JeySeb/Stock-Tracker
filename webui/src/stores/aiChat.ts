import { defineStore } from 'pinia'
import { ref, computed, readonly } from 'vue'
import { useAuthStore } from './auth'

export interface ChatMessage {
  id: string
  type: 'user' | 'assistant' | 'system'
  content: string
  timestamp: Date
  isTyping?: boolean
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
        content: 'Hello! I\'m your AI investment assistant. I can help you analyze stocks, understand market trends, and provide personalized recommendations. What would you like to know?',
        timestamp: new Date()
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
      metadata
    }

    currentSession.value.messages.push(userMessage)
    currentSession.value.lastMessageAt = new Date()

    // Show typing indicator
    isTyping.value = true

    try {
      // TODO: Replace with actual AI API call
      const response = await simulateAIResponse(content)
      
      const assistantMessage: ChatMessage = {
        id: crypto.randomUUID(),
        type: 'assistant',
        content: response.content,
        timestamp: new Date(),
        metadata: response.metadata
      }

      currentSession.value.messages.push(assistantMessage)
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
        timestamp: new Date()
      }

      currentSession.value.messages.push(errorMessage)
    } finally {
      isTyping.value = false
      saveSessionsToStorage()
    }
  }

  // Simulate AI response (replace with actual API call)
  async function simulateAIResponse(userMessage: string) {
    // Simulate network delay
    await new Promise(resolve => setTimeout(resolve, 1000 + Math.random() * 2000))

    const lowerMessage = userMessage.toLowerCase()

    if (lowerMessage.includes('aapl') || lowerMessage.includes('apple')) {
      return {
        content: `Based on my analysis of Apple (AAPL), here's what I found:

📊 **Current Analysis:**
- Recent price target increases from major brokerages
- Strong institutional buying pressure
- Q4 earnings beat expectations by 12%

💡 **Key Insights:**
- 73% of brokerages rate it as "Buy" or "Strong Buy"
- Average price target: $185 (current: $175)
- Risk level: Moderate

🔮 **My Recommendation:**
Consider this stock for your portfolio. The fundamentals look strong, and technical indicators suggest continued upward momentum.

Would you like me to analyze any specific aspects like competitors, valuation metrics, or market sentiment?`,
        metadata: {
          ticker: 'AAPL',
          actionType: 'analysis' as const,
          confidence: 0.85
        }
      }
    }

    if (lowerMessage.includes('recommendation') || lowerMessage.includes('suggest')) {
      return {
        content: `Here are my top 3 stock recommendations based on current market conditions:

🥇 **Microsoft (MSFT)** - Strong Buy
- AI growth story with solid fundamentals
- Price target: $420 (upside: 15%)
- Risk: Low

🥈 **NVIDIA (NVDA)** - Buy
- Leading AI chip manufacturer
- Price target: $520 (upside: 22%)
- Risk: Moderate-High

🥉 **Amazon (AMZN)** - Buy
- AWS growth and retail recovery
- Price target: $165 (upside: 18%)
- Risk: Moderate

These recommendations are based on broker consensus, technical analysis, and my proprietary sentiment scoring. Would you like detailed analysis on any of these?`,
        metadata: {
          actionType: 'recommendation' as const,
          confidence: 0.78
        }
      }
    }

    // Default response
    return {
      content: `I understand you're asking about "${userMessage}". Let me help you with that.

As your AI investment assistant, I can provide insights on:
• Stock analysis and recommendations
• Market trend interpretation
• Portfolio optimization suggestions
• Risk assessment
• Earnings analysis

Could you be more specific about what you'd like to know? For example:
- "Analyze TSLA stock"
- "What are the best tech stocks right now?"
- "Should I buy or sell Netflix?"`,
      metadata: {
        actionType: 'general' as const
      }
    }
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
            timestamp: new Date(message.timestamp)
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
    sendMessage
  }
})