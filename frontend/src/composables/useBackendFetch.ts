import type { AlbumMetadata, StandardResponse, TrackMetadata, TrackMetadataWithImageUrl } from '../types'
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

  const postPlaycount = async (musicbrainz_track_id: string): Promise<StandardResponse> => {
    const formData = new FormData()
    formData.append('musicbrainz_track_id', musicbrainz_track_id)

    const response = await backendFetchRequest('playcounts', {
      body: formData,
      method: 'POST',
    })

    const json = await response.json() as StandardResponse
    return json
  }

  return {
    backendFetchRequest,
    getAlbumTracks,
    getArtistTracks,
    getArtistAlbums,
    getCurrentUser,
    getUsers,
    postPlaycount,
  }
}
