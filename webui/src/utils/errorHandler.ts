export interface AppError {
  id: string
  type: 'network' | 'validation' | 'auth' | 'permission' | 'unknown'
  message: string
  details?: unknown
  timestamp: Date
  userAgent?: string
  url?: string
  userId?: string
}

class ErrorHandler {
  private errors: AppError[] = []
  private maxErrors = 50

  logError(error: Error | AppError, context?: unknown) {
    const appError: AppError = this.normalizeError(error, context)
    
    // Add to error log
    this.errors.unshift(appError)
    if (this.errors.length > this.maxErrors) {
      this.errors = this.errors.slice(0, this.maxErrors)
    }

    // Log to console in development
    if (import.meta.env.DEV) {
      console.error('🚨 Application Error:', appError)
    }

    // In production, send to monitoring service
    if (import.meta.env.PROD) {
      this.sendToMonitoring(appError)
    }

    return appError
  }

  private normalizeError(error: Error | AppError, context?: unknown): AppError {
    if ('id' in error && 'type' in error) {
      return error as AppError
    }

    const appError: AppError = {
      id: crypto.randomUUID(),
      type: this.categorizeError(error),
      message: error.message || 'An unknown error occurred',
      details: context,
      timestamp: new Date(),
      userAgent: navigator.userAgent,
      url: window.location.href
    }

    return appError
  }

  private categorizeError(error: Error): AppError['type'] {
    const message = error.message.toLowerCase()
    
    if (message.includes('network') || message.includes('fetch')) {
      return 'network'
    }
    if (message.includes('unauthorized') || message.includes('401')) {
      return 'auth'
    }
    if (message.includes('forbidden') || message.includes('403')) {
      return 'permission'
    }
    if (message.includes('validation') || message.includes('invalid')) {
      return 'validation'
    }
    
    return 'unknown'
  }

  private async sendToMonitoring(error: AppError) {
    try {
      // Replace with actual monitoring service (Sentry, LogRocket, etc.)
      await fetch('/api/errors', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(error)
      })
    } catch (monitoringError) {
      console.error('Failed to send error to monitoring:', monitoringError)
    }
  }

  getRecentErrors(limit = 10): AppError[] {
    return this.errors.slice(0, limit)
  }

  clearErrors() {
    this.errors = []
  }

  // User-friendly error messages
  getUserMessage(error: AppError): string {
    switch (error.type) {
      case 'network':
        return 'Please check your internet connection and try again.'
      case 'auth':
        return 'Your session has expired. Please log in again.'
      case 'permission':
        return 'You don\'t have permission to perform this action.'
      case 'validation':
        return 'Please check your input and try again.'
      default:
        return 'Something went wrong. Please try again.'
    }
  }
}

export const errorHandler = new ErrorHandler()

// Global error boundary for uncaught errors
window.addEventListener('error', (event) => {
  errorHandler.logError(new Error(event.message), {
    filename: event.filename,
    lineno: event.lineno,
    colno: event.colno
  })
})

window.addEventListener('unhandledrejection', (event) => {
  errorHandler.logError(new Error(event.reason?.message || 'Unhandled promise rejection'), {
    reason: event.reason
  })
})