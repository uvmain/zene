<script setup lang="ts">
import { audioElement, currentVolume, toggleMute, volumeInput } from '~/logic/playbackQueue'

function handleInput(e: Event) {
  const value = (e.target as HTMLInputElement).value
  volumeInput(value)
}
</script>

<template>
  <div v-if="audioElement" id="volume-range-input" class="flex flex-row cursor-pointer items-center gap-2 lg:gap-2">
    <div @click="toggleMute()">
      <icon-nrk-media-volume-3 v-if="audioElement.volume > 0.66" class="text-sm text-muted sm:text-sm" />
      <icon-nrk-media-volume-2 v-else-if="audioElement.volume > 0.33" class="text-sm text-muted sm:text-sm" />
      <icon-nrk-media-volume-1 v-else class="text-sm text-muted sm:text-sm" />
    </div>
    <input
      type="range"
      class="h-2 w-20 cursor-pointer background-1 accent-primary2 lg:w-30 sm:w-24 active:accent-primary1"
      max="1"
      step="0.01"
      :value="currentVolume"
      @input="handleInput($event)"
    />
  </div>
</template>
