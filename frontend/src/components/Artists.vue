<script setup lang="ts">
import type { ArtistOrder } from '~/logic/store'
import type { SubsonicArtist } from '~/types/subsonicArtist'
import { fetchArtistList } from '~/logic/backendFetch'
import { generateSeed } from '~/logic/common'
import { artistOrder, ArtistOrders, artistSeed, artistsStore } from '~/logic/store'

const props = defineProps({
  title: { type: String, default: 'Artists' },
  artists: { type: Array as PropType<SubsonicArtist[]>, required: false },
  orderDisabled: { type: Boolean, default: false },
  limitRows: { type: Boolean, default: false },
  sortKey: { type: String, default: 'currentArtistOrder' },
})

const artists = ref<SubsonicArtist[]>(props.artists || [] as SubsonicArtist[])

let fetchType: string

watch(() => props.artists, (newArtists) => {
  if (newArtists) {
    artists.value = newArtists
  }
})

watchEffect(() => {
  if (!Object.values(ArtistOrders).includes(artistOrder.value as ArtistOrder)) {
    artistOrder.value = ArtistOrders.RecentlyUpdated
    getArtists()
  }
})

function setOrder(order: ArtistOrder) {
  if (artistOrder.value === order) {
    return
  }
  artistOrder.value = order
  setFetchType(order)
  getArtists()
}

function setFetchType(order: ArtistOrder) {
  switch (order) {
    case ArtistOrders.RecentlyUpdated:
      fetchType = 'newest'
      break
    case ArtistOrders.Random:
      fetchType = 'random'
      break
    case ArtistOrders.Alphabetical:
      fetchType = 'alphabetical'
      break
    case ArtistOrders.Starred:
      fetchType = 'starred'
      break
    case ArtistOrders.RecentlyPlayed:
      fetchType = 'recent'
      break
  }
}

async function getArtists() {
  if (artistsStore.value.length > 0) {
    artists.value = artistsStore.value
  }
  const fetchOptions = {
    type: fetchType,
    seed: artistSeed.value,
    limit: props.limitRows ? 50 : undefined,
  }
  const fetchedArtists = await fetchArtistList(fetchOptions)
  if (fetchedArtists && JSON.stringify(fetchedArtists) !== JSON.stringify(artists.value)) {
    artists.value = fetchedArtists
    artistsStore.value = artists.value
  }
}

async function refresh() {
  if (fetchType === 'random') {
    artistSeed.value = generateSeed()
  }
  getArtists()
}

onBeforeMount(async () => {
  if (!props.orderDisabled) {
    setFetchType(artistOrder.value)
  }
  if (!props.artists) {
    await getArtists()
  }
})
</script>

<template>
  <div class="flex flex-col gap-y-2 lg:gap-y-4">
    <div class="flex flex-row gap-x-4 items-center justify-between">
      <div class="flex flex-row gap-x-2 items-center">
        <h2 class="text-lg font-semibold uppercase lg:text-xl">
          {{ title }}
        </h2>
        <Refresher @refreshed="refresh" />
      </div>
      <hr v-if="!props.orderDisabled" class="mx-2 border-t border-primary-400/20 flex-1 lg:mx-4" />
      <DropdownMenu
        v-if="!props.orderDisabled"
        :title="artistOrder"
        :options="Object.values(ArtistOrders)"
        align="right"
        @select="setOrder"
      />
    </div>
    <div
      v-if="artists.length > 0"
      class="auto-grid w-full"
      :class="{ 'limit-rows': limitRows }"
    >
      <ArtistThumb
        v-for="(artist, index) in artists"
        :key="artist.id"
        :artist="artist"
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
  @apply grid gap-4 lg:gap-6 mx-auto;
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
