<script setup lang="ts">
import type { ArtistMetadata } from '~/types'
import { useLocalStorage } from '@vueuse/core'
import { useBackendFetch } from '~/composables/useBackendFetch'
import { useRandomSeed } from '~/composables/useRandomSeed'

const props = defineProps({
  limit: { type: Number, default: 30 },
})

const router = useRouter()
const { backendFetchRequest } = useBackendFetch()
const { refreshRandomArtistSeed, getRandomArtistSeed, randomArtistSeed } = useRandomSeed()

const artists = ref<ArtistMetadata[]>()
const showOrderOptions = ref(false)
const currentOrder = useLocalStorage<'recentlyUpdated' | 'random' | 'alphabetical'>('currentArtistsOrder', 'recentlyUpdated')

const headerTitle = computed(() => {
  switch (currentOrder.value) {
    case 'recentlyUpdated':
      return 'Artists: Recently Updated'
    case 'random':
      return 'Artists: Random'
    case 'alphabetical':
      return 'Artists: Alphabetical'
    default:
      return 'Artists'
  }
})

function setOrder(order: 'recentlyUpdated' | 'random' | 'alphabetical') {
  currentOrder.value = order
  showOrderOptions.value = false
  getArtists()
}

async function refreshArtists() {
  await refreshRandomArtistSeed()
  await getArtists()
}

async function getArtists() {
  const formData = new FormData()
  formData.append('recent', currentOrder.value === 'recentlyUpdated' ? 'true' : 'false')
  formData.append('random', randomArtistSeed.value.toString())
  formData.append('limit', props.limit.toString())
  const response = await backendFetchRequest('artists', {
    method: 'POST',
    body: formData,
  })
  const json = await response.json() as ArtistMetadata[]
  artists.value = json
}

onBeforeMount(async () => {
  getRandomArtistSeed()
  await getArtists()
})
</script>

<template>
  <div class="relative">
    <RefreshHeader :title="headerTitle" @refreshed="refreshArtists()" @title-click="showOrderOptions = !showOrderOptions" />
    <div v-if="showOrderOptions" class="absolute left-0 top-0 z-10 w-auto border-1 border-white rounded-md border-solid bg-zene-800 text-white shadow-lg">
      <div class="cursor-pointer px-4 py-2 hover:bg-zene-600" @click="setOrder('recentlyUpdated')">
        Recently Updated
      </div>
      <div class="cursor-pointer px-4 py-2 hover:bg-zene-600" @click="setOrder('random')">
        Random
      </div>
      <div class="cursor-pointer px-4 py-2 hover:bg-zene-600" @click="setOrder('alphabetical')">
        Alphabetical
      </div>
    </div>    <div class="flex flex-wrap justify-center gap-6 md:justify-start">
      <div v-for="artist in artists" :key="artist.musicbrainz_artist_id" class="flex flex-col gap-y-1 overflow-hidden transition duration-200 hover:scale-110">
        <ArtistThumb :artist="artist" class="h-40 cursor-pointer" @click="() => router.push(`/artists/${artist.musicbrainz_artist_id}`)" />
      </div>
    </div>
  </div>
</template>
