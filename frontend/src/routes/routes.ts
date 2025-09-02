import type { RouteRecordRaw } from 'vue-router'
import { createRouter, createWebHistory } from 'vue-router'
import Admin from './AdminRoute.vue'
import Album from './AlbumRoute.vue'
import Albums from './AlbumsRoute.vue'
import Artist from './ArtistRoute.vue'
import Artists from './ArtistsRoute.vue'
import Genre from './GenreRoute.vue'
import Genres from './GenresRoute.vue'
import Home from './HomeRoute.vue'
import Login from './LoginRoute.vue'
import Playlists from './PlaylistsRoute.vue'
import Queue from './Queue.vue'
import Track from './TrackRoute.vue'
import Tracks from './TracksRoute.vue'

export const routes: RouteRecordRaw[] = [
  { path: '/', name: 'Home', component: Home as Component, meta: { requiresAuth: true } },
  { path: '/login', name: 'Login', component: Login as Component },
  { path: '/admin', name: 'Admin', component: Admin as Component, meta: { requiresAuth: true } },
  { path: '/albums', name: 'Albums', component: Albums as Component, meta: { requiresAuth: true } },
  { path: '/albums/:musicbrainz_album_id', name: 'Album', component: Album as Component, meta: { requiresAuth: true } },
  { path: '/artists', name: 'Artists', component: Artists as Component, meta: { requiresAuth: true } },
  { path: '/artists/:musicbrainz_artist_id', name: 'Artist', component: Artist as Component, meta: { requiresAuth: true } },
  { path: '/genres', name: 'Genres', component: Genres as Component, meta: { requiresAuth: true } },
  { path: '/genres/:genre', name: 'Genre', component: Genre as Component, meta: { requiresAuth: true } },
  { path: '/tracks', name: 'Tracks', component: Tracks as Component, meta: { requiresAuth: true } },
  { path: '/tracks/:musicbrainz_track_id', name: 'Track', component: Track as Component, meta: { requiresAuth: true } },
  { path: '/queue', name: 'Queue', component: Queue as Component, meta: { requiresAuth: true } },
  { path: '/playlists', name: 'Playlists', component: Playlists as Component, meta: { requiresAuth: true } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, _from, next) => {
  const authed = localStorage.getItem('apiKey') != null && localStorage.getItem('apiKey') !== ''
  if (to.meta.requiresAuth === true && !authed)
    return next({ name: 'login' })
  if (to.name === 'login' && authed)
    return next({ name: 'home' })
  next()
})

export default router
