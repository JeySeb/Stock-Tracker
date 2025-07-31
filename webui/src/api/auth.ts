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
    console.log('🔐 AuthAPI: Attempting login with:', { ...credentials, password: '[HIDDEN]' })
    try {
      const response = await apiClient.post<AuthResponse>('/auth/login', credentials)
      console.log('✅ AuthAPI: Login successful:', response)
      return response
    } catch (error) {
      console.error('❌ AuthAPI: Login failed:', error)
      throw error
    }
  },

  async register(userData: RegisterRequest): Promise<AuthResponse> {
    console.log('🔐 AuthAPI: Attempting registration with:', { ...userData, password: '[HIDDEN]' })
    try {
      const response = await apiClient.post<AuthResponse>('/auth/register', userData)
      console.log('✅ AuthAPI: Registration successful:', response)
      return response
    } catch (error) {
      console.error('❌ AuthAPI: Registration failed:', error)
      throw error
    }
  },

  async refreshToken(refreshToken: string): Promise<{ tokens: AuthResponse['tokens'] }> {
    return apiClient.post('/auth/refresh', { refresh_token: refreshToken })
  },

  async getCurrentUser(): Promise<User> {
    return apiClient.get('/auth/me')
  }
}