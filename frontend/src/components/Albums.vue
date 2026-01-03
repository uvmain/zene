<script setup lang="ts">
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import { useElementSize, useElementVisibility, useLocalStorage } from '@vueuse/core'
import { fetchAlbums } from '~/composables/backendFetch'
import { generateSeed } from '~/composables/logic'

const props = defineProps({
  limit: { type: Number, default: 30 },
  offset: { type: Number, default: 0 },
  scrollable: { type: Boolean, default: false },
  limitRows: { type: Boolean, default: false },
  sortKey: { type: String, default: 'currentAlbumOrder' },
})

const loading = ref(false)
const seed = useLocalStorage<number>('albumSeed', 0)
const currentOffset = ref<number>(0)
const canLoadMore = ref(true)
const observer = useTemplateRef('observer')
const firstAlbumElement = ref<HTMLElement | null>(null)
const observerIsVisible = useElementVisibility(observer)
const albums = ref<SubsonicAlbum[]>([])
const showOrderOptions = ref(false)

const allowedOrders = ['recentlyUpdated', 'random', 'alphabetical', 'releaseDate', 'recentlyPlayed'] as const
type AlbumOrder = typeof allowedOrders[number]
const currentOrder = useLocalStorage<AlbumOrder>(props.sortKey, 'recentlyUpdated')

const sortOptions = [
  { label: 'Recently Updated', emitValue: 'recentlyUpdated' },
  { label: 'Recently Played', emitValue: 'recentlyPlayed' },
  { label: 'Random', emitValue: 'random' },
  { label: 'Alphabetical', emitValue: 'alphabetical' },
  { label: 'Release Date', emitValue: 'releaseDate' },
]

const { height: firstAlbumHeight } = useElementSize(firstAlbumElement)

const heightStyle = computed(() => {
  if (props.limitRows && firstAlbumHeight.value > 0) {
    const smHeight = (firstAlbumHeight.value * 3) + 24
    const lgHeight = (firstAlbumHeight.value * 2) + 24
    return {
      'maxHeight': `${smHeight}px`,
      '--albums-lg-max-height': `${lgHeight}px`,
    }
  }
  return {}
})

let fetchType: string

watchEffect(() => {
  if (!allowedOrders.includes(currentOrder.value as AlbumOrder)) {
    currentOrder.value = 'recentlyUpdated'
  }
})

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
    case 'recentlyPlayed':
      return 'Albums: Recently Played'
    default:
      return 'Albums'
  }
})

function setOrder(order: AlbumOrder) {
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
      fetchType = 'newest'
      break
    case 'random':
      fetchType = 'random'
      break
    case 'alphabetical':
      fetchType = 'alphabeticalbyname'
      break
    case 'releaseDate':
      fetchType = 'release'
      break
    case 'recentlyPlayed':
      fetchType = 'recent'
      break
  }
  const albumsResponse = await fetchAlbums(fetchType, props.limit, currentOffset.value, seed.value)
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
  if (fetchType === 'random') {
    seed.value = generateSeed()
  }
  resetAlbumsArray()
  await getAlbums()
}

onBeforeMount(async () => {
  await getAlbums()
})
</script>

<template>
  <div class="relative">
    <RefreshHeader :title="headerTitle" @refreshed="refresh()" @title-click="showOrderOptions = !showOrderOptions" />
    <RefreshOptions v-if="showOrderOptions" :options="sortOptions" @set-order="setOrder" />
    <div
      v-if="albums.length > 0"
      class="auto-grid-6 overflow-hidden pr-1"
      :style="heightStyle"
    >
      <Album
        v-for="(album, index) in albums" :key="album.id"
        :ref="index === 0 ? (el => firstAlbumElement = el as HTMLElement) : undefined"
        :album="album"
        :index="index"
        size="sm"
        class="transition duration-200 hover:scale-100 md:scale-95"
        :show-date="false"
      />
    </div>
    <div v-if="canLoadMore && props.scrollable" ref="observer" class="h-16px">
      Loading more albums...
    </div>
  </div>
</template>

<style scoped>
@media (min-width: 1024px) {
  .auto-grid-6 {
    max-height: var(--albums-lg-max-height, none) !important;
  }
}
</style>
