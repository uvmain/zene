<script setup lang="ts">
import { useSessionStorage } from '@vueuse/core'
import dayjs from 'dayjs'
import { useBackendFetch } from '../composables/useBackendFetch'
import { useRandomSeed } from '../composables/useRandomSeed'

const props = defineProps({
  limit: { type: Number, default: 30 },
})

const { backendFetchRequest } = useBackendFetch()
const { refreshRandomAlbumSeed, getRandomAlbumSeed, randomAlbumSeed } = useRandomSeed()

const recentlyAddedAlbums = ref()
const showOrderOptions = ref(false)
const currentOrder = useSessionStorage<'recentlyUpdated' | 'random' | 'alphabetical'>('currentAlbumOrder', 'recentlyUpdated')

const headerTitle = computed(() => {
  switch (currentOrder.value) {
    case 'recentlyUpdated':
      return 'Albums: Recently Updated'
    case 'random':
      return 'Albums: Random'
    case 'alphabetical':
      return 'Albums: Alphabetical'
    default:
      return 'Albums'
  }
})

function setOrder(order: 'recentlyUpdated' | 'random' | 'alphabetical') {
  currentOrder.value = order
  showOrderOptions.value = false
  getAlbums()
}

async function refreshAlbums() {
  await refreshRandomAlbumSeed()
  await getAlbums()
}

async function getAlbums() {
  const params = new URLSearchParams()
  params.set('recent', currentOrder.value === 'recentlyUpdated' ? 'true' : 'false')
  params.set('random', currentOrder.value === 'random' ? `${randomAlbumSeed.value}` : '')
  params.set('limit', `${props.limit}`)
  const response = await backendFetchRequest(`albums?${params.toString()}`)
  const json = await response.json()
  const albums = json.map((album: any) => ({
    album: album.album,
    artist: album.artist,
    album_artist: album.album_artist ?? album.artist,
    musicbrainz_album_id: album.musicbrainz_album_id,
    release_date: dayjs(album.release_date).format('YYYY'),
    image_url: `/api/albums/${album.musicbrainz_album_id}/art?size=lg`,
  }))
  recentlyAddedAlbums.value = albums
}

onBeforeMount(async () => {
  getRandomAlbumSeed()
  await getAlbums()
})
</script>

<template>
  <div class="relative">
    <RefreshHeader :title="headerTitle" @refreshed="refreshAlbums()" @title-click="showOrderOptions = !showOrderOptions" />
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
    </div>
    <div class="flex flex-wrap justify-center gap-6 md:justify-start">
      <div v-for="album in recentlyAddedAlbums" :key="album.album" class="flex flex-col gap-y-1 overflow-hidden transition duration-200 hover:scale-110">
        <Album :album="album" size="lg" />
      </div>
    </div>
  </div>
</template>
