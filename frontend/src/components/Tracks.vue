<script setup lang="ts">
import type { TrackMetadataWithImageUrl } from '../types'
import { formatTime, getAlbumUrl, getArtistUrl, getTrackUrl } from '../composables/logic'
import { usePlaybackQueue } from '../composables/usePlaybackQueue'
import { useRouteTracks } from '../composables/useRouteTracks'

const props = defineProps({
  showAlbum: { type: Boolean, default: false },
  tracks: { type: Object as PropType<TrackMetadataWithImageUrl[]>, required: true },
})

const { currentlyPlayingTrack, currentQueue, play, setCurrentlyPlayingTrackInQueue } = usePlaybackQueue()
const { routeTracks, setCurrentlyPlayingTrackInRouteTracks } = useRouteTracks()

const rowRefs = ref<any[]>([])
const currentRow = ref()

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
    <table class="w-full table-auto text-left">
      <thead>
        <tr class="text-lg text-white/70">
          <th class="w-15 text-center">
            #
          </th>
          <th class="px-2">
            Title
          </th>
          <th v-if="showAlbum" class="px-2">
            Album
          </th>
          <th class="w-15 text-center">
            <icon-tabler-clock-hour-3 class="inline" />
          </th>
        </tr>
        <tr>
          <td><hr class="border-1 border-white/40 border-solid" /></td>
          <td><hr class="border-1 border-white/40 border-solid" /></td>
          <td v-if="showAlbum">
            <hr class="border-1 border-white/40 border-solid" />
          </td>
          <td><hr class="border-1 border-white/40 border-solid" /></td>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(track, index) in tracks"
          :key="track.title"
          :ref="el => rowRefs[index] = el"
          class="group transition-colors duration-200 ease-out hover:bg-zene-200/20"
          :class="{ 'bg-white/02': index % 2 === 0, 'bg-zene-200/40': isTrackPlaying(track.musicbrainz_track_id) }"
        >
          <td
            class="w-15 cursor-pointer text-center"
            @click="handlePlay(track)"
          >
            <span class="group-hover:hidden">{{ track.track_number }}</span>
            <icon-tabler-player-play-filled class="hidden text-xl group-hover:inline" />
          </td>
          <td>
            <div class="flex flex-row cursor-pointer px-2" @click="handlePlay(track)">
              <div class="flex flex-col">
                <RouterLink
                  class="cursor-pointer text-lg text-white/80 no-underline hover:underline hover:underline-white"
                  :to="getTrackUrl(track.musicbrainz_track_id)"
                >
                  {{ track.title }}
                </RouterLink>
                <RouterLink
                  class="cursor-pointer text-sm text-white/80 no-underline hover:underline hover:underline-white"
                  :to="getArtistUrl(track.musicbrainz_artist_id)"
                >
                  {{ track.artist }}
                </RouterLink>
              </div>
            </div>
          </td>
          <td v-if="showAlbum" class="cursor-pointer" @click="handlePlay(track)">
            <RouterLink
              class="cursor-pointer px-2 text-sm text-white/80 no-underline hover:underline hover:underline-white"
              :to="getAlbumUrl(track.musicbrainz_album_id)"
            >
              {{ track.album }}
            </RouterLink>
          </td>
          <td class="w-15 cursor-pointer text-center" @click="handlePlay(track)">
            {{ formatTime(Number.parseInt(track.duration)) }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
