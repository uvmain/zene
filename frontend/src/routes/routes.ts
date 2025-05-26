import type { RouteRecordRaw } from 'vue-router'
import Album from './AlbumRoute.vue'
import Albums from './Albums.vue'
import Artists from './Artists.vue'
import Home from './Home.vue'

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
    path: '/albums/:musicbrainz_album_id',
    name: 'Album',
    component: Album as Component,
    props: true, // Ensure route params are passed as props
  },
  {
    path: '/artists',
    name: 'Artists',
    component: Artists as Component,
  },
]
