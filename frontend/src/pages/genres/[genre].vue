<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { fetchSongsByGenre } from '~/logic/backendFetch'
import { routeTracks } from '~/logic/routeTracks'

const route = useRoute('/genres/[genre]')

const tracks = ref<SubsonicSong[]>()

const genre = computed(() => `${route.params.genre}`)

function resetRefs() {
  tracks.value = []
  routeTracks.value = []
}

async function getData() {
  const genreTracks = await fetchSongsByGenre(genre.value)
  tracks.value = genreTracks
  routeTracks.value = tracks.value
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
  <Tracks v-if="tracks" :tracks="tracks" :show-album="true" />
</template>
