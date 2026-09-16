<script setup>
import { onMounted, ref } from 'vue'
import BannerPage from '@/components/BannerPage.vue'
import api from '@/api'

const loading = ref(true)
const tagList = ref([])

onMounted(() => {
  api.getTags().then((resp) => {
    tagList.value = resp.data || []
    loading.value = false
  })
})

const tagColors = ['var(--brand)', 'var(--text-primary)', 'var(--text-muted)']
</script>

<template>
  <BannerPage :loading="loading" title="标签" label="tag" card>
    <h2 class="text-center text-2xl leading-8 lg:text-3xl">
      标签 - {{ tagList.length }}
    </h2>
    <div class="mt-6 text-center">
      <RouterLink
        v-for="(t, index) of tagList" :key="t.id" :to="`tags/${t.id}?name=${t.name}`"
        :style="{
          'font-size': `${18 + (index % 3) * 4}px`,
          'color': tagColors[index % tagColors.length],
        }"
        class="inline-block px-2 leading-11 transition-300 hover:scale-110 !hover:text-brand-hover"
      >
        {{ t.name }}
      </RouterLink>
    </div>
  </BannerPage>
</template>

<style scoped>
/* 实现截断文字效果, 即不会在结束处将一个词语拆开 */
a {
  display: inline-block;
}
</style>
