import type { Component } from 'vue'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { routes } from '../src/routes/routes'

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
