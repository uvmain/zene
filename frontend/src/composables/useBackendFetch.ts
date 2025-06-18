import type { AlbumMetadata, TrackMetadata, TrackMetadataWithImageUrl } from '../types'
import type { User, UsersResponse } from '../types/auth'
import { trackWithImageUrl } from '../composables/logic'

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
    let url = `artists/${musicbrainz_artist_id}/tracks?random=true`
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

  return {
    backendFetchRequest,
    getAlbumTracks,
    getArtistTracks,
    getArtistAlbums,
    getCurrentUser,
    getGenreTracks,
    getUsers,
  }
}
