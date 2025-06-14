import type { AlbumMetadata, TrackMetadata, TrackMetadataWithImageUrl } from '../types'
import type { User } from '../types/auth'
import { trackWithImageUrl } from '../composables/logic'

export async function backendFetchRequest(path: string, options = {}): Promise<Response> {
  const url = `/api/${path}`
  const response = await fetch(url, options)
  return response
}

export async function getAlbumTracks(musicbrainz_album_id: string): Promise<TrackMetadataWithImageUrl[]> {
  const response = await backendFetchRequest(`albums/${musicbrainz_album_id}/tracks`)
  const json = await response.json() as TrackMetadata[]
  const trackArray: TrackMetadataWithImageUrl[] = []
  json.forEach((track) => {
    trackArray.push(trackWithImageUrl(track))
  })
  return trackArray
}

export async function getArtistTracks(musicbrainz_artist_id: string, limit = 0): Promise<TrackMetadataWithImageUrl[]> {
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

export async function getArtistAlbums(musicbrainz_artist_id: string): Promise<AlbumMetadata[]> {
  const response = await backendFetchRequest(`artists/${musicbrainz_artist_id}/albums?chronological=true`)
  const json = await response.json() as AlbumMetadata[]
  return json
}

export async function getCurrentUser(): Promise<User> {
  const response = await backendFetchRequest('user')
  const json = await response.json() as User
  return json
}

export async function getUsers(): Promise<User[]> {
  const response = await backendFetchRequest('users')
  const json = await response.json() as User[]
  return json
}
