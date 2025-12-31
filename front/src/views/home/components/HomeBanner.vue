<template>
  <div 
    class="banner-container"
    @mousemove="handleGlobalMove"
    @mouseleave="handleGlobalLeave"
  >
    <!-- 背景层 -->
    <div class="banner-bg"></div>

    <!-- 
      3D 内容层 (包含 Gopher 轨道) 
      注意：这里保留 global 视差，但 OrbitingCircles 放底层
    -->
    <div class="banner-content" :style="globalContentStyle">
      
      <!-- Gopher 轨道层 (在文字后面) -->
      <div class="orbit-layer">
        <OrbitingCircles
          :items="innerIcons"
          :radius="180"
          :duration="20"
          :icon-size="80"
          :reverse="true"
        >
          <template #default="{ item }">
            <div class="icon-wrapper glass">
              <img :src="item.icon" alt="" />
            </div>
          </template>
        </OrbitingCircles>
        
        <OrbitingCircles
          :items="outerIcons"
          :radius="250"
          :duration="40"
          :icon-size="120"
        >
          <template #default="{ item }">
             <div class="icon-wrapper glass">
               <img :src="item.icon" alt="" />
             </div>
          </template>
        </OrbitingCircles>
      </div>

      <!-- 
        === 核心修改：使用 ThreeDCard 包裹文字 === 
        Z轴前置，确保可以交互
      -->
      <div class="text-card-layer">
        <ThreeDCard>
          <ThreeDItem translateZ="60">
            <div class="title-wrapper">
              <Transition name="fade-slide" mode="out-in">
                <h1 :key="currentText" class="cute-title">
                  {{ currentText }}
                </h1>
              </Transition>
            </div>
          </ThreeDItem>

          <ThreeDItem translateZ="40" class="mt-4">
             <div class="subtitle">
               {{ typer?.output || 'Full Stack Developer' }} <span class="cursor">|</span>
             </div>
          </ThreeDItem>
        </ThreeDCard>
      </div>

    </div>

    <!-- 滚动按钮 -->
    <div class="scroll-down-btn" @click="scrollDown">
      <span class="i-ep:arrow-down-bold arrow-icon" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import OrbitingCircles from './OrbitingCircles.vue';
import ThreeDCard from './ThreeDCard.vue';
import ThreeDItem from './ThreeDItem.vue';

// --- Props ---
const props = defineProps({
  blogConfig: { type: Object, default: () => ({ website_name: 'SamuelQZQ' }) },
  typer: { type: Object, default: () => ({ output: '' }) }
});
const emit = defineEmits(['scroll-down']);

// --- 图标数据 (使用上一轮优化的 CDN 链接) ---
const innerIcons = [
  { icon: 'https://go.dev/images/gophers/machine-colorized.svg' },
  { icon: 'https://go.dev/images/gophers/biplane.svg' },
  { icon: 'https://go.dev/images/gophers/motorcycle.svg' },
];

const outerIcons = [
  { icon: '/images/docker-original.svg' },   // 正确路径
  { icon: '/images/git-original.svg' },      // 正确路径
  { icon: '/images/go-original-wordmark.svg' } // 正确路径
];

// --- 文字轮播逻辑 ---
const textList = [props.blogConfig.website_name || 'SamuelQZQ', 'Developer', 'Dreamer', 'Creator'];
const currentText = ref(textList[0]);
let textInterval = null;

onMounted(() => {
  let idx = 0;
  textInterval = setInterval(() => {
    idx = (idx + 1) % textList.length;
    currentText.value = textList[idx];
  }, 3000); // 每3秒切换
});

onUnmounted(() => {
  if (textInterval) clearInterval(textInterval);
});

// --- 全局 Banner 视差 (保留给 Orbit 使用) ---
const containerRef = ref(null);
const mouseX = ref(0);
const mouseY = ref(0);
const isHovering = ref(false);

const handleGlobalMove = (e) => {
  // 只做轻微的背景移动，把强烈的 3D 交互留给 ThreeDCard
  const w = window.innerWidth;
  const h = window.innerHeight;
  mouseX.value = (e.clientX / w) * 2 - 1;
  mouseY.value = (e.clientY / h) * 2 - 1;
};

const handleGlobalLeave = () => {
  mouseX.value = 0;
  mouseY.value = 0;
};

const globalContentStyle = computed(() => {
  // 这里只对整个容器做非常轻微的倾斜，避免和 ThreeDCard 冲突
  return {
    transform: `perspective(1000px) rotateX(${-mouseY.value * 2}deg) rotateY(${mouseX.value * 2}deg)`,
    transition: 'transform 0.5s ease-out'
  };
});

const scrollDown = () => {
  emit('scroll-down');
  window.scrollTo({ top: window.innerHeight, behavior: 'smooth' });
};
</script>

<style lang="scss" scoped>
/* 引入可爱的 Google 字体 */
@import url('https://fonts.googleapis.com/css2?family=Titan+One&display=swap');

$text-color: #ffffff;

.banner-container {
  position: relative;

  height: 100vh;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
  background-color: #000;
}

.banner-bg {
  position: absolute;
  inset: 0;
  z-index: 0;
  background-image: url('https://img.heliar.top/file/1767183149047_网站背景.jpg'); 
  background-size: cover;
  background-position: center;
  pointer-events: none;
  &::after {
    content: '';
    position: absolute;
    inset: 0;
    background: rgba(0, 0, 0, 0.4); 
  }
}

.banner-content {
  position: relative;
  z-index: 10;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  transform-style: preserve-3d;
}

.orbit-layer {
  position: absolute;
  inset: 0;
  z-index: 10;
  transform: translateZ(0px);
  pointer-events: none; /* 让鼠标穿透去触发中间的卡片 */
}

.icon-wrapper {
  width: 100%; height: 100%;
  display: flex; align-items: center; justify-content: center;
  img { width: 70%; height: 70%; object-fit: contain; }
  &.glass {
    background: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(2px);
    border-radius: 50%;
    border: 1px solid rgba(255, 255, 255, 0.2);
  }
}

/* === 卡片与文字层 === */
.text-card-layer {
  z-index: 50; /* 确保在最上层 */
  transform: translateZ(50px);
}

.title-wrapper {
  // 固定高度，防止切换文字时抖动
  height: 120px; 
  display: flex;
  align-items: center;
  justify-content: center;
  
  @media (min-width: 1024px) { height: 160px; }
}

/* === 可爱 3D 字体样式 === */
.cute-title {
  font-family: 'Titan One', cursive; /* 圆润的卡通字体 */
  font-size: 4rem;
  line-height: 1;
  text-align: center;
  
  /* 颜色方案：参考图片中的蓝/粉/白 */
  background: linear-gradient(to bottom, #ffffff 0%, #dbeafe 100%);
  -webkit-background-clip: text;
  color: transparent; /* 为了让背景渐变显示 */
  
  /* 核心：描边 + 多重阴影制造 3D 厚度感 */
  -webkit-text-stroke: 3px #3b82f6; /* 蓝色外边框 */
  paint-order: stroke fill; /* 确保描边不遮挡内部填充 */
  
  /* 这里的阴影是关键：层层叠加模拟厚度 */
  filter: drop-shadow(0px 4px 0px #2563eb) 
          drop-shadow(0px 8px 0px #1d4ed8)
          drop-shadow(0px 12px 10px rgba(0,0,0,0.5));
          
  @media (min-width: 1024px) {
    font-size: 6rem;
    -webkit-text-stroke: 4px #3b82f6;
  }
}

.subtitle {
  font-family: 'Consolas', monospace;
  font-size: 1.2rem;
  color: #fff;
  background: rgba(0,0,0,0.4);
  padding: 0.5rem 1.2rem;
  border-radius: 99px;
  backdrop-filter: blur(4px);
  border: 1px solid rgba(255,255,255,0.2);
  box-shadow: 0 4px 10px rgba(0,0,0,0.3);
}

.cursor {
  animation: blink 1s infinite;
  color: #60a5fa;
}

/* === 切换动画 === */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.5s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(20px) scale(0.8);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-20px) scale(1.1);
}

.scroll-down-btn {
  position: absolute; bottom: 2rem; z-index: 30;
  color: #fff; font-size: 2rem; cursor: pointer;
  animation: bounce 2s infinite;
}

@keyframes blink { 0%, 100% { opacity: 1; } 50% { opacity: 0; } }
@keyframes bounce {
  0%, 20%, 50%, 80%, 100% { transform: translateY(0); }
  40% { transform: translateY(-10px); }
  60% { transform: translateY(-5px); }
}
</style>
