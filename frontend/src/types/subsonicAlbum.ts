import type { SubsonicSong } from './subsonicSong'

export interface SubsonicAlbum {
  id: string
  title: string
  name: string
  artist: string
  artistId: string
  coverArt: string
  songCount: number
  duration: number
  playCount: number
  created: string
  starred: string
  year: number
  genre: string
  userRating: number
  recordLabels: RecordLabel[]
  musicBrainzId: string
  genres: Genre[]
  displayArtist: string
  sortName: string
  releaseDate: ReleaseDate
  song: SubsonicSong[]
  artists: Artist2[]
  albumArtists: AlbumArtist2[]
}

export interface RecordLabel {
  name: string
}

export interface Genre {
  name: string
}

export interface ReleaseDate {
  year: number
  month: number
  day: number
}

export interface Artist2 {
  id: string
  name: string
  userRating: number
  averageRating: number
}

export interface AlbumArtist2 {
  id: string
  name: string
  userRating: number
  averageRating: number
}
