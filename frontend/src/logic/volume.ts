import { audioElement } from '~/logic/audioElement'
import * as Chromecast from '~/logic/chromecast'
import { volumeStore } from '~/stores/main'
import { debugLog } from './logger'

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
  if (Chromecast.connected.value) {
    Chromecast.setVolume(volume)
  }
  if (!audioElement.value && !Chromecast.connected.value) {
    debugLog('No audio element or Chromecast connected, cannot change volume')
  }
  volumeStore.value = volumeString
}
