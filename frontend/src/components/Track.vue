<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { formatTimeFromSeconds, getCoverArtUrl, onImageError } from '~/composables/logic'
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

const isTrackPlaying = computed(() => {
  const isPlaying = (currentlyPlayingTrack.value && currentlyPlayingTrack.value?.id === props.track.id) ?? false
  if (isPlaying && props.autoScrolling) {
    trackElement.value?.scrollIntoView({ behavior: 'smooth', block: 'center' })
  }
  return isPlaying
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

watch(playcountUpdatedMusicbrainzTrackId, (newtrack) => {
  if (props.track.musicBrainzId === newtrack) {
    playCount.value = (playCount.value ?? 0) + 1
  }
})
</script>

<template>
  <div
    ref="trackElement"
    class="group flex flex-row cursor-pointer items-center justify-between transition-colors duration-200 ease-out"
    :class="{
      'hover:bg-primary2/40': !isTrackPlaying,
      'background-3 bg-opacity-40': !isTrackPlaying && trackIndex % 2 === 0,
      'bg-primary1/40': isTrackPlaying,
    }"
    @click="handlePlay"
  >
    <div
      class="relative h-full w-15 flex items-center justify-center p-1"
    >
      <div class="relative translate-x-0 opacity-100 transition-all duration-300 group-hover:translate-x-[1rem] group-hover:opacity-0">
        <div v-if="!showAlbum">
          <div v-if="track.discNumber > 1" class="absolute left--4 text-sm text-muted opacity-50">
            {{ track.discNumber }}
          </div>
          <div>{{ track.track }}</div>
        </div>
        <span v-else>{{ trackIndex + 1 }}</span>
      </div>
      <icon-nrk-media-play
        class="absolute m-auto translate-x-[-1rem] text-xl opacity-0 transition-all duration-300 group-hover:translate-x-0 group-hover:opacity-100"
      />
    </div>
    <div>
      <div class="flex shrink">
        <div class="flex flex-col px-2">
          <RouterLink
            class="line-clamp-1 text-ellipsis text-primary no-underline md:text-lg hover:underline hover:underline-white"
            :to="`/tracks/${track.id}`"
            @click.stop
          >
            {{ track.title }}
          </RouterLink>
          <RouterLink
            class="hidden text-sm text-muted no-underline lg:block hover:underline hover:underline-white"
            :to="`/artists/${track.artistId}`"
            @click.stop
          >
            {{ track.artist }}
          </RouterLink>
        </div>
      </div>
    </div>

    <div v-if="showAlbum" class="relative w-15 flex items-center justify-center">
      <div v-if="track.discNumber > 1" class="absolute left-2 text-sm text-muted opacity-60">
        {{ track.discNumber }}
      </div>
      <div>
        {{ track.track }}
      </div>
    </div>

    <div v-if="showAlbum">
      <div class="flex flex-row items-center gap-2 px-1">
        <RouterLink
          :to="`/albums/${track.albumId}`"
          class="flex items-center"
          @click.stop
        >
          <img
            class="hidden size-10 rounded-sm object-cover md:block"
            :src="getCoverArtUrl(track.albumId, 40)"
            alt="Album Cover"
            :loading="trackIndex < 20 ? 'eager' : 'lazy'"
            width="40"
            height="40"
            @error="onImageError"
          />
        </RouterLink>
        <RouterLink
          class="text-muted no-underline hover:underline hover:underline-white"
          :to="`/albums/${track.albumId}`"
          @click.stop
        >
          {{ track.album }}
        </RouterLink>
      </div>
    </div>

    <div class="w-15 cursor-pointer text-center" @click="handlePlay">
      {{ playCount ?? 0 }}
    </div>
    <div class="w-15 cursor-pointer text-center" @click="handlePlay">
      {{ formatTimeFromSeconds(track.duration) }}
    </div>
  </div>
</template>
