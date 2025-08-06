import { defineStore } from 'pinia'
import { ref, computed, readonly } from 'vue'
import { brokersAPI } from '@/api/brokers'
import type { BrokerScore } from '@/api/brokers'

export const useBrokersStore = defineStore('brokers', () => {
  // State
  const brokerScores = ref<BrokerScore[]>([])
  const isLoading = ref(false)

  // Getters
  const topBrokersByReportCount = computed(() => {
    if (!brokerScores.value || brokerScores.value.length === 0) return []
    return [...brokerScores.value]
      .sort((a, b) => b.report_count - a.report_count)
      .slice(0, 10)
  })

  const topBrokersByScore = computed(() => {
    if (!brokerScores.value || brokerScores.value.length === 0) return []
    return [...brokerScores.value]
      .sort((a, b) => b.calculated_score - a.calculated_score)
      .slice(0, 10)
  })

  // Actions
  async function fetchBrokerScores() {
    isLoading.value = true
    try {
      const response = await brokersAPI.getBrokerScores()
      brokerScores.value = response.data
      return response
    } catch (error) {
      console.error('Failed to fetch broker scores:', error)
      brokerScores.value = []
    } finally {
      isLoading.value = false
    }
  }

  return {
    // State
    brokerScores: readonly(brokerScores),
    isLoading: readonly(isLoading),
    
    // Getters
    topBrokersByReportCount,
    topBrokersByScore,
    
    // Actions
    fetchBrokerScores
  }
})