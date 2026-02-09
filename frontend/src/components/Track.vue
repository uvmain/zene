<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { postTrackStarred } from '~/composables/backendFetch'
import { albumArtSizes, formatTimeFromSeconds, getCoverArtUrl, onImageError } from '~/composables/logic'
import { usePlaybackQueue } from '~/composables/usePlaybackQueue'
import { usePlaycounts } from '~/composables/usePlaycounts'
import { useRouteTracks } from '~/composables/useRouteTracks'

const props = defineProps({
  track: { type: Object as PropType<SubsonicSong>, required: true },
  showAlbum: { type: Boolean, default: false },
  trackIndex: { type: Number, required: true },
  autoScrolling: { type: Boolean, default: true },
})

const { currentlyPlayingTrack, currentQueue, play, setCurrentlyPlayingTrack } = usePlaybackQueue()
const { routeTracks, setCurrentlyPlayingTrackInRouteTracks } = useRouteTracks()
const { playcountUpdatedMusicbrainzTrackId } = usePlaycounts()

const trackElement = useTemplateRef('trackElement')
const isStarred = ref<string | undefined>(props.track.starred)

const isTrackPlaying = computed(() => {
  const isPlaying = (currentlyPlayingTrack.value && currentlyPlayingTrack.value?.id === props.track.id) ?? false
  if (isPlaying && props.autoScrolling) {
    trackElement.value?.scrollIntoView({ behavior: 'smooth', block: 'center' })
  }
  return isPlaying
})

const trackGenres = computed(() => {
  const genres = props.track.genres.length ? props.track.genres.map(g => g.name.trim()) : []
  return genres
})

function handlePlay() {
  if (currentQueue.value?.tracks.some(queueTrack => queueTrack.id === props.track.id)) {
    setCurrentlyPlayingTrack(props.track)
  }
  else if (routeTracks.value?.some(queueTrack => queueTrack.musicBrainzId === props.track.musicBrainzId)) {
    setCurrentlyPlayingTrackInRouteTracks(props.track)
  }
  else {
    play(undefined, undefined, props.track)
  }
}

const playCount = ref(props.track.playCount ?? 0)

function toggleStarred() {
  if (isStarred.value) {
    postTrackStarred(props.track.id, false)
    isStarred.value = undefined
  }
  else {
    postTrackStarred(props.track.id, true)
    isStarred.value = new Date().toDateString()
  }
}

watch(playcountUpdatedMusicbrainzTrackId, (newtrack) => {
  if (props.track.musicBrainzId === newtrack) {
    playCount.value = (playCount.value ?? 0) + 1
  }
})
</script>

<template>
  <div
    ref="trackElement"
    class="group corner-cut max-w-100% flex flex-row cursor-pointer items-center gap-2 p-2 text-base transition-colors duration-300 ease-out"
    :class="{
      'hover:bg-primary2/40': !isTrackPlaying,
      'dark:bg-zshade-700/60 bg-zshade-100/60': !isTrackPlaying && trackIndex % 2 === 0,
      'dark:bg-zshade-700/20 bg-zshade-100/20': !isTrackPlaying && trackIndex % 2 !== 0,
      'bg-primary1/40': isTrackPlaying,
    }"
    @click="handlePlay"
  >
    <!-- track number and play button -->
    <div class="relative h-full w-15 flex items-center justify-center">
      <div class="relative translate-x-0 opacity-100 transition-all duration-300 group-hover:(translate-x-[1rem] opacity-0)">
        <div v-if="!showAlbum">
          <div v-if="track.discNumber > 1" class="absolute left--4 text-sm text-muted opacity-50">
            {{ track.discNumber }}
          </div>
          <div>{{ track.track }}</div>
        </div>
        <span v-else>{{ trackIndex + 1 }}</span>
      </div>
      <icon-nrk-media-play
        class="absolute m-auto translate-x-[-1rem] text-xl opacity-0 transition-all duration-300 group-hover:(translate-x-0 opacity-100)"
      />
    </div>
    <!-- album art, title and artist -->
    <div class="flex flex-row items-center gap-2">
      <div v-if="showAlbum" class="flex flex-row items-center gap-2">
        <RouterLink
          :to="`/albums/${track.albumId}`"
          class="flex items-center"
          @click.stop
        >
          <img
            class="size-60px rounded-sm object-cover shadow-sm shadow-zshade-500 dark:shadow-zshade-900"
            :src="getCoverArtUrl(track.albumId, albumArtSizes.size60)"
            alt="Album Cover"
            :loading="trackIndex < 20 ? 'eager' : 'lazy'"
            width="60"
            height="60"
            @error="onImageError"
          />
        </RouterLink>
      </div>
      <div class="flex flex-col px-2">
        <RouterLink
          class="line-clamp-1 truncate text-lg text-primary no-underline hover:(underline underline-white)"
          :to="`/tracks/${track.id}`"
          @click.stop
        >
          {{ track.title }}
        </RouterLink>
        <RouterLink
          class="line-clamp-1 truncate text-sm text-muted no-underline hover:(underline underline-white)"
          :to="`/artists/${track.artistId}`"
          @click.stop
        >
          {{ track.artist }}
        </RouterLink>
      </div>
    </div>
    <!-- track duration -->
    <div class="w-15 cursor-pointer text-center">
      {{ formatTimeFromSeconds(track.duration) }}
    </div>
    <!-- album -->
    <div v-if="showAlbum" class="flex-1 cursor-pointer text-center">
      <RouterLink
        :to="`/albums/${track.albumId}`"
        class="line-clamp-1 truncate text-primary no-underline hover:(underline underline-white)"
        @click.stop
      >
        {{ track.album }}
      </RouterLink>
    </div>
    <!-- track genres -->
    <div class="flex-1 cursor-pointer text-center">
      <RouterLink
        :to="`/genres/${trackGenres.join(',')}`"
        class="line-clamp-1 truncate text-primary no-underline hover:(underline underline-white)"
        @click.stop
      >
        {{ trackGenres.join(', ') }}
      </RouterLink>
    </div>
    <!-- year -->
    <div class="w-15 cursor-pointer text-center">
      {{ track.year }}
    </div>
    <!-- starred -->
    <div class="w-15 flex cursor-pointer items-center justify-center" @click="toggleStarred" @click.stop>
      <icon-nrk-star-active v-if="isStarred" class="text-muted" />
      <icon-nrk-star v-else class="text-muted opacity-40 hover:opacity-100" />
    </div>
    <!-- play count -->
    <div class="w-15 cursor-pointer text-center">
      {{ playCount ?? 0 }}
    </div>
  </div>
</template>
