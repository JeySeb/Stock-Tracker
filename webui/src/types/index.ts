// User & Authentication Types
export type UserTier = 'guest' | 'basic' | 'premium'

export interface User {
  id: string
  email: string
  first_name: string
  last_name: string
  tier: UserTier
  is_verified: boolean
  last_login?: string | null
  created_at: string
  updated_at: string
}

export interface AuthTokens {
  access_token: string
  refresh_token: string
  expires_in: number
}

export interface AuthResponse {
  user: User
  tokens: AuthTokens
}

// Stock Types
export interface StockEvent {
  id: string
  ticker: string
  company: string
  brokerage: string
  action: string
  rating_from: string
  rating_to: string
  target_from: number
  target_to: number
  event_time: string
  price_close: number | null
  created_at: string
}

// Scoring Factor Type
export interface ScoringFactor {
  name: string
  score: number
  weight: number
  explanation: string
  tier: 'basic' | 'enriched' | 'premium'
}

// Recommendation Types
export interface Recommendation {
  id: string
  ticker: string
  company_name: string
  total_events: number
  positive_events: number
  negative_events: number
  avg_target_change: number
  latest_target_price: number
  broker_consensus: string
  basic_score: number
  confidence: number
  recommendation_type: 'Strong Buy' | 'Buy' | 'Hold' | 'Sell' | 'Strong Sell'
  scoring_factors: readonly ScoringFactor[]
  tier: 'basic' | 'enriched' | 'premium'
  external_data?: ExternalData
  ai_insights?: AIInsights
  last_event_time: string
  created_at: string
  expires_at: string
}
//TODO: Check this model.
export interface ExternalData {
  current_price: number
  day_change: number
  day_change_percent: number
  price_change_24h: number
  volume: number
  market_cap: number
  pe_ratio?: number
  dividend_yield?: number
  week_52_high?: number
  week_52_low?: number
  avg_volume?: number
  last_updated: string
}

// TODO: CHECK THIS AT FINAL PHASE
export interface AIInsights {
  sentiment_score: number
  news_sentiment: string
  social_media_buzz: number
  technical_indicators: {
    rsi: number
    macd: string
    moving_averages: string
  }
  ai_prediction: string
  risk_assessment: string
}

// API Response Types
export interface PaginatedResponse<T> {
  data: T[]
  pagination: {
    page: number
    limit: number
    total_pages: number
    total_items: number
    has_next: boolean
    has_prev: boolean
  }
}

export interface RecommendationResponse {
  data: Recommendation[]
  meta: {
    count: number
    user_tier: UserTier
    features: string[]
    cache_hit: boolean
    generation_time: number
    rate_limit_remaining?: number
  }
}

// Subscription Types
export interface Subscription {
  id: string
  user_id: string
  plan: 'monthly' | 'yearly'
  status: 'pending' | 'active' | 'cancelled' | 'expired'
  price: number
  currency: string
  start_date: string
  end_date: string
  payment_reference?: string
  created_at: string
  updated_at: string
}