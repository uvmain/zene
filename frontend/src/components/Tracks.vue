<script setup lang="ts">
import type { TrackMetadata, TrackMetadataWithImageUrl } from '../types'
import { formatTime, getAlbumUrl, getArtistUrl, getTrackUrl } from '../composables/logic'
import { usePlaybackQueue } from '../composables/usePlaybackQueue'

defineProps({
  showAlbum: { type: Boolean, default: false },
  tracks: { type: Object as PropType<TrackMetadata[] | TrackMetadataWithImageUrl[]>, required: true },
})

const { currentlyPlayingTrack, currentPlaylist, play } = usePlaybackQueue()

function isTrackPlaying(trackId: string): boolean {
  return (currentlyPlayingTrack.value && currentlyPlayingTrack.value?.musicbrainz_track_id === trackId) ?? false
}
</script>

<template>
  <div class="rounded-lg bg-black/20 p-4">
    {{ currentPlaylist?.position }} {{ currentPlaylist?.tracks.length }}
    <table class="w-full table-auto text-left">
      <thead>
        <tr class="text-lg text-white/70">
          <th class="w-15 text-center">
            #
          </th>
          <th>Title</th>
          <th v-if="showAlbum">
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
          v-for="track, index in tracks"
          :key="track.title"
          class="group transition-colors duration-200 ease-out hover:bg-zene-200/20"
          :class="{ 'bg-white/02': index % 2 === 0, 'bg-white/40': isTrackPlaying(track.musicbrainz_track_id) }"
        >
          <td
            class="w-15 cursor-pointer text-center"
            @click="play(undefined, undefined, track)"
          >
            <span class="group-hover:hidden">{{ track.track_number }}</span>
            <icon-tabler-player-play-filled class="hidden text-xl group-hover:inline" />
          </td>
          <td class="flex flex-col">
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
          </td>
          <td v-if="showAlbum">
            <RouterLink
              class="cursor-pointer text-sm text-white/80 no-underline hover:underline hover:underline-white"
              :to="getAlbumUrl(track.musicbrainz_album_id)"
            >
              {{ track.album }}
            </RouterLink>
          </td>
          <td class="w-15 text-center">
            {{ formatTime(Number.parseInt(track.duration)) }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
