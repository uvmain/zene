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
    <RefreshHeader title="Recently Updated Artists" @refreshed="getArtists()" />
    <div class="flex flex-wrap justify-center gap-6 md:justify-start">
      <div v-for="artist in artists" :key="artist.musicbrainz_artist_id" class="flex flex-col gap-y-1 overflow-hidden transition duration-200 hover:scale-110">
        <ArtistThumb :artist="artist" class="h-40 cursor-pointer" @click="() => router.push(`/artists/${artist.musicbrainz_artist_id}`)" />
      </div>
    </div>
  </div>
</template>
