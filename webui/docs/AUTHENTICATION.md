# Authentication System Documentation

## Overview

The Stock Tracker web application implements a comprehensive authentication system with user registration, login, and session management. The system is built using Vue 3 Composition API, Pinia for state management, and integrates with a RESTful API backend.

## Architecture

### Components Structure

```
src/
├── views/auth/
│   ├── LoginView.vue          # User login interface
│   └── RegisterView.vue       # User registration interface
├── stores/
│   └── auth.ts               # Authentication state management
├── api/
│   └── auth.ts               # Authentication API endpoints
├── composables/
│   └── useAuthValidation.ts  # Form validation logic
└── types/
    └── index.ts              # TypeScript type definitions
```

## Features

### 1. User Registration

**Endpoint:** `POST /api/v1/auth/register`

**Features:**
- ✅ Complete form validation
- ✅ Password strength indicator
- ✅ Real-time validation feedback
- ✅ Terms and conditions acceptance
- ✅ Email format validation
- ✅ Name length validation (1-100 characters)
- ✅ Password confirmation
- ✅ Error handling for duplicate emails
- ✅ Loading states and disabled form during submission

**Form Fields:**
- First Name (required, 1-100 characters)
- Last Name (required, 1-100 characters)
- Email (required, valid email format)
- Password (required, minimum 8 characters)
- Confirm Password (required, must match password)
- Terms and Conditions (required checkbox)

**Password Strength Indicator:**
- Visual strength bar (red/yellow/green)
- Text indicator (Very Weak/Weak/Medium/Strong)
- Real-time strength calculation based on:
  - Length (8+ characters)
  - Character variety (uppercase, lowercase, numbers, symbols)

### 2. User Login

**Endpoint:** `POST /api/v1/auth/login`

**Features:**
- ✅ Email and password validation
- ✅ Error handling for invalid credentials
- ✅ Loading states
- ✅ Automatic redirect after successful login
- ✅ Session persistence

### 3. Session Management

**Features:**
- ✅ JWT token storage in localStorage
- ✅ Automatic token refresh
- ✅ Token expiration handling
- ✅ Persistent sessions across browser restarts
- ✅ Automatic logout on token expiration

## API Integration

### Registration Request

```typescript
interface RegisterRequest {
  email: string
  password: string
  first_name: string
  last_name: string
}
```

**Example:**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!",
  "first_name": "John",
  "last_name": "Doe"
}
```

### Registration Response

```typescript
interface AuthResponse {
  user: User
  tokens: AuthTokens
}
```

**Example:**
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "tier": "basic",
    "is_verified": false,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "tokens": {
    "access_token": "eyJ...",
    "refresh_token": "eyJ...",
    "expires_in": 3600
  }
}
```

## Validation System

### Form Validation Composable

The `useAuthValidation` composable provides:

- **Email Validation:** Format checking and required field validation
- **Password Validation:** Length and complexity requirements
- **Name Validation:** Length constraints (1-100 characters)
- **Password Strength:** Real-time strength calculation
- **API Error Handling:** Mapping server errors to form fields
- **Real-time Validation:** Instant feedback on field changes

### Validation Rules

| Field | Rules |
|-------|-------|
| Email | Required, valid email format |
| Password | Required, minimum 8 characters |
| First Name | Required, 1-100 characters |
| Last Name | Required, 1-100 characters |
| Confirm Password | Required, must match password |
| Terms | Required checkbox |

## Error Handling

### Client-side Validation Errors

- **Field-specific errors:** Displayed below each input field
- **General errors:** Displayed in a prominent error banner
- **Real-time validation:** Errors appear on field blur
- **Visual indicators:** Red borders on invalid fields

### API Error Handling

- **409 Conflict:** Email already exists
- **401 Unauthorized:** Invalid credentials
- **400 Bad Request:** Validation errors from server
- **500 Server Error:** Generic error handling

### Error Display

```vue
<!-- Field-specific error -->
<p v-if="errors.email" class="mt-1 text-sm text-red-600">
  {{ errors.email }}
</p>

<!-- General error banner -->
<div v-if="errors.general" class="rounded-md bg-red-50 p-4">
  <p class="text-sm text-red-800">{{ errors.general }}</p>
</div>
```

## State Management

### Auth Store (Pinia)

The authentication store manages:

- **User data:** Current user information
- **Tokens:** Access and refresh tokens
- **Loading states:** Form submission states
- **Authentication status:** Login/logout state
- **User tier:** Subscription level (guest/basic/premium)

### Store Actions

```typescript
// Login
await authStore.login(email, password)

// Register
await authStore.register(userData)

// Logout
authStore.logout()

// Refresh tokens
await authStore.refreshToken()

// Demo login (for testing)
authStore.demoLogin()
```

### Store Getters

```typescript
// Check if user is authenticated
const isAuthenticated = authStore.isAuthenticated

// Get user tier
const userTier = authStore.userTier

// Check feature access
const hasFeature = authStore.hasFeature('ai_insights')
```

## Security Features

### Token Management

- **Secure Storage:** Tokens stored in localStorage
- **Automatic Refresh:** Background token refresh
- **Expiration Handling:** Automatic logout on token expiry
- **Token Rotation:** Refresh token rotation for security

### Form Security

- **CSRF Protection:** Built into API client
- **Input Sanitization:** Automatic trimming and validation
- **Password Hiding:** Toggle password visibility
- **Rate Limiting:** Handled by backend API

## UI/UX Features

### Responsive Design

- **Mobile-first:** Optimized for all screen sizes
- **Accessibility:** Proper ARIA labels and keyboard navigation
- **Loading States:** Visual feedback during API calls
- **Error States:** Clear error messaging and recovery

### Visual Design

- **Consistent Styling:** Tailwind CSS with custom primary colors
- **Form Validation:** Real-time visual feedback
- **Password Strength:** Color-coded strength indicator
- **Loading Indicators:** Spinner animations during submission

### User Experience

- **Progressive Enhancement:** Works without JavaScript
- **Keyboard Navigation:** Full keyboard accessibility
- **Auto-focus:** Smart focus management
- **Form Persistence:** Form data preserved on validation errors

## Testing Features

### Demo Mode

For development and testing purposes:

- **Demo Login:** Instant login with premium user
- **Demo Registration:** Pre-filled form with demo data
- **UI Testing:** Test styling without backend

### Error Simulation

- **Network Errors:** Simulate API failures
- **Validation Errors:** Test form validation
- **Token Expiration:** Test session handling

## Integration Points

### Router Integration

```typescript
// Route guards for protected pages
const requireAuth = () => {
  const authStore = useAuthStore()
  if (!authStore.isAuthenticated) {
    return { name: 'login' }
  }
}

// Guest-only routes
meta: { requiresGuest: true }
```

### API Client Integration

```typescript
// Automatic token injection
apiClient.setAccessToken(tokens.access_token)

// Token expiration handling
window.addEventListener('token-expired', async () => {
  await authStore.refreshToken()
})
```

## Future Enhancements

### Planned Features

- [ ] Email verification flow
- [ ] Password reset functionality
- [ ] Two-factor authentication
- [ ] Social login integration
- [ ] Account deletion
- [ ] Profile management
- [ ] Session management UI

### Security Improvements

- [ ] Biometric authentication
- [ ] Device fingerprinting
- [ ] Suspicious activity detection
- [ ] Enhanced password policies
- [ ] Audit logging

## Usage Examples

### Basic Registration

```vue
<template>
  <RegisterView />
</template>

<script setup>
import RegisterView from '@/views/auth/RegisterView.vue'
</script>
```

### Custom Validation

```vue
<script setup>
import { useAuthValidation } from '@/composables/useAuthValidation'

const { validateRegistrationForm, errors } = useAuthValidation()

const handleSubmit = () => {
  if (validateRegistrationForm(formData)) {
    // Proceed with registration
  }
}
</script>
```

### Authentication Check

```vue
<script setup>
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

// Check if user can access premium features
if (authStore.hasFeature('ai_insights')) {
  // Show premium content
}
</script>
```

## Troubleshooting

### Common Issues

1. **Token Expiration:** Automatic refresh should handle this
2. **Validation Errors:** Check form field requirements
3. **Network Errors:** Verify API endpoint availability
4. **Duplicate Email:** Use different email for registration

### Debug Mode

Enable debug logging in development:

```typescript
// In auth store
if (import.meta.env.DEV) {
  console.log('Auth state:', { user, isAuthenticated })
}
```

## Performance Considerations

- **Lazy Loading:** Auth components loaded on demand
- **Token Caching:** Efficient token storage and retrieval
- **Form Debouncing:** Prevents excessive validation calls
- **Bundle Splitting:** Auth code separated from main bundle

This authentication system provides a robust, secure, and user-friendly foundation for the Stock Tracker application, with comprehensive validation, error handling, and modern UI/UX patterns.