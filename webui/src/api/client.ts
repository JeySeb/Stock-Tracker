import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig } from 'axios'
import { cacheManager } from './cacheManager'

// Extend AxiosRequestConfig to include metadata
interface ExtendedAxiosRequestConfig extends AxiosRequestConfig {
  metadata?: {
    requestId: string
    timestamp: number
  }
  _retry?: boolean
}

class APIClient {
  private client: AxiosInstance
  private accessToken: string | null = null

  constructor() {
    console.log('Current MODE:', import.meta.env.MODE)
    // Use relative URL in development to leverage Vite proxy
    const baseURL = import.meta.env.DEV 
      ? '/api/v1' 
      : (import.meta.env.VITE_API_BASE_URL || import.meta.env.VITE_API_BASE_URL_DEV || 'http://localhost:8080/api/v1')
    console.log('Using API baseURL:', baseURL)
    
    this.client = axios.create({
      baseURL,
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json',
      },
    })

    this.setupInterceptors()
  }

  setAccessToken(token: string | null) {
    this.accessToken = token
  }

  private setupInterceptors() {
    // Request interceptor
    this.client.interceptors.request.use((config) => {
      if (this.accessToken) {
        config.headers.Authorization = `Bearer ${this.accessToken}`
      }

      // Add request ID and timestamp for tracking
      (config as ExtendedAxiosRequestConfig).metadata = { 
        requestId: crypto.randomUUID(),
        timestamp: Date.now()
      }

      return config
    })

    // Response interceptor
    this.client.interceptors.response.use(
      (response) => {
        // Log successful requests in development
        if (import.meta.env.DEV) {
          const extendedConfig = response.config as ExtendedAxiosRequestConfig
          const duration = Date.now() - (extendedConfig.metadata?.timestamp || 0)
          console.log(`✅ ${response.config.method?.toUpperCase()} ${response.config.url}`, {
            status: response.status,
            duration: `${duration}ms`,
            requestId: extendedConfig.metadata?.requestId
          })
        }
        return response
      },
      async (error) => {
        const originalRequest = error.config as ExtendedAxiosRequestConfig
        
        // Log errors in development
        if (import.meta.env.DEV) {
          const duration = Date.now() - (originalRequest.metadata?.timestamp || 0)
          console.error(`❌ ${originalRequest.method?.toUpperCase()} ${originalRequest.url}`, {
            status: error.response?.status,
            duration: `${duration}ms`,
            message: error.message,
            requestId: originalRequest.metadata?.requestId
          })
        }

        if (error.response?.status === 401 && !originalRequest._retry) {
          originalRequest._retry = true
          
          // Emit token expired event
          window.dispatchEvent(new CustomEvent('token-expired'))
          
          return Promise.reject(error)
        }
        
        return Promise.reject(error)
      }
    )
  }

  async get<T>(url: string, config?: AxiosRequestConfig & { cacheKey?: string; cacheTTL?: number }): Promise<T> {
    // Check cache if cacheKey provided
    if (config?.cacheKey && cacheManager.has(config.cacheKey)) {
      const cachedData = cacheManager.get<T>(config.cacheKey)
      if (cachedData) {
        if (import.meta.env.DEV) {
          console.log(`🎯 Cache hit for ${config.cacheKey}`)
        }
        return cachedData
      }
    }

    const response = await this.client.get(url, config)
    const data = response.data

    // Store in cache if cacheKey provided
    if (config?.cacheKey) {
      cacheManager.set(config.cacheKey, data, config.cacheTTL || 300)
      if (import.meta.env.DEV) {
        console.log(`💾 Cached ${config.cacheKey}`)
      }
    }

    return data
  }

  async post<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.client.post(url, data, config)
    return response.data
  }

  async put<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.client.put(url, data, config)
    return response.data
  }

  async delete<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.client.delete(url, config)
    return response.data
  }

  // Cache management methods
  invalidateCache(pattern: string): void {
    const stats = cacheManager.getStats()
    const keysToDelete = stats.keys.filter(key => key.includes(pattern))
    keysToDelete.forEach(key => cacheManager.delete(key))
    if (import.meta.env.DEV) {
      console.log(`🗑️ Invalidated ${keysToDelete.length} cache entries matching "${pattern}"`)
    }
  }

  clearCache(): void {
    cacheManager.clear()
    if (import.meta.env.DEV) {
      console.log('🗑️ Cleared entire cache')
    }
  }
}

export const apiClient = new APIClient()