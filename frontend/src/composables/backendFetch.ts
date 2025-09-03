import type {
  SubsonicAlbumListResponse,
  SubsonicAlbumResponse,
  SubsonicApiKeyResponse,
  SubsonicArtistResponse,
  SubsonicArtistsResponse,
  SubsonicGenresResponse,
  SubsonicLyricsListResponse,
  SubsonicRandomSongsResponse,
  SubsonicResponse,
  SubsonicResponseWrapper,
  SubsonicSearchResponse,
  SubsonicSongResponse,
  SubsonicSongsByGenreResponse,
  SubsonicTopSongsResponse,
} from '../types/subsonic'
import type { SubsonicAlbum } from '../types/subsonicAlbum'
import type { SearchResult } from '~/types'
import type { SubsonicArtist, SubsonicIndexArtist } from '~/types/subsonicArtist'
import type { SubsonicGenre } from '~/types/subsonicGenres'
import type { StructuredLyric } from '~/types/subsonicLyrics'
import type { SubsonicSong } from '~/types/subsonicSong'
import { useLocalStorage } from '@vueuse/core'
import { useDebug } from '~/composables/useDebug'

const { debugLog } = useDebug()
const apiKey = useLocalStorage('apiKey', '')

export async function fetchNewApiKeyWithTokenAndSalt(username: string, token: string, salt: string): Promise<string> {
  try {
    const formData = new FormData()
    formData.append('u', username)
    formData.append('t', token)
    formData.append('s', salt)
    formData.append('v', '1.16.1')
    formData.append('c', 'zeneclient')
    formData.append('f', 'json')

    const url = 'rest/createApiKey.view'
    const response = await fetch(url, {
      method: 'POST',
      body: formData,
    })
    const data = await response.json() as SubsonicResponseWrapper

    if (data['subsonic-response'].status !== 'ok') {
      throw new Error(data['subsonic-response'].error?.message ?? 'Failed to create new API key')
    }
    const apiKeysResponse = data['subsonic-response'] as SubsonicApiKeyResponse
    return apiKeysResponse.apiKeys.apiKey[0].api_key
  }
  catch (error) {
    debugLog(error as string)
    return ''
  }
}

export async function fetchApiKeysWithTokenAndSalt(username: string, token: string, salt: string): Promise<SubsonicApiKeyResponse> {
  try {
    const formData = new FormData()
    formData.append('u', username)
    formData.append('t', token)
    formData.append('s', salt)
    formData.append('v', '1.16.1')
    formData.append('c', 'zeneclient')
    formData.append('f', 'json')

    const url = 'rest/getApiKeys.view'
    const response = await fetch(url, {
      method: 'POST',
      body: formData,
    })
    const data = await response.json() as SubsonicResponseWrapper

    if (data['subsonic-response'].status !== 'ok') {
      throw new Error(data['subsonic-response'].error?.message ?? 'Failed to fetch existing API keys')
    }
    return data['subsonic-response'] as SubsonicApiKeyResponse
  }
  catch (error) {
    debugLog(error as string)
    return {} as SubsonicApiKeyResponse
  }
}

export async function openSubsonicFetchRequest(path: string, options: RequestInit = {}): Promise<SubsonicResponse> {
  const formData = new FormData()
  if (apiKey.value !== null && apiKey.value.length > 0) {
    formData.append('apiKey', apiKey.value)
    formData.append('f', 'json')
    formData.append('v', '1.16.0')
    formData.append('c', 'zene-frontend')
  }
  else {
    const router = useRouter()
    await router.push('/login')
  }

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

  try {
    const jsonData = await response.json() as SubsonicResponseWrapper
    if (jsonData['subsonic-response']?.status === 'error' && [40, 44].includes(jsonData['subsonic-response']?.error?.code ?? 0)) {
      // user is not authenticated
      apiKey.value = ''
      const router = useRouter()
      await router.push('/login')
    }
    else {
      return jsonData['subsonic-response']
    }
  }
  catch (error) {
    console.error('Error fetching data:', error)
  }
  return {} as SubsonicResponse
}

export async function fetchAlbum(musicbrainz_album_id: string): Promise<SubsonicAlbum> {
  const formData = new FormData()
  formData.append('id', musicbrainz_album_id)
  const response = await openSubsonicFetchRequest('getAlbum', {
    body: formData,
  }) as SubsonicAlbumResponse
  return response.album
}

export async function fetchAlbums(type: string, size = 50, offset = 0): Promise<SubsonicAlbum[]> {
  const formData = new FormData()
  formData.append('type', type)
  formData.append('size', size.toString())
  formData.append('offset', offset.toString())
  const response = await openSubsonicFetchRequest('getAlbumList', {
    body: formData,
  }) as SubsonicAlbumListResponse
  return response.albumList.album
}

export async function fetchRandomTracks(size?: number): Promise<SubsonicSong[]> {
  const options: RequestInit = {}

  if (size != null && size > 0) {
    const formData = new FormData()
    formData.append('size', size.toString())
    options.body = formData
  }
  const response = await openSubsonicFetchRequest('getRandomSongs', options) as SubsonicRandomSongsResponse
  return response.randomSongs.song
}

export async function fetchArtist(musicbrainz_artist_id: string): Promise<SubsonicArtist> {
  const formData = new FormData()
  formData.append('id', musicbrainz_artist_id)
  const response = await openSubsonicFetchRequest('getArtist', {
    body: formData,
  }) as SubsonicArtistResponse
  return response.artist
}

export async function fetchArtists(limit = 0): Promise<SubsonicIndexArtist[]> {
  const response = await openSubsonicFetchRequest('getArtists') as SubsonicArtistsResponse
  const artists = response.artists
  const artistArray: SubsonicIndexArtist[] = []
  for (const index of artists.index) {
    for (const artist of index.artist) {
      artistArray.push(artist)
    }
  }
  if (limit > 0) {
    return artistArray.slice(0, limit)
  }
  return artistArray
}

export async function fetchArtistTopSongs(musicbrainz_artist_id: string, count = 50): Promise<SubsonicSong[]> {
  const formData = new FormData()
  formData.append('id', musicbrainz_artist_id)
  formData.append('count', count.toString())
  const response = await openSubsonicFetchRequest('getTopSongs', {
    body: formData,
  }) as SubsonicTopSongsResponse
  return response.topSongs.song
}

export async function postScrobble(musicbrainz_track_id: string): Promise<boolean> {
  const formData = new FormData()
  formData.append('id', musicbrainz_track_id)
  const response = await openSubsonicFetchRequest('scrobble', {
    body: formData,
  })
  return response.status === 'ok'
}

export async function fetchGenres(count?: number) {
  const response = await openSubsonicFetchRequest('getGenres') as SubsonicGenresResponse
  const allGenres = response.genres.genre
  if (allGenres.length === 0) {
    return [] as SubsonicGenre[]
  }
  if (count !== undefined) {
    return allGenres.slice(0, count)
  }
  return allGenres
}

export async function fetchLyrics(musicbrainz_song_id: string): Promise<StructuredLyric> {
  const formData = new FormData()
  formData.append('id', musicbrainz_song_id)
  const response = await openSubsonicFetchRequest('getLyricsBySongId', {
    body: formData,
  }) as SubsonicLyricsListResponse
  return response.lyricsList.structuredLyrics[0]
}

export async function fetchAlbumsForArtist(artistId: string): Promise<SubsonicAlbum[]> {
  const formData = new FormData()
  formData.append('id', artistId)
  const response = await openSubsonicFetchRequest('getArtist', {
    body: formData,
  }) as SubsonicArtistResponse
  const albums: SubsonicAlbum[] = []
  for (const album of response.artist.album) {
    const albumWithSongs = await fetchAlbum(album.id)
    albums.push(albumWithSongs)
  }
  return albums
}

export async function fetchSongsByGenre(genre: string, count: number, offset: number): Promise<SubsonicSong[]> {
  const formData = new FormData()
  formData.append('genre', genre)
  formData.append('count', count.toString())
  formData.append('offset', offset.toString())
  const response = await openSubsonicFetchRequest('getSongsByGenre', {
    body: formData,
  }) as SubsonicSongsByGenreResponse
  return response.songsByGenre.song
}

export async function fetchSong(musicbrainz_track_id: string): Promise<SubsonicSong> {
  const formData = new FormData()
  formData.append('id', musicbrainz_track_id)
  const response = await openSubsonicFetchRequest('getSong', {
    body: formData,
  }) as SubsonicSongResponse
  return response.song
}

export async function fetchSearchResults(query: string, limit = 50): Promise<SearchResult> {
  const formData = new FormData()
  formData.append('query', query)
  formData.append('artistcount', limit.toString())
  formData.append('albumcount', limit.toString())
  formData.append('songcount', limit.toString())
  const response = await openSubsonicFetchRequest('search2', {
    body: formData,
  }) as SubsonicSearchResponse
  const artists = response.searchResult2.artist
  const albums = response.searchResult2.album
  const songs = response.searchResult2.song

  const allGenres = await fetchGenres()
  const searchedGenres = allGenres.filter(genre => genre.value.toLowerCase().includes(query.toLowerCase())).splice(0, limit)
  return {
    artists,
    albums,
    songs,
    genres: searchedGenres,
  }
}
