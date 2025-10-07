import type { SubsonicAlbum } from './subsonicAlbum'
import type { SubsonicArtist } from './subsonicArtist'
import type { SubsonicGenre } from './subsonicGenres'
import type { SubsonicSong } from './subsonicSong'

export interface Queue {
  tracks: SubsonicSong[]
  position: number
}

export type LoadingAttribute = 'lazy' | 'eager'

export interface SearchResult { artists: SubsonicArtist[], albums: SubsonicAlbum[], songs: SubsonicSong[], genres: SubsonicGenre[] }
