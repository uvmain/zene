import type { SubsonicGenre } from './subsonicGenres'

export interface SubsonicPodcastChannel {
  id: string
  parent: string
  isDir: string
  title: string
  url: string
  description: string
  coverArt: string
  originalImageUrl: string
  status: string
  type: string
  isVideo: string
  streamId: string
  channelId: string
  lastRefresh: string
  created: string
  episode: SubsonicPodcastEpisode[]
}

export interface SubsonicPodcastEpisode {
  id: string
  streamId: string
  channelId: string
  title: string
  description: string
  publishDate: string
  status: string
  parent: string
  isDir: string
  year: string
  genre: string
  genres: SubsonicGenre[]
  coverArt: string
  size: string
  contentType: string
  suffix: string
  duration: string
  bitRate: string
  path: string
  sourceUrl: string
}
