<template>
  <div 
    class="banner-container"
    @mousemove="handleGlobalMove"
    @mouseleave="handleGlobalLeave"
  >
    <!-- 背景层 -->
    <div class="banner-bg"></div>

    <!-- 3D 内容层 -->
    <div class="banner-content" :style="globalContentStyle">
      
      <!-- Gopher 轨道层 -->
      <div class="orbit-layer">
        <OrbitingCircles
          :items="innerIcons"
          :radius="180"
          :duration="20"
          :icon-size="80"
          :reverse="true"
        >
          <template #default="{ item }">
            <div class="icon-wrapper glass-icon">
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
             <div class="icon-wrapper glass-icon">
               <img :src="item.icon" alt="" />
             </div>
          </template>
        </OrbitingCircles>
      </div>

      <!-- 文字层 -->
      <div class="text-card-layer">
        <ThreeDCard>
          <ThreeDItem translateZ="80" class="floating-item">
            <div class="title-wrapper">
              <!-- 
                 注意：这里去掉了 mode="out-in" 
                 配合下方的 Grid 布局实现无缝重叠切换 
              -->
              <Transition name="fade-blur">
                <h1 :key="currentText" class="dreamy-title">
                  {{ currentText }}
                </h1>
              </Transition>
            </div>
          </ThreeDItem>

          <ThreeDItem translateZ="50" class="mt-6">
          <div class="subtitle-glass">
      <span class="text-content">
     
      {{ (typer?.output || '日日自新，步步强于昨日') }}
    </span>
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

const props = defineProps({
  blogConfig: { type: Object, default: () => ({ website_name: 'Zane' }) },
  typer: { type: Object, default: () => ({ output: '' }) }
});
const emit = defineEmits(['scroll-down']);

// 图标数据
const innerIcons = [
  { icon: 'https://go.dev/images/gophers/machine-colorized.svg' },
  { icon: 'https://go.dev/images/gophers/biplane.svg' },
  { icon: 'https://go.dev/images/gophers/motorcycle.svg' },
];
const outerIcons = [
  { icon: '/images/docker-original.svg' },
  { icon: '/images/git-original.svg' },
  { icon: '/images/go-original-wordmark.svg' }
];

// 文字轮播
const textList = [props.blogConfig.website_name || 'Zane', 'Developer', 'Dreamer', 'Creator'];
const currentText = ref(textList[0]);
let textInterval = null;

onMounted(() => {
  let idx = 0;
  textInterval = setInterval(() => {
    idx = (idx + 1) % textList.length;
    currentText.value = textList[idx];
  }, 3500); 
});

onUnmounted(() => {
  if (textInterval) clearInterval(textInterval);
});

// 全局视差
const containerRef = ref(null);
const mouseX = ref(0);
const mouseY = ref(0);

const handleGlobalMove = (e) => {
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
  return {
    transform: `perspective(1000px) rotateX(${-mouseY.value * 2}deg) rotateY(${mouseX.value * 2}deg)`,
    transition: 'transform 0.8s cubic-bezier(0.2, 0.8, 0.2, 1)'
  };
});

const scrollDown = () => {
  emit('scroll-down');
  window.scrollTo({ top: window.innerHeight, behavior: 'smooth' });
};
</script>

<style lang="scss" scoped>
@import url('https://fonts.googleapis.com/css2?family=Titan+One&display=swap');

.banner-container {
  position: relative;
  height: 100vh;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
  background-color: #1a1b26;
}

.banner-bg {
  position: absolute;
  inset: 0;
  z-index: 0;
  background-image: url('https://img.heliar.top/file/1767411799555_wallhaven-wej9vp_2560x1440.png'); 
  background-size: cover;
  background-position: center;
  pointer-events: none;
  &::after {
    content: '';
    position: absolute;
    inset: 0;
    background: linear-gradient(
      to bottom, 
      rgba(20, 30, 60, 0.2) 0%, 
      rgba(10, 15, 40, 0.5) 100%
    );
    backdrop-filter: blur(1px);
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
  pointer-events: none;
}

.glass-icon {
  width: 100%; height: 100%;
  display: flex; align-items: center; justify-content: center;
  background: rgba(255, 255, 255, 0.03);
  backdrop-filter: blur(3px);
  border-radius: 50%;
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 4px 10px rgba(0,0,0,0.1);
  img { 
    width: 65%; height: 65%; object-fit: contain; 
    filter: drop-shadow(0 2px 4px rgba(0,0,0,0.2));
  }
}

.text-card-layer {
  z-index: 50;
  transform: translateZ(50px);
}

/* === 容器优化：Grid 布局 === */
.title-wrapper {
  height: 140px; 
  /* 关键：Grid 布局 + 居中 */
  display: grid;
  place-items: center;
  @media (min-width: 1024px) { height: 180px; }
}

/* === 标题优化：柔和光效 + 强制重叠 === */
.dreamy-title {
  /* 关键：让进入和离开的文字都在同一个 Grid 格子里 */
  grid-area: 1 / 1;
  
  font-family: 'Titan One', cursive;
  font-size: 4.5rem;
  line-height: 1.1;
  text-align: center;
  color: #ffffff;
  
  /* 柔和描边 (0.4 透明度) */
  -webkit-text-stroke: 1px rgba(255, 255, 255, 0.4);
  
  /* 柔和阴影 (去掉了刺眼的高光) */
  text-shadow: 
    0 5px 15px rgba(0, 0, 0, 0.25), 
    0 0 20px rgba(100, 180, 255, 0.3);

  animation: breathe 4s ease-in-out infinite;

  @media (min-width: 1024px) {
    font-size: 7rem;
    text-shadow: 
      0 8px 20px rgba(0, 0, 0, 0.25),
      0 0 30px rgba(100, 180, 255, 0.3);
  }
}

.subtitle-glass {
  font-family: 'Consolas', monospace;
  font-size: 1.1rem;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.95);
  background: rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(6px);
  padding: 0.6rem 1.5rem;
  border-radius: 99px;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: all 0.3s;
  &:hover {
    background: rgba(0, 0, 0, 0.3);
    border-color: rgba(255, 255, 255, 0.3);
  }
}

.cursor {
  display: inline-block; width: 2px; height: 1.2em;
  background-color: #60a5fa; animation: blink 1s step-end infinite;
}

@keyframes breathe {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

/* === 关键修复：无缝重叠切换动画 === */
.fade-blur-enter-active,
.fade-blur-leave-active {
  /* 
     这里千万不要写 opacity: 0; 
     只定义过渡时间和曲线
  */
  transition: all 0.6s ease;
}

/* 离开状态：变透明 + 放大 + 模糊 */
.fade-blur-leave-to {
  opacity: 0;
  transform: scale(1.1);
  filter: blur(10px);
}

/* 进入前状态：透明 + 缩小 + 模糊 */
.fade-blur-enter-from {
  opacity: 0;
  transform: scale(0.9);
  filter: blur(10px);
}

@keyframes blink { 0%, 100% { opacity: 1; } 50% { opacity: 0; } }

.scroll-down-btn {
  position: absolute; bottom: 3rem; z-index: 30;
  color: rgba(255,255,255,0.7); 
  font-size: 2rem; cursor: pointer;
  animation: float-btn 2s infinite ease-in-out;
  transition: color 0.3s;
  &:hover { color: #fff; }
}

@keyframes float-btn {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(10px); }
}
</style>
