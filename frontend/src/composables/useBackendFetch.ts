import type { AlbumMetadata, TrackMetadata, TrackMetadataWithImageUrl } from '../types'
import type { TokenResponse, User, UsersResponse } from '../types/auth'
import { useRandomSeed } from '../composables/useRandomSeed'
import { useLogic } from './useLogic'

const { getRandomSeed } = useRandomSeed()
const { trackWithImageUrl } = useLogic()

export function useBackendFetch() {
  const backendFetchRequest = async (path: string, options = {}): Promise<Response> => {
    const url = `/api/${path}`
    const response = await fetch(url, options)
    return response
  }

  const getAlbumTracks = async (musicbrainz_album_id: string): Promise<TrackMetadataWithImageUrl[]> => {
    const response = await backendFetchRequest(`albums/${musicbrainz_album_id}/tracks`)
    const json = await response.json() as TrackMetadata[]
    const trackArray: TrackMetadataWithImageUrl[] = []
    json.forEach((track) => {
      trackArray.push(trackWithImageUrl(track))
    })
    return trackArray
  }

  const getArtistTracks = async (musicbrainz_artist_id: string, limit = 0): Promise<TrackMetadataWithImageUrl[]> => {
    const randomSeed = getRandomSeed()
    let url = `artists/${musicbrainz_artist_id}/tracks?random=${randomSeed}`
    if (limit > 0) {
      url = `${url}&limit=${limit}`
    }
    const response = await backendFetchRequest(url)
    const json = await response.json() as TrackMetadata[]
    const trackArray: TrackMetadataWithImageUrl[] = []
    json.forEach((track) => {
      trackArray.push(trackWithImageUrl(track))
    })
    return trackArray
  }

  const getArtistAlbums = async (musicbrainz_artist_id: string): Promise<AlbumMetadata[]> => {
    const response = await backendFetchRequest(`artists/${musicbrainz_artist_id}/albums?chronological=true`)
    const json = await response.json() as AlbumMetadata[]
    return json
  }

  const getCurrentUser = async (): Promise<User> => {
    const response = await backendFetchRequest('user')
    const json = await response.json() as User
    return json
  }

  const getUsers = async (): Promise<User[]> => {
    const response = await backendFetchRequest('users')
    const json = await response.json() as UsersResponse
    return json.users
  }

  const getGenreTracks = async (genre: string, limit = 0, random = false): Promise<TrackMetadataWithImageUrl[]> => {
    const response = await backendFetchRequest(`genres/tracks?genres=${genre}&limit=${limit}&random=${random}`)
    return await response.json() as TrackMetadataWithImageUrl[]
  }

  const getTemporaryToken = async (duration = 30): Promise<TokenResponse> => {
    const response = await backendFetchRequest(`temporary_token?duration=${duration}`)
    return await response.json() as TokenResponse
  }

  const refreshTemporaryToken = async (currentToken: string, duration = 30): Promise<TokenResponse> => {
    const formData = new FormData()
    formData.append('token', currentToken)
    formData.append('duration', duration.toString())

    const response = await backendFetchRequest('temporary_token', {
      method: 'POST',
      body: formData,
    })
    return await response.json() as TokenResponse
  }

  const getMimeType = async (url: string): Promise<string> => {
    const response = await fetch(url, { method: 'HEAD' })
    const contentType = response.headers.get('content-type') ?? response.headers.get('Content-Type') ?? ''
    return contentType
  }

  const getLyrics = async (musicbrainzTrackId: string): Promise<{ plainLyrics: string, syncedLyrics: string }> => {
    const response = await backendFetchRequest(`tracks/${musicbrainzTrackId}/lyrics`)
    if (!response.ok) {
      throw new Error(`Failed to fetch lyrics for track ${musicbrainzTrackId}: ${response.statusText}`)
    }
    const data = await response.json() as { plainLyrics: string, syncedLyrics: string }
    return data
  }

  return {
    backendFetchRequest,
    getAlbumTracks,
    getArtistTracks,
    getArtistAlbums,
    getCurrentUser,
    getGenreTracks,
    getUsers,
    getTemporaryToken,
    refreshTemporaryToken,
    getMimeType,
    getLyrics,
  }
}
