import type { AlbumMetadata } from '../../types'
import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import Album from '../Album.vue'

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

const mockAlbum: AlbumMetadata = {
  album_artist: 'Test album artist',
  album: 'Test Album',
  genres: ['Test Genre'],
  release_date: '2023-01-01',
  image_url: 'https://example.com/test-album.jpg',
  artist: 'Test Artist',
  musicbrainz_album_id: 'test-album-id',
  musicbrainz_artist_id: 'test-artist-id',
}

describe('album', () => {
  it('should render correctly', () => {
    const wrapper = mount(Album, {
      props: { album: mockAlbum },
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
    const wrapper = mount(Album, {
      props: { album: mockAlbum },
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
