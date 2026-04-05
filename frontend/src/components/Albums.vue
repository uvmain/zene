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

const albums = ref<SubsonicAlbum[]>(albumsStore.value)
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
  const fetchedAlbums = await fetchAlbums(fetchOptions)
  if (fetchedAlbums && fetchedAlbums.length > 0 && JSON.stringify(fetchedAlbums) !== JSON.stringify(albums.value)) {
    albums.value = fetchedAlbums
    albumsStore.value = fetchedAlbums
  }
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
        v-for="(album, index) in albums"
        :key="album.id"
        :album="album"
        :index="index"
        class="transition duration-200 hover:scale-100 lg:(scale-95)"
      />
    </div>
    <Loading v-else />
  </div>
</template>

<style scoped>
.auto-grid {
  @apply grid gap-1rem mx-auto lg:mx-0;
  @apply grid-cols-[repeat(auto-fit,minmax(min(6rem,100%),1fr))];
  @apply md:grid-cols-[repeat(auto-fit,minmax(min(8rem,100%),1fr))];
  @apply lg:grid-cols-[repeat(auto-fit,minmax(min(10rem,100%),1fr))];
}

.limit-rows {
  @apply grid-rows-[repeat(3,auto)] auto-rows-0 gap-y-0 -mb-1rem;
  @apply lg:grid-rows-[repeat(2,auto)] auto-rows-0 gap-y-0 -mb-1rem;
}

.limit-rows > * {
  @apply mb-1rem overflow-hidden;
}
</style>
