<script setup lang="ts">
import type { Visualizer, VisualizerOptions } from 'butterchurn'
import { onKeyStroke } from '@vueuse/core'
import butterchurn from 'butterchurn'
import butterchurnPresets from 'butterchurn-presets'
import { audioContext, audioNode } from '~/logic/playbackQueue'

const canvas = useTemplateRef('canvas') as Ref<HTMLCanvasElement>
const visualizer = ref<Visualizer | null>(null)

const blendSeconds = 2.0
const allPresets = butterchurnPresets.getPresets()
let originalWidth = 800
let originalHeight = 600
let animationFrameId: number | null = null

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

  if (document.fullscreenElement === canvas.value) {
    document.exitFullscreen()
    canvas.value.width = originalWidth
    canvas.value.height = originalHeight
    visualizer.value.setRendererSize(originalWidth, originalHeight)
  }
  else {
    canvas.value.requestFullscreen()
    canvas.value.width = screen.width
    canvas.value.height = screen.height
    visualizer.value.setRendererSize(screen.width, screen.height)
  }
}

onKeyStroke(['F', 'f'], (e) => {
  e.preventDefault()
  toggleFullscreen()
})

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
    // onlyUseWASM: true,
  }

  visualizer.value = butterchurn.createVisualizer(audioContext.value, canvas.value, options) as Visualizer

  visualizer.value.connectAudio(audioNode.value)

  const randomPresetInt = Math.floor(Math.random() * Object.keys(allPresets).length)
  const preset = allPresets[Object.keys(allPresets)[randomPresetInt]]
  visualizer.value.loadPreset(preset, blendSeconds)

  visualizer.value.setRendererSize(width, height)

  renderLoop()
}

onMounted(() => {
  if (canvas.value) {
    canvas.value.addEventListener('dblclick', toggleFullscreen)
  }
  createVisualizer()
})

onUnmounted(() => {
  if (canvas.value) {
    canvas.value.removeEventListener('dblclick', toggleFullscreen)
  }
  stopRenderLoop()
  visualizer.value = null
})
</script>

<template>
  <div class="h-100dvh w-full">
    <canvas ref="canvas" class="h-full w-full" />
  </div>
</template>
