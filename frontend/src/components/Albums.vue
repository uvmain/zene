<script setup lang="ts">
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import { useElementVisibility, useLocalStorage } from '@vueuse/core'
import { fetchAlbums } from '~/composables/backendFetch'
import { generateSeed } from '~/composables/logic'

const props = defineProps({
  limit: { type: Number, default: 30 },
  offset: { type: Number, default: 0 },
  scrollable: { type: Boolean, default: false },
})

const loading = ref(false)
const seed = useLocalStorage<number>('albumSeed', 0)
const currentOffset = ref<number>(0)
const canLoadMore = ref(true)
const observer = useTemplateRef<HTMLDivElement>('observer')
const observerIsVisible = useElementVisibility(observer)

let type: string

const albums = ref<SubsonicAlbum[]>([] as SubsonicAlbum[])
const showOrderOptions = ref(false)
const currentOrder = useLocalStorage<'recentlyUpdated' | 'random' | 'alphabetical' | 'releaseDate'>('currentAlbumOrder', 'recentlyUpdated')

watch(observerIsVisible, (newValue) => {
  if (newValue && props.scrollable) {
    getAlbums()
  }
})

const headerTitle = computed(() => {
  switch (currentOrder.value) {
    case 'recentlyUpdated':
      return 'Albums: Recently Updated'
    case 'random':
      return 'Albums: Random'
    case 'alphabetical':
      return 'Albums: Alphabetical'
    case 'releaseDate':
      return 'Albums: Release Date'
    default:
      return 'Albums'
  }
})

function setOrder(order: 'recentlyUpdated' | 'random' | 'alphabetical' | 'releaseDate') {
  if (currentOrder.value === order) {
    showOrderOptions.value = false
    return
  }
  currentOrder.value = order
  showOrderOptions.value = false
  resetAlbumsArray()
  getAlbums()
}

function resetAlbumsArray() {
  canLoadMore.value = true
  loading.value = false
  currentOffset.value = props.offset
  albums.value = [] as SubsonicAlbum[]
}

async function getAlbums() {
  if (loading.value) {
    return
  }
  loading.value = true
  if (!canLoadMore.value) {
    return
  }
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
    case 'releaseDate':
      type = 'release'
      break
  }
  const albumsResponse = await fetchAlbums(type, props.limit, currentOffset.value, seed.value)
  if (albumsResponse.length > 0) {
    currentOffset.value += albumsResponse.length
    albums.value?.push(...albumsResponse)
  }
  if (albumsResponse.length < props.limit) {
    canLoadMore.value = false
  }
  loading.value = false
  if (observerIsVisible.value) {
    getAlbums()
  }
}

async function refresh() {
  if (type === 'random') {
    seed.value = generateSeed()
  }
  resetAlbumsArray()
  await getAlbums()
}

onBeforeMount(async () => {
  resetAlbumsArray()
  await getAlbums()
})
</script>

<template>
  <div class="relative">
    <RefreshHeader :title="headerTitle" @refreshed="refresh()" @title-click="showOrderOptions = !showOrderOptions" />
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
      <div class="cursor-pointer px-4 py-2 hover:background-3" @click="setOrder('releaseDate')">
        Release Date
      </div>
    </div>
    <div v-if="albums.length > 0" class="flex flex-wrap justify-center gap-6 md:justify-start">
      <div v-for="(album, index) in albums" :key="album.id" class="transition duration-200 hover:scale-110">
        <Album :album="album" :index="index" size="sm" />
      </div>
    </div>
  </div>
  <div v-if="canLoadMore && props.scrollable" ref="observer" class="h-16">
    Loading more albums...
  </div>
</template>
