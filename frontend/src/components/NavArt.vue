<script setup lang="ts">
import { getCoverArtUrl, onImageError } from '~/composables/logic'
import { usePlaybackQueue } from '~/composables/usePlaybackQueue'

const { currentlyPlayingTrack } = usePlaybackQueue()
const router = useRouter()

const coverArtUrl = computed(() => {
  return currentlyPlayingTrack.value ? getCoverArtUrl(currentlyPlayingTrack.value?.id) : '/default-square.png'
})
</script>

<template>
  <div v-if="currentlyPlayingTrack" class="flex flex-col gap-2">
    <RouterLink
      class="cursor-pointer text-lg text-white/80 no-underline hover:underline hover:underline-white"
      :to="getCoverArtUrl(currentlyPlayingTrack.id)"
    >
      {{ currentlyPlayingTrack?.title }}
    </RouterLink>
    <RouterLink
      class="cursor-pointer text-sm text-white/80 no-underline hover:underline hover:underline-white"
      :to="getCoverArtUrl(currentlyPlayingTrack.artistId)"
    >
      {{ currentlyPlayingTrack?.artist }}
    </RouterLink>
    <RouterLink
      class="cursor-pointer text-sm text-white/80 no-underline hover:underline hover:underline-white"
      :to="getCoverArtUrl(currentlyPlayingTrack.albumId)"
    >
      {{ currentlyPlayingTrack?.album }}
    </RouterLink>
    <img :src="coverArtUrl" class="w-full cursor-pointer rounded-lg object-cover" @error="onImageError" @click="() => router.push(`/albums/${currentlyPlayingTrack?.albumId}`)">
  </div>
</template>
