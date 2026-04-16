<script setup lang="ts">
import { fetchPodcastChannel } from '~/logic/backendFetch'
import { artSizes, getCoverArtUrl, onImageError } from '~/logic/common'
import { currentlyPlayingItem } from '~/logic/playbackQueue'

const props = defineProps({
  large: { type: Boolean, default: false },
})

const router = useRouter()
const podcastChannelName = ref<string>('')
const showTrackModal = ref(false)

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

function handleTrackClick() {
  if (currentlyPlayingItem.value.track) {
    showTrackModal.value = true
  }
  else if (currentlyPlayingItem.value.podcastEpisode) {
    router.push(`/podcasts/${currentlyPlayingItem.value.podcastEpisode.id}`)
  }
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
    class="mt-auto flex flex-col space-y-2"
  >
    <div class="flex flex-col space-y-1">
      <div
        class="text-lg text-primary"
        :class="{ 'link': currentlyPlayingItem.track, 'podcast-link': currentlyPlayingItem.podcastEpisode }"
        @click="handleTrackClick"
      >
        {{ currentlyPlayingItem.track?.title || currentlyPlayingItem.podcastEpisode?.title }}
      </div>
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
    <TrackInfo v-if="currentlyPlayingItem.track" v-model="showTrackModal" :track="currentlyPlayingItem.track" @click.stop />
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
