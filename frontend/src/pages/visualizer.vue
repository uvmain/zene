<script setup lang="ts">
import type { Visualizer, VisualizerOptions } from 'butterchurn'
import { onKeyStroke } from '@vueuse/core'
import butterchurn from 'butterchurn'
import butterchurnPresets from 'butterchurn-presets'
import { audioContext, audioNode } from '~/logic/playbackQueue'

const canvas = useTemplateRef('canvas') as Ref<HTMLCanvasElement>
const gridParent = useTemplateRef('grid') as Ref<HTMLDivElement>
const visualizer = ref<Visualizer | null>(null)
const currentVisualizerIndex = ref(0)
const initialFadeIn = ref(true)
const isFullScreen = ref(false)

const allPresets = butterchurnPresets.getPresets()
let originalWidth = 800
let originalHeight = 600
const meshSize = { x: 48, y: 36 }
let animationFrameId: number | null = null
let presetInterval: NodeJS.Timeout | null = null
const intervalSeconds = 25.0
const blendSeconds = 2.7

function renderLoop() {
  if (visualizer.value != null) {
    visualizer.value.render()
    requestAnimationFrame(renderLoop)
  }
}

function stopRenderLoop() {
  if (animationFrameId !== null) {
    cancelAnimationFrame(animationFrameId)
    animationFrameId = null
  }
}

watch([audioContext, audioNode], () => {
  stopRenderLoop()
  visualizer.value = null
  if (audioContext.value != null && audioNode.value != null) {
    createVisualizer()
  }
})

function toggleFullscreen() {
  if (!canvas.value || !visualizer.value) {
    return
  }

  if (document.fullscreenElement === gridParent.value) {
    document.exitFullscreen()
    setWindowed()
  }
  else {
    gridParent.value.requestFullscreen()
    setFullscreen()
  }
}

function setWindowed() {
  if (!canvas.value || !visualizer.value) {
    return
  }
  canvas.value.width = originalWidth
  canvas.value.height = originalHeight
  visualizer.value.setRendererSize(originalWidth, originalHeight)
  isFullScreen.value = false
}

function setFullscreen() {
  if (!canvas.value || !visualizer.value) {
    return
  }
  canvas.value.width = screen.width
  canvas.value.height = screen.height
  visualizer.value.setRendererSize(screen.width, screen.height)
  isFullScreen.value = true
}

function createVisualizer() {
  if (!canvas.value || !audioContext.value || !audioNode.value) {
    return
  }

  let width = 800
  let height = 600
  const parent = canvas.value.parentElement
  if (parent) {
    width = parent.clientWidth
    height = parent.clientHeight
  }
  else {
    width = window.innerWidth
    height = window.innerHeight
  }
  originalWidth = width
  originalHeight = height
  canvas.value.width = width
  canvas.value.height = height

  const options: VisualizerOptions = {
    width,
    height,
    pixelRatio: window.devicePixelRatio || 1,
    meshWidth: meshSize.x,
    meshHeight: meshSize.y,
  }

  visualizer.value = butterchurn.createVisualizer(audioContext.value, canvas.value, options) as Visualizer

  visualizer.value.connectAudio(audioNode.value)

  loadRandomPreset()

  visualizer.value.setRendererSize(width, height)

  // load next preset on 10 second timer
  presetInterval = setInterval(loadRandomPreset, intervalSeconds * 1000)

  renderLoop()
}

function loadRandomPreset() {
  if (!visualizer.value) {
    return
  }
  currentVisualizerIndex.value = Math.floor(Math.random() * Object.keys(allPresets).length)
  const preset = allPresets[Object.keys(allPresets)[currentVisualizerIndex.value]]
  visualizer.value.loadPreset(preset, blendSeconds)
}

onKeyStroke(['F', 'f'], (e) => {
  e.preventDefault()
  toggleFullscreen()
})

onMounted(() => {
  if (canvas.value) {
    canvas.value.addEventListener('dblclick', toggleFullscreen)
  }
  document.addEventListener('fullscreenchange', () => {
    if (document.fullscreenElement === gridParent.value) {
      setFullscreen()
    }
    else {
      setWindowed()
    }
  })
  createVisualizer()

  setTimeout(() => {
    initialFadeIn.value = false
  }, 1000)
})

onUnmounted(() => {
  if (canvas.value) {
    canvas.value.removeEventListener('dblclick', toggleFullscreen)
  }
  stopRenderLoop()
  visualizer.value = null
  if (presetInterval) {
    clearInterval(presetInterval)
  }
})
</script>

<template>
  <div ref="grid" class="group grid h-100dvh w-full">
    <canvas ref="canvas" class="z-1 col-span-full row-span-full h-full w-full" />
    <div
      class="corner-cut z-2 col-span-full row-span-full mb-2 ml-auto mr-2 mt-auto w-80 bg-cover bg-center text-primary transition-opacity duration-1000 transition-ease-out group-hover:opacity-100"
      :class="{
        'opacity-100': initialFadeIn,
        'opacity-0': !initialFadeIn,
      }"
    >
      <!-- info panel -->
      <div class="corner-cut z-3 flex flex-col bg-zshade-300/60 px-4 py-2 backdrop-blur-xl dark:bg-zshade-900/60">
        <NavArt v-if="isFullScreen" />
        <PlayerProgressBar v-if="isFullScreen" class="mt-2" :compact="true" />
        <PlayerMediaControls v-if="isFullScreen" class="mt-2" :compact="true" />
        <p class="text-wrap text-sm">
          Press F or double-click to toggle fullscreen.
        </p>
      </div>
    </div>
  </div>
</template>
