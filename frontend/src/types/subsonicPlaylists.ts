import type { SubsonicSong } from './subsonicSong'

export interface SubsonicPlaylist {
  id: number
  name: string
  owner: string
  public: boolean
  songCount: number
  duration: number
  created: string
  changed: string
  coverArt: string
  allowedUser: string[]
  entry: SubsonicSong[]
}

export interface PlaylistGenre {
  name: string
}

export interface PlaylistArtist {
  id: string
  name: string
}

export interface PlaylistAlbumArtist {
  id: string
  name: string
}
