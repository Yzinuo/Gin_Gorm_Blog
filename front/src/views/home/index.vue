<script setup>
import { onMounted, reactive, ref } from 'vue'

// 无限轮播图
import InfiniteLoading from 'v3-infinite-loading'

// Markdown => Html
import { marked } from 'marked'

import ArticleCard from './components/ArticleCard.vue'
import AuthorInfo from './components/AuthorInfo.vue'
import WebsiteInfo from './components/WebsiteInfo.vue'
import HomeBanner from './components/HomeBanner.vue'
import Announcement from './components/Announcement.vue'
import TalkingCarousel from './components/TalkingCarousel.vue'
import AppFooter from '@/components/layout/AppFooter.vue'
import api from '@/api'

const articleList = ref([])
const loading = ref(true)
const loadError = ref(false)

// 无限加载文章
const params = reactive({ page_size: 5, page_num: 1 }) // 列表加载参数
async function getArticlesInfinite($state) {
  if (!loading.value) {
    try {
      const resp = await api.getArticles(params)
      // 加载完成
      if (!resp.data.length) {
        $state.complete()
        return
      }
      // 非首次加载, 都是往列表中添加数据
      articleList.value.push(...resp.data)
      // 过滤 Markdown 符号
      articleList.value.forEach(e => e.content = filterMdSymbol(e.content))
      params.page_num++
      $state.loaded()
    }
    catch (error) {
      $state.error()
    }
  }
}

async function loadInitialArticles() {
  loading.value = true
  loadError.value = false
  try {
    const resp = await api.getArticles(params)
    articleList.value = resp.data
    articleList.value.forEach(e => e.content = filterMdSymbol(e.content))
    params.page_num++
  }
  catch {
    loadError.value = true
  }
  finally {
    loading.value = false
  }
}
onMounted(loadInitialArticles)

// 卡片只展示纯文本摘要，保留文章详情页的 Markdown 渲染。
function filterMdSymbol(md) {
  const template = document.createElement('template')
  template.innerHTML = marked(md || '')
  template.content.querySelectorAll('pre, script, style, table').forEach(element => element.remove())
  return (template.content.textContent || '').replace(/\s+/g, ' ').trim()
}

function backTop() {
  window.scrollTo({ behavior: 'smooth', top: 0 })
}
</script>

<template>
  <div class="home-page-wrapper">
    <!-- 首页全屏纯暗黑赛博武士封面 (100svh 填满首屏，纯净暗黑绝无杂色) -->
    <HomeBanner />

    <!-- 下方阅读区域：平滑渐变到暖纸羊皮纸质感 -->
    <div class="home-reading-stage">
      <div class="home-transition-veil" />
      <!-- 内容 -->
      <main id="articles" class="home-articles mx-auto mb-8 max-w-[1230px] w-full flex flex-col justify-center px-3">
        <div class="home-section-heading">
          <h2>最近的记录</h2>
          <RouterLink to="/about">
            认识我 · 3D 履历 <span aria-hidden="true">↗</span>
          </RouterLink>
        </div>
        <div class="grid grid-cols-12 gap-4">
          <!-- 左半部分 -->
          <div class="col-span-12 lg:col-span-9 space-y-5">
            <!-- 说说轮播 -->
            <TalkingCarousel />
            <p v-if="loading" class="card-view text-muted" role="status">
              正在加载文章…
            </p>
            <div v-else-if="loadError" class="card-view flex flex-wrap items-center justify-between gap-4" role="status">
              <p class="text-muted">
                文章暂时未能加载，你可以先浏览我的自我介绍。
              </p>
              <button type="button" class="the-button" @click="loadInitialArticles">
                重新加载文章
              </button>
            </div>
            <!-- 文章列表 -->
            <div class="space-y-5">
              <ArticleCard v-for="(item, idx) in articleList" :key="item.id" :article="item" :idx="idx" />
            </div>
            <!-- 无限加载 -->
            <div class="f-c-c">
              <InfiniteLoading v-if="!loading && !loadError" class="mt-2 lg:mt-5" @infinite="getArticlesInfinite">
                <!-- TODO: 优化界面 -->
                <template #spinner>
                  <span class="animate-pulse text-xl">
                    loading...
                  </span>
                </template>
                <template #complete>
                  <span class="flex gap-2 text-gray">
                    没有更多文章啦!
                    <button class="flex items-center" @click="backTop">
                      点击回到顶部 <span class="i-mdi:arrow-up-bold-box ml-1 inline-block text-xl" />
                    </button>
                  </span>
                </template>
              </InfiniteLoading>
            </div>
          </div>
          <!-- 右半部分 -->
          <div class="col-span-0 lg:col-span-3">
            <!-- sticky 实现悬浮固定效果 -->
            <div class="sticky top-5 space-y-5">
              <!-- 博主信息 -->
              <AuthorInfo />
              <!-- 公告 -->
              <Announcement />
              <!-- 网站资讯 -->
              <WebsiteInfo />
            </div>
          </div>
        </div>
      </main>
      <!-- 底部 -->
      <AppFooter />
    </div>
  </div>
</template>

<style scoped>
.home-page-wrapper {
  background-color: #170f14;
  min-height: 100vh;
}
.home-reading-stage {
  position: relative;
  background-color: var(--bg-page);
}
.home-transition-veil {
  height: 180px;
  width: 100%;
  margin-top: -1px;
  background: linear-gradient(
    180deg,
    #170f14 0%,
    rgba(23, 15, 20, 0.95) 20%,
    rgba(145, 120, 115, 0.28) 55%,
    rgba(247, 244, 238, 0.75) 85%,
    var(--bg-page) 100%
  );
  pointer-events: none;
}
.home-articles { padding-top: 10px; scroll-margin-top: 76px; }
.home-section-heading { display: flex; flex-wrap: wrap; align-items: center; justify-content: space-between; gap: 16px; margin-bottom: 24px; }
.home-section-heading h2 { font: 28px Georgia, 'Songti SC', SimSun, serif; }
.home-section-heading a { color: var(--brand); font-size: 13px; }
</style>
