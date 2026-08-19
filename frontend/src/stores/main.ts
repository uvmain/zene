import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import type { SubsonicArtist } from '~/types/subsonicArtist'
import { useLocalStorage, useSessionStorage } from '@vueuse/core'

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

export const streamQualities: readonly (string | number)[] = ['Unlimited', 96, 128, 160, 192, 256, 320, 512, 1024]

export const apiKey = useLocalStorage<string>('apiKey', '')
export const shareToken = useSessionStorage<string>('shareToken', '')
export const albumSeed = useLocalStorage<number>('albumSeed', 0)
export const albumOrder = useLocalStorage<AlbumOrder>('albumOrder', AlbumOrders.RecentlyUpdated)
export const artistSeed = useLocalStorage<number>('artistSeed', 0)
export const artistOrder = useLocalStorage<ArtistOrder>('artistOrder', ArtistOrders.RecentlyUpdated)
export const debugEnabled = useLocalStorage('debugEnabled', false)
export const shuffleEnabled = useLocalStorage<boolean>('shuffleEnabled', false)
export const repeatStatus = useLocalStorage<'off' | '1' | 'all'>('repeatStatus', 'off')
export const transcodeStreamQuality = useLocalStorage<string | number>('transcodeStreamQuality', 256)
export const directPlayMaxQuality = useLocalStorage<string | number>('directPlayMaxQuality', 'Unlimited')
export const randomTracksSeed = useLocalStorage<number>('randomTracksSeed', 0)
export const albumsStore = useLocalStorage<SubsonicAlbum[]>('albumsStore', [])
export const artistsStore = useLocalStorage<SubsonicArtist[]>('artistsStore', [])
export const volumeStore = useLocalStorage<string>('volumeStore', '1')
export const accentColour = useLocalStorage<string>('accentColour', 'hsla(22 95% 60% / 1)')
export const autoSwitchColours = useLocalStorage<boolean>('autoSwitchColours', true)
export const backendUrl = useLocalStorage<string>('backendUrl', new URL(window.location.href).origin)
export const wakeLockEnabled = useLocalStorage<boolean>('wakeLockEnabled', true)
