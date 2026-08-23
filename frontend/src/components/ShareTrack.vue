<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { artSizes, formatTimeFromSeconds, getCoverArtUrl, onImageError } from '~/logic/common'
import { currentlyPlayingItem, handlePlay } from '~/logic/playbackQueue'
import { playcountUpdatedMusicbrainzTrackId } from '~/logic/playerUtils'

const props = defineProps({
  track: { type: Object as PropType<SubsonicSong>, required: true },
  trackIndex: { type: Number, required: true },
})

const route = useRoute()

const playCount = ref(props.track.playCount ?? 0)
const isStarred = ref<string | undefined>(props.track.starred)
const showTrackModal = ref(false)

const isTrackPlaying = computed(() => {
  return (currentlyPlayingItem.value.track && currentlyPlayingItem.value.track.id === props.track.id)
})

const trackGenres = computed(() => {
  const genres = props.track.genres.length ? props.track.genres.map(g => g.name.trim()) : []
  return genres
})

watch(playcountUpdatedMusicbrainzTrackId, (newtrack) => {
  if (props.track.musicBrainzId === newtrack) {
    playCount.value = (playCount.value ?? 0) + 1
  }
})
</script>

<template>
  <div>
    <div
      class="group text-base px-2 py-1 track-grid cursor-pointer transition-colors duration-300 ease-out"
      :class="{
        'hover:bg-main-500/20': !isTrackPlaying,
        'dark:bg-background-700/60 bg-background-100/60': !isTrackPlaying && trackIndex % 2 === 0,
        'dark:bg-background-700/20 bg-background-100/20': !isTrackPlaying && trackIndex % 2 !== 0,
        'bg-gradient-to-r from-main-500/40 via-main-500/20 corner-cut': isTrackPlaying,
        'corner-cut': trackIndex === 0,
      }"
      @click="handlePlay(track, route.path)"
    >
      <!-- track number and play button -->
      <div class="flex items-center justify-center relative">
        <div class="opacity-100 translate-x-0 transition-all duration-300 relative group-hover:(opacity-0 translate-x-[1rem])">
          <div>
            <div v-if="track.discNumber > 1" class="text-sm text-muted opacity-40 bottom-1px left--4 absolute">
              {{ track.discNumber }}:
            </div>
            <div>{{ track.track }}</div>
          </div>
        </div>
        <icon-nrk-media-play
          class="text-xl m-auto opacity-0 translate-x-[-1rem] transition-all duration-300 absolute group-hover:(opacity-100 translate-x-0)"
        />
      </div>
      <!-- album art, title and artist -->
      <div class="flex flex-row gap-4 min-h-60px min-w-0 items-center overflow-hidden">
        <div class="flex flex-shrink-0 items-center">
          <img
            class="rounded-sm size-60px shadow-background-500 shadow-sm object-cover dark:shadow-background-900"
            :src="getCoverArtUrl(track.albumId, artSizes.size60)"
            :alt="`Album art for ${track.title} by ${track.artist}`"
            :loading="trackIndex < 20 ? 'eager' : 'lazy'"
            width="60"
            height="60"
            @error="onImageError"
          />
        </div>
        <div class="flex flex-shrink-1 flex-col min-w-0">
          <div
            class="text-lg text-primary text-left outline-none link no-underline truncate line-clamp-1"
            @click="showTrackModal = true"
            @click.stop
          >
            {{ track.title }}
          </div>
          <div>
            {{ track.artist }}
          </div>
        </div>
      </div>
      <!-- track duration -->
      <div class="text-center">
        {{ formatTimeFromSeconds(track.duration) }}
      </div>

      <!-- album -->
      <div class="hidden md:(min-w-0 block)">
        <div class="text-primary link truncate line-clamp-1">
          {{ track.album }}
        </div>
      </div>

      <!-- track genres -->
      <div class="hidden lg:(min-w-0 truncate line-clamp-1)">
        <div
          v-for="genre in trackGenres"
          :key="genre"
          :title="genre"
          class="text-primary no-underline"
        >
          <span class="link">{{ genre }}</span><span v-if="genre !== trackGenres[trackGenres.length - 1]" class="text-muted">, </span>
        </div>
      </div>

      <!-- year -->
      <div class="hidden md:(text-center block cursor-pointer)">
        {{ track.year }}
      </div>

      <div class="text-center">
        <icon-nrk-more @click.stop="showTrackModal = true" />
      </div>
    </div>
    <TrackInfo v-model="showTrackModal" :track="track" />
  </div>
</template>
