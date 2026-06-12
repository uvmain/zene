<script setup lang="ts">
import * as ChromeCast from '~/logic/chromecast'
import { debugLog } from '~/logic/logger'

const isChrome = ref(false)

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
    <button v-if="ChromeCast.connected.value" class="footer-icon-on" @click="ChromeCast.connect">
      <icon-nrk-media-chromecast-active />
    </button>
    <button v-else class="footer-icon" @click="ChromeCast.connect">
      <icon-nrk-media-chromecast />
    </button>
    <button v-if="ChromeCast.connected.value" @click="ChromeCast.loadMedia">
      Load Media
    </button>
    <button v-if="ChromeCast.connected.value" @click="ChromeCast.stop">
      Stop
    </button>
  </div>

  <div v-if="ChromeCast.connected.value">
    <button v-if="ChromeCast.playing.value" @click="ChromeCast.pause">
      pause_arrow
    </button>
    <button v-else @click="ChromeCast.play">
      play_arrow
    </button>
    <input type="range" step="any" min="0" :max="ChromeCast.duration.value" :value="ChromeCast.currentTime.value" @change="ChromeCast.onSeekChange">
    <button v-if="ChromeCast.muted.value" @click="ChromeCast.setMute">
      volume_mute
    </button>
    <button v-else @click="ChromeCast.setMute">
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
