<script setup lang="ts">
// import type { SubsonicSong } from '~/types/subsonicSong'
// import { debugLog } from '~/logic/logger'

// const castPlayer = ref<cast.framework.RemotePlayer | null>(null)
// const castPlayerController = ref<cast.framework.RemotePlayerController | null>(null)
// const session = ref<cast.framework.CastSession | null>(null)
// const isCasting = ref(false)
// const castProgressInterval = ref<NodeJS.Timeout | null>(null)

// function initializeCast() {
//   const context = cast.framework.CastContext.getInstance()
//   context.setOptions({
//     receiverApplicationId: chrome.cast.media.DEFAULT_MEDIA_RECEIVER_APP_ID,
//     autoJoinPolicy: chrome.cast.AutoJoinPolicy.ORIGIN_SCOPED,
//   })

//   debugLog('CastContext initialized')
// }

// function toggleMute() {
//   if (isCasting.value && castPlayerController.value) {
//     castPlayerController.value.muteOrUnmute()
//   }
// }

// function setupCastPlayer() {
//   if (!castPlayer.value) {
//     castPlayer.value = new cast.framework.RemotePlayer()
//     castPlayerController.value = new cast.framework.RemotePlayerController(castPlayer.value)

//     castPlayerController.value.addEventListener(cast.framework.CastContextEventType.CAST_STATE_CHANGED, onCastStateChanged)
//     castPlayerController.value.addEventListener(cast.framework.CastContextEventType.SESSION_STATE_CHANGED, onSessionStateChanged)

//     // Listen for remote player events
//     castPlayerController.value.addEventListener(cast.framework.RemotePlayerEventType.IS_PAUSED_CHANGED, onCastPlayerStateChanged)
//     castPlayerController.value.addEventListener(cast.framework.RemotePlayerEventType.CURRENT_TIME_CHANGED, onCastTimeChanged)
//     castPlayerController.value.addEventListener(cast.framework.RemotePlayerEventType.VOLUME_LEVEL_CHANGED, onCastVolumeChanged)
//     castPlayerController.value.addEventListener(cast.framework.RemotePlayerEventType.IS_MUTED_CHANGED, onCastMuteChanged)
//     castPlayerController.value.addEventListener(cast.framework.RemotePlayerEventType.IS_CONNECTED_CHANGED, onCastConnectionChanged)

//     // Listen for media info changes to detect track changes on the remote player
//     castPlayerController.value.addEventListener(cast.framework.RemotePlayerEventType.MEDIA_INFO_CHANGED, onCastMediaInfoChanged)

//     // Start progress tracking for cast
//     startCastProgressTracking()
//   }
// }

// function cleanupCastPlayer() {
//   if (castPlayerController.value) {
//     castPlayerController.value.removeEventListener(cast.framework.RemotePlayerEventType.IS_PAUSED_CHANGED, onCastPlayerStateChanged)
//     castPlayerController.value.removeEventListener(cast.framework.RemotePlayerEventType.CURRENT_TIME_CHANGED, onCastTimeChanged)
//     castPlayerController.value.removeEventListener(cast.framework.RemotePlayerEventType.VOLUME_LEVEL_CHANGED, onCastVolumeChanged)
//     castPlayerController.value.removeEventListener(cast.framework.RemotePlayerEventType.IS_MUTED_CHANGED, onCastMuteChanged)
//     castPlayerController.value.removeEventListener(cast.framework.RemotePlayerEventType.IS_CONNECTED_CHANGED, onCastConnectionChanged)
//     castPlayerController.value.removeEventListener(cast.framework.RemotePlayerEventType.MEDIA_INFO_CHANGED, onCastMediaInfoChanged)
//   }

//   stopCastProgressTracking()
//   castPlayer.value = null
//   castPlayerController.value = null
// }

// function onCastConnectionChanged() {
//   if (castPlayer.value && !castPlayer.value.isConnected) {
//     debugLog('Cast device disconnected unexpectedly')
//     // Capture position before losing connection
//     if (castPlayer.value.currentTime > 0) {
//       savedLocalPosition.value = castPlayer.value.currentTime
//       debugLog(`Position captured on disconnect: ${savedLocalPosition.value}s`)
//     }
//   }
//   // Update cast state to handle potential transition back to local playback
//   updateCastState()
// }

// function startCastProgressTracking() {
//   if (castProgressInterval.value) {
//     clearInterval(castProgressInterval.value)
//   }

//   castProgressInterval.value = setInterval(() => {
//     if (castPlayer.value && castPlayer.value.isConnected) {
//       updateProgress()
//     }
//   }, 1000) // Update every second
// }

// function stopCastProgressTracking() {
//   if (castProgressInterval.value) {
//     clearInterval(castProgressInterval.value)
//     castProgressInterval.value = null
//   }
// }

// function onCastPlayerStateChanged() {
//   if (castPlayer.value) {
//     isPlaying.value = !castPlayer.value.isPaused
//   }
// }

// function onCastTimeChanged() {
//   if (castPlayer.value) {
//     currentTime.value = castPlayer.value.currentTime
//   }
// }

// function onCastVolumeChanged() {
//   if (castPlayer.value) {
//     currentVolume.value = castPlayer.value.volumeLevel
//   }
// }

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

// function updateCastState() {
//   const context = cast.framework.CastContext.getInstance()
//   const currentSession = context.getCurrentSession()

//   // Check if we're transitioning from casting to local
//   if (session.value && !currentSession && isCasting.value) {
//     // Cast session ended, prepare to resume local playback
//     isTransitioningFromCast.value = true
//     debugLog('Cast session ended, preparing to resume local playback')

//     // Capture the last known cast position before cleanup
//     if (castPlayer.value && castPlayer.value.isConnected) {
//       savedLocalPosition.value = castPlayer.value.currentTime
//       debugLog(`Captured cast position: ${savedLocalPosition.value}s`)
//     }
//   }

//   session.value = currentSession

//   if (session.value) {
//     isCasting.value = true
//     setupCastPlayer()
//   }
//   else {
//     const wasPlayingBeforeTransition = isCasting.value && isPlaying.value
//     isCasting.value = false
//     cleanupCastPlayer()

//     // Resume local playback if we were playing before and have a track
//     if (isTransitioningFromCast.value && wasPlayingBeforeTransition && currentlyPlayingTrack.value && audioPlayer.value) {
//       resumeLocalPlayback()
//     }

//     isTransitioningFromCast.value = false
//   }
// }

// function onCastStateChanged(event: any) {
//   debugLog(`Cast state changed:', ${event.castState}`)
//   updateCastState()
// }

// function onSessionStateChanged(event: any) {
//   debugLog(`Session state changed: ${event.sessionState}`)

//   // Handle specific session states for better transition management
//   if (event.sessionState === cast.framework.SessionState.SESSION_ENDED
//     || event.sessionState === cast.framework.SessionState.SESSION_ENDING) {
//     debugLog('Cast session ending/ended')

//     // Capture final cast position if available
//     if (castPlayer.value && castPlayer.value.isConnected) {
//       savedLocalPosition.value = castPlayer.value.currentTime
//       debugLog(`Final cast position captured: ${savedLocalPosition.value}s`)
//     }
//   }

//   updateCastState()
// }

// onMounted(async () => {
//   debugLog('Waiting for Cast SDK...')

//   // Hook for when the SDK becomes available (in case it's not already)
//   window.__onGCastApiAvailable = (isAvailable: boolean) => {
//     debugLog(`Cast API available (async): ${isAvailable}`)
//     if (isAvailable) {
//       initializeCast()
//       updateCastState() // Initialize cast state
//     }
//     else {
//       console.warn('Cast API not available')
//     }
//   }

//   // If the SDK already loaded and called __onGCastApiAvailable BEFORE this script ran
//   // we have to check manually and initialize right now
//   if ((window.cast && window.cast.isAvailable) || (window.chrome?.cast && window.chrome.cast.isAvailable)) {
//     debugLog('Cast API already available (sync)')
//     initializeCast()
//     updateCastState() // Initialize cast state
//   }
// })

// onUnmounted(() => {
//   // Clean up cast resources
//   cleanupCastPlayer()

//   // Remove cast context listeners
//   const context = cast.framework.CastContext.getInstance()
//   if (context) {
//     context.removeEventListener(cast.framework.CastContextEventType.CAST_STATE_CHANGED, onCastStateChanged)
//     context.removeEventListener(cast.framework.CastContextEventType.SESSION_STATE_CHANGED, onSessionStateChanged)
//   }
// })

// function onCastMuteChanged() {
//   if (castPlayer.value) {
//     currentVolume.value = castPlayer.value.isMuted ? 0 : castPlayer.value.volumeLevel
//   }
// }

// async function onCastMediaInfoChanged() {
//   if (!castPlayer.value || !castPlayer.value.isConnected || !castPlayer.value.mediaInfo) {
//     return
//   }

//   debugLog('Cast media info changed, checking for track changes')

//   const remoteMediaUrl = castPlayer.value.mediaInfo.contentId
//   const currentTrackUrl = trackUrl.value

//   // Check if the media URL on the remote player is different from our current track
//   if (remoteMediaUrl && currentTrackUrl && !remoteMediaUrl.includes(currentTrackUrl)) {
//     debugLog('Remote player track changed, syncing local queue')

//     // Extract track ID from the remote media URL
//     // The URL format should be something like: /api/v1/tracks/{track_id}/audio
//     const trackIdMatch = remoteMediaUrl.match(/\/tracks\/([^/]+)\/audio/)
//     if (trackIdMatch && trackIdMatch[1]) {
//       const remoteTrackId = trackIdMatch[1]

//       // Check if this track change matches what we expect from our queue
//       if (currentQueue.value && currentQueue.value.tracks.length > 0) {
//         const nextTrackIndex = currentQueue.value.position + 1
//         const prevTrackIndex = currentQueue.value.position - 1

//         let targetTrack: SubsonicSong | undefined

//         // Check if the remote track matches the next track in our queue
//         if (nextTrackIndex < currentQueue.value.tracks.length) {
//           const nextTrack = currentQueue.value.tracks[nextTrackIndex]
//           if (nextTrack.id === remoteTrackId) {
//             currentQueue.value.position = nextTrackIndex
//             targetTrack = nextTrack
//             debugLog('Remote player advanced to next track in queue')
//           }
//         }

//         // Check if the remote track matches the previous track in our queue
//         if (!targetTrack && prevTrackIndex >= 0) {
//           const prevTrack = currentQueue.value.tracks[prevTrackIndex]
//           if (prevTrack.id === remoteTrackId) {
//             currentQueue.value.position = prevTrackIndex
//             targetTrack = prevTrack
//             debugLog('Remote player went back to previous track in queue')
//           }
//         }

//         // Check if it's a track at the beginning or end (queue wrapping)
//         if (!targetTrack) {
//           for (let i = 0; i < currentQueue.value.tracks.length; i++) {
//             if (currentQueue.value.tracks[i].id === remoteTrackId) {
//               currentQueue.value.position = i
//               targetTrack = currentQueue.value.tracks[i]
//               debugLog(`Remote player jumped to track at position ${i}`)
//               break
//             }
//           }
//         }

//         // Update the currently playing track if we found a match
//         if (targetTrack) {
//           setCurrentlyPlayingTrack(targetTrack)
//           debugLog('Local queue synced with remote player')
//         }
//         else {
//           debugLog('Remote track not found in current queue, may need to handle queue changes')
//         }
//       }
//       else {
//         debugLog('No current queue available for syncing')
//       }
//     }
//     else {
//       debugLog(`Could not extract track ID from remote media URL: ${remoteMediaUrl}`)
//     }
//   }
// }
</script>

<template>
  <div class="inline-block size-22px flex cursor-pointer items-center sm:size-24px">
    <google-cast-launcher />
  </div>
</template>
