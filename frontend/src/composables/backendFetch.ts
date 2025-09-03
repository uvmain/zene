import type {
  SubsonicAlbumListResponse,
  SubsonicAlbumResponse,
  SubsonicArtistResponse,
  SubsonicArtistsResponse,
  SubsonicGenresResponse,
  SubsonicLyricsListResponse,
  SubsonicRandomSongsResponse,
  SubsonicResponse,
  SubsonicSongResponse,
  SubsonicSongsByGenreResponse,
  SubsonicTopSongsResponse,
} from '../types/subsonic'
import type { SubsonicAlbum } from '../types/subsonicAlbum'
import type { SubsonicArtist } from '~/types/subsonicArtist'
import type { StructuredLyric } from '~/types/subsonicLyrics'
import type { SubsonicSong } from '~/types/subsonicSong'

export async function openSubsonicFetchRequest(path: string, options: RequestInit = {}): Promise<Response> {
  const apiKey = localStorage.getItem('apiKey')
  const formData = new FormData()
  if (apiKey !== null && apiKey?.length > 0) {
    formData.append('apiKey', apiKey)
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
    const jsonData = await response.json() as SubsonicResponse
    if (jsonData['subsonic-response']?.status === 'error' && [40, 44].includes(jsonData['subsonic-response']?.error?.code ?? 0)) {
      // user is not authenticated
      localStorage.removeItem('apiKey')
      const router = useRouter()
      await router.push('/login')
    }
  }
  catch (error) {
    console.error('Error fetching data:', error)
  }
  return response
}

export async function fetchAlbum(musicbrainz_album_id: string): Promise<SubsonicAlbum> {
  const formData = new FormData()
  formData.append('id', musicbrainz_album_id)
  const response = await openSubsonicFetchRequest('getAlbum', {
    body: formData,
  })
  const json = await response.json() as SubsonicAlbumResponse
  const album = json.album
  return album
}

export async function fetchAlbums(type: string, size = 50, offset = 0): Promise<SubsonicAlbum[]> {
  const formData = new FormData()
  formData.append('type', type)
  formData.append('size', size.toString())
  formData.append('offset', offset.toString())
  const response = await openSubsonicFetchRequest('GetAlbumList', {
    body: formData,
  })
  const json = await response.json() as SubsonicAlbumListResponse
  const albums = json.albumList.album
  return albums
}

export async function fetchRandomTracks(size?: number): Promise<SubsonicSong[]> {
  const options: RequestInit = {}

  if (size != null && size > 0) {
    const formData = new FormData()
    formData.append('size', size.toString())
    options.body = formData
  }
  const response = await openSubsonicFetchRequest('getRandomSongs', options)
  const json = await response.json() as SubsonicRandomSongsResponse
  const randomTracks = json.randomSongs.song
  return randomTracks
}

export async function fetchArtist(musicbrainz_artist_id: string): Promise<SubsonicArtist> {
  const formData = new FormData()
  formData.append('id', musicbrainz_artist_id)
  const response = await openSubsonicFetchRequest('getArtist', {
    body: formData,
  })
  const json = await response.json() as SubsonicArtistResponse
  const artist = json.artist
  return artist
}

export async function fetchArtists(limit = 0): Promise<SubsonicArtist[]> {
  const response = await openSubsonicFetchRequest('getArtists')
  const json = await response.json() as SubsonicArtistsResponse
  const artists = json.artists
  const artistArray: SubsonicArtist[] = []
  for (const index of artists.index) {
    for (const artist of index.artist) {
      const data = await fetchArtist(artist.musicBrainzId)
      artistArray.push(data)
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
  })
  const json = await response.json() as SubsonicTopSongsResponse
  const topSongs = json.topSongs.song
  return topSongs
}

export async function postScrobble(musicbrainz_track_id: string): Promise<boolean> {
  const formData = new FormData()
  formData.append('id', musicbrainz_track_id)
  const response = await openSubsonicFetchRequest('scrobble', {
    body: formData,
  })
  return response.ok
}

export async function fetchGenres(count?: number) {
  const response = await openSubsonicFetchRequest('getGenres')
  const json = await response.json() as SubsonicGenresResponse
  const allGenres = json.genres.genre
  if (count !== undefined) {
    return allGenres.slice(0, count)
  }
  return allGenres
}

export async function fetchLyrics(musicbrainz_artist_id: string): Promise<StructuredLyric> {
  const formData = new FormData()
  formData.append('id', musicbrainz_artist_id)
  const response = await openSubsonicFetchRequest('getArtist', {
    body: formData,
  })
  const json = await response.json() as SubsonicLyricsListResponse
  return json.lyricsList.structuredLyrics[0]
}

export async function fetchAlbumsForArtist(artistId: string): Promise<SubsonicAlbum[]> {
  const formData = new FormData()
  formData.append('id', artistId)
  const response = await openSubsonicFetchRequest('getArtist', {
    body: formData,
  })
  const json = await response.json() as SubsonicArtistResponse
  const albums: SubsonicAlbum[] = []
  for (const album of json.artist.album) {
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
  })
  const json = await response.json() as SubsonicSongsByGenreResponse
  return json.songsByGenre.song
}

export async function fetchSong(musicbrainz_track_id: string): Promise<SubsonicSong> {
  const formData = new FormData()
  formData.append('id', musicbrainz_track_id)
  const response = await openSubsonicFetchRequest('getSong', {
    body: formData,
  })
  const json = await response.json() as SubsonicSongResponse
  return json.song
}
