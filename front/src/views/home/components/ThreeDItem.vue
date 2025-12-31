<template>
  <div 
    ref="itemRef"
    class="threed-item"
    :style="itemStyle"
  >
    <slot></slot>
  </div>
</template>

<script setup>
import { inject, ref, watch, computed } from 'vue';

const props = defineProps({
  translateZ: { type: [Number, String], default: 0 },
  rotateX: { type: [Number, String], default: 0 },
  rotateY: { type: [Number, String], default: 0 },
  rotateZ: { type: [Number, String], default: 0 },
});

const isMouseEntered = inject('threed-mouse-state');
const itemRef = ref(null);

const itemStyle = computed(() => {
  if (isMouseEntered && isMouseEntered.value) {
    return {
      transform: `
        translateZ(${props.translateZ}px) 
        rotateX(${props.rotateX}deg) 
        rotateY(${props.rotateY}deg) 
        rotateZ(${props.rotateZ}deg)
      `,
      transition: 'transform 0.2s linear'
    };
  } else {
    return {
      transform: 'translateZ(0) rotate(0)',
      transition: 'transform 0.2s linear'
    };
  }
});
</script>

<style scoped>
.threed-item {
  width: fit-content;
  display: block;
}
</style>
