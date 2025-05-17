<script setup lang="ts">
import { backendFetchRequest } from './composables/fetchFromBackend'

const randomTrack = ref()
const topAlbums = ref()
const topTracks = ref()

// async function getArtists() {
//   const response = await backendFetchRequest('artists')
//   const json = await response.json()
//   topArtists.value = json.map((artist: any) => ({
//     name: artist.artist,
//     plays: artist.musicbrainz_artist_id,
//   }))
// }

const genres = ref([
  'Pop',
  'Rock',
  'Hip-Hop',
  'Jazz',
  'Classical',
  'Electronic',
  'Country',
  'Reggae',
  'Blues',
])

async function getAlbums() {
  const response = await backendFetchRequest('albums')
  const json = await response.json()
  const albums = json.map((album: any) => ({
    name: album.album,
    artist: album.artist,
    musicbrainz_album_id: album.musicbrainz_album_id,
  }))
  topAlbums.value = albums.slice(0, 6)
}

async function getRandomTrack() {
  const response = await backendFetchRequest('metadata')
  const json = await response.json()
  const tracks = json.map((track: any) => ({
    name: track.title,
    artist: track.artist,
    album: track.album,
    musicbrainz_track_id: track.musicbrainz_track_id,
    musicbrainz_album_id: track.musicbrainz_album_id,
    musicbrainz_artist_id: track.musicbrainz_artist_id,
  }))
  randomTrack.value = tracks[Math.floor(Math.random() * tracks.length)]
}

async function getTopTracks() {
  const response = await backendFetchRequest('metadata')
  const json = await response.json()
  const tracks = json.map((track: any) => ({
    name: track.title,
    artist: track.artist,
    album: track.album,
    musicbrainz_track_id: track.musicbrainz_track_id,
    musicbrainz_album_id: track.musicbrainz_album_id,
    musicbrainz_artist_id: track.musicbrainz_artist_id,
  }))
  topTracks.value = tracks.slice(0, 6)
}

function getAlbumArt(musicbrainz_album_id: string, size = 'md') {
  if (!musicbrainz_album_id)
    return ''
  return `/api/art/albums/${musicbrainz_album_id}?size=${size}`
}

onBeforeMount(async () => {
  await getAlbums()
  await getRandomTrack()
  await getTopTracks()
})
</script>

<template>
  <div class="grid grid-cols-[250px_1fr] h-screen from-zenegray-900 to-zenegray-700 bg-gradient-to-b text-white">
    <!-- Sidebar -->
    <aside class="bg-zenegray-700 p-4 space-y-6">
      <div class="space-2 flex">
        <img class="h-10 w-10 rounded-full" src="/logo.png" alt="Logo" />
        <div class="text-xl font-bold">
          Zene
        </div>
      </div>
      <nav class="space-y-2">
        <div class="text-sm text-gray-400 font-semibold">
          MAIN
        </div>
        <ul class="space-y-1">
          <li><a class="block rounded px-2 py-1 hover:bg-gray-700">Home</a></li>
          <li><a class="block rounded px-2 py-1 hover:bg-gray-700">Genres</a></li>
          <li><a class="block rounded px-2 py-1 hover:bg-gray-700">Albums</a></li>
          <li><a class="block rounded px-2 py-1 hover:bg-gray-700">Artists</a></li>
          <li><a class="block rounded px-2 py-1 hover:bg-gray-700">Radio</a></li>
        </ul>
        <div class="mt-4 text-sm text-gray-400 font-semibold">
          LIBRARY
        </div>
        <ul class="space-y-1">
          <li><a class="block rounded px-2 py-1 hover:bg-gray-700">Recently Added</a></li>
          <li><a class="block rounded px-2 py-1 hover:bg-gray-700">Recently Played</a></li>
          <li><a class="block rounded px-2 py-1 hover:bg-gray-700">Albums</a></li>
          <li><a class="block rounded px-2 py-1 hover:bg-gray-700">Random</a></li>
        </ul>
        <div class="mt-4 text-sm text-gray-400 font-semibold">
          PLAYLISTS
        </div>
        <ul class="space-y-1">
          <li><a class="block rounded px-2 py-1 hover:bg-gray-700">Create New</a></li>
          <li><a class="block rounded px-2 py-1 hover:bg-gray-700">Running club</a></li>
          <li><a class="block rounded px-2 py-1 hover:bg-gray-700">Head bangers</a></li>
          <li><a class="block rounded px-2 py-1 hover:bg-gray-700">Tomato</a></li>
        </ul>
      </nav>
    </aside>

    <!-- Main Content -->
    <main class="overflow-y-auto p-6 space-y-6">
      <Header />

      <section v-if="randomTrack">
        <img :src="getAlbumArt(randomTrack?.musicbrainz_album_id, 'xl')" class="h-80 w-full object-cover">
        <h1 class="text-4xl font-bold">
          {{ randomTrack?.name }}
        </h1>
        <p class="text-gray-400">
          {{ randomTrack?.artist }} - {{ randomTrack?.album }}
        </p>
        <button class="mt-2 rounded border-none bg-sky-300 px-4 py-2 outline-none">
          Listen Now
        </button>
      </section>

      <section class="grid grid-cols-2 gap-6">
        <div>
          <h2 class="mb-2 text-lg font-semibold">
            Top Albums
          </h2>
          <div class="flex gap-6 overflow-x-auto">
            <div v-for="album in topAlbums" :key="album.album" class="w-20 text-center">
              <img class="h-20 w-20 rounded-lg bg-gray-700" :src="getAlbumArt(album.musicbrainz_album_id)" alt="Album Cover" />
              <div class="mt-1 text-sm">
                {{ album.name }}
              </div>
              <div class="text-xs text-gray-400">
                {{ album.artist }}
              </div>
            </div>
          </div>
        </div>

        <!-- Player -->
        <div class="rounded bg-gray-800 p-4">
          <div class="text-center">
            <div class="mx-auto mb-2 h-24 w-24 bg-blue-500"></div>
            <div class="font-semibold">
              Track name
            </div>
            <div class="text-sm text-gray-400">
              Artist · Album
            </div>
          </div>
          <div class="mt-4 flex justify-between text-xs">
            <span>2:34</span><span>3:21</span>
          </div>
          <div class="my-1 h-1 rounded bg-gray-600"></div>
          <div class="mt-2 flex justify-center space-x-4">
            <button>
              <icon-tabler-player-skip-back-filled />
            </button>
            <button>
              <icon-tabler-player-play-filled class="text-3xl" />
            </button>
            <button>
              <icon-tabler-player-skip-forward-filled />
            </button>
          </div>
          <button class="mt-2 text-xs underline">
            LYRICS
          </button>
        </div>
      </section>

      <section class="grid grid-cols-3 gap-6">
        <!-- Genres -->
        <div>
          <h2 class="mb-2 text-lg font-semibold">
            Genres
          </h2>
          <div class="flex flex-wrap gap-2">
            <span v-for="genre in genres" :key="genre" class="rounded bg-gray-700 px-3 py-1 text-sm">{{ genre }}</span>
          </div>
        </div>

        <!-- Top Tracks -->
        <div>
          <h2 class="mb-2 text-lg font-semibold">
            Top Tracks
          </h2>
          <ul class="space-y-2">
            <li v-for="(track, i) in topTracks" :key="track.title" class="flex justify-between text-sm">
              <span>{{ i + 1 }}. {{ track.name }}</span>
              <span>{{ track.artist }}</span>
            </li>
          </ul>
        </div>
      </section>
    </main>
  </div>
</template>

<style>
html, body, #app {
  margin: 0;
  padding: 0;
  border: 0;
  font-family: 'Montserrat', sans-serif;
  min-height: 100%;
  @apply standard;
}
</style>
