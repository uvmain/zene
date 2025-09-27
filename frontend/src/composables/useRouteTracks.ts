import type { SubsonicSong } from '~/types/subsonicSong'
import { useLocalStorage } from '@vueuse/core'
import { usePlaybackQueue } from './usePlaybackQueue'

const routeTracks = useLocalStorage<SubsonicSong[]>('routeTracks', [] as SubsonicSong[])
const { setCurrentQueue, setCurrentlyPlayingTrack } = usePlaybackQueue()

export function useRouteTracks() {
  const clearRouteTracks = () => {
    routeTracks.value = [] as SubsonicSong[]
  }

  const setCurrentlyPlayingTrackInRouteTracks = (track: SubsonicSong) => {
    setCurrentQueue(routeTracks.value)
    setCurrentlyPlayingTrack(track)
  }

  return {
    clearRouteTracks,
    routeTracks,
    setCurrentlyPlayingTrackInRouteTracks,
  }
}
