import type { ReleaseDate } from '~/types/subsonicAlbum'
import { apiKey, streamQuality } from '~/logic/store'
import { chromecastConnected } from './castRefs'

const chromecastDevMode = computed(() => {
  return chromecastConnected.value === true
    && import.meta.env.MODE === 'development'
    && import.meta.env.DEV === true
    && import.meta.env.VITE_CHROMECAST_BASE_URL !== undefined
    && import.meta.env.VITE_CHROMECAST_API_KEY !== undefined
})

export function formatTimeFromSeconds(time: number): string {
  const minutes = Math.floor(time / 60)
  const seconds = Math.floor(time % 60).toString().padStart(2, '0')
  return `${minutes}:${seconds}`
}

export async function clearApiKey() {
  apiKey.value = ''
}

export function getAuthenticatedTrackUrl(musicbrainz_track_id: string, raw = false): string {
  const queryParams = new URLSearchParams({
    apiKey: chromecastDevMode.value ? import.meta.env.VITE_CHROMECAST_API_KEY as string : apiKey.value,
    c: 'zene-frontend',
    v: '1.6.0',
    id: musicbrainz_track_id,
    format: chromecastDevMode.value ? 'mp3' : 'aac',
  })
  if (!raw) {
    queryParams.append('maxBitRate', streamQuality.value.toString())
  }
  else {
    queryParams.append('raw', 'true')
  }
  return chromecastDevMode.value ? `${import.meta.env.VITE_CHROMECAST_BASE_URL}/rest/stream.view?${queryParams.toString()}` : `/rest/stream.view?${queryParams.toString()}`
}

export function onImageError(event: Event) {
  const target = event.target as HTMLImageElement
  target.onerror = null
  target.src = '/default-square.png'
}

export function parseReleaseDate(releaseDate: ReleaseDate): string {
  return `${(releaseDate.year ?? 1).toString().padStart(4, '0')}`
}

export function generateSeed() {
  return Math.floor(Math.random() * 1000000)
}

export enum artSizes {
  size60 = 60,
  size120 = 120,
  size150 = 150,
  size200 = 200,
  size400 = 400,
}

export function getCoverArtUrl(musicbrainzId: string, size: number = artSizes.size400, timeUpdated?: string): string {
  if (timeUpdated != null) {
    return size === 0 ? `/share/img/${musicbrainzId}?time=${timeUpdated}` : `/share/img/${musicbrainzId}?size=${size}&time=${timeUpdated}`
  }
  return size === 0 ? `/share/img/${musicbrainzId}` : `/share/img/${musicbrainzId}?size=${size}`
}

export async function cacheBustAlbumArt(albumId: string) {
  const promises = []
  promises.push(fetch(getCoverArtUrl(albumId, 0), { method: 'POST' }))
  for (const size of Object.values(artSizes).filter(value => typeof value === 'number')) {
    promises.push(fetch(getCoverArtUrl(albumId, size), { method: 'POST' }))
  }
  await Promise.all(promises)
}
