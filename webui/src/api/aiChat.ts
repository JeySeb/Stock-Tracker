import axios from 'axios'
import type { AxiosResponse } from 'axios'

// LangGraph-Interpreter Chat API Types
export interface ChatInteractionRequest {
  user_message: string
  input_format: 'plain_text'
  output_format: 'webui'
  user_id: string
}

export interface WhatsAppTextMessage {
  type: 'text'
  text: {
    body: string
  }
}

export interface WhatsAppButtonAction {
  type: 'reply'
  reply: {
    id: string
    title: string
  }
}

export interface WhatsAppInteractiveButton {
  type: 'interactive'
  interactive: {
    type: 'button'
    body: {
      text: string
    }
    action: {
      buttons: WhatsAppButtonAction[]
    }
  }
}

export interface WhatsAppInteractiveURL {
  type: 'interactive'
  interactive: {
    type: 'cta_url'
    header: {
      type: 'text'
      text: string
    }
    body: {
      text: string
    }
    action: {
      name: 'cta_url'
      parameters: {
        display_text: string
        url: string
      }
    }
  }
}

export type WhatsAppMessage = WhatsAppTextMessage | WhatsAppInteractiveButton | WhatsAppInteractiveURL

export interface ChatInteractionResponse {
  response: WhatsAppMessage[]
  format: string
  user_id: string
  processed_count: number
}

// Enhanced Chat Message Types for Frontend
export interface EnhancedChatMessage {
  id: string
  type: 'user' | 'assistant' | 'system'
  content: string
  timestamp: Date
  messageType: 'text' | 'interactive_button' | 'interactive_url'
  
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
  
  // Original metadata
  metadata?: {
    ticker?: string
    actionType?: 'analysis' | 'recommendation' | 'general'
    confidence?: number
  }
}

class AIChatAPI {
  private baseURL: string
  
  constructor() {
    // LangGraph-Interpreter API URL - use proxy in development
    this.baseURL = import.meta.env.DEV 
      ? '' // Use Vite proxy in development
      : (import.meta.env.VITE_LANGGRAPH_API_URL || 'http://0.0.0.0:8000')
  }

  /**
   * Send a chat message to the LangGraph-Interpreter API
   */
  async sendMessage(message: string, userId: string): Promise<ChatInteractionResponse> {
    try {
      const requestData: ChatInteractionRequest = {
        user_message: message,
        input_format: 'plain_text',
        output_format: 'webui',
        user_id: userId
      }

      // Allow long-running responses (default 3 minutes, configurable)
      const timeoutMs = Number((import.meta as any).env?.VITE_CHAT_REQUEST_TIMEOUT_MS || 180000)

      const response: AxiosResponse<ChatInteractionResponse> = await axios.post(
        `${this.baseURL}/chat/chat-interaction`,
        requestData,
        {
          headers: {
            'Content-Type': 'application/json'
          },
          timeout: timeoutMs
        }
      )

      return response.data
    } catch (error) {
      console.error('Error sending chat message:', error)
      throw new Error('Failed to send message to AI chat service')
    }
  }

  /**
   * Transform WhatsApp message format to enhanced chat message
   */
  transformToEnhancedMessage(
    whatsappMessage: WhatsAppMessage, 
    messageId: string,
    onButtonClick: (buttonId: string, buttonTitle: string) => void
  ): EnhancedChatMessage {
    const baseMessage: Partial<EnhancedChatMessage> = {
      id: messageId,
      type: 'assistant',
      timestamp: new Date()
    }

    switch (whatsappMessage.type) {
      case 'text':
        return {
          ...baseMessage,
          messageType: 'text',
          content: whatsappMessage.text.body
        } as EnhancedChatMessage

      case 'interactive':
        if (whatsappMessage.interactive.type === 'button') {
          const buttons = whatsappMessage.interactive.action.buttons.map(button => ({
            id: button.reply.id,
            title: button.reply.title,
            action: () => onButtonClick(button.reply.id, button.reply.title)
          }))

          return {
            ...baseMessage,
            messageType: 'interactive_button',
            content: whatsappMessage.interactive.body.text,
            buttons
          } as EnhancedChatMessage
        } else if (whatsappMessage.interactive.type === 'cta_url') {
          return {
            ...baseMessage,
            messageType: 'interactive_url',
            content: whatsappMessage.interactive.body.text,
            urlAction: {
              headerText: whatsappMessage.interactive.header.text,
              buttonText: whatsappMessage.interactive.action.parameters.display_text,
              url: whatsappMessage.interactive.action.parameters.url
            }
          } as EnhancedChatMessage
        }
        break
    }

    // Fallback for unknown message types
    return {
      ...baseMessage,
      messageType: 'text',
      content: 'Unsupported message format received'
    } as EnhancedChatMessage
  }

  /**
   * Process raw API response and transform all messages
   */
  processAPIResponse(
    response: ChatInteractionResponse,
    onButtonClick: (buttonId: string, buttonTitle: string) => void
  ): EnhancedChatMessage[] {
    return response.response.map((message, index) => 
      this.transformToEnhancedMessage(
        message, 
        crypto.randomUUID(), 
        onButtonClick
      )
    )
  }
}

export const aiChatAPI = new AIChatAPI()