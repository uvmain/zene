import { debugLog } from '~/logic/logger'

export const audioElement = ref<HTMLAudioElement | null>(null)
export const audioNode = ref<AudioNode | null>(null)
export const audioContext = ref<AudioContext | null>(null)

let contextCreated: boolean = false

interface ExtendedWindow extends Window {
  AudioContext?: typeof AudioContext
  webkitAudioContext?: typeof AudioContext
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
  audio.removeEventListener('play', createContextOnPlay)
}
