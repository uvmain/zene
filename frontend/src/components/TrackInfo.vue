<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { formatTimeFromSeconds, getCoverArtUrl, onImageError } from '~/logic/common'

const props = defineProps({
  track: { type: Object as PropType<SubsonicSong>, required: true },
})

const coverArtUrl = computed(() => {
  return getCoverArtUrl(props.track.musicBrainzId)
})
</script>

<template>
  <dialog class="m-auto p-4 corner-cut background-1">
    <div class="flex flex-col gap-8 lg:flex-row">
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
    <button :commandfor="`track-modal-${track.id}`" command="close">
      Close
    </button>
  </dialog>
</template>

<style lang="css" scoped>
dialog::backdrop {
  @apply backdrop-blur-lg;
}
</style>
