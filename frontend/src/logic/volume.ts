import { audioElement } from '~/logic/audioElement'
import { castPlayerController, chromecastConnected } from '~/logic/castRefs'

export const previousVolume = ref(1)
export const currentVolume = ref(1)

export function toggleMute() {
  if (currentVolume.value > 0) {
    previousVolume.value = currentVolume.value
    changeVolume('0')
  }
  else {
    changeVolume(previousVolume.value.toString())
  }
}

export function changeVolume(volumeString: string) {
  const volume = Number.parseFloat(volumeString)
  currentVolume.value = volume
  if (chromecastConnected.value && castPlayerController.value) {
    castPlayerController.value.setVolumeLevel(volume)
  }
  if (audioElement.value) {
    audioElement.value.volume = volume
  }
}
