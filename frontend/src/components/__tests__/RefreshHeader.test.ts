import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import RefreshHeader from '../RefreshHeader.vue'

// Mock router
const mockRouter = {
  push: vi.fn(),
  replace: vi.fn(),
}

describe('refreshHeader', () => {
  it('should render correctly', () => {
    const wrapper = mount(RefreshHeader, {
      props: { title: 'Test Title' },
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
    const wrapper = mount(RefreshHeader, {
      props: { title: 'Test Title' },
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
