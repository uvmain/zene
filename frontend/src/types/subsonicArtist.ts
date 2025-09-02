export interface SubsonicIndexArtist {
  id: string
  name: string
  coverArt: string
  artistImageUrl: string
  albumCount: number
  musicBrainzId: string
  sortName: string
  userRating: number
  averageRating: number
}

export interface SubsonicIndexArtists {
  index: Index[]
}

export interface Index {
  name: string
  artist: SubsonicIndexArtist[]
}

export interface SubsonicArtist {
  id: string
  name: string
  coverArt: string
  artistImageUrl: string
  albumCount: number
  album: Album[]
  musicBrainzId: string
  sortName: string
  userRating: number
  averageRating: number
}

export interface Album {
  id: string
  parent: string
  isDir: boolean
  title: string
  album: string
  artist: string
  year: number
  genre: string
  coverArt: string
  duration: number
  recordLabels: RecordLabel[]
  songCount: number
  created: string
  artistId: string
  genres: Genre[]
  artists: Artist2[]
  displayArtist: string
  albumArtists: AlbumArtist[]
  displayAlbumArtist: string
}

export interface RecordLabel {
  name: string
}

export interface Genre {
  name: string
}

export interface Artist2 {
  id: string
  name: string
}

export interface AlbumArtist {
  id: string
  name: string
}
