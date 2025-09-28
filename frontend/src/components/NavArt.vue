<script setup lang="ts">
import { getCoverArtUrl, onImageError } from '~/composables/logic'
import { usePlaybackQueue } from '~/composables/usePlaybackQueue'

const { currentlyPlayingTrack } = usePlaybackQueue()
const router = useRouter()

const coverArtUrl = computed(() => {
  return currentlyPlayingTrack.value ? getCoverArtUrl(currentlyPlayingTrack.value?.albumId) : '/default-square.png'
})
</script>

<template>
  <div v-if="currentlyPlayingTrack" class="flex flex-col gap-2">
    <RouterLink
      class="cursor-pointer text-lg text-white/80 no-underline hover:underline hover:underline-white"
      :to="`/tracks/${currentlyPlayingTrack?.id}`"
    >
      {{ currentlyPlayingTrack?.title }}
    </RouterLink>
    <RouterLink
      class="cursor-pointer text-sm text-white/80 no-underline hover:underline hover:underline-white"
      :to="`/artists/${currentlyPlayingTrack.artistId}`"
    >
      {{ currentlyPlayingTrack?.artist }}
    </RouterLink>
    <RouterLink
      class="cursor-pointer text-sm text-white/80 no-underline hover:underline hover:underline-white"
      :to="`/albums/${currentlyPlayingTrack.albumId}`"
    >
      {{ currentlyPlayingTrack?.album }}
    </RouterLink>
    <img :src="coverArtUrl" class="w-full cursor-pointer object-cover" @error="onImageError" @click="() => router.push(`/albums/${currentlyPlayingTrack?.albumId}`)">
  </div>
</template>
