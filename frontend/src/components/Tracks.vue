<script setup lang="ts">
import type { TrackMetadata, TrackMetadataWithImageUrl } from '../types'
import { formatTime, getArtistUrl } from '../composables/logic'
import { play } from '../composables/play'

defineProps({
  showAlbum: { type: Boolean, default: false },
  tracks: { type: Object as PropType<TrackMetadata[] | TrackMetadataWithImageUrl[]>, required: true },
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
          :class="{ 'bg-white/02': index % 2 === 0 }"
        >
          <td
            class="w-15 cursor-pointer text-center"
            @click="play(undefined, undefined, track)"
          >
            <span class="group-hover:hidden">{{ track.track_number }}</span>
            <icon-tabler-player-play-filled class="hidden text-xl group-hover:inline" />
          </td>
          <td>
            <div class="text-lg">
              {{ track.title }}
            </div>
            <a
              class="cursor-pointer text-sm text-white/80 no-underline hover:underline hover:underline-white"
              :href="getArtistUrl(track.musicbrainz_artist_id)"
            >
              {{ track.artist }}
            </a>
          </td>
          <td v-if="showAlbum">
            {{ track.album }}
          </td>
          <td class="w-15 text-center">
            {{ formatTime(Number.parseInt(track.duration)) }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
