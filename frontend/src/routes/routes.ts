import type { RouteRecordRaw } from 'vue-router'
import Album from './AlbumRoute.vue'
import Albums from './AlbumsRoute.vue'
import Artists from './ArtistsRoute.vue'
import Home from './HomeRoute.vue'

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
  },
  {
    path: '/artists',
    name: 'Artists',
    component: Artists as Component,
  },
]
