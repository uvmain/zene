export interface SubsonicUserResponse {
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
    user: SubsonicUser
  }
}

export interface SubsonicUsersResponse {
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
    users: SubsonicUsers
  }
}

export interface SubsonicUsers {
  user: SubsonicUser[]
}

export interface SubsonicUser {
  folder: number[]
  username: string
  email: string
  password: string
  scrobblingEnabled: string
  adminRole: string
  settingsRole: string
  downloadRole: string
  uploadRole: string
  playlistRole: string
  coverArtRole: string
  commentRole: string
  podcastRole: string
  streamRole: string
  jukeboxRole: string
  shareRole: string
}
