import { audioElement } from '~/logic/audioElement'
import * as Chromecast from '~/logic/chromecast'
import { volumeStore } from '~/stores/main'

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
  Chromecast.toggleMute()
}

export function changeVolume(volumeString: string) {
  const volume = Number.parseFloat(volumeString)
  currentVolume.value = volume
  if (audioElement.value) {
    audioElement.value.volume = volume
  }
  else if (Chromecast.connected.value) {
    Chromecast.setVolume(volume)
  }
  else {
    console.warn('Audio element not found when trying to change volume')
  }
  volumeStore.value = volumeString
}
