import type { TrackMetadataWithImageUrl } from '../types'
import { useLocalStorage } from '@vueuse/core'
import { usePlaybackQueue } from './usePlaybackQueue'

const routeTracks = useLocalStorage<TrackMetadataWithImageUrl[]>('routeTracks', [] as TrackMetadataWithImageUrl[])
const { setCurrentQueue, setCurrentlyPlayingTrackInQueue } = usePlaybackQueue()

export function useRouteTracks() {
  const clearRouteTracks = () => {
    routeTracks.value = [] as TrackMetadataWithImageUrl[]
  }

  const setCurrentlyPlayingTrackInRouteTracks = (track: TrackMetadataWithImageUrl) => {
    setCurrentQueue(routeTracks.value)
    setCurrentlyPlayingTrackInQueue(track)
  }

  return {
    clearRouteTracks,
    routeTracks,
    setCurrentlyPlayingTrackInRouteTracks,
  }
}
