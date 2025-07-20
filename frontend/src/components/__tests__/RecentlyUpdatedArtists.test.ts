import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import RecentlyUpdatedArtists from '../RecentlyUpdatedArtists.vue'
import { mockArtistsResponse } from '../../../test/mocks/artists'

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

describe('recentlyUpdatedArtists', () => {
  it('should render correctly', () => {

    vi.spyOn(globalThis, 'fetch').mockResolvedValue({
      json: async () => Promise.resolve(mockArtistsResponse),
    } as Response)

    const wrapper = mount(RecentlyUpdatedArtists, {
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
    const wrapper = mount(RecentlyUpdatedArtists, {
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
