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
  displayAlbumArtist: string
  sortName: string
  releaseDate: ReleaseDate
  song: SubsonicSong[]
  artists: Artist2[]
  albumArtists: AlbumArtist2[]
}

interface RecordLabel {
  name: string
}

interface Genre {
  name: string
}

export interface ReleaseDate {
  year: number
  month: number
  day: number
}

interface Artist2 {
  id: string
  name: string
  userRating: number
  averageRating: number
}

interface AlbumArtist2 {
  id: string
  name: string
  userRating: number
  averageRating: number
}
