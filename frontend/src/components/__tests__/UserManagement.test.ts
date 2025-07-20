import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import { mockUserResponse } from '../../../test/mocks/user'
import UserManagement from '../UserManagement.vue'

// Mock router
const mockRouter = {
  push: vi.fn(),
  replace: vi.fn(),
}

describe('userManagement', () => {
  it('should render correctly', () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue({
      json: async () => Promise.resolve(mockUserResponse),
    } as Response)

    const wrapper = mount(UserManagement, {
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
    const wrapper = mount(UserManagement, {
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
