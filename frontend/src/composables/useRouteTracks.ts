import type { SubsonicSong } from '~/types/subsonicSong'
import { useLocalStorage } from '@vueuse/core'
import { usePlaybackQueue } from './usePlaybackQueue'

const routeTracks = useLocalStorage<SubsonicSong[]>('routeTracks', [] as SubsonicSong[])
const { setCurrentQueue, setCurrentlyPlayingTrackInQueue } = usePlaybackQueue()

export function useRouteTracks() {
  const clearRouteTracks = () => {
    routeTracks.value = [] as SubsonicSong[]
  }

  const setCurrentlyPlayingTrackInRouteTracks = (track: SubsonicSong) => {
    setCurrentQueue(routeTracks.value)
    setCurrentlyPlayingTrackInQueue(track)
  }

  return {
    clearRouteTracks,
    routeTracks,
    setCurrentlyPlayingTrackInRouteTracks,
  }
}
