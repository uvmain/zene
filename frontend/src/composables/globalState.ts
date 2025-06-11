import type { Playlist, TrackMetadata, TrackMetadataWithImageUrl } from '../types'
import { trackWithImageUrl } from '../composables/logic'

export const currentlyPlayingTrack = ref<TrackMetadataWithImageUrl | undefined>()

export const currentPlaylist = ref<Playlist | undefined>()

export function resetCurrentlyPlayingTrack() {
  currentlyPlayingTrack.value = undefined
}

export function setCurrentlyPlayingTrack(track: TrackMetadata | TrackMetadataWithImageUrl) {
  console.log(`setting current track to ${track.filename}`)
  currentlyPlayingTrack.value = trackWithImageUrl(track)
}

export function setCurrentPlaylist(tracks: TrackMetadata[] | TrackMetadataWithImageUrl[]) {
  currentPlaylist.value = {
    tracks: tracks.map(track => trackWithImageUrl(track)),
    position: 0,
  }
  console.log(`setting current track to ${tracks[0].filename}`)
  currentlyPlayingTrack.value = trackWithImageUrl(tracks[0])
}
