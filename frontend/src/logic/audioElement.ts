import type { PlayItem } from '~/types'
import * as Chromecast from '~/logic/chromecast'
import { debugLog } from '~/logic/logger'

export const audioElement = ref<HTMLAudioElement | null>(null)
export const audioNode = ref<AudioNode | null>(null)
export const audioContext = ref<AudioContext | null>(null)
const previousPlayItem = ref<PlayItem | null>(null)

let contextCreated: boolean = false

interface ExtendedWindow extends Window {
  AudioContext?: typeof AudioContext
  webkitAudioContext?: typeof AudioContext
}

function isSamePlayItem(playItem1: PlayItem | null, playItem2: PlayItem): boolean {
  if (playItem1?.track && playItem2.track) {
    return playItem1.track.musicBrainzId === playItem2.track.musicBrainzId
  }

  if (playItem1?.podcastEpisode && playItem2.podcastEpisode) {
    return playItem1.podcastEpisode.streamId === playItem2.podcastEpisode.streamId
  }

  return false
}

function setElementSource(audio: HTMLAudioElement, src: string) {
  if (audio.getAttribute('src') === src) {
    return
  }

  audio.setAttribute('src', src)
}

function cleanElementSource(audio: HTMLAudioElement) {
  audio.pause()
  audio.currentTime = 0
  audio.removeAttribute('src')
  audio.load()
}

async function waitForCanPlayThrough(audio: HTMLAudioElement): Promise<boolean> {
  if (audio.readyState >= HTMLMediaElement.HAVE_ENOUGH_DATA) {
    return true
  }

  return new Promise((resolve) => {
    function onCanPlayThrough() {
      cleanup()
      resolve(true)
    }

    function onError() {
      cleanup()
      resolve(false)
    }

    function cleanup() {
      audio.removeEventListener('canplaythrough', onCanPlayThrough)
      audio.removeEventListener('error', onError)
    }

    audio.addEventListener('canplaythrough', onCanPlayThrough, { once: true })
    audio.addEventListener('error', onError, { once: true })
  })
}

// One-time play event to create AudioContext after user interaction
export function createContextOnPlay() {
  const audio = audioElement.value
  if (!audio) {
    return
  }
  if (!contextCreated && typeof window !== 'undefined') {
    const extendedWindow = window as ExtendedWindow
    const AudioContextConstructor = extendedWindow.AudioContext || extendedWindow.webkitAudioContext
    if (AudioContextConstructor) {
      audioContext.value = new AudioContextConstructor()
    }
    if (audioContext.value) {
      audioNode.value = audioContext.value.createMediaElementSource(audio)
      audioNode.value.connect(audioContext.value.destination)
      contextCreated = true
      debugLog('Audio context created')
    }
    else {
      debugLog('Failed to create audio context')
    }
  }
}

export async function playWhenReady(playItem: PlayItem, src: string): Promise<boolean> {
  if (Chromecast.connected.value) {
    const mediaUrl = src
    await Chromecast.loadMedia(mediaUrl)
    void Chromecast.play()
    previousPlayItem.value = playItem
    return true
  }

  const audio = audioElement.value
  if (!audio) {
    return false
  }

  const shouldReload = !isSamePlayItem(previousPlayItem.value, playItem) || audio.getAttribute('src') !== src

  if (shouldReload) {
    audio.pause()
    setElementSource(audio, src)
    audio.load()
  }

  const ready = await waitForCanPlayThrough(audio)
  if (!ready) {
    return false
  }

  await audio.play()
  previousPlayItem.value = playItem
  return true
}

export function clearActiveAudio() {
  const audio = audioElement.value
  if (!audio) {
    return
  }

  cleanElementSource(audio)
  previousPlayItem.value = null
}

export function seek(seekSeconds: number) {
  if (audioElement.value) {
    audioElement.value.currentTime = seekSeconds
  }
}
