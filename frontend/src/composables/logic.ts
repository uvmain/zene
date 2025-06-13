import type { TrackMetadata, TrackMetadataWithImageUrl } from '../types'
import { useSessionStorage } from '@vueuse/core'
import dayjs from 'dayjs'

const inputText = useSessionStorage<string>('searchInput', '')

export function closeSearch() {
  inputText.value = ''
}

export function niceDate(dateString: string): string {
  const date = dayjs(dateString)
  return date.isValid() ? date.format('DD/MM/YYYY') : 'Invalid Date'
}

export function formatTime(time: number): string {
  const minutes = Math.floor(time / 60)
  const seconds = Math.floor(time % 60)
  return `${minutes}:${seconds.toString().padStart(2, '0')}`
}

export function trackWithImageUrl(track: TrackMetadata | TrackMetadataWithImageUrl): TrackMetadataWithImageUrl {
  let trackMetadataWithImageUrl: TrackMetadataWithImageUrl
  if (Object.hasOwn(track, 'image_url')) {
    trackMetadataWithImageUrl = track as TrackMetadataWithImageUrl
  }
  else {
    trackMetadataWithImageUrl = track as TrackMetadataWithImageUrl
    trackMetadataWithImageUrl.image_url = `/api/albums/${track.musicbrainz_album_id}/art`
  }
  return trackMetadataWithImageUrl
}

export function getArtistUrl(musicbrainz_artist_id: string) {
  return `/artists/${musicbrainz_artist_id}`
}

export function getAlbumUrl(musicbrainz_album_id: string) {
  return `/albums/${musicbrainz_album_id}`
}

export function getTrackUrl(musicbrainz_track_id: string) {
  return `/tracks/${musicbrainz_track_id}`
}
