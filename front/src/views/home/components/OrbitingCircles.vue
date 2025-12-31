<template>
  <!-- 容器：绝对定位，确保与父容器重叠 -->
  <div class="orbiting-circles-container">
    
    <!-- 1. SVG 轨道路径 (可选) -->
    <svg
      v-if="path"
      xmlns="http://www.w3.org/2000/svg"
      version="1.1"
      class="orbit-path"
    >
      <circle
        class="path-circle"
        cx="50%"
        cy="50%"
        :r="radius"
        fill="none"
      />
    </svg>

    <!-- 2. 子元素渲染 (通过 Slot) -->
    <!-- 
      注意：Vue 的 Slot 机制与 React.Children 不同。
      我们在这里并不直接循环 children，而是让外部传递一个 list，
      或者外部直接使用 v-for 调用 OrbitingCircles（如果是多轨道的场景）。
      
      但为了完美复刻 React 版“自动分布子元素”的逻辑，
      最 Vue 的方式是：这个组件作为一个“轨道”，你往里面放一个 items 数组，
      组件负责把它们均匀分布并旋转。
    -->
    <div
      v-for="(item, index) in items"
      :key="index"
      class="orbit-item"
      :class="{ 'reverse-orbit': reverse }"
      :style="getItemStyle(index)"
    >
      <!-- 插槽：将当前 item 数据回传给父组件自定义渲染 -->
      <slot :item="item" :index="index"></slot>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  // 需要在轨道上展示的数据列表
  items: {
    type: Array,
    default: () => []
  },
  reverse: {
    type: Boolean,
    default: false
  },
  duration: {
    type: Number,
    default: 20 // 秒
  },
  delay: {
    type: Number,
    default: 0
  },
  radius: {
    type: Number,
    default: 160
  },
  path: {
    type: Boolean,
    default: true
  },
  iconSize: {
    type: Number,
    default: 30
  },
  speed: {
    type: Number,
    default: 1
  }
});

// 计算总时长
const calculatedDuration = computed(() => props.duration / props.speed);

// 计算每个 Item 的样式
const getItemStyle = (index) => {
  // 如果 items 为空，避免除以0
  const count = props.items.length || 1;
  // 计算初始角度偏移，让图标均匀分布
  const angleStep = 360 / count;
  const initialAngle = angleStep * index;

  return {
    '--duration': `${calculatedDuration.value}s`,
    '--radius': `${props.radius}px`,
    '--icon-size': `${props.iconSize}px`,
    // 我们通过 delay 负值来模拟“分布在圆周不同位置”的效果
    // 动画本质是一样的，只是起跑时间不同
    'animation-delay': `-${(calculatedDuration.value / count) * index}s`
  };
};
</script>

<style scoped>
.orbiting-circles-container {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  pointer-events: none; /* 让鼠标穿透，不影响文字交互 */
  display: flex;
  justify-content: center;
  align-items: center;
}

.orbit-path {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
}

.path-circle {
  stroke: rgba(255, 255, 255, 0.1); /* 白色半透明轨道 */
  stroke-width: 1;
}

.orbit-item {
  position: absolute;
  display: flex;
  align-items: center;
  justify-content: center;
  width: var(--icon-size);
  height: var(--icon-size);
  border-radius: 50%;
  
  /* 核心动画 */
  animation: orbit var(--duration) linear infinite;
}

/* 反向旋转类 */
.reverse-orbit {
  animation-direction: reverse;
}

/* 
  关键帧动画：Orbit
  原理：利用 transform 将元素从中心移到半径处，然后绕中心旋转。
  
  注意：React 版通常用的是 CSS Variables + calc() 计算角度。
  这里为了简化且高性能，采用通用的做法：
  旋转父容器，同时反向旋转子元素（如果不想让图标倒立）。
  
  但 React Magic UI 的 OrbitingCircles 实现略有不同，它是让 transform-origin 保持在中心，
  然后旋转自身。
*/
@keyframes orbit {
  0% {
    transform: rotate(0deg) translateX(var(--radius)) rotate(0deg);
  }
  100% {
    transform: rotate(360deg) translateX(var(--radius)) rotate(-360deg);
  }
}
</style>
