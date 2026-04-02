<script setup lang="ts">
import type { AlbumOrder } from '~/logic/store'
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import { fetchAlbums } from '~/logic/backendFetch'
import { generateSeed } from '~/logic/common'
import { albumOrder, albumOrders, albumSeed, albumsStore } from '~/logic/store'

const props = defineProps({
  limitRows: { type: Boolean, default: false },
  sortKey: { type: String, default: 'currentAlbumOrder' },
})

const albums = ref<SubsonicAlbum[]>([])
const showOrderOptions = ref(false)

const sortOptions = [
  { label: 'Recently Updated', emitValue: 'recentlyUpdated' },
  { label: 'Recently Played', emitValue: 'recentlyPlayed' },
  { label: 'Random', emitValue: 'random' },
  { label: 'Alphabetical', emitValue: 'alphabetical' },
  { label: 'Release Date', emitValue: 'releaseDate' },
]

let fetchType: string

watchEffect(() => {
  if (!albumOrders.includes(albumOrder.value as AlbumOrder)) {
    albumOrder.value = 'recentlyUpdated'
    getAlbums()
  }
})

const headerTitle = computed(() => {
  switch (albumOrder.value) {
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
  if (albumOrder.value === order) {
    showOrderOptions.value = false
    return
  }
  albumOrder.value = order
  showOrderOptions.value = false
  getAlbums()
}

async function getAlbums() {
  if (albumsStore.value.length > 0) {
    albums.value = albumsStore.value
  }
  switch (albumOrder.value) {
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

  const fetchOptions = {
    type: fetchType,
    seed: albumSeed.value,
    limit: props.limitRows ? 50 : undefined,
  }
  albums.value = await fetchAlbums(fetchOptions)
  albumsStore.value = albums.value
}

async function refresh() {
  if (fetchType === 'random') {
    albumSeed.value = generateSeed()
  }
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
      class="auto-grid overflow-hidden"
      :class="{ 'limit-rows': limitRows }"
    >
      <Album
        v-for="(album, index) in albums" :key="album.id"
        :album="album"
        :index="index"
        size="sm"
        class="transition duration-200 hover:scale-100 lg:scale-95"
        :show-date="false"
      />
    </div>
    <Loading v-else />
  </div>
</template>

<style scoped>
.auto-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(min(150px, 100%), 1fr));
  gap: calc(var(--spacing) * 4);
}

.auto-grid.limit-rows {
  grid-template-rows: repeat(3, auto);
  grid-auto-rows: 0;
  row-gap: 0;
  overflow: hidden;
}

.auto-grid.limit-rows > * {
  margin-bottom: calc(var(--spacing) * 4);
  overflow: hidden;
}
</style>
