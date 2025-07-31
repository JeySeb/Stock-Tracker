interface CacheItem<T> {
  data: T
  timestamp: number
  ttl: number
}

class CacheManager {
  private cache = new Map<string, CacheItem<unknown>>()
  private maxSize = 100

  set<T>(key: string, data: T, ttlMinutes: number = 30): void {
    // Clean up expired items if cache is getting full
    if (this.cache.size >= this.maxSize) {
      this.cleanup()
    }

    this.cache.set(key, {
      data,
      timestamp: Date.now(),
      ttl: ttlMinutes * 60 * 1000
    })
  }

  get<T>(key: string): T | null {
    const item = this.cache.get(key)
    
    if (!item) return null
    
    // Check if expired
    if (Date.now() - item.timestamp > item.ttl) {
      this.cache.delete(key)
      return null
    }
    
    return item.data as T
  }

  has(key: string): boolean {
    const item = this.cache.get(key)
    if (!item) return false
    
    // Check if expired
    if (Date.now() - item.timestamp > item.ttl) {
      this.cache.delete(key)
      return false
    }
    
    return true
  }

  delete(key: string): void {
    this.cache.delete(key)
  }

  clear(): void {
    this.cache.clear()
  }

  private cleanup(): void {
    const now = Date.now()
    for (const [key, item] of this.cache.entries()) {
      if (now - item.timestamp > item.ttl) {
        this.cache.delete(key)
      }
    }
  }

  getStats() {
    return {
      size: this.cache.size,
      maxSize: this.maxSize,
      keys: Array.from(this.cache.keys())
    }
  }
}

export const cacheManager = new CacheManager()

// Cache key generators
export const generateCacheKey = {
  stocks: (filters: Record<string, unknown>) => `stocks:${JSON.stringify(filters)}`,
  recommendations: (filters: Record<string, unknown>, tier: string) => `recommendations:${tier}:${JSON.stringify(filters)}`,
  stocksByTicker: (ticker: string) => `stocks_ticker:${ticker.toUpperCase()}`,
  recommendationByTicker: (ticker: string, tier: string) => `recommendation:${ticker.toUpperCase()}:${tier}`,
  stats: () => 'stats:general'
}