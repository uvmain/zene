import type { ReleaseDate } from '../types/subsonicAlbum'
import { useLocalStorage } from '@vueuse/core'
import dayjs from 'dayjs'
import { useSettings } from '~/composables/useSettings'

const apiKey = useLocalStorage('apiKey', '')
const { streamQuality } = useSettings()

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
  const seconds = Math.floor(time % 60)
  return `${minutes}:${seconds.toString().padStart(2, '0')}`
}

export function getCoverArtUrl(musicbrainzId: string, size = 400, timeUpdated?: string): string {
  if (timeUpdated != null) {
    return `/share/img/${musicbrainzId}?size=${size}&time=${timeUpdated}`
  }
  return `/share/img/${musicbrainzId}?size=${size}`
}

export function getAuthenticatedTrackUrl(musicbrainz_track_id: string): string {
  const queryParamString = `apiKey=${apiKey.value}&c=zene-frontend&v=1.6.0&maxBitRate=${streamQuality.value}&id=${musicbrainz_track_id}&format=aac`
  return `/rest/stream.view?${queryParamString}`
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

export async function cacheBustAlbumArt(albumId: string) {
  const promises = []
  promises.push(fetch(getCoverArtUrl(albumId), { method: 'POST' }))
  promises.push(fetch(getCoverArtUrl(albumId, 120), { method: 'POST' }))
  promises.push(fetch(getCoverArtUrl(albumId, 150), { method: 'POST' }))
  promises.push(fetch(getCoverArtUrl(albumId, 200), { method: 'POST' }))
  await Promise.all(promises)
}
