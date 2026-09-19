<script setup>
import { onMounted, ref } from 'vue'
import { NButton, NTabPane, NTabs } from 'naive-ui'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import ResumeEditor from './ResumeEditor.vue'

import CommonPage from '@/components/common/CommonPage.vue'
import api from '@/api'

defineOptions({ name: '关于我' })

const aboutContent = ref('')
const btnLoading = ref(false)
const activeTab = ref('resume')

onMounted(async () => {
  const resp = await api.getAbout()
  aboutContent.value = resp.data
})

async function handleSave() {
  try {
    btnLoading.value = true
    await api.updateAbout({ content: aboutContent.value })
    window.$message.success('更新成功')
  }
  finally {
    btnLoading.value = false
  }
}
</script>

<template>
  <CommonPage :show-header="false">
    <div class="mb-4 flex items-center justify-between">
      <div class="mx-1 text-2xl font-bold">
        关于我
      </div>
      <NButton v-if="activeTab === 'text'" type="primary" :loading="btnLoading" @click="handleSave">
        <template #icon>
          <span v-if="!btnLoading" class="i-line-md:confirm-circle" />
        </template>
        保存
      </NButton>
    </div>
    <NTabs v-model:value="activeTab" type="line">
      <NTabPane name="resume" tab="成就贴纸" display-directive="show">
        <ResumeEditor />
      </NTabPane>
      <NTabPane name="text" tab="更多关于我">
        <MdEditor v-model="aboutContent" style="height: calc(100vh - 290px)" />
      </NTabPane>
    </NTabs>
  </CommonPage>
</template>

<style lang="scss">
.md-preview {
  ul,
  ol {
    list-style: revert;
  }
}
</style>
