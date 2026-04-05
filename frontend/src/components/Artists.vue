<script setup lang="ts">
import type { ArtistOrder } from '~/logic/store'
import type { SubsonicArtist } from '~/types/subsonicArtist'
import { fetchArtistList } from '~/logic/backendFetch'
import { generateSeed } from '~/logic/common'
import { artistOrder, artistSeed, artistsStore } from '~/logic/store'

const props = defineProps({
  limitRows: { type: Boolean, default: false },
  sortKey: { type: String, default: 'currentArtistOrder' },
})

const artists = ref<SubsonicArtist[]>([] as SubsonicArtist[])
const showOrderOptions = ref(false)

const sortOptions = [
  { label: 'Recently Updated', emitValue: 'newest' },
  { label: 'Recently Played', emitValue: 'recent' },
  { label: 'Random', emitValue: 'random' },
  { label: 'Alphabetical', emitValue: 'alphabetical' },
  { label: 'Starred', emitValue: 'starred' },
]

const headerTitle = computed(() => {
  switch (artistOrder.value) {
    case 'newest':
      return 'Artists: Recently Updated'
    case 'random':
      return 'Artists: Random'
    case 'alphabetical':
      return 'Artists: Alphabetical'
    case 'starred':
      return 'Artists: Starred'
    case 'recent':
      return 'Artists: Recently Played'
    default:
      return 'Artists'
  }
})

function setOrder(order: ArtistOrder) {
  if (artistOrder.value === order) {
    showOrderOptions.value = false
    return
  }
  artistOrder.value = order
  showOrderOptions.value = false
  getArtists()
}

async function getArtists() {
  if (artistsStore.value.length > 0) {
    artists.value = artistsStore.value
  }
  const fetchOptions = {
    type: artistOrder.value,
    seed: artistSeed.value,
    limit: props.limitRows ? 50 : undefined,
  }
  const fetchedArtists = await fetchArtistList(fetchOptions)
  if (fetchedArtists && fetchedArtists.length > 0 && JSON.stringify(fetchedArtists) !== JSON.stringify(artists.value)) {
    artists.value = fetchedArtists
    artistsStore.value = artists.value
  }
}

async function refresh() {
  if (artistOrder.value === 'random') {
    artistSeed.value = generateSeed()
  }
  getArtists()
}

onBeforeMount(async () => {
  await getArtists()
})
</script>

<template>
  <div class="relative">
    <RefreshHeader :title="headerTitle" @refreshed="refresh()" @title-click="showOrderOptions = !showOrderOptions" />
    <RefreshOptions v-if="showOrderOptions" :options="sortOptions" @set-order="setOrder" />
    <div
      v-if="artists.length > 0"
      class="auto-grid overflow-hidden"
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
