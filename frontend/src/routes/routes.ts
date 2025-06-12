import type { RouteRecordRaw } from 'vue-router'
import Album from './AlbumRoute.vue'
import Albums from './AlbumsRoute.vue'
import Artist from './ArtistRoute.vue'
import Artists from './ArtistsRoute.vue'
import Genres from './GenresRoute.vue'
import Home from './HomeRoute.vue'
import Login from './LoginRoute.vue'
import Playlist from './Playlist.vue'
import Track from './TrackRoute.vue'
import Tracks from './TracksRoute.vue'

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
  {
    path: '/artists/:musicbrainz_artist_id',
    name: 'Artist',
    component: Artist as Component,
  },
  {
    path: '/genres',
    name: 'Genres',
    component: Genres as Component,
  },
  {
    path: '/tracks',
    name: 'Tracks',
    component: Tracks as Component,
  },
  {
    path: '/tracks/:musicbrainz_track_id',
    name: 'Track',
    component: Track as Component,
  },
  {
    path: '/login',
    name: 'Login',
    component: Login as Component,
  },
  {
    path: '/playlist',
    name: 'Playlist',
    component: Playlist as Component,
  },
]
