import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import { mockAlbumsResponse } from '../../../test/mocks/albums'
import RecentlyAddedAlbums from '../Albums.vue'

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

describe('recentlyAddedAlbums', () => {
  it('should render correctly', () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue({
      json: async () => Promise.resolve(mockAlbumsResponse),
    } as Response)

    const wrapper = mount(RecentlyAddedAlbums, {
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
    const wrapper = mount(RecentlyAddedAlbums, {
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
