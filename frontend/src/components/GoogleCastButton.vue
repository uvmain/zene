<script setup lang="ts">
import { cleanupCastPlayer, initializeCast } from '~/logic/castFunctions'
import { castPlayer, chromecastAvailable } from '~/logic/castRefs'
import { debugLog } from '~/logic/logger'
import { currentVolume } from '~/logic/volume'

const castButton = useTemplateRef('castButton')

watchEffect(() => {
  if (castButton.value && castButton.value.shadowRoot) {
    const shadow = castButton.value.shadowRoot
    if (!shadow)
      return
    const styleNode = shadow.querySelector('style')
    if (styleNode) {
      styleNode.textContent = '.cast_caf_state_c {fill: hsl(32 100% 50%);}.cast_caf_state_d {fill: var(--disconnected-color, #7d7d7d);}.cast_caf_state_h {opacity: 0;}'
    }
  }
})

watchEffect(() => {
  if (castPlayer.value && castPlayer.value.volumeLevel !== currentVolume.value) {
    castPlayer.value.volumeLevel = currentVolume.value
  }
})

onMounted(async () => {
  debugLog('Waiting for Cast SDK...')
  window.__onGCastApiAvailable = (isAvailable: boolean) => {
    if (isAvailable) {
      debugLog('Cast API available')
      chromecastAvailable.value = true
      initializeCast()
    }
    else {
      debugLog('Cast API unavailable')
    }
  }

  // If the SDK already loaded, manually check and initialize
  if ((window.cast && window.cast.framework) || (window.chrome?.cast && window.chrome.cast.isAvailable)) {
    debugLog('Cast API already available')
    chromecastAvailable.value = true
    initializeCast()
  }
})

onUnmounted(() => {
  cleanupCastPlayer()
})
</script>

<template>
  <div v-if="chromecastAvailable" class="border-none flex cursor-pointer items-center justify-center">
    <google-cast-launcher ref="castButton" class="size-7" />
  </div>
</template>

<style lang="css" scoped>
:root {
  --connected-color: rgb(209, 147, 12);
}
#google-cast-launcher .cast_caf_state_c {
  fill: var(--connected-color);
}
</style>
