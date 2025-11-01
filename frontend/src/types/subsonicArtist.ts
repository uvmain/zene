import type { SubsonicAlbum } from '~/types/subsonicAlbum'

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
  album: SubsonicAlbum[]
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
