import { postScrobble } from '~/logic/backendFetch'

export const playcountUpdatedMusicbrainzTrackId = ref<string | undefined>()

export async function postPlaycount(musicbrainz_track_id: string): Promise<void> {
  const responseOk = await postScrobble(musicbrainz_track_id)
  if (!responseOk) {
    throw new Error(`Failed to post playcount`)
  }
  playcountUpdatedMusicbrainzTrackId.value = musicbrainz_track_id
}
