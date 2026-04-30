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
  <div v-if="currentlyPlayingItem.track || currentlyPlayingItem.podcastEpisode" class="mb-2 flex flex-row gap-2 max-w-xs w-full items-center lg:gap-2 lg:max-w-200 sm:max-w-md">
    <span id="currentTime" class="text-sm text-muted text-left w-12 sm:text-sm">
      {{ currentTimeFormatted }}
    </span>
    <input
      type="range"
      class="accent-primary-500 background-2 h-2 w-full cursor-pointer lg:h-1"
      :max="maxDuration"
      :value="currentTime"
      @input="handleSeek($event)"
    />
    <span id="duration" class="text-sm text-muted text-right w-12 sm:text-sm">
      {{ duration }}
    </span>
  </div>
</template>
