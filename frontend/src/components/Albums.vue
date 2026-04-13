<script setup lang="ts">
import type { AlbumOrder } from '~/logic/store'
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import { fetchAlbums } from '~/logic/backendFetch'
import { generateSeed } from '~/logic/common'
import { albumOrder, AlbumOrders, albumSeed, albumsStore } from '~/logic/store'
import DropdownMenu from './DropdownMenu.vue'

const props = defineProps({
  title: { type: String, default: 'Albums' },
  orderDisabled: { type: Boolean, default: false },
  albums: { type: Array as PropType<SubsonicAlbum[]>, required: false },
  limitRows: { type: Boolean, default: false },
})

const albums = ref<SubsonicAlbum[]>(props.albums || albumsStore.value)

let fetchType: string

watch(() => props.albums, (newAlbums) => {
  if (newAlbums) {
    albums.value = newAlbums
  }
})

watchEffect(() => {
  if (!Object.values(AlbumOrders).includes(albumOrder.value as AlbumOrder)) {
    albumOrder.value = AlbumOrders.RecentlyUpdated
    getAlbums()
  }
})

function setOrder(order: AlbumOrder) {
  if (albumOrder.value === order) {
    return
  }
  albumOrder.value = order
  setFetchType(order)
  getAlbums()
}

function setFetchType(order: AlbumOrder) {
  switch (order) {
    case AlbumOrders.RecentlyUpdated:
      fetchType = 'newest'
      break
    case AlbumOrders.Random:
      fetchType = 'random'
      break
    case AlbumOrders.Alphabetical:
      fetchType = 'alphabeticalbyname'
      break
    case AlbumOrders.ReleaseDate:
      fetchType = 'release'
      break
    case AlbumOrders.RecentlyPlayed:
      fetchType = 'recent'
      break
  }
}

async function getAlbums() {
  if (albumsStore.value.length > 0) {
    albums.value = albumsStore.value
  }
  const fetchOptions = {
    type: fetchType,
    seed: albumSeed.value,
    limit: props.limitRows ? 50 : undefined,
  }
  const fetchedAlbums = await fetchAlbums(fetchOptions)
  if (fetchedAlbums && JSON.stringify(fetchedAlbums) !== JSON.stringify(albums.value)) {
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
  if (!props.orderDisabled) {
    setFetchType(albumOrder.value)
  }
  if (!props.albums) {
    await getAlbums()
  }
})
</script>

<template>
  <div class="flex flex-col gap-y-2 lg:gap-y-4">
    <div class="flex flex-row gap-x-4 items-center justify-between">
      <div class="flex flex-row gap-x-2 items-center">
        <h2 class="text-lg font-semibold lg:text-xl">
          {{ title }}
        </h2>
        <Refresher @refreshed="refresh" />
      </div>
      <hr v-if="!props.orderDisabled" class="mx-2 border-t border-primary-400/20 flex-1 lg:mx-4" />
      <DropdownMenu
        v-if="!props.orderDisabled"
        :title="albumOrder"
        :options="Object.values(AlbumOrders)"
        align="right"
        @select="setOrder"
      />
    </div>
    <div
      v-if="albums.length > 0"
      class="auto-grid w-full"
      :class="{ 'limit-rows': limitRows }"
    >
      <Album
        v-for="(album, index) in albums"
        :key="album.id"
        :album="album"
        :index="index"
        class="scale-100 transition duration-200 hover:scale-105"
      />
      <!-- <div v-for="index in 30" id="push-non-full-grid-left" :key="index" aria-none class="size-full" /> -->
    </div>
    <Loading v-else />
  </div>
</template>

<style scoped>
.auto-grid {
  @apply grid gap-4 lg:gap-6 mx-auto lg:mx-0;
  @apply grid-cols-[repeat(auto-fill,minmax(min(6rem,100%),1fr))];
  @apply md:grid-cols-[repeat(auto-fill,minmax(min(8rem,100%),1fr))];
  @apply lg:grid-cols-[repeat(auto-fill,minmax(min(10rem,100%),1fr))];
}

.limit-rows {
  @apply grid-rows-[repeat(3,auto)] auto-rows-0 gap-y-0 -mb-4 lg:-mb-6;
  @apply lg:grid-rows-[repeat(2,auto)];
}

.limit-rows > * {
  @apply mb-4 lg:mb-6 overflow-hidden;
}
</style>
