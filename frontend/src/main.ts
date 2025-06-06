import type { RouterScrollBehavior } from 'vue-router'
import { ViteSSG } from 'vite-ssg'
import App from './App.vue'
import { routes } from './routes/routes'
import 'virtual:uno.css'

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
)
