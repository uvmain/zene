import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import RecentlyAddedAlbums from '../RecentlyAddedAlbums.vue'

// Mock router
const mockRouter = {
  push: vi.fn(),
  replace: vi.fn(),
}

describe('recentlyAddedAlbums', () => {
  it('should render correctly', () => {
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
