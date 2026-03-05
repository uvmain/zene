<script setup lang="ts">
import { chromecastAvailable, cleanupCastPlayer, initializeCast } from '~/logic/chromecast'
import { debugLog } from '~/logic/logger'

const castButton = useTemplateRef('castButton')

// async function castAudio() {
//   const context = cast.framework.CastContext.getInstance()
//   session.value = context.getCurrentSession()

//   if (!session.value) {
//     console.error('No active cast session found')
//     return
//   }

//   if (!trackUrl.value) {
//     console.error('No track URL available for casting')
//     return
//   }

//   // Capture current local playback state and position before switching to cast
//   const wasPlayingLocally = audioPlayer.value?.audioRef && !audioPlayer.value.audioRef.paused
//   if (wasPlayingLocally && audioPlayer.value?.audioRef) {
//     savedLocalPosition.value = audioPlayer.value?.audioRef?.currentTime
//     debugLog(`Captured local position: ${savedLocalPosition.value}s`)
//   }
//   else {
//     savedLocalPosition.value = currentTime.value
//   }

//   isTransitioningToCast.value = true

//   // prefix base url to requestUrl
//   if (window) {
//     const protocol = window.location.protocol
//     const host = window.location.host
//     castUrl.value = `${protocol}//${host}${trackUrl.value}`
//   }

//   const mediaInfo = new chrome.cast.media.MediaInfo(castUrl.value ?? '', 'audio/aac')

//   // Add metadata for better cast experience
//   if (currentlyPlayingTrack.value) {
//     mediaInfo.metadata = {
//       title: currentlyPlayingTrack.value.title,
//       artist: currentlyPlayingTrack.value.artist || currentlyPlayingTrack.value.artist || '',
//       images: [trackArtUrl.value],
//     }
//   }

//   const request = new chrome.cast.media.LoadRequest(mediaInfo)

//   // Set the starting position for cast playback
//   if (savedLocalPosition.value > 0) {
//     request.currentTime = savedLocalPosition.value
//     debugLog(`Setting cast start position to: ${savedLocalPosition.value}s`)
//   }

//   try {
//     await session.value.loadMedia(request)
//     debugLog('Media loaded to cast device')
//     isCasting.value = true

//     // Pause local audio when casting starts
//     if (audioPlayer.value?.audioRef) {
//       audioPlayer.value.audioRef.pause()
//     }

//     // Update cast state
//     updateCastState()

//     // Automatically start cast playback if local audio was playing
//     if (wasPlayingLocally && castPlayerController.value) {
//       debugLog('Auto-starting cast playback since local audio was playing')
//       castPlayerController.value.playOrPause()
//     }

//     isTransitioningToCast.value = false
//   }
//   catch (err) {
//     console.error('Error loading media:', err)
//     isTransitioningToCast.value = false
//   }
// }
function updateStyle(elem: HTMLElement) {
  const shadow = elem.shadowRoot
  if (!shadow)
    return
  const styleNode = shadow.querySelector('style')
  if (styleNode) {
    styleNode.textContent = '.cast_caf_state_c {fill: hsl(32 100% 50%);}.cast_caf_state_d {fill: var(--disconnected-color, #7d7d7d);}.cast_caf_state_h {opacity: 0;}'
  }
}

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

  updateStyle(castButton.value)
})

onUnmounted(() => {
  cleanupCastPlayer()
})
</script>

<template>
  <div v-if="chromecastAvailable" class="flex cursor-pointer items-center justify-center border-none">
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
