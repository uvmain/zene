<script setup lang="ts">
import type { HeroMetadata } from '../types'
import dayjs from 'dayjs'
import { backendFetchRequest } from '../composables/fetchFromBackend'

const METADATA_COUNT = 20
const isShaking = ref(false)
const albumArray = ref<HeroMetadata[]>([])
const index = ref(0)

const indexCount = computed(() => {
  return albumArray.value.length < METADATA_COUNT ? albumArray.value.length : METADATA_COUNT
})

function nextIndex() {
  if (index.value < indexCount.value - 1) {
    index.value += 1
  }
}

function prevIndex() {
  if (index.value > 0) {
    index.value -= 1
  }
}

async function getRandomAlbums(limit: number): Promise<HeroMetadata[]> {
  const response = await backendFetchRequest(`albums?random=true&limit=${limit}`)
  const json = await response.json()
  const heroMetadata: HeroMetadata[] = []
  json.forEach((metadata: any) => {
    const metadataInstance = {
      artist: metadata.artist,
      album: metadata.album,
      musicbrainz_track_id: metadata.musicbrainz_track_id,
      musicbrainz_album_id: metadata.musicbrainz_album_id as string,
      musicbrainz_artist_id: metadata.musicbrainz_artist_id,
      genres: metadata.genres.split(';').filter((genre: string) => genre !== ''),
      release_date: dayjs(metadata.release_date).format('YYYY'),
      image_url: `/api/art/albums/${metadata.musicbrainz_album_id}?size=xl`,
    }
    heroMetadata.push(metadataInstance)
  })
  return heroMetadata
}

async function getNewRandomMetadata() {
  albumArray.value = await getRandomAlbums(METADATA_COUNT)
  index.value = 0
}

function handleDiceClick() {
  isShaking.value = true
  setTimeout(() => {
    isShaking.value = false
  }, 200)
  getNewRandomMetadata()
}

onBeforeMount(async () => {
  await getNewRandomMetadata()
})
</script>

<template>
  <section v-if="albumArray.length" class="relative h-80 overflow-hidden rounded-xl">
    <img :src="albumArray[index].image_url" class="absolute z-0 h-80 w-full object-cover">
    <div class="absolute z-10 h-full w-full backdrop-blur-md">
      <div class="h-60 flex flex-row from-black to-opacity-0 bg-gradient-to-r p-10">
        <img :src="albumArray[index].image_url" class="rounded-lg">
        <div class="m-6 flex flex-col gap-5">
          <div class="text-4xl text-white font-bold">
            {{ albumArray[index].album }}
          </div>
          <div class="text-white">
            {{ albumArray[index].artist }} • {{ albumArray[index].release_date }}
          </div>
          <div v-if="albumArray[index].genres.length > 0" class="flex flex-row gap-x-2">
            <GenreBottle v-for="genre in albumArray[index].genres" :key="genre" :genre />
          </div>
          <button class="bg-zene-600/70 hover:bg-zene-200/70 w-30 border-1 border-white rounded-full border-solid px-4 py-2 text-xl text-white outline-none">
            Play
          </button>
        </div>
      </div>
    </div>
    <div class="bg-zene-800/50 absolute right-3 top-3 z-10 flex gap-4 rounded-full p-2 text-white">
      <icon-tabler-chevron-left class="text-4xl" :class="{ 'opacity-40': index === 0 }" @click="prevIndex" />
      <icon-tabler-dice-6 class="text-4xl" :class="{ shake: isShaking }" @click="handleDiceClick" />
      <icon-tabler-chevron-right class="text-4xl" :class="{ 'opacity-40': index === indexCount - 1 }" @click="nextIndex" />
    </div>
  </section>
</template>

<style scoped>
@keyframes shake {
  0% { transform: rotate(0deg); }
  25% { transform: rotate(-10deg); }
  50% { transform: rotate(10deg); }
  75% { transform: rotate(-10deg); }
  100% { transform: rotate(0deg); }
}

.shake {
  animation: shake 0.2s ease-in-out;
}
</style>
