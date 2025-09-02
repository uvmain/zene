export interface SubsonicArtist {
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

export interface SubsonicArtists {
  index: Index[]
}

export interface Index {
  name: string
  artist: SubsonicArtist[]
}
