import type { SubsonicAlbum } from './subsonicAlbum'
import type { SubsonicApiKey } from './subsonicApiKey'
import type { SubsonicArtist, SubsonicIndexArtists } from './subsonicArtist'
import type { SubsonicGenres } from './subsonicGenres'
import type { LyricsList } from './subsonicLyrics'
import type { SubsonicPlaylist } from './subsonicPlaylists'
import type { SubsonicPodcastChannel } from './subsonicPodcasts'
import type { SubsonicRadioStation } from './subsonicRadioStations'
import type { SubsonicSong } from './subsonicSong'
import type { SubsonicUser } from './subsonicUser'

export interface SubsonicResponseWrapper {
  'subsonic-response': SubsonicResponse
}

export interface SubsonicResponse {
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

export interface SubsonicSongResponse extends SubsonicResponse {
  song: SubsonicSong
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
  artists: SubsonicIndexArtists
}

export interface SubsonicUserResponse extends SubsonicResponse {
  user: SubsonicUser
}

export interface SubsonicUsersResponse extends SubsonicResponse {
  users: {
    user: SubsonicUser[]
  }
}

export interface SubsonicSongsByGenreResponse extends SubsonicResponse {
  songsByGenre: {
    song: SubsonicSong[]
  }
}

export interface SubsonicSearchResponse extends SubsonicResponse {
  searchResult2: {
    artist: SubsonicArtist[]
    album: SubsonicAlbum[]
    song: SubsonicSong[]
  }
}

export interface SubsonicApiKeyResponse extends SubsonicResponse {
  apiKeys: {
    apiKey: SubsonicApiKey[]
  }
}

export interface SubsonicRadioStationsResponse extends SubsonicResponse {
  internetRadioStations: {
    internetRadioStation: SubsonicRadioStation[]
  }
}

export interface SubsonicPodcastChannelsResponse extends SubsonicResponse {
  podcasts: {
    channel: SubsonicPodcastChannel[]
  }
}

export interface SubsonicPlaylistsResponse extends SubsonicResponse {
  playlists: {
    playlist: SubsonicPlaylist[]
  }
}

export interface SubsonicPlaylistResponse extends SubsonicResponse {
  playlist: SubsonicPlaylist
}
