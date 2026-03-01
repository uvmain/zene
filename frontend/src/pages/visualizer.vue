<script setup lang="ts">
import type { Visualizer, VisualizerOptions } from 'butterchurn'
import type { ButterchurnPreset } from '~/types'
import { onKeyStroke } from '@vueuse/core'
import butterchurn from 'butterchurn'
import { getButterchurnPresets } from '~/logic/backendFetch'
import { audioContext, audioNode, currentlyPlayingTrack } from '~/logic/playbackQueue'

const canvas = useTemplateRef('canvas') as Ref<HTMLCanvasElement>
const gridParent = useTemplateRef('grid') as Ref<HTMLDivElement>
const visualizer = ref<Visualizer | null>(null)
const currentVisualizerIndex = ref(0)
const initialFadeIn = ref(true)
const presetNameFadeIn = ref(true)
const isFullScreen = ref(false)
const fetchedPresets = ref<ButterchurnPreset[]>([])
const presetNameFadeTimeout = ref<NodeJS.Timeout | null>(null)

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
  }
  else {
    gridParent.value.requestFullscreen()
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
  if (!canvas.value || !audioContext.value || !audioNode.value || !fetchedPresets.value || fetchedPresets.value.length === 0) {
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

  presetInterval = setInterval(loadRandomPreset, intervalSeconds * 1000)

  renderLoop()
}

async function loadRandomPreset(blendTimeSeconds = blendSeconds) {
  if (presetNameFadeTimeout.value) {
    clearTimeout(presetNameFadeTimeout.value)
  }
  presetNameFadeIn.value = true
  if (!visualizer.value || !fetchedPresets.value || fetchedPresets.value.length === 0) {
    return
  }
  currentVisualizerIndex.value = Math.floor(Math.random() * fetchedPresets.value.length)
  const preset = fetchedPresets.value[currentVisualizerIndex.value]
  await visualizer.value.loadPreset(preset.preset, blendTimeSeconds)
  presetNameFadeTimeout.value = setTimeout(() => {
    presetNameFadeIn.value = false
  }, Math.max(blendTimeSeconds * 1000, blendSeconds * 1000))
}

onKeyStroke(['F', 'f'], (e) => {
  e.preventDefault()
  toggleFullscreen()
})

watch(currentlyPlayingTrack, (old, current) => {
  if (old !== current) {
    loadRandomPreset()
  }
})

onMounted(async () => {
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

  const presets = await getButterchurnPresets({ random: true, count: 200 })
  fetchedPresets.value = presets

  createVisualizer()

  setTimeout(() => {
    initialFadeIn.value = false
  }, 1000)
})

onUnmounted(() => {
  if (canvas.value) {
    canvas.value.removeEventListener('dblclick', toggleFullscreen)
  }
  document.removeEventListener('fullscreenchange', () => {
    if (document.fullscreenElement === gridParent.value) {
      setFullscreen()
    }
    else {
      setWindowed()
    }
  })
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
      class="corner-cut z-2 col-span-full row-span-full w-80 bg-cover bg-center text-primary transition-opacity duration-1000 transition-ease-out group-hover:opacity-100"
      :class="{
        'opacity-100': initialFadeIn,
        'opacity-0': !initialFadeIn,
        'fixed bottom-4 right-4': isFullScreen,
        'mb-2 ml-auto mr-2 mt-auto': !isFullScreen,
      }"
    >
      <!-- info panel -->
      <div class="corner-cut z-3 flex flex-col bg-zshade-300/60 px-4 py-2 backdrop-blur-xl dark:bg-zshade-900/60">
        <div v-if="isFullScreen" class="flex flex-col gap-2">
          <NavArt />
          <PlayerProgressBar :compact="true" />
          <PlayerMediaControls :compact="true" />
        </div>
        <div class="group/next h-10 flex flex-row cursor-pointer items-center justify-between">
          <div class="flex items-center text-wrap text-sm" @click="loadRandomPreset(0)">
            <p class="fixed opacity-100 transition-opacity duration-500 group-hover/next:opacity-0">
              Press F or double-click to toggle fullscreen.
            </p>
            <p class="fixed opacity-0 transition-opacity duration-500 group-hover/next:opacity-100">
              Next preset
            </p>
          </div>
          <icon-nrk-media-ffw class="size-10 min-w-10 text-muted group-hover/next:text-primary1" />
        </div>
      </div>
    </div>
    <!-- preset name -->
    <div
      v-if="fetchedPresets && fetchedPresets.length > 0 && visualizer"
      class="corner-cut corner-cut z-2 z-4 col-span-full row-span-full mx-auto mb-auto mt-2 overflow-hidden text-ellipsis whitespace-nowrap bg-zshade-300/60 px-4 py-2 text-sm backdrop-blur-xl transition-opacity duration-500 dark:bg-zshade-900/60"
      :class="{
        'opacity-100': presetNameFadeIn,
        'opacity-0': !presetNameFadeIn,
      }"
    >
      preset: {{ fetchedPresets[currentVisualizerIndex].name }}
    </div>
  </div>
</template>
