<script setup lang="ts">
import type { TrackMetadataWithImageUrl } from '../types'
import { backendFetchRequest } from '../composables/fetchFromBackend'
import { setCurrentlyPlayingTrack } from '../composables/globalState'

const tracks = ref<TrackMetadataWithImageUrl[]>([])
const loading = ref(true)
const error = ref<string | null>(null)

onMounted(async () => {
  try {
    const response = await backendFetchRequest('tracks')
    const json = await response.json()
    tracks.value = json as TrackMetadataWithImageUrl[]
  }
  catch (err) {
    error.value = 'Failed to fetch tracks.'
    console.error(err)
  }
  finally {
    loading.value = false
  }
})

function playTrack(track: TrackMetadataWithImageUrl) {
  setCurrentlyPlayingTrack(track)
}
</script>

<template>
  <div class="mx-auto p-4 container">
    <h1 class="mb-4 text-2xl font-bold">
      Tracks
    </h1>
    <div v-if="loading" class="text-center">
      Loading...
    </div>
    <div v-if="error" class="text-center text-red-500">
      {{ error }}
    </div>
    <div v-if="!loading && !error && tracks.length === 0" class="text-center">
      No tracks found.
    </div>
    <table v-if="!loading && !error && tracks.length > 0" class="min-w-full bg-white">
      <thead>
        <tr>
          <th class="border-b px-4 py-2">
            Play
          </th>
          <th class="border-b px-4 py-2">
            Title
          </th>
          <th class="border-b px-4 py-2">
            Artist
          </th>
          <th class="border-b px-4 py-2">
            Album
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="track in tracks" :key="track.musicbrainz_track_id">
          <td class="border-b px-4 py-2 text-center">
            <button class="rounded bg-blue-500 px-2 py-1 text-white font-bold hover:bg-blue-700" @click="playTrack(track)">
              Play
            </button>
          </td>
          <td class="border-b px-4 py-2">
            <router-link :to="`/tracks/${track.musicbrainz_track_id}`" class="text-blue-600 hover:underline">
              {{ track.title }}
            </router-link>
          </td>
          <td class="border-b px-4 py-2">
            {{ track.artist }}
          </td>
          <td class="border-b px-4 py-2">
            {{ track.album }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
/* Add any component-specific styles here if needed */
</style>
