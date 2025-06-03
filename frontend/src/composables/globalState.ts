import type { TrackMetadata, TrackMetadataWithImageUrl } from '../types'
import { trackWithImageUrl } from '../composables/logic'

export const currentlyPlayingTrack = ref<TrackMetadataWithImageUrl | undefined>()

export function resetCurrentlyPlayingTrack() {
  currentlyPlayingTrack.value = undefined
}

export function setCurrentlyPlayingTrack(track: TrackMetadata | TrackMetadataWithImageUrl) {
  console.log(`setting current track to ${track.filename}`)
  currentlyPlayingTrack.value = trackWithImageUrl(track)
}
