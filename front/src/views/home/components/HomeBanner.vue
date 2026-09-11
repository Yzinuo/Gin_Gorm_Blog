<template>
  <div class="banner-container">
    <!-- WebGL X-Ray 核心画布 (透视人物) -->
    <XRayCanvas class="xray-canvas-layer" />

    <!-- 顶部/前景层：文字与信息浮于左方 -->
    <div class="banner-content-overlay">
      <div class="hero-left-col">
        <!-- 切换标题 (原样式，居左排版) -->
        <div class="title-wrapper">
          <Transition name="fade-blur">
            <h1 :key="currentText" class="dreamy-title">
              {{ currentText }}
            </h1>
          </Transition>
        </div>

        <!-- 副标题 / 打字机 (无毛玻璃，纯净质感) -->
        <div class="subtitle-clean">
          <span class="subtitle-bullet">///</span>
          <span class="text-content">
            {{ (typer?.output || '日日自新，步步强于昨日') }}
          </span>
          <span class="cyber-cursor"></span>
        </div>
      </div>
    </div>

    <!-- 底部滚动按钮 -->
    <div class="scroll-down-btn" title="向下滚动" @click="scrollDown">
      <span class="i-ep:arrow-down-bold arrow-icon" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import XRayCanvas from './XRayCanvas.vue';

const props = defineProps({
  blogConfig: { type: Object, default: () => ({ website_name: 'Zane' }) },
  typer: { type: Object, default: () => ({ output: '' }) }
});
const emit = defineEmits(['scroll-down']);

// 文字轮播列表
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
  overflow: hidden;
  background-color: #0d0e15;
}

/* X-Ray 画布填满整屏 */
.xray-canvas-layer {
  position: absolute;
  inset: 0;
  z-index: 1;
}

/* 前景内容层：穿透鼠标事件，保证 X-Ray 探针随时感应 */
.banner-content-overlay {
  position: absolute;
  inset: 0;
  z-index: 10;
  pointer-events: none;
  display: flex;
  align-items: center;
  padding: 0 4rem;
  max-width: 1600px;
  margin: 0 auto;

  @media (max-width: 768px) {
    padding: 0 1.5rem;
  }
}

/* 左侧标题区域排版 */
.hero-left-col {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  text-align: left;
  max-width: 650px;
}

/* 标题容器：Grid 布局保持重叠切换 */
.title-wrapper {
  height: 120px;
  display: grid;
  place-items: start;
  align-items: center;
  @media (min-width: 1024px) {
    height: 160px;
  }
}

/* 原样式 Dreamy 标题 (保留 Titan One、描边、光效与呼吸动效，居左对齐) */
.dreamy-title {
  grid-area: 1 / 1;
  font-family: 'Titan One', cursive;
  font-size: 4.5rem;
  line-height: 1.1;
  text-align: left;
  color: #ffffff;
  pointer-events: auto;
  user-select: none;

  /* 柔和描边 (0.4 透明度) */
  -webkit-text-stroke: 1px rgba(255, 255, 255, 0.4);

  /* 原柔和光效阴影 */
  text-shadow:
    0 5px 15px rgba(0, 0, 0, 0.4),
    0 0 25px rgba(255, 60, 60, 0.35);

  animation: breathe 4s ease-in-out infinite;

  @media (min-width: 1024px) {
    font-size: 6.8rem;
    text-shadow:
      0 8px 24px rgba(0, 0, 0, 0.45),
      0 0 35px rgba(255, 60, 60, 0.4);
  }
}

/* 纯净无毛玻璃副标题 */
.subtitle-clean {
  font-family: 'Consolas', 'Fira Code', monospace;
  font-size: 1.15rem;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.9);
  margin-top: 1rem;
  display: flex;
  align-items: center;
  gap: 10px;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.8);
  pointer-events: auto;

  @media (min-width: 1024px) {
    font-size: 1.3rem;
  }
}

.subtitle-bullet {
  color: #ff3b3b;
  font-weight: 700;
  letter-spacing: 2px;
}

.cyber-cursor {
  display: inline-block;
  width: 2px;
  height: 1.2em;
  background-color: #ff3b3b;
  animation: blink 1s step-end infinite;
  box-shadow: 0 0 8px #ff3b3b;
}

@keyframes breathe {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-8px); }
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0; }
}

/* 标题切换过渡动画 */
.fade-blur-enter-active,
.fade-blur-leave-active {
  transition: all 0.6s cubic-bezier(0.2, 0.8, 0.2, 1);
}

.fade-blur-leave-to {
  opacity: 0;
  transform: scale(1.08) translateY(-10px);
  filter: blur(8px);
}

.fade-blur-enter-from {
  opacity: 0;
  transform: scale(0.92) translateY(10px);
  filter: blur(8px);
}

/* 滚动按钮 */
.scroll-down-btn {
  position: absolute;
  bottom: 2.5rem;
  left: 50%;
  transform: translateX(-50%);
  z-index: 20;
  color: rgba(255, 255, 255, 0.75);
  font-size: 2rem;
  cursor: pointer;
  pointer-events: auto;
  animation: float-btn 2.2s infinite ease-in-out;
  transition: all 0.3s;

  &:hover {
    color: #ff3b3b;
    filter: drop-shadow(0 0 8px rgba(255, 59, 59, 0.8));
  }
}

@keyframes float-btn {
  0%, 100% { transform: translate(-50%, 0); }
  50% { transform: translate(-50%, 10px); }
}
</style>
