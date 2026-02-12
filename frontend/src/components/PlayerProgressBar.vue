<script setup lang="ts">
import type { SubsonicPodcastEpisode } from '~/types/subsonicPodcasts'
import type { SubsonicSong } from '~/types/subsonicSong'
import { formatTimeFromSeconds } from '~/logic/common'
import { seek } from '~/logic/playbackQueue'

const props = defineProps({
  currentTimeInSeconds: { type: Number, default: 0 },
  currentlyPlayingTrack: { type: Object as PropType<SubsonicSong>, default: () => ({}) },
  currentlyPlayingPodcastEpisode: { type: Object as PropType<SubsonicPodcastEpisode>, default: () => ({}) },
})

const currentTime = computed(() => formatTimeFromSeconds(props.currentTimeInSeconds))

const duration = computed(() => {
  if (props.currentlyPlayingPodcastEpisode && props.currentlyPlayingPodcastEpisode.duration) {
    return formatTimeFromSeconds(Number(props.currentlyPlayingPodcastEpisode.duration))
  }
  else if (props.currentlyPlayingTrack && props.currentlyPlayingTrack.duration) {
    return formatTimeFromSeconds(props.currentlyPlayingTrack.duration)
  }
  return '0:00'
})

const inputDuration = computed(() => {
  if (props.currentlyPlayingPodcastEpisode && props.currentlyPlayingPodcastEpisode.duration) {
    return Number(props.currentlyPlayingPodcastEpisode.duration)
  }
  else if (props.currentlyPlayingTrack && props.currentlyPlayingTrack.duration) {
    return props.currentlyPlayingTrack.duration
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
  <div v-if="currentlyPlayingTrack.duration || currentlyPlayingPodcastEpisode.duration" class="max-w-xs w-full flex flex-row items-center gap-2 lg:max-w-200 lg:max-w-lg sm:max-w-md lg:gap-2">
    <span id="currentTime" class="w-10 text-right text-sm text-muted lg:w-12 sm:w-10 sm:text-sm">
      {{ currentTime }}
    </span>
    <input
      type="range"
      class="h-2 w-full cursor-pointer background-2 accent-primary1 lg:h-1"
      :max="inputDuration"
      :value="currentTimeInSeconds"
      @input="handleSeek($event)"
    />
    <span id="duration" class="w-10 text-sm text-muted lg:w-12 sm:w-10 sm:text-sm">
      {{ duration }}
    </span>
  </div>
</template>
