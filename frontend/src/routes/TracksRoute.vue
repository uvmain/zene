<script setup lang="ts">
import type { TrackMetadataWithImageUrl } from '../types'
import { useBackendFetch } from '../composables/useBackendFetch'
import { useLogic } from '../composables/useLogic'
import { usePlaycounts } from '../composables/usePlaycounts'
import { useRandomSeed } from '../composables/useRandomSeed'
import { useRouteTracks } from '../composables/useRouteTracks'

const { routeTracks, clearRouteTracks } = useRouteTracks()
const { backendFetchRequest } = useBackendFetch()
const { getRandomSeed } = useRandomSeed()
const { trackWithImageUrl } = useLogic()
const { playcount_updated_musicbrainz_track_id } = usePlaycounts()

const LIMIT = 100 as const
const offset = ref(0)
const error = ref<string | null>(null)
const canLoadMore = ref(true)

async function loadMore() {
  offset.value += LIMIT
  const randomSeed = getRandomSeed()
  const response = await backendFetchRequest(`tracks?random=${randomSeed}&limit=${LIMIT}&offset=${offset.value}`)
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

watch(playcount_updated_musicbrainz_track_id, (newTrack) => {
  routeTracks.value?.forEach((track) => {
    if (track.musicbrainz_track_id === newTrack) {
      track.user_play_count = track.user_play_count + 1
      track.global_play_count = track.global_play_count + 1
    }
  })
})

onMounted(async () => {
  try {
    const randomSeed = getRandomSeed()
    const response = await backendFetchRequest(`tracks?random=${randomSeed}&limit=${LIMIT}`)
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
