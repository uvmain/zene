import type { SubsonicAlbum } from './subsonicAlbum'
import type { SubsonicArtist, SubsonicArtists } from './subsonicArtist'
import type { SubsonicGenres } from './subsonicGenres'
import type { LyricsList } from './subsonicLyrics'
import type { SubsonicSong } from './subsonicSong'

export interface SubsonicResponse {
  'subsonic-response': {
    status: string
    version: string
    type: string
    serverVersion: string
    openSubsonic: boolean
    error?: {
      code: number
      message: string
      helpUrl?: string
    }
    [key: string]: any
  }
}
export interface SubsonicRandomSongsResponse extends SubsonicResponse {
  randomSongs: {
    song: SubsonicSong[]
  }
}

export interface SubsonicTopSongsResponse extends SubsonicResponse {
  topSongs: {
    song: SubsonicSong[]
  }
}

export interface SubsonicGenresResponse extends SubsonicResponse {
  genres: SubsonicGenres
}

export interface SubsonicAlbumListResponse extends SubsonicResponse {
  albumList: {
    album: SubsonicAlbum[]
  }
}

export interface SubsonicLyricsListResponse extends SubsonicResponse {
  lyricsList: LyricsList
}

export interface SubsonicAlbumResponse extends SubsonicResponse {
  album: SubsonicAlbum
}

export interface SubsonicArtistResponse extends SubsonicResponse {
  artist: SubsonicArtist
}

export interface SubsonicArtistsResponse extends SubsonicResponse {
  artists: SubsonicArtists
}
