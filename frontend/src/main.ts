import type { RouteLocationNormalized, RouterScrollBehavior } from 'vue-router'
import { useDark } from '@vueuse/core'
import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import { routes } from 'vue-router/auto-routes'
import App from '~/App.vue'
import { initializeAccentColour } from '~/logic/colours'
import { apiKey } from '~/stores/main'
import { createKVStoreIfNotExists } from '~/stores/keyValueIdbStore'
import { createEpisodeStoreIfNotExists } from '~/stores/podcastStore'
import { debugLog } from './logic/logger'
import 'virtual:uno.css'
import '~/styles/main.css'

useDark()
initializeAccentColour()
createEpisodeStoreIfNotExists()
createKVStoreIfNotExists()

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

let totalLength = 0
for (const item of Object.values(localStorage) as string[]) {
  const itemLength = item.length * 2 // each character is 2 bytes in UTF-16
  totalLength += itemLength
}

debugLog(`Localstorage space used: ${(totalLength / 1024).toFixed(2)} KB`)

const app = createApp(App as Component)

app.use(router)

app.mount('#app')
