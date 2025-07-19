import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import Navbar from '../Navbar.vue'

// Mock Vue Router composables
vi.mock('vue-router', () => ({
  useRoute: vi.fn(() => ({
    path: '/',
    params: {},
    query: {},
  })),
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

describe('navbar', () => {
  it('should render correctly', () => {
    const wrapper = mount(Navbar, {
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
    const wrapper = mount(Navbar, {
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
