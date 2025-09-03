<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { fetchSongsByGenre } from '~/composables/backendFetch'
import { useRouteTracks } from '~/composables/useRouteTracks'

const route = useRoute()
const { routeTracks } = useRouteTracks()

const tracks = ref<SubsonicSong[]>()

const genre = computed(() => `${route.params.genre}`)

onMounted(async () => {
  const genreTracks = await fetchSongsByGenre(genre.value, 30, 0)
  routeTracks.value = genreTracks
  tracks.value = genreTracks
})
</script>

<template>
  <Tracks v-if="tracks" :tracks="tracks" :show-album="true" />
</template>
