import type {
  SubsonicAlbumListResponse,
  SubsonicAlbumResponse,
  SubsonicApiKeyResponse,
  SubsonicArtistInfoResponse,
  SubsonicArtistListResponse,
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
import type { SearchResult } from '~/types'
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import type { SubsonicArtist, SubsonicArtistInfo, SubsonicIndexArtist } from '~/types/subsonicArtist'
import type { SubsonicGenre } from '~/types/subsonicGenres'
import type { StructuredLyric } from '~/types/subsonicLyrics'
import type { SubsonicPodcastChannel } from '~/types/subsonicPodcasts'
import type { SubsonicSong } from '~/types/subsonicSong'
import { useLocalStorage } from '@vueuse/core'
import { debugLog } from '~/logic/logger'

const apiKey = useLocalStorage('apiKey', '')

const concurrencyMap = new Map<string, Promise<any>>()

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

export async function openSubsonicFetchRequest<T>(path: string, options: RequestInit = {}): Promise<T> {
  if (apiKey.value == null || apiKey.value.length === 0) {
    const router = useRouter()
    await router.push('/login')
  }

  const formDataEntries = [] as [string, string][]
  (options?.body as FormData)?.forEach((value: FormDataEntryValue, key: string) => {
    formDataEntries.push([key, value.toString()])
  })

  const concurrencyKey = `${path}|${options.method}|${JSON.stringify(formDataEntries)}|${JSON.stringify(options.body)}`
  if (concurrencyMap.has(concurrencyKey)) {
    return concurrencyMap.get(concurrencyKey) as Promise<T>
  }

  const formData = new FormData()
  formData.append('apiKey', apiKey.value)
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

  const promise = async <T>(path: string, options: RequestInit): Promise<T> => {
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
        return jsonData['subsonic-response'] as T
      }
    }
    catch (error) {
      console.error('Error fetching data:', error)
    }
    return {} as T
  }

  const promiseInstance = promise(path, options)
  concurrencyMap.set(concurrencyKey, promiseInstance)
  void promiseInstance.finally(() => {
    concurrencyMap.delete(concurrencyKey)
  })

  return promiseInstance as Promise<T>
}

export function getStreamUrl(path: string, params: URLSearchParams = {} as URLSearchParams): string {
  if (apiKey.value !== null && apiKey.value.length > 0) {
    params.append('apiKey', apiKey.value)
    params.append('f', 'json')
    params.append('v', '1.16.0')
    params.append('c', 'zene-frontend')
  }
  else {
    return ''
  }

  return `/rest/${path}?${params.toString()}`
}

export async function fetchAlbum(musicbrainz_album_id: string): Promise<SubsonicAlbum> {
  const formData = new FormData()
  formData.append('id', musicbrainz_album_id)
  const response = await openSubsonicFetchRequest<SubsonicAlbumResponse>('getAlbum', {
    body: formData,
  })
  return response.album
}

export async function fetchAlbums(type: string, size = 50, offset = 0, seed?: number): Promise<SubsonicAlbum[]> {
  const formData = new FormData()
  formData.append('type', type)
  formData.append('size', size.toString())
  formData.append('offset', offset.toString())
  if (type === 'random' && seed !== undefined && seed > 0) {
    formData.append('seed', seed.toString())
  }
  const response = await openSubsonicFetchRequest<SubsonicAlbumListResponse>('getAlbumList', {
    body: formData,
  })
  return response.albumList.album
}

export async function fetchRandomTracks(limit?: number, offset?: number, seed?: number): Promise<SubsonicSong[]> {
  const options: RequestInit = {}

  if (limit != null && limit > 0) {
    const formData = new FormData()
    if (offset !== undefined && offset > 0) {
      formData.append('offset', offset.toString())
    }
    if (limit !== undefined && limit > 0) {
      formData.append('size', limit.toString())
    }
    if (seed !== undefined && seed > 0) {
      formData.append('seed', seed.toString())
    }
    options.body = formData
  }
  const response = await openSubsonicFetchRequest<SubsonicRandomSongsResponse>('getRandomSongs', options)
  return response.randomSongs.song
}

export async function fetchArtist(musicbrainz_artist_id: string): Promise<SubsonicArtist> {
  const formData = new FormData()
  formData.append('id', musicbrainz_artist_id)
  const response = await openSubsonicFetchRequest<SubsonicArtistResponse>('getArtist', {
    body: formData,
  })
  return response.artist
}

export async function fetchArtistInfo(musicbrainz_artist_id: string, limit = 20): Promise<SubsonicArtistInfo> {
  const formData = new FormData()
  formData.append('id', musicbrainz_artist_id)
  formData.append('count', limit.toString())
  const response = await openSubsonicFetchRequest<SubsonicArtistInfoResponse>('getArtistInfo', {
    body: formData,
  })
  return response.artistInfo
}

export async function fetchArtists(limit = 0): Promise<SubsonicIndexArtist[]> {
  const response = await openSubsonicFetchRequest<SubsonicArtistsResponse>('getArtists')
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

export async function fetchArtistList(type: string, limit: number, offset: number, seed?: number): Promise<SubsonicArtist[]> {
  const formData = new FormData()
  formData.append('type', type)
  formData.append('size', limit.toString())
  formData.append('offset', offset.toString())
  if (seed !== undefined && seed > 0) {
    formData.append('seed', seed.toString())
  }
  const response = await openSubsonicFetchRequest<SubsonicArtistListResponse>('getArtistList', {
    body: formData,
  })
  const artists = response.artistList.artist
  return artists
}

export async function fetchArtistTopSongs(musicbrainz_artist_id: string, limit = 50, offset = 0): Promise<SubsonicSong[]> {
  const formData = new FormData()
  formData.append('id', musicbrainz_artist_id)
  formData.append('count', limit.toString())
  if (offset > 0) {
    formData.append('offset', offset.toString())
  }
  const response = await openSubsonicFetchRequest<SubsonicTopSongsResponse>('getTopSongs', {
    body: formData,
  })
  return response.topSongs.song
}

export async function postScrobble(musicbrainz_track_id: string): Promise<boolean> {
  const formData = new FormData()
  formData.append('id', musicbrainz_track_id)
  const response = await openSubsonicFetchRequest<SubsonicResponse>('scrobble', {
    body: formData,
  })
  return response.status === 'ok'
}

export async function fetchGenres(count?: number) {
  const response = await openSubsonicFetchRequest<SubsonicGenresResponse>('getGenres')
  const allGenres = response.genres.genre
  if (allGenres.length === 0) {
    return [] as SubsonicGenre[]
  }
  if (count !== undefined) {
    return allGenres.slice(0, count)
  }
  return allGenres
}

export async function fetchLyrics(musicbrainz_song_id: string): Promise<StructuredLyric | null> {
  const formData = new FormData()
  formData.append('id', musicbrainz_song_id)
  const response = await openSubsonicFetchRequest<SubsonicLyricsListResponse>('getLyricsBySongId', {
    body: formData,
  })
  if (response.lyricsList?.structuredLyrics?.length > 0) {
    return response.lyricsList.structuredLyrics[0]
  }
  else {
    return null
  }
}

export async function fetchAlbumsForArtist(artistId: string): Promise<SubsonicAlbum[]> {
  const formData = new FormData()
  formData.append('id', artistId)
  const response = await openSubsonicFetchRequest<SubsonicArtistResponse>('getArtist', {
    body: formData,
  })
  const albums: SubsonicAlbum[] = []
  for (const album of response.artist.album) {
    const albumWithSongs = await fetchAlbum(album.id)
    albums.push(albumWithSongs)
  }
  return albums
}

export async function fetchSongsByGenre(genre: string, limit: number, offset: number): Promise<SubsonicSong[]> {
  const formData = new FormData()
  formData.append('genre', genre)
  formData.append('count', limit.toString())
  formData.append('offset', offset.toString())
  const response = await openSubsonicFetchRequest<SubsonicSongsByGenreResponse>('getSongsByGenre', {
    body: formData,
  })
  return response.songsByGenre.song
}

export async function fetchSong(musicbrainz_track_id: string): Promise<SubsonicSong> {
  const formData = new FormData()
  formData.append('id', musicbrainz_track_id)
  const response = await openSubsonicFetchRequest<SubsonicSongResponse>('getSong', {
    body: formData,
  })
  return response.song
}

export async function fetchSearchResults(query: string, limit = 50): Promise<SearchResult> {
  const formData = new FormData()
  formData.append('query', query)
  formData.append('artistcount', limit.toString())
  formData.append('albumcount', limit.toString())
  formData.append('songcount', limit.toString())
  const response = await openSubsonicFetchRequest<SubsonicSearchResponse>('search2', {
    body: formData,
  })
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

export type AlbumArtSseMessage = | { source: 'Deezer', data: string }
  | { source: 'CoverArtArchive', data: string }
  | { source: 'LocalArt', data: { folderArt: string, embeddedArt: string } }

export async function useServerSentEventsForAlbumArt(artist: string, album: string, onMessage: (data: AlbumArtSseMessage) => void, onError: (error: any) => void): Promise<EventSource> {
  const params = new URLSearchParams()
  params.append('apiKey', apiKey.value)
  params.append('f', 'json')
  params.append('v', '1.16.0')
  params.append('c', 'zene-frontend')
  params.append('artist', artist)
  params.append('album', album)

  const eventSource = new EventSource(`/rest/getalbumartssse?${params.toString()}`)

  eventSource.addEventListener('message', (event) => {
    const data = JSON.parse(event.data as string) as AlbumArtSseMessage
    onMessage(data)
  })

  eventSource.addEventListener('done', () => {
    console.log('SSE completed — closing stream')
    eventSource.close()
  })

  eventSource.onerror = (err) => {
    onError(err)
  }

  return eventSource
}

export async function postNewAlbumArt(musicbrainz_song_id: string, image: Blob): Promise<SubsonicResponse> {
  const formData = new FormData()
  formData.append('id', musicbrainz_song_id)
  formData.append('file', image)

  const response = await openSubsonicFetchRequest<SubsonicResponse>('updateAlbumArt', {
    body: formData,
  })
  return response
}

export async function useServerSentEventsForPodcast(podcastId: string, onMessage: (data: SubsonicPodcastChannel) => void, onError: (error: any) => void): Promise<EventSource> {
  const params = new URLSearchParams()
  params.append('apiKey', apiKey.value)
  params.append('f', 'json')
  params.append('v', '1.16.0')
  params.append('c', 'zene-frontend')
  params.append('id', podcastId)

  const eventSource = new EventSource(`/rest/getpodcastssse?${params.toString()}`)

  eventSource.addEventListener('message', (event) => {
    const data = JSON.parse(event.data as string) as SubsonicPodcastChannel
    onMessage(data)
  })

  eventSource.addEventListener('done', () => {
    console.log('SSE completed — closing stream')
    eventSource.close()
  })

  eventSource.onerror = (err) => {
    onError(err)
  }

  return eventSource
}

export async function downloadMediaBlob(mediaId: string): Promise<Blob> {
  const formData = new FormData()
  formData.append('apiKey', apiKey.value)
  formData.append('f', 'json')
  formData.append('v', '1.16.0')
  formData.append('c', 'zene-frontend')
  formData.append('id', mediaId)

  const url = '/rest/download.view'
  const response = await fetch(url, {
    method: 'POST',
    body: formData,
  })

  if (!response.ok) {
    throw new Error(`Failed to download media from ${url}: ${response.statusText}`)
  }
  const blob = await response.blob()
  return blob
}

export async function postTrackStarred(musicbrainz_track_id: string, starred: boolean): Promise<SubsonicResponse> {
  const formData = new FormData()
  formData.append('id', musicbrainz_track_id)

  if (starred) {
    return openSubsonicFetchRequest<SubsonicResponse>('star', {
      body: formData,
    })
  }
  else {
    return openSubsonicFetchRequest<SubsonicResponse>('unstar', {
      body: formData,
    })
  }
}
