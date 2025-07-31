<template>
    <div
      ref="container"
      class="virtual-list-container"
      :style="{ height: `${containerHeight}px`, overflow: 'auto' }"
      @scroll="handleScroll"
    >
      <div :style="{ height: `${totalHeight}px`, position: 'relative' }">
        <div
          v-for="item in visibleItems"
          :key="getItemKey(item.data)"
          :style="{
            position: 'absolute',
            top: `${item.top}px`,
            left: 0,
            right: 0,
            height: `${itemHeight}px`
          }"
        >
          <slot :item="item.data" :index="item.index" />
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, computed } from 'vue'
  
  interface Props<T = unknown> {
    items: T[]
    itemHeight: number
    containerHeight: number
    buffer?: number
    getItemKey?: (item: T) => string | number
  }
  
  const props = withDefaults(defineProps<Props>(), {
    buffer: 5,
    getItemKey: (item: unknown) => (item as { id?: string | number })?.id || Math.random()
  })
  
  const container = ref<HTMLElement>()
  const scrollTop = ref(0)
  
  const totalHeight = computed(() => props.items.length * props.itemHeight)
  
  const visibleStart = computed(() => {
    return Math.max(0, Math.floor(scrollTop.value / props.itemHeight) - props.buffer)
  })
  
  const visibleEnd = computed(() => {
    const visibleCount = Math.ceil(props.containerHeight / props.itemHeight)
    return Math.min(props.items.length, visibleStart.value + visibleCount + props.buffer * 2)
  })
  
  const visibleItems = computed(() => {
    return props.items.slice(visibleStart.value, visibleEnd.value).map((item, index) => ({
      data: item,
      index: visibleStart.value + index,
      top: (visibleStart.value + index) * props.itemHeight
    }))
  })
  
  const handleScroll = (event: Event) => {
    const target = event.target as HTMLElement
    scrollTop.value = target.scrollTop
  }
  
  // Scroll to specific item
  const scrollToItem = (index: number) => {
    if (container.value) {
      container.value.scrollTop = index * props.itemHeight
    }
  }
  
  // Scroll to top
  const scrollToTop = () => {
    if (container.value) {
      container.value.scrollTop = 0
    }
  }
  
  defineExpose({
    scrollToItem,
    scrollToTop
  })
  </script>
  
  <style scoped>
  .virtual-list-container {
    will-change: scroll-position;
  }
  </style>