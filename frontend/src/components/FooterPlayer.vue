<script setup lang="ts">
import { apiKey } from '~/logic/store'

const loggedIn = computed(() => {
  return !!apiKey.value
})
</script>

<template>
  <footer
    v-if="loggedIn"
    class="border-0 border-t-1 border-background-600 border-solid background-2 w-full bottom-0 sticky"
  >
    <div class="px-8 py-4 flex flex-col gap-6 items-center lg:flex-row">
      <div
        class="flex flex-grow flex-col h-full w-full items-center justify-center"
      >
        <FooterTrack
          class="mb-3 flex max-w-xs w-full items-center lg:max-w-200 sm:max-w-md lg:hidden"
        />
        <PlayerAudio />
        <PlayerProgressBar />
        <PlayerMediaControls />
      </div>
      <div class="flex flex-row gap-4 items-center lg:gap-6">
        <!-- <GoogleCastButton /> -->
        <PlayerVisualizerButton />
        <PlayerLyricsButton />
        <PlayerQueueButton />
        <PlayerVolumeSlider />
      </div>
    </div>
  </footer>
</template>
