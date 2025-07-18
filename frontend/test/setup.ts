import { beforeAll, afterEach, afterAll } from 'vitest'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import '@testing-library/jest-dom'

// Mock API handlers
const handlers = [
  // Mock authentication endpoints
  http.get('/api/check-session', () => {
    return HttpResponse.json({ loggedIn: true })
  }),
  
  http.get('/api/logout', () => {
    return HttpResponse.json({ success: true })
  }),
  
  // Mock temporary token for audio streaming
  http.get('/api/temporary_token', () => {
    return HttpResponse.json({ token: 'mock-token', expires_at: Date.now() + 30000 })
  }),
  
  // Mock albums endpoints
  http.get('/api/albums', () => {
    return HttpResponse.json([
      {
        id: 1,
        title: 'Test Album',
        artist: 'Test Artist',
        release_date: '2023-01-01',
        musicbrainz_album_id: 'test-album-id',
        musicbrainz_artist_id: 'test-artist-id',
        genres: 'Rock;Alternative',
        album_art_id: 'test-art-id',
      },
    ])
  }),
  
  http.get('/api/albums/random', () => {
    return HttpResponse.json([
      {
        id: 1,
        title: 'Random Album',
        artist: 'Random Artist',
        release_date: '2023-01-01',
        musicbrainz_album_id: 'random-album-id',
        musicbrainz_artist_id: 'random-artist-id',
        genres: 'Rock;Pop',
        album_art_id: 'random-art-id',
      },
    ])
  }),
  
  // Mock artists endpoints
  http.get('/api/artists', () => {
    return HttpResponse.json([
      {
        id: 1,
        name: 'Test Artist',
        musicbrainz_artist_id: 'test-artist-id',
        image_url: '/api/artists/test-artist-id/art',
      },
    ])
  }),
  
  // Mock tracks endpoints
  http.get('/api/tracks', () => {
    return HttpResponse.json([
      {
        id: 1,
        title: 'Test Track',
        artist: 'Test Artist',
        album: 'Test Album',
        duration: 180,
        track_number: 1,
        musicbrainz_track_id: 'test-track-id',
      },
    ])
  }),
  
  // Mock genres endpoints
  http.get('/api/genres', () => {
    return HttpResponse.json([
      {
        id: 1,
        name: 'Rock',
        count: 100,
      },
    ])
  }),
  
  // Mock search endpoint
  http.get('/api/search', () => {
    return HttpResponse.json({
      albums: [],
      artists: [],
      tracks: [],
    })
  }),
  
  // Mock settings endpoints
  http.get('/api/settings', () => {
    return HttpResponse.json({
      theme: 'dark',
      volume: 0.8,
    })
  }),
  
  http.post('/api/settings', () => {
    return HttpResponse.json({ success: true })
  }),
  
  // Mock user endpoints
  http.get('/api/users', () => {
    return HttpResponse.json([
      {
        id: 1,
        username: 'testuser',
        is_admin: false,
      },
    ])
  }),
  
  http.get('/api/current-user', () => {
    return HttpResponse.json({
      id: 1,
      username: 'testuser',
      is_admin: false,
    })
  }),
  
  // Mock art endpoints
  http.get(/\/api\/albums\/.*\/art/, () => {
    return new HttpResponse('mock-image-data', {
      headers: { 'Content-Type': 'image/jpeg' },
    })
  }),
  
  http.get(/\/api\/artists\/.*\/art/, () => {
    return new HttpResponse('mock-image-data', {
      headers: { 'Content-Type': 'image/jpeg' },
    })
  }),
]

// Setup MSW server
const server = setupServer(...handlers)

// Setup and teardown
beforeAll(() => {
  server.listen({ onUnhandledRequest: 'warn' })
})

afterEach(() => {
  server.resetHandlers()
})

afterAll(() => {
  server.close()
})

// Global test utilities
global.server = server

// Mock window.matchMedia for CSS queries
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(), // deprecated
    removeListener: vi.fn(), // deprecated
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
})

// Mock IntersectionObserver
global.IntersectionObserver = vi.fn().mockImplementation(() => ({
  observe: vi.fn(),
  unobserve: vi.fn(),
  disconnect: vi.fn(),
}))

// Mock ResizeObserver
global.ResizeObserver = vi.fn().mockImplementation(() => ({
  observe: vi.fn(),
  unobserve: vi.fn(),
  disconnect: vi.fn(),
}))