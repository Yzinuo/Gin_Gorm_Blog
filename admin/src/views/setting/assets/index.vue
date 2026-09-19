<script setup>
import { computed, onMounted, reactive } from 'vue'
import { NAlert, NButton, NCard, NProgress, NSpin, NTag } from 'naive-ui'
import { onBeforeRouteLeave } from 'vue-router'
import AssetModelPreview from './AssetModelPreview.vue'
import api from '@/api'

const definitions = [
  { key: 'home.xray', title: '首页 X-Ray', help: 'Before / After 必须尺寸一致，宽度至少 1280 px。', parts: [{ name: 'before', label: 'Before' }, { name: 'after', label: 'After' }], accept: 'image/png,image/jpeg,image/webp' },
  { key: 'resume.model', title: 'Resume 3D 模型', help: '保留原始 GLB；若原文件超过 6 MiB，须同时上传固定脚本生成的 delivery。', parts: [{ name: 'file', label: '原始 GLB' }, { name: 'delivery', label: '优化后 GLB（原文件 ≤6 MiB 时可选）', optional: true }], accept: '.glb,model/gltf-binary' },
  ...Array.from({ length: 8 }, (_, index) => ({
    key: `resume.sticker.${String(index + 1).padStart(2, '0')}`,
    title: `Resume 贴纸 ${String(index + 1).padStart(2, '0')}`,
    help: 'PNG / JPEG / WebP，最大 5 MiB、4096 × 4096。',
    parts: [{ name: 'file', label: '图片' }],
    accept: 'image/png,image/jpeg,image/webp',
  })),
]

const state = reactive(Object.fromEntries(definitions.map(item => [item.key, { loading: true, uploading: false, progress: 0, versions: [], files: {}, error: '', loadedPreviews: {}, confirmedPreviews: {} }])))
const totalBytes = computed(() => Object.values(state).reduce((sum, item) => sum + item.versions.reduce((versionSum, version) => versionSum + [...version.source_manifest.files, ...version.delivery_manifest.files].reduce((fileSum, file) => fileSum + file.bytes, 0), 0), 0))

function formatBytes(value) {
  if (!value)
    return '0 B'
  const units = ['B', 'KiB', 'MiB', 'GiB']
  const power = Math.min(Math.floor(Math.log(value) / Math.log(1024)), units.length - 1)
  return `${(value / 1024 ** power).toFixed(power ? 2 : 0)} ${units[power]}`
}

async function load(key) {
  const current = state[key]
  current.loading = true
  current.error = ''
  try {
    const response = await api.getAssetVersions(key)
    current.versions = response.data || []
  }
  catch (error) {
    current.error = error?.data || '版本列表加载失败。'
  }
  finally {
    current.loading = false
  }
}

function selectFile(key, part, event) {
  state[key].files[part] = event.target.files?.[0] || null
}

async function upload(definition) {
  const current = state[definition.key]
  if (definition.parts.some(part => !part.optional && !current.files[part.name])) {
    window.$message.error('请先选择全部必需文件。')
    return
  }
  const body = new FormData()
  definition.parts.forEach((part) => {
    if (current.files[part.name])
      body.append(part.name, current.files[part.name])
  })
  current.uploading = true
  current.progress = 0
  current.error = ''
  try {
    await api.createAssetVersion(definition.key, body, ({ loaded, total }) => current.progress = total ? Math.round(loaded / total * 100) : 0)
    window.$message.success('上传与校验完成。请预览后手动发布。')
    current.files = {}
    await load(definition.key)
  }
  catch (error) {
    current.error = error?.data || error?.message || '上传处理失败。'
  }
  finally {
    current.uploading = false
  }
}

function confirmPublish(definition, version) {
  const current = state[definition.key].versions.find(item => item.status === 'published')
  window.$dialog.warning({
    title: version.status === 'archived' ? '确认回滚资源？' : '确认发布资源？',
    content: `${definition.title}：${current?.version?.slice(0, 8) || '尚无线上版本'} → ${version.version.slice(0, 8)}。发布只切换数据库指针，不覆盖对象。`,
    positiveText: version.status === 'archived' ? '确认回滚' : '确认发布',
    negativeText: '取消',
    async onPositiveClick() {
      await api.publishAssetVersion(definition.key, version.id)
      window.$message.success('资源指针已切换。')
      await load(definition.key)
    },
  })
}

function previewFiles(version) {
  return version.delivery_manifest.files || []
}

function markPreviewLoaded(key, version, file) {
  state[key].loadedPreviews[`${version.id}:${file.key}`] = true
}

function previewReady(key, version) {
  const files = previewFiles(version)
  return files.length > 0 && files.every(file => state[key].loadedPreviews[`${version.id}:${file.key}`])
}

onMounted(() => Promise.all(definitions.map(item => load(item.key))))

onBeforeRouteLeave(() => {
  const hasUnpublishedVersion = Object.values(state).some(item => item.versions.some(version => version.status === 'ready'))
  // eslint-disable-next-line no-alert
  if (hasUnpublishedVersion && !window.confirm('仍有已上传但未发布的资源版本，确定离开吗？'))
    return false
})
</script>

<template>
  <div class="asset-page">
    <header>
      <div>
        <h1>资源管理</h1>
        <p>上传 → 自动校验 → 预览 → 手动发布。对象 URL 不可变，回滚只切换当前版本。</p>
      </div>
      <NTag :type="totalBytes >= 8 * 1024 ** 3 ? 'error' : 'info'">
        受管对象：{{ formatBytes(totalBytes) }}
      </NTag>
    </header>
    <NAlert v-if="totalBytes >= 8 * 1024 ** 3" type="error" title="R2 存储已达到 8 GiB 成本警戒线">
      请先运行 GC dry-run 并核对引用关系，禁止直接删除当前或最近三个发布版本。
    </NAlert>

    <NCard v-for="definition in definitions" :key="definition.key" :title="definition.title" size="small">
      <NSpin :show="state[definition.key].loading">
        <p class="hint">
          {{ definition.help }} · <code>{{ definition.key }}</code>
        </p>
        <div class="upload-row">
          <label v-for="part in definition.parts" :key="part.name">
            <span>{{ part.label }}</span>
            <input type="file" :accept="definition.accept" :disabled="state[definition.key].uploading" @change="selectFile(definition.key, part.name, $event)">
          </label>
          <NButton type="primary" :loading="state[definition.key].uploading" @click="upload(definition)">
            上传并校验
          </NButton>
        </div>
        <NProgress v-if="state[definition.key].uploading" type="line" :percentage="state[definition.key].progress" processing />
        <NAlert v-if="state[definition.key].error" type="error">
          {{ state[definition.key].error }}
        </NAlert>

        <div v-if="state[definition.key].versions.length" class="versions">
          <article v-for="version in state[definition.key].versions.slice(0, 3)" :key="version.id" class="version" :class="version.status">
            <div class="version-head">
              <div>
                <NTag :type="version.status === 'published' ? 'success' : version.status === 'ready' ? 'warning' : 'default'">
                  {{ version.status }}
                </NTag>
                <strong>{{ version.version.slice(0, 12) }}</strong>
                <small>{{ new Date(version.created_at).toLocaleString() }}</small>
              </div>
              <NButton v-if="['ready', 'archived'].includes(version.status)" size="small" :type="version.status === 'ready' ? 'primary' : 'warning'" :disabled="!state[definition.key].confirmedPreviews[version.id]" @click="confirmPublish(definition, version)">
                {{ version.status === 'archived' ? '回滚到此版本' : '发布' }}
              </NButton>
            </div>
            <ul class="checks">
              <li v-for="check in version.validation_report.checks" :key="check">
                ✓ {{ check }}
              </li>
            </ul>
            <div class="previews">
              <div v-for="file in previewFiles(version)" :key="file.key" class="preview-item">
                <img v-if="file.content_type.startsWith('image/')" :src="file.url" :alt="file.role" @load="markPreviewLoaded(definition.key, version, file)">
                <AssetModelPreview v-else-if="file.content_type === 'model/gltf-binary'" :src="file.url" @loaded="markPreviewLoaded(definition.key, version, file)" />
                <a :href="file.url" target="_blank" rel="noopener">{{ file.role }} · {{ formatBytes(file.bytes) }}<template v-if="file.width"> · {{ file.width }}×{{ file.height }}</template></a>
              </div>
            </div>
            <NButton v-if="['ready', 'archived'].includes(version.status)" class="preview-confirm" size="tiny" :disabled="!previewReady(definition.key, version)" @click="state[definition.key].confirmedPreviews[version.id] = true">
              {{ state[definition.key].confirmedPreviews[version.id] ? '已确认预览' : '确认预览无误' }}
            </NButton>
          </article>
        </div>
        <p v-else class="empty">
          尚无版本。线上将继续使用仓库 fallback。
        </p>
      </NSpin>
    </NCard>
  </div>
</template>

<style scoped>
.asset-page { display: grid; gap: 18px; padding: 18px; }
header { display: flex; align-items: center; justify-content: space-between; gap: 20px; }
h1 { margin: 0 0 6px; font-size: 24px; }
header p, .hint, .empty { color: #777; }
.upload-row { display: flex; align-items: end; flex-wrap: wrap; gap: 14px; margin: 16px 0; }
.upload-row label { display: grid; gap: 6px; font-size: 12px; }
.versions { display: grid; gap: 12px; margin-top: 18px; }
.version { border: 1px solid #e6e8eb; border-radius: 8px; padding: 14px; }
.version.published { border-color: #18a05866; }
.version-head, .version-head > div { display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 10px; }
.version-head small { color: #888; }
.checks { display: flex; flex-wrap: wrap; gap: 8px 16px; padding: 0; list-style: none; color: #39785a; font-size: 12px; }
.previews { display: flex; flex-wrap: wrap; gap: 10px; }
.preview-item { display: grid; gap: 6px; color: inherit; font-size: 11px; }
.preview-item a { color: inherit; }
.previews img { width: 150px; height: 90px; object-fit: cover; border-radius: 5px; background: #222; }
.preview-confirm { margin-top: 10px; }
@media (max-width: 700px) { header { align-items: flex-start; flex-direction: column; } }
</style>
