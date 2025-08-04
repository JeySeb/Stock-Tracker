<template>
  <div class="relative w-full h-full flex items-center justify-center overflow-hidden">
    <!-- Professional Chart Container with Glass Effect -->
    <div class="relative w-full h-full flex items-center justify-center overflow-hidden rounded-2xl chart-container p-8 mx-auto" style="max-width: 900px; height: 600px;">
      <!-- Main Chart Container -->
      <div class="relative w-72 h-72 lg:w-80 lg:h-80 xl:w-96 xl:h-96">
        <!-- Background Grid -->
        <svg class="absolute inset-0 w-full h-full" viewBox="0 0 224 224" fill="none">
          <defs>
            <pattern id="chartGrid" width="28" height="28" patternUnits="userSpaceOnUse">
              <path d="M 28 0 L 0 0 0 28" fill="none" stroke="currentColor" stroke-width="0.5" opacity="0.1"/>
            </pattern>
            <!-- Enhanced Gradients for chart areas -->
            <linearGradient id="chartGradient1" x1="0%" y1="0%" x2="0%" y2="100%">
              <stop offset="0%" stop-color="rgba(59, 130, 246, 0.2)"/>
              <stop offset="50%" stop-color="rgba(59, 130, 246, 0.1)"/>
              <stop offset="100%" stop-color="rgba(59, 130, 246, 0.02)"/>
            </linearGradient>
            <linearGradient id="chartGradient2" x1="0%" y1="0%" x2="0%" y2="100%">
              <stop offset="0%" stop-color="rgba(34, 197, 94, 0.2)"/>
              <stop offset="50%" stop-color="rgba(34, 197, 94, 0.1)"/>
              <stop offset="100%" stop-color="rgba(34, 197, 94, 0.02)"/>
            </linearGradient>
            <linearGradient id="chartGradient3" x1="0%" y1="0%" x2="0%" y2="100%">
              <stop offset="0%" stop-color="rgba(168, 85, 247, 0.2)"/>
              <stop offset="50%" stop-color="rgba(168, 85, 247, 0.1)"/>
              <stop offset="100%" stop-color="rgba(168, 85, 247, 0.02)"/>
            </linearGradient>
          </defs>
          <rect width="224" height="224" fill="url(#chartGrid)" class="grid-line"/>
        </svg>

        <!-- Stock Price Chart (Main Trend) -->
        <svg class="absolute inset-0 w-full h-full" viewBox="0 0 224 224" fill="none">
          <!-- Main Stock Price Line -->
          <path 
            d="M 14 168 Q 42 140 70 150 T 126 112 T 182 84 T 210 56" 
            stroke="rgb(59, 130, 246)" 
            stroke-width="4" 
            fill="none"
            class="animate-draw-primary"
            stroke-linecap="round"
          />
          <path 
            d="M 14 168 Q 42 140 70 150 T 126 112 T 182 84 T 210 56 L 210 224 L 14 224 Z" 
            fill="url(#chartGradient1)"
            class="animate-fill-area"
          />
        </svg>

        <!-- Volume Chart (Secondary) -->
        <svg class="absolute inset-0 w-full h-full" viewBox="0 0 224 224" fill="none">
          <path 
            d="M 21 182 Q 49 154 77 164 T 133 126 T 189 98 T 217 70" 
            stroke="rgb(34, 197, 94)" 
            stroke-width="3" 
            fill="none"
            class="animate-draw-secondary"
            stroke-linecap="round"
          />
          <path 
            d="M 21 182 Q 49 154 77 164 T 133 126 T 189 98 T 217 70 L 217 224 L 21 224 Z" 
            fill="url(#chartGradient2)"
            class="animate-fill-area-delayed"
          />
        </svg>

        <!-- Market Sentiment (Tertiary) -->
        <svg class="absolute inset-0 w-full h-full" viewBox="0 0 224 224" fill="none">
          <path 
            d="M 28 196 Q 56 168 84 178 T 140 140 T 196 112 T 224 84" 
            stroke="rgb(168, 85, 247)" 
            stroke-width="3" 
            fill="none"
            class="animate-draw-tertiary"
            stroke-linecap="round"
          />
          <path 
            d="M 28 196 Q 56 168 84 178 T 140 140 T 196 112 T 224 84 L 224 224 L 28 224 Z" 
            fill="url(#chartGradient3)"
            class="animate-fill-area-delayed-2"
          />
        </svg>

        <!-- Candlestick Chart Elements -->
        <div class="absolute inset-0">
          <!-- Bullish Candlesticks (Green) -->
          <div class="absolute w-2 h-8 candle-bull rounded-sm animate-candle" style="left: 15%; top: 60%; animation-delay: 0.5s;"></div>
          <div class="absolute w-2 h-6 candle-bull rounded-sm animate-candle" style="left: 25%; top: 65%; animation-delay: 0.7s;"></div>
          <div class="absolute w-2 h-10 candle-bull rounded-sm animate-candle" style="left: 35%; top: 55%; animation-delay: 0.9s;"></div>
          <div class="absolute w-2 h-8 candle-bull rounded-sm animate-candle" style="left: 45%; top: 60%; animation-delay: 1.1s;"></div>
          <div class="absolute w-2 h-12 candle-bull rounded-sm animate-candle" style="left: 55%; top: 50%; animation-delay: 1.3s;"></div>
          
          <!-- Bearish Candlesticks (Red) -->
          <div class="absolute w-2 h-6 candle-bear rounded-sm animate-candle" style="left: 65%; top: 70%; animation-delay: 1.5s;"></div>
          <div class="absolute w-2 h-8 candle-bear rounded-sm animate-candle" style="left: 75%; top: 65%; animation-delay: 1.7s;"></div>
          <div class="absolute w-2 h-4 candle-bear rounded-sm animate-candle" style="left: 85%; top: 75%; animation-delay: 1.9s;"></div>
        </div>

        <!-- Enhanced Data Points with Stock-Specific Styling -->
        <div class="absolute inset-0">
          <div class="absolute w-5 h-5 bg-blue-500 rounded-full animate-ping shadow-lg" style="left: 6.25%; top: 75%; animation-delay: 0.5s; box-shadow: 0 0 8px rgba(59, 130, 246, 0.3);"></div>
          <div class="absolute w-5 h-5 bg-blue-500 rounded-full animate-ping shadow-lg" style="left: 31.25%; top: 62.5%; animation-delay: 1s; box-shadow: 0 0 8px rgba(59, 130, 246, 0.3);"></div>
          <div class="absolute w-5 h-5 bg-blue-500 rounded-full animate-ping shadow-lg" style="left: 56.25%; top: 50%; animation-delay: 1.5s; box-shadow: 0 0 8px rgba(59, 130, 246, 0.3);"></div>
          <div class="absolute w-5 h-5 bg-blue-500 rounded-full animate-ping shadow-lg" style="left: 81.25%; top: 37.5%; animation-delay: 2s; box-shadow: 0 0 8px rgba(59, 130, 246, 0.3);"></div>
          <div class="absolute w-5 h-5 bg-blue-500 rounded-full animate-ping shadow-lg" style="left: 93.75%; top: 25%; animation-delay: 2.5s; box-shadow: 0 0 8px rgba(59, 130, 246, 0.3);"></div>
        </div>
      </div>

      <!-- Professional Financial Indicators -->
      <div class="absolute top-8 right-8 financial-indicator rounded-xl px-4 py-3 animate-fade-in">
        <div class="flex items-center">
          <div class="bg-blue-100 p-2 rounded-lg mr-3">
            <i class="fas fa-chart-line text-blue-600 text-lg"></i>
          </div>
          <div>
            <p class="text-gray-500 text-sm"></p>
            <p class="text-xl font-bold text-gray-800">$1,234.56</p>
            <p class="text-green-500 text-sm font-medium flex items-center">
              <i class="fas fa-arrow-up mr-1"></i> +2.5% Today
            </p>
          </div>
        </div>
      </div>
      
      <div class="absolute bottom-8 left-8 financial-indicator rounded-xl px-4 py-3 animate-fade-in" style="animation-delay: 0.3s;">
        <div class="flex items-center">
          <div class="bg-green-100 p-2 rounded-lg mr-3">
            <i class="fas fa-arrow-trend-up text-green-600 text-lg"></i>
          </div>
          <div>
            <p class="text-gray-500 text-sm">52W High</p>
            <p class="text-xl font-bold text-gray-800">$1,378.90</p>
            <p class="text-green-500 text-sm font-medium">+11.6% from current</p>
          </div>
        </div>
      </div>
      
      <div class="absolute top-1/2 -right-12 w-16 h-16 bg-white rounded-full flex items-center justify-center animate-bounce shadow-xl financial-indicator" style="animation-delay: 1s;">
        <div class="w-12 h-12 bg-purple-100 rounded-full flex items-center justify-center">
          <i class="fas fa-chart-pie text-purple-600 text-xl"></i>
        </div>
      </div>
      
      <div class="absolute top-1/4 -left-12 w-16 h-16 bg-white rounded-full flex items-center justify-center animate-bounce shadow-xl financial-indicator" style="animation-delay: 1.5s;">
        <div class="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
          <i class="fas fa-bullhorn text-blue-600 text-xl"></i>
        </div>
      </div>

      <!-- Enhanced Pulse Rings -->
      <div class="absolute inset-0 rounded-full border-3 border-blue-200 animate-ping opacity-75 pulse-ring"></div>
      <div class="absolute inset-0 rounded-full border-3 border-green-200 animate-ping opacity-50 pulse-ring" style="animation-delay: 1s;"></div>
      <div class="absolute inset-0 rounded-full border-3 border-purple-200 animate-ping opacity-30 pulse-ring" style="animation-delay: 2s;"></div>

      <!-- Stock Market Data Particles -->
      <div class="absolute inset-0">
        <div class="absolute w-3 h-3 bg-blue-400 rounded-full animate-float shadow-lg" style="left: 20%; top: 30%; animation-delay: 0s; box-shadow: 0 0 8px rgba(59, 130, 246, 0.3);"></div>
        <div class="absolute w-2 h-2 bg-green-400 rounded-full animate-float shadow-lg" style="left: 40%; top: 20%; animation-delay: 0.5s; box-shadow: 0 0 6px rgba(34, 197, 94, 0.3);"></div>
        <div class="absolute w-2.5 h-2.5 bg-purple-400 rounded-full animate-float shadow-lg" style="left: 60%; top: 40%; animation-delay: 1s; box-shadow: 0 0 7px rgba(168, 85, 247, 0.3);"></div>
        <div class="absolute w-2 h-2 bg-blue-400 rounded-full animate-float shadow-lg" style="left: 80%; top: 25%; animation-delay: 1.5s; box-shadow: 0 0 6px rgba(59, 130, 246, 0.3);"></div>
        <div class="absolute w-3 h-3 bg-green-400 rounded-full animate-float shadow-lg" style="left: 30%; top: 60%; animation-delay: 2s; box-shadow: 0 0 8px rgba(34, 197, 94, 0.3);"></div>
      </div>
      
      <!-- Enhanced Loading Animation -->
      <div class="absolute inset-0 flex items-center justify-center">
        <div class="w-20 h-20 border-4 border-blue-100 border-t-blue-500 rounded-full animate-spin opacity-50"></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// Component logic can be added here if needed
</script>

<style scoped>
@keyframes drawPrimary {
  0% {
    stroke-dasharray: 0 300;
    stroke-dashoffset: 0;
    opacity: 0;
  }
  25% {
    opacity: 1;
  }
  50% {
    stroke-dasharray: 300 0;
    stroke-dashoffset: 0;
    opacity: 1;
  }
  100% {
    stroke-dasharray: 300 0;
    stroke-dashoffset: 0;
    opacity: 1;
  }
}

@keyframes drawSecondary {
  0% {
    stroke-dasharray: 0 300;
    stroke-dashoffset: 0;
    opacity: 0;
  }
  25% {
    opacity: 1;
  }
  50% {
    stroke-dasharray: 300 0;
    stroke-dashoffset: 0;
    opacity: 1;
  }
  100% {
    stroke-dasharray: 300 0;
    stroke-dashoffset: 0;
    opacity: 1;
  }
}

@keyframes drawTertiary {
  0% {
    stroke-dasharray: 0 300;
    stroke-dashoffset: 0;
    opacity: 0;
  }
  25% {
    opacity: 1;
  }
  50% {
    stroke-dasharray: 300 0;
    stroke-dashoffset: 0;
    opacity: 1;
  }
  100% {
    stroke-dasharray: 300 0;
    stroke-dashoffset: 0;
    opacity: 1;
  }
}

@keyframes fillArea {
  0% {
    opacity: 0;
    transform: scaleY(0);
  }
  50% {
    opacity: 0.4;
    transform: scaleY(0.5);
  }
  100% {
    opacity: 0.6;
    transform: scaleY(1);
  }
}

@keyframes fillAreaDelayed {
  0% {
    opacity: 0;
    transform: scaleY(0);
  }
  50% {
    opacity: 0.3;
    transform: scaleY(0.5);
  }
  100% {
    opacity: 0.4;
    transform: scaleY(1);
  }
}

@keyframes fillAreaDelayed2 {
  0% {
    opacity: 0;
    transform: scaleY(0);
  }
  50% {
    opacity: 0.2;
    transform: scaleY(0.5);
  }
  100% {
    opacity: 0.3;
    transform: scaleY(1);
  }
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px) scale(1) rotate(0deg);
    opacity: 0.7;
  }
  25% {
    transform: translateY(-8px) scale(1.1) rotate(5deg);
    opacity: 0.9;
  }
  50% {
    transform: translateY(-15px) scale(1.2) rotate(0deg);
    opacity: 1;
  }
  75% {
    transform: translateY(-8px) scale(1.1) rotate(-5deg);
    opacity: 0.9;
  }
}

@keyframes candle {
  0%, 100% {
    opacity: 0.6;
    transform: scaleY(1);
    box-shadow: 0 0 4px rgba(0,0,0,0.1);
  }
  50% {
    opacity: 1;
    transform: scaleY(1.1);
    box-shadow: 0 0 8px rgba(0,0,0,0.15);
  }
}

@keyframes fadeIn {
  0% {
    opacity: 0;
    transform: translateY(-10px);
  }
  100% {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes ping {
  75%, 100% {
    transform: scale(2);
    opacity: 0;
  }
}

@keyframes bounce {
  0%, 20%, 53%, 80%, 100% {
    transform: translate3d(0, 0, 0);
  }
  40%, 43% {
    transform: translate3d(0, -8px, 0);
  }
  70% {
    transform: translate3d(0, -4px, 0);
  }
  90% {
    transform: translate3d(0, -2px, 0);
  }
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: .5;
  }
}

.animate-draw-primary {
  animation: drawPrimary 5s ease-in-out infinite;
  stroke-dasharray: 300;
  stroke-dashoffset: 300;
}

.animate-draw-secondary {
  animation: drawSecondary 5s ease-in-out infinite;
  stroke-dasharray: 300;
  stroke-dashoffset: 300;
  animation-delay: 0.8s;
}

.animate-draw-tertiary {
  animation: drawTertiary 5s ease-in-out infinite;
  stroke-dasharray: 300;
  stroke-dashoffset: 300;
  animation-delay: 1.6s;
}

.animate-fill-area {
  animation: fillArea 5s ease-in-out infinite;
  transform-origin: bottom;
}

.animate-fill-area-delayed {
  animation: fillAreaDelayed 5s ease-in-out infinite;
  animation-delay: 0.8s;
  transform-origin: bottom;
}

.animate-fill-area-delayed-2 {
  animation: fillAreaDelayed2 5s ease-in-out infinite;
  animation-delay: 1.6s;
  transform-origin: bottom;
}

.animate-float {
  animation: float 5s ease-in-out infinite;
}

.animate-candle {
  animation: candle 3s ease-in-out infinite;
}

.animate-fade-in {
  animation: fadeIn 2s ease-out forwards;
}

.animate-ping {
  animation: ping 2s cubic-bezier(0, 0, 0.2, 1) infinite;
}

.animate-bounce {
  animation: bounce 2s infinite;
}

.animate-pulse {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

.border-3 {
  border-width: 3px;
}

.chart-container {
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(0, 0, 0, 0.05);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
  backdrop-filter: blur(10px);
}

.financial-indicator {
  background: white;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.07);
  border: 1px solid rgba(0, 0, 0, 0.03);
}

.grid-line {
  stroke: rgba(0, 0, 0, 0.05);
}

.candle-bull {
  background: linear-gradient(to bottom, #10b981, #059669);
}

.candle-bear {
  background: linear-gradient(to bottom, #ef4444, #dc2626);
}

.pulse-ring {
  border-color: rgba(59, 130, 246, 0.15);
}
</style>