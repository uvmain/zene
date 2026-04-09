<script setup lang="ts">
import { changeVolume, currentVolume, toggleMute } from '~/logic/volume'

function handleInput(e: Event) {
  const value = (e.target as HTMLInputElement).value
  changeVolume(value)
}
</script>

<template>
  <div id="volume-range-input" class="flex flex-row gap-2 cursor-pointer items-center lg:gap-2">
    <button :title="currentVolume === 0 ? 'Unmute' : 'Mute'" class="flex flex-row items-center" @click="toggleMute()">
      <icon-nrk-media-volume-3 v-if="currentVolume > 0.66" class="text-sm text-muted sm:text-sm" />
      <icon-nrk-media-volume-2 v-else-if="currentVolume > 0.33" class="text-sm text-muted sm:text-sm" />
      <icon-nrk-media-volume-1 v-else class="text-sm text-muted sm:text-sm" />
    </button>
    <input
      type="range"
      :title="`Volume: ${Math.round(currentVolume * 100)}%`"
      class="accent-primary-500 background-1 h-2 w-20 cursor-pointer active:accent-accent-500 lg:w-30 sm:w-24"
      max="1"
      step="0.01"
      :value="currentVolume"
      @input="handleInput($event)"
    />
  </div>
</template>
