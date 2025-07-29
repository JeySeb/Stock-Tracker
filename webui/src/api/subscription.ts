import { apiClient } from './client'
import type { Subscription } from '@/types'

export interface CreateSubscriptionRequest {
  plan: 'monthly' | 'yearly'
}

export const subscriptionAPI = {
  async getCurrentSubscription(): Promise<Subscription> {
    return apiClient.get('/subscriptions/current')
  },

  async createSubscription(data: CreateSubscriptionRequest): Promise<Subscription> {
    return apiClient.post('/subscriptions', data)
  },

  async processPayment(subscriptionId: string): Promise<{ message: string }> {
    return apiClient.post(`/subscriptions/${subscriptionId}/payment`)
  }
}