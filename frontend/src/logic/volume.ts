import { audioElement } from '~/logic/audioElement'
import { castPlayer, castPlayerController } from '~/logic/castRefs'
import { debugLog } from '~/logic/logger'

export const previousVolume = ref(1)
export const currentVolume = ref(1)

export function toggleMute() {
  if (audioElement.value) {
    debugLog('Changing volume')
    if (audioElement.value.volume !== 0) {
      previousVolume.value = audioElement.value.volume
      audioElement.value.volume = 0
      currentVolume.value = 0
    }
    else {
      audioElement.value.volume = previousVolume.value
      currentVolume.value = previousVolume.value
    }
  }
}

export function changeVolume(volumeString: string) {
  if (castPlayer.value && castPlayerController.value) {
    castPlayerController.value.setVolumeLevel(Number.parseFloat(volumeString))
    return
  }
  if (!audioElement.value) {
    return
  }
  const volume = Number.parseFloat(volumeString)
  audioElement.value.volume = volume
  currentVolume.value = volume
}
