<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { postTrackStarred } from '~/logic/backendFetch'
import { artSizes, formatTimeFromSeconds, getCoverArtUrl, onImageError } from '~/logic/common'
import { currentlyPlayingItem, handlePlay } from '~/logic/playbackQueue'
import { playcountUpdatedMusicbrainzTrackId } from '~/logic/playerUtils'

const props = defineProps({
  track: { type: Object as PropType<SubsonicSong>, required: true },
  showAlbum: { type: Boolean, default: false },
  primaryArtist: { type: String, required: false },
  trackIndex: { type: Number, required: true },
})

const route = useRoute()

const isStarred = ref<string | undefined>(props.track.starred)
const playCount = ref(props.track.playCount ?? 0)

const artistIsAlbumArtist = computed(() => {
  if (!props.primaryArtist) {
    const albumArtists = props.track.albumArtists?.map(artist => artist.name.trim()) ?? []
    return albumArtists.includes(props.track.artist.trim()) && route.path.startsWith('/artists/')
  }
  return props.track.artist.trim() === props.primaryArtist.trim()
})

const isTrackPlaying = computed(() => {
  return (currentlyPlayingItem.value.track && currentlyPlayingItem.value.track.id === props.track.id)
})

const trackGenres = computed(() => {
  const genres = props.track.genres.length ? props.track.genres.map(g => g.name.trim()) : []
  return genres
})

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
    class="group text-base px-2 py-1 gap-4 grid cursor-pointer transition-colors duration-300 ease-out items-center"
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
    @click="handlePlay(track)"
  >
    <!-- track number and play button -->
    <div class="flex items-center justify-center relative">
      <div class="opacity-100 translate-x-0 transition-all duration-300 relative group-hover:(opacity-0 translate-x-[1rem])">
        <div v-if="!showAlbum">
          <div v-if="track.discNumber > 1" class="text-sm text-muted opacity-40 bottom-1px left--4 absolute">
            {{ track.discNumber }}:
          </div>
          <div>{{ track.track }}</div>
        </div>
        <span v-else>{{ trackIndex + 1 }}</span>
      </div>
      <icon-nrk-media-play
        class="text-xl m-auto opacity-0 translate-x-[-1rem] transition-all duration-300 absolute group-hover:(opacity-100 translate-x-0)"
      />
    </div>
    <!-- album art, title and artist -->
    <div class="flex flex-row gap-4 min-h-60px min-w-0 items-center overflow-hidden">
      <div v-if="showAlbum" class="flex flex-shrink-0 items-center">
        <RouterLink
          :to="`/albums/${track.albumId}`"
          class="flex items-center"
          @click.stop
        >
          <img
            class="rounded-sm size-60px shadow-sm shadow-zshade-500 object-cover dark:shadow-zshade-900"
            :src="getCoverArtUrl(track.albumId, artSizes.size60)"
            alt="Album Cover"
            :loading="trackIndex < 20 ? 'eager' : 'lazy'"
            width="60"
            height="60"
            @error="onImageError"
          />
        </RouterLink>
      </div>
      <div class="flex flex-shrink-1 flex-col min-w-0">
        <RouterLink
          class="text-lg text-primary no-underline truncate line-clamp-1 hover:(underline underline-white)"
          :to="`/tracks/${track.id}`"
          @click.stop
        >
          {{ track.title }}
        </RouterLink>
        <RouterLink
          class="text-sm text-muted no-underline truncate line-clamp-1 hover:(underline underline-white)"
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
    <div v-if="showAlbum" class="min-w-0">
      <RouterLink
        :to="`/albums/${track.albumId}`"
        class="text-primary no-underline truncate line-clamp-1 hover:(underline underline-white)"
        @click.stop
      >
        {{ track.album }}
      </RouterLink>
    </div>
    <!-- track genres -->
    <div class="fade-out gap-1 min-w-0 truncate line-clamp-1">
      <RouterLink
        v-for="genre in trackGenres"
        :key="genre"
        :to="`/genres/${genre}`"
        class="text-primary no-underline hover:(underline underline-white)"
        @click.stop
      >
        {{ genre }}<span v-if="genre !== trackGenres[trackGenres.length - 1]" class="text-muted">, </span>
      </RouterLink>
    </div>
    <!-- year -->
    <div class="text-center cursor-pointer">
      {{ track.year }}
    </div>
    <!-- starred -->
    <div class="flex cursor-pointer items-center justify-center" @click="toggleStarred" @click.stop>
      <icon-nrk-star-active v-if="isStarred" class="text-muted" />
      <icon-nrk-star v-else class="text-muted opacity-40 hover:opacity-100" />
    </div>
    <!-- play count -->
    <div class="text-center cursor-pointer">
      {{ playCount ?? 0 }}
    </div>
  </div>
</template>

<style lang="css" scoped>
.fade-out {
  mask: linear-gradient(to right, rgba(0,0,0,1) 60%, rgba(0,0,0,0.4) 100%);
  -webkit-mask: linear-gradient(to right, rgba(0,0,0,1) 60%, rgba(0,0,0,0.4) 100%);
}
</style>
