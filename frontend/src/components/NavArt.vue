<script setup lang="ts">
import { fetchPodcastChannel } from '~/logic/backendFetch'
import { artSizes, getCoverArtUrl, onImageError } from '~/logic/common'
import { currentlyPlayingItem } from '~/logic/playbackQueue'

const props = defineProps({
  large: { type: Boolean, default: false },
})

const podcastChannelName = ref<string>('')

const coverArtUrl = computed(() => {
  const artSize = props.large ? artSizes.size400 : artSizes.size200
  if (currentlyPlayingItem.value.track) {
    return getCoverArtUrl(currentlyPlayingItem.value.track.albumId, artSize)
  }
  else if (currentlyPlayingItem.value.podcastEpisode) {
    return getCoverArtUrl(currentlyPlayingItem.value.podcastEpisode.coverArt, artSize)
  }
  else {
    return '/default-square.png'
  }
})

const trackTarget = computed(() => {
  return currentlyPlayingItem.value.track ? `/tracks/${currentlyPlayingItem.value.track.id}` : `/podcasts/${currentlyPlayingItem.value.podcastEpisode?.id}`
})

const artistTarget = computed(() => {
  return currentlyPlayingItem.value.track ? `/artists/${currentlyPlayingItem.value.track.artistId}` : `/podcasts/${currentlyPlayingItem.value.podcastEpisode?.channelId}`
})

const artTarget = computed(() => {
  return currentlyPlayingItem.value.track ? `/albums/${currentlyPlayingItem.value.track.albumId}` : `/podcasts/${currentlyPlayingItem.value.podcastEpisode?.channelId}`
})

async function fetchPodcastChannelName(channelId: string) {
  const response = await fetchPodcastChannel(channelId)
  podcastChannelName.value = response?.podcasts.channel[0].title || ''
}

watch(currentlyPlayingItem, (newItem) => {
  if (newItem && newItem.podcastEpisode) {
    fetchPodcastChannelName(newItem.podcastEpisode.channelId)
  }
})
</script>

<template>
  <div
    v-if="currentlyPlayingItem.track || currentlyPlayingItem.podcastEpisode"
    class="hidden lg:(mt-auto flex flex-col space-y-2)"
  >
    <div class="flex flex-col space-y-1">
      <RouterLink
        class="text-lg text-primary"
        :class="{ 'router-link': currentlyPlayingItem.track, 'podcast-link': currentlyPlayingItem.podcastEpisode }"
        :to="trackTarget"
      >
        {{ currentlyPlayingItem.track?.title || currentlyPlayingItem.podcastEpisode?.title }}
      </RouterLink>
      <RouterLink
        class="router-link text-sm text-muted"
        :to="artistTarget"
      >
        {{ currentlyPlayingItem.track?.artist || podcastChannelName }}
      </RouterLink>
      <RouterLink
        v-if="currentlyPlayingItem.track"
        class="router-link text-sm text-muted"
        :to="`/albums/${currentlyPlayingItem.track.albumId}`"
      >
        {{ currentlyPlayingItem.track?.album }}
      </RouterLink>
    </div>
    <RouterLink
      class="router-link text-sm text-muted"
      :to="artTarget"
    >
      <img
        :src="coverArtUrl"
        class="rounded-md size-full aspect-square cursor-pointer"
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
