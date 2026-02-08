import { postScrobble } from './backendFetch'

const playcountUpdatedMusicbrainzTrackId = ref<string | undefined>()

export function usePlaycounts() {
  const postPlaycount = async (musicbrainz_track_id: string): Promise<void> => {
    const responseOk = await postScrobble(musicbrainz_track_id)
    if (!responseOk) {
      throw new Error(`Failed to post playcount`)
    }
    playcountUpdatedMusicbrainzTrackId.value = musicbrainz_track_id
  }

  return {
    postPlaycount,
    playcountUpdatedMusicbrainzTrackId,
  }
}
