import type { TrackMetadata, TrackMetadataWithImageUrl } from '../types'
import { useSessionStorage } from '@vueuse/core'
import dayjs from 'dayjs'

export function useLogic() {
  const inputText = useSessionStorage<string>('searchInput', '')

  const closeSearch = () => {
    inputText.value = ''
  }

  const niceDate = (dateString: string): string => {
    const date = dayjs(dateString)
    return date.isValid() ? date.format('DD/MM/YYYY') : 'Invalid Date'
  }

  const formatTime = (time: number): string => {
    const minutes = Math.floor(time / 60)
    const seconds = Math.floor(time % 60)
    return `${minutes}:${seconds.toString().padStart(2, '0')}`
  }

  const trackWithImageUrl = (track: TrackMetadata | TrackMetadataWithImageUrl): TrackMetadataWithImageUrl => {
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

  const getArtistUrl = (musicbrainz_artist_id: string): string => {
    return `/artists/${musicbrainz_artist_id}`
  }

  const getAlbumUrl = (musicbrainz_album_id: string): string => {
    return `/albums/${musicbrainz_album_id}`
  }

  const getTrackUrl = (musicbrainz_track_id: string): string => {
    return `/tracks/${musicbrainz_track_id}`
  }

  const getRandomInteger = (): number => {
    return Math.floor(Math.random() * 1000000)
  }

  return {
    inputText,
    closeSearch,
    niceDate,
    formatTime,
    trackWithImageUrl,
    getArtistUrl,
    getAlbumUrl,
    getTrackUrl,
    getRandomInteger,
  }
}
