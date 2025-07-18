import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import HeroAlbum from '../HeroAlbum.vue'

// Mock router
const mockRouter = {
  push: vi.fn(),
  replace: vi.fn(),
}

describe('HeroAlbum', () => {
  it('should render correctly', () => {
    const wrapper = mount(HeroAlbum, {
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
    const wrapper = mount(HeroAlbum, {
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
