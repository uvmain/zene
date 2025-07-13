<script setup lang="ts">
import type { AlbumMetadata } from '../types'
import dayjs from 'dayjs'
import { useBackendFetch } from '../composables/useBackendFetch'
import { usePlaybackQueue } from '../composables/usePlaybackQueue'

const { backendFetchRequest } = useBackendFetch()
const { getRandomSeed, refreshRandomSeed } = usePlaybackQueue()

const METADATA_COUNT = 20
const isShaking = ref(false)
const albumArray = ref<AlbumMetadata[]>([])
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

async function getRandomAlbums(limit: number): Promise<AlbumMetadata[]> {
  const randomSeed = getRandomSeed()
  const response = await backendFetchRequest(`albums?random=${randomSeed}&limit=${limit}`)
  const json = await response.json()
  const albumMetadata: AlbumMetadata[] = []
  json.forEach((metadata: any) => {
    const metadataInstance = {
      artist: metadata.artist,
      album: metadata.album,
      album_artist: metadata.album_artist,
      musicbrainz_album_id: metadata.musicbrainz_album_id as string,
      musicbrainz_artist_id: metadata.musicbrainz_artist_id,
      genres: metadata.genres.split(';').filter((genre: string) => genre !== ''),
      release_date: dayjs(metadata.release_date).format('YYYY'),
      image_url: `/api/albums/${metadata.musicbrainz_album_id}/art?size=xl`,
    }
    albumMetadata.push(metadataInstance)
  })
  return albumMetadata
}

async function getNewRandomMetadata() {
  refreshRandomSeed()
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
  albumArray.value = await getRandomAlbums(METADATA_COUNT)
  index.value = 0
})
</script>

<template>
  <section v-if="albumArray.length" class="h-48 overflow-hidden rounded-lg md:h-65">
    <div
      class="h-full w-full bg-cover bg-center"
      :style="{ backgroundImage: `url(${albumArray[index].image_url})` }"
    >
      <div class="h-full w-full flex items-center justify-between backdrop-blur-md">
        <Album :album="albumArray[index]" size="xl" />
        <div class="m-3 mb-auto flex gap-2 rounded-full bg-zene-800/50 p-3 text-white md:m-6 md:p-2">
          <icon-tabler-chevron-left
            class="cursor-pointer text-2xl opacity-80 md:text-3xl active:opacity-100"
            :class="{ 'text-gray': index === 0 }"
            @click="prevIndex"
          />
          <icon-tabler-dice-6
            class="cursor-pointer text-2xl opacity-80 md:text-3xl active:opacity-100"
            :class="{ shake: isShaking }"
            @click="handleDiceClick()"
          />
          <icon-tabler-chevron-right
            class="cursor-pointer text-2xl opacity-80 md:text-3xl active:opacity-100"
            :class="{ 'text-gray': index === indexCount - 1 }"
            @click="nextIndex"
          />
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
@keyframes shake {
  0% {
    transform: rotate(0deg);
  }
  25% {
    transform: rotate(-15deg);
  }
  50% {
    transform: rotate(15deg);
  }
  75% {
    transform: rotate(-15deg);
  }
  100% {
    transform: rotate(0deg);
  }
}

.shake {
  animation: shake 0.2s ease-in-out;
}
</style>
