<template>
  <div class="relative">
    <!-- Label -->
    <label v-if="label" :for="id" class="block text-sm font-medium text-gray-700 mb-1">
      {{ label }}
    </label>

    <!-- Selected Items Display -->
    <div
      :id="id"
      @click="toggleDropdown"
      class="relative w-full min-h-[38px] px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm cursor-pointer focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
      :class="{ 'border-primary-500 ring-2 ring-primary-500': isOpen }"
    >
      <div class="flex flex-wrap gap-1">
        <template v-if="selectedItems.length">
          <span
            v-for="item in selectedItems"
            :key="item.value"
            class="inline-flex items-center px-2 py-0.5 rounded text-sm font-medium bg-primary-100 text-primary-800"
          >
            {{ item.label }}
            <button
              type="button"
              @click.stop="removeItem(item)"
              class="ml-1 inline-flex items-center justify-center w-4 h-4 text-primary-400 hover:bg-primary-200 hover:text-primary-500 rounded-full focus:outline-none"
            >
              <span class="sr-only">Remove</span>
              ×
            </button>
          </span>
        </template>
        <span v-else class="text-gray-500">{{ placeholder }}</span>
      </div>

      <!-- Dropdown arrow -->
      <span class="absolute inset-y-0 right-0 flex items-center pr-2 pointer-events-none">
        <svg class="h-5 w-5 text-gray-400" viewBox="0 0 20 20" fill="currentColor">
          <path
            fill-rule="evenodd"
            d="M10 3a1 1 0 01.707.293l3 3a1 1 0 01-1.414 1.414L10 5.414 7.707 7.707a1 1 0 01-1.414-1.414l3-3A1 1 0 0110 3zm-3.707 9.293a1 1 0 011.414 0L10 14.586l2.293-2.293a1 1 0 011.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z"
            clip-rule="evenodd"
          />
        </svg>
      </span>
    </div>

    <!-- Dropdown Menu -->
    <div
      v-show="isOpen"
      class="absolute z-10 mt-1 w-full bg-white shadow-lg max-h-60 rounded-md py-1 text-base ring-1 ring-black ring-opacity-5 overflow-auto focus:outline-none sm:text-sm"
    >
      <!-- Search Input -->
      <div class="sticky top-0 bg-white px-3 py-2">
        <input
          type="search"
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          placeholder="Search..."
          v-model="searchQuery"
          @click.stop
        />
      </div>

      <!-- Options List -->
      <div class="py-1">
        <template v-if="filteredOptions.length">
          <button
            v-for="option in filteredOptions"
            :key="option.value"
            @click.stop="toggleItem(option)"
            class="w-full text-left px-3 py-2 hover:bg-primary-50 flex items-center justify-between"
            :class="{ 'bg-primary-50': isSelected(option) }"
          >
            <span>{{ option.label }}</span>
            <svg
              v-if="isSelected(option)"
              class="h-5 w-5 text-primary-600"
              viewBox="0 0 20 20"
              fill="currentColor"
            >
              <path
                fill-rule="evenodd"
                d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                clip-rule="evenodd"
              />
            </svg>
          </button>
        </template>
        <div v-else class="px-3 py-2 text-sm text-gray-500">
          No options found
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'

interface Option {
  label: string
  value: string
}

interface Props {
  modelValue: string[]
  options: Option[]
  label?: string
  placeholder?: string
  id?: string
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: 'Select options...',
  id: () => `multiselect-${Math.random().toString(36).substr(2, 9)}`
})

const emit = defineEmits<{
  'update:modelValue': [value: string[]]
}>()

const isOpen = ref(false)
const searchQuery = ref('')

const selectedItems = computed(() => {
  return props.options.filter(option => props.modelValue.includes(option.value))
})

const filteredOptions = computed(() => {
  if (!searchQuery.value) return props.options
  const query = searchQuery.value.toLowerCase()
  return props.options.filter(option => 
    option.label.toLowerCase().includes(query) || 
    option.value.toLowerCase().includes(query)
  )
})

function toggleDropdown() {
  isOpen.value = !isOpen.value
}

function closeDropdown(e: MouseEvent) {
  const target = e.target as HTMLElement
  if (!target.closest(`#${props.id}`)) {
    isOpen.value = false
  }
}

function isSelected(option: Option): boolean {
  return props.modelValue.includes(option.value)
}

function toggleItem(option: Option) {
  const newValue = [...props.modelValue]
  const index = newValue.indexOf(option.value)
  
  if (index === -1) {
    newValue.push(option.value)
  } else {
    newValue.splice(index, 1)
  }
  
  emit('update:modelValue', newValue)
}

function removeItem(option: Option) {
  const newValue = props.modelValue.filter(value => value !== option.value)
  emit('update:modelValue', newValue)
}

// Event listeners for clicking outside
onMounted(() => {
  document.addEventListener('click', closeDropdown)
})

onUnmounted(() => {
  document.removeEventListener('click', closeDropdown)
})
</script>