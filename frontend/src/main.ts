import type { RouteLocationNormalized, RouterScrollBehavior } from 'vue-router'
import { useDark } from '@vueuse/core'
import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import { routes } from 'vue-router/auto-routes'
import App from '~/App.vue'
import { apiKey } from '~/logic/store'
import { createEpisodeStoreIfNotExists } from '~/stores/usePodcastStore'

import 'virtual:uno.css'
import '~/styles/main.css'

useDark()
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

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
  scrollBehavior,
})

router.beforeEach((to: RouteLocationNormalized) => {
  if ((apiKey.value == null || apiKey.value.length === 0) && to.path !== '/login') {
    return { path: '/login', replace: true }
  }
})

createApp(App as Component).use(router).mount('#app')
