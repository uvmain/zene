<script setup lang="ts">
import dayjs from 'dayjs'
import { backendFetchRequest } from '../composables/fetchFromBackend'

const previousMetadataResults = ref<string[]>([])
const randomMetadata = ref()

async function getRandomMetadata() {
  const response = await backendFetchRequest('metadata/random?limit=1')
  const json = await response.json()
  randomMetadata.value = {
    name: json[0].title,
    artist: json[0].artist,
    album: json[0].album,
    musicbrainz_track_id: json[0].musicbrainz_track_id,
    musicbrainz_album_id: json[0].musicbrainz_album_id as string,
    musicbrainz_artist_id: json[0].musicbrainz_artist_id,
    genres: json[0].genre.split(';').filter((genre: string) => genre !== ''),
    release_date: dayjs(json[0].release_date).format('YYYY'),
    image_url: `/api/art/albums/${json[0].musicbrainz_album_id}?size=xl`,
  }
}

async function getNewRandomMetadata() {
  while (!randomMetadata.value || previousMetadataResults.value.includes(randomMetadata.value.musicbrainz_album_id)) {
    await getRandomMetadata()
  }
  previousMetadataResults.value.push(randomMetadata.value.musicbrainz_album_id)
  while (previousMetadataResults.value.length >= 5) {
    previousMetadataResults.value.shift()
  }
}

onBeforeMount(async () => {
  await getNewRandomMetadata()
})
</script>

<template>
  <section v-if="randomMetadata" class="relative h-80 overflow-hidden rounded-xl">
    <img :src="randomMetadata.image_url" class="absolute z-0 h-80 w-full object-cover">
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
          <div v-if="randomMetadata.genres.length > 0" class="flex flex-row gap-x-2">
            <GenreBottle v-for="genre in randomMetadata.genres" :key="genre" :genre />
          </div>
          <button class="w-30 border-1 border-white rounded-full border-solid bg-zenegray-800/30 px-4 py-2 text-white font-semibold outline-none hover:bg-sky-400/50">
            Play
          </button>
        </div>
      </div>
    </div>
    <icon-tabler-dice-6 class="absolute right-3 top-3 z-10 transform-rotate-z-10 text-4xl text-zenegray-50/90" @click="getNewRandomMetadata" />
  </section>
</template>
