<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { formatTimeFromSeconds } from '~/composables/logic'

const props = defineProps({
  currentTimeInSeconds: { type: Number, default: 0 },
  currentlyPlayingTrack: { type: Object as PropType<SubsonicSong>, default: () => ({}) },
})

const emits = defineEmits(['seek'])

const currentTime = computed(() => formatTimeFromSeconds(props.currentTimeInSeconds))
const duration = computed(() => {
  return formatTimeFromSeconds(props.currentlyPlayingTrack ? props.currentlyPlayingTrack.duration : 0)
})

function seek(event: Event) {
  const target = event.target as HTMLInputElement
  const seekTime = Number.parseFloat(target.value)
  emits('seek', seekTime)
}
</script>

<template>
  <div v-if="currentlyPlayingTrack.duration" class="max-w-xs w-full flex flex-row items-center gap-2 lg:max-w-200 md:max-w-lg sm:max-w-md md:gap-2">
    <span id="currentTime" class="w-10 text-right text-sm text-muted md:w-12 sm:w-10 sm:text-sm">
      {{ currentTime }}
    </span>
    <input
      type="range"
      class="h-2 w-full cursor-pointer background-2 accent-primary1 md:h-1"
      :max="currentlyPlayingTrack ? currentlyPlayingTrack.duration : 0"
      :value="currentTimeInSeconds"
      @input="seek($event)"
    />
    <span id="duration" class="w-10 text-sm text-muted md:w-12 sm:w-10 sm:text-sm">
      {{ duration }}
    </span>
  </div>
</template>
