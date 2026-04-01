import type { SubsonicAlbum } from './subsonicAlbum'
import type { SubsonicArtist } from './subsonicArtist'
import type { SubsonicGenre } from './subsonicGenres'
import type { SubsonicPodcastEpisode } from './subsonicPodcasts'
import type { SubsonicSong } from './subsonicSong'

export type LoadingAttribute = 'lazy' | 'eager'

export interface SearchResult { artists: SubsonicArtist[], albums: SubsonicAlbum[], songs: SubsonicSong[], genres: SubsonicGenre[] }

export interface ButterchurnPreset {
  name: string
  preset: any
}

export interface PlayItem {
  track?: SubsonicSong
  podcastEpisode?: SubsonicPodcastEpisode
}
