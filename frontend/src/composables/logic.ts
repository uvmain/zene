import type { ReleaseDate } from '../types/subsonicAlbum'
import dayjs from 'dayjs'

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

export function getCoverArtUrl(musicbrainzId: string, size = 400): string {
  return `/share/img/${musicbrainzId}?size=${size}`
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
