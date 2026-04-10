import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import type { SubsonicArtist } from '~/types/subsonicArtist'
import type { SubsonicGenre } from '~/types/subsonicGenres'
import { useLocalStorage } from '@vueuse/core'

export enum AlbumOrders {
  RecentlyUpdated = 'Recently Updated',
  Random = 'Random',
  Alphabetical = 'Alphabetical',
  ReleaseDate = 'Release Date',
  RecentlyPlayed = 'Recently Played',
}

export type AlbumOrder = typeof AlbumOrders[keyof typeof AlbumOrders]

export enum ArtistOrders {
  RecentlyUpdated = 'Recently Updated',
  RecentlyPlayed = 'Recently Played',
  Random = 'Random',
  Alphabetical = 'Alphabetical',
  Starred = 'Starred',
}

export type ArtistOrder = typeof ArtistOrders[keyof typeof ArtistOrders]

export const streamQualities = [96, 128, 160, 192, 256, 'native'] as const
export type StreamQuality = typeof streamQualities[number]

export const apiKey = useLocalStorage('apiKey', '')
export const albumSeed = useLocalStorage<number>('albumSeed', 0)
export const albumOrder = useLocalStorage<AlbumOrder>('albumOrder', AlbumOrders.RecentlyUpdated)
export const artistSeed = useLocalStorage<number>('artistSeed', 0)
export const artistOrder = useLocalStorage<ArtistOrder>('artistOrder', ArtistOrders.RecentlyUpdated)
export const debugEnabled = useLocalStorage('debugEnabled', false)
export const shuffleEnabled = useLocalStorage<boolean>('shuffleEnabled', false)
export const repeatStatus = useLocalStorage<'off' | '1' | 'all'>('repeatStatus', 'off')
export const streamQuality = useLocalStorage<StreamQuality>('streamQuality', 160)
export const randomTracksSeed = useLocalStorage<number>('randomTracksSeed', 0)

export const albumsStore = useLocalStorage<SubsonicAlbum[]>('albumsStore', [])
export const artistsStore = useLocalStorage<SubsonicArtist[]>('artistsStore', [])
export const genresStore = useLocalStorage<SubsonicGenre[]>('genresStore', [])
