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

  describe('cast Media Info Change Logic', () => {
    it('should extract track ID from media URL correctly', () => {
      // Test the regex pattern used in onCastMediaInfoChanged
      const testUrls = [
        '/api/tracks/test-track-123/audio?quality=medium',
        '/api/v1/tracks/another-track-456/audio',
        '/tracks/simple-track/audio',
      ]

      const expectedIds = ['test-track-123', 'another-track-456', 'simple-track']

      testUrls.forEach((url, index) => {
        const trackIdMatch = url.match(/\/tracks\/([^/]+)\/audio/)
        expect(trackIdMatch).toBeTruthy()
        expect(trackIdMatch![1]).toBe(expectedIds[index])
      })
    })

    it('should handle invalid URLs gracefully', () => {
      const invalidUrls = [
        '/invalid/url/format',
        '/api/not-tracks/test/audio',
        '',
        null,
        undefined,
      ]

      invalidUrls.forEach((url) => {
        const trackIdMatch = url?.match(/\/tracks\/([^/]+)\/audio/)
        expect(trackIdMatch).toBeFalsy()
      })
    })

    it('should have the onCastMediaInfoChanged method defined', () => {
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

      // Check that the method exists on the component instance
      // eslint-disable-next-line ts/no-unsafe-member-access
      expect((wrapper.vm as any).onCastMediaInfoChanged).toBeDefined()
      // eslint-disable-next-line ts/no-unsafe-member-access
      expect(typeof (wrapper.vm as any).onCastMediaInfoChanged).toBe('function')
    })
  })
})
