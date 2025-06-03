<script setup lang="ts">
import type { ArtistMetadata } from '../types'
import { backendFetchRequest } from '../composables/fetchFromBackend'

const route = useRoute()
const artist = ref<ArtistMetadata>()

const musicbrainz_artist_id = computed(() => `${route.params.musicbrainz_artist_id}`)

async function getArtists() {
  const response = await backendFetchRequest(`artists/${musicbrainz_artist_id.value}`)
  const json = await response.json() as ArtistMetadata
  artist.value = json
}

function onImageError(event: Event) {
  const target = event.target as HTMLImageElement
  target.onerror = null
  target.src = '/default-square.png'
}

onBeforeMount(async () => {
  await getArtists()
})
</script>

<template>
  <section v-if="artist" class="h-80 overflow-hidden rounded-lg">
    <div
      class="h-full w-full bg-cover bg-center"
      :style="{ backgroundImage: `url(${artist.image_url})` }"
    >
      <div class="h-full w-full flex items-center justify-center gap-6 align-middle backdrop-blur-md">
        <div class="mx-auto flex items-center justify-center gap-6 rounded-lg bg-black/40 p-4 align-middle">
          <div class="size-60">
            <img
              class="h-full w-full rounded-md object-cover"
              :src="artist.image_url"
              @error="onImageError"
            />
          </div>
          <div class="text-7xl text-gray-300 font-bold">
            {{ artist.artist }}
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
