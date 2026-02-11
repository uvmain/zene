import { debugLog } from '~/logic/logger'
import { clearQueue, currentlyPlayingPodcastEpisode, currentlyPlayingTrack, currentQueue, getNextTrack, getPreviousShuffledTrack, getPreviousTrack, getShuffledTrack, resetCurrentlyPlayingTrack, setCurrentQueue } from '~/logic/playbackQueue'
import { postPlaycount } from '~/logic/playCounts'
import { routeTracks } from '~/logic/routeTracks'

export const isPlaying = ref(false)
export const playcountPosted = ref(false)
export const currentTime = ref(0)
export const previousVolume = ref(1)
export const currentVolume = ref(1)
export const trackUrl = ref('')
export const shuffleEnabled = ref(false)
export const repeatStatus = ref<'off' | '1' | 'all'>('off')

type AudioElement = HTMLAudioElement | null | undefined

export function toggleShuffle() {
  shuffleEnabled.value = !shuffleEnabled.value
}

export function toggleMute(audioElement: AudioElement) {
  if (audioElement) {
    debugLog('Changing volume')
    if (audioElement.volume !== 0) {
      previousVolume.value = audioElement.volume
      audioElement.volume = 0
      currentVolume.value = 0
    }
    else {
      audioElement.volume = previousVolume.value
      currentVolume.value = previousVolume.value
    }
  }
}

export async function togglePlayback(audioElement: AudioElement) {
  if (!audioElement) {
    console.error('Audio element not found')
    return
  }
  else if (!(currentQueue.value && currentQueue.value.tracks.length) && routeTracks.value.length) {
    setCurrentQueue(routeTracks.value)
  }
  else if ((currentQueue.value && currentQueue.value.tracks.length) && !currentlyPlayingTrack.value) {
    setCurrentQueue(currentQueue.value?.tracks)
  }

  if (isPlaying.value) {
    audioElement.pause()
    isPlaying.value = false
  }
  else {
    await audioElement.play()
    isPlaying.value = true
  }
}

export async function stopPlayback(audioElement: AudioElement) {
  if (audioElement) {
    if (audioElement.currentTime < 1) {
      resetCurrentlyPlayingTrack()
      clearQueue()
    }
    audioElement.pause()
    audioElement.load()
  }

  isPlaying.value = false
}

export async function handleNextTrack() {
  if (currentQueue.value && currentQueue.value.tracks.length > 0) {
    shuffleEnabled.value ? await getShuffledTrack() : await getNextTrack()
  }
}

export async function handlePreviousTrack() {
  shuffleEnabled.value ? await getPreviousShuffledTrack() : await getPreviousTrack()
}

export function seek(audioElement: AudioElement, seekSeconds: number) {
  if (audioElement) {
    audioElement.currentTime = seekSeconds
  }
}

export async function updateProgress(audioElement: AudioElement) {
  if (!audioElement) {
    return
  }
  currentTime.value = audioElement.currentTime

  if (currentlyPlayingTrack.value && !playcountPosted.value) {
    const halfwayPoint = currentlyPlayingTrack.value.duration / 2
    if (currentTime.value >= halfwayPoint) {
      await postPlaycount(currentlyPlayingTrack.value.musicBrainzId)
      playcountPosted.value = true
    }
  }
  else if (currentlyPlayingPodcastEpisode.value && !playcountPosted.value) {
    const halfwayPoint = Number(currentlyPlayingPodcastEpisode.value.duration) / 2
    if (currentTime.value >= halfwayPoint) {
      await postPlaycount(currentlyPlayingPodcastEpisode.value.streamId)
      playcountPosted.value = true
    }
  }
}

export function volumeInput(audioElement: AudioElement, volumeString: string) {
  const volume = Number.parseFloat(volumeString)

  if (audioElement) {
    audioElement.volume = volume
  }

  currentVolume.value = volume
}

export function toggleRepeat() {
  switch (repeatStatus.value) {
    case 'off':
      repeatStatus.value = '1'
      break
    case '1':
      repeatStatus.value = 'all'
      break
    case 'all':
      repeatStatus.value = 'off'
      break
  }
  console.log('Toggle repeat')
}
