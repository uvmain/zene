export interface SubsonicSong {
  id: string
  parent: string
  isDir: boolean
  title: string
  album: string
  artist: string
  track: number
  year: number
  genre: string
  coverArt: string
  size: number
  duration: number
  bitRate: number
  samplingRate: number
  channelCount: number
  path: string
  discNumber: number
  created: string
  albumId: string
  artistId: string
  musicBrainzId: string
  genres: Genre2[]
  artists: Artist[]
  displayArtist: string
  albumArtists: AlbumArtist[]
  displayAlbumArtist: string
  playcount?: number
}

export interface Genre2 {
  name: string
}

export interface Artist {
  id: string
  name: string
}

export interface AlbumArtist {
  id: string
  name: string
}
