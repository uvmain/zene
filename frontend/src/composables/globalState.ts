import type { TrackMetadata } from '../types'

export const currentlyPlayingTrack = ref<TrackMetadata | undefined>()

export function resetCurrentlyPlayingTrack() {
  currentlyPlayingTrack.value = undefined
}
