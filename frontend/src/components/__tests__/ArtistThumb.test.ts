import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import ArtistThumb from '../ArtistThumb.vue'

// Mock router
const mockRouter = {
  push: vi.fn(),
  replace: vi.fn(),
}

const mockArtist = {
  id: 1,
  name: 'Test Artist',
  musicbrainz_artist_id: 'test-artist-id',
  image_url: '/api/artists/test-artist-id/art',
}

describe('ArtistThumb', () => {
  it('should render correctly', () => {
    const wrapper = mount(ArtistThumb, {
      props: { artist: mockArtist },
      global: {
        mocks: {
          $router: mockRouter,
          $route: { path: '/', params: {}, query: {} },
        },
        stubs: {
          'RouterLink': true,
          'RouterView': true,
        },
      },
    })
    expect(wrapper.exists()).toBe(true)
  })

  it('should be a Vue instance', () => {
    const wrapper = mount(ArtistThumb, {
      props: { artist: mockArtist },
      global: {
        mocks: {
          $router: mockRouter,
          $route: { path: '/', params: {}, query: {} },
        },
        stubs: {
          'RouterLink': true,
          'RouterView': true,
        },
      },
    })
    expect(wrapper.vm).toBeTruthy()
  })
})
