import { SubsonicResponse } from "./subsonic"

export type Share = {
  id: number
  url: string
  description: string
  username: string
  created: string
  visitCount: number
  expires?: string
  lastVisited?: string

  entry: {
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
    contentType: string
    suffix: string
    duration: number
    bitRate: number
    samplingRate: number
    channelCount: number
    path: string
    discNumber: number
    created: string
    albumId: string
    artistId: string
    type: string
    mediaType: string
    sortName: string
    musicBrainzId: string
    genres: {
      name: string
    }[]
    artists: {
      id: string
      name: string
    }[]
    displayArtist: string
    albumArtists: {
      id: string
      name: string
    }[]
    displayAlbumArtist: string
  }[]
}

export interface SubsonicSharesResponse extends SubsonicResponse {
  shares: {
    share: Share[]
  }
}