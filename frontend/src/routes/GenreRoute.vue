<script setup lang="ts">
import type { TrackMetadataWithImageUrl } from '../types'
import { useBackendFetch } from '../composables/useBackendFetch'
import { useRouteTracks } from '../composables/useRouteTracks'

const route = useRoute()
const { routeTracks } = useRouteTracks()
const { getGenreTracks } = useBackendFetch()

const tracks = ref<TrackMetadataWithImageUrl[]>()

const genre = computed(() => `${route.params.genre}`)

onMounted(async () => {
  routeTracks.value = await getGenreTracks(genre.value, 30, true)
  tracks.value = routeTracks.value
})
</script>

<template>
  <Tracks v-if="tracks" :tracks="tracks" :show-album="true" />
</template>
