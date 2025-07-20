import { mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import GenreBottle from '../GenreBottle.vue'

// Mock Vue Router composables
const mockRouterPush = vi.fn()
vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: mockRouterPush,
  }),
}))

// Mock useSearch composable
vi.mock('../composables/useSearch', () => ({
  useSearch: () => ({
    closeSearch: vi.fn(),
  }),
}))

const mockGenre = 'Rock'

describe('genreBottle', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('should render correctly', () => {
    const wrapper = mount(GenreBottle, {
      props: { genre: mockGenre },
    })
    expect(wrapper.exists()).toBe(true)
  })

  it('should be a Vue instance', () => {
    const wrapper = mount(GenreBottle, {
      props: { genre: mockGenre },
    })
    expect(wrapper.vm).toBeTruthy()
  })

  it('should display the genre name', () => {
    const wrapper = mount(GenreBottle, {
      props: { genre: mockGenre },
    })
    expect(wrapper.text()).toContain(mockGenre)
  })

  it('should navigate to genre page on click', async () => {
    const wrapper = mount(GenreBottle, {
      props: { genre: mockGenre },
    })

    await wrapper.trigger('click')
    expect(mockRouterPush).toHaveBeenCalledWith(`/genres/${mockGenre}`)
  })
})
