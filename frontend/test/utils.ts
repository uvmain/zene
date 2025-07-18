import { mount, VueWrapper } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { routes } from '../src/routes/routes'
import type { Component } from 'vue'

// Test router setup
export function createTestRouter() {
  return createRouter({
    history: createWebHistory(),
    routes,
  })
}

// Default mount options for components
export function mountComponent(component: Component, options: any = {}) {
  const router = createTestRouter()
  
  return mount(component, {
    global: {
      plugins: [router],
      stubs: {
        // Stub out complex components that aren't the focus of the test
        'RouterLink': true,
        'RouterView': true,
        'icon-carbon-search': true,
        'icon-tabler-music': true,
        'icon-tabler-user': true,
        'icon-carbon-play': true,
        'icon-carbon-pause': true,
        'icon-carbon-previous': true,
        'icon-carbon-next': true,
      },
      mocks: {
        $router: {
          push: vi.fn(),
          replace: vi.fn(),
          go: vi.fn(),
          back: vi.fn(),
          forward: vi.fn(),
        },
        $route: {
          path: '/',
          params: {},
          query: {},
        },
      },
    },
    ...options,
  })
}

// Helper to wait for component to be fully mounted and updated
export async function waitForComponent(wrapper: VueWrapper<any>) {
  await wrapper.vm.$nextTick()
  await new Promise(resolve => setTimeout(resolve, 0))
}

// Mock implementation for commonly used composables - these should be used inside vi.mock calls
export const mockUseAuth = () => ({
  checkIfLoggedIn: vi.fn().mockResolvedValue(true),
  logout: vi.fn().mockResolvedValue(undefined),
  userLoginState: { value: true },
})

export const mockUseBackendFetch = () => ({
  backendFetchRequest: vi.fn().mockResolvedValue({
    json: () => Promise.resolve({}),
    ok: true,
  }),
  getCurrentUser: vi.fn().mockResolvedValue({ id: 1, username: 'test' }),
  getUsers: vi.fn().mockResolvedValue([]),
})

export const mockUseSearch = () => ({
  searchQuery: { value: '' },
  searchResults: { value: [] },
  isSearchOpen: { value: false },
  openSearch: vi.fn(),
  closeSearch: vi.fn(),
  performSearch: vi.fn(),
})

export const mockUseNavbar = () => ({
  isNavOpen: { value: false },
  toggleNav: vi.fn(),
  closeNav: vi.fn(),
})

export const mockUseSettings = () => ({
  settings: { value: { theme: 'dark', volume: 0.8 } },
  updateSettings: vi.fn(),
  loadSettings: vi.fn(),
})

// Helper to mock router push/replace
export function mockRouter() {
  return {
    push: vi.fn(),
    replace: vi.fn(),
    go: vi.fn(),
    back: vi.fn(),
    forward: vi.fn(),
    currentRoute: { value: { path: '/', params: {}, query: {} } },
  }
}