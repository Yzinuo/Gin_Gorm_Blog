<script setup>
import { ref } from 'vue'
import { useWindowScroll, watchThrottled } from '@vueuse/core'
import { Icon } from '@iconify/vue'

const { y } = useWindowScroll()
const styleVal = ref('')
watchThrottled(y, () => {
  styleVal.value = (y.value > 20) ? 'opacity: 1; transform: translateX(-40px);' : ''
}, { throttle: 100 })

const options = [
  {
    icon: 'fluent:arrow-up-12-filled',
    fn: () => window.scrollTo({ behavior: matchMedia('(prefers-reduced-motion: reduce)').matches ? 'instant' : 'smooth', top: 0 }),
  },
]
</script>

<template>
  <div class="fixed bottom-20 z-4 text-white transition-600 -right-9 space-y-1" :style="styleVal">
    <button
      v-for="item of options" :key="item.icon"
      type="button"
      aria-label="回到顶部"
      class="f-c-c cursor-pointer rounded-sm bg-brand p-2 text-on-brand duration-300 hover:bg-brand-hover"
      @click="item.fn"
    >
      <Icon class="h-5 w-5" :icon="item.icon" />
    </button>
  </div>
</template>
