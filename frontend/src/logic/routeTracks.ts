import type { SubsonicSong } from '~/types/subsonicSong'
import { useLocalStorage } from '@vueuse/core'
import { setCurrentlyPlayingTrack, setCurrentQueue } from '~/logic/playbackQueue'

export const routeTracks = useLocalStorage<SubsonicSong[]>('routeTracks', [] as SubsonicSong[])

export function clearRouteTracks() {
  routeTracks.value = [] as SubsonicSong[]
}

export function setCurrentlyPlayingTrackInRouteTracks(track: SubsonicSong) {
  setCurrentQueue(routeTracks.value)
  setCurrentlyPlayingTrack(track)
}
