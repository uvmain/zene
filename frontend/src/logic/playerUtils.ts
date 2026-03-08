import { postScrobble } from '~/logic/backendFetch'
import { debugLog } from '~/logic/logger'
import { repeatStatus, shuffleEnabled } from '~/logic/store'

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
