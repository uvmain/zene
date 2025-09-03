<script setup lang="ts">
import type { SubsonicArtist } from '~/types/subsonicArtist'
import { fetchArtists } from '~/composables/backendFetch'

const router = useRouter()

const artists = ref<SubsonicArtist[]>()

onBeforeMount(async () => {
  artists.value = await fetchArtists()
})
</script>

<template>
  <div>
    <h2 class="mb-2 text-lg font-semibold">
      Recently Added Artists
    </h2>
    <div class="flex flex-wrap gap-6">
      <div v-for="artist in artists" :key="artist.id" class="w-30 flex flex-col gap-y-1 overflow-hidden transition duration-200 hover:scale-110">
        <ArtistThumb :artist="artist" class="cursor-pointer" @click="() => router.push(`/artists/${artist.musicBrainzId}`)" />
      </div>
    </div>
  </div>
</template>
