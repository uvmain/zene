<script setup lang="ts">
import { getCoverArtUrl, onImageError } from '~/composables/logic'
import { usePlaybackQueue } from '~/composables/usePlaybackQueue'

const { currentlyPlayingTrack, currentlyPlayingPodcastEpisode } = usePlaybackQueue()
const router = useRouter()

const coverArtUrl = computed(() => {
  if (currentlyPlayingTrack.value) {
    return getCoverArtUrl(currentlyPlayingTrack.value?.albumId, 200)
  }
  else {
    return '/default-square.png'
  }
})
</script>

<template>
  <div>
    <div v-if="currentlyPlayingTrack" class="flex flex-col gap-2">
      <RouterLink
        class="cursor-pointer text-lg text-primary no-underline hover:underline hover:underline-white"
        :to="`/tracks/${currentlyPlayingTrack?.id}`"
      >
        {{ currentlyPlayingTrack?.title }}
      </RouterLink>
      <RouterLink
        class="cursor-pointer text-sm text-muted no-underline hover:underline hover:underline-white"
        :to="`/artists/${currentlyPlayingTrack.artistId}`"
      >
        {{ currentlyPlayingTrack?.artist }}
      </RouterLink>
      <RouterLink
        class="cursor-pointer text-sm text-muted no-underline hover:underline hover:underline-white"
        :to="`/albums/${currentlyPlayingTrack.albumId}`"
      >
        {{ currentlyPlayingTrack?.album }}
      </RouterLink>
      <img
        :src="coverArtUrl"
        class="w-full cursor-pointer object-cover"
        @error="onImageError"
        @click="() => router.push(`/albums/${currentlyPlayingTrack?.albumId}`)"
      />
    </div>
    <div v-else-if="currentlyPlayingPodcastEpisode" class="flex flex-col gap-2">
      <RouterLink
        class="cursor-pointer text-lg text-primary no-underline hover:underline hover:underline-white"
        :to="`/podcasts/${currentlyPlayingPodcastEpisode?.channelId}`"
      >
        {{ currentlyPlayingPodcastEpisode?.title }}
      </RouterLink>
      <img
        :src="currentlyPlayingPodcastEpisode.coverArt"
        class="w-full cursor-pointer object-cover"
        @error="onImageError"
        @click="() => router.push(`/podcasts/${currentlyPlayingPodcastEpisode?.channelId}`)"
      />
    </div>
  </div>
</template>
