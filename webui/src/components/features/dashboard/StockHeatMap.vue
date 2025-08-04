<template>
  <div class="divide-y divide-gray-200">
    <!-- Loading State -->
    <div v-if="loading" class="space-y-3 p-6">
      <div class="grid grid-cols-4 gap-4">
        <div v-for="n in 8" :key="n" class="animate-pulse">
          <div class="h-32 bg-gray-200 shimmer"></div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div 
      v-else-if="!events || events.length === 0" 
      class="p-6 text-center text-gray-500"
    >
      No recent stock events
    </div>

    <!-- Heat Map -->
    <div v-else class="bg-white shadow-xl overflow-hidden">
      <!-- Header -->
      <div class="bg-gradient-to-r from-blue-600 to-indigo-700 p-6 text-white">
        <div class="flex flex-col md:flex-row md:items-center md:justify-between">
          <div>
            <h1 class="text-2xl font-bold">Stock Performance Heatmap</h1>
            <p class="opacity-90 mt-1">Real‑time visualization of <span class="font-semibold">most recent</span> market movements</p>
          </div>
          <!-- Time filter buttons - DISABLED
          <div class="mt-4 md:mt-0 flex space-x-2">
            <button class="time‑btn" :class="timeFrame==='1D' && 'active'" @click="timeFrame='1D'">1D</button>
            <button class="time‑btn" :class="timeFrame==='1W' && 'active'" @click="timeFrame='1W'">1W</button>
            <button class="time‑btn" :class="timeFrame==='1M' && 'active'" @click="timeFrame='1M'">1M</button>
            <button class="time‑btn" :class="timeFrame==='3M' && 'active'" @click="timeFrame='3M'">3M</button>
          </div>
          -->
        </div>
      </div>

      <!-- Legend + Filters -->
      <div class="p-6 border-b border-gray-200">
        <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
          <!-- Legend -->
          <div class="flex items-center space-x-4">
            <div class="flex items-center space-x-2">
              <div class="w-5 h-5 bg-gradient-to-br from-green-500 to-green-600"></div>
              <span class="text-sm font-medium">Growth (5%+)</span>
            </div>
            <div class="flex items-center space-x-2">
              <div class="w-5 h-5 bg-gradient-to-br from-green-300 to-green-400"></div>
              <span class="text-sm font-medium">Slight Growth</span>
            </div>
            <div class="flex items-center space-x-2">
              <div class="w-5 h-5 bg-gray-200"></div>
              <span class="text-sm font-medium">Neutral</span>
            </div>
            <div class="flex items-center space-x-2">
              <div class="w-5 h-5 bg-gradient-to-br from-red-300 to-red-400"></div>
              <span class="text-sm font-medium">Slight Decline</span>
            </div>
            <div class="flex items-center space-x-2">
              <div class="w-5 h-5 bg-gradient-to-br from-red-500 to-red-600"></div>
              <span class="text-sm font-medium">Decline (5%+)</span>
            </div>
          </div>
          <!-- Filters -->
          <div class="flex flex-col sm:flex-row sm:items-center space-y-2 sm:space-y-0 sm:space-x-4">
            <div class="flex items-center space-x-2">
              <button 
                @click="filterType = 'all'"
                :class="filterType === 'all' ? 'bg-blue-600 text-white' : 'bg-gray-200 text-gray-700'"
                class="px-3 py-1 text-xs transition-colors"
              >
                All
              </button>
              <button 
                @click="filterType = 'positive'"
                :class="filterType === 'positive' ? 'bg-green-600 text-white' : 'bg-gray-200 text-gray-700'"
                class="px-3 py-1 text-xs transition-colors"
              >
                Positive
              </button>
              <button 
                @click="filterType = 'negative'"
                :class="filterType === 'negative' ? 'bg-red-600 text-white' : 'bg-gray-200 text-gray-700'"
                class="px-3 py-1 text-xs transition-colors"
              >
                Negative
              </button>
            </div>
            <div class="text-xs text-gray-500">
              {{ filteredEvents.length }} of {{ props.events?.length || 0 }} stocks
            </div>
          </div>
        </div>
      </div>

      <!-- Metro Heatmap (d3‑treemap / squarify) -->
      <div class="p-6">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8">
          <div
            v-for="(box, i) in bigBoxesWithLayout"
            :key="`bigbox-${i}`"
            class="big-box-container"
          >
            <div v-if="props.categories" class="mb-3">
              <h3 class="text-lg font-semibold text-gray-800">{{ box.category }}</h3>
              <p class="text-sm text-gray-500">{{ box.leaves.length }} stocks</p>
            </div>

            <div class="big-box relative bg-gray-50 border border-gray-300 overflow-hidden">
              <div
                v-for="leaf in box.leaves"
                :key="leaf.event.id"
                class="absolute tile cursor-pointer group"
                :style="getLeafStyle(leaf)"
                @click="handleStockClick(leaf.event)"
              >
                <div class="flex flex-col justify-between h-full p-1">
                  <span class="font-bold text-[0.65rem] truncate" :class="getTextClass(leaf.event)">
                    {{ leaf.event.ticker }}
                  </span>
                  <span class="text-[0.6rem]" :class="getSubTextClass(leaf.event)">
                    {{ getChangePercentage(leaf.event) }}
                  </span>
                </div>
                <!-- Tooltip -->
                <div class="tooltip" >
                  <div class="font-bold">{{ leaf.event.ticker }} - {{ leaf.event.company }}</div>
                  <div class="text-gray-600">${{ leaf.event.target_to?.toFixed(2) || 'N/A' }} ({{ getChangePercentage(leaf.event) }})</div>
                  <div class="text-gray-600">{{ leaf.event.brokerage }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Summary -->
      <div class="p-6 bg-gray-50 border-t border-gray-200">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div class="bg-white p-4 shadow-sm border border-gray-100">
            <div class="text-gray-500 text-sm">Displayed Stocks</div>
            <div class="text-2xl font-bold text-gray-800">{{ filteredEvents.length }}</div>
          </div>
          <div class="bg-white p-4 shadow-sm border border-gray-100">
            <div class="text-gray-500 text-sm">Positive Movers</div>
            <div class="text-2xl font-bold text-green-600">{{ positiveCount }}</div>
          </div>
          <div class="bg-white p-4 shadow-sm border border-gray-100">
            <div class="text-gray-500 text-sm">Negative Movers</div>
            <div class="text-2xl font-bold text-red-600">{{ negativeCount }}</div>
          </div>
          <div class="bg-white p-4 shadow-sm border border-gray-100">
            <div class="text-gray-500 text-sm">Neutral</div>
            <div class="text-2xl font-bold text-gray-600">{{ neutralCount }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal -->
    <StockDetailModal
      :is-open="showModal"
      :ticker="selectedStock?.ticker"
      :stock-data="selectedStock"
      @close="closeModal"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
// @ts-ignore
import { treemap, hierarchy, treemapSquarify } from 'd3-hierarchy'
import StockDetailModal from '@/components/features/stocks/StockDetailModal.vue'
import type { StockEvent } from '@/types'

/* ──────────────────────────────────────────────────────────────────────── */
/* ░░  Props & Interfaces                                                   */
/* ──────────────────────────────────────────────────────────────────────── */
interface Props {
  events?: StockEvent[]
  loading?: boolean
  categories?: CategoryGroup[]
}
interface CategoryGroup { name: string; tickers: string[] }
interface LeafRect { event: StockEvent; x: number; y: number; w: number; h: number }

const props = withDefaults(defineProps<Props>(), { loading: false })

/* ──────────────────────────────────────────────────────────────────────── */
/* ░░  UI State                                                             */
/* ──────────────────────────────────────────────────────────────────────── */
// const timeFrame = ref<'1D'|'1W'|'1M'|'3M'>('1W') // DISABLED
const filterType = ref<'all'|'positive'|'negative'>('all')
const showModal     = ref(false)
const selectedStock = ref<StockEvent | null>(null)

/* ──────────────────────────────────────────────────────────────────────── */
/* ░░  Helpers                                                              */
/* ──────────────────────────────────────────────────────────────────────── */
function getChangeValue (e: StockEvent): number {
  if (!e.target_from || !e.target_to) return 0
  return ((e.target_to - e.target_from) / e.target_from) * 100
}
function getChangePercentage (e: StockEvent): string {
  const v = getChangeValue(e)
  return v === 0 ? '0.0%' : `${v > 0 ? '+' : ''}${v.toFixed(1)}%`
}

/* Text color helpers */
function getTextClass (e: StockEvent) {
  const c = getChangeValue(e)
  if (c === 0) return 'text-gray-800'
  if (Math.abs(c) >= 5) return 'text-white'
  return 'text-gray-800'
}
function getSubTextClass (e: StockEvent) {
  const c = getChangeValue(e)
  if (c === 0) return 'text-gray-600'
  if (Math.abs(c) >= 5) return c > 0 ? 'text-green-100' : 'text-red-100'
  return 'text-gray-600'
}
/* Background‑color generator (diverging palette) */
function tileBg (e: StockEvent): string {
  const c = getChangeValue(e)
  if (c === 0) return '#d1d5db' // neutral gray‑300
  const intensity = Math.min(Math.abs(c) / 10, 1) // cap at 10%
  if (c > 0) {
    return `rgba(34,197,94,${0.4 + 0.6 * intensity})` // green‑500 hue varying alpha
  }
  return `rgba(239,68,68,${0.4 + 0.6 * intensity})` // red‑500
}

/* ──────────────────────────────────────────────────────────────────────── */
/* ░░  Event Lists (sort & filter)                                          */
/* ──────────────────────────────────────────────────────────────────────── */

// Helper function to parse dates more robustly - DISABLED
/*
function parseDate(dateString: string): Date | null {
  if (!dateString) return null
  
  try {
    // Try parsing as ISO string first
    const date = new Date(dateString)
    if (!isNaN(date.getTime())) {
      return date
    }
    
    // Try parsing common date formats
    const formats = [
      /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/, // ISO-like
      /^\d{4}-\d{2}-\d{2}/, // YYYY-MM-DD
      /^\d{2}\/\d{2}\/\d{4}/, // MM/DD/YYYY
      /^\d{1,2}\/\d{1,2}\/\d{4}/, // M/D/YYYY
    ]
    
    for (const format of formats) {
      if (format.test(dateString)) {
        const parsed = new Date(dateString)
        if (!isNaN(parsed.getTime())) {
          return parsed
        }
      }
    }
    
    return null
  } catch (error) {
    console.warn('Failed to parse date:', dateString, error)
    return null
  }
}

// Helper function to filter events by timeframe - DISABLED
function filterByTimeFrame(events: StockEvent[], timeFrame: string): StockEvent[] {
  if (!events || events.length === 0) return []
  
  const now = new Date()
  let cutoffDate: Date
  
  switch (timeFrame) {
    case '1D':
      cutoffDate = new Date(now.getTime() - 24 * 60 * 60 * 1000) // 24 hours ago
      break
    case '1W':
      cutoffDate = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000) // 7 days ago
      break
    case '1M':
      cutoffDate = new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000) // 30 days ago
      break
    case '3M':
      cutoffDate = new Date(now.getTime() - 90 * 24 * 60 * 60 * 1000) // 90 days ago
      break
    default:
      return events // Return all events if timeframe is not recognized
  }
  
  // Debug: Log first few events to see their date formats
  if (events.length > 0) {
    console.log('Debug - First 3 events:', events.slice(0, 3).map(e => ({
      ticker: e.ticker,
      event_time: e.event_time,
      created_at: e.created_at,
      event_time_parsed: e.event_time ? parseDate(e.event_time)?.toISOString() || 'invalid' : 'null',
      created_at_parsed: e.created_at ? parseDate(e.created_at)?.toISOString() || 'invalid' : 'null'
    })))
    console.log('Debug - Cutoff date:', cutoffDate.toISOString(), 'Timeframe:', timeFrame)
  }
  
  const filtered = events.filter(event => {
    // Try event_time first, then fallback to created_at
    const dateString = event.event_time || event.created_at
    if (!dateString) {
      return false // Skip events without any timestamp
    }
    
    const eventDate = parseDate(dateString)
    if (!eventDate) {
      console.warn('Could not parse date:', dateString, 'for ticker:', event.ticker)
      return false
    }
    
    return eventDate >= cutoffDate
  })
  
  console.log('Debug - Filtered count:', filtered.length, 'of', events.length)
  return filtered
}

const timeFilteredEvents = computed(() => {
  if (!props.events) return []
  const filtered = filterByTimeFrame(props.events, timeFrame.value)
  
  // If no events match the timeframe, show all events (fallback)
  if (filtered.length === 0 && props.events.length > 0) {
    console.log(`No events match timeframe ${timeFrame.value}, showing all ${props.events.length} events as fallback`)
    return props.events
  }
  
  return filtered
})
*/

const sortedEvents = computed(() => {
  return [...(props.events || [])]
    .sort((a, b) => Math.abs(getChangeValue(b)) - Math.abs(getChangeValue(a)))
    .slice(0, 60)
})

const filteredEvents = computed(() => {
  if (filterType.value === 'all') return sortedEvents.value
  return sortedEvents.value.filter(e => filterType.value === 'positive' ? getChangeValue(e) > 0 : getChangeValue(e) < 0)
})

/* ──────────────────────────────────────────────────────────────────────── */
/* ░░  Grouping + d3‑treemap layout                                         */
/* ──────────────────────────────────────────────────────────────────────── */
function createBigBoxes (stocks: StockEvent[], cats?: CategoryGroup[]) {
  const groups: { category: string; items: StockEvent[] }[] = []
  if (cats && cats.length) {
    cats.forEach(c => {
      const items = stocks.filter(s => c.tickers.includes(s.ticker))
      if (items.length) groups.push({ category: c.name, items })
    })
  } else {
    // fallback: 8 roughly even groups
    const copy = [...stocks]
    const n = 8
    const per = Math.ceil(copy.length / n)
    for (let i = 0; i < n && copy.length; i++) {
      groups.push({ category: `Group ${i+1}`, items: copy.splice(0, per) })
    }
  }
  return groups
}

function layoutTreemap (items: StockEvent[]): LeafRect[] {
  if (!items.length) return []
  
  // Create a wrapper object that contains the items as children
  const wrapper = { children: items }
  const root = hierarchy(wrapper)
    .sum((d: any) => {
      // Handle both the wrapper object and StockEvent objects
      if (d.children) return 0 // wrapper object has no value
      return Math.abs(getChangeValue(d as StockEvent)) + 0.1 // ensure non‑zero
    })
  
  const treemapLayout = treemap()
    .size([100, 100])
    .paddingInner(1)
    .round(true)
    .tile(treemapSquarify)
  
  // Use type assertion to bypass TypeScript strict checking for d3-hierarchy
  treemapLayout(root as any)
  
  return root.leaves().map((l: any) => ({
    event: l.data as StockEvent,
    x: l.x0,
    y: l.y0,
    w: l.x1 - l.x0,
    h: l.y1 - l.y0
  }))
}

const bigBoxesWithLayout = computed(() => {
  return createBigBoxes(filteredEvents.value, props.categories).map(g => ({
    category: g.category,
    leaves: layoutTreemap(g.items)
  }))
})

/* ──────────────────────────────────────────────────────────────────────── */
/* ░░  Summary counts                                                       */
/* ──────────────────────────────────────────────────────────────────────── */
const positiveCount = computed(() => filteredEvents.value.filter(e => getChangeValue(e) > 0).length)
const negativeCount = computed(() => filteredEvents.value.filter(e => getChangeValue(e) < 0).length)
const neutralCount  = computed(() => filteredEvents.value.filter(e => getChangeValue(e) === 0).length)

// Watch for timeframe changes to trigger reactivity - DISABLED
/*
watch(timeFrame, () => {
  // This will trigger recomputation of all dependent computed properties
})
*/

/* ──────────────────────────────────────────────────────────────────────── */
/* ░░  Style helpers                                                        */
/* ──────────────────────────────────────────────────────────────────────── */
function getLeafStyle (leaf: LeafRect) {
  return {
    left: `${leaf.x}%`,
    top: `${leaf.y}%`,
    width: `${leaf.w}%`,
    height: `${leaf.h}%`,
    background: tileBg(leaf.event)
  }
}

/* ──────────────────────────────────────────────────────────────────────── */
/* ░░  Modal handlers                                                       */
/* ──────────────────────────────────────────────────────────────────────── */
function handleStockClick (e: StockEvent) {
  selectedStock.value = e
  showModal.value = true
}
function closeModal () {
  showModal.value = false
  selectedStock.value = null
}

/* Misc */
function formatLastUpdated () {
  return new Date().toLocaleString('en-US', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit', timeZoneName: 'short' })
}
function getNextRefreshTime () {
  const now = new Date()
  const next = new Date(now.getTime() + 60_000)
  const diff = Math.ceil((next.getTime() - now.getTime()) / 1000)
  return `${Math.floor(diff / 60)}:${String(diff % 60).padStart(2, '0')}`
}
</script>

<style scoped>
/***** Utility Components *****/
.time‑btn { @apply px-4 py-2 bg-white/20 hover:bg-white/30 transition text-sm rounded; }
.time‑btn.active { @apply bg-white/30 font-semibold; }

.big-box { aspect-ratio: 1/1; border-radius: 0; transition: border-color .2s ease; }
.big-box:hover { border-color: #3b82f6; }

.tile { border-radius: 0; transition: transform .15s ease; overflow: hidden; border: 1px solid rgba(0,0,0,0.05); }
.tile:hover { transform: scale(1.03); z-index: 2; }

.tooltip {
  @apply absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-3 py-2 bg-white border border-gray-200 text-gray-800 text-xs shadow-lg opacity-0 group-hover:opacity-100 transition-all duration-150 pointer-events-none whitespace-nowrap z-30;
}

@keyframes shimmer { 0% { background-position: -1000px 0; } 100% { background-position: 1000px 0; } }
.shimmer { background: linear-gradient(to right, #f0f0f0 4%, #e0e0e0 25%, #f0f0f0 36%); background-size: 1000px 100%; animation: shimmer 2s infinite linear; }
</style>