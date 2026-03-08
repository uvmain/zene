import { postScrobble } from '~/logic/backendFetch'
import { audioElement, currentVolume, previousVolume } from '~/logic/playbackQueue'
import { repeatStatus, shuffleEnabled } from '~/logic/store'
import { debugLog } from './logger'

export const playcountUpdatedMusicbrainzTrackId = ref<string | undefined>()

export async function postPlaycount(musicbrainz_track_id: string): Promise<void> {
  const responseOk = await postScrobble(musicbrainz_track_id)
  if (!responseOk) {
    debugLog(`Failed to post playcount for ${musicbrainz_track_id}`)
  }
  playcountUpdatedMusicbrainzTrackId.value = musicbrainz_track_id
}

export function toggleShuffle() {
  shuffleEnabled.value = !shuffleEnabled.value
}

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
  if (!audioElement.value) {
    return
  }
  const volume = Number.parseFloat(volumeString)
  audioElement.value.volume = volume
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
}
