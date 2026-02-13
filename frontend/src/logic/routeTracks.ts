import type { SubsonicSong } from '~/types/subsonicSong'
import { setCurrentlyPlayingTrack, setCurrentQueue } from '~/logic/playbackQueue'
import { routeTracks } from '~/logic/store'

export function clearRouteTracks() {
  routeTracks.value = [] as SubsonicSong[]
}

export function setCurrentlyPlayingTrackInRouteTracks(track: SubsonicSong) {
  setCurrentQueue(routeTracks.value)
  setCurrentlyPlayingTrack(track)
}
