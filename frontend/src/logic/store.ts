import { useLocalStorage } from '@vueuse/core'

export const albumOrders = ['recentlyUpdated', 'random', 'alphabetical', 'releaseDate', 'recentlyPlayed'] as const
export type AlbumOrder = typeof albumOrders[number]

export const artistOrders = ['newest', 'random', 'alphabetical', 'starred', 'recent'] as const
export type ArtistOrder = typeof artistOrders[number]

export const streamQualities = [96, 128, 160, 192, 256, 'native'] as const
export type StreamQuality = typeof streamQualities[number]

export const apiKey = useLocalStorage('apiKey', '')
export const albumSeed = useLocalStorage<number>('albumSeed', 0)
export const albumOrder = useLocalStorage<AlbumOrder>('albumOrder', 'recentlyUpdated')
export const artistSeed = useLocalStorage<number>('artistSeed', 0)
export const artistOrder = useLocalStorage<ArtistOrder>('artistOrder', 'newest')
export const debugEnabled = useLocalStorage('debugEnabled', false)
export const shuffleEnabled = useLocalStorage<boolean>('shuffleEnabled', false)
export const repeatStatus = useLocalStorage<'off' | '1' | 'all'>('repeatStatus', 'off')
export const streamQuality = useLocalStorage<StreamQuality>('streamQuality', 160)
export const randomTracksSeed = useLocalStorage<number>('randomTracksSeed', 0)
