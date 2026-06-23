import type { ReleaseDate } from '~/types/subsonicAlbum'
import { apiKey, overrideBackendUrl, streamQuality } from '~/stores/main'

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
    apiKey: apiKey.value,
    c: 'zene-frontend',
    v: '1.6.0',
    id: musicbrainz_track_id,
    format: 'aac',
  })
  if (!raw) {
    queryParams.append('maxBitRate', streamQuality.value.toString())
  }
  else {
    queryParams.append('raw', 'true')
  }
  const path = `/rest/stream.view?${queryParams.toString()}`
  const url = overrideBackendUrl.value === '' ? path : `${overrideBackendUrl.value}${path}`
  return url
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
  let path = ''
  if (timeUpdated != null) {
    path = size === 0 ? `/share/img/${musicbrainzId}?time=${timeUpdated}` : `/share/img/${musicbrainzId}?size=${size}&time=${timeUpdated}`
  }
  else {
    path = size === 0 ? `/share/img/${musicbrainzId}` : `/share/img/${musicbrainzId}?size=${size}`
  }
  const url = overrideBackendUrl.value === '' ? path : `${overrideBackendUrl.value}${path}`
  return url
}

export async function cacheBustArt(musicbrainz_id: string) {
  const promises = []
  promises.push(fetch(getCoverArtUrl(musicbrainz_id, 0), { method: 'POST' }))
  for (const size of Object.values(artSizes).filter(value => typeof value === 'number')) {
    promises.push(fetch(getCoverArtUrl(musicbrainz_id, size), { method: 'POST' }))
  }
  await Promise.all(promises)
}
