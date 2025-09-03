<script setup lang="ts">
import type { SubsonicIndexArtist } from '~/types/subsonicArtist'
import { fetchArtists } from '~/composables/backendFetch'

const props = defineProps({
  limit: { type: Number, default: 30 },
})

const router = useRouter()

const artists = ref<SubsonicIndexArtist[]>()

const headerTitle = computed(() => {
  return 'Artists: Alphabetical'
})

async function getArtists() {
  artists.value = await fetchArtists(props.limit)
}

onBeforeMount(async () => {
  await getArtists()
})
</script>

<template>
  <div class="relative">
    <RefreshHeader :title="headerTitle" @refreshed="getArtists()" />
    <div class="flex flex-wrap justify-center gap-6 md:justify-start">
      <div v-for="artist in artists" :key="artist.musicBrainzId" class="flex flex-col gap-y-1 overflow-hidden transition duration-200 hover:scale-110">
        <ArtistThumb :artist="artist" class="h-40 cursor-pointer" @click="() => router.push(`/artists/${artist.musicBrainzId}`)" />
      </div>
    </div>
  </div>
</template>
