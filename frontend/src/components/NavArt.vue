<script setup lang="ts">
import { fetchPodcastChannel } from '~/logic/backendFetch'
import { artSizes, getCoverArtUrl, onImageError } from '~/logic/common'
import { currentlyPlayingPodcastEpisode, currentlyPlayingTrack } from '~/logic/playbackQueue'

const podcastChannelName = ref<string>('')

const coverArtUrl = computed(() => {
  if (currentlyPlayingTrack.value) {
    return getCoverArtUrl(currentlyPlayingTrack.value?.albumId, artSizes.size200)
  }
  else if (currentlyPlayingPodcastEpisode.value) {
    return getCoverArtUrl(currentlyPlayingPodcastEpisode.value.coverArt, artSizes.size200)
  }
  else {
    return '/default-square.png'
  }
})

const trackTarget = computed(() => {
  return currentlyPlayingTrack.value ? `/tracks/${currentlyPlayingTrack.value.id}` : `/podcasts/${currentlyPlayingPodcastEpisode.value?.id}`
})

const artistTarget = computed(() => {
  return currentlyPlayingTrack.value ? `/artists/${currentlyPlayingTrack.value.artistId}` : `/podcasts/${currentlyPlayingPodcastEpisode.value?.channelId}`
})

const artTarget = computed(() => {
  return currentlyPlayingTrack.value ? `/albums/${currentlyPlayingTrack.value.albumId}` : `/podcasts/${currentlyPlayingPodcastEpisode.value?.channelId}`
})

async function fetchPodcastChannelName(channelId: string) {
  const response = await fetchPodcastChannel(channelId)
  podcastChannelName.value = response?.podcasts.channel[0].title || ''
}

watch(currentlyPlayingPodcastEpisode, (newEpisode) => {
  if (newEpisode) {
    fetchPodcastChannelName(newEpisode.channelId)
  }
})
</script>

<template>
  <div
    v-if="currentlyPlayingTrack || currentlyPlayingPodcastEpisode"
    class="mt-auto hidden flex flex-col lg:block space-y-2"
  >
    <div class="flex flex-col space-y-1">
      <RouterLink
        class="text-lg text-primary"
        :class="{ 'router-link': currentlyPlayingTrack, 'podcast-link': currentlyPlayingPodcastEpisode }"
        :to="trackTarget"
      >
        {{ currentlyPlayingTrack?.title || currentlyPlayingPodcastEpisode?.title }}
      </RouterLink>
      <RouterLink
        class="router-link text-sm text-muted"
        :to="artistTarget"
      >
        {{ currentlyPlayingTrack?.artist || podcastChannelName }}
      </RouterLink>
      <RouterLink
        v-if="currentlyPlayingTrack"
        class="router-link text-sm text-muted"
        :to="`/albums/${currentlyPlayingTrack.albumId}`"
      >
        {{ currentlyPlayingTrack?.album }}
      </RouterLink>
    </div>
    <RouterLink
      class="router-link text-sm text-muted"
      :to="artTarget"
    >
      <img
        :src="coverArtUrl"
        class="aspect-square w-full cursor-pointer rounded-md object-cover"
        @error="onImageError"
      />
    </RouterLink>
  </div>
</template>

<style scoped lang="css">
.router-link {
  @apply line-clamp-1 truncate cursor-pointer no-underline hover:underline hover:underline-white;
}
.podcast-link {
  @apply line-clamp-3 truncate cursor-pointer no-underline hover:underline hover:underline-white text-wrap;
}
</style>
