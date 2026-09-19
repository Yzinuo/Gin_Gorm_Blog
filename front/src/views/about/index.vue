<script setup>
import { computed, nextTick, onMounted, onUnmounted, ref } from 'vue'
import { marked } from 'marked'
import hljs from 'highlight.js/lib/core'
import 'highlight.js/styles/a11y-dark.css'
import go from 'highlight.js/lib/languages/go'
import json from 'highlight.js/lib/languages/json'
import javascript from 'highlight.js/lib/languages/javascript'
import bash from 'highlight.js/lib/languages/bash'

import ResumeScene from './components/ResumeScene.vue'
import cameraMap from './camera-map.json'
import { content } from './content'
import api from '@/api'
import { sanitizeHtml } from '@/utils/sanitize'
import { enhanceCodeBlocks } from '@/utils/codeBlock'

hljs.registerLanguage('go', go)
hljs.registerLanguage('bash', bash)
hljs.registerLanguage('json', json)
hljs.registerLanguage('javascript', javascript)

const story = ref(null)
const scene = ref(null)
const chapter = ref(0)
const html = ref('')
const textStatus = ref('loading')
let disposed = false
const profile = ref([])
const entries = computed(() => cameraMap.stickers.map((sticker, index) => {
  const saved = profile.value.find(item => item.id === index + 1)
  return { ...sticker, ...content[index], ...saved, image: saved?.image || '', originalImage: sticker.image }
}))
const active = computed(() => entries.value[chapter.value - 1])
const caption = computed(() => active.value?.title || (chapter.value === 9 ? '下一章，待续。' : '每段经历，\n都有自己的坐标。'))
const chapterLabel = computed(() => active.value?.category || (chapter.value === 9 ? 'THE NEXT CHAPTER' : 'INTRODUCTION'))
const base = import.meta.env.BASE_URL
const resumeStatus = ref('loading')
const managedAssets = ref({})
const modelURL = computed(() => managedAssets.value['resume.model']?.model?.src || `${base}resume/resume-ready.optimized.glb`)

function stickerURL(entry) {
  const managed = managedSticker(entry)
  if (managed)
    return managed
  if (!entry.image)
    return `${base}resume/${entry.originalImage}`
  if (/^https?:\/\//i.test(entry.image))
    return entry.image
  return `${(import.meta.env.VITE_BACKEND_URL || '').replace(/\/$/, '')}/${entry.image.replace(/^\//, '')}`
}

function managedSticker(entry) {
  return managedAssets.value[`resume.sticker.${String(entry.index).padStart(2, '0')}`]?.image?.src || ''
}

const stickers = computed(() => entries.value.map(entry => ({ object: entry.object, image: managedSticker(entry) || entry.image ? stickerURL(entry) : '' })))

function fallbackSticker(event, entry) {
  const fallback = new URL(`${base}resume/${entry.originalImage}`, window.location.origin).href
  if (event.target.src !== fallback)
    event.target.src = fallback
}

async function loadResume() {
  resumeStatus.value = 'loading'
  try {
    const { data } = await api.getResume()
    if (!Array.isArray(data?.entries) || data.entries.length !== 8)
      throw new Error('Invalid résumé configuration')
    if (disposed)
      return
    profile.value = data.entries
    resumeStatus.value = 'ready'
  }
  catch {
    if (!disposed)
      resumeStatus.value = 'error'
  }
}

async function loadAssets() {
  try {
    const { data } = await api.getAssets('resume')
    if (!disposed)
      managedAssets.value = data?.assets || {}
  }
  catch {
    // Keep repository/profile fallbacks readable when manifest or R2 is down.
  }
}

function explore(id) {
  const element = story.value?.querySelector(`#${id}`) || story.value?.querySelector('#next-chapter')
  if (!element)
    return
  const bounds = element.getBoundingClientRect()
  const stage = document.querySelector('.resume-stage').getBoundingClientRect()
  const focus = innerWidth <= 640 ? stage.bottom + (innerHeight - stage.bottom) / 2 : innerHeight / 2
  window.scrollTo({ top: window.scrollY + bounds.top + bounds.height / 2 - focus, behavior: matchMedia('(prefers-reduced-motion: reduce)').matches ? 'instant' : 'smooth' })
}

async function loadIntroduction() {
  textStatus.value = 'loading'
  try {
    const { data } = await api.about()
    if (disposed)
      return
    html.value = sanitizeHtml(await marked.parse(typeof data === 'string' ? data : '', { async: true }))
    textStatus.value = 'ready'
    await nextTick()
    if (disposed)
      return
    story.value?.querySelectorAll('.about-copy pre code').forEach(element => hljs.highlightElement(element))
    enhanceCodeBlocks(story.value?.querySelector('.about-copy'))
    window.MathJax?.typeset?.()
  }
  catch {
    if (!disposed)
      textStatus.value = 'error'
  }
}
onMounted(loadIntroduction)
onMounted(loadResume)
onMounted(loadAssets)
onUnmounted(() => disposed = true)
</script>

<template>
  <main class="personal-atlas">
    <ResumeScene ref="scene" :story="story" :stickers="stickers" :model-url="modelURL" @chapter="chapter = $event">
      <div class="stage-caption" :class="{ compact: chapter !== 0 }">
        <span>{{ String(chapter).padStart(2, '0') }} / {{ chapterLabel }}</span>
        <h2>{{ caption }}</h2>
        <p>{{ chapter === 0 ? '向下滚动，走近我的世界。' : '一段经历，一个坐标。' }}</p>
      </div>
      <div class="stage-bottom">
        <span class="status-dot" /> <span>MOVE YOUR CURSOR · 眼神随你</span>
      </div>
    </ResumeScene>
    <RouterLink to="/" class="atlas-wordmark" aria-label="Zane，返回博客首页">
      ZANE<span>PERSONAL ATLAS / 个人图志</span>
    </RouterLink>
    <nav class="atlas-nav" aria-label="自我介绍导航">
      <RouterLink to="/">
        返回博客 ↗
      </RouterLink>
      <a href="#entry-1" @click.prevent="explore('entry-1')">经历</a>
      <a href="#more-about" @click.prevent="explore('more-about')">更多关于我</a>
    </nav>
    <div ref="story" class="atlas-story">
      <section id="introduction" class="story-section introduction" data-frame="0">
        <span class="atlas-eyebrow">A LIVING RÉSUMÉ / 自我介绍</span>
        <h1>不止一页<br>个人履历。</h1>
        <p>从研究到实践，从一行代码到一个作品。<br>把走过的路，留在这个小小的世界里。</p>
        <div class="intro-rule" />
        <p class="intro-tags">
          研究 · 工程 · 开源
        </p>
        <a href="#entry-1" class="atlas-scroll-link" @click.prevent="explore('entry-1')">开始探索 <span aria-hidden="true">↓</span></a>
        <small>08 EXPERIENCES / 一段持续生长的旅程</small>
        <p v-if="resumeStatus === 'error'" role="status">
          最新经历暂未加载，正在展示默认介绍。
          <button type="button" class="entry-focus" @click="loadResume">
            重新加载
          </button>
        </p>
      </section>
      <section v-for="(entry, index) in entries" :id="`entry-${entry.index}`" :key="entry.index" class="story-section story-entry" :class="{ active: chapter === index + 1 }" :data-frame="entry.frame">
        <div class="entry-meta">
          <span class="entry-number">{{ String(entry.index).padStart(2, '0') }}</span><span>{{ entry.category }}</span>
        </div>
        <img class="sticker-image" :src="stickerURL(entry)" :alt="`${entry.title}贴纸`" width="180" height="140" loading="lazy" @error="fallbackSticker($event, entry)">
        <h2>{{ entry.title }}</h2>
        <p>{{ entry.description }}</p>
        <ul class="entry-tags" aria-label="经历标签">
          <li v-for="tag in entry.tags" :key="tag">
            {{ tag }}
          </li>
        </ul>
        <button type="button" class="entry-focus" @click="scene?.focus(entry.frame)">
          聚焦这个坐标 <span aria-hidden="true">↗</span>
        </button>
      </section>
      <section id="next-chapter" class="story-section next-chapter" data-frame="500">
        <span class="atlas-eyebrow">THE NEXT CHAPTER</span>
        <h2>继续探索，<br>继续创造。</h2>
        <p>新的想法，正在成为下一个作品。</p>
        <RouterLink to="/" class="atlas-return">
          回到博客，阅读我的记录 <span aria-hidden="true">↗</span>
        </RouterLink>
        <a href="#introduction" class="entry-focus" @click.prevent="explore('introduction')">回到起点 ↑</a>
      </section>
      <section v-if="html || textStatus !== 'ready'" id="more-about" class="about-copy" aria-label="更多关于我">
        <h2>更多关于我</h2>
        <p v-if="textStatus === 'loading'" role="status">
          正在加载介绍…
        </p>
        <div v-else-if="textStatus === 'error'" class="text-muted" role="status">
          文字介绍暂时未能加载。<button type="button" class="entry-focus" @click="loadIntroduction">
            重新加载
          </button>
        </div>
        <article v-else class="max-w-none prose prose-truegray" v-html="html" />
      </section>
      <footer class="atlas-footer">
        © {{ new Date().getFullYear() }} Zane · 研究 / 工程 / 开源
      </footer>
    </div>
  </main>
</template>

<style scoped>
.personal-atlas {
  background: var(--bg-page);
  color: var(--text-primary);
}
.atlas-wordmark {
  position: fixed;
  top: 30px;
  left: 38px;
  z-index: 8;
  display: flex;
  align-items: center;
  gap: 22px;
  font: 31px Georgia, serif;
  letter-spacing: 3px;
}
.atlas-wordmark span {
  max-width: 140px;
  font: 10px/1.8 Georgia, serif;
  letter-spacing: 2px;
  color: #e1d3cf;
}
.atlas-nav {
  position: fixed;
  top: 0;
  right: 0;
  width: 36%;
  z-index: 8;
  min-height: 80px;
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 26px;
  padding: 18px 38px;
  background: linear-gradient(var(--bg-page) 70%, transparent);
  font-size: 12px;
}
.atlas-nav a:hover, .entry-focus:hover {
  color: var(--brand);
}
.stage-caption {
  position: absolute;
  left: 38px;
  right: 28px;
  bottom: 96px;
  pointer-events: none;
  text-shadow: 0 2px 24px #190e1bcc;
}
.stage-caption > span {
  font-size: 10px;
  letter-spacing: 2.5px;
  color: #eddbd5;
}
.stage-caption h2 {
  white-space: pre-line;
  font: 400 clamp(26px, 3vw, 48px)/1.4 Georgia, 'Songti SC', SimSun, serif;
  letter-spacing: 2px;
  margin: 15px 0 12px;
}
.stage-caption p {
  font-size: 12px;
  letter-spacing: 1px;
}
.stage-bottom {
  position: absolute;
  bottom: 25px;
  left: 38px;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 10px;
  letter-spacing: 1px;
  text-shadow: 0 1px 8px #000;
}
.status-dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: var(--brand);
}
.atlas-story {
  position: relative;
  margin-left: 64%;
  background: var(--bg-page);
}
.story-section {
  min-height: 100svh;
  padding: 120px 48px 75px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  position: relative;
  border-bottom: 1px solid var(--border-color);
}
.atlas-eyebrow {
  font-size: 10px;
  letter-spacing: 2.5px;
  color: var(--text-muted);
}
.story-section h1, .story-section h2 {
  font: 400 clamp(28px, 3.1vw, 46px)/1.55 Georgia, 'Songti SC', SimSun, serif;
  letter-spacing: 1px;
  margin: 24px 0;
}
.story-section p {
  white-space: pre-wrap;
  overflow-wrap: anywhere;
  font-size: 14px;
  line-height: 2;
  color: var(--text-muted);
}
.intro-rule {
  height: 1px;
  width: 36px;
  background: var(--brand);
  margin: 35px 0 22px;
}
.intro-tags {
  letter-spacing: 3px;
  font-size: 12px !important;
}
.atlas-scroll-link {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid var(--border-color);
  margin-top: 55px;
  padding-bottom: 17px;
  font-size: 13px;
}
.atlas-scroll-link span {
  font-size: 23px;
  color: var(--brand);
}
.story-section small {
  font-size: 10px;
  line-height: 1.9;
  color: var(--text-muted);
  margin-top: 40px;
  letter-spacing: 1px;
}
.entry-number {
  font: italic 66px/1 Georgia, serif;
  color: var(--text-muted);
}
.entry-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 15px;
  margin-bottom: 28px;
}
.entry-meta > span:last-child {
  font-size: 10px;
  letter-spacing: 2px;
}
.sticker-image {
  width: 180px;
  height: 140px;
  object-fit: contain;
  object-position: left center;
}
.story-entry h2 {
  font-size: 30px;
  margin: 15px 0;
}
.entry-tags {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin: 24px 0;
  padding: 0;
  list-style: none;
}
.entry-tags li {
  border: 1px solid var(--border-color);
  border-radius: 20px;
  padding: 6px 12px;
  font-size: 12px;
  color: var(--text-muted);
}
.entry-focus {
  border: 0;
  border-bottom: 1px solid var(--border-color);
  align-self: flex-start;
  background: none;
  padding: 12px 0;
  font-size: 12px;
  color: var(--text-muted);
}
.story-entry::before {
  content: '';
  position: absolute;
  width: 5px;
  height: 5px;
  border-radius: 50%;
  left: 21px;
  top: 50%;
  background: var(--border-color);
}
.story-entry.active::before {
  background: var(--brand);
  box-shadow: 0 0 0 5px var(--brand-soft);
}
.atlas-return {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  background: var(--brand);
  color: var(--on-brand);
  padding: 18px 20px;
  font-size: 13px;
  margin: 34px 0 18px;
}
.about-copy {
  padding: 56px 48px;
  border-bottom: 1px solid var(--border-color);
  overflow-wrap: anywhere;
}
.about-copy > h2 {
  font: 28px Georgia, 'Songti SC', SimSun, serif;
  margin-bottom: 28px;
}
.atlas-footer {
  padding: 32px;
  font-size: 11px;
  text-align: center;
  color: var(--text-muted);
}
@media (max-width: 900px) {
  .atlas-nav {
    width: 42%;
    padding: 18px 25px;
    gap: 18px;
  }
  .atlas-story {
    margin-left: 58%;
  }
  .story-section {
    padding: 110px 30px 65px;
  }
  .atlas-wordmark {
    left: 25px;
    top: 25px;
  }
  .atlas-wordmark span {
    display: none;
  }
  .stage-caption {
    left: 25px;
  }
  .about-copy {
    padding: 45px 30px;
  }
}
@media (max-width: 640px) {
  .atlas-wordmark {
    font-size: 21px;
    top: 20px;
    left: 20px;
  }
  .atlas-nav {
    width: auto;
    min-height: 60px;
    background: none;
    padding: 12px 18px;
    gap: 16px;
    font-size: 11px;
  }
  .atlas-story {
    margin-left: 0;
    padding-top: 54svh;
  }
  .story-section {
    min-height: 70svh;
    padding: 50px 32px 40px;
    scroll-margin-top: 54svh;
  }
  .story-section h1, .story-section h2 {
    font-size: 30px;
  }
  .stage-caption {
    left: 20px;
    bottom: 45px;
  }
  .stage-caption h2 {
    font-size: 24px;
    margin: 9px 0 0;
  }
  .stage-caption.compact h2 {
    font-size: 21px;
  }
  .stage-caption p {
    display: none;
  }
  .stage-caption > span {
    font-size: 9px;
  }
  .stage-bottom {
    bottom: 12px;
    left: 20px;
    font-size: 9px;
  }
  .intro-rule {
    margin: 22px 0;
  }
  .atlas-scroll-link {
    margin-top: 25px;
  }
  .entry-number {
    font-size: 44px;
  }
  .sticker-image {
    width: 140px;
    height: 100px;
  }
  .entry-meta {
    margin-bottom: 18px;
  }
  .story-entry::before {
    left: 15px;
  }
  .story-entry h2 {
    font-size: 27px;
  }
  .story-section small {
    margin-top: 24px;
  }
}
</style>
