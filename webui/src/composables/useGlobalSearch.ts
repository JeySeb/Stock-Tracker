import { ref, watch, readonly } from 'vue'
import { debounce } from 'lodash-es'
import { useStocksStore } from '@/stores/stocks'
import { useRecommendationsStore } from '@/stores/recommendations'
import type { StockEvent, Recommendation } from '@/types'

export interface SearchResult {
  type: 'stock' | 'recommendation' | 'company'
  id: string
  title: string
  subtitle: string
  description: string
  data: StockEvent | Recommendation | { name: string }
  relevance: number
}

export function useGlobalSearch() {
  const searchQuery = ref('')
  const isSearching = ref(false)
  const searchResults = ref<SearchResult[]>([])
  const recentSearches = ref<string[]>([])
  const searchHistory = ref<SearchResult[]>([])

  const stocksStore = useStocksStore()
  const recommendationsStore = useRecommendationsStore()

  // Load recent searches from localStorage
  const loadRecentSearches = () => {
    const stored = localStorage.getItem('stock_tracker_recent_searches')
    if (stored) {
      recentSearches.value = JSON.parse(stored).slice(0, 10)
    }
  }

  // Save search to history
  const saveSearch = (query: string, results: SearchResult[]) => {
    if (query.trim()) {
      recentSearches.value = [
        query,
        ...recentSearches.value.filter(s => s !== query)
      ].slice(0, 10)
      
      localStorage.setItem('stock_tracker_recent_searches', JSON.stringify(recentSearches.value))
      
      if (results.length > 0) {
        searchHistory.value = [
          results[0],
          ...searchHistory.value.filter(r => r.id !== results[0].id)
        ].slice(0, 20)
      }
    }
  }

  // Advanced search algorithm
  const performSearch = async (query: string): Promise<SearchResult[]> => {
    if (!query.trim()) return []

    isSearching.value = true
    const results: SearchResult[] = []

    try {
      // Search stocks
      const stockResults = await searchStocks(query)
      results.push(...stockResults)

      // Search recommendations
      const recommendationResults = await searchRecommendations(query)
      results.push(...recommendationResults)

      // Search companies
      const companyResults = await searchCompanies(query)
      results.push(...companyResults)

      // Sort by relevance
      return results.sort((a, b) => b.relevance - a.relevance).slice(0, 50)
    } finally {
      isSearching.value = false
    }
  }

  const searchStocks = async (query: string): Promise<SearchResult[]> => {
    await stocksStore.fetchStocks()
    const stocks = stocksStore.stocks

    return stocks
      .filter(stock => 
        stock.ticker.toLowerCase().includes(query.toLowerCase()) ||
        stock.company.toLowerCase().includes(query.toLowerCase()) ||
        stock.brokerage.toLowerCase().includes(query.toLowerCase())
      )
      .map(stock => ({
        type: 'stock' as const,
        id: stock.id,
        title: stock.ticker,
        subtitle: stock.company,
        description: `${stock.action} by ${stock.brokerage}`,
        data: stock,
        relevance: calculateStockRelevance(stock, query)
      }))
  }
  const searchRecommendations = async (query: string): Promise<SearchResult[]> => {
    await recommendationsStore.fetchRecommendations()
    const recommendations = recommendationsStore.recommendations

    return recommendations
      .filter(rec => 
        rec.ticker.toLowerCase().includes(query.toLowerCase()) ||
        rec.company_name.toLowerCase().includes(query.toLowerCase())
      )
      .map(rec => ({
        type: 'recommendation' as const,
        id: rec.id,
        title: rec.ticker,
        subtitle: rec.company_name,
        description: `${rec.recommendation_type} - Score: ${(rec.basic_score * 100).toFixed(0)}%`,
        data: {
          ...rec,
          scoring_factors: rec.scoring_factors.map(factor => ({ ...factor })) // Create mutable copy of each scoring factor
        },
        relevance: calculateRecommendationRelevance(rec, query)
      }))
  }

  const searchCompanies = async (query: string): Promise<SearchResult[]> => {
    // This would integrate with a company database or API
    // For now, we'll extract unique companies from stocks
    const companies = new Set(stocksStore.stocks.map(s => s.company))
    
    return Array.from(companies)
      .filter(company => company.toLowerCase().includes(query.toLowerCase()))
      .map(company => ({
        type: 'company' as const,
        id: company,
        title: company,
        subtitle: 'Company',
        description: 'View all events and recommendations',
        data: { name: company },
        relevance: calculateTextRelevance(company, query)
      }))
  }

  const calculateStockRelevance = (stock: StockEvent, query: string): number => {
    let score = 0
    const q = query.toLowerCase()
    
    if (stock.ticker.toLowerCase() === q) score += 100
    else if (stock.ticker.toLowerCase().startsWith(q)) score += 80
    else if (stock.ticker.toLowerCase().includes(q)) score += 60
    
    if (stock.company.toLowerCase().includes(q)) score += 40
    if (stock.brokerage.toLowerCase().includes(q)) score += 20
    
    // Boost recent events
    const eventDate = new Date(stock.event_time)
    const daysSinceEvent = (Date.now() - eventDate.getTime()) / (1000 * 60 * 60 * 24)
    if (daysSinceEvent < 7) score += 30
    else if (daysSinceEvent < 30) score += 15
    
    return score
  }

  const calculateRecommendationRelevance = (rec: Recommendation, query: string): number => {
    let score = 0
    const q = query.toLowerCase()
    
    if (rec.ticker.toLowerCase() === q) score += 100
    else if (rec.ticker.toLowerCase().startsWith(q)) score += 80
    else if (rec.ticker.toLowerCase().includes(q)) score += 60
    
    if (rec.company_name.toLowerCase().includes(q)) score += 40
    
    // Boost high-confidence recommendations
    score += rec.confidence * 20
    score += rec.basic_score * 15
    
    return score
  }

  const calculateTextRelevance = (text: string, query: string): number => {
    const t = text.toLowerCase()
    const q = query.toLowerCase()
    
    if (t === q) return 100
    if (t.startsWith(q)) return 80
    if (t.includes(q)) return 60
    return 0
  }

  // Debounced search
  const debouncedSearch = debounce(async (query: string) => {
    if (query.length < 2) {
      searchResults.value = []
      return
    }

    const results = await performSearch(query)
    searchResults.value = results
    saveSearch(query, results)
  }, 300)

  // Watch for search query changes
  watch(searchQuery, (newQuery) => {
    debouncedSearch(newQuery)
  })

  // Clear search
  const clearSearch = () => {
    searchQuery.value = ''
    searchResults.value = []
  }

  // Initialize
  loadRecentSearches()

  return {
    searchQuery,
    isSearching: readonly(isSearching),
    searchResults: readonly(searchResults),
    recentSearches: readonly(recentSearches),
    searchHistory: readonly(searchHistory),
    clearSearch,
    performSearch
  }
}