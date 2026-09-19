<script setup>
import { onBeforeUnmount, onMounted, ref } from 'vue'
import * as THREE from 'three'
import { GLTFLoader } from 'three/addons/loaders/GLTFLoader.js'
import { MeshoptDecoder } from 'three/addons/libs/meshopt_decoder.module.js'

const props = defineProps({ src: { type: String, required: true } })
const emit = defineEmits(['loaded'])
const canvas = ref()
const error = ref('')
let frame = 0
let renderer
let scene

onMounted(() => {
  scene = new THREE.Scene()
  scene.background = new THREE.Color(0x16181D)
  const camera = new THREE.PerspectiveCamera(38, 5 / 3, 0.01, 1000)
  renderer = new THREE.WebGLRenderer({ antialias: true, canvas: canvas.value })
  renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2))
  renderer.setSize(300, 180, false)
  renderer.outputColorSpace = THREE.SRGBColorSpace
  scene.add(new THREE.HemisphereLight(0xFFFFFF, 0x303744, 3))

  const loader = new GLTFLoader()
  loader.setMeshoptDecoder(MeshoptDecoder)
  loader.load(props.src, (gltf) => {
    const model = gltf.scene
    const clock = new THREE.Clock()
    const mixer = new THREE.AnimationMixer(model)
    scene.add(model)
    const box = new THREE.Box3().setFromObject(model)
    const center = box.getCenter(new THREE.Vector3())
    const size = Math.max(...box.getSize(new THREE.Vector3()).toArray(), 0.1)
    model.position.sub(center)
    camera.position.set(size * 1.2, size * 0.75, size * 1.8)
    camera.lookAt(0, 0, 0)
    camera.near = Math.max(size / 100, 0.01)
    camera.far = size * 20
    camera.updateProjectionMatrix()
    const resumeCamera = gltf.cameras.find(item => item.name === 'ResumeCamera')
    if (resumeCamera?.isPerspectiveCamera) {
      resumeCamera.aspect = 5 / 3
      resumeCamera.updateProjectionMatrix()
    }
    const cameraAction = gltf.animations.find(item => item.name === 'CameraAction')
    if (cameraAction)
      mixer.clipAction(cameraAction).play()
    emit('loaded')
    const render = () => {
      frame = requestAnimationFrame(render)
      mixer.update(clock.getDelta())
      renderer.render(scene, resumeCamera || camera)
    }
    render()
  }, undefined, (reason) => {
    error.value = `GLB 预览失败：${reason?.message || '资源无法加载'}`
  })
})

onBeforeUnmount(() => {
  cancelAnimationFrame(frame)
  scene?.traverse((object) => {
    object.geometry?.dispose?.()
    const materials = Array.isArray(object.material) ? object.material : [object.material]
    materials.filter(Boolean).forEach((material) => {
      Object.values(material).forEach(value => value?.isTexture && value.dispose())
      material.dispose?.()
    })
  })
  renderer?.dispose()
})
</script>

<template>
  <div class="model-preview">
    <canvas ref="canvas" />
    <span v-if="error">{{ error }}</span>
  </div>
</template>

<style scoped>
.model-preview { display: grid; gap: 5px; }
canvas { width: 300px; max-width: 100%; aspect-ratio: 5 / 3; border-radius: 5px; }
span { color: #d03050; font-size: 11px; }
</style>
