import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import RecentlyUpdatedArtists from '../RecentlyUpdatedArtists.vue'

// Mock router
const mockRouter = {
  push: vi.fn(),
  replace: vi.fn(),
}

describe('recentlyUpdatedArtists', () => {
  it('should render correctly', () => {
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
