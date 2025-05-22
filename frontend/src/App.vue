<script setup lang="ts">
import RecentlyAddedAlbums from './components/RecentlyAddedAlbums.vue'

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

// async function getTopTracks() {
//   const response = await backendFetchRequest('metadata')
//   const json = await response.json()
//   const tracks = json.map((track: any) => ({
//     name: track.title,
//     artist: track.artist,
//     album: track.album,
//     musicbrainz_track_id: track.musicbrainz_track_id,
//     musicbrainz_album_id: track.musicbrainz_album_id,
//     musicbrainz_artist_id: track.musicbrainz_artist_id,
//   }))
//   topTracks.value = tracks.slice(0, 6)
// }
</script>

<template>
  <div class="grid grid-cols-[250px_1fr] h-screen from-zene-800 to-zene-600 bg-gradient-to-b text-white">
    <Navbar />

    <main class="overflow-y-auto p-6 space-y-6">
      <Header />

      <HeroAlbum />

      <section class="grid grid-cols-2 gap-6">
        <div>
          <RecentlyAddedAlbums />
          <section class="grid grid-cols-3 gap-6">
            <!-- Genres -->
            <div>
              <h2 class="mb-2 text-lg font-semibold">
                Genres
              </h2>
              <div class="flex flex-wrap gap-2">
                <GenreBottle v-for="genre in genres" :key="genre" :genre />
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
