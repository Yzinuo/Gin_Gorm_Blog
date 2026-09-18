import * as THREE from 'three'
import { GLTFLoader } from 'three/addons/loaders/GLTFLoader.js'
import { RoomEnvironment } from 'three/addons/environments/RoomEnvironment.js'

export function scrollFrame(sections, focus, totalFrames) {
  if (!sections.length || focus <= sections[0].y)
    return 0
  for (let i = 1; i < sections.length; i++) {
    const a = sections[i - 1]
    const b = sections[i]
    if (focus <= b.y)
      return THREE.MathUtils.lerp(a.frame, b.frame, (focus - a.y) / Math.max(1, b.y - a.y))
  }
  return totalFrames
}

function disposeModel(model) {
  const geometries = new Set()
  const materials = new Set()
  const textures = new Set()
  model?.traverse((object) => {
    if (object.geometry)
      geometries.add(object.geometry)
    if (object.material)
      [].concat(object.material).forEach(material => materials.add(material))
    object.skeleton?.dispose()
  })
  materials.forEach((material) => {
    Object.values(material).forEach(value => value?.isTexture && textures.add(value))
    material.dispose()
  })
  textures.forEach((texture) => {
    texture.source?.data?.close?.()
    texture.dispose()
  })
  geometries.forEach(geometry => geometry.dispose())
}

export function createResumeScene({ canvas, stage, map, getSections, getStickers = () => [], onStickerError = () => {}, onProgress, onReady, onChapter, onError }) {
  const controller = new AbortController()
  const { signal } = controller
  const motion = matchMedia('(prefers-reduced-motion: reduce)')
  const scene = new THREE.Scene()
  let renderer, environment, model, mixer, action, resizeObserver, camera, animationId
  let previousTime = 0
  let smoothFrame = 0
  let targetFrame = 0
  let chapter = -1
  let manualFrame = null
  let ready = false
  let disposed = false
  let pointerActive = false
  const eyes = []
  const stickerMaterials = new Map()
  const sourceMaterials = new Set()
  let stickerRevision = 0
  const mouse = new THREE.Vector2()
  const eyeCenter = new THREE.Vector3()
  const position = new THREE.Vector3()
  const offset = new THREE.Quaternion()
  const desired = new THREE.Quaternion()
  const rotation = new THREE.Euler(0, 0, 0, 'YXZ')

  async function updateStickers() {
    if (!model || disposed || !stickerMaterials.size)
      return
    const revision = ++stickerRevision
    const results = await Promise.allSettled(getStickers().map(async (sticker) => {
      const slots = stickerMaterials.get(sticker.object)
      if (!slots)
        return
      if (!sticker.image) {
        slots.forEach(({ material, original }) => {
          if (material.map !== original)
            material.map?.dispose()
          material.map = original
          material.needsUpdate = true
        })
        return
      }
      const timeout = new AbortController()
      const abort = () => timeout.abort()
      signal.addEventListener('abort', abort, { once: true })
      const timer = setTimeout(abort, 15000)
      let bitmap
      try {
        const response = await fetch(sticker.image, { signal: timeout.signal })
        if (!response.ok)
          throw new Error('Sticker unavailable')
        bitmap = await createImageBitmap(await response.blob())
        if (disposed || revision !== stickerRevision)
          return
        slots.forEach(({ material, original }) => {
          // Preserve the original UV transform and aspect ratio. Fit new images
          // inside the texture without stretching the sticker geometry.
          const surface = document.createElement('canvas')
          const ratio = (original?.image?.width || 1) / (original?.image?.height || 1)
          surface.width = ratio >= 1 ? 1024 : Math.round(1024 * ratio)
          surface.height = ratio >= 1 ? Math.round(1024 / ratio) : 1024
          const scale = Math.min(surface.width / bitmap.width, surface.height / bitmap.height)
          surface.getContext('2d').drawImage(bitmap, (surface.width - bitmap.width * scale) / 2, (surface.height - bitmap.height * scale) / 2, bitmap.width * scale, bitmap.height * scale)
          const texture = new THREE.CanvasTexture(surface)
          texture.flipY = original?.flipY ?? false
          texture.colorSpace = THREE.SRGBColorSpace
          if (original) {
            texture.channel = original.channel
            texture.offset.copy(original.offset)
            texture.repeat.copy(original.repeat)
            texture.center.copy(original.center)
            texture.rotation = original.rotation
            texture.wrapS = original.wrapS
            texture.wrapT = original.wrapT
          }
          if (material.map !== original)
            material.map?.dispose()
          material.map = texture
          material.needsUpdate = true
        })
      }
      catch (error) {
        if (!disposed && revision === stickerRevision) {
          slots.forEach(({ material, original }) => {
            if (material.map !== original)
              material.map?.dispose()
            material.map = original
            material.needsUpdate = true
          })
        }
        throw error
      }
      finally {
        bitmap?.close()
        clearTimeout(timer)
        signal.removeEventListener('abort', abort)
      }
    }))
    if (!disposed && revision === stickerRevision)
      onStickerError(results.some(result => result.status === 'rejected'))
  }

  function updateScroll() {
    manualFrame = null
    const bounds = stage.getBoundingClientRect()
    const focus = innerWidth <= 640 ? bounds.bottom + (innerHeight - bounds.bottom) / 2 : innerHeight / 2
    const sections = [...getSections()].map((element) => {
      const rect = element.getBoundingClientRect()
      return { y: rect.top + rect.height / 2, frame: Number(element.dataset.frame) }
    })
    targetFrame = scrollFrame(sections, focus, map.totalFrames)
  }

  function resize() {
    if (!renderer || disposed)
      return
    renderer.setPixelRatio(Math.min(devicePixelRatio || 1, innerWidth <= 640 ? 1.25 : 1.6))
    renderer.setSize(stage.clientWidth, stage.clientHeight, false)
    if (camera) {
      camera.aspect = stage.clientWidth / Math.max(1, stage.clientHeight)
      camera.updateProjectionMatrix()
    }
    updateScroll()
  }

  function animate(time) {
    if (disposed || !ready || document.hidden)
      return
    const delta = Math.min((time - (previousTime || time)) / 1000, 0.05)
    previousTime = time
    const target = manualFrame ?? targetFrame
    smoothFrame = motion.matches ? target : THREE.MathUtils.damp(smoothFrame, target, 7, delta)
    action.time = THREE.MathUtils.clamp(smoothFrame / map.fps, 0, action.getClip().duration)
    mixer.update(0)
    scene.updateMatrixWorld(true)
    eyeCenter.set(0, 0, 0)
    eyes.forEach((eye) => {
      eye.object.getWorldPosition(position)
      eyeCenter.add(position)
    })
    eyeCenter.multiplyScalar(1 / eyes.length).project(camera)
    const trackPointer = pointerActive && !motion.matches
    const yaw = trackPointer ? THREE.MathUtils.clamp((mouse.x - eyeCenter.x) * 0.11, -0.14, 0.14) : 0
    const pitch = trackPointer ? THREE.MathUtils.clamp(-(mouse.y - eyeCenter.y) * 0.065, -0.087, 0.087) : 0
    offset.setFromEuler(rotation.set(pitch, yaw, 0))
    eyes.forEach((eye) => {
      desired.copy(offset).multiply(eye.base)
      eye.object.quaternion.slerp(desired, motion.matches ? 1 : 1 - Math.exp(-12 * delta))
    })
    const current = smoothFrame < 25 ? 0 : smoothFrame > 445 ? 9 : Math.min(8, Math.max(1, Math.round(smoothFrame / 50)))
    if (current !== chapter) {
      chapter = current
      onChapter(current)
    }
    renderer.render(scene, camera)
    animationId = requestAnimationFrame(animate)
  }

  function dispose() {
    if (disposed)
      return
    disposed = true
    controller.abort()
    cancelAnimationFrame(animationId)
    resizeObserver?.disconnect()
    mixer?.stopAllAction()
    if (model)
      mixer?.uncacheRoot(model)
    stickerMaterials.forEach(slots => slots.forEach(({ material, original }) => {
      if (material.map !== original)
        material.map?.dispose()
      material.map = original
    }))
    stickerMaterials.clear()
    disposeModel(model)
    sourceMaterials.forEach(material => material.dispose())
    sourceMaterials.clear()
    environment?.dispose()
    renderer?.dispose()
    renderer?.forceContextLoss()
  }

  async function load() {
    try {
      renderer = new THREE.WebGLRenderer({ canvas, antialias: true, alpha: true, powerPreference: 'high-performance' })
      renderer.outputColorSpace = THREE.SRGBColorSpace
      renderer.toneMapping = THREE.ACESFilmicToneMapping
      renderer.toneMappingExposure = 0.9
      const pmrem = new THREE.PMREMGenerator(renderer)
      const room = new RoomEnvironment()
      environment = pmrem.fromScene(room, 0.04)
      scene.environment = environment.texture
      room.dispose()
      pmrem.dispose()
      scene.environmentIntensity = 0.65
      scene.add(new THREE.HemisphereLight(0xEEE5DF, 0x45383B, 0.7))
      const key = new THREE.DirectionalLight(0xFFF0E3, 1.7)
      key.position.set(-1, 2, 3)
      scene.add(key)
      const fill = new THREE.DirectionalLight(0xE3E8FF, 0.8)
      fill.position.set(2, 1, 1)
      scene.add(fill)
      resizeObserver = new ResizeObserver(resize)
      resizeObserver.observe(stage)
      resize()
      const response = await fetch(`${import.meta.env.BASE_URL}resume/resume-ready.glb`, { signal })
      if (!response.ok)
        throw new Error(`Model request failed: ${response.status}`)
      const size = Number(response.headers.get('content-length'))
      let buffer
      if (response.body && size > 0) {
        const reader = response.body.getReader()
        const chunks = []
        let received = 0
        while (true) {
          const { done, value } = await reader.read()
          if (done)
            break
          chunks.push(value)
          received += value.byteLength
          onProgress(Math.min(99, Math.round(received / size * 100)))
        }
        buffer = await new Blob(chunks).arrayBuffer()
      }
      else buffer = await response.arrayBuffer()
      if (disposed)
        return
      const gltf = await new GLTFLoader().parseAsync(buffer, '')
      if (disposed) {
        disposeModel(gltf.scene)
        return
      }
      model = gltf.scene
      scene.add(model)
      camera = gltf.cameras.find(item => item.name === map.camera)
      const clip = gltf.animations.find(item => item.name === map.animation)
      if (!camera || !clip)
        throw new Error('The résumé camera or animation is missing')
      mixer = new THREE.AnimationMixer(model)
      action = mixer.clipAction(clip)
      action.setLoop(THREE.LoopOnce, 1)
      action.clampWhenFinished = true
      action.play()
      action.paused = true
      for (const name of ['eye_L', 'eye_R']) {
        const object = model.getObjectByName(name)
        if (!object)
          throw new Error(`Missing eye: ${name}`)
        eyes.push({ object, base: object.quaternion.clone() })
      }
      model.traverse((object) => {
        if (object.isMesh && object.name.startsWith('sticker_')) {
          const originals = [].concat(object.material)
          originals.forEach(material => sourceMaterials.add(material))
          object.material = Array.isArray(object.material) ? object.material.map(material => material.clone()) : object.material.clone()
          const materials = [].concat(object.material)
          stickerMaterials.set(object.name, materials.map(material => ({ material, original: material.map })))
          materials.forEach((material) => {
            material.transparent = false
            material.alphaTest = 0.5
            material.side = THREE.DoubleSide
            material.needsUpdate = true
          })
        }
      })
      updateStickers()
      window.addEventListener('scroll', updateScroll, { passive: true, signal })
      window.addEventListener('resize', resize, { passive: true, signal })
      window.addEventListener('pointermove', (event) => {
        const rect = stage.getBoundingClientRect()
        mouse.set((event.clientX - rect.left) / rect.width * 2 - 1, -(event.clientY - rect.top) / rect.height * 2 + 1)
        pointerActive = event.pointerType !== 'touch'
      }, { passive: true, signal })
      window.addEventListener('pointerout', (event) => {
        if (!event.relatedTarget)
          pointerActive = false
      }, { signal })
      document.addEventListener('visibilitychange', () => {
        cancelAnimationFrame(animationId)
        previousTime = 0
        if (!document.hidden)
          animationId = requestAnimationFrame(animate)
      }, { signal })
      canvas.addEventListener('webglcontextlost', (event) => {
        event.preventDefault()
        dispose()
        onError()
      }, { signal })
      resize()
      ready = true
      onProgress(100)
      onReady()
      animationId = requestAnimationFrame(animate)
    }
    catch (error) {
      if (disposed)
        return
      console.error('Unable to load the résumé scene:', error)
      dispose()
      onError()
    }
  }
  load()
  return { dispose, updateStickers, focus: frame => manualFrame = frame }
}
