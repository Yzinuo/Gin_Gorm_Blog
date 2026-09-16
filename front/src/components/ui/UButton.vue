<script setup lang="ts">
import { defineProps, withDefaults } from 'vue'

const variants = {
  default: 'text-foreground bg-raised hover:bg-surface focus:ring-brand',
  primary: 'text-on-brand bg-brand hover:bg-brand-hover focus:ring-brand',
  secondary: 'text-foreground bg-raised hover:bg-surface focus:ring-brand',
  accent: 'text-on-brand bg-brand hover:bg-brand-hover focus:ring-brand',
  success: 'text-white bg-green-700 hover:bg-green-800 focus:ring-green-500',
  info: 'text-on-brand bg-brand hover:bg-brand-hover focus:ring-brand',
  warning: 'text-black bg-amber-400 hover:bg-amber-500 focus:ring-amber-400',
  error: 'text-white bg-red-700 hover:bg-red-800 focus:ring-red-500',
}

withDefaults(defineProps<{
  as?: 'button' | 'a'
  type?: 'default' | 'success' | 'info' | 'warning' | 'error' | 'primary' | 'secondary' | 'accent'
  size?: 'sm' | 'md' | 'lg'
  disabled?: boolean
}>(), {
  as: 'button',
  size: 'md',
  type: 'default',
  disabled: false,
})
</script>

<script lang="ts">
export default {
  name: 'UButton',
  inheritAttrs: false,
}
</script>

<template>
  <div class="flex items-center">
    <component
      :is="as"
      v-bind="$attrs"
      type="button"
      :disabled="disabled"
      :aria-disabled="disabled"
      class="inline-flex cursor-pointer items-center justify-center whitespace-nowrap rounded-md font-sans text-xs font-semibold leading-4 shadow-sm disabled:cursor-not-allowed disabled:opacity-30 focus:outline-none focus:ring-2 focus:ring-offset-2"
      :class="[
        variants[type],
        {
          'px-2 py-1': size === 'sm',
          'px-3 py-2': size === 'md',
          'px-4 py-3': size === 'lg',
        }]"
    >
      <slot />
    </component>
  </div>
</template>
