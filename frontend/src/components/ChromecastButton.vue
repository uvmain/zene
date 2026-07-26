<script setup lang="ts">
import * as ChromeCast from '~/logic/chromecast'
import { debugLog } from '~/logic/logger'

const isChrome = ref(false)

const chromecastConnected = computed(() => {
  return ChromeCast.connected.value
})

const buttonTitle = computed(() => {
  return chromecastConnected.value ? 'Chromecast connected' : 'Connect to Chromecast'
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
    <button :title="buttonTitle" :aria-label="buttonTitle" v-if="chromecastConnected" class="footer-icon-on" @click="ChromeCast.connect">
      <icon-nrk-media-chromecast-active />
    </button>
    <button :title="buttonTitle" :aria-label="buttonTitle" v-else class="footer-icon" @click="ChromeCast.connect">
      <icon-nrk-media-chromecast />
    </button>
  </div>
</template>
