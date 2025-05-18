<script setup lang="ts">
import dayjs from 'dayjs'
import { backendFetchRequest } from '../composables/fetchFromBackend'

const randomMetadata = ref()

async function getRandomMetadata() {
  const response = await backendFetchRequest('metadata/random')
  const json = await response.json()
  randomMetadata.value = {
    name: json.title,
    artist: json.artist,
    album: json.album,
    musicbrainz_track_id: json.musicbrainz_track_id,
    musicbrainz_album_id: json.musicbrainz_album_id,
    musicbrainz_artist_id: json.musicbrainz_artist_id,
    genres: json.genre.split(';'),
    release_date: dayjs(json.release_date).format('YYYY'),
    image_url: `/api/art/albums/${json.musicbrainz_album_id}?size=xl`,
  }
}

onBeforeMount(async () => {
  await getRandomMetadata()
})
</script>

<template>
  <section v-if="randomMetadata" class="relative h-80">
    <img :src="randomMetadata.image_url" class="absolute z-0 h-80 w-full rounded-xl object-cover">
    <div class="absolute z-10 h-full w-full backdrop-blur-md">
      <div class="h-60 flex flex-row from-black to-opacity-0 bg-gradient-to-r p-10">
        <img :src="randomMetadata.image_url" class="rounded-lg">
        <div class="m-6 flex flex-col gap-5">
          <div class="text-4xl text-white font-bold">
            {{ randomMetadata.album }}
          </div>
          <div class="text-white">
            {{ randomMetadata.artist }} • {{ randomMetadata.release_date }}
          </div>
          <div class="flex flex-row gap-x-2">
            <span v-for="genre in randomMetadata.genres" :key="genre" class="rounded-full bg-zenegray-200 px-4 py-1 text-xs text-zenegray-50 font-semibold">
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
