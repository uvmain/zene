<script setup lang="ts">
import type { TrackMetadataWithImageUrl } from '~/types'
import { useBackendFetch } from '~/composables/useBackendFetch'
import { useLogic } from '~/composables/useLogic'
import { useRandomSeed } from '~/composables/useRandomSeed'
import { useRouteTracks } from '~/composables/useRouteTracks'

const { routeTracks, clearRouteTracks } = useRouteTracks()
const { backendFetchRequest } = useBackendFetch()
const { getRandomSeed } = useRandomSeed()
const { trackWithImageUrl } = useLogic()

const LIMIT = 100 as const
const offset = ref(0)
const error = ref<string | null>(null)
const canLoadMore = ref(true)

async function loadMore() {
  offset.value += LIMIT
  const randomSeed = getRandomSeed()
  const formData = new FormData()
  formData.append('random', randomSeed.toString())
  formData.append('limit', LIMIT.toString())
  formData.append('offset', offset.value.toString())
  const response = await backendFetchRequest('albums', {
    method: 'POST',
    body: formData,
  })
  const json = await response.json() as TrackMetadataWithImageUrl[]
  if (json.length === 0) {
    canLoadMore.value = false
  }
  else {
    const tracksWithImages = json.map((track) => {
      const newTrack = trackWithImageUrl(track)
      return newTrack
    })
    routeTracks.value.push(...tracksWithImages)
  }
}

onMounted(async () => {
  try {
    const randomSeed = getRandomSeed()
    const formData = new FormData()
    formData.append('random', randomSeed.toString())
    formData.append('limit', LIMIT.toString())
    const response = await backendFetchRequest('tracks', {
      method: 'POST',
      body: formData,
    })
    const json = await response.json() as TrackMetadataWithImageUrl[]

    const tracksWithImages = json.map((track) => {
      const newTrack = trackWithImageUrl(track)
      return newTrack
    })
    routeTracks.value = tracksWithImages
  }
  catch (err) {
    error.value = 'Failed to fetch tracks.'
    console.error(err)
  }
})

onUnmounted(() => clearRouteTracks())
</script>

<template>
  <Tracks :tracks="routeTracks" :show-album="true" :can-load-more="canLoadMore" @load-more="loadMore()" />
</template>
