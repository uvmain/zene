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

let animationFrameId: number | null = null
const targetFPS = 60
const frameDuration = 1000 / targetFPS

function renderLoop() {
  if (visualizer.value != null) {
    visualizer.value.render()
    animationFrameId = window.setTimeout(renderLoop, frameDuration)
  }
}

function stopRenderLoop() {
  if (animationFrameId !== null) {
    clearTimeout(animationFrameId)
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
  if (!canvas.value) {
    return
  }

  if (document.fullscreenElement) {
    document.exitFullscreen()
  }
  else {
    canvas.value.requestFullscreen()
  }
}

onKeyStroke(['F', 'f'], (e) => {
  e.preventDefault()
  toggleFullscreen()
})

onMounted(() => {
  if (canvas.value) {
    canvas.value.addEventListener('dblclick', toggleFullscreen)
  }
})

function createVisualizer() {
  if (!canvas.value || !audioContext.value || !audioNode.value) {
    return
  }

  const options: VisualizerOptions = {
    width: 800,
    height: 600,
    pixelRatio: window.devicePixelRatio || 1,
    onlyUseWASM: true,
    hideControls: false,
  }

  visualizer.value = butterchurn.createVisualizer(audioContext.value, canvas.value, options) as Visualizer

  visualizer.value.connectAudio(audioNode.value)

  const randomPresetInt = Math.floor(Math.random() * Object.keys(allPresets).length)
  const preset = allPresets[Object.keys(allPresets)[randomPresetInt]]
  visualizer.value.loadPreset(preset, blendSeconds)

  visualizer.value.setRendererSize(600, 600)

  renderLoop()
}

onMounted(() => {
  createVisualizer()
})

onUnmounted(() => {
  stopRenderLoop()
  visualizer.value = null
})
</script>

<template>
  <canvas ref="canvas" width="600" height="600"></canvas>
</template>
