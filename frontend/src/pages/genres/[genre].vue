<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { fetchSongsByGenre } from '~/composables/backendFetch'
import { useRouteTracks } from '~/composables/useRouteTracks'

const route = useRoute()
const { routeTracks } = useRouteTracks()

const tracks = ref<SubsonicSong[]>()
const canLoadMore = ref<boolean>(true)
const limit = ref<number>(100)
const offset = ref<number>(0)

const genre = computed(() => `${route.params.genre}`)

function resetRefs() {
  tracks.value = []
  routeTracks.value = []
  offset.value = 0
  canLoadMore.value = true
}

async function getData() {
  const genreTracks = await fetchSongsByGenre(genre.value, limit.value, offset.value)
  tracks.value = tracks.value?.concat(genreTracks) ?? genreTracks
  routeTracks.value = tracks.value
  offset.value += genreTracks.length

  if (genreTracks.length < limit.value) {
    canLoadMore.value = false
  }
}

watch(genre, async () => {
  resetRefs()
  await getData()
})

onMounted(async () => {
  await getData()
})
</script>

<template>
  <Tracks v-if="tracks" :tracks="tracks" :show-album="true" :observer-enabled="canLoadMore" @observer-visible="getData" />
</template>
