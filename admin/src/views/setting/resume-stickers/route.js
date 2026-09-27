const Layout = () => import('@/layout/index.vue')

export default {
  name: 'ResumeStickers',
  path: '/resume-stickers',
  component: Layout,
  redirect: '/resume-stickers/edit',
  isCatalogue: true,
  meta: {
    title: '3D 履历贴纸',
    icon: 'mdi:cloud-upload-outline',
    order: 7,
  },
  children: [
    {
      name: 'ResumeStickerEditor',
      path: 'edit',
      component: () => import('./index.vue'),
      meta: {
        title: '3D 履历贴纸',
        icon: 'mdi:cloud-upload-outline',
        order: 1,
      },
    },
  ],
}
