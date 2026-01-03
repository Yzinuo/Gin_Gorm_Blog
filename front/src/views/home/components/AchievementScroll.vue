<template>
  <div class="velocity-scroll-container">
    
    <!-- === 新增：标题区域 === -->
    <div class="section-header">
      <h2 class="section-title">
        <span class="icon">🏆</span> My Achievements
      </h2>
      <p class="section-desc">
        Record of my academic awards, open source contributions.
      </p>
    </div>
    <!-- ==================== -->

    <!-- 第一行：向左滚动 -->
    <div class="scroll-row left">
      <div class="scroll-track">
        <!-- 循环 4 次确保无缝 -->
        <div 
          v-for="(item, idx) in [...rowA, ...rowA, ...rowA, ...rowA]" 
          :key="'a-' + idx" 
          class="img-card"
        >
          <img :src="item.img" loading="lazy" class="card-image" />
          <div class="noise-layer"></div>
          <div class="card-overlay">
            <span class="achievement-tag">{{ item.tag }}</span>
            <span class="achievement-title">{{ item.title }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 第二行：向右滚动 -->
    <div class="scroll-row right">
      <div class="scroll-track">
        <div 
          v-for="(item, idx) in [...rowB, ...rowB, ...rowB, ...rowB]" 
          :key="'b-' + idx" 
          class="img-card"
        >
          <img :src="item.img" loading="lazy" class="card-image" />
          <div class="noise-layer"></div>
          <div class="card-overlay">
            <span class="achievement-tag">{{ item.tag }}</span>
            <span class="achievement-title">{{ item.title }}</span>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 两侧遮罩 -->
    <div class="fade-side left"></div>
    <div class="fade-side right"></div>
  </div>
</template>

<script setup>
// 数据保持不变...
const rowA = [
  { title: 'SCI Q2 一作论文发表', tag: 'Academic', img: 'https://img.heliar.top/file/1767414513796_6aa246b350f30f954a693f1b8f143029.png' },
  { title: '蓝桥杯全国三等奖', tag: 'Contest', img: 'https://img.heliar.top/file/1767414286247_52a06c29606e55b44aa0d08fe9c55182.jpg' },
  { title: '湖南省一等奖', tag: 'Contest', img: 'https://img.heliar.top/file/1767414296535_b6aff72e24aaa5816864754f060717ef.jpg' },
  { title: '字节跳动青训营', tag: 'Project', img: 'https://img.heliar.top/file/1767415028684_c134ee9bee6dcb7a2e43c1d516a1f339.jpg'}
];

const rowB = [
  { title: 'Github 700 Stars贡献者', tag: 'Open Source', img: 'https://img.heliar.top/file/1767414641143_image.png' },
  { title: 'MIT6.824', tag: 'Project', img: 'https://img.heliar.top/file/1767414443640_image.png' },
  { title: '后端项目：Zanecode', tag: 'Project', img: 'https://img.heliar.top/file/1767414735473_生成编程卡皮巴拉图标.png' },
  { title: '英语六级近500分', tag: 'Language', img: 'https://img.heliar.top/file/1767414879705_069b660cf57d140e7fff2291caa91667.jpg' },
];
</script>

<style lang="scss" scoped>
.velocity-scroll-container {
  position: relative;
  width: 100%;
  padding: 3rem 0 2rem; /* 上方留多点空间给标题 */
  background: transparent; 
  display: flex;
  flex-direction: column;
  gap: 2rem; /* 增加标题和滚动条之间的间距 */
  overflow: hidden;
}

/* === 新增：标题样式 === */
.section-header {
  text-align: center;
  padding: 0 1rem;
  margin-bottom: 0.5rem;
  z-index: 10; /* 确保在遮罩之上 */
}

.section-title {
  font-size: 2rem;
  font-weight: 800;
  color: #1f2937; /* 深灰黑色 */
  margin-bottom: 0.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  
  /* 移动端适配 */
  @media (min-width: 768px) {
    font-size: 2.5rem;
  }
}

.icon {
  font-size: 2rem;
  animation: wave 2s infinite;
  display: inline-block;
}

.section-desc {
  font-size: 2rem;
  color: #6b7280; /* 灰色副标题 */
  max-width: 1500px;
  margin: 0 auto;
  line-height: 2;
}

@keyframes wave {
  0%, 100% { transform: rotate(0deg); }
  25% { transform: rotate(-10deg); }
  75% { transform: rotate(10deg); }
}

/* === 滚动条样式 (保持不变) === */
.scroll-row {
  display: flex;
  width: 100%;
  overflow: hidden;
  transform: translate3d(0, 0, 0);
}

.scroll-track {
  display: flex;
  width: max-content;
  gap: 1.5rem; 
  padding: 0 0.75rem;
}

.scroll-row.left .scroll-track { animation: scroll-left 30s linear infinite; }
.scroll-row.right .scroll-track { animation: scroll-right 35s linear infinite; }

.velocity-scroll-container:hover .scroll-track { animation-play-state: paused; }

@keyframes scroll-left { 0% { transform: translateX(0); } 100% { transform: translateX(-25%); } }
@keyframes scroll-right { 0% { transform: translateX(-25%); } 100% { transform: translateX(0); } }

/* === 卡片样式 === */
.img-card {
  position: relative;
  width: 350px;
  height: 160px;
  border-radius: 12px;
  overflow: hidden;
  flex-shrink: 0;
  box-shadow: 0 4px 6px rgba(0,0,0,0.05);
  cursor: pointer;
  transition: transform 0.3s;
  background: #f3f4f6;

  &:hover {
    transform: scale(1.02);
    z-index: 10;
    box-shadow: 0 10px 20px rgba(0,0,0,0.1);
  }
}

.card-image { width: 100%; height: 100%; object-fit: cover; }

.noise-layer {
  position: absolute; inset: 0; pointer-events: none; z-index: 2;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noiseFilter'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.8' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noiseFilter)' opacity='1'/%3E%3C/svg%3E");
  background-size: 30%; opacity: 0.15; mix-blend-mode: overlay;
}

.card-overlay {
  position: absolute; inset: 0; z-index: 3;
  background: linear-gradient(to top, rgba(0,0,0,0.8) 0%, rgba(0,0,0,0.1) 60%, rgba(0,0,0,0) 100%);
  display: flex; flex-direction: column; justify-content: flex-end; padding: 1rem;
}

.achievement-tag {
  font-size: 0.75rem; color: #fbbf24; font-weight: bold; text-transform: uppercase; margin-bottom: 0.2rem;
}

.achievement-title {
  font-size: 1rem; color: #fff; font-weight: 600; text-shadow: 0 2px 4px rgba(0,0,0,0.8);
}

/* 遮罩：白色渐变 */
.fade-side {
  position: absolute; top: 0; bottom: 0; width: 10%; z-index: 20; pointer-events: none;
  &.left { left: 0; background: linear-gradient(to right, #ffffff, transparent); }
  &.right { right: 0; background: linear-gradient(to left, #ffffff, transparent); }
}
</style>
