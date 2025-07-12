import type { AlbumMetadata, ArtistMetadata, GenreMetadata, TrackMetadata, TrackMetadataWithImageUrl } from '../types'
import { useSessionStorage } from '@vueuse/core'
import dayjs from 'dayjs'
import { useBackendFetch } from './useBackendFetch'

const searchInput = useSessionStorage<string>('searchInput', '')
const { backendFetchRequest } = useBackendFetch()

const searchResults = ref<TrackMetadataWithImageUrl[]>([])
const searchResultsGenres = ref<any[]>([])
const searchResultsArtists = ref<ArtistMetadata[]>([])

export function useSearch() {
  const closeSearch = () => {
    searchInput.value = ''
  }

  const getGenres = async () => {
    const response = await backendFetchRequest(`genres?search=${searchInput.value}`)
    const json = await response.json() as GenreMetadata[]
    if (json.length === 0) {
      searchResultsGenres.value = []
      return
    }
    searchResultsGenres.value = json.slice(0, 12)
  }

  const getArtists = async () => {
    const response = await backendFetchRequest(`artists?search=${searchInput.value}`)
    const json = await response.json() as ArtistMetadata[]
    if (json.length === 0) {
      searchResultsArtists.value = []
      return
    }
    searchResultsArtists.value = json
  }

  const search = async () => {
    if (!searchInput.value || searchInput.value.length < 3) {
      searchResults.value = []
      return
    }
    const response = await backendFetchRequest(`search?search=${searchInput.value}`)
    const json = await response.json() as TrackMetadata[]
    const albumMetadata: TrackMetadataWithImageUrl[] = []
    json.forEach((metadata) => {
      const metadataInstance = {
        file_path: metadata.file_path,
        file_name: metadata.file_name,
        date_added: metadata.date_added,
        date_modified: metadata.date_modified,
        format: metadata.format,
        duration: metadata.duration,
        size: metadata.size,
        bitrate: metadata.bitrate,
        title: metadata.title,
        artist: metadata.artist,
        album: metadata.album,
        album_artist: metadata.album_artist ?? metadata.artist,
        track_number: metadata.track_number,
        total_tracks: metadata.total_tracks,
        disc_number: metadata.disc_number,
        total_discs: metadata.total_discs,
        musicbrainz_artist_id: metadata.musicbrainz_artist_id,
        musicbrainz_album_id: metadata.musicbrainz_album_id,
        musicbrainz_track_id: metadata.musicbrainz_track_id,
        label: metadata.label,
        genre: metadata.genre,
        release_date: dayjs(metadata.release_date).format('YYYY'),
        image_url: `/api/albums/${metadata.musicbrainz_album_id}/art?size=lg`,
        user_play_count: metadata.user_play_count,
        global_play_count: metadata.global_play_count,
      }
      albumMetadata.push(metadataInstance)
    })
    searchResults.value = albumMetadata
    await getGenres()
    await getArtists()
  }

  const searchResultsAlbums = computed(() => {
    const uniqueAlbums = new Map<string, AlbumMetadata>()
    searchResults.value.forEach((album: TrackMetadataWithImageUrl) => {
      if (!uniqueAlbums.has(album.musicbrainz_album_id)) {
        uniqueAlbums.set(album.musicbrainz_album_id, {
          artist: album.artist,
          album_artist: album.album_artist ?? album.artist,
          album: album.album,
          musicbrainz_album_id: album.musicbrainz_album_id,
          musicbrainz_artist_id: album.musicbrainz_artist_id,
          genres: album.genre.split(';').filter((genre: string) => genre !== ''),
          release_date: album.release_date,
          image_url: album.image_url,
        })
      }
    })
    return Array.from(uniqueAlbums.values())
  })

  const searchResultsTracks = computed(() => {
    return searchResults.value
  })

  return {
    search,
    searchInput,
    searchResults,
    searchResultsGenres,
    searchResultsArtists,
    searchResultsTracks,
    searchResultsAlbums,
    closeSearch,
  }
}
