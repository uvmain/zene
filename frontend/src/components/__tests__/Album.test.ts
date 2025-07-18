import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import Album from '../Album.vue'

// Mock router
const mockRouter = {
  push: vi.fn(),
  replace: vi.fn(),
}

const mockAlbum = {
  id: 1,
  title: 'Test Album',
  artist: 'Test Artist',
  release_date: '2023-01-01',
  musicbrainz_album_id: 'test-album-id',
  musicbrainz_artist_id: 'test-artist-id',
}

describe('Album', () => {
  it('should render correctly', () => {
    const wrapper = mount(Album, {
      props: { album: mockAlbum },
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
    const wrapper = mount(Album, {
      props: { album: mockAlbum },
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
