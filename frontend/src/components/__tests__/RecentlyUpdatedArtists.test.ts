import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import RecentlyUpdatedArtists from '../RecentlyUpdatedArtists.vue'

// Mock router
const mockRouter = {
  push: vi.fn(),
  replace: vi.fn(),
}

describe('RecentlyUpdatedArtists', () => {
  it('should render correctly', () => {
    const wrapper = mount(RecentlyUpdatedArtists, {
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
    const wrapper = mount(RecentlyUpdatedArtists, {
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
