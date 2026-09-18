<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { onBeforeRouteLeave } from 'vue-router'
import { NAlert, NButton, NDynamicTags, NFormItem, NInput, NSpin } from 'naive-ui'
import api from '@/api'
import { useAuthStore } from '@/store'

const placements = ['额头中央', '额头左侧', '额头右侧', '画面左上脸颊', '画面左下脸颊', '画面右上脸颊', '画面右下脸颊', '下巴中央']
const originals = ['中科院', '草花', '指色', 'MIT6.824', '国家励志', '字节青训营', '本科SCI', 'Github']
const entries = ref([])
const selected = ref(0)
const loading = ref(true)
const saving = ref(false)
const uploading = ref(false)
const loadError = ref(false)
const baseline = ref('')
const fileInput = ref(null)
const imageError = ref('')
const auth = useAuthStore()
const entry = computed(() => entries.value[selected.value])
const dirty = computed(() => baseline.value && JSON.stringify(entries.value) !== baseline.value)
let disposed = false
const controller = new AbortController()

function imageURL(value) {
  if (!value)
    return `${import.meta.env.VITE_BLOG_URL || (import.meta.env.DEV ? 'http://localhost:3333' : window.location.origin)}/resume/stickers/${encodeURIComponent(originals[selected.value])}.png`
  if (/^https?:\/\//i.test(value))
    return value
  return `${import.meta.env.VITE_SERVER_URL.replace(/\/$/, '')}/${value.replace(/^\//, '')}`
}

async function load() {
  loading.value = true
  loadError.value = false
  try {
    const { data } = await api.getResume()
    if (!data?.entries || data.entries.length !== 8)
      throw new Error('Invalid profile')
    if (disposed)
      return
    entries.value = data.entries.sort((a, b) => a.id - b.id)
    baseline.value = JSON.stringify(entries.value)
  }
  catch { loadError.value = true }
  finally { loading.value = false }
}

async function upload(event) {
  const file = event.target.files?.[0]
  event.target.value = ''
  if (!file)
    return
  const target = entry.value
  const uploadController = new AbortController()
  const abort = () => uploadController.abort()
  controller.signal.addEventListener('abort', abort, { once: true })
  const timer = setTimeout(abort, 30000)
  uploading.value = true
  imageError.value = ''
  try {
    if (!['image/png', 'image/jpeg', 'image/webp'].includes(file.type) || file.size > 5 * 1024 * 1024)
      throw new Error('请选择不超过 5 MB 的 PNG、JPEG 或 WebP 图片。')
    const bitmap = await createImageBitmap(file)
    const valid = bitmap.width <= 4096 && bitmap.height <= 4096
    bitmap.close()
    if (!valid)
      throw new Error('图片宽高请控制在 4096 像素以内。')
    const body = new FormData()
    body.append('file', file)
    const response = await fetch(`${import.meta.env.VITE_BASE_API || '/api'}/upload`, {
      method: 'POST', headers: { Authorization: `Bearer ${auth.token}` }, body, signal: uploadController.signal,
    })
    if (!response.ok)
      throw new Error('上传失败，请重试。')
    const result = await response.json()
    if (result.code !== 0 || typeof result.data !== 'string')
      throw new Error(result.message || '上传失败，请重试。')
    // Confirm the uploaded asset is readable as a 3D texture (including CORS).
    const image = await fetch(imageURL(result.data), { signal: uploadController.signal })
    if (!image.ok)
      throw new Error('图片已上传但无法读取，请检查图片服务。')
    const decoded = await createImageBitmap(await image.blob())
    decoded.close()
    if (!disposed)
      target.image = result.data
  }
  catch (error) {
    if (!disposed)
      imageError.value = error.name === 'AbortError' ? '上传或图片读取超时，请重试。' : error.message || '无法读取图片，请重试。'
  }
  finally {
    clearTimeout(timer)
    controller.signal.removeEventListener('abort', abort)
    uploading.value = false
  }
}

async function save() {
  for (const item of entries.value) {
    if (!item.title.trim() || !item.category.trim() || !item.description.trim()
      || item.title.length > 60 || item.category.length > 60 || item.description.length > 1200
      || item.tags.length > 5 || item.tags.some(tag => !tag.trim() || tag.length > 24)) {
      selected.value = item.id - 1
      window.$message.error('请补全标题、英文分类和介绍；标签最多 5 个，每个不超过 24 字。')
      return
    }
  }
  saving.value = true
  try {
    const result = await api.updateResume({ entries: entries.value })
    if (result?.code !== 0)
      throw new Error('Save failed')
    baseline.value = JSON.stringify(entries.value)
    window.$message.success('已保存，前台刷新后即可看到更新')
  }
  catch { window.$message.error('保存失败，修改已保留，请重试。') }
  finally { saving.value = false }
}

function beforeUnload(event) {
  if (dirty.value || uploading.value) {
    event.preventDefault()
    event.returnValue = ''
  }
}
onBeforeRouteLeave(() => {
  if (!dirty.value && !uploading.value)
    return true
  return new Promise((resolve) => {
    window.$dialog.warning({
      title: '离开自我介绍编辑？',
      content: '有尚未保存的修改，离开后会丢失。',
      positiveText: '离开',
      negativeText: '继续编辑',
      onPositiveClick: () => resolve(true),
      onNegativeClick: () => resolve(false),
      onClose: () => resolve(false),
      onMaskClick: () => resolve(false),
    })
  })
onMounted(() => {
  load()
  window.addEventListener('beforeunload', beforeUnload)
})
onBeforeUnmount(() => {
  disposed = true
  controller.abort()
  window.removeEventListener('beforeunload', beforeUnload)
})
</script>

<template>
  <NSpin :show="loading">
    <div v-if="loadError" role="alert">
      成就配置加载失败，请重试。
      <NButton @click="load">重新加载</NButton>
    </div>
    <template v-else-if="entry">
      <div class="resume-toolbar">
        <p>选择一个位置，更换贴纸或修改介绍。保存后会同步更新前台卡片和人脸贴纸。</p>
        <NButton type="primary" :loading="saving" :disabled="uploading || !dirty" @click="save">保存全部成就</NButton>
      </div>
      <div class="resume-slots" aria-label="选择贴纸位置">
        <NButton v-for="(item, index) in entries" :key="item.id" :type="selected === index ? 'primary' : 'default'" :disabled="saving || uploading" @click="selected = index; imageError = ''">
          {{ item.id }} · {{ placements[index] }}
        </NButton>
      </div>
      <div class="resume-editor">
        <fieldset :disabled="saving || uploading" class="resume-fields">
          <NFormItem label="贴纸图片">
            <div>
              <input ref="fileInput" type="file" accept="image/png,image/jpeg,image/webp" hidden @change="upload">
              <NButton :loading="uploading" :disabled="saving || uploading" @click="fileInput.click()">上传新贴纸</NButton>
              <NButton :disabled="saving || uploading || !entry.image" @click="entry.image = ''; imageError = ''">恢复原贴纸</NButton>
              <p class="resume-hint">推荐透明背景 PNG；最大 5 MB、4096 × 4096 像素，图片会等比居中显示。</p>
              <NAlert v-if="imageError" type="error">{{ imageError }}</NAlert>
            </div>
          </NFormItem>
          <NFormItem label="标题">
            <NInput v-model:value="entry.title" :maxlength="60" show-count :disabled="saving" />
          </NFormItem>
          <NFormItem label="英文分类">
            <NInput v-model:value="entry.category" :maxlength="60" show-count :disabled="saving" />
          </NFormItem>
          <NFormItem label="介绍">
            <NInput v-model:value="entry.description" type="textarea" :autosize="{ minRows: 6, maxRows: 14 }" :maxlength="1200" show-count :disabled="saving" />
          </NFormItem>
          <NFormItem label="标签（最多 5 个，每个不超过 24 字）">
            <NDynamicTags v-model:value="entry.tags" :max="5" :disabled="saving || uploading" />
          </NFormItem>
        </fieldset>
        <aside class="resume-preview" aria-label="成就卡片预览">
          <small>卡片预览 · {{ placements[selected] }}</small>
          <img :key="`${selected}-${entry.image}`" :src="imageURL(entry.image)" :alt="`${entry.title}贴纸`">
          <span>{{ entry.category }}</span>
          <h2>{{ entry.title }}</h2>
          <p>{{ entry.description }}</p>
          <ul><li v-for="tag in entry.tags" :key="tag">{{ tag }}</li></ul>
          <small v-if="dirty">修改尚未保存</small>
        </aside>
      </div>
    </template>
  </NSpin>
</template>

<style scoped>
.resume-toolbar { display: flex; align-items: center; justify-content: space-between; gap: 20px; margin: 12px 0 20px; }
.resume-slots { display: flex; flex-wrap: wrap; gap: 10px; margin-bottom: 24px; }
.resume-editor { display: grid; grid-template-columns: minmax(0, 1.2fr) minmax(0, 1fr); gap: 32px; align-items: start; }
.resume-fields { border: 0; padding: 0; min-width: 0; }
.resume-hint { margin-top: 10px; opacity: .7; font-size: 12px; }
.resume-preview { border-radius: 14px; padding: 28px; background: #211e21; color: #f0e9e4; overflow-wrap: anywhere; }
.resume-preview img { display: block; width: 140px; height: 120px; object-fit: contain; margin: 20px 0; }
.resume-preview span, .resume-preview small { display: block; color: #df635f; }
.resume-preview h2 { font-size: 26px; margin: 12px 0; }
.resume-preview p { white-space: pre-wrap; line-height: 1.9; }
.resume-preview ul { display: flex; flex-wrap: wrap; gap: 8px; margin: 20px 0; padding: 0; list-style: none; }
.resume-preview li { border: 1px solid #6b4245; padding: 4px 10px; border-radius: 20px; }
@media (max-width: 850px) { .resume-editor { grid-template-columns: 1fr; } .resume-toolbar { align-items: flex-start; flex-direction: column; } }
</style>
