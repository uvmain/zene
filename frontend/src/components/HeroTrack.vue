<script setup lang="ts">
import dayjs from 'dayjs'
import { backendFetchRequest } from '../composables/fetchFromBackend'

const randomAlbum = ref()

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
    genre: track.genre,
    release_date: dayjs(track.release_date).format('YYYY'),
    image_url: `/api/art/albums/${track.musicbrainz_album_id}?size=xl`,
  }))
  randomAlbum.value = tracks[Math.floor(Math.random() * tracks.length)]
}

const genres = computed(() => {
  return randomAlbum.value.genre.split(';')
})

onBeforeMount(async () => {
  await getRandomTrack()
})
</script>

<template>
  <section v-if="randomAlbum" class="relative h-80">
    <img :src="randomAlbum?.image_url" class="absolute z-0 h-80 w-full rounded-xl object-cover">
    <div class="absolute z-10 h-full w-full backdrop-blur-md">
      <div class="h-60 flex flex-row from-black to-opacity-0 bg-gradient-to-r p-10">
        <img :src="randomAlbum?.image_url" class="rounded-lg">
        <div class="m-6 flex flex-col gap-5">
          <div class="text-4xl text-white font-bold">
            {{ randomAlbum?.album }}
          </div>
          <div class="text-white">
            {{ randomAlbum?.artist }} • {{ randomAlbum.release_date }}
          </div>
          <div class="flex flex-row gap-x-2">
            <span v-for="genre in genres" :key="genre" class="bg-zenegray-200 rounded-full px-4 py-1 text-xs text-zenegray-50 font-semibold">
              {{ genre }}
            </span>
          </div>
          <button class="w-30 border-1 border-white rounded-full border-solid bg-zenegray-800/30 px-4 py-2 text-white font-semibold outline-none hover:bg-sky-400/50">
            Play
          </button>
        </div>
      </div>
    </div>
  </section>
</template>
