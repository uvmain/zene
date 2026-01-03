<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { useElementVisibility } from '@vueuse/core'
import { formatTimeFromSeconds, getCoverArtUrl, onImageError } from '~/composables/logic'
import { usePlaybackQueue } from '~/composables/usePlaybackQueue'
import { usePlaycounts } from '~/composables/usePlaycounts'
import { useRouteTracks } from '~/composables/useRouteTracks'

const props = defineProps({
  showAlbum: { type: Boolean, default: false },
  tracks: { type: Object as PropType<SubsonicSong[]>, required: true },
  observerEnabled: { type: Boolean, default: false },
  autoScrolling: { type: Boolean, default: true },
})

const emits = defineEmits(['observerVisible'])

const { currentlyPlayingTrack, currentQueue, play, setCurrentlyPlayingTrack } = usePlaybackQueue()
const { routeTracks, setCurrentlyPlayingTrackInRouteTracks } = useRouteTracks()
const { playcount_updated_musicbrainz_track_id } = usePlaycounts()

const rowRefs = ref<any[]>([])
const currentRow = ref()
const observer = useTemplateRef('observer')
const observerIsVisible = useElementVisibility(observer)

watch(observerIsVisible, (newValue) => {
  if (newValue && props.observerEnabled) {
    emits('observerVisible')
  }
})

function isTrackPlaying(trackId: string): boolean {
  return (currentlyPlayingTrack.value && currentlyPlayingTrack.value?.id === trackId) ?? false
}

function handlePlay(track: SubsonicSong) {
  if (currentQueue.value?.tracks.some(queueTrack => queueTrack.id === track.id)) {
    setCurrentlyPlayingTrack(track)
  }
  else if (routeTracks.value?.some(queueTrack => queueTrack.musicBrainzId === track.musicBrainzId)) {
    setCurrentlyPlayingTrackInRouteTracks(track)
  }
  else {
    play(undefined, undefined, track)
  }
}

type SortOptions = 'titleAsc' | 'titleDesc' | 'artistAsc' | 'artistDesc' | 'albumAsc' | 'albumDesc' | 'playCount' | 'durationAsc' | 'durationDesc' | 'trackNumberAsc' | 'trackNumberDesc'
const currentSortOption = ref<SortOptions>('trackNumberAsc')

function sortTracksBy(sortOption: SortOptions) {
  switch (sortOption) {
    case 'titleAsc':
      routeTracks.value.sort((a, b) => a.title.localeCompare(b.title))
      currentSortOption.value = 'titleAsc'
      break
    case 'titleDesc':
      routeTracks.value.sort((a, b) => b.title.localeCompare(a.title))
      currentSortOption.value = 'titleDesc'
      break
    case 'artistAsc':
      routeTracks.value.sort((a, b) => a.artist.localeCompare(b.artist))
      currentSortOption.value = 'artistAsc'
      break
    case 'artistDesc':
      routeTracks.value.sort((a, b) => b.artist.localeCompare(a.artist))
      currentSortOption.value = 'artistDesc'
      break
    case 'albumAsc':
      routeTracks.value.sort((a, b) => a.album.localeCompare(b.album))
      currentSortOption.value = 'albumAsc'
      break
    case 'albumDesc':
      routeTracks.value.sort((a, b) => b.album.localeCompare(a.album))
      currentSortOption.value = 'albumDesc'
      break
    case 'playCount':
      routeTracks.value.sort((a, b) => (b.playCount ?? 0) - (a.playCount ?? 0))
      currentSortOption.value = 'playCount'
      break
    case 'durationAsc':
      routeTracks.value.sort((a, b) => a.duration - b.duration)
      currentSortOption.value = 'durationAsc'
      break
    case 'durationDesc':
      routeTracks.value.sort((a, b) => b.duration - a.duration)
      currentSortOption.value = 'durationDesc'
      break
    case 'trackNumberAsc':
      routeTracks.value.sort((a, b) => a.track - b.track)
      currentSortOption.value = 'trackNumberAsc'
      break
    case 'trackNumberDesc':
      routeTracks.value.sort((a, b) => b.track - a.track)
      currentSortOption.value = 'trackNumberDesc'
      break
  }
}

watch(currentlyPlayingTrack, async (newTrack) => {
  if (!newTrack)
    return
  await nextTick()
  const index = props.tracks.findIndex(track => track.musicBrainzId === newTrack.musicBrainzId)
  currentRow.value = rowRefs.value[index]
  if (props.autoScrolling && currentRow.value) {
    currentRow.value.scrollIntoView({ behavior: 'smooth', block: 'nearest' })
  }
})

watch(() => props.tracks, (newTracks) => {
  routeTracks.value = newTracks
}, { immediate: true })

watch(playcount_updated_musicbrainz_track_id, (newTrack) => {
  routeTracks.value?.forEach((track) => {
    if (track.musicBrainzId === newTrack) {
      track.playCount = (track.playCount ?? 0) + 1
    }
  })
})
</script>

<template>
  <div class="corner-cut-large background-2">
    <table class="h-full w-full p-4 text-left">
      <thead>
        <tr class="text-lg text-muted">
          <th class="w-15 cursor-pointer text-center" @click="currentSortOption === 'trackNumberAsc' ? sortTracksBy('trackNumberDesc') : sortTracksBy('trackNumberAsc')">
            #
          </th>
          <th class="cursor-pointer px-2" @click="currentSortOption === 'titleAsc' ? sortTracksBy('titleDesc') : sortTracksBy('titleAsc')">
            Title
          </th>
          <th v-if="showAlbum" class="w-16 cursor-pointer text-center" @click="currentSortOption === 'trackNumberAsc' ? sortTracksBy('trackNumberDesc') : sortTracksBy('trackNumberAsc')">
            Track
          </th>
          <th v-if="showAlbum" class="cursor-pointer px-2" @click="currentSortOption === 'albumAsc' ? sortTracksBy('albumDesc') : sortTracksBy('albumAsc')">
            Album
          </th>
          <th class="w-16 cursor-pointer text-center text-sm" @click="sortTracksBy('playCount')">
            Play Count
          </th>
          <th class="w-16 cursor-pointer text-center" @click="currentSortOption === 'durationAsc' ? sortTracksBy('durationDesc') : sortTracksBy('durationAsc')">
            <icon-nrk-clock class="inline" />
          </th>
        </tr>
        <tr>
          <td><HRule /></td>
          <td><HRule /></td>
          <td v-if="showAlbum">
            <HRule />
          </td>
          <td v-if="showAlbum">
            <HRule />
          </td>
          <td><HRule /></td>
          <td><HRule /></td>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(track, index) in routeTracks"
          :key="track.path"
          :ref="el => rowRefs[index] = el"
          class="group cursor-pointer transition-colors duration-200 ease-out"
          :class="{
            'hover:bg-primary2/40': !isTrackPlaying(track.id),
            'background-3 bg-opacity-50': !isTrackPlaying(track.id) && index % 2 === 0,
            'bg-primary1/40': isTrackPlaying(track.id),
          }"
          @click="handlePlay(track)"
        >
          <td
            class="relative h-full w-15 flex items-center justify-center"
          >
            <div class="relative translate-x-0 opacity-100 transition-all duration-300 group-hover:translate-x-[1rem] group-hover:opacity-0">
              <div v-if="!showAlbum">
                <div v-if="track.discNumber > 1" class="absolute left--4 text-sm text-muted opacity-50">
                  {{ track.discNumber }}
                </div>
                <div>{{ track.track }}</div>
              </div>
              <span v-else>{{ index + 1 }}</span>
            </div>
            <icon-nrk-media-play
              class="absolute m-auto translate-x-[-1rem] text-xl opacity-0 transition-all duration-300 group-hover:translate-x-0 group-hover:opacity-100"
            />
          </td>
          <td>
            <div class="flex shrink">
              <div class="flex flex-col px-2">
                <RouterLink
                  class="line-clamp-1 text-ellipsis text-lg text-primary no-underline hover:underline hover:underline-white"
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
          </td>

          <td v-if="showAlbum" class="relative w-15 flex items-center justify-center">
            <div v-if="track.discNumber > 1" class="absolute left-2 text-sm text-muted opacity-60">
              {{ track.discNumber }}
            </div>
            <div>
              {{ track.track }}
            </div>
          </td>

          <td v-if="showAlbum">
            <div class="flex flex-row items-center gap-2 px-1">
              <RouterLink
                :to="`/albums/${track.albumId}`"
                class="flex items-center"
                @click.stop
              >
                <img
                  class="size-10 object-cover"
                  :src="getCoverArtUrl(track.albumId, 40)"
                  alt="Album Cover"
                  :loading="index < 20 ? 'eager' : 'lazy'"
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
          </td>

          <td class="w-15 cursor-pointer text-center" @click="handlePlay(track)">
            {{ track.playCount ?? 0 }}
          </td>
          <td class="w-15 cursor-pointer text-center" @click="handlePlay(track)">
            {{ formatTimeFromSeconds(track.duration) }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
  <Loading v-if="observerEnabled" ref="observer" class="mb-6 text-center text-muted" />
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
