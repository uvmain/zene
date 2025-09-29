<script setup lang="ts">
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import { useLocalStorage } from '@vueuse/core'
import { fetchAlbums } from '~/composables/backendFetch'

const props = defineProps({
  limit: { type: Number, default: 30 },
})

const albums = ref<SubsonicAlbum[]>()
const showOrderOptions = ref(false)
const currentOrder = useLocalStorage<'recentlyUpdated' | 'random' | 'alphabetical'>('currentAlbumOrder', 'recentlyUpdated')

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

async function getAlbums() {
  let type: string
  switch (currentOrder.value) {
    case 'recentlyUpdated':
      type = 'newest'
      break
    case 'random':
      type = 'random'
      break
    case 'alphabetical':
      type = 'alphabeticalbyname'
      break
  }
  albums.value = await fetchAlbums(type, props.limit, 0)
}

onBeforeMount(async () => {
  await getAlbums()
})
</script>

<template>
  <div class="relative">
    <RefreshHeader :title="headerTitle" @refreshed="getAlbums()" @title-click="showOrderOptions = !showOrderOptions" />
    <div v-if="showOrderOptions" class="corner-cut absolute left-0 top-0 z-10 w-auto background-2">
      <div class="cursor-pointer px-4 py-2 hover:background-3" @click="setOrder('recentlyUpdated')">
        Recently Updated
      </div>
      <div class="cursor-pointer px-4 py-2 hover:background-3" @click="setOrder('random')">
        Random
      </div>
      <div class="cursor-pointer px-4 py-2 hover:background-3" @click="setOrder('alphabetical')">
        Alphabetical
      </div>
    </div>
    <div class="flex flex-wrap justify-center gap-6 md:justify-start">
      <div v-for="album in albums" :key="album.id" class="flex flex-col gap-y-1 overflow-hidden transition duration-200 hover:scale-110">
        <Album :album="album" size="lg" />
      </div>
    </div>
  </div>
</template>
