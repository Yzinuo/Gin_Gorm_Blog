<script setup>
import { storeToRefs } from 'pinia'
import { useAppStore } from '@/store'

const { blogInfo } = storeToRefs(useAppStore())
</script>

<template>
  <footer class="f-c-c">
    <div class="footer-wrap w-full px-5 py-10 text-center text-white leading-8">
      <p> &copy;{{ 2022 }} - {{ new Date().getFullYear() }} By Zane</p>
      <a href="https://beian.miit.gov.cn/" target="_blank">
        {{ blogInfo.blog_config?.website_record }}
      </a>
    </div>
  </footer>
</template>

<style scoped>
.footer-wrap {
  background: linear-gradient(145deg, var(--bg-stage), #301d24);
  border-top: 1px solid var(--border-color);
  color: var(--text-muted);
}

@keyframes Gradient {
  0% {
    background-position: 0 50%;
  }
  50% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0 50%;
  }
}
</style>
