import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import { useAuthStore } from '@/stores/auth'

// ECharts configuration
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { PieChart } from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent
} from 'echarts/components'

// Register ECharts components
use([
  CanvasRenderer,
  PieChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent
])

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

// Initialize auth store after pinia is set up
const authStore = useAuthStore()

// Wait for auth initialization before mounting the app
async function initializeApp() {
  try {
    console.log('🚀 Initializing app...')
    await authStore.initializeAuth()
    console.log('✅ Auth initialized, mounting app...')
    app.mount('#app')
  } catch (error) {
    console.log('⚠️ Auth initialization failed, mounting app anyway...', error)
    app.mount('#app')
  }
}

initializeApp()
