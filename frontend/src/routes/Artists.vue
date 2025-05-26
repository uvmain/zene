<script setup lang="ts">
import type { ArtistMetadata } from '../types'
import { backendFetchRequest } from '../composables/fetchFromBackend'

const artists = ref()

async function getAlbums() {
  const response = await backendFetchRequest('artists')
  const json = await response.json() as ArtistMetadata[]
  artists.value = json
}

onBeforeMount(async () => {
  await getAlbums()
})
</script>

<template>
  <div>
    <h2 class="mb-2 text-lg font-semibold">
      Recently Added Artists
    </h2>
    <div class="flex flex-wrap gap-6">
      <div v-for="artist in artists" :key="artist.musicbrainz_artist_id" class="w-30 flex flex-col gap-y-1 overflow-hidden transition duration-200 hover:scale-110">
        <ArtistThumb :artist="artist" />
      </div>
    </div>
  </div>
</template>
