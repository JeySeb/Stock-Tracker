import { defineStore } from 'pinia'
import { ref, computed, readonly } from 'vue'
import { authAPI } from '@/api/auth'
import { apiClient } from '@/api/client'
import type { User, AuthTokens, UserTier } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref<User | null>(null)
  const accessToken = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)
  const isLoading = ref(false)
  const isInitialized = ref(false)

  // Getters
  const isAuthenticated = computed(() => !!user.value && !!accessToken.value)
  const userTier = computed((): UserTier => user.value?.tier || 'guest')
  const hasFeature = computed(() => (feature: string) => {
    const tierFeatures = {
      guest: ['basic_recommendations', 'market_analytics'],
      basic: ['basic_recommendations', 'market_analytics', 'real_time_data', 'external_apis'],
      premium: ['basic_recommendations', 'market_analytics', 'real_time_data', 'external_apis', 'ai_insights', 'sentiment_analysis']
    }
    return tierFeatures[userTier.value]?.includes(feature) || false
  })

  // Actions
  async function login(email: string, password: string) {
    isLoading.value = true
    try {
      const response = await authAPI.login({ email, password })
      setAuthData(response.user, response.tokens)
      return response
    } catch (error) {
      console.error('Login error in store:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  async function register(userData: Parameters<typeof authAPI.register>[0]) {
    isLoading.value = true
    try {
      const response = await authAPI.register(userData)
      setAuthData(response.user, response.tokens)
      return response
    } catch (error) {
      console.error('Registration error in store:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  async function refreshTokens() {
    if (!refreshToken.value) throw new Error('No refresh token available')
    
    try {
      const response = await authAPI.refreshToken(refreshToken.value)
      setTokens(response.tokens)
    } catch (error) {
      logout()
      throw error
    }
  }

  function setAuthData(userData: User, tokens: AuthTokens) {
    console.log('🔧 Setting auth data:', { userData, tokens: { ...tokens, access_token: tokens.access_token?.substring(0, 20) + '...' } })
    user.value = userData
    setTokens(tokens)
    console.log('✅ Auth data set successfully:', {
      userSet: !!user.value,
      tokenSet: !!accessToken.value,
      isAuthenticated: isAuthenticated.value
    })
  }

  function updateUserTier(tier: UserTier) {
    if (user.value) {
      user.value = { ...user.value, tier }
    }
  }

  function setTokens(tokens: AuthTokens) {
    console.log('🔑 Setting tokens:', { access_token: tokens.access_token?.substring(0, 20) + '...', refresh_token: tokens.refresh_token?.substring(0, 10) + '...' })
    accessToken.value = tokens.access_token
    refreshToken.value = tokens.refresh_token
    
    // Update API client with new token
    apiClient.setAccessToken(tokens.access_token)
    
    // Store in localStorage for persistence
    localStorage.setItem('stock_tracker_access_token', tokens.access_token)
    localStorage.setItem('stock_tracker_refresh_token', tokens.refresh_token)
    
    console.log('✅ Tokens set and stored:', {
      accessTokenSet: !!accessToken.value,
      refreshTokenSet: !!refreshToken.value,
      localStorage: {
        access: !!localStorage.getItem('stock_tracker_access_token'),
        refresh: !!localStorage.getItem('stock_tracker_refresh_token')
      }
    })
  }

  function logout() {
    user.value = null
    accessToken.value = null
    refreshToken.value = null
    isInitialized.value = false
    
    // Clear API client token
    apiClient.setAccessToken(null)
    
    // Clear localStorage
    localStorage.removeItem('stock_tracker_access_token')
    localStorage.removeItem('stock_tracker_refresh_token')
  }

  function initializeAuth() {
    if (isInitialized.value) {
      console.log('🔧 Auth store already initialized, skipping...')
      return Promise.resolve(user.value)
    }
    
    console.log('🔧 Initializing auth store...')
    
    // Restore tokens from localStorage
    const storedAccessToken = localStorage.getItem('stock_tracker_access_token')
    const storedRefreshToken = localStorage.getItem('stock_tracker_refresh_token')
    
    console.log('🔍 Found stored tokens:', {
      hasAccessToken: !!storedAccessToken,
      hasRefreshToken: !!storedRefreshToken
    })
    
    if (storedAccessToken && storedRefreshToken) {
      accessToken.value = storedAccessToken
      refreshToken.value = storedRefreshToken
      
      // Update API client with stored token
      apiClient.setAccessToken(storedAccessToken)
      
      console.log('🔄 Attempting to restore user session...')
      
      // Try to fetch current user data
      return authAPI.getCurrentUser()
        .then(userData => {
          console.log('✅ User session restored:', userData)
          user.value = userData
          isInitialized.value = true
          return userData
        })
        .catch((error) => {
          console.error('❌ Failed to restore user session:', error)
          logout()
          isInitialized.value = true
          throw new Error('Failed to restore user session')
        })
    } else {
      console.log('📝 No stored tokens found, user will need to login')
      isInitialized.value = true
      return Promise.resolve(null)
    }

    // Listen for token expiration events
    window.addEventListener('token-expired', async () => {
      try {
        await refreshTokens()
      } catch {
        logout()
      }
    })
  }

  function demoLogin() {
    // Simulate a demo user for testing UI styling
    const demoUser: User = {
      id: 'demo-user-123',
      email: 'demo@example.com',
      first_name: 'Demo',
      last_name: 'User',
      tier: 'premium' as UserTier,
      is_verified: true,
      last_login: new Date().toISOString(),
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    }
    
    const demoTokens: AuthTokens = {
      access_token: 'demo-access-token',
      refresh_token: 'demo-refresh-token',
      expires_in: 3600
    }
    
    setAuthData(demoUser, demoTokens)
  }

      return {
      // State
      user: readonly(user),
      isLoading: readonly(isLoading),
      isInitialized: readonly(isInitialized),
      
      // Getters
      isAuthenticated,
      userTier,
      hasFeature,
      accessToken: readonly(accessToken),
      
      // Actions
      login,
      register,
      refreshToken: refreshTokens,
      logout,
      initializeAuth,
      updateUserTier,
      demoLogin
    }
})