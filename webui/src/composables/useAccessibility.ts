import { ref, onMounted, readonly } from 'vue'

export function useAccessibility() {
  const isHighContrast = ref(false)
  const isReducedMotion = ref(false)
  const fontSize = ref('normal')

  const checkAccessibilityPreferences = () => {
    // Check for high contrast preference
    isHighContrast.value = window.matchMedia('(prefers-contrast: high)').matches

    // Check for reduced motion preference
    isReducedMotion.value = window.matchMedia('(prefers-reduced-motion: reduce)').matches

    // Load font size preference
    const savedFontSize = localStorage.getItem('font-size-preference')
    if (savedFontSize) {
      fontSize.value = savedFontSize
      applyFontSize(savedFontSize)
    }
  }

  const applyFontSize = (size: string) => {
    const root = document.documentElement
    switch (size) {
      case 'small':
        root.style.setProperty('--font-size-multiplier', '0.875')
        break
      case 'large':
        root.style.setProperty('--font-size-multiplier', '1.125')
        break
      case 'extra-large':
        root.style.setProperty('--font-size-multiplier', '1.25')
        break
      default:
        root.style.setProperty('--font-size-multiplier', '1')
    }
  }

  const setFontSize = (size: string) => {
    fontSize.value = size
    applyFontSize(size)
    localStorage.setItem('font-size-preference', size)
  }

  const announceToScreenReader = (message: string) => {
    const announcement = document.createElement('div')
    announcement.setAttribute('aria-live', 'polite')
    announcement.setAttribute('aria-atomic', 'true')
    announcement.style.position = 'absolute'
    announcement.style.left = '-10000px'
    announcement.style.width = '1px'
    announcement.style.height = '1px'
    announcement.style.overflow = 'hidden'
    
    document.body.appendChild(announcement)
    announcement.textContent = message
    
    setTimeout(() => {
      document.body.removeChild(announcement)
    }, 1000)
  }

  const trapFocus = (element: HTMLElement) => {
    const focusableElements = element.querySelectorAll(
      'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
    )
    const firstElement = focusableElements[0] as HTMLElement
    const lastElement = focusableElements[focusableElements.length - 1] as HTMLElement

    const handleTab = (e: KeyboardEvent) => {
      if (e.key === 'Tab') {
        if (e.shiftKey) {
          if (document.activeElement === firstElement) {
            lastElement.focus()
            e.preventDefault()
          }
        } else {
          if (document.activeElement === lastElement) {
            firstElement.focus()
            e.preventDefault()
          }
        }
      }
    }

    element.addEventListener('keydown', handleTab)
    firstElement?.focus()

    return () => {
      element.removeEventListener('keydown', handleTab)
    }
  }

  onMounted(() => {
    checkAccessibilityPreferences()
  })

  return {
    isHighContrast: readonly(isHighContrast),
    isReducedMotion: readonly(isReducedMotion),
    fontSize: readonly(fontSize),
    setFontSize,
    announceToScreenReader,
    trapFocus
  }
}