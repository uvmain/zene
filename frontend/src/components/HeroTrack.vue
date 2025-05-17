<script setup lang="ts">
import { backendFetchRequest } from '../composables/fetchFromBackend'

const randomTrack = ref()

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
    image_url: `/api/art/albums/${track.musicbrainz_album_id}?size=xl`,
  }))
  randomTrack.value = tracks[Math.floor(Math.random() * tracks.length)]
}

onBeforeMount(async () => {
  await getRandomTrack()
})
</script>

<template>
  <section v-if="randomTrack" class="relative h-80">
    <img :src="randomTrack?.image_url" class="absolute z-0 h-80 w-full object-cover">
    <div class="absolute z-10 h-full w-full backdrop-blur-lg">
      <div class="h-full flex flex-col justify-between">
        <div class="m-6 text-xl">
          Random Track
        </div>
        <div class="m-6">
          <h1 class="text-4xl font-bold">
            {{ randomTrack?.name }}
          </h1>
          <p class="text-gray-400">
            {{ randomTrack?.artist }} - {{ randomTrack?.album }}
          </p>
          <button class="mt-2 w-40 rounded border-none bg-sky-300 px-4 py-2 outline-none">
            Listen Now
          </button>
        </div>
      </div>
    </div>
    <img :src="randomTrack?.image_url" class="absolute size-40">
  </section>
</template>
