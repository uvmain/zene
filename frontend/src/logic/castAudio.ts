import { audioElement } from '~/logic/audioElement'
import { castPlayerController, castSession, isCasting, savedLocalPosition } from '~/logic/castRefs'
import { getCoverArtUrl } from '~/logic/common'
import { debugLog } from '~/logic/logger'
import { currentlyPlayingItem, currentTime, trackUrl } from '~/logic/playbackQueue'

let castUrl = ''

export async function castAudio() {
  if (!castSession.value) {
    console.error('No active cast session found')
    return
  }

  if (trackUrl.value === undefined) {
    console.error('No track URL available for casting')
    return
  }

  // Capture current local playback state and position before switching to cast
  const wasPlayingLocally = audioElement.value && !audioElement.value.paused
  if (wasPlayingLocally && audioElement.value) {
    savedLocalPosition.value = audioElement.value.currentTime
    debugLog(`Captured local position: ${savedLocalPosition.value}s`)
  }
  else {
    savedLocalPosition.value = currentTime.value
  }

  // isTransitioningToCast.value = true

  // prefix base url to requestUrl
  if (typeof window !== 'undefined') {
    const protocol = window.location.protocol
    const host = window.location.host
    castUrl = `${protocol}//${host}${trackUrl.value}`
  }
  debugLog(`cast URL: ${castUrl}`)

  const mediaInfo = new chrome.cast.media.MediaInfo(castUrl ?? '', 'audio/aac')

  // add metadata
  if (currentlyPlayingItem.value.track) {
    mediaInfo.metadata = {
      title: currentlyPlayingItem.value.track.title,
      artist: currentlyPlayingItem.value.track.artist || currentlyPlayingItem.value.track.artist || '',
      images: [
        getCoverArtUrl(currentlyPlayingItem.value.track?.musicBrainzId),
      ],
    }
  }

  const request = new chrome.cast.media.LoadRequest(mediaInfo)

  // Set the starting position for cast playback
  if (savedLocalPosition.value > 0) {
    request.currentTime = savedLocalPosition.value
    debugLog(`Setting cast start position to: ${savedLocalPosition.value}s`)
  }

  try {
    await castSession.value.loadMedia(request)
    debugLog('Media loaded to cast device')
    isCasting.value = true

    // Pause local audio when casting starts
    if (audioElement.value) {
      audioElement.value.pause()
    }

    // // Update cast state
    // updateCastState()

    // Automatically start cast playback if local audio was playing
    if (wasPlayingLocally && castPlayerController.value) {
      debugLog('Auto-starting cast playback since local audio was playing')
      castPlayerController.value.playOrPause()
    }

    // isTransitioningToCast.value = false
  }
  catch (err) {
    console.error('Error loading media:', err)
    // isTransitioningToCast.value = false
  }
}
