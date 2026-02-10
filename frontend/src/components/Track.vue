<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { postTrackStarred } from '~/logic/backendFetch'
import { albumArtSizes, formatTimeFromSeconds, getCoverArtUrl, onImageError } from '~/logic/common'
import { currentlyPlayingTrack, currentQueue, play, setCurrentlyPlayingTrack } from '~/logic/playbackQueue'
import { playcountUpdatedMusicbrainzTrackId } from '~/logic/playCounts'
import { routeTracks, setCurrentlyPlayingTrackInRouteTracks } from '~/logic/routeTracks'

const props = defineProps({
  track: { type: Object as PropType<SubsonicSong>, required: true },
  showAlbum: { type: Boolean, default: false },
  primaryArtist: { type: String, required: false },
  trackIndex: { type: Number, required: true },
  autoScrolling: { type: Boolean, default: true },
})

const trackElement = useTemplateRef('trackElement')
const isStarred = ref<string | undefined>(props.track.starred)

const artistIsAlbumArtist = computed(() => {
  if (!props.primaryArtist) {
    const albumArtists = props.track.albumArtists?.map(artist => artist.name.trim()) ?? []
    return albumArtists.includes(props.track.artist.trim())
  }
  return props.track.artist.trim() === props.primaryArtist.trim()
})

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
    class="group grid max-w-100% cursor-pointer items-center gap-4 px-2 py-1 text-base transition-colors duration-300 ease-out"
    :class="{
      'hover:bg-primary2/40': !isTrackPlaying,
      'dark:bg-zshade-700/60 bg-zshade-100/60': !isTrackPlaying && trackIndex % 2 === 0,
      'dark:bg-zshade-700/20 bg-zshade-100/20': !isTrackPlaying && trackIndex % 2 !== 0,
      'bg-primary1/40': isTrackPlaying,
      'corner-cut': trackIndex === 0,
      // grid-cols-[200px_1fr]
      'grid-cols-[60px_minmax(0,_1.2fr)_60px_minmax(0,_0.9fr)_minmax(0,_0.9fr)_60px_60px_60px]': showAlbum,
      'grid-cols-[60px_minmax(0,_1fr)_60px_minmax(0,_1fr)_60px_60px_60px]': !showAlbum,
    }"
    @click="handlePlay"
  >
    <!-- track number and play button -->
    <div class="relative flex items-center justify-center">
      <div class="relative translate-x-0 opacity-100 transition-all duration-300 group-hover:(translate-x-[1rem] opacity-0)">
        <div v-if="!showAlbum">
          <div v-if="track.discNumber > 1" class="absolute bottom-1px left--4 text-sm text-muted opacity-40">
            {{ track.discNumber }}:
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
    <div class="min-h-60px min-w-0 flex flex-row items-center gap-4 overflow-hidden">
      <div v-if="showAlbum" class="flex flex-shrink-0 items-center">
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
      <div class="min-w-0 flex-1">
        <RouterLink
          class="line-clamp-1 truncate text-lg text-primary no-underline hover:(underline underline-white)"
          :to="`/tracks/${track.id}`"
          @click.stop
        >
          {{ track.title }}
        </RouterLink>
        <RouterLink
          class="line-clamp-1 truncate text-sm text-muted no-underline hover:(underline underline-white)"
          :class="{ hidden: artistIsAlbumArtist }"
          :to="`/artists/${track.artistId}`"
          @click.stop
        >
          {{ track.artist }}
        </RouterLink>
      </div>
    </div>
    <!-- track duration -->
    <div class="text-center">
      {{ formatTimeFromSeconds(track.duration) }}
    </div>
    <!-- album -->
    <div v-if="showAlbum" class="min-w-0 flex flex-shrink-1">
      <RouterLink
        :to="`/albums/${track.albumId}`"
        class="line-clamp-1 truncate text-primary no-underline hover:(underline underline-white)"
        @click.stop
      >
        {{ track.album }}
      </RouterLink>
    </div>
    <!-- track genres -->
    <div class="fade-to-zero line-clamp-1 min-w-0 flex flex-row gap-1 truncate">
      <RouterLink
        v-for="genre in trackGenres"
        :key="genre"
        :to="`/genres/${genre}`"
        class="text-primary no-underline hover:(underline underline-white)"
        @click.stop
      >
        {{ genre }}<span v-if="genre !== trackGenres[trackGenres.length - 1]" class="text-muted">,</span>
      </RouterLink>
    </div>
    <!-- year -->
    <div class="cursor-pointer text-center">
      {{ track.year }}
    </div>
    <!-- starred -->
    <div class="flex cursor-pointer items-center justify-center" @click="toggleStarred" @click.stop>
      <icon-nrk-star-active v-if="isStarred" class="text-muted" />
      <icon-nrk-star v-else class="text-muted opacity-40 hover:opacity-100" />
    </div>
    <!-- play count -->
    <div class="cursor-pointer text-center">
      {{ playCount ?? 0 }}
    </div>
  </div>
</template>

<style lang="css" scoped>
.fade-to-zero {
  mask: linear-gradient(to right, rgba(0,0,0,1) 50%, rgba(0,0,0,0.4) 100%);
  -webkit-mask: linear-gradient(to right, rgba(0,0,0,1) 50%, rgba(0,0,0,0.4) 100%);
}
</style>
