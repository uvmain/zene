<script setup lang="ts">
import { onKeyStroke } from '@vueuse/core'

defineProps({
  isPlaying: { type: Boolean, default: false },
})

const emits = defineEmits(['togglePlayback', 'stopPlayback', 'nextTrack', 'previousTrack', 'getRandomTracks'])

onKeyStroke('MediaPlayPause', (e) => {
  e.preventDefault()
  emits('togglePlayback')
})

onKeyStroke('MediaTrackPrevious', (e) => {
  e.preventDefault()
  emits('previousTrack')
})

onKeyStroke('MediaTrackNext', (e) => {
  e.preventDefault()
  emits('nextTrack')
})

onKeyStroke('MediaStop', (e) => {
  e.preventDefault()
  emits('stopPlayback')
})
</script>

<template>
  <div class="mt-2 flex flex-row items-center justify-center gap-x-2 lg:mt-2 lg:gap-x-4 sm:gap-x-2">
    <button id="repeat" class="media-control-button" @click="emits('stopPlayback')">
      <icon-nrk-media-stop class="footer-icon" />
    </button>
    <button id="shuffle" class="media-control-button" @click="emits('togglePlayback')">
      <icon-nrk-reorder class="footer-icon" />
    </button>
    <button id="back" class="media-control-button" @click="emits('previousTrack')">
      <icon-nrk-media-previous class="footer-icon" />
    </button>
    <ZButton
      id="play-pause"
      class="group/button"
      :primary="true"
      :size12="true"
      hover-text="Play/Pause"
      @click="emits('togglePlayback')"
    >
      <icon-nrk-media-play v-if="!isPlaying" class="footer-icon" />
      <icon-nrk-media-pause v-else class="footer-icon" />
    </ZButton>
    <button id="forward" class="media-control-button" @click="emits('nextTrack')">
      <icon-nrk-media-next class="footer-icon" />
    </button>
    <button id="repeat" class="media-control-button" @click="emits('togglePlayback')">
      <icon-nrk-media-jumpto class="footer-icon" />
    </button>
    <button
      id="shuffle"
      class="media-control-button"
      @click="emits('getRandomTracks')"
    >
      <icon-nrk-dice-3-active class="footer-icon" />
    </button>
  </div>
</template>

<style scoped>
.media-control-button {
  @apply h-10 w-10 flex cursor-pointer items-center justify-center border-none bg-white/0 font-semibold outline-none lg:h-12 lg:w-12 sm:h-10 sm:w-10;
}
</style>
