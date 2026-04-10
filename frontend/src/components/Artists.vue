<script setup lang="ts">
import type { ArtistOrder } from '~/logic/store'
import type { SubsonicArtist } from '~/types/subsonicArtist'
import { fetchArtistList } from '~/logic/backendFetch'
import { generateSeed } from '~/logic/common'
import { artistOrder, ArtistOrders, artistSeed, artistsStore } from '~/logic/store'

const props = defineProps({
  limitRows: { type: Boolean, default: false },
  sortKey: { type: String, default: 'currentArtistOrder' },
})

const artists = ref<SubsonicArtist[]>([] as SubsonicArtist[])

let fetchType: string

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
  // typeParam != "starred"  && typeParam != "highest" &&
  // typeParam != "frequent" && typeParam != "recent"
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
  getArtists()
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
  await getArtists()
})
</script>

<template>
  <div>
    <div class="flex flex-row gap-x-4 items-center justify-between">
      <div class="flex flex-row gap-x-2 items-center">
        <h2 class="text-lg font-semibold lg:text-xl">
          Artists
        </h2>
        <Refresher @refreshed="refresh" />
      </div>
      <hr class="mx-4 border-t border-primary-400/20 flex-1" />
      <DropdownMenu
        :title="artistOrder"
        :options="Object.values(ArtistOrders)"
        align="right"
        @select="setOrder"
      />
    </div>
    <div
      v-if="artists.length > 0"
      class="auto-grid mt-4 overflow-hidden"
      :class="{ 'limit-rows': limitRows }"
    >
      <ArtistThumb
        v-for="(artist, index) in artists"
        :key="artist.musicBrainzId"
        :artist="artist"
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
