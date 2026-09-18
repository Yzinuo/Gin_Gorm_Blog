<script setup>
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

import UploadOne from './UploadOne.vue'
import BannerPage from '@/components/BannerPage.vue'

import { useUserStore } from '@/store'
import api from '@/api'

const userStore = useUserStore()
const router = useRouter()
const loading = ref(true)

const form = reactive({
  avatar: '',
  nickname: '',
  intro: '',
  website: '',
})

function syncForm() {
  Object.assign(form, {
    avatar: userStore.avatar,
    nickname: userStore.nickname,
    intro: userStore.intro,
    website: userStore.website,
  })
}

async function loadUserInfo() {
  loading.value = true
  try {
    await userStore.getUserInfo()
    if (!userStore.userId) {
      router.replace('/')
      return
    }
    syncForm()
  }
  catch {
    router.replace('/')
  }
  finally {
    loading.value = false
  }
}

onMounted(loadUserInfo)

async function updateUserInfo() {
  try {
    await api.updateUser(form)
    window.$message?.success('修改成功!')
    await userStore.getUserInfo()
    syncForm()
  }
  catch (err) {
    console.error(err)
  }
}
</script>

<template>
  <BannerPage label="user" title="个人中心" card :loading="loading">
    <p class="mb-6 text-xl font-bold">
      基本信息
    </p>
    <div class="grid grid-cols-12 gap-4">
      <div class="col-span-4 f-c-c">
        <UploadOne v-model:preview="form.avatar" />
      </div>
      <div class="col-span-8 lg:col-span-7">
        <div class="my-6 space-y-3">
          <div
            v-for="item of [
              { label: '昵称', key: 'nickname' },
              { label: '个人网站', key: 'website' },
              { label: '简介', key: 'intro' },
            ]" :key="item.label"
          >
            <div class="mb-2">
              {{ item.label }}
            </div>
            <input
              v-model="form[item.key]" required :placeholder="`请输入${item.label}`"
              class="block w-full border-0 rounded-md p-2 text-foreground shadow-sm outline-none ring-1 ring-line ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-brand"
            >
          </div>
          <div>
            <div class="mb-2">
              邮箱
            </div>
            <input
              :value="userStore.email" disabled
              class="block w-full cursor-not-allowed border-0 rounded-md bg-raised p-2 text-muted shadow-sm outline-none ring-1 ring-line ring-inset"
            >
            <p class="mt-1 text-xs text-muted">
              邮箱用于登录，目前不能在这里修改。
            </p>
          </div>
        </div>
        <button class="the-button mt-2" @click="updateUserInfo">
          修改
        </button>
      </div>
      <div class="col-span-0 lg:col-span-1" />
    </div>
  </BannerPage>
</template>
