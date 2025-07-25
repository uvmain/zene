import type { MockInstance } from 'vitest'
import type { TrackMetadata } from '~/types'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { useBackendFetch } from '../useBackendFetch'

// Mock dependencies
vi.mock('../useRandomSeed', () => ({
  useRandomSeed: () => ({
    getRandomSeed: (): number => 12345,
  }),
}))

vi.mock('../useLogic', () => ({
  useLogic: () => ({
    trackWithImageUrl: (track: TrackMetadata) => ({
      ...track,
      image_url: `/api/albums/${track.musicbrainz_album_id}/art`,
    }),
  }),
}))

describe('useBackendFetch', () => {
  let fetchSpy: MockInstance

  beforeEach(() => {
    vi.clearAllMocks()
    fetchSpy = vi.spyOn(globalThis, 'fetch')
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  describe('backendFetchRequest', () => {
    it('should make a fetch request to the correct URL', async () => {
      const mockResponse = { ok: true, json: async () => Promise.resolve({}) }
      fetchSpy.mockResolvedValue(mockResponse as Response)

      const { backendFetchRequest } = useBackendFetch()
      await backendFetchRequest('test-path')

      expect(fetchSpy).toHaveBeenCalledWith('/api/test-path', {})
    })

    it('should pass through options to fetch', async () => {
      const mockResponse = { ok: true, json: async () => Promise.resolve({}) }
      fetchSpy.mockResolvedValue(mockResponse as Response)

      const { backendFetchRequest } = useBackendFetch()
      const options = { method: 'POST', body: 'test' }
      await backendFetchRequest('test-path', options)

      expect(fetchSpy).toHaveBeenCalledWith('/api/test-path', options)
    })
  })

  describe('getAlbumTracks', () => {
    it('should fetch album tracks and return tracks with image URLs', async () => {
      const mockTracks = [
        { id: 1, musicbrainz_album_id: 'album1', title: 'Track 1' },
        { id: 2, musicbrainz_album_id: 'album1', title: 'Track 2' },
      ]
      const mockResponse = {
        json: async () => Promise.resolve(mockTracks),
      }
      fetchSpy.mockResolvedValue(mockResponse as Response)

      const { getAlbumTracks } = useBackendFetch()
      const result = await getAlbumTracks('album1')

      expect(fetchSpy).toHaveBeenCalledWith('/api/albums/album1/tracks', {})
      expect(result).toHaveLength(2)
      expect(result[0]).toEqual({
        ...mockTracks[0],
        image_url: '/api/albums/album1/art',
      })
    })
  })

  describe('getArtistTracks', () => {
    it('should fetch artist tracks with random seed and no limit', async () => {
      const mockTracks = [
        { id: 1, musicbrainz_album_id: 'album1', title: 'Track 1' },
      ]
      const mockResponse = {
        json: async () => Promise.resolve(mockTracks),
      }
      fetchSpy.mockResolvedValue(mockResponse as Response)

      const { getArtistTracks } = useBackendFetch()
      const result = await getArtistTracks('artist1')

      expect(fetchSpy).toHaveBeenCalledWith('/api/artists/artist1/tracks?random=12345', {})
      expect(result).toHaveLength(1)
      expect(result[0]).toEqual({
        ...mockTracks[0],
        image_url: '/api/albums/album1/art',
      })
    })

    it('should fetch artist tracks with limit when provided', async () => {
      const mockTracks = [
        { id: 1, musicbrainz_album_id: 'album1', title: 'Track 1' },
      ]
      const mockResponse = {
        json: async () => Promise.resolve(mockTracks),
      }
      fetchSpy.mockResolvedValue(mockResponse as Response)

      const { getArtistTracks } = useBackendFetch()
      await getArtistTracks('artist1', 10)

      expect(fetchSpy).toHaveBeenCalledWith('/api/artists/artist1/tracks?random=12345&limit=10', {})
    })
  })

  describe('getArtistAlbums', () => {
    it('should fetch artist albums chronologically', async () => {
      const mockAlbums = [
        { id: 1, title: 'Album 1', year: 2020 },
        { id: 2, title: 'Album 2', year: 2021 },
      ]
      const mockResponse = {
        json: async () => Promise.resolve(mockAlbums),
      }
      fetchSpy.mockResolvedValue(mockResponse as Response)

      const { getArtistAlbums } = useBackendFetch()
      const result = await getArtistAlbums('artist1')

      expect(fetchSpy).toHaveBeenCalledWith('/api/artists/artist1/albums?chronological=true', {})
      expect(result).toEqual(mockAlbums)
    })
  })

  describe('getCurrentUser', () => {
    it('should fetch current user data', async () => {
      const mockUser = { id: 1, username: 'testuser', email: 'test@example.com' }
      const mockResponse = {
        json: async () => Promise.resolve(mockUser),
      }
      fetchSpy.mockResolvedValue(mockResponse as Response)

      const { getCurrentUser } = useBackendFetch()
      const result = await getCurrentUser()

      expect(fetchSpy).toHaveBeenCalledWith('/api/user', {})
      expect(result).toEqual(mockUser)
    })
  })

  describe('getUsers', () => {
    it('should fetch users list', async () => {
      const mockUsers = [
        { id: 1, username: 'user1' },
        { id: 2, username: 'user2' },
      ]
      const mockResponse = {
        json: async () => Promise.resolve({ users: mockUsers }),
      }
      fetchSpy.mockResolvedValue(mockResponse as Response)

      const { getUsers } = useBackendFetch()
      const result = await getUsers()

      expect(fetchSpy).toHaveBeenCalledWith('/api/users', {})
      expect(result).toEqual(mockUsers)
    })
  })

  describe('getGenreTracks', () => {
    it('should fetch genre tracks with default parameters', async () => {
      const mockTracks = [
        { id: 1, title: 'Rock Track 1' },
        { id: 2, title: 'Rock Track 2' },
      ]
      const mockResponse = {
        json: async () => Promise.resolve(mockTracks),
      }
      fetchSpy.mockResolvedValue(mockResponse as Response)

      const { getGenreTracks } = useBackendFetch()
      const result = await getGenreTracks('rock')

      expect(fetchSpy).toHaveBeenCalledWith('/api/genres/tracks?genres=rock&limit=0&random=false', {})
      expect(result).toEqual(mockTracks)
    })

    it('should fetch genre tracks with limit and random parameters', async () => {
      const mockTracks = [{ id: 1, title: 'Rock Track 1' }]
      const mockResponse = {
        json: async () => Promise.resolve(mockTracks),
      }
      fetchSpy.mockResolvedValue(mockResponse as Response)

      const { getGenreTracks } = useBackendFetch()
      await getGenreTracks('rock', 5, true)

      expect(fetchSpy).toHaveBeenCalledWith('/api/genres/tracks?genres=rock&limit=5&random=true', {})
    })
  })

  describe('getTemporaryToken', () => {
    it('should fetch temporary token with default duration', async () => {
      const mockToken = { token: 'abc123', expires_at: '2025-07-20T12:00:00Z' }
      const mockResponse = {
        json: async () => Promise.resolve(mockToken),
      }
      fetchSpy.mockResolvedValue(mockResponse as Response)

      const { getTemporaryToken } = useBackendFetch()
      const result = await getTemporaryToken()

      expect(fetchSpy).toHaveBeenCalledWith('/api/temporary_token?duration=30', {})
      expect(result).toEqual(mockToken)
    })

    it('should fetch temporary token with custom duration', async () => {
      const mockToken = { token: 'abc123', expires_at: '2025-07-20T12:00:00Z' }
      const mockResponse = {
        json: async () => Promise.resolve(mockToken),
      }
      fetchSpy.mockResolvedValue(mockResponse as Response)

      const { getTemporaryToken } = useBackendFetch()
      await getTemporaryToken(60)

      expect(fetchSpy).toHaveBeenCalledWith('/api/temporary_token?duration=60', {})
    })
  })

  describe('refreshTemporaryToken', () => {
    it('should refresh temporary token with form data', async () => {
      const mockToken = { token: 'def456', expires_at: '2025-07-20T13:00:00Z' }
      const mockResponse = {
        json: async () => Promise.resolve(mockToken),
      }
      fetchSpy.mockResolvedValue(mockResponse as Response)

      const { refreshTemporaryToken } = useBackendFetch()
      const result = await refreshTemporaryToken('old-token')

      expect(fetchSpy).toHaveBeenCalledWith('/api/temporary_token', {
        method: 'POST',
        body: expect.any(FormData) as FormData,
      })

      // Verify FormData content
      const callArgs = fetchSpy.mock.calls[0] as [string, RequestInit]
      const formData = callArgs[1].body as FormData
      expect(formData.get('token')).toBe('old-token')
      expect(formData.get('duration')).toBe('30')
      expect(result).toEqual(mockToken)
    })

    it('should refresh temporary token with custom duration', async () => {
      const mockToken = { token: 'def456', expires_at: '2025-07-20T13:00:00Z' }
      const mockResponse = {
        json: async () => Promise.resolve(mockToken),
      }
      fetchSpy.mockResolvedValue(mockResponse as Response)

      const { refreshTemporaryToken } = useBackendFetch()
      await refreshTemporaryToken('old-token', 120)

      const callArgs = fetchSpy.mock.calls[0] as [string, RequestInit]
      const formData = callArgs[1].body as FormData
      expect(formData.get('duration')).toBe('120')
    })
  })

  describe('getMimeType', () => {
    it('should fetch mime type from Content-Type header', async () => {
      const mockResponse = {
        headers: {
          get: vi.fn((header: string) => {
            if (header === 'content-type')
              return 'audio/mpeg'
            if (header === 'Content-Type')
              return null
            return null
          }),
        },
      } as unknown as Response
      fetchSpy.mockResolvedValue(mockResponse)

      const { getMimeType } = useBackendFetch()
      const result = await getMimeType('https://example.com/song.mp3')

      expect(fetchSpy).toHaveBeenCalledWith('https://example.com/song.mp3', { method: 'HEAD' })
      expect(result).toBe('audio/mpeg')
    })

    it('should fallback to Content-Type header if content-type is not available', async () => {
      const mockResponse = {
        headers: {
          get: vi.fn((header: string) => {
            if (header === 'content-type')
              return null
            if (header === 'Content-Type')
              return 'audio/wav'
            return null
          }),
        },
      } as unknown as Response
      fetchSpy.mockResolvedValue(mockResponse)

      const { getMimeType } = useBackendFetch()
      const result = await getMimeType('https://example.com/song.wav')

      expect(result).toBe('audio/wav')
    })

    it('should return empty string if no content type headers are found', async () => {
      const mockResponse = {
        headers: {
          get: vi.fn(() => null),
        },
      } as unknown as Response
      fetchSpy.mockResolvedValue(mockResponse)

      const { getMimeType } = useBackendFetch()
      const result = await getMimeType('https://example.com/song.unknown')

      expect(result).toBe('')
    })
  })
})
