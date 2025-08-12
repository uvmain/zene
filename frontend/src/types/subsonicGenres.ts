export interface SubsonicGenresResponse {
  'subsonic-response': {
    status: string
    version: string
    type: string
    serverVersion: string
    openSubsonic: boolean
    genres: Genres
  }
}

export interface Genres {
  genre: Genre[]
}

export interface Genre {
  song_count: string
  album_count: string
  value: string
}
