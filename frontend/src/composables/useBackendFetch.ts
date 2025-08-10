import type { AlbumMetadata, TrackMetadata, TrackMetadataWithImageUrl } from '~/types'
import type { TokenResponse, User, UsersResponse } from '~/types/auth'
import type { SubsonicUserResponse } from '~/types/getUser'
import { useAuth } from '~/composables/useAuth'
import { useRandomSeed } from '~/composables/useRandomSeed'
import { useLogic } from './useLogic'

const { userApiKey, userSalt, userToken, userUsername, userLoginState, userIsAdminState } = useAuth()
const { getRandomSeed } = useRandomSeed()
const { trackWithImageUrl } = useLogic()
const router = useRouter()

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

  const openSubsonicFetchRequest = async (path: string, options = {}): Promise<Response> => {
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

  const getTemporaryToken = async (duration = 30): Promise<TokenResponse> => {
    const formData = new FormData()
    formData.append('duration', duration.toString())
    const response = await backendFetchRequest('temporary_token', {
      method: 'POST',
      body: formData,
    })
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

  const checkIfLoggedIn = async (): Promise<boolean> => {
    try {
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
        await router.push('/login')
      }
      formData.append('f', 'json')
      formData.append('v', '1.16.0')
      formData.append('c', 'zene-frontend')

      const response = await openSubsonicFetchRequest('getUser.view', {
        method: 'POST',
        body: formData,
      })

      const json = await response.json() as SubsonicUserResponse
      const subsonicResponse = json['subsonic-response']
      if (subsonicResponse.error) {
        throw new Error(subsonicResponse.error.message)
      }
      userLoginState.value = subsonicResponse.status === 'ok'
      userIsAdminState.value = subsonicResponse.user.adminRole === 'true'
      return userLoginState.value
    }
    catch {
      userLoginState.value = false
      userIsAdminState.value = false
      return false
    }
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
    getTemporaryToken,
    refreshTemporaryToken,
    getMimeType,
    getLyrics,
    checkIfLoggedIn,
  }
}
