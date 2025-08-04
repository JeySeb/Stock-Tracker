<template>
  <div class="min-h-screen bg-white">
    <!-- Hero Section -->
    <div class="relative bg-white overflow-hidden">
      <div class="max-w-7xl mx-auto">
        <div class="relative z-10 pb-8 bg-white sm:pb-16 md:pb-20 lg:max-w-2xl lg:w-full lg:pb-28 xl:pb-32">
          <main class="mt-10 mx-auto max-w-7xl px-4 sm:mt-12 sm:px-6 md:mt-16 lg:mt-20 lg:px-8 xl:mt-28">
            <div class="sm:text-center lg:text-left">
              <h1 class="text-4xl tracking-tight font-extrabold text-gray-900 sm:text-5xl md:text-6xl">
                <span class="block xl:inline">AI-Powered</span>
                <span class="block text-primary-600 xl:inline">Stock Analysis</span>
              </h1>
              <p class="mt-3 text-base text-gray-500 sm:mt-5 sm:text-lg sm:max-w-xl sm:mx-auto md:mt-5 md:text-xl lg:mx-0">
                Get intelligent stock recommendations powered by AI and real-time market data. 
                Make informed investment decisions with our advanced analytics platform.
              </p>
              <div class="mt-5 sm:mt-8 sm:flex sm:justify-center lg:justify-start space-y-3 sm:space-y-0 sm:space-x-3">
                <div class="rounded-md shadow">
                  <router-link
                    to="/dashboard"
                    class="w-full flex items-center justify-center px-8 py-3 border border-transparent text-base font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 md:py-4 md:text-lg md:px-10 transition-colors"
                  >
                    Explore Now
                  </router-link>
                </div>
                <div class="rounded-md shadow">
                  <router-link
                    to="/register"
                    class="w-full flex items-center justify-center px-8 py-3 border border-transparent text-base font-medium rounded-md text-primary-700 bg-primary-100 hover:bg-primary-200 md:py-4 md:text-lg md:px-10 transition-colors"
                  >
                    Start Free Trial
                  </router-link>
                </div>
                <div>
                  <button
                    @click="scrollToPricing"
                    class="w-full flex items-center justify-center px-8 py-3 border border-gray-300 text-base font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 md:py-4 md:text-lg md:px-10 transition-colors"
                  >
                    View Pricing
                  </button>
                </div>
              </div>
              
              <!-- Guest Access Prominent Call-to-Action
              <div class="mt-8 p-4 bg-gradient-to-r from-primary-50 to-blue-50 rounded-lg border border-primary-200">
                <div class="flex items-center space-x-3">
                  <div class="flex-shrink-0">
                    <div class="w-8 h-8 bg-primary-600 rounded-full flex items-center justify-center">
                      <span class="text-white text-sm">🚀</span>
                    </div>
                  </div>
                  <div class="flex-1">
                    <p class="text-sm text-primary-800">
                      <strong>New to Stock Tracker?</strong> Explore our dashboard with sample data - no registration required!
                    </p>
                  </div>
                  <router-link
                    to="/dashboard"
                    class="flex-shrink-0 bg-primary-600 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-primary-700 transition-colors"
                  >
                    Explore Now
                  </router-link>
                </div>
              </div> -->
            </div>
          </main>
        </div>
      </div>
      <div class="lg:absolute lg:inset-y-0 lg:right-0 lg:w-1/2">
        <div class="h-56 w-full bg-gradient-to-br from-primary-400 to-primary-600 sm:h-72 md:h-96 lg:w-full lg:h-full flex items-center justify-center">
          <div class="text-center text-white">
            <AnimatedStockIcon />
            <div class="text-xl font-semibold mt-4">Real-time Analytics</div>
            <div class="text-sm opacity-90 mt-2">Live market insights at your fingertips</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Features Section -->
    <div class="py-12 bg-gray-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="lg:text-center">
          <h2 class="text-base text-primary-600 font-semibold tracking-wide uppercase">Features</h2>
          <p class="mt-2 text-3xl leading-8 font-extrabold tracking-tight text-gray-900 sm:text-4xl">
            Everything you need to make smart investment decisions
          </p>
          <p class="mt-4 max-w-2xl text-xl text-gray-500 lg:mx-auto">
            Our platform provides comprehensive tools and insights to help you navigate the stock market with confidence.
          </p>
        </div>

        <div class="mt-10">
          <dl class="space-y-10 md:space-y-0 md:grid md:grid-cols-2 md:gap-x-8 md:gap-y-10">
            <div v-for="feature in features" :key="feature.name" class="relative">
              <dt>
                <div class="absolute flex items-center justify-center h-12 w-12 rounded-md bg-primary-500 text-white">
                  <span class="text-2xl">{{ feature.icon }}</span>
                </div>
                <p class="ml-16 text-lg leading-6 font-medium text-gray-900">{{ feature.name }}</p>
              </dt>
              <dd class="mt-2 ml-16 text-base text-gray-500">{{ feature.description }}</dd>
            </div>
          </dl>
        </div>
      </div>
    </div>

    <!-- Tier Benefits Section -->
    <div class="bg-white py-16">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center">
          <h2 class="text-3xl font-extrabold text-gray-900 sm:text-4xl">
            Choose Your Access Level
          </h2>
          <p class="mt-3 max-w-2xl mx-auto text-xl text-gray-500 sm:mt-4">
            Start with our free features and upgrade when you're ready for more advanced capabilities.
          </p>
        </div>

        <div class="mt-12 space-y-4 sm:mt-16 sm:space-y-0 sm:grid sm:grid-cols-3 sm:gap-6">
          <!-- Guest Tier -->
          <div class="border border-gray-200 rounded-lg shadow-sm divide-y divide-gray-200">
            <div class="p-6">
              <h3 class="text-lg leading-6 font-medium text-gray-900">Guest Access</h3>
              <p class="mt-4 text-sm text-gray-500">Perfect for exploring our platform</p>
              <p class="mt-8">
                <span class="text-4xl font-extrabold text-gray-900">Free</span>
              </p>
              <router-link
                to="/dashboard"
                class="mt-8 block w-full bg-gray-600 border border-transparent rounded-md py-2 text-sm font-semibold text-white text-center hover:bg-gray-700"
              >
                Try Now
              </router-link>
            </div>
            <div class="pt-6 pb-8 px-6">
              <h4 class="text-xs font-medium text-gray-900 tracking-wide uppercase">Included features</h4>
              <ul class="mt-6 space-y-4">
                <li v-for="feature in guestFeatures" :key="feature" class="flex space-x-3">
                  <svg class="flex-shrink-0 h-5 w-5 text-green-500" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                  </svg>
                  <span class="text-sm text-gray-500">{{ feature }}</span>
                </li>
              </ul>
            </div>
          </div>

          <!-- Registered Tier -->
          <div class="border border-gray-200 rounded-lg shadow-sm divide-y divide-gray-200">
            <div class="p-6">
              <h3 class="text-lg leading-6 font-medium text-gray-900">Registered User</h3>
              <p class="mt-4 text-sm text-gray-500">Enhanced access with external data</p>
              <p class="mt-8">
                <span class="text-4xl font-extrabold text-gray-900">Free</span>
              </p>
              <router-link
                to="/register"
                class="mt-8 block w-full bg-primary-600 border border-transparent rounded-md py-2 text-sm font-semibold text-white text-center hover:bg-primary-700"
              >
                Sign Up Free
              </router-link>
            </div>
            <div class="pt-6 pb-8 px-6">
              <h4 class="text-xs font-medium text-gray-900 tracking-wide uppercase">Everything in Guest, plus</h4>
              <ul class="mt-6 space-y-4">
                <li v-for="feature in registeredFeatures" :key="feature" class="flex space-x-3">
                  <svg class="flex-shrink-0 h-5 w-5 text-green-500" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                  </svg>
                  <span class="text-sm text-gray-500">{{ feature }}</span>
                </li>
              </ul>
            </div>
          </div>

          <!-- Premium Tier -->
          <div class="border border-primary-200 rounded-lg shadow-sm divide-y divide-gray-200 relative">
            <div class="absolute top-0 right-6 transform -translate-y-1/2">
              <span class="inline-flex rounded-full bg-primary-100 px-4 py-1 text-sm font-semibold text-primary-600">
                Most Popular
              </span>
            </div>
            <div class="p-6">
              <h3 class="text-lg leading-6 font-medium text-gray-900">Premium</h3>
              <p class="mt-4 text-sm text-gray-500">AI-powered insights and advanced features</p>
              <p class="mt-8">
                <span class="text-4xl font-extrabold text-gray-900">$29.99</span>
                <span class="text-base font-medium text-gray-500">/month</span>
              </p>
              <router-link
                to="/subscription"
                class="mt-8 block w-full bg-primary-600 border border-transparent rounded-md py-2 text-sm font-semibold text-white text-center hover:bg-primary-700"
              >
                Upgrade to Premium
              </router-link>
            </div>
            <div class="pt-6 pb-8 px-6">
              <h4 class="text-xs font-medium text-gray-900 tracking-wide uppercase">Everything in Registered, plus</h4>
              <ul class="mt-6 space-y-4">
                <li v-for="feature in premiumFeatures" :key="feature" class="flex space-x-3">
                  <svg class="flex-shrink-0 h-5 w-5 text-green-500" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                  </svg>
                  <span class="text-sm text-gray-500">{{ feature }}</span>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Pricing Section -->
    <div id="pricing" class="bg-gray-50 py-16">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center">
          <h2 class="text-3xl font-extrabold text-gray-900 sm:text-4xl">
            Premium Subscription Plans
          </h2>
          <p class="mt-3 max-w-2xl mx-auto text-xl text-gray-500 sm:mt-4">
            Choose the plan that best fits your investment strategy.
          </p>
        </div>

        <div class="mt-12 space-y-4 sm:mt-16 sm:space-y-0 sm:grid sm:grid-cols-2 sm:gap-6 lg:max-w-4xl lg:mx-auto">
          <!-- Monthly Plan -->
          <div class="border border-gray-200 rounded-lg shadow-sm divide-y divide-gray-200 bg-white">
            <div class="p-6">
              <h3 class="text-lg leading-6 font-medium text-gray-900">Monthly Plan</h3>
              <p class="mt-4 text-sm text-gray-500">Perfect for trying premium features</p>
              <p class="mt-8">
                <span class="text-4xl font-extrabold text-gray-900">$29.99</span>
                <span class="text-base font-medium text-gray-500">/month</span>
              </p>
              <router-link
                to="/subscription"
                class="mt-8 block w-full bg-primary-600 border border-transparent rounded-md py-2 text-sm font-semibold text-white text-center hover:bg-primary-700"
              >
                Start Monthly Plan
              </router-link>
            </div>
            <div class="pt-6 pb-8 px-6">
              <h4 class="text-xs font-medium text-gray-900 tracking-wide uppercase">What's included</h4>
              <ul class="mt-6 space-y-4">
                <li v-for="feature in monthlyPlanFeatures" :key="feature" class="flex space-x-3">
                  <svg class="flex-shrink-0 h-5 w-5 text-green-500" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                  </svg>
                  <span class="text-sm text-gray-500">{{ feature }}</span>
                </li>
              </ul>
            </div>
          </div>

          <!-- Yearly Plan -->
          <div class="border border-primary-200 rounded-lg shadow-sm divide-y divide-gray-200 bg-white relative">
            <div class="absolute top-0 right-6 transform -translate-y-1/2">
              <span class="inline-flex rounded-full bg-green-100 px-4 py-1 text-sm font-semibold text-green-600">
                Save $59.89
              </span>
            </div>
            <div class="p-6">
              <h3 class="text-lg leading-6 font-medium text-gray-900">Yearly Plan</h3>
              <p class="mt-4 text-sm text-gray-500">Best value for committed investors</p>
              <p class="mt-8">
                <span class="text-4xl font-extrabold text-gray-900">$299.99</span>
                <span class="text-base font-medium text-gray-500">/year</span>
              </p>
              <p class="text-sm text-green-600 mt-1">$24.99/month when paid annually</p>
              <router-link
                to="/subscription"
                class="mt-8 block w-full bg-primary-600 border border-transparent rounded-md py-2 text-sm font-semibold text-white text-center hover:bg-primary-700"
              >
                Start Yearly Plan
              </router-link>
            </div>
            <div class="pt-6 pb-8 px-6">
              <h4 class="text-xs font-medium text-gray-900 tracking-wide uppercase">Everything in Monthly, plus</h4>
              <ul class="mt-6 space-y-4">
                <li v-for="feature in yearlyPlanFeatures" :key="feature" class="flex space-x-3">
                  <svg class="flex-shrink-0 h-5 w-5 text-green-500" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                  </svg>
                  <span class="text-sm text-gray-500">{{ feature }}</span>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- CTA Section -->
    <div class="bg-primary-600">
      <div class="max-w-2xl mx-auto text-center py-16 px-4 sm:py-20 sm:px-6 lg:px-8">
        <h2 class="text-3xl font-extrabold text-white sm:text-4xl">
          <span class="block">Ready to start investing smarter?</span>
        </h2>
        <p class="mt-4 text-lg leading-6 text-primary-200">
          Join thousands of investors who trust our AI-powered platform for their investment decisions.
        </p>
        <router-link
          to="/register"
          class="mt-8 w-full inline-flex items-center justify-center px-5 py-3 border border-transparent text-base font-medium rounded-md text-primary-600 bg-white hover:bg-primary-50 sm:w-auto"
        >
          Start Your Free Account
        </router-link>
      </div>
    </div>

    <!-- Footer -->
    <footer class="bg-white">
      <div class="max-w-7xl mx-auto py-12 px-4 sm:px-6 lg:px-8">
        <div class="text-center">
          <h3 class="text-lg font-semibold text-gray-900">Stock Tracker</h3>
          <p class="mt-2 text-sm text-gray-500">
            AI-powered investment insights for smarter decisions.
          </p>
          <div class="mt-4 flex justify-center space-x-6">
            <router-link to="/login" class="text-gray-400 hover:text-gray-500">
              Sign In
            </router-link>
            <router-link to="/register" class="text-gray-400 hover:text-gray-500">
              Sign Up
            </router-link>
            <router-link to="/subscription" class="text-gray-400 hover:text-gray-500">
              Pricing
            </router-link>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import AnimatedStockIcon from '@/components/icons/AnimatedStockIcon.vue'

const features = [
  {
    name: 'Real-time Market Data',
    description: 'Access live stock prices, market trends, and financial indicators from multiple sources.',
    icon: '📊'
  },
  {
    name: 'AI-Powered Recommendations',
    description: 'Get intelligent stock recommendations based on advanced machine learning algorithms.',
    icon: '🤖'
  },
  {
    name: 'Sentiment Analysis',
    description: 'Track news and social media sentiment to gauge market emotions around stocks.',
    icon: '💭'
  },
  {
    name: 'Risk Assessment',
    description: 'Understand investment risks with comprehensive analysis and scoring.',
    icon: '⚖️'
  },
  {
    name: 'Portfolio Analytics',
    description: 'Advanced analytics to track and optimize your investment portfolio performance.',
    icon: '📈'
  },
  {
    name: 'External API Integration',
    description: 'Connect with Yahoo Finance, Alpha Vantage, and other financial data providers.',
    icon: '🔗'
  }
]

const guestFeatures = [
  'Basic stock recommendations',
  'Market analytics system',
  'Filter and browse stock events',
  'Basic recommendation algorithm'
]

const registeredFeatures = [
  'Real-time market data',
  'External API access',
  'Yahoo Finance integration',
  'Alpha Vantage data',
  'Enhanced analytics'
]

const premiumFeatures = [
  'AI-powered insights',
  'Advanced sentiment analysis',
  'AI chatbot assistant',
  'Custom alerts',
  'Priority support'
]

const monthlyPlanFeatures = [
  'All registered features',
  'Real-time data access',
  'External API integration',
  'Advanced analytics',
  'Email support'
]

const yearlyPlanFeatures = [
  'Save 17% vs monthly',
  'AI-powered insights',
  'Sentiment analysis',
  'Priority support',
  'Advanced AI features'
]

function scrollToPricing() {
  const element = document.getElementById('pricing')
  if (element) {
    element.scrollIntoView({ behavior: 'smooth' })
  }
}
</script> 