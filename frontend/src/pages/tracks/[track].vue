<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { fetchSong } from '~/logic/backendFetch'
import { formatTimeFromSeconds, getCoverArtUrl, onImageError } from '~/logic/common'

const route = useRoute('/tracks/[track]')

const track = ref<SubsonicSong>()
const loading = ref(true)
const error = ref<string | null>(null)

const musicbrainzTrackId = computed(() => route.params.track as string)

const coverArtUrl = computed(() => {
  return track.value ? getCoverArtUrl(track.value?.musicBrainzId) : ''
})

onMounted(async () => {
  if (!musicbrainzTrackId.value) {
    error.value = 'Track ID is missing from the route.'
    loading.value = false
    return
  }
  try {
    const response = await fetchSong(musicbrainzTrackId.value)
    track.value = response
  }
  catch (err) {
    error.value = 'Failed to fetch track details.'
    console.error(err)
  }
  finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="mx-auto p-4 container">
    <div v-if="loading" class="text-center">
      Loading track details...
    </div>
    <div v-if="error" class="text-accent1 text-center">
      {{ error }}
    </div>
    <div v-if="track && !loading && !error" class="flex flex-col gap-8 lg:flex-row">
      <img
        v-if="coverArtUrl"
        :src="coverArtUrl"
        alt="Album Art"
        class="h-auto max-w-30vw w-full shadow-lg"
        @error="onImageError"
      >
      <div class="flex flex-col lg:w-2/3">
        <h1 class="text-3xl text-primary font-bold mb-2">
          {{ track.title }}
        </h1>
        <RouterLink
          class="text-xl text-muted mb-1 no-underline cursor-pointer hover:underline hover:underline-white"
          :to="`/artists/${track.artistId}`"
        >
          Artist: {{ track.artist }}
        </RouterLink>
        <RouterLink
          class="text-lg text-muted mb-1 no-underline cursor-pointer hover:underline hover:underline-white"
          :to="`/albums/${track.albumId}`"
        >
          Album: {{ track.album }}
        </RouterLink>
        <p class="text-muted mb-1">
          Duration: {{ formatTimeFromSeconds(track.duration) }}
        </p>
        <p v-if="track" class="text-muted mb-4">
          Released: {{ track.year }}
        </p>

        <PlayButton :track="track" />
      </div>
    </div>
    <div v-if="!track && !loading && !error" class="text-center">
      Track not found.
    </div>
  </div>
</template>
