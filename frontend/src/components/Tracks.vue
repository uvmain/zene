<script setup lang="ts">
import type { TrackMetadataWithImageUrl } from '../types'
import { useIntersectionObserver } from '@vueuse/core'
import { formatTime, getAlbumUrl, getArtistUrl, getTrackUrl } from '../composables/logic'
import { usePlaybackQueue } from '../composables/usePlaybackQueue'
import { useRouteTracks } from '../composables/useRouteTracks'

const props = defineProps({
  showAlbum: { type: Boolean, default: false },
  tracks: { type: Object as PropType<TrackMetadataWithImageUrl[]>, required: true },
  canLoadMore: { type: Boolean, default: true },
})

const emits = defineEmits(['loadMore'])

const { currentlyPlayingTrack, currentQueue, play, setCurrentlyPlayingTrackInQueue } = usePlaybackQueue()
const { routeTracks, setCurrentlyPlayingTrackInRouteTracks } = useRouteTracks()

const rowRefs = ref<any[]>([])
const currentRow = ref()

const observer = ref<HTMLDivElement>()
const observerIndex = computed(() => {
  return props.tracks.length ? props.tracks.length - 3 : 0
})

useIntersectionObserver(
  observer,
  ([entry], _) => {
    if (entry?.isIntersecting && props.canLoadMore) {
      emits('loadMore')
    }
  },
)

function isTrackPlaying(trackId: string): boolean {
  return (currentlyPlayingTrack.value && currentlyPlayingTrack.value?.musicbrainz_track_id === trackId) ?? false
}

function handlePlay(track: TrackMetadataWithImageUrl) {
  if (currentQueue.value?.tracks.some(queueTrack => queueTrack.musicbrainz_track_id === track.musicbrainz_track_id)) {
    setCurrentlyPlayingTrackInQueue(track)
  }
  else if (routeTracks.value?.some(queueTrack => queueTrack.musicbrainz_track_id === track.musicbrainz_track_id)) {
    setCurrentlyPlayingTrackInRouteTracks(track)
  }
  else {
    play(undefined, undefined, track)
  }
}

function onImageError(event: Event) {
  const target = event.target as HTMLImageElement
  target.onerror = null
  target.src = '/default-square.png'
}

watch(currentlyPlayingTrack, async (newTrack) => {
  if (!newTrack)
    return
  await nextTick()
  const index = props.tracks.findIndex(track => track.musicbrainz_track_id === newTrack.musicbrainz_track_id)
  currentRow.value = rowRefs.value[index]
  currentRow.value.scrollIntoView({ behavior: 'smooth', block: 'nearest' })
})
</script>

<template>
  <div class="rounded-lg bg-black/20 p-4">
    <table class="h-full w-full table-auto text-left">
      <thead>
        <tr class="text-lg text-white/70">
          <th class="w-15 text-center">
            #
          </th>
          <th class="px-2">
            Title
          </th>
          <th v-if="showAlbum" class="w-16 text-center">
            Track
          </th>
          <th v-if="showAlbum" class="px-2">
            Album
          </th>
          <th class="w-16 text-center text-sm">
            My Plays
          </th>
          <th class="w-16 text-center text-sm">
            All Plays
          </th>
          <th class="w-16 text-center">
            <icon-tabler-clock-hour-3 class="inline" />
          </th>
        </tr>
        <tr>
          <td><hr class="border-1 border-white/40 border-solid" /></td>
          <td><hr class="border-1 border-white/40 border-solid" /></td>
          <td v-if="showAlbum">
            <hr class="border-1 border-white/40 border-solid" />
          </td>
          <td v-if="showAlbum">
            <hr class="border-1 border-white/40 border-solid" />
          </td>
          <td><hr class="border-1 border-white/40 border-solid" /></td>
          <td><hr class="border-1 border-white/40 border-solid" /></td>
          <td><hr class="border-1 border-white/40 border-solid" /></td>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(track, index) in tracks"
          :key="track.file_name"
          :ref="el => rowRefs[index] = el"
          class="group cursor-pointer transition-colors duration-200 ease-out hover:bg-zene-200/20"
          :class="{ 'bg-white/02': index % 2 === 0, 'bg-zene-200/40': isTrackPlaying(track.musicbrainz_track_id) }"
          @click="handlePlay(track)"
        >
          <td
            class="relative h-full w-15 flex items-center justify-center"
          >
            <div v-if="canLoadMore && index === observerIndex" ref="observer" class="invisible" />
            <div class="relative translate-x-0 opacity-100 transition-all duration-300 group-hover:translate-x-[1rem] group-hover:opacity-0">
              <div v-if="!showAlbum">
                <div v-if="Number.parseInt(track.total_discs) > 1" class="absolute bottom-0 text-xs opacity-50 -left-3">
                  {{ track.disc_number }}
                </div>
                <span>{{ track.track_number }}</span>
              </div>
              <span v-else>{{ index }}</span>
            </div>
            <icon-tabler-player-play-filled
              class="absolute m-auto translate-x-[-1rem] text-xl opacity-0 transition-all duration-300 group-hover:translate-x-0 group-hover:opacity-100"
            />
          </td>
          <td>
            <div class="flex shrink">
              <div class="flex flex-col px-2">
                <RouterLink
                  class="text-ellipsis text-lg text-white/80 no-underline hover:underline hover:underline-white"
                  :to="getTrackUrl(track.musicbrainz_track_id)"
                >
                  {{ track.title }}
                </RouterLink>
                <RouterLink
                  class="text-sm text-white/80 no-underline hover:underline hover:underline-white"
                  :to="getArtistUrl(track.musicbrainz_artist_id)"
                >
                  {{ track.artist }}
                </RouterLink>
              </div>
            </div>
          </td>

          <td v-if="showAlbum" class="relative w-15 flex items-center justify-center">
            <div v-if="Number.parseInt(track.total_discs) > 1" class="absolute bottom-0 left-3 text-xs opacity-50">
              {{ track.disc_number }}
            </div>
            <div>
              {{ track.track_number }}
            </div>
          </td>

          <td v-if="showAlbum">
            <div class="flex flex-row items-center gap-2 px-1">
              <RouterLink
                :to="getAlbumUrl(track.musicbrainz_album_id)"
                class="flex items-center"
              >
                <img class="size-10 rounded-lg rounded-md object-cover" :src="track.image_url" alt="Album Cover" @error="onImageError" />
              </RouterLink>
              <RouterLink
                class="text-white/80 no-underline hover:underline hover:underline-white"
                :to="getAlbumUrl(track.musicbrainz_album_id)"
              >
                {{ track.album }}
              </RouterLink>
            </div>
          </td>

          <td class="w-15 cursor-pointer text-center" @click="handlePlay(track)">
            {{ track.user_play_count }}
          </td>
          <td class="w-15 cursor-pointer text-center" @click="handlePlay(track)">
            {{ track.global_play_count }}
          </td>
          <td class="w-15 cursor-pointer text-center" @click="handlePlay(track)">
            {{ formatTime(Number.parseInt(track.duration)) }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.rotate-in {
  transform: rotate(-90deg);
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
.group:hover .rotate-in {
  transform: rotate(0deg);
}
</style>
