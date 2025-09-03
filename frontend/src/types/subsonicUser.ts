export interface SubsonicUser {
  folder: number[]
  username: string
  password?: string
  email: string
  adminRole: boolean
  scrobblingEnabled: boolean
  streamRole: boolean
  settingsRole: boolean
  jukeboxRole: boolean
  downloadRole: boolean
  uploadRole: boolean
  playlistRole: boolean
  coverArtRole: boolean
  commentRole: boolean
  podcastRole: boolean
  shareRole: boolean
  videoConversionRole: boolean
  maxBitRate: number
}
