import type { TrackMetadata, TrackMetadataWithImageUrl } from '../types'
import dayjs from 'dayjs'

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
