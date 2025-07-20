<script setup lang="ts">
import { useLogic } from '../composables/useLogic'
import { usePlaybackQueue } from '../composables/usePlaybackQueue'

const { currentlyPlayingTrack } = usePlaybackQueue()
const { getTrackUrl, getArtistUrl, getAlbumUrl } = useLogic()
const router = useRouter()

function onImageError(event: Event) {
  const target = event.target as HTMLImageElement
  target.onerror = null
  target.src = '/default-square.png'
}
</script>

<template>
  <div v-if="currentlyPlayingTrack" class="flex flex-col gap-2">
    <RouterLink
      class="cursor-pointer text-lg text-white/80 no-underline hover:underline hover:underline-white"
      :to="getTrackUrl(currentlyPlayingTrack.musicbrainz_track_id)"
    >
      {{ currentlyPlayingTrack?.title }}
    </RouterLink>
    <RouterLink
      class="cursor-pointer text-sm text-white/80 no-underline hover:underline hover:underline-white"
      :to="getArtistUrl(currentlyPlayingTrack.musicbrainz_artist_id)"
    >
      {{ currentlyPlayingTrack?.artist }}
    </RouterLink>
    <RouterLink
      class="cursor-pointer text-sm text-white/80 no-underline hover:underline hover:underline-white"
      :to="getAlbumUrl(currentlyPlayingTrack.musicbrainz_album_id)"
    >
      {{ currentlyPlayingTrack?.album }}
    </RouterLink>
    <img :src="currentlyPlayingTrack?.image_url" class="w-full cursor-pointer rounded-lg object-cover" @error="onImageError" @click="() => router.push(`/albums/${currentlyPlayingTrack?.musicbrainz_album_id}`)">
  </div>
</template>
