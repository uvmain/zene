<script setup lang="ts">
import type { SubsonicArtist } from '~/types/subsonicArtist'
import { useElementVisibility, useLocalStorage } from '@vueuse/core'
import { fetchArtistList } from '~/composables/backendFetch'
import { generateSeed } from '~/composables/logic'

const props = defineProps({
  limit: { type: Number, default: 30 },
  offset: { type: Number, default: 0 },
  scrollable: { type: Boolean, default: false },
  limitRows: { type: Boolean, default: false },
  sortKey: { type: String, default: 'currentArtistOrder' },
})

type OrderType = 'newest' | 'random' | 'alphabetical' | 'starred' | 'recent'

const loading = ref(false)
const seed = useLocalStorage<number>('artistSeed', 0)
const router = useRouter()
const currentOffset = ref<number>(0)
const canLoadMore = ref(true)
const observer = useTemplateRef<HTMLDivElement>('observer')
const observerIsVisible = useElementVisibility(observer)
const artists = ref<SubsonicArtist[]>([] as SubsonicArtist[])
const showOrderOptions = ref(false)

const _allowedOrders = ['newest', 'random', 'alphabetical', 'starred', 'recent'] as const
type AllowedOrder = typeof _allowedOrders[number]
const currentOrder = useLocalStorage<AllowedOrder>(props.sortKey, 'newest')

const sortOptions = [
  { label: 'Recently Updated', emitValue: 'newest' },
  { label: 'Recently Played', emitValue: 'recent' },
  { label: 'Random', emitValue: 'random' },
  { label: 'Alphabetical', emitValue: 'alphabetical' },
  { label: 'Starred', emitValue: 'starred' },
]

const heightStyle = computed(() => {
  if (props.limitRows) {
    const smHeight = (174 * 3) + 48
    const lgHeight = (174 * 2) + 24
    return {
      'maxHeight': `${smHeight}px`,
      '--albums-lg-max-height': `${lgHeight}px`,
    }
  }
  return {}
})

watch(observerIsVisible, (newValue) => {
  if (newValue && props.scrollable) {
    getArtists()
  }
})

const headerTitle = computed(() => {
  switch (currentOrder.value) {
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

function setOrder(order: OrderType) {
  if (currentOrder.value === order) {
    showOrderOptions.value = false
    return
  }
  currentOrder.value = order
  showOrderOptions.value = false
  resetArtistsArray()
  getArtists()
}

function resetArtistsArray() {
  canLoadMore.value = true
  loading.value = false
  currentOffset.value = props.offset
  artists.value = [] as SubsonicArtist[]
}

async function getArtists() {
  if (loading.value) {
    return
  }
  loading.value = true
  if (!canLoadMore.value) {
    return
  }
  const artistsResponse = await fetchArtistList(currentOrder.value, props.limit, currentOffset.value, seed.value)
  if (artistsResponse.length > 0) {
    currentOffset.value += artistsResponse.length
    artists.value?.push(...artistsResponse)
  }
  if (artistsResponse.length < props.limit) {
    canLoadMore.value = false
  }
  loading.value = false
  if (observerIsVisible.value) {
    getArtists()
  }
}

async function refresh() {
  if (currentOrder.value === 'random') {
    seed.value = generateSeed()
  }
  resetArtistsArray()
  getArtists()
}

onBeforeMount(async () => {
  resetArtistsArray()
  await getArtists()
})
</script>

<template>
  <div class="relative">
    <RefreshHeader :title="headerTitle" @refreshed="refresh()" @title-click="showOrderOptions = !showOrderOptions" />
    <RefreshOptions v-if="showOrderOptions" :options="sortOptions" @set-order="setOrder" />
    <div
      v-if="artists.length > 0"
      class="auto-grid-6 overflow-hidden"
      :style="heightStyle"
    >
      <ArtistThumb
        v-for="(artist, index) in artists"
        :key="artist.musicBrainzId"
        :artist="artist"
        :index="index"
        class="mx-auto transition duration-200 lg:(mx-none scale-95) hover:scale-100"
        @click="() => router.push(`/artists/${artist.musicBrainzId}`)"
      />
    </div>
    <div v-if="canLoadMore && props.scrollable" ref="observer" class="h-16">
      Loading more artists...
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
