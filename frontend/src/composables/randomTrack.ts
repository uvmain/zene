import type { TrackMetadata } from '../types'
import { backendFetchRequest } from './fetchFromBackend'

export async function getRandomTrack(): Promise<TrackMetadata> {
  const response = await backendFetchRequest('tracks?random=true&limit=1')
  const json = await response.json() as TrackMetadata[]
  return json[0]
}
