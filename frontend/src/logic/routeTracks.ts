import type { SubsonicSong } from '~/types/subsonicSong'
import { setCurrentlyPlayingTrack, setCurrentQueue } from '~/logic/playbackQueue'

export const routeTracks = ref<SubsonicSong[]>([])

export function clearRouteTracks() {
  routeTracks.value = [] as SubsonicSong[]
}

export function setCurrentlyPlayingTrackInRouteTracks(track: SubsonicSong) {
  setCurrentQueue(routeTracks.value)
  setCurrentlyPlayingTrack(track)
}
