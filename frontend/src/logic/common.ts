import type { ReleaseDate } from '~/types/subsonicAlbum'
import dayjs from 'dayjs'
import { apiKey, streamQuality } from '~/logic/store'

export function niceDate(dateString: string): string {
  const date = dayjs(dateString)
  return date.isValid() ? date.format('DD/MM/YYYY') : 'Invalid Date'
}

export function formatDate(dateString: string): string {
  if (!dateString)
    return 'N/A'
  const date = new Date(dateString)
  return date.toLocaleDateString(undefined, { year: 'numeric', month: 'long', day: 'numeric' })
}

export function formatTimeFromSeconds(time: number): string {
  const minutes = Math.floor(time / 60)
  const seconds = Math.floor(time % 60).toString().padStart(2, '0')
  return `${minutes}:${seconds}`
}

export function getAuthenticatedTrackUrl(musicbrainz_track_id: string, raw = false): string {
  const queryParams = new URLSearchParams({
    apiKey: apiKey.value,
    c: 'zene-frontend',
    v: '1.6.0',
    id: musicbrainz_track_id,
    format: raw ? 'raw' : 'aac',
  })
  if (!raw) {
    queryParams.append('maxBitRate', streamQuality.value.toString())
  }
  else {
    queryParams.append('raw', 'true')
  }
  // const queryParamString = `v=1.6.0&maxBitRate=${streamQuality.value}&id=${musicbrainz_track_id}&format=aac`
  return `/rest/stream.view?${queryParams.toString()}`
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
    return `/share/img/${musicbrainzId}?size=${size}&time=${timeUpdated}`
  }
  return `/share/img/${musicbrainzId}?size=${size}`
}

export async function cacheBustAlbumArt(albumId: string) {
  const promises = []
  promises.push(fetch(getCoverArtUrl(albumId), { method: 'POST' }))
  for (const size of Object.values(artSizes).filter(value => typeof value === 'number')) {
    promises.push(fetch(getCoverArtUrl(albumId, size), { method: 'POST' }))
  }
  await Promise.all(promises)
}
