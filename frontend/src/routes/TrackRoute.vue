<script setup lang="ts">
import type { TrackMetadataWithImageUrl } from '../types'
import { backendFetchRequest } from '../composables/fetchFromBackend'
import { formatTime, getAlbumUrl, getArtistUrl } from '../composables/logic'

const route = useRoute()
const track = ref<TrackMetadataWithImageUrl | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)

const musicbrainzTrackId = computed(() => route.params.musicbrainz_track_id as string)

onMounted(async () => {
  if (!musicbrainzTrackId.value) {
    error.value = 'Track ID is missing from the route.'
    loading.value = false
    return
  }
  try {
    const response = await backendFetchRequest(`tracks/${musicbrainzTrackId.value}`)
    const trackMetadataWithImageUrl = await response.json() as TrackMetadataWithImageUrl
    trackMetadataWithImageUrl.image_url = `/api/albums/${trackMetadataWithImageUrl.musicbrainz_album_id}/art`
    track.value = trackMetadataWithImageUrl
  }
  catch (err) {
    error.value = 'Failed to fetch track details.'
    console.error(err)
  }
  finally {
    loading.value = false
  }
})

function formatDate(dateString: string): string {
  if (!dateString)
    return 'N/A'
  const date = new Date(dateString)
  return date.toLocaleDateString(undefined, { year: 'numeric', month: 'long', day: 'numeric' })
}

function onImageError(event: Event) {
  const target = event.target as HTMLImageElement
  target.onerror = null
  target.src = '/default-square.png'
}
</script>

<template>
  <div class="mx-auto p-4 container">
    <div v-if="loading" class="text-center">
      Loading track details...
    </div>
    <div v-if="error" class="text-center text-red-500">
      {{ error }}
    </div>
    <div v-if="track && !loading && !error" class="flex flex-col gap-8 md:flex-row">
      <img v-if="track.image_url" :src="track.image_url" alt="Album Art" class="h-auto max-w-30vw w-full rounded-lg shadow-lg" @error="onImageError">
      <div class="flex flex-col md:w-2/3">
        <h1 class="mb-2 text-3xl font-bold">
          {{ track.title }}
        </h1>
        <RouterLink
          class="mb-1 cursor-pointer text-xl text-gray-300 no-underline hover:underline hover:underline-white"
          :to="getArtistUrl(track.musicbrainz_artist_id)"
        >
          Artist: {{ track.artist }}
        </RouterLink>
        <RouterLink
          class="mb-1 cursor-pointer text-lg text-gray-300 no-underline hover:underline hover:underline-white"
          :to="getAlbumUrl(track.musicbrainz_album_id)"
        >
          Album: {{ track.album }}
        </RouterLink>
        <p class="mb-1 text-gray-300">
          Duration: {{ formatTime(Number.parseFloat(track.duration)) }}
        </p>
        <p v-if="track.release_date" class="mb-4 text-gray-300">
          Released: {{ formatDate(track.release_date) }}
        </p>

        <PlayButton :track="track" />
      </div>
    </div>
    <div v-if="!track && !loading && !error" class="text-center">
      Track not found.
    </div>
  </div>
</template>
