import { apiClient } from './client'
import type { AuthResponse, User } from '@/types'

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
  first_name: string
  last_name: string
}

export interface RefreshTokenRequest {
  refresh_token: string
}

export const authAPI = {
  async login(credentials: LoginRequest): Promise<AuthResponse> {
    return apiClient.post('/auth/login', credentials)
  },

  async register(userData: RegisterRequest): Promise<AuthResponse> {
    return apiClient.post('/auth/register', userData)
  },

  async refreshToken(refreshToken: string): Promise<{ tokens: AuthResponse['tokens'] }> {
    return apiClient.post('/auth/refresh', { refresh_token: refreshToken })
  },

  async getCurrentUser(): Promise<User> {
    return apiClient.get('/auth/me')
  }
}