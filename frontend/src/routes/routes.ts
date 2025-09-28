import type { RouteRecordRaw } from 'vue-router'
import Admin from './AdminRoute.vue'
import Album from './AlbumRoute.vue'
import Albums from './AlbumsRoute.vue'
import Artist from './ArtistRoute.vue'
import Artists from './ArtistsRoute.vue'
import Genre from './GenreRoute.vue'
import Genres from './GenresRoute.vue'
import Home from './HomeRoute.vue'
import Login from './LoginRoute.vue'
import Playlist from './PlaylistRoute.vue'
import Playlists from './PlaylistsRoute.vue'
import Podcast from './Podcast.vue'
import Podcasts from './Podcasts.vue'
import Queue from './Queue.vue'
import Radio from './Radio.vue'
import Track from './TrackRoute.vue'
import Tracks from './TracksRoute.vue'

export const routes: RouteRecordRaw[] = [
  { path: '/', name: 'Home', component: Home as Component },
  { path: '/login', name: 'Login', component: Login as Component },
  { path: '/admin', name: 'Admin', component: Admin as Component },
  { path: '/albums', name: 'Albums', component: Albums as Component },
  { path: '/albums/:musicbrainz_album_id', name: 'Album', component: Album as Component },
  { path: '/artists', name: 'Artists', component: Artists as Component },
  { path: '/artists/:musicbrainz_artist_id', name: 'Artist', component: Artist as Component },
  { path: '/genres', name: 'Genres', component: Genres as Component },
  { path: '/genres/:genre', name: 'Genre', component: Genre as Component },
  { path: '/tracks', name: 'Tracks', component: Tracks as Component },
  { path: '/tracks/:musicbrainz_track_id', name: 'Track', component: Track as Component },
  { path: '/queue', name: 'Queue', component: Queue as Component },
  { path: '/playlists', name: 'Playlists', component: Playlists as Component },
  { path: '/playlists/:playlist_id', name: 'Playlist', component: Playlist as Component },
  { path: '/radio', name: 'Radio', component: Radio as Component },
  { path: '/podcasts', name: 'Podcasts', component: Podcasts as Component },
  { path: '/podcasts/:podcast_id', name: 'Podcast', component: Podcast as Component },
]
