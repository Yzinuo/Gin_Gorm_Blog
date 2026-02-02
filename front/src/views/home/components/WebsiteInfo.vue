<script setup>
import { onMounted, onUnmounted, ref } from 'vue'
import { storeToRefs } from 'pinia'

import dayjs from 'dayjs'

import { useAppStore } from '@/store'

const { blogConfig, viewCount } = storeToRefs(useAppStore())

// 每秒刷新时间
const runTime = ref('')

const formatRuntime = (seconds) => {
  if (!Number.isFinite(seconds) || seconds < 0) {
    return '0 天 0 时 0 分'
  }
  const days = Math.floor(seconds / 86400)
  const hours = Math.floor((seconds % 86400) / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  return `${days} 天 ${hours} 时 ${minutes} 分`
}

const getRuntimeSeconds = () => {
  const config = blogConfig.value || {}
  if (config.website_createtime_unix) {
    const unix = Number.parseInt(config.website_createtime_unix, 10)
    if (Number.isFinite(unix)) {
      return Math.max(0, Math.floor((Date.now() - unix * 1000) / 1000))
    }
  }

  const createtime = config.website_createtime_rfc3339 || config.website_createtime
  if (createtime) {
    const createTime = dayjs(createtime)
    if (createTime.isValid()) {
      return Math.max(0, dayjs().diff(createTime, 'second'))
    }
  }
  return 0
}

const refreshRuntime = () => {
  runTime.value = formatRuntime(getRuntimeSeconds())
}

// 每 30 秒刷新当前时间
const timer = setInterval(refreshRuntime, 30 * 1000)

onMounted(() => {
  refreshRuntime()
})

onUnmounted(() => {
  clearInterval(timer)
})
</script>

<template>
  <div class="card-view hidden animate-zoom-in animate-duration-600 lg:block space-y-2">
    <p class="flex items-center text-lg">
      <span class="i-icon-park:analysis mr-1.5" />
      <span> 网站咨询 </span>
    </p>
    <div class="space-y-1">
      <p>
        <span> 运行时间： </span>
        <span class="float-right"> {{ runTime }} </span>
      </p>
      <p>
        <span> 总访问量： </span>
        <span class="float-right"> {{ viewCount }} </span>
      </p>
    </div>
  </div>
</template>
