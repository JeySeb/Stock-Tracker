import { reactive, computed } from 'vue'

export interface AuthFormData {
  email: string
  password: string
  first_name?: string
  last_name?: string
  confirm_password?: string
  acceptTerms?: boolean
}

export interface AuthErrors {
  email: string
  password: string
  first_name?: string
  last_name?: string
  confirm_password?: string
  general?: string
}

export function useAuthValidation() {
  const errors = reactive<AuthErrors>({
    email: '',
    password: '',
    first_name: '',
    last_name: '',
    confirm_password: '',
    general: ''
  })

  // Email validation
  const validateEmail = (email: string): string => {
    if (!email.trim()) {
      return 'Email is required'
    }
    
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!emailRegex.test(email)) {
      return 'Please enter a valid email address'
    }
    
    return ''
  }

  // Password validation
  const validatePassword = (password: string): string => {
    if (!password) {
      return 'Password is required'
    }
    
    if (password.length < 8) {
      return 'Password must be at least 8 characters long'
    }
    
    return ''
  }

  // Name validation
  const validateName = (name: string, fieldName: string): string => {
    if (!name.trim()) {
      return `${fieldName} is required`
    }
    
    if (name.length > 100) {
      return `${fieldName} must be 100 characters or less`
    }
    
    return ''
  }

  // Confirm password validation
  const validateConfirmPassword = (password: string, confirmPassword: string): string => {
    if (!confirmPassword) {
      return 'Please confirm your password'
    }
    
    if (password !== confirmPassword) {
      return 'Passwords do not match'
    }
    
    return ''
  }

  // Password strength calculation
  const calculatePasswordStrength = (password: string): number => {
    if (!password) return 0
    
    let strength = 0
    
    // Length check
    if (password.length >= 8) strength++
    
    // Character variety checks
    if (/[a-z]/.test(password) && /[A-Z]/.test(password)) strength++
    if (/[0-9]/.test(password) && /[^A-Za-z0-9]/.test(password)) strength++
    
    return strength
  }

  // Get password strength color
  const getPasswordStrengthColor = (strength: number): string => {
    if (strength >= 3) return 'text-green-600'
    if (strength >= 2) return 'text-yellow-600'
    if (strength >= 1) return 'text-red-600'
    return 'text-gray-400'
  }

  // Get password strength text
  const getPasswordStrengthText = (strength: number): string => {
    if (strength >= 3) return 'Strong'
    if (strength >= 2) return 'Medium'
    if (strength >= 1) return 'Weak'
    return 'Very Weak'
  }

  // Validate registration form
  const validateRegistrationForm = (form: AuthFormData): boolean => {
    let isValid = true
    
    // Clear previous errors
    Object.keys(errors).forEach(key => {
      errors[key as keyof AuthErrors] = ''
    })
    
    // Validate each field
    errors.first_name = validateName(form.first_name || '', 'First name')
    errors.last_name = validateName(form.last_name || '', 'Last name')
    errors.email = validateEmail(form.email)
    errors.password = validatePassword(form.password)
    errors.confirm_password = validateConfirmPassword(form.password, form.confirm_password || '')
    
    // Check if any errors exist
    Object.values(errors).forEach(error => {
      if (error) isValid = false
    })
    
    // Check terms acceptance
    if (!form.acceptTerms) {
      errors.general = 'You must accept the terms and conditions'
      isValid = false
    }
    
    return isValid
  }

  // Validate login form
  const validateLoginForm = (form: AuthFormData): boolean => {
    let isValid = true
    
    // Clear previous errors
    errors.email = ''
    errors.password = ''
    errors.general = ''
    
    // Validate fields
    errors.email = validateEmail(form.email)
    errors.password = validatePassword(form.password)
    
    // Check if any errors exist
    if (errors.email || errors.password) {
      isValid = false
    }
    
    return isValid
  }

  // Handle API errors
  const handleApiError = (error: unknown): void => {
    // Clear previous errors
    Object.keys(errors).forEach(key => {
      errors[key as keyof AuthErrors] = ''
    })
    
    console.log('API Error received:', error)
    
    const errorObj = error as { 
      response?: { 
        status?: number; 
        data?: { 
          message?: string;
          error?: string;
          details?: unknown;
        } 
      };
      message?: string;
      code?: string;
    }
    
    // Handle 409 Conflict (user already exists)
    if (errorObj.response?.status === 409) {
      errors.email = 'An account with this email already exists'
      return
    }
    
    // Handle 401 Unauthorized
    if (errorObj.response?.status === 401) {
      errors.general = 'Invalid email or password'
      return
    }
    
    // Handle 400 Bad Request
    if (errorObj.response?.status === 400) {
      const message = errorObj.response.data?.message || errorObj.response.data?.error || errorObj.message
      if (message) {
        // Map API error messages to specific fields
        if (message.toLowerCase().includes('email') || message.toLowerCase().includes('already exists')) {
          errors.email = message
        } else if (message.toLowerCase().includes('password')) {
          errors.password = message
        } else if (message.toLowerCase().includes('first_name') || message.toLowerCase().includes('first name')) {
          errors.first_name = message
        } else if (message.toLowerCase().includes('last_name') || message.toLowerCase().includes('last name')) {
          errors.last_name = message
        } else {
          errors.general = message
        }
      } else {
        errors.general = 'Please check your input and try again.'
      }
      return
    }
    
    // Handle other response errors
    if (errorObj.response?.data?.message) {
      const message = errorObj.response.data.message
      if (message.toLowerCase().includes('email') || message.toLowerCase().includes('already exists')) {
        errors.email = message
      } else if (message.toLowerCase().includes('password')) {
        errors.password = message
      } else if (message.toLowerCase().includes('first_name') || message.toLowerCase().includes('first name')) {
        errors.first_name = message
      } else if (message.toLowerCase().includes('last_name') || message.toLowerCase().includes('last name')) {
        errors.last_name = message
      } else {
        errors.general = message
      }
      return
    }
    
    // Handle direct error messages
    if (errorObj.message) {
      if (errorObj.message.toLowerCase().includes('email') || errorObj.message.toLowerCase().includes('already exists')) {
        errors.email = errorObj.message
      } else {
        errors.general = errorObj.message
      }
      return
    }
    
    // Handle network errors
    if (errorObj.message?.includes('Network Error') || errorObj.message?.includes('ERR_NETWORK')) {
      errors.general = 'Network error: Unable to connect to the server. Please check your internet connection and try again.'
      return
    }
    
    // Handle timeout errors
    if (errorObj.message?.includes('timeout') || errorObj.code === 'ECONNABORTED') {
      errors.general = 'Request timeout: The server is taking too long to respond. Please try again.'
      return
    }
    
    // Fallback error
    errors.general = 'An unexpected error occurred. Please try again.'
  }

  // Clear all errors
  const clearErrors = (): void => {
    Object.keys(errors).forEach(key => {
      errors[key as keyof AuthErrors] = ''
    })
  }

  // Check if form has any errors
  const hasErrors = computed(() => {
    return Object.values(errors).some(error => error !== '')
  })

  return {
    errors,
    validateEmail,
    validatePassword,
    validateName,
    validateConfirmPassword,
    calculatePasswordStrength,
    getPasswordStrengthColor,
    getPasswordStrengthText,
    validateRegistrationForm,
    validateLoginForm,
    handleApiError,
    clearErrors,
    hasErrors
  }
}