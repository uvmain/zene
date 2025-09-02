import type { SubsonicSong } from './subsonicSong'

export interface Queue {
  tracks: SubsonicSong[]
  position: number
}
