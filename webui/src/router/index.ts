import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  {
    path: '/',
    name: 'home',
    component: () => import('@/views/LandingView.vue'),
    meta: { requiresGuest: true }
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
    component: () => import('@/views/SubscriptionView.vue')
  },
  {
    path: '/real-time-data',
    name: 'real-time-data',
    component: () => import('@/views/RealTimeDataView.vue'),
    meta: { requiresAuth: true }
  }/** 
  {
    path: '/profile',
    name: 'profile',
    component: () => import('@/views/ProfileView.vue')
  }*/
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Global navigation guards
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  console.log('🛡️ Router guard:', {
    to: to.name,
    from: from.name,
    isAuthenticated: authStore.isAuthenticated,
    user: !!authStore.user,
    accessToken: !!authStore.accessToken,
    isInitialized: authStore.isInitialized,
    requiresGuest: to.meta.requiresGuest,
    requiresAuth: to.meta.requiresAuth
  })
  
  // If auth store is not initialized yet, wait for it
  if (!authStore.isInitialized) {
    console.log('⏳ Auth store not initialized, waiting...')
    try {
      await authStore.initializeAuth()
      console.log('✅ Auth store initialized')
    } catch (error) {
      console.log('❌ Auth initialization failed:', error)
    }
  }
  
  // Redirect authenticated users away from auth pages and landing page
  if (to.meta.requiresGuest && authStore.isAuthenticated) {
    console.log('🔄 Redirecting authenticated user away from guest-only page')
    return next({ name: 'dashboard' })
  }
  
  // Redirect unauthenticated users away from protected pages
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    console.log('🔄 Redirecting unauthenticated user to login')
    return next({ name: 'login' })
  }
  
  console.log('✅ Router guard allowing navigation')
  next()
})

export default router