import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import SearchResults from '../SearchResults.vue'

// Mock router
const mockRouter = {
  push: vi.fn(),
  replace: vi.fn(),
}

describe('SearchResults', () => {
  it('should render correctly', () => {
    const wrapper = mount(SearchResults, {
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
    const wrapper = mount(SearchResults, {
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
