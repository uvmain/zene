<script setup lang="ts">
import type { ButterchurnVisualizer } from 'butterchurn'
import butterchurn from 'butterchurn'
import butterchurnPresets from 'butterchurn-presets'
import { audioContext, audioNode } from '~/logic/playbackQueue'

const canvas = useTemplateRef('canvas') as Ref<HTMLCanvasElement>
const visualizer = ref<ButterchurnVisualizer | null>(null)
const blendSeconds = 2.0
const allPresets = butterchurnPresets.getPresets()

let animationFrameId: number | null = null

function renderLoop() {
  if (visualizer.value) {
    visualizer.value.render()
    animationFrameId = requestAnimationFrame(renderLoop)
  }
}

function stopRenderLoop() {
  if (animationFrameId !== null) {
    cancelAnimationFrame(animationFrameId)
    animationFrameId = null
  }
}

watch([audioContext, audioNode], () => {
  if (visualizer.value) {
    stopRenderLoop()
    visualizer.value.destroy()
    visualizer.value = null
  }
  if (audioContext.value && audioNode.value) {
    createVisualizer()
  }
})

function createVisualizer() {
  if (!canvas.value || !audioContext.value || !audioNode.value) {
    return
  }

  visualizer.value = butterchurn.createVisualizer(audioContext.value, canvas.value, {
    width: 800,
    height: 600,
  })

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
  if (visualizer.value) {
    visualizer.value.destroy()
    visualizer.value = null
  }
})
</script>

<template>
  <canvas ref="canvas" width="600" height="600"></canvas>
</template>
