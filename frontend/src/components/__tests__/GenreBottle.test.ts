import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import GenreBottle from '../GenreBottle.vue'

// Mock router
const mockRouter = {
  push: vi.fn(),
  replace: vi.fn(),
}

const mockGenre = {
  id: 1,
  name: 'Rock',
  count: 100,
}

describe('GenreBottle', () => {
  it('should render correctly', () => {
    const wrapper = mount(GenreBottle, {
      props: { genre: mockGenre },
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
    const wrapper = mount(GenreBottle, {
      props: { genre: mockGenre },
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
