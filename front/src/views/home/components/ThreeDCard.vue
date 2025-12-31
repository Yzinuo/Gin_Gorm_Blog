<template>
  <div 
    class="threed-container" 
    style="perspective: 1000px;"
  >
    <div
      ref="containerRef"
      class="threed-body"
      :style="bodyStyle"
      @mouseenter="handleMouseEnter"
      @mousemove="handleMouseMove"
      @mouseleave="handleMouseLeave"
    >
      <slot></slot>
    </div>
  </div>
</template>

<script setup>
import { ref, provide, computed } from 'vue';

const containerRef = ref(null);
const isMouseEntered = ref(false);
const rotateX = ref(0);
const rotateY = ref(0);

// 提供给子组件 (Item) 使用
provide('threed-mouse-state', isMouseEntered);

const handleMouseMove = (e) => {
  if (!containerRef.value) return;
  const { left, top, width, height } = containerRef.value.getBoundingClientRect();
  
  // 计算旋转角度 (类似 React 版的逻辑)
  const x = (e.clientX - left - width / 2) / 25;
  const y = (e.clientY - top - height / 2) / 25;
  
  rotateX.value = x;
  rotateY.value = y;
};

const handleMouseEnter = () => {
  isMouseEntered.value = true;
};

const handleMouseLeave = () => {
  isMouseEntered.value = false;
  // 复位
  rotateX.value = 0;
  rotateY.value = 0;
};

const bodyStyle = computed(() => {
  return {
    transform: `rotateY(${rotateX.value}deg) rotateX(${rotateY.value}deg)`,
    transformStyle: 'preserve-3d',
    transition: 'transform 0.1s ease-out' // 鼠标移动时平滑一点
  };
});
</script>

<style scoped>
.threed-container {
  display: flex;
  align-items: center;
  justify-content: center;
  /* 允许鼠标事件穿透容器边缘，只响应 Body */
  padding: 2rem; 
}

.threed-body {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  transform-style: preserve-3d;
  /* 这里不设背景色，保持透明，或者你可以加 glass 效果 */
}
</style>
