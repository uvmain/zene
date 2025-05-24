import { backendFetchRequest} from './fetchFromBackend'
import type { TrackMetadata } from '../types'

export async function getRandomTrack(): Promise<TrackMetadata> {
  const response = await backendFetchRequest('tracks?random=true&limit=1')
  const json = await response.json()
  return json[0]
}