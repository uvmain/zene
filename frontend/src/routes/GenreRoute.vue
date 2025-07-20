<script setup lang="ts">
import type { TrackMetadataWithImageUrl } from '../types'
import { useBackendFetch } from '../composables/useBackendFetch'
import { useLogic } from '../composables/useLogic'
import { useRouteTracks } from '../composables/useRouteTracks'

const route = useRoute()
const { routeTracks } = useRouteTracks()
const { getGenreTracks } = useBackendFetch()
const { trackWithImageUrl } = useLogic()

const tracks = ref<TrackMetadataWithImageUrl[]>()

const genre = computed(() => `${route.params.genre}`)

onMounted(async () => {
  const genreTracks = await getGenreTracks(genre.value, 30, true)
  const tracksWithImages = genreTracks.map((element) => {
    const newElement = trackWithImageUrl(element)
    return newElement
  })
  routeTracks.value = tracksWithImages
  tracks.value = tracksWithImages
})
</script>

<template>
  <Tracks v-if="tracks" :tracks="tracks" :show-album="true" />
</template>
