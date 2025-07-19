import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import Tracks from '../Tracks.vue'

// Mock router
const mockRouter = {
  push: vi.fn(),
  replace: vi.fn(),
}

const mockTracks = [
  {
    id: 1,
    title: 'Test Track',
    artist: 'Test Artist',
    album: 'Test Album',
    duration: 180,
    track_number: 1,
    total_tracks: 10,
    musicbrainz_track_id: 'test-track-id',
    musicbrainz_album_id: 'test-album-id',
    musicbrainz_artist_id: 'test-artist-id',
    image_url: '/api/albums/test-album-id/art',
  },
]

describe('tracks', () => {
  it('should render correctly', () => {
    const wrapper = mount(Tracks, {
      props: { tracks: mockTracks },
      global: {
        mocks: {
          $router: mockRouter,
          $route: { path: '/', params: {}, query: {} },
        },
        stubs: {
          RouterLink: true,
          RouterView: true,
        },
      },
    })
    expect(wrapper.exists()).toBe(true)
  })

  it('should be a Vue instance', () => {
    const wrapper = mount(Tracks, {
      props: { tracks: mockTracks },
      global: {
        mocks: {
          $router: mockRouter,
          $route: { path: '/', params: {}, query: {} },
        },
        stubs: {
          RouterLink: true,
          RouterView: true,
        },
      },
    })
    expect(wrapper.vm).toBeTruthy()
  })
})
