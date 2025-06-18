import type { StandardResponse } from '../types'
import { useBackendFetch } from './useBackendFetch'
import { usePlaybackQueue } from './usePlaybackQueue'
import { useRouteTracks } from './useRouteTracks'

const { backendFetchRequest } = useBackendFetch()
const { currentQueue } = usePlaybackQueue()
const { routeTracks } = useRouteTracks()

const last_updated_musicbrainz_track_id = ref<string | undefined>()

export function usePlaycounts() {
  const postPlaycount = async (musicbrainz_track_id: string): Promise<StandardResponse> => {
    const formData = new FormData()
    formData.append('musicbrainz_track_id', musicbrainz_track_id)

    const response = await backendFetchRequest('playcounts', {
      body: formData,
      method: 'POST',
    })

    const json = await response.json() as StandardResponse
    return json
  }

  const updatePlaycount = (musicbrainz_track_id: string) => {
    if (currentQueue.value && currentQueue.value.tracks.length) {
      currentQueue.value.tracks.forEach((t) => {
        if (t.musicbrainz_track_id === musicbrainz_track_id) {
          t.user_play_count = t.user_play_count + 1
          t.global_play_count = t.global_play_count + 1
        }
      })
    }
    else if (routeTracks.value.length > 0) {
      routeTracks.value.forEach((t) => {
        if (t.musicbrainz_track_id === musicbrainz_track_id) {
          t.user_play_count = t.user_play_count + 1
          t.global_play_count = t.global_play_count + 1
        }
      })
    }
  }

  return {
    postPlaycount,
    updatePlaycount,
    last_updated_musicbrainz_track_id,
  }
}
