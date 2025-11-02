<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { useLocalStorage } from '@vueuse/core'
import { getCoverArtUrl } from '~/composables/logic'
import { useDebug } from '~/composables/useDebug'
import { usePlaybackQueue } from '~/composables/usePlaybackQueue'
import { usePlaycounts } from '~/composables/usePlaycounts'
import { useRouteTracks } from '~/composables/useRouteTracks'
import { useSettings } from '~/composables/useSettings'
import PlayerAudio from './PlayerAudio.vue'

const { debugLog } = useDebug()
const { clearQueue, currentlyPlayingTrack, resetCurrentlyPlayingTrack, getNextTrack, getPreviousTrack, getRandomTracks, currentQueue, setCurrentQueue, setCurrentlyPlayingTrack } = usePlaybackQueue()
const { streamQuality } = useSettings()
const { routeTracks } = useRouteTracks()
const { postPlaycount } = usePlaycounts()
const router = useRouter()
const apiKey = useLocalStorage('apiKey', '')

const audioPlayer = useTemplateRef('playerAudio')
const isPlaying = ref(false)
const playcountPosted = ref(false)
const currentTime = ref(0)
const previousVolume = ref(1)
const currentVolume = ref(1)
const isPlayPauseActive = ref(false)
const session = ref<cast.framework.CastSession | null>(null)
const castPlayer = ref<cast.framework.RemotePlayer | null>(null)
const castPlayerController = ref<cast.framework.RemotePlayerController | null>(null)
const isCasting = ref(false)
const castProgressInterval = ref<NodeJS.Timeout | null>(null)
const castUrl = ref<string | null>()
const savedLocalPosition = ref<number>(0)
const isTransitioningToCast = ref<boolean>(false)
const isTransitioningFromCast = ref<boolean>(false)
const showLyrics = ref<boolean>(false)

const trackUrl = computed(() => {
  const queryParamString = `apiKey=${apiKey.value}&c=zene-frontend&v=1.6.0&maxBitRate=${streamQuality.value}&id=${currentlyPlayingTrack.value?.musicBrainzId}&format=aac`
  return currentlyPlayingTrack.value ? `/rest/stream.view?${queryParamString}` : ''
})

const trackArtUrl = computed(() => {
  return currentlyPlayingTrack.value ? getCoverArtUrl(currentlyPlayingTrack.value?.musicBrainzId) : ''
})

async function togglePlayback() {
  if (!audioPlayer.value?.audioRef) {
    console.error('Audio element not found')
    return
  }

  if (!currentQueue.value?.tracks?.length && routeTracks.value.length) {
    setCurrentQueue(routeTracks.value)
  }
  else if (currentQueue.value?.tracks?.length && !currentlyPlayingTrack.value) {
    setCurrentQueue(currentQueue.value?.tracks)
  }

  // Handle casting
  const context = cast.framework.CastContext.getInstance()
  session.value = context.getCurrentSession()

  if (session.value && isCasting.value) {
    // Control cast playback
    if (castPlayerController.value) {
      castPlayerController.value.playOrPause()
    }
  }
  else if (session.value && !isCasting.value) {
    // Start casting
    debugLog('Starting cast session')
    await castAudio()
    return
  }
  else {
    // Control local playback
    if (isPlaying.value) {
      audioPlayer.value?.audioRef.pause()
    }
    else {
      audioPlayer.value?.audioRef.play()
    }
  }

  if (!isPlayPauseActive.value) {
    isPlayPauseActive.value = true
    setTimeout(() => {
      isPlayPauseActive.value = false
    }, 200)
  }

  if (!isCasting.value) {
    updateIsPlaying()
  }
}

function toggleMute() {
  if (isCasting.value && castPlayerController.value) {
    // Control cast device mute
    castPlayerController.value.muteOrUnmute()
  }
  else if (audioPlayer.value?.audioRef) {
    // Control local audio mute
    debugLog('Changing volume')
    if (audioPlayer.value.audioRef.volume !== 0) {
      previousVolume.value = audioPlayer.value.audioRef.volume
      audioPlayer.value.audioRef.volume = 0
      currentVolume.value = 0
    }
    else {
      audioPlayer.value.audioRef.volume = previousVolume.value
      currentVolume.value = previousVolume.value
    }
  }
}

async function stopPlayback() {
  if (isCasting.value && castPlayerController.value) {
    // Stop cast playback
    castPlayerController.value.stop()
  }
  else if (audioPlayer.value?.audioRef) {
    // Stop local playback
    if (audioPlayer.value.audioRef.currentTime < 1) {
      resetCurrentlyPlayingTrack()
      clearQueue()
    }
    audioPlayer.value.audioRef.pause()
    audioPlayer.value.audioRef.load()
  }

  isPlaying.value = false
}

function updateIsPlaying() {
  if (!audioPlayer.value?.audioRef)
    return
  isPlaying.value = !audioPlayer.value.audioRef.paused
}

function updateProgress() {
  // Get progress from cast player if casting, otherwise from local audio
  if (isCasting.value && castPlayer.value) {
    currentTime.value = castPlayer.value.currentTime
  }
  else if (audioPlayer.value?.audioRef) {
    currentTime.value = audioPlayer.value.audioRef.currentTime
  }
  else {
    return
  }

  if (currentlyPlayingTrack.value && !playcountPosted.value) {
    const halfwayPoint = currentlyPlayingTrack.value.duration / 2
    if (currentTime.value >= halfwayPoint) {
      postPlaycount(currentlyPlayingTrack.value.musicBrainzId)
      playcountPosted.value = true
    }
  }
}

watch(currentlyPlayingTrack, (newTrack, oldTrack) => {
  if (newTrack !== oldTrack) {
    playcountPosted.value = false
  }
})

function seek(seekSeconds: number) {
  if (isCasting.value && castPlayer.value && castPlayerController.value) {
    // Seek cast playback
    castPlayer.value.currentTime = seekSeconds
    castPlayerController.value.seek()
  }
  else if (audioPlayer.value?.audioRef) {
    // Seek local playback
    audioPlayer.value.audioRef.currentTime = seekSeconds
  }
}

function volumeInput(volumeString: string) {
  const volume = Number.parseFloat(volumeString)

  if (isCasting.value && castPlayer.value && castPlayerController.value) {
    // Control cast device volume
    castPlayer.value.volumeLevel = volume
    castPlayerController.value.setVolumeLevel()
  }
  else if (audioPlayer.value?.audioRef) {
    // Control local volume
    audioPlayer.value.audioRef.volume = volume
  }

  currentVolume.value = volume
}

async function handleNextTrack() {
  if (currentQueue.value && currentQueue.value.tracks.length > 0) {
    await getNextTrack()

    // If casting, load the new track to the cast device
    if (isCasting.value && currentlyPlayingTrack.value) {
      await castAudio()
    }
  }
  else {
    isPlaying.value = false
  }
}

async function handlePreviousTrack() {
  await getPreviousTrack()

  // If casting, load the new track to the cast device
  if (isCasting.value && currentlyPlayingTrack.value) {
    await castAudio()
  }
}

async function handleGetRandomTracks() {
  await getRandomTracks(500)
  router.push('/queue')
}

watch(currentlyPlayingTrack, (newTrack, oldTrack) => {
  const audio = audioPlayer.value?.audioRef
  if (!audio) {
    return
  }
  if (newTrack && newTrack.musicBrainzId !== oldTrack?.musicBrainzId) {
    audio.pause()
    audio.load()
  }
  else if (!newTrack) {
    audio.pause()
    audio.removeAttribute('src')
    audio.load()
    currentTime.value = 0
    isPlaying.value = false
  }
})

async function castAudio() {
  const context = cast.framework.CastContext.getInstance()
  session.value = context.getCurrentSession()

  if (!session.value) {
    console.error('No active cast session found')
    return
  }

  if (!trackUrl.value) {
    console.error('No track URL available for casting')
    return
  }

  // Capture current local playback state and position before switching to cast
  const wasPlayingLocally = audioPlayer.value?.audioRef && !audioPlayer.value.audioRef.paused
  if (wasPlayingLocally && audioPlayer.value?.audioRef) {
    savedLocalPosition.value = audioPlayer.value?.audioRef?.currentTime
    debugLog(`Captured local position: ${savedLocalPosition.value}s`)
  }
  else {
    savedLocalPosition.value = currentTime.value
  }

  isTransitioningToCast.value = true

  // prefix base url to requestUrl
  if (window) {
    const protocol = window.location.protocol
    const host = window.location.host
    castUrl.value = `${protocol}//${host}${trackUrl.value}`
  }

  const mediaInfo = new chrome.cast.media.MediaInfo(castUrl.value ?? '', 'audio/aac')

  // Add metadata for better cast experience
  if (currentlyPlayingTrack.value) {
    mediaInfo.metadata = {
      title: currentlyPlayingTrack.value.title,
      artist: currentlyPlayingTrack.value.artist || currentlyPlayingTrack.value.artist || '',
      images: [trackArtUrl.value],
    }
  }

  const request = new chrome.cast.media.LoadRequest(mediaInfo)

  // Set the starting position for cast playback
  if (savedLocalPosition.value > 0) {
    request.currentTime = savedLocalPosition.value
    debugLog(`Setting cast start position to: ${savedLocalPosition.value}s`)
  }

  try {
    await session.value.loadMedia(request)
    debugLog('Media loaded to cast device')
    isCasting.value = true

    // Pause local audio when casting starts
    if (audioPlayer.value?.audioRef) {
      audioPlayer.value.audioRef.pause()
    }

    // Update cast state
    updateCastState()

    // Automatically start cast playback if local audio was playing
    if (wasPlayingLocally && castPlayerController.value) {
      debugLog('Auto-starting cast playback since local audio was playing')
      castPlayerController.value.playOrPause()
    }

    isTransitioningToCast.value = false
  }
  catch (err) {
    console.error('Error loading media:', err)
    isTransitioningToCast.value = false
  }
}

function initializeCast() {
  const context = cast.framework.CastContext.getInstance()
  context.setOptions({
    receiverApplicationId: chrome.cast.media.DEFAULT_MEDIA_RECEIVER_APP_ID,
    autoJoinPolicy: chrome.cast.AutoJoinPolicy.ORIGIN_SCOPED,
  })

  debugLog('CastContext initialized')
}

function onCastStateChanged(event: any) {
  debugLog(`Cast state changed:', ${event.castState}`)
  updateCastState()
}

function onSessionStateChanged(event: any) {
  debugLog(`Session state changed: ${event.sessionState}`)

  // Handle specific session states for better transition management
  if (event.sessionState === cast.framework.SessionState.SESSION_ENDED
    || event.sessionState === cast.framework.SessionState.SESSION_ENDING) {
    debugLog('Cast session ending/ended')

    // Capture final cast position if available
    if (castPlayer.value && castPlayer.value.isConnected) {
      savedLocalPosition.value = castPlayer.value.currentTime
      debugLog(`Final cast position captured: ${savedLocalPosition.value}s`)
    }
  }

  updateCastState()
}

function updateCastState() {
  const context = cast.framework.CastContext.getInstance()
  const currentSession = context.getCurrentSession()

  // Check if we're transitioning from casting to local
  if (session.value && !currentSession && isCasting.value) {
    // Cast session ended, prepare to resume local playback
    isTransitioningFromCast.value = true
    debugLog('Cast session ended, preparing to resume local playback')

    // Capture the last known cast position before cleanup
    if (castPlayer.value && castPlayer.value.isConnected) {
      savedLocalPosition.value = castPlayer.value.currentTime
      debugLog(`Captured cast position: ${savedLocalPosition.value}s`)
    }
  }

  session.value = currentSession

  if (session.value) {
    isCasting.value = true
    setupCastPlayer()
  }
  else {
    const wasPlayingBeforeTransition = isCasting.value && isPlaying.value
    isCasting.value = false
    cleanupCastPlayer()

    // Resume local playback if we were playing before and have a track
    if (isTransitioningFromCast.value && wasPlayingBeforeTransition && currentlyPlayingTrack.value && audioPlayer.value) {
      resumeLocalPlayback()
    }

    isTransitioningFromCast.value = false
  }
}

function setupCastPlayer() {
  if (!castPlayer.value) {
    castPlayer.value = new cast.framework.RemotePlayer()
    castPlayerController.value = new cast.framework.RemotePlayerController(castPlayer.value)

    castPlayerController.value.addEventListener(cast.framework.CastContextEventType.CAST_STATE_CHANGED, onCastStateChanged)
    castPlayerController.value.addEventListener(cast.framework.CastContextEventType.SESSION_STATE_CHANGED, onSessionStateChanged)

    // Listen for remote player events
    castPlayerController.value.addEventListener(cast.framework.RemotePlayerEventType.IS_PAUSED_CHANGED, onCastPlayerStateChanged)
    castPlayerController.value.addEventListener(cast.framework.RemotePlayerEventType.CURRENT_TIME_CHANGED, onCastTimeChanged)
    castPlayerController.value.addEventListener(cast.framework.RemotePlayerEventType.VOLUME_LEVEL_CHANGED, onCastVolumeChanged)
    castPlayerController.value.addEventListener(cast.framework.RemotePlayerEventType.IS_MUTED_CHANGED, onCastMuteChanged)
    castPlayerController.value.addEventListener(cast.framework.RemotePlayerEventType.IS_CONNECTED_CHANGED, onCastConnectionChanged)

    // Listen for media info changes to detect track changes on the remote player
    castPlayerController.value.addEventListener(cast.framework.RemotePlayerEventType.MEDIA_INFO_CHANGED, onCastMediaInfoChanged)

    // Start progress tracking for cast
    startCastProgressTracking()
  }
}

function resumeLocalPlayback() {
  if (!audioPlayer.value?.audioRef || !currentlyPlayingTrack.value) {
    debugLog('Cannot resume local playback: missing audio element or track')
    return
  }

  debugLog(`Resuming local playback from position: ${savedLocalPosition.value}s`)

  // Set the position and play
  audioPlayer.value.audioRef.currentTime = savedLocalPosition.value
  currentTime.value = savedLocalPosition.value

  // Use a promise to handle the play() call properly
  audioPlayer.value.audioRef.play().then(() => {
    debugLog('Local playback resumed successfully')
    updateIsPlaying()
  }).catch((error) => {
    console.error('Error resuming local playback:', error)
    isPlaying.value = false
  })
}

function onCastConnectionChanged() {
  if (castPlayer.value && !castPlayer.value.isConnected) {
    debugLog('Cast device disconnected unexpectedly')
    // Capture position before losing connection
    if (castPlayer.value.currentTime > 0) {
      savedLocalPosition.value = castPlayer.value.currentTime
      debugLog(`Position captured on disconnect: ${savedLocalPosition.value}s`)
    }
  }
  // Update cast state to handle potential transition back to local playback
  updateCastState()
}

function cleanupCastPlayer() {
  if (castPlayerController.value) {
    castPlayerController.value.removeEventListener(cast.framework.RemotePlayerEventType.IS_PAUSED_CHANGED, onCastPlayerStateChanged)
    castPlayerController.value.removeEventListener(cast.framework.RemotePlayerEventType.CURRENT_TIME_CHANGED, onCastTimeChanged)
    castPlayerController.value.removeEventListener(cast.framework.RemotePlayerEventType.VOLUME_LEVEL_CHANGED, onCastVolumeChanged)
    castPlayerController.value.removeEventListener(cast.framework.RemotePlayerEventType.IS_MUTED_CHANGED, onCastMuteChanged)
    castPlayerController.value.removeEventListener(cast.framework.RemotePlayerEventType.IS_CONNECTED_CHANGED, onCastConnectionChanged)
    castPlayerController.value.removeEventListener(cast.framework.RemotePlayerEventType.MEDIA_INFO_CHANGED, onCastMediaInfoChanged)
  }

  stopCastProgressTracking()
  castPlayer.value = null
  castPlayerController.value = null
}

function startCastProgressTracking() {
  if (castProgressInterval.value) {
    clearInterval(castProgressInterval.value)
  }

  castProgressInterval.value = setInterval(() => {
    if (castPlayer.value && castPlayer.value.isConnected) {
      updateProgress()
    }
  }, 1000) // Update every second
}

function stopCastProgressTracking() {
  if (castProgressInterval.value) {
    clearInterval(castProgressInterval.value)
    castProgressInterval.value = null
  }
}

function onCastPlayerStateChanged() {
  if (castPlayer.value) {
    isPlaying.value = !castPlayer.value.isPaused
  }
}

function onCastTimeChanged() {
  if (castPlayer.value) {
    currentTime.value = castPlayer.value.currentTime
  }
}

function onCastVolumeChanged() {
  if (castPlayer.value) {
    currentVolume.value = castPlayer.value.volumeLevel
  }
}

function onCastMuteChanged() {
  if (castPlayer.value) {
    currentVolume.value = castPlayer.value.isMuted ? 0 : castPlayer.value.volumeLevel
  }
}

async function onCastMediaInfoChanged() {
  if (!castPlayer.value || !castPlayer.value.isConnected || !castPlayer.value.mediaInfo) {
    return
  }

  debugLog('Cast media info changed, checking for track changes')

  const remoteMediaUrl = castPlayer.value.mediaInfo.contentId
  const currentTrackUrl = trackUrl.value

  // Check if the media URL on the remote player is different from our current track
  if (remoteMediaUrl && currentTrackUrl && !remoteMediaUrl.includes(currentTrackUrl)) {
    debugLog('Remote player track changed, syncing local queue')

    // Extract track ID from the remote media URL
    // The URL format should be something like: /api/v1/tracks/{track_id}/audio
    const trackIdMatch = remoteMediaUrl.match(/\/tracks\/([^/]+)\/audio/)
    if (trackIdMatch && trackIdMatch[1]) {
      const remoteTrackId = trackIdMatch[1]

      // Check if this track change matches what we expect from our queue
      if (currentQueue.value && currentQueue.value.tracks.length > 0) {
        const nextTrackIndex = currentQueue.value.position + 1
        const prevTrackIndex = currentQueue.value.position - 1

        let targetTrack: SubsonicSong | undefined

        // Check if the remote track matches the next track in our queue
        if (nextTrackIndex < currentQueue.value.tracks.length) {
          const nextTrack = currentQueue.value.tracks[nextTrackIndex]
          if (nextTrack.id === remoteTrackId) {
            currentQueue.value.position = nextTrackIndex
            targetTrack = nextTrack
            debugLog('Remote player advanced to next track in queue')
          }
        }

        // Check if the remote track matches the previous track in our queue
        if (!targetTrack && prevTrackIndex >= 0) {
          const prevTrack = currentQueue.value.tracks[prevTrackIndex]
          if (prevTrack.id === remoteTrackId) {
            currentQueue.value.position = prevTrackIndex
            targetTrack = prevTrack
            debugLog('Remote player went back to previous track in queue')
          }
        }

        // Check if it's a track at the beginning or end (queue wrapping)
        if (!targetTrack) {
          for (let i = 0; i < currentQueue.value.tracks.length; i++) {
            if (currentQueue.value.tracks[i].id === remoteTrackId) {
              currentQueue.value.position = i
              targetTrack = currentQueue.value.tracks[i]
              debugLog(`Remote player jumped to track at position ${i}`)
              break
            }
          }
        }

        // Update the currently playing track if we found a match
        if (targetTrack) {
          setCurrentlyPlayingTrack(targetTrack)
          debugLog('Local queue synced with remote player')
        }
        else {
          debugLog('Remote track not found in current queue, may need to handle queue changes')
        }
      }
      else {
        debugLog('No current queue available for syncing')
      }
    }
    else {
      debugLog(`Could not extract track ID from remote media URL: ${remoteMediaUrl}`)
    }
  }
}

onMounted(async () => {
  debugLog('Waiting for Cast SDK...')

  // Hook for when the SDK becomes available (in case it's not already)
  window.__onGCastApiAvailable = (isAvailable: boolean) => {
    debugLog(`Cast API available (async): ${isAvailable}`)
    if (isAvailable) {
      initializeCast()
      updateCastState() // Initialize cast state
    }
    else {
      console.warn('Cast API not available')
    }
  }

  // If the SDK already loaded and called __onGCastApiAvailable BEFORE this script ran
  // we have to check manually and initialize right now
  if ((window.cast && window.cast.isAvailable) || (window.chrome?.cast && window.chrome.cast.isAvailable)) {
    debugLog('Cast API already available (sync)')
    initializeCast()
    updateCastState() // Initialize cast state
  }
})

onUnmounted(() => {
  // Clean up cast resources
  cleanupCastPlayer()

  // Remove cast context listeners
  const context = cast.framework.CastContext.getInstance()
  if (context) {
    context.removeEventListener(cast.framework.CastContextEventType.CAST_STATE_CHANGED, onCastStateChanged)
    context.removeEventListener(cast.framework.CastContextEventType.SESSION_STATE_CHANGED, onSessionStateChanged)
  }
})
</script>

<template>
  <footer
    class="sticky bottom-0 mt-auto w-full border-0 border-t-1 border-zshade-400 border-zshade-600 border-solid background-2"
  >
    <div class="flex flex-col items-center px-2 md:flex-row space-y-2 md:px-4 md:space-x-2 md:space-y-0">
      <div
        class="h-full w-full flex flex-grow flex-col items-center justify-center py-2 space-y-2 md:py-2 md:space-y-2"
      >
        <PlayerAudio
          ref="playerAudio"
          :track-url="trackUrl"
          @play="updateIsPlaying()"
          @pause="updateIsPlaying()"
          @time-update="updateProgress"
          @ended="handleNextTrack()"
        />
        <div>
          <PlayerProgressBar
            :current-time-in-seconds="currentTime"
            :currently-playing-track="currentlyPlayingTrack"
            @seek="seek"
          />

          <PlayerMediaControls
            :is-playing="isPlaying"
            @stop-playback="stopPlayback()"
            @toggle-playback="togglePlayback()"
            @next-track="handleNextTrack()"
            @previous-track="handlePreviousTrack()"
            @get-random-tracks="handleGetRandomTracks()"
          />
        </div>
      </div>

      <div class="flex flex-row items-center gap-x-3 md:gap-x-4">
        <PlayerCastButton />
        <PlayerLyricsButton
          :track="currentlyPlayingTrack"
          :is-active="showLyrics"
          :current-time="currentTime"
          :currently-playing-track="currentlyPlayingTrack"
          @toggle-lyrics="showLyrics = !showLyrics"
        />
        <PlayerPlaylistButton />
        <PlayerVolumeSlider
          :audio-ref="audioPlayer?.audioRef"
          :model-value="currentVolume"
          @toggle-mute="toggleMute()"
          @update:model-value="volumeInput"
        />
      </div>
    </div>
  </footer>
</template>
