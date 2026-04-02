<script setup lang="ts">
import type { ArtistOrder } from '~/logic/store'
import type { SubsonicArtist } from '~/types/subsonicArtist'
import { fetchArtistList } from '~/logic/backendFetch'
import { generateSeed } from '~/logic/common'
import { artistOrder, artistSeed } from '~/logic/store'

const props = defineProps({
  limitRows: { type: Boolean, default: false },
  sortKey: { type: String, default: 'currentArtistOrder' },
})

const router = useRouter()
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
  const fetchOptions = {
    type: artistOrder.value,
    offset: artists.value.length,
    seed: artistSeed.value,
    limit: props.limitRows ? 50 : undefined,
  }
  artists.value = await fetchArtistList(fetchOptions)
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
        class="lg:mx-none mx-auto transition duration-200 hover:scale-100 lg:(scale-95)"
        @click="() => router.push(`/artists/${artist.musicBrainzId}`)"
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
