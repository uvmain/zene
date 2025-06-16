<script setup lang="ts">
import type { ArtistMetadata } from '../types'
import { useBackendFetch } from '../composables/useBackendFetch'

const router = useRouter()
const { backendFetchRequest } = useBackendFetch()

const artists = ref<ArtistMetadata[]>()

async function getArtists() {
  const response = await backendFetchRequest('artists?recent=true&limit=15')
  const json = await response.json() as ArtistMetadata[]
  artists.value = json
}

onBeforeMount(async () => {
  await getArtists()
})
</script>

<template>
  <div>
    <h2 class="py-2 text-lg font-semibold">
      Recently Updated Artists
    </h2>
    <div class="flex flex-wrap gap-6">
      <div v-for="artist in artists" :key="artist.musicbrainz_artist_id" class="flex flex-col gap-y-1 overflow-hidden transition duration-200 hover:scale-110">
        <ArtistThumb :artist="artist" class="h-40 cursor-pointer" @click="() => router.push(`/artists/${artist.musicbrainz_artist_id}`)" />
      </div>
    </div>
  </div>
</template>
