<script setup lang="ts">
import { formatTimeFromSeconds } from '~/logic/common'
import { currentlyPlayingItem, currentTime, seek } from '~/logic/playbackQueue'

const currentTimeFormatted = computed(() => formatTimeFromSeconds(currentTime.value))

const duration = computed(() => {
  if (currentlyPlayingItem.value.podcastEpisode && currentlyPlayingItem.value.podcastEpisode.duration) {
    return formatTimeFromSeconds(Number(currentlyPlayingItem.value.podcastEpisode.duration))
  }
  else if (currentlyPlayingItem.value.track && currentlyPlayingItem.value.track.duration) {
    return formatTimeFromSeconds(currentlyPlayingItem.value.track.duration)
  }
  return '0:00'
})

const maxDuration = computed(() => {
  if (currentlyPlayingItem.value.podcastEpisode && currentlyPlayingItem.value.podcastEpisode.duration) {
    return Number(currentlyPlayingItem.value.podcastEpisode.duration)
  }
  else if (currentlyPlayingItem.value.track && currentlyPlayingItem.value.track.duration) {
    return currentlyPlayingItem.value.track.duration
  }
  return 0
})

function handleSeek(event: Event) {
  const target = event.target as HTMLInputElement
  const seekTime = Number.parseFloat(target.value)
  seek(seekTime)
}
</script>

<template>
  <div v-if="currentlyPlayingItem.track || currentlyPlayingItem.podcastEpisode" class="mb-2 max-w-xs w-full flex flex-row items-center gap-2 lg:max-w-200 lg:max-w-lg sm:max-w-md lg:gap-2">
    <span id="currentTime" class="w-10 text-right text-sm text-muted lg:w-12 sm:w-10 sm:text-sm">
      {{ currentTimeFormatted }}
    </span>
    <input
      type="range"
      class="h-2 w-full cursor-pointer background-2 accent-primary1 lg:h-1"
      :max="maxDuration"
      :value="currentTime"
      @input="handleSeek($event)"
    />
    <span id="duration" class="w-10 text-sm text-muted lg:w-12 sm:w-10 sm:text-sm">
      {{ duration }}
    </span>
  </div>
</template>
