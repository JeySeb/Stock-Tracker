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
  }
  /**,
  {
    path: '/register',
    name: 'register',
    component: () => import('@/views/auth/RegisterView.vue'),
    meta: { requiresGuest: true }
  }*/,
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
  /** 
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
  } */
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