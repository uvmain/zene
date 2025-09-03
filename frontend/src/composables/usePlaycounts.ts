import { postScrobble } from './backendFetch'

const playcount_updated_musicbrainz_track_id = ref<string | undefined>()

export function usePlaycounts() {
  const postPlaycount = async (musicbrainz_track_id: string): Promise<void> => {
    const responseOk = await postScrobble(musicbrainz_track_id)
    if (!responseOk) {
      throw new Error(`Failed to post playcount`)
    }
    playcount_updated_musicbrainz_track_id.value = musicbrainz_track_id
  }

  return {
    postPlaycount,
    playcount_updated_musicbrainz_track_id,
  }
}
