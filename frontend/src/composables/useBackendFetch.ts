import type { AlbumMetadata, TrackMetadata, TrackMetadataWithImageUrl } from '~/types'
import type { SubsonicLyricsResponse } from '~/types/subsonicLyrics'
import type { SubsonicUser, SubsonicUserResponse, SubsonicUsersResponse } from '~/types/subsonicUser'
import { useAuth } from '~/composables/useAuth'
import { useRandomSeed } from '~/composables/useRandomSeed'
import { useLogic } from './useLogic'

const { userApiKey, userSalt, userToken, userUsername } = useAuth()
const { getRandomSeed } = useRandomSeed()
const { trackWithImageUrl } = useLogic()

export function useBackendFetch() {
  const backendFetchRequest = async (path: string, options: RequestInit = {}): Promise<Response> => {
    const url = `/api/${path}`
    const formData = new FormData()
    if (userApiKey.value) {
      formData.append('apiKey', userApiKey.value)
    }
    else if (userSalt.value && userToken.value) {
      formData.append('u', userUsername.value)
      formData.append('s', userSalt.value)
      formData.append('t', userToken.value)
    }
    else {
      const router = useRouter()
      await router.push('/login')
    }
    formData.append('f', 'json')
    formData.append('v', '1.16.0')
    formData.append('c', 'zene-frontend')

    // append formdata to existing body
    if (options.body instanceof FormData) {
      formData.forEach((value, key) => {
        (options.body as FormData).append(key, value)
      })
    }
    else {
      options.body = formData
    }

    options.method = 'POST'

    const response = await fetch(url, options)
    return response
  }

  const openSubsonicFetchRequest = async (path: string, options: RequestInit = {}): Promise<Response> => {
    const formData = new FormData()
    if (userApiKey.value) {
      formData.append('apiKey', userApiKey.value)
    }
    else if (userSalt.value && userToken.value) {
      formData.append('u', userUsername.value)
      formData.append('s', userSalt.value)
      formData.append('t', userToken.value)
    }
    else {
      const router = useRouter()
      await router.push('/login')
    }
    formData.append('f', 'json')
    formData.append('v', '1.16.0')
    formData.append('c', 'zene-frontend')

    // append formdata to existing body
    if (options.body instanceof FormData) {
      formData.forEach((value, key) => {
        (options.body as FormData).append(key, value)
      })
    }
    else {
      options.body = formData
    }

    options.method = options.method ?? 'POST'

    const url = `/rest/${path}`

    const response = await fetch(url, options)
    return response
  }

  const getAlbum = async (musicbrainz_album_id: string): Promise<AlbumMetadata> => {
    const response = await backendFetchRequest(`albums/${musicbrainz_album_id}`)
    if (!response.ok) {
      throw new Error(`Failed to fetch album with ID ${musicbrainz_album_id}: ${response.statusText}`)
    }
    const json = await response.json() as AlbumMetadata
    return json
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
    const formData = new FormData()
    formData.append('limit', limit.toString())
    formData.append('random', randomSeed.toString())
    const url = `artists/${musicbrainz_artist_id}/tracks`
    const response = await backendFetchRequest(url, {
      method: 'POST',
      body: formData,
    })
    const json = await response.json() as TrackMetadata[]
    const trackArray: TrackMetadataWithImageUrl[] = []
    json.forEach((track) => {
      trackArray.push(trackWithImageUrl(track))
    })
    return trackArray
  }

  const getArtistAlbums = async (musicbrainz_artist_id: string): Promise<AlbumMetadata[]> => {
    const formData = new FormData()
    formData.append('chronological', 'true')
    const response = await backendFetchRequest(`artists/${musicbrainz_artist_id}/albums`, {
      method: 'POST',
      body: formData,
    })
    const json = await response.json() as AlbumMetadata[]
    return json
  }

  const getCurrentUser = async (): Promise<SubsonicUser> => {
    const response = await openSubsonicFetchRequest('getUser.view')
    const json = await response.json() as SubsonicUserResponse
    return json['subsonic-response'].user
  }

  const getUsers = async (): Promise<SubsonicUser[]> => {
    const response = await openSubsonicFetchRequest('getUsers.view')
    const json = await response.json() as SubsonicUsersResponse
    return json['subsonic-response'].users.user
  }

  const getGenreTracks = async (genre: string, limit = 0, random = false): Promise<TrackMetadataWithImageUrl[]> => {
    const formData = new FormData()
    formData.append('genres', genre)
    formData.append('limit', limit.toString())
    formData.append('random', random.toString())
    const response = await backendFetchRequest('genres/tracks', {
      method: 'POST',
      body: formData,
    })
    return await response.json() as TrackMetadataWithImageUrl[]
  }

  const getMimeType = async (url: string): Promise<string> => {
    const response = await fetch(url, { method: 'HEAD' })
    const contentType = response.headers.get('content-type') ?? response.headers.get('Content-Type') ?? ''
    return contentType
  }

  const getLyrics = async (musicbrainzTrackId: string): Promise<SubsonicLyricsResponse> => {
    const formData = new FormData()
    formData.append('id', musicbrainzTrackId)
    const response = await openSubsonicFetchRequest('getLyricsBySongId.view', {
      body: formData,
    })
    if (!response.ok) {
      throw new Error(`Failed to fetch lyrics for track ${musicbrainzTrackId}: ${response.statusText}`)
    }
    const data = await response.json() as SubsonicLyricsResponse
    return data
  }

  return {
    backendFetchRequest,
    openSubsonicFetchRequest,
    getAlbum,
    getAlbumTracks,
    getArtistTracks,
    getArtistAlbums,
    getCurrentUser,
    getGenreTracks,
    getUsers,
    getMimeType,
    getLyrics,
  }
}
