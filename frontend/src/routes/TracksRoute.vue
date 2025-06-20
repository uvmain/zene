<script setup lang="ts">
import type { TrackMetadataWithImageUrl } from '../types'
import { trackWithImageUrl } from '../composables/logic'
import { useBackendFetch } from '../composables/useBackendFetch'
import { useRouteTracks } from '../composables/useRouteTracks'

const { routeTracks, clearRouteTracks } = useRouteTracks()
const { backendFetchRequest } = useBackendFetch()

const tracks = ref<TrackMetadataWithImageUrl[]>([])
const loading = ref(true)
const error = ref<string | null>(null)

onMounted(async () => {
  try {
    const response = await backendFetchRequest('tracks?random=true&limit=100')
    const json = await response.json() as TrackMetadataWithImageUrl[]

    const tracksWithImages = json.map((element) => {
      const newElement = trackWithImageUrl(element, 'sm')
      return newElement
    })
    tracks.value = tracksWithImages
    routeTracks.value = tracksWithImages
  }
  catch (err) {
    error.value = 'Failed to fetch tracks.'
    console.error(err)
  }
  finally {
    loading.value = false
  }
})

onUnmounted(() => clearRouteTracks())
</script>

<template>
  <Tracks :tracks="tracks" :show-album="true" />
</template>
