import type { AlbumMetadata, TrackMetadata, TrackMetadataWithImageUrl } from '../types'
import type { TokenResponse, User, UsersResponse } from '../types/auth'
import { trackWithImageUrl } from '../composables/logic'
import { usePlaybackQueue } from '../composables/usePlaybackQueue'

const { getRandomSeed } = usePlaybackQueue()

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

  const getMimeType = async (url: string): Promise<string> => {
    const response = await fetch(url, { method: 'HEAD' })
    const contentType = response.headers.get('content-type') ?? response.headers.get('Content-Type') ?? ''
    return contentType
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
    getMimeType,
  }
}
