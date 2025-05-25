import type { RouteRecordRaw } from 'vue-router'
import Albums from './components/routes/Albums.vue'
import Artists from './components/routes/Artists.vue'
import Home from './components/routes/Home.vue'

export const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: Home as Component,
  },
  {
    path: '/albums',
    name: 'Albums',
    component: Albums as Component,
  },
  {
    path: '/artists',
    name: 'Artists',
    component: Artists as Component,
  },
]
