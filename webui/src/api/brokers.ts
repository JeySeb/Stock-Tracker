import { apiClient } from './client'

export interface BrokerScore {
  id: string
  name: string
  credibility_score: number
  report_count: number
  calculated_score: number
  created_at: string
  updated_at: string
}

export interface BrokerScoresResponse {
  data: BrokerScore[]
}

export const brokersAPI = {
  async getBrokerScores(): Promise<BrokerScoresResponse> {
    return apiClient.get('/brokers/scores')
  },

  async getBrokerStocks(brokerName: string): Promise<any> {
    // This would be used for broker-specific stock analysis
    const params = new URLSearchParams({ brokerage: brokerName })
    return apiClient.get(`/stocks?${params.toString()}`)
  }
}