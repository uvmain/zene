import type { RouterScrollBehavior } from 'vue-router'
import { useLocalStorage } from '@vueuse/core'
import { ViteSSG } from 'vite-ssg'
import { routes } from 'vue-router/auto-routes'
import { createEpisodeStoreIfNotExists } from '~/stores/usePodcastStore'
import App from './App.vue'
import 'virtual:uno.css'

const apiKey = useLocalStorage('apiKey', '')

createEpisodeStoreIfNotExists()

const scrollBehavior: RouterScrollBehavior = async (to, from, savedPosition) => {
  if (to.hash) {
    return { el: to.hash }
  }

  if (savedPosition) {
    await new Promise(resolve => setTimeout(resolve, 500))
    window.scrollTo(savedPosition.left, savedPosition.top)
    return savedPosition
  }

  return { top: 0 }
}

export const createApp = ViteSSG(
  App as Component,
  {
    routes,
    scrollBehavior,
    base: import.meta.env.BASE_URL,
  },
  ({ router }) => {
    router.beforeEach(async (to) => {
      if ((apiKey.value == null || apiKey.value.length === 0) && to.path !== '/login') {
        return { path: '/login', replace: true }
      }
      return true
    })
  },
)
