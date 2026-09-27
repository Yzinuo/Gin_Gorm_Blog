<script setup>
import { nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { createResumeScene } from '../scene'
import cameraMap from '../camera-map.json'

const props = defineProps({ story: Object, stickers: { type: Array, default: () => [] }, modelUrl: { type: String, default: '' } })
const emit = defineEmits(['chapter'])
const stage = ref(null)
const canvas = ref(null)
const status = ref('idle')
const progress = ref(0)
const attempt = ref(0)
let scene
let disposed = false

async function start() {
  scene?.dispose()
  const current = ++attempt.value
  status.value = 'loading'
  progress.value = 0
  await nextTick()
  if (disposed || current !== attempt.value)
    return
  scene = createResumeScene({
    canvas: canvas.value,
    stage: stage.value,
    map: cameraMap,
    getSections: () => props.story?.querySelectorAll('[data-frame]') || [],
    getStickers: () => props.stickers,
    modelURL: props.modelUrl || `${import.meta.env.BASE_URL}resume/resume-ready.glb`,
    onProgress: value => progress.value = value,
    onReady: () => status.value = 'ready',
    onChapter: value => emit('chapter', value),
    onError: () => status.value = 'error',
  })
}

defineExpose({ focus: frame => scene?.focus(frame) })
watch(() => props.stickers, () => scene?.updateStickers(), { deep: true })
onMounted(() => {
  const connection = navigator.connection
  const constrainedNetwork = connection?.saveData || ['slow-2g', '2g', '3g'].includes(connection?.effectiveType)
  // The model is large enough that on a constrained connection visitors should
  // decide whether to download it; all written content remains available.
  if (constrainedNetwork || matchMedia('(max-width: 640px)').matches)
    return
  const schedule = window.requestIdleCallback || (callback => setTimeout(callback, 1200))
  schedule(() => !disposed && start(), { timeout: 4000 })
})
onUnmounted(() => {
  disposed = true
  scene?.dispose()
})
</script>

<template>
  <div ref="stage" class="resume-stage" :data-status="status">
    <canvas :key="attempt" ref="canvas" aria-label="随滚动切换镜头、眼神跟随鼠标的个人 3D 模型" />
    <div class="stage-shade" />
    <div v-if="status !== 'ready'" class="scene-status" role="status" aria-live="polite">
      <span class="loading-monogram" aria-hidden="true">Z.</span>
      <template v-if="status === 'idle'">
        <p>互动场景已准备好</p>
        <button type="button" class="the-button" @click="start">
          加载 3D 场景
        </button>
        <small>移动网络下按需加载，不影响阅读下方经历。</small>
      </template>
      <template v-else-if="status === 'loading'">
        <p>正在展开我的世界</p>
        <progress :value="progress" max="100" aria-label="3D 场景加载进度" />
        <span>{{ progress }}%</span>
        <small>下方经历可以先行阅读</small>
      </template>
      <template v-else>
        <p>3D 场景暂时无法显示</p>
        <small>你仍然可以阅读完整经历。</small>
        <button type="button" class="the-button" @click="start">
          重新加载
        </button>
      </template>
    </div>
    <slot v-if="status === 'ready'" />
  </div>
</template>

<style scoped>
.resume-stage {
  position: fixed;
  inset: 0 36% 0 0;
  overflow: hidden;
  background: radial-gradient(ellipse at 54% 44%, #82655f 0, #503c3d 38%, #251c22 78%);
  color: #ffffff;
}
canvas {
  display: block;
  width: 100%;
  height: 100%;
}
.stage-shade {
  position: absolute;
  inset: 0;
  pointer-events: none;
  background: linear-gradient(0deg, #170e1677, transparent 33%);
}
.scene-status {
  position: absolute;
  inset: 0;
  padding: 80px 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  gap: 18px;
  text-align: center;
  background: var(--bg-stage);
}
.loading-monogram {
  font: 72px Georgia, serif;
  color: var(--brand);
}
.scene-status p {
  font-size: 16px;
  letter-spacing: 2px;
  color: #ffffff;
}
.scene-status .the-button { padding: 10px 18px; border-radius: 6px; color: #171114; background: var(--brand); cursor: pointer; }
.scene-status .the-button:focus-visible { outline: 2px solid #ffffff; outline-offset: 4px; }
.scene-status small, .scene-status > span:last-of-type {
  color: rgba(255, 255, 255, 0.7);
  font-size: 12px;
}
progress {
  width: min(220px, 80%);
  height: 5px;
  accent-color: var(--brand);
}
@media (max-width: 900px) {
  .resume-stage {
    right: 42%;
  }
}
@media (max-width: 640px) {
  .resume-stage {
    inset: 0 0 auto;
    height: 54svh;
    z-index: 5;
  }
}
</style>
