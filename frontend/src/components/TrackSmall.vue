<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { artSizes, formatTimeFromSeconds, getCoverArtUrl, onImageError } from '~/logic/common'
import { currentlyPlayingItem, handlePlay } from '~/logic/playbackQueue'
import { playcountUpdatedMusicbrainzTrackId } from '~/logic/playerUtils'

const props = defineProps({
  track: { type: Object as PropType<SubsonicSong>, required: true },
  showAlbum: { type: Boolean, default: false },
  primaryArtist: { type: String, required: false },
  trackIndex: { type: Number, required: true },
  cornerCut: { type: Boolean, default: false },
})

const route = useRoute()

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

watch(playcountUpdatedMusicbrainzTrackId, (newtrack) => {
  if (props.track.musicBrainzId === newtrack) {
    playCount.value = (playCount.value ?? 0) + 1
  }
})
</script>

<template>
  <div
    class="group text-base px-2 py-1 gap-4 grid grid-cols-[40px_minmax(0,_1fr)_40px_20px] cursor-pointer transition-colors duration-300 ease-out items-center"
    :class="{
      'hover:bg-accent-500/30': !isTrackPlaying,
      'dark:bg-background-700/60 bg-background-100/60': !isTrackPlaying && trackIndex % 2 === 0,
      'dark:bg-background-700/20 bg-background-100/20': !isTrackPlaying && trackIndex % 2 !== 0,
      'bg-primary-500/30 corner-cut': isTrackPlaying,
      'corner-cut': trackIndex === 0 || cornerCut,
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
            class="rounded-sm size-60px shadow-background-500 shadow-sm object-cover dark:shadow-background-900"
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

    <div class="text-center">
      <icon-nrk-more />
    </div>
  </div>
</template>
