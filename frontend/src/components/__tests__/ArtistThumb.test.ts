import type { ArtistMetadata } from '~/types'
import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import ArtistThumb from '../ArtistThumb.vue'

// Mock Vue Router composables
vi.mock('vue-router', () => ({
  useRouter: vi.fn(() => ({
    push: vi.fn(),
    replace: vi.fn(),
  })),
}))

// Mock router
const mockRouter = {
  push: vi.fn(),
  replace: vi.fn(),
}

const mockArtist: ArtistMetadata = {
  artist: 'Test Artist',
  musicbrainz_artist_id: 'test-artist-id',
  image_url: '/api/artists/test-artist-id/art',
}

describe('artistThumb', () => {
  it('should render correctly', () => {
    const wrapper = mount(ArtistThumb, {
      props: { artist: mockArtist },
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
    const wrapper = mount(ArtistThumb, {
      props: { artist: mockArtist },
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
