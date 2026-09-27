<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'

const props = defineProps({ asset: { type: Object, default: null } })

const containerRef = ref(null)
const canvasRef = ref(null)
const isReady = ref(false)
const hasError = ref(false)
let disposed = false
const reducedMotion = matchMedia('(prefers-reduced-motion: reduce)').matches
const beforeURL = computed(() => props.asset?.before?.src || '/images/Before-1672.webp')
const afterURL = computed(() => props.asset?.after?.src || '/images/After-1672.webp')
const beforeSrcset = computed(() => {
  const variants = props.asset?.before?.srcset || [{ src: '/images/Before-960.webp', width: 960 }, { src: '/images/Before-1672.webp', width: 1672 }]
  return variants.map(item => `${item.src} ${item.width}w`).join(', ')
})
let positionBuffer = null
let isVisible = true
let intersectionObserver = null
const hasStarted = ref(false)

// 配置参数
const CONFIG = {
  imgWidth: 1672,
  imgHeight: 941,
  baseRadius: 180, // 基础透视半径 (px)
  maxStretchRadius: 260, // 运动时最大拉伸半径
  lerpFactor: 0.1, // 惯性平滑度 (越小越柔和流体感越强)
  rimWidthRatio: 0.055, // 边缘光圈宽度比例
  rimColor: [1.0, 0.18, 0.18], // 霓虹激光红
  rimGlowIntensity: 1.8,
}

let gl = null
let program = null
let textureBefore = null
let textureAfter = null
let animationFrameId = null

// 物理状态追踪
const state = {
  targetX: 0,
  targetY: 0,
  currentX: 0,
  currentY: 0,
  prevX: 0,
  prevY: 0,
  velocityX: 0,
  velocityY: 0,
  speed: 0,
  currentRadius: CONFIG.baseRadius,
  targetRadius: CONFIG.baseRadius,
  hoverFactor: 0.0, // 0: 完全淡出, 1: 完全显现
  targetHover: 0.0,
  hasInitialMoved: false,
  startTime: performance.now(),
  lastFrameTime: performance.now(),
}

// 顶点着色器
const VS_SOURCE = `
attribute vec2 a_position;
varying vec2 v_uv;

void main() {
  v_uv = (a_position + 1.0) * 0.5;
  // 翻转 Y 轴匹配 WebGL 纹理坐标
  v_uv.y = 1.0 - v_uv.y;
  gl_Position = vec4(a_position, 0.0, 1.0);
}
`

// 片段着色器 (核心液态 X-Ray 渲染器)
const FS_SOURCE = `
precision highp float;

varying vec2 v_uv;

uniform sampler2D u_texBefore;
uniform sampler2D u_texAfter;
uniform vec2 u_resolution;
uniform vec2 u_imageResolution;
uniform float u_imageAlign;
uniform vec2 u_mouse;       // 归一化屏幕坐标 [0, 1]
uniform vec2 u_velocity;    // 移动速度矢量
uniform float u_radius;     // 像素半径
uniform float u_time;
uniform float u_hover;      // 显隐过渡因子 [0, 1]

// 计算保持 Cover 比例居中的 UV 坐标
vec2 getCoverUV(vec2 uv, vec2 screenRes, vec2 imgRes) {
  float screenRatio = screenRes.x / screenRes.y;
  float imgRatio = imgRes.x / imgRes.y;
  
  vec2 newUV = uv;
  if (screenRatio > imgRatio) {
    // 屏幕更宽，按宽度撑满，Y方向上下居中裁切
    float scale = imgRatio / screenRatio;
    newUV.y = (uv.y - 0.5) * scale + 0.5;
  } else {
    // 屏幕更高，按高度撑满，X方向左右居中裁切
    float scale = screenRatio / imgRatio;
    newUV.x = uv.x * scale + (1.0 - scale) * u_imageAlign;
  }
  return newUV;
}

void main() {
  vec2 uv = getCoverUV(v_uv, u_resolution, u_imageResolution);

  // 如果超出图片范围，渲染黑边底色
  if (uv.x < 0.0 || uv.x > 1.0 || uv.y < 0.0 || uv.y > 1.0) {
    gl_FragColor = vec4(0.05, 0.05, 0.08, 1.0);
    return;
  }

  // 像素空间坐标计算
  vec2 pixelPos = v_uv * u_resolution;
  vec2 mousePixel = u_mouse * u_resolution;
  vec2 diff = pixelPos - mousePixel;

  // 速度拉伸投影 (向速度反方向形变，形成类似流体水滴的拖尾感)
  float speed = length(u_velocity);
  vec2 dir = speed > 0.001 ? normalize(u_velocity) : vec2(0.0);
  
  // 沿着速度方向进行轴向拉伸
  float distAlong = dot(diff, dir);
  vec2 distPerp = diff - dir * distAlong;
  
  // 拉伸系数：当高速移动时，顺向伸长
  float stretch = clamp(speed * 0.04, 0.0, 0.7);
  vec2 stretchedDiff = distPerp + dir * (distAlong / (1.0 + stretch));
  float rawDist = length(stretchedDiff);

  // 极坐标有机流体波纹 (Wobble & Noise)
  float angle = atan(diff.y, diff.x);
  float wave = sin(angle * 5.0 + u_time * 3.5) * 0.045 + 
               cos(angle * 3.0 - u_time * 2.0) * 0.035;
  
  // 有效半径引入呼吸波纹
  float effectiveRadius = u_radius * (1.0 + wave) * u_hover;

  // 探针边缘柔和羽化遮罩 (Feathered Mask)
  float softness = u_radius * 0.18 + 10.0;
  float mask = smoothstep(effectiveRadius + softness, effectiveRadius - softness, rawDist);

  // 边缘霓虹光环 (Neon Rim Glow)
  // 当距离接近 effectiveRadius 时产生强烈的赛博红光散射
  float rimDist = abs(rawDist - effectiveRadius);
  float rimWidth = u_radius * 0.075 + 4.0;
  float rim = exp(-pow(rimDist / rimWidth, 2.0)) * u_hover;

  // 采样双层纹理
  vec4 colBefore = texture2D(u_texBefore, uv);
  vec4 colAfter = texture2D(u_texAfter, uv);

  // X-Ray 探针内部轻微的科技扫描线 (Subtle Cyber HUD Scanline)
  float scanline = sin(gl_FragCoord.y * 1.5) * 0.035;
  vec3 xRayColor = colAfter.rgb - vec3(scanline);

  // 双层平滑合成
  vec3 blended = mix(colBefore.rgb, xRayColor, mask);

  // 叠加边缘高能光环 (Crimson Neon Rim)
  vec3 rimColor = vec3(1.0, 0.16, 0.16) * rim * 1.9;
  // 核心高光
  vec3 rimCore = vec3(1.0, 0.65, 0.4) * pow(rim, 3.5) * 1.6;

  blended += rimColor + rimCore;

  gl_FragColor = vec4(blended, 1.0);
}
`

function createShader(gl, type, source) {
  const shader = gl.createShader(type)
  gl.shaderSource(shader, source)
  gl.compileShader(shader)
  if (!gl.getShaderParameter(shader, gl.COMPILE_STATUS)) {
    console.error('Shader compile error:', gl.getShaderInfoLog(shader))
    gl.deleteShader(shader)
    return null
  }
  return shader
}

function initGL() {
  const canvas = canvasRef.value
  if (!canvas)
    return false

  gl = canvas.getContext('webgl', {
    antialias: true,
    alpha: false,
    powerPreference: 'high-performance',
  })

  if (!gl) {
    console.error('WebGL not supported')
    return false
  }

  const vs = createShader(gl, gl.VERTEX_SHADER, VS_SOURCE)
  const fs = createShader(gl, gl.FRAGMENT_SHADER, FS_SOURCE)
  if (!vs || !fs)
    return false

  program = gl.createProgram()
  gl.attachShader(program, vs)
  gl.attachShader(program, fs)
  gl.linkProgram(program)

  if (!gl.getProgramParameter(program, gl.LINK_STATUS)) {
    console.error('Program link error:', gl.getProgramInfoLog(program))
    return false
  }

  gl.useProgram(program)

  // 全屏四边形顶点数据
  gl.deleteShader(vs)
  gl.deleteShader(fs)
  positionBuffer = gl.createBuffer()
  gl.bindBuffer(gl.ARRAY_BUFFER, positionBuffer)
  gl.bufferData(
    gl.ARRAY_BUFFER,
    new Float32Array([
      -1,
      -1,
      1,
      -1,
      -1,
      1,
      -1,
      1,
      1,
      -1,
      1,
      1,
    ]),
    gl.STATIC_DRAW,
  )

  const posLoc = gl.getAttribLocation(program, 'a_position')
  gl.enableVertexAttribArray(posLoc)
  gl.vertexAttribPointer(posLoc, 2, gl.FLOAT, false, 0, 0)

  return true
}

function loadTexture(gl, url) {
  return new Promise((resolve, reject) => {
    const texture = gl.createTexture()
    const image = new Image()
    image.crossOrigin = 'anonymous'
    image.onload = () => {
      if (disposed) {
        gl.deleteTexture(texture)
        reject(new Error('Page disposed'))
        return
      }
      gl.bindTexture(gl.TEXTURE_2D, texture)
      gl.pixelStorei(gl.UNPACK_FLIP_Y_WEBGL, false)
      gl.texImage2D(gl.TEXTURE_2D, 0, gl.RGBA, gl.RGBA, gl.UNSIGNED_BYTE, image)

      // 设置滤波
      gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
      gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
      gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
      gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

      resolve(texture)
    }
    image.onerror = (e) => {
      if (url.endsWith('.webp')) {
        const fallbackUrl = url.includes('Before') ? '/images/Before.png' : '/images/After.png'
        image.src = fallbackUrl
        return
      }
      gl.deleteTexture(texture)
      reject(e)
    }
    image.src = url
  })
}

function resizeCanvas() {
  const container = containerRef.value
  const canvas = canvasRef.value
  if (!container || !canvas || !gl)
    return

  const dpr = Math.min(window.devicePixelRatio || 1, 2)
  const width = container.clientWidth
  const height = container.clientHeight

  if (canvas.width !== width * dpr || canvas.height !== height * dpr) {
    canvas.width = width * dpr
    canvas.height = height * dpr
    gl.viewport(0, 0, canvas.width, canvas.height)
  }
}

function render() {
  if (disposed)
    return
  if (!isVisible || document.hidden) {
    animationFrameId = requestAnimationFrame(render)
    return
  }
  if (!gl || !program || !isReady.value) {
    animationFrameId = requestAnimationFrame(render)
    return
  }

  const now = performance.now()
  state.lastFrameTime = now
  const elapsed = reducedMotion ? 0 : (now - state.startTime) / 1000

  const canvas = canvasRef.value
  const width = canvas.width
  const height = canvas.height

  // 1. 物理位置插值 (Lerp)
  state.currentX += (state.targetX - state.currentX) * CONFIG.lerpFactor
  state.currentY += (state.targetY - state.currentY) * CONFIG.lerpFactor

  // 2. 速度与形变计算
  const dx = state.currentX - state.prevX
  const dy = state.currentY - state.prevY
  state.prevX = state.currentX
  state.prevY = state.currentY

  // 速度矢量做平滑衰减
  state.velocityX += (dx - state.velocityX) * 0.25
  state.velocityY += (dy - state.velocityY) * 0.25
  const rawSpeed = Math.sqrt(state.velocityX * state.velocityX + state.velocityY * state.velocityY)
  state.speed += (rawSpeed - state.speed) * 0.2

  // 3. 动态探针半径计算 (运动速度越快，光斑适当扩充)
  const speedExpand = Math.min(state.speed * 3.5, CONFIG.maxStretchRadius - CONFIG.baseRadius)
  const targetRad = CONFIG.baseRadius + speedExpand
  state.currentRadius += (targetRad - state.currentRadius) * 0.15

  // 4. 显隐渐变因子 (Hover / Exit)
  state.hoverFactor += (state.targetHover - state.hoverFactor) * 0.08

  // 5. 闲置微呼吸 (Idle Breathing)
  let activeRadius = state.currentRadius
  if (state.speed < 0.5) {
    activeRadius += Math.sin(elapsed * 2.2) * 8.0
  }

  gl.useProgram(program)

  // 设置 Uniform 变量
  const uRes = gl.getUniformLocation(program, 'u_resolution')
  const uImgRes = gl.getUniformLocation(program, 'u_imageResolution')
  const uImageAlign = gl.getUniformLocation(program, 'u_imageAlign')
  const uMouse = gl.getUniformLocation(program, 'u_mouse')
  const uVel = gl.getUniformLocation(program, 'u_velocity')
  const uRad = gl.getUniformLocation(program, 'u_radius')
  const uTime = gl.getUniformLocation(program, 'u_time')
  const uHov = gl.getUniformLocation(program, 'u_hover')

  gl.uniform2f(uRes, width, height)
  gl.uniform2f(uImgRes, CONFIG.imgWidth, CONFIG.imgHeight)
  gl.uniform1f(uImageAlign, window.innerWidth <= 900 ? 0.5 : 0.0)
  // 传入归一化的物理坐标
  gl.uniform2f(uMouse, state.currentX / width, state.currentY / height)
  gl.uniform2f(uVel, state.velocityX, state.velocityY)
  // 半径按 DPR 适配像素
  const dpr = Math.min(window.devicePixelRatio || 1, 2)
  gl.uniform1f(uRad, activeRadius * dpr)
  gl.uniform1f(uTime, elapsed)
  gl.uniform1f(uHov, state.hoverFactor)

  // 绑定纹理
  const uTexBeforeLoc = gl.getUniformLocation(program, 'u_texBefore')
  const uTexAfterLoc = gl.getUniformLocation(program, 'u_texAfter')

  gl.activeTexture(gl.TEXTURE0)
  gl.bindTexture(gl.TEXTURE_2D, textureBefore)
  gl.uniform1i(uTexBeforeLoc, 0)

  gl.activeTexture(gl.TEXTURE1)
  gl.bindTexture(gl.TEXTURE_2D, textureAfter)
  gl.uniform1i(uTexAfterLoc, 1)

  gl.drawArrays(gl.TRIANGLES, 0, 6)

  if (!reducedMotion)
    animationFrameId = requestAnimationFrame(render)
}

// 鼠标交互事件
function handleMouseMove(e) {
  if (reducedMotion)
    return
  const container = containerRef.value
  if (!container)
    return

  const rect = container.getBoundingClientRect()
  const dpr = Math.min(window.devicePixelRatio || 1, 2)
  const clientX = e.clientX - rect.left
  const clientY = e.clientY - rect.top

  state.targetX = clientX * dpr
  state.targetY = clientY * dpr
  state.targetHover = 1.0

  if (!state.hasInitialMoved) {
    state.hasInitialMoved = true
    state.currentX = state.targetX
    state.currentY = state.targetY
    state.prevX = state.targetX
    state.prevY = state.targetY
  }
}

function handleMouseEnter() {
  state.targetHover = 1.0
}

function handleMouseLeave() {
  // 鼠标移出时，光圈优雅淡出
  state.targetHover = 0.0
}

function handleTouchStart(e) {
  if (e.touches && e.touches.length > 0)
    handleMouseMove(e.touches[0])
}

function handleTouchMove(e) {
  if (e.touches && e.touches.length > 0)
    handleMouseMove(e.touches[0])
}

let resizeObserver = null

async function startEnhancement() {
  if (hasStarted.value || disposed || reducedMotion)
    return
  hasStarted.value = true
  if (!initGL()) {
    hasError.value = true
    return
  }
  resizeCanvas()

  // 初始将探针预置在右侧人物肩膀/背部核心区域，呼吸待命
  const canvas = canvasRef.value
  state.targetX = canvas.width * 0.58
  state.targetY = canvas.height * 0.46
  state.currentX = state.targetX
  state.currentY = state.targetY
  state.prevX = state.targetX
  state.prevY = state.targetY
  state.targetHover = 0.85 // 页面加载后默认显露一处 X-Ray 唤起好奇心
  if (reducedMotion)
    state.hoverFactor = 0.85

  // 监听容器尺寸调整
  resizeObserver = new ResizeObserver(() => {
    resizeCanvas()
    if (reducedMotion && isReady.value)
      render()
  })
  if (containerRef.value) {
    resizeObserver.observe(containerRef.value)
  }

  // 加载双图纹理
  try {
    const [tBefore, tAfter] = await Promise.all([
      loadTexture(gl, beforeURL.value),
      loadTexture(gl, afterURL.value),
    ])
    if (disposed) {
      gl.deleteTexture(tBefore)
      gl.deleteTexture(tAfter)
      return
    }
    textureBefore = tBefore
    textureAfter = tAfter
    isReady.value = true
  }
  catch (err) {
    if (disposed)
      return
    hasError.value = true
    console.error('Failed to load X-Ray textures:', err)
    return
  }

  animationFrameId = requestAnimationFrame(render)
}

onMounted(() => {
  // 默认启动 X-Ray 效果
  startEnhancement()
  intersectionObserver = new IntersectionObserver(([entry]) => {
    isVisible = entry.isIntersecting
    if (entry.isIntersecting && !hasStarted.value) {
      startEnhancement()
    }
  })
  if (containerRef.value) {
    intersectionObserver.observe(containerRef.value)
  }
})

onUnmounted(() => {
  disposed = true
  intersectionObserver?.disconnect()
  if (animationFrameId) {
    cancelAnimationFrame(animationFrameId)
  }
  if (resizeObserver) {
    resizeObserver.disconnect()
  }
  if (gl) {
    if (textureBefore)
      gl.deleteTexture(textureBefore)
    if (textureAfter)
      gl.deleteTexture(textureAfter)
    if (program)
      gl.deleteProgram(program)
    if (positionBuffer)
      gl.deleteBuffer(positionBuffer)
    gl.getExtension('WEBGL_lose_context')?.loseContext()
  }
})
</script>

<template>
  <div
    ref="containerRef"
    class="xray-canvas-container"
    @mousemove="handleMouseMove"
    @mouseenter="handleMouseEnter"
    @mouseleave="handleMouseLeave"
    @touchstart.passive="handleTouchStart"
    @touchmove.passive="handleTouchMove"
  >
    <picture v-if="!isReady || hasError" class="xray-fallback">
      <source v-if="beforeSrcset" :srcset="beforeSrcset" sizes="100vw" type="image/webp">
      <img :src="beforeURL" alt="博客封面" width="1672" height="941" fetchpriority="high">
    </picture>
    <canvas ref="canvasRef" class="xray-canvas" :class="{ 'is-unavailable': hasError }" aria-label="随鼠标移动显示透视光效的博客封面" />

    <!-- 加载中指示器 (初始资源较大时优雅过渡) -->
    <Transition name="fade">
      <div v-if="hasStarted && !isReady && !hasError" class="loading-overlay">
        <div class="cyber-spinner" />
        <span class="loading-text">封面加载中…</span>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.xray-canvas-container {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  background: var(--bg-stage);
  cursor: crosshair;
  z-index: 1;
}

.xray-fallback { position: absolute; inset: 0; display: block; }
.xray-fallback img { width: 100%; height: 100%; display: block; object-fit: cover; object-position: left center; }

@media (max-width: 900px) {
  .xray-fallback img { object-position: center; }
}

.xray-canvas {
  position: relative;
  width: 100%;
  height: 100%;
  display: block;
}
.xray-canvas.is-unavailable { visibility: hidden; }

.loading-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  justify-content: flex-end;
  padding: 32px;
  background: transparent;
  pointer-events: none;
  z-index: 10;
  gap: 16px;
}

.cyber-spinner {
  width: 46px;
  height: 46px;
  border: 2px solid rgba(255, 50, 50, 0.2);
  border-top-color: #ff3333;
  border-radius: 50%;
  animation: cyber-spin 0.9s linear infinite;
  box-shadow: 0 0 15px rgba(255, 50, 50, 0.4);
}

.loading-text {
  font-family: monospace;
  font-size: 0.85rem;
  letter-spacing: 0.2em;
  color: rgba(255, 255, 255, 0.7);
}

@keyframes cyber-spin {
  to {
    transform: rotate(360deg);
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.6s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
