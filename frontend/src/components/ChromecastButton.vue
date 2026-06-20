<script setup lang="ts">
import * as ChromeCast from '~/logic/chromecast'
import { debugLog } from '~/logic/logger'

const isChrome = ref(false)

const chromecastDuration = computed(() => {
  return ChromeCast.duration.value
})

const chromecastConnected = computed(() => {
  return ChromeCast.connected.value
})

onMounted(() => {
  isChrome.value = ChromeCast.isBrowserChrome()
  if (isChrome.value) {
    debugLog('[ChromecastButton] Browser is Chrome, initialising Chromecast')
    ChromeCast.initialiseChromecast()
  }
  else {
    debugLog('[ChromecastButton] Browser is not Chrome, Chromecast will not be available')
  }
})
</script>

<template>
  <div v-if="isChrome" class="flex items-center justify-center">
    <button v-if="chromecastConnected" class="footer-icon-on" @click="ChromeCast.connect">
      <icon-nrk-media-chromecast-active />
    </button>
    <button v-else class="footer-icon" @click="ChromeCast.connect">
      <icon-nrk-media-chromecast />
    </button>
    <button v-if="chromecastConnected" @click="ChromeCast.loadMedia">
      Load Media
    </button>
    <button v-if="chromecastConnected" @click="ChromeCast.stop">
      Stop
    </button>
  </div>

  <div v-if="chromecastConnected">
    <button v-if="ChromeCast.playing" @click="ChromeCast.pause">
      pause_arrow
    </button>
    <button v-else @click="ChromeCast.play">
      play_arrow
    </button>
    <input type="range" step="any" min="0" :max="chromecastDuration" :value="ChromeCast.currentTime" @change="ChromeCast.onSeekChange">
    <button v-if="ChromeCast.muted" @click="ChromeCast.toggleMute">
      volume_mute
    </button>
    <button v-else @click="ChromeCast.toggleMute">
      volume_up
    </button>
    <input
      type="range"
      step="any"
      min="0"
      max="1"
      :value="ChromeCast.volume.value"
      @change="ChromeCast.onVolumeChange"
    >
  </div>
</template>
