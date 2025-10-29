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

interface Index {
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

export interface SubsonicArtistInfo {
  musicBrainzId: string
  smallImageUrl: string
  mediumImageUrl: string
  largeImageUrl: string
  similarArtists: {
    id: string
    name: string
    coverArt: string
    artistImageUrl: string
    albumCount: number
    sortName: string
    userRating: number
    averageRating: number
  }[]
}

interface Album {
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

interface RecordLabel {
  name: string
}

interface Genre {
  name: string
}

interface Artist2 {
  id: string
  name: string
}

interface AlbumArtist {
  id: string
  name: string
}
