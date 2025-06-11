import type { TrackMetadata, TrackMetadataWithImageUrl } from '../types'
import { currentPlaylist } from '../composables/globalState'
import { backendFetchRequest } from './fetchFromBackend'
import { trackWithImageUrl } from './logic'

export async function getRandomTrack(): Promise<TrackMetadataWithImageUrl> {
  const response = await backendFetchRequest('tracks?random=true&limit=1')
  const json = await response.json() as TrackMetadata[]
  return trackWithImageUrl(json[0])
}

export async function getNextTrack(): Promise<TrackMetadataWithImageUrl> {
  if (currentPlaylist.value && currentPlaylist.value.tracks.length) {
    const currentIndex = currentPlaylist.value.position
    let nextTrack: TrackMetadataWithImageUrl
    if (currentIndex < currentPlaylist.value.tracks.length - 1) {
      nextTrack = currentPlaylist.value.tracks[currentIndex + 1]
      currentPlaylist.value.position = currentPlaylist.value.position + 1
    }
    else {
      nextTrack = currentPlaylist.value.tracks[0]
      currentPlaylist.value.position = 1
    }
    return nextTrack
  }
  else {
    const response = await backendFetchRequest('tracks?random=true&limit=1')
    const json = await response.json() as TrackMetadata[]
    return trackWithImageUrl(json[0])
  }
}

export async function getPreviousTrack(): Promise<TrackMetadataWithImageUrl> {
  if (currentPlaylist.value && currentPlaylist.value.tracks.length) {
    const currentIndex = currentPlaylist.value.position
    let nextTrack: TrackMetadataWithImageUrl
    if (currentIndex > 0) {
      nextTrack = currentPlaylist.value.tracks[currentIndex - 1]
      currentPlaylist.value.position = currentPlaylist.value.position - 1
    }
    else {
      nextTrack = currentPlaylist.value.tracks[currentPlaylist.value.tracks.length - 1]
      currentPlaylist.value.position = currentPlaylist.value.tracks.length - 1
    }
    return nextTrack
  }
  else {
    const response = await backendFetchRequest('tracks?random=true&limit=1')
    const json = await response.json() as TrackMetadata[]
    return trackWithImageUrl(json[0])
  }
}
