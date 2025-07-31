interface PerformanceMetric {
    name: string
    value: number
    timestamp: number
    type: 'navigation' | 'resource' | 'custom'
  }
  
  class PerformanceMonitor {
    private metrics: PerformanceMetric[] = []
    private observer?: PerformanceObserver
  
    init() {
      // Monitor Core Web Vitals
      this.observeWebVitals()
      
      // Monitor resource loading
      this.observeResources()
      
      // Monitor navigation timing
      this.observeNavigation()
    }
  
    private observeWebVitals() {
      // Observe Largest Contentful Paint
      this.observeMetric('largest-contentful-paint', (entry) => {
        const lcpEntry = entry as PerformanceEntry & { value: number }
        this.recordMetric('LCP', lcpEntry.value, 'navigation')
      })

      // Observe First Input Delay
      this.observeMetric('first-input', (entry) => {
        const fidEntry = entry as PerformanceEntry & { processingStart: number; startTime: number }
        this.recordMetric('FID', fidEntry.processingStart - fidEntry.startTime, 'navigation')
      })

      // Observe Cumulative Layout Shift
      this.observeMetric('layout-shift', (entry) => {
        const clsEntry = entry as PerformanceEntry & { value: number; hadRecentInput: boolean }
        if (!clsEntry.hadRecentInput) {
          this.recordMetric('CLS', clsEntry.value, 'navigation')
        }
      })
    }

    private observeResources() {
      this.observeMetric('resource', (entry) => {
        const resourceEntry = entry as PerformanceEntry & { loadEventEnd: number; fetchStart: number }
        this.recordMetric('Resource Load Time', resourceEntry.loadEventEnd - resourceEntry.fetchStart, 'resource')
      })
    }

    private observeNavigation() {
      this.observeMetric('navigation', (entry) => {
        const navEntry = entry as PerformanceEntry & { 
          domContentLoadedEventEnd: number; 
          fetchStart: number;
          responseEnd: number;
          loadEventEnd: number;
        }
        
        this.recordMetric('Page Load Time', navEntry.loadEventEnd - navEntry.fetchStart, 'navigation')
        this.recordMetric('DOM Content Loaded', navEntry.domContentLoadedEventEnd - navEntry.fetchStart, 'navigation')
        this.recordMetric('First Paint', navEntry.responseEnd - navEntry.fetchStart, 'navigation')
      })
    }
  
    private observeMetric(type: string, callback: (entry: PerformanceEntry) => void) {
      if (!PerformanceObserver) return
  
      try {
        const observer = new PerformanceObserver((list) => {
          list.getEntries().forEach(callback)
        })
        observer.observe({ entryTypes: [type] })
      } catch (error) {
        console.warn(`Failed to observe ${type}:`, error)
      }
    }
  
    recordMetric(name: string, value: number, type: PerformanceMetric['type'] = 'custom') {
      const metric: PerformanceMetric = {
        name,
        value,
        timestamp: Date.now(),
        type
      }
  
      this.metrics.push(metric)
      
      // Keep only recent metrics (last 100)
      if (this.metrics.length > 100) {
        this.metrics = this.metrics.slice(-100)
      }
  
      // Log critical metrics
      if (value > this.getThreshold(name)) {
        console.warn(`⚠️ Performance issue detected: ${name} = ${value}ms`)
      }
  
      // Send to analytics in production
      if (import.meta.env.PROD) {
        this.sendToAnalytics(metric)
      }
    }
  
    private getThreshold(metricName: string): number {
      const thresholds = {
        'LCP': 2500,
        'FID': 100,
        'CLS': 0.1,
        'Page Load Time': 3000,
        'DOM Content Loaded': 1500
      }
      return thresholds[metricName as keyof typeof thresholds] || Infinity
    }
  
    private async sendToAnalytics(metric: PerformanceMetric) {
      try {
        // Replace with actual analytics service
        await fetch('/api/analytics/performance', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(metric)
        })
      } catch (error) {
        console.error('Failed to send performance metric:', error)
      }
    }
  
    getMetrics(): PerformanceMetric[] {
      return [...this.metrics]
    }
  
    getAverageMetric(name: string): number {
      const relevantMetrics = this.metrics.filter(m => m.name === name)
      if (relevantMetrics.length === 0) return 0
      
      const sum = relevantMetrics.reduce((acc, metric) => acc + metric.value, 0)
      return sum / relevantMetrics.length
    }
  
    generateReport(): string {
      const report = {
        timestamp: new Date().toISOString(),
        metrics: this.metrics.reduce((acc, metric) => {
          acc[metric.name] = {
            latest: metric.value,
            average: this.getAverageMetric(metric.name),
            count: this.metrics.filter(m => m.name === metric.name).length
          }
          return acc
        }, {} as Record<string, { latest: number; average: number; count: number }>)
      }
  
      return JSON.stringify(report, null, 2)
    }
  }
  
  export const performanceMonitor = new PerformanceMonitor()
  
  // Initialize when the module loads
  performanceMonitor.init()