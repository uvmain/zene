import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import FooterPlayer from '../FooterPlayer.vue'

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

describe('footerPlayer', () => {
  it('should render correctly', () => {
    const mockResponse = {
      Token: 'mock-token',
      ExpiresAt: '2023-10-01T00:00:00Z',
    }
    vi.spyOn(globalThis, 'fetch').mockResolvedValue({
      json: async () => Promise.resolve(mockResponse),
    } as Response)

    const wrapper = mount(FooterPlayer, {
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
    const wrapper = mount(FooterPlayer, {
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
