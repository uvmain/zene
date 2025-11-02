<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'

defineProps({
  currentlyPlayingTrack: { type: Object as PropType<SubsonicSong | null>, default: null },
  currentTime: { type: Number, default: 0 },
})

const showLyrics = ref(false)
</script>

<template>
  <div>
    <button
      id="lyrics"
      class="h-10 w-10 flex cursor-pointer items-center justify-center border-none bg-white/0 text-muted font-semibold outline-none md:h-12 md:w-12 sm:h-10 sm:w-10"
      @click="showLyrics = !showLyrics"
    >
      <icon-nrk-mening class="footer-icon" />
    </button>
    <LyricsModal
      v-if="showLyrics && currentlyPlayingTrack"
      :track="currentlyPlayingTrack"
      :current-seconds="currentTime"
      @close="showLyrics = false"
    />
  </div>
</template>
