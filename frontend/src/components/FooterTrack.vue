<script setup lang="ts">
import { artSizes, getCoverArtUrl, onImageError } from '~/logic/common'
import { currentlyPlayingItem } from '~/logic/playbackQueue'

const route = useRoute()
const showTrackModal = ref(false)
</script>

<template>
  <div
    v-if="currentlyPlayingItem.track"
    class="p-2 corner-cut bg-primary-500/30 flex flex-row gap-4 items-center justify-center overflow-hidden"
  >
    <div>
      <RouterLink
        :to="`/albums/${currentlyPlayingItem.track.albumId}`"
        @click.stop
      >
        <img
          class="rounded-sm size-60px shadow-background-500 shadow-sm object-cover dark:shadow-background-900"
          :src="getCoverArtUrl(currentlyPlayingItem.track.albumId, artSizes.size60)"
          alt="Album Cover"
          loading="lazy"
          width="60"
          height="60"
          @error="onImageError"
        />
      </RouterLink>
    </div>
    <div class="flex flex-shrink-1 flex-col min-w-0">
      <div
        class="text-lg text-primary text-left outline-none link no-underline truncate line-clamp-1"
        @click="showTrackModal = true"
        @click.stop
      >
        {{ currentlyPlayingItem.track.title }}
      </div>
      <RouterLink
        v-if="!route.path.startsWith('/artists/')"
        class="text-sm text-muted link no-underline truncate line-clamp-1"
        :to="`/artists/${currentlyPlayingItem.track.artistId}`"
        @click.stop
      >
        {{ currentlyPlayingItem.track.artist }}
      </RouterLink>
      <RouterLink
        v-else
        class="text-sm text-muted link no-underline truncate line-clamp-1"
        :to="`/albums/${currentlyPlayingItem.track.albumId}`"
        @click.stop
      >
        {{ currentlyPlayingItem.track.album }}
      </RouterLink>
    </div>
    <div class="text-center" @click="showTrackModal = true" @click.stop>
      <icon-nrk-more />
    </div>

    <TrackInfo v-model="showTrackModal" :track="currentlyPlayingItem.track" @click.stop />
  </div>
</template>
