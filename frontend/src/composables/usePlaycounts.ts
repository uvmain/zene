import { useBackendFetch } from './useBackendFetch'

const { backendFetchRequest } = useBackendFetch()

const playcount_updated_musicbrainz_track_id = ref<string | undefined>()

export function usePlaycounts() {
  const postPlaycount = async (musicbrainz_track_id: string): Promise<void> => {
    const formData = new FormData()
    formData.append('musicbrainz_track_id', musicbrainz_track_id)

    const response = await backendFetchRequest('playcounts', {
      body: formData,
    })
    if (!response.ok) {
      throw new Error(`Failed to post playcount: ${response.statusText}`)
    }
    playcount_updated_musicbrainz_track_id.value = musicbrainz_track_id
  }

  return {
    postPlaycount,
    playcount_updated_musicbrainz_track_id,
  }
}
