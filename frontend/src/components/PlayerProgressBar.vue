<script setup lang="ts">
import { formatTimeFromSeconds } from '~/logic/common'
import { currentlyPlayingPodcastEpisode, currentlyPlayingTrack, currentTime, seek } from '~/logic/playbackQueue'

const currentTimeFormatted = computed(() => formatTimeFromSeconds(currentTime.value))

const duration = computed(() => {
  if (currentlyPlayingPodcastEpisode.value && currentlyPlayingPodcastEpisode.value.duration) {
    return formatTimeFromSeconds(Number(currentlyPlayingPodcastEpisode.value.duration))
  }
  else if (currentlyPlayingTrack.value && currentlyPlayingTrack.value.duration) {
    return formatTimeFromSeconds(currentlyPlayingTrack.value.duration)
  }
  return '0:00'
})

const maxDuration = computed(() => {
  if (currentlyPlayingPodcastEpisode.value && currentlyPlayingPodcastEpisode.value.duration) {
    return Number(currentlyPlayingPodcastEpisode.value.duration)
  }
  else if (currentlyPlayingTrack.value && currentlyPlayingTrack.value.duration) {
    return currentlyPlayingTrack.value.duration
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
  <div v-if="currentlyPlayingTrack || currentlyPlayingPodcastEpisode" class="mb-2 max-w-xs w-full flex flex-row items-center gap-2 lg:max-w-200 lg:max-w-lg sm:max-w-md lg:gap-2">
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
