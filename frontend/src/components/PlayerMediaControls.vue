<script setup lang="ts">
import { onKeyStroke } from '@vueuse/core'
import { shuffleEnabled, toggleShuffle } from '~/logic/playerActions'

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
    <button id="shuffle" class="media-control-button" @click="toggleShuffle">
      <icon-ion-shuffle-sharp
        :class="{
          'footer-icon': !shuffleEnabled,
          'footer-icon-on': shuffleEnabled,
        }"
      />
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
