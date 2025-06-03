import type { TrackMetadata, TrackMetadataWithImageUrl } from '../types'
import { trackWithImageUrl } from '../composables/logic'

export async function backendFetchRequest(path: string, options = {}): Promise<Response> {
  const url = `/api/${path}`
  const response = await fetch(url, options)
  return response
}

export async function getAlbumTracks(musicbrainz_album_id: string): Promise<TrackMetadataWithImageUrl[]> {
  const response = await backendFetchRequest(`albums/${musicbrainz_album_id}/tracks`)
  const json = await response.json() as TrackMetadata[]
  const trackArray: TrackMetadataWithImageUrl[] = []
  json.forEach((track) => {
    trackArray.push(trackWithImageUrl(track))
  })
  return trackArray
}
