import type { TrackMetadata, TrackMetadataWithImageUrl } from '../types'
import { backendFetchRequest } from './fetchFromBackend'
import { trackWithImageUrl } from './logic'

export async function getRandomTrack(): Promise<TrackMetadataWithImageUrl> {
  const response = await backendFetchRequest('tracks?random=true&limit=1')
  const json = await response.json() as TrackMetadata[]
  return trackWithImageUrl(json[0])
}
