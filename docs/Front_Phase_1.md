# Frontend Development - Phase 1: Foundation & Core Infrastructure

## Phase Overview

**Duration:** 2-3 weeks  
**Focus:** Establish solid foundation with authentication, routing, state management, and API integration layer  
**Goal:** Create a secure, scalable base that supports tier-based functionality and future AI integration

## Justification for Phase 1 Structure

### Why Start Here?
1. **Authentication First:** Financial platforms require robust security from day one
2. **Tier System Foundation:** User tiers affect every feature, so this must be established early
3. **API Integration Layer:** Centralized API handling ensures consistency and maintainability
4. **Scalable Architecture:** Proper project structure prevents technical debt as features grow

### Critical Success Factors
- Secure JWT token management with refresh mechanism
- Tier-aware component rendering system
- Modular API service architecture
- Responsive layout foundation

---

## Step-by-Step Implementation

### Step 1: Project Setup & Configuration

#### 1.1 Initialize Vue 3 Project
```bash
cd Stock-Tracker
npm create vue@latest webui -- --typescript --router --pinia --eslint --prettier
cd webui
npm install
```

#### 1.2 Install Additional Dependencies
```bash
# UI & Styling
npm install @headlessui/vue @heroicons/vue tailwindcss @tailwindcss/forms @tailwindcss/typography

# HTTP Client & Utilities
npm install axios @vueuse/core date-fns

# Charts (for future phases)
npm install echarts vue-echarts

# Development & Testing
npm install -D @types/node vitest jsdom @vitest/ui
```

#### 1.3 Configure Tailwind CSS
Create `tailwind.config.js`:
```javascript
/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#eff6ff',
          500: '#3b82f6',
          600: '#2563eb',
          700: '#1d4ed8',
        },
        success: {
          50: '#f0fdf4',
          500: '#22c55e',
          600: '#16a34a',
        },
        danger: {
          50: '#fef2f2',
          500: '#ef4444',
          600: '#dc2626',
        },
        financial: {
          buy: '#10b981',
          sell: '#ef4444',
          hold: '#f59e0b',
          neutral: '#6b7280',
        }
      }
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ],
}
```

#### 1.4 Project Structure Setup
```
src/
├── api/                    # API integration layer
│   ├── client.ts          # Axios configuration
│   ├── auth.ts            # Authentication endpoints
│   ├── stocks.ts          # Stocks endpoints
│   ├── recommendations.ts # Recommendations endpoints
│   └── subscriptions.ts   # Subscription endpoints
├── components/
│   ├── ui/                # Reusable UI components
│   ├── layout/            # Layout components
│   └── features/          # Feature-specific components
├── composables/           # Vue composables
├── stores/                # Pinia stores
├── types/                 # TypeScript type definitions
├── utils/                 # Utility functions
├── views/                 # Route views
└── assets/               # Static assets
```

### Step 2: Type Definitions & API Models

#### 2.1 Core Types (`src/types/index.ts`)
```typescript
// User & Authentication Types
export type UserTier = 'guest' | 'basic' | 'premium'

export interface User {
  id: string
  email: string
  first_name: string
  last_name: string
  tier: UserTier
  is_verified: boolean
  last_login: string | null
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
  scoring_factors: ScoringFactor[]
  tier: 'basic' | 'enriched' | 'premium'
  external_data?: ExternalData
  ai_insights?: AIInsights
  last_event_time: string
  created_at: string
  expires_at: string
}

export interface ExternalData {
  current_price: number
  price_change_24h: number
  volume: number
  market_cap: number
  pe_ratio: number
  analyst_ratings: {
    strong_buy: number
    buy: number
    hold: number
    sell: number
    strong_sell: number
  }
}

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
```

### Step 3: API Integration Layer

#### 3.1 HTTP Client Setup (`src/api/client.ts`)
```typescript
import axios, { AxiosInstance, AxiosRequestConfig } from 'axios'
import { useAuthStore } from '@/stores/auth'

class APIClient {
  private client: AxiosInstance

  constructor() {
    this.client = axios.create({
      baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json',
      },
    })

    this.setupInterceptors()
  }

  private setupInterceptors() {
    // Request interceptor to add auth token
    this.client.interceptors.request.use((config) => {
      const authStore = useAuthStore()
      if (authStore.accessToken) {
        config.headers.Authorization = `Bearer ${authStore.accessToken}`
      }
      return config
    })

    // Response interceptor for token refresh
    this.client.interceptors.response.use(
      (response) => response,
      async (error) => {
        const originalRequest = error.config
        
        if (error.response?.status === 401 && !originalRequest._retry) {
          originalRequest._retry = true
          
          const authStore = useAuthStore()
          try {
            await authStore.refreshToken()
            return this.client(originalRequest)
          } catch (refreshError) {
            authStore.logout()
            return Promise.reject(refreshError)
          }
        }
        
        return Promise.reject(error)
      }
    )
  }

  async get<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.client.get(url, config)
    return response.data
  }

  async post<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.client.post(url, data, config)
    return response.data
  }

  async put<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.client.put(url, data, config)
    return response.data
  }

  async delete<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.client.delete(url, config)
    return response.data
  }
}

export const apiClient = new APIClient()
```

#### 3.2 Authentication API (`src/api/auth.ts`)
```typescript
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
```

### Step 4: State Management with Pinia

#### 4.1 Authentication Store (`src/stores/auth.ts`)
```typescript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authAPI } from '@/api/auth'
import type { User, AuthTokens, UserTier } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref<User | null>(null)
  const accessToken = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)
  const isLoading = ref(false)

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
    user.value = userData
    setTokens(tokens)
  }

  function setTokens(tokens: AuthTokens) {
    accessToken.value = tokens.access_token
    refreshToken.value = tokens.refresh_token
    
    // Store in localStorage for persistence
    localStorage.setItem('stock_tracker_access_token', tokens.access_token)
    localStorage.setItem('stock_tracker_refresh_token', tokens.refresh_token)
  }

  function logout() {
    user.value = null
    accessToken.value = null
    refreshToken.value = null
    
    // Clear localStorage
    localStorage.removeItem('stock_tracker_access_token')
    localStorage.removeItem('stock_tracker_refresh_token')
  }

  function initializeAuth() {
    // Restore tokens from localStorage
    const storedAccessToken = localStorage.getItem('stock_tracker_access_token')
    const storedRefreshToken = localStorage.getItem('stock_tracker_refresh_token')
    
    if (storedAccessToken && storedRefreshToken) {
      accessToken.value = storedAccessToken
      refreshToken.value = storedRefreshToken
      
      // Try to fetch current user data
      authAPI.getCurrentUser()
        .then(userData => user.value = userData)
        .catch(() => logout())
    }
  }

  return {
    // State
    user: readonly(user),
    isLoading: readonly(isLoading),
    
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
    initializeAuth
  }
})
```

### Step 5: Router Configuration

#### 5.1 Router Setup (`src/router/index.ts`)
```typescript
import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import type { UserTier } from '@/types'

// Route guards
const requireAuth = () => {
  const authStore = useAuthStore()
  if (!authStore.isAuthenticated) {
    return { name: 'login' }
  }
}

const requireTier = (minTier: UserTier) => () => {
  const authStore = useAuthStore()
  const tierLevels = { guest: 0, basic: 1, premium: 2 }
  
  if (!authStore.isAuthenticated) {
    return { name: 'login' }
  }
  
  if (tierLevels[authStore.userTier] < tierLevels[minTier]) {
    return { name: 'subscription' }
  }
}

const routes = [
  {
    path: '/',
    name: 'home',
    redirect: '/dashboard'
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/auth/LoginView.vue'),
    meta: { requiresGuest: true }
  },
  {
    path: '/register',
    name: 'register',
    component: () => import('@/views/auth/RegisterView.vue'),
    meta: { requiresGuest: true }
  },
  {
    path: '/dashboard',
    name: 'dashboard',
    component: () => import('@/views/DashboardView.vue')
  },
  {
    path: '/stocks',
    name: 'stocks',
    component: () => import('@/views/StocksView.vue')
  },
  {
    path: '/recommendations',
    name: 'recommendations',
    component: () => import('@/views/RecommendationsView.vue')
  },
  {
    path: '/subscription',
    name: 'subscription',
    component: () => import('@/views/SubscriptionView.vue'),
    beforeEnter: requireAuth()
  },
  {
    path: '/profile',
    name: 'profile',
    component: () => import('@/views/ProfileView.vue'),
    beforeEnter: requireAuth()
  },
  // AI Placeholder routes for future phases
  {
    path: '/ai-insights',
    name: 'ai-insights',
    component: () => import('@/views/placeholder/AIInsightsView.vue'),
    beforeEnter: requireTier('premium')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Global navigation guards
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  // Redirect authenticated users away from auth pages
  if (to.meta.requiresGuest && authStore.isAuthenticated) {
    return next({ name: 'dashboard' })
  }
  
  next()
})

export default router
```

### Step 6: Core Layout Components

#### 6.1 Main Layout (`src/components/layout/AppLayout.vue`)
```vue
<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <AppHeader />
    
    <!-- Main Content -->
    <div class="flex">
      <!-- Sidebar -->
      <AppSidebar v-if="showSidebar" />
      
      <!-- Content Area -->
      <main class="flex-1 p-6">
        <router-view />
      </main>
    </div>
    
    <!-- AI Chat Placeholder (Future Phase) -->
    <AIChatPlaceholder />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import AppHeader from './AppHeader.vue'
import AppSidebar from './AppSidebar.vue'
import AIChatPlaceholder from '../placeholder/AIChatPlaceholder.vue'

const route = useRoute()

const showSidebar = computed(() => {
  const hiddenSidebarRoutes = ['login', 'register']
  return !hiddenSidebarRoutes.includes(route.name as string)
})
</script>
```

#### 6.2 Header Component (`src/components/layout/AppHeader.vue`)
```vue
<template>
  <header class="bg-white shadow-sm border-b border-gray-200">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center h-16">
        <!-- Logo -->
        <div class="flex items-center">
          <h1 class="text-2xl font-bold text-gray-900">
            Stock Tracker
          </h1>
        </div>
        
        <!-- User Controls -->
        <div class="flex items-center space-x-4">
          <!-- Tier Badge -->
          <TierBadge v-if="authStore.isAuthenticated" :tier="authStore.userTier" />
          
          <!-- User Menu -->
          <UserMenu v-if="authStore.isAuthenticated" />
          
          <!-- Auth Buttons -->
          <div v-else class="space-x-2">
            <router-link
              to="/login"
              class="text-gray-700 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium"
            >
              Login
            </router-link>
            <router-link
              to="/register"
              class="bg-primary-600 text-white hover:bg-primary-700 px-3 py-2 rounded-md text-sm font-medium"
            >
              Sign Up
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import TierBadge from '@/components/ui/TierBadge.vue'
import UserMenu from '@/components/ui/UserMenu.vue'

const authStore = useAuthStore()
</script>
```

### Step 7: Essential UI Components

#### 7.1 Tier Badge (`src/components/ui/TierBadge.vue`)
```vue
<template>
  <span :class="badgeClasses">
    {{ tierDisplay }}
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { UserTier } from '@/types'

interface Props {
  tier: UserTier
}

const props = defineProps<Props>()

const tierDisplay = computed(() => {
  const displays = {
    guest: 'Guest',
    basic: 'Basic',
    premium: 'Premium'
  }
  return displays[props.tier]
})

const badgeClasses = computed(() => {
  const baseClasses = 'inline-flex items-center px-3 py-1 rounded-full text-xs font-medium'
  const tierClasses = {
    guest: 'bg-gray-100 text-gray-800',
    basic: 'bg-blue-100 text-blue-800',
    premium: 'bg-purple-100 text-purple-800'
  }
  return `${baseClasses} ${tierClasses[props.tier]}`
})
</script>
```

### Step 8: Authentication Views

#### 8.1 Login View (`src/views/auth/LoginView.vue`)
```vue
<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
          Sign in to your account
        </h2>
      </div>
      
      <form class="mt-8 space-y-6" @submit.prevent="handleLogin">
        <div class="space-y-4">
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700">
              Email address
            </label>
            <input
              id="email"
              v-model="form.email"
              type="email"
              required
              class="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500 focus:z-10 sm:text-sm"
              placeholder="Enter your email"
            />
          </div>
          
          <div>
            <label for="password" class="block text-sm font-medium text-gray-700">
              Password
            </label>
            <input
              id="password"
              v-model="form.password"
              type="password"
              required
              class="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500 focus:z-10 sm:text-sm"
              placeholder="Enter your password"
            />
          </div>
        </div>

        <div>
          <button
            type="submit"
            :disabled="authStore.isLoading"
            class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50"
          >
            <span v-if="!authStore.isLoading">Sign in</span>
            <span v-else>Signing in...</span>
          </button>
        </div>

        <div class="text-center">
          <router-link to="/register" class="text-primary-600 hover:text-primary-500">
            Don't have an account? Sign up
          </router-link>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const form = reactive({
  email: '',
  password: ''
})

async function handleLogin() {
  try {
    await authStore.login(form.email, form.password)
    router.push('/dashboard')
  } catch (error) {
    console.error('Login failed:', error)
    // Handle error (show toast, etc.)
  }
}
</script>
```

---

## Testing Strategy for Phase 1

### Unit Tests
- Authentication store actions and getters
- API client interceptors
- Utility functions

### Integration Tests
- Login/register flow
- Token refresh mechanism
- Route guards

### E2E Tests
- Complete authentication workflow
- Navigation between authenticated/unauthenticated states

---

## Phase 1 Deliverables

✅ **Completed Setup:**
- Project initialization with proper tooling
- Type-safe API integration layer
- Secure authentication system with JWT refresh
- Tier-aware state management
- Responsive layout foundation
- Router with proper guards

✅ **Key Features:**
- User registration and login
- Automatic token refresh
- Tier-based access control
- Responsive header and navigation
- Protected route handling

✅ **AI Readiness:**
- Placeholder components for future AI features
- Extensible store architecture
- Reserved routes for AI functionality

---

## Next Steps to Phase 2

Phase 1 provides the solid foundation needed for Phase 2, which will focus on:
- Stocks data visualization and filtering
- Recommendations system with tier-specific features
- Dashboard with charts and analytics
- Subscription management UI

The authentication and state management established in Phase 1 will seamlessly support the data-heavy features coming in Phase 2. 