<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { formatTime, getCoverArtUrl, onImageError } from '~/composables/logic'
import { usePlaybackQueue } from '~/composables/usePlaybackQueue'
import { usePlaycounts } from '~/composables/usePlaycounts'
import { useRouteTracks } from '~/composables/useRouteTracks'

const props = defineProps({
  showAlbum: { type: Boolean, default: false },
  tracks: { type: Object as PropType<SubsonicSong[]>, required: true },
})

const { currentlyPlayingTrack, currentQueue, play, setCurrentlyPlayingTrackInQueue } = usePlaybackQueue()
const { routeTracks, setCurrentlyPlayingTrackInRouteTracks } = useRouteTracks()
const { playcount_updated_musicbrainz_track_id } = usePlaycounts()

const rowRefs = ref<any[]>([])
const currentRow = ref()

function isTrackPlaying(trackId: string): boolean {
  return (currentlyPlayingTrack.value && currentlyPlayingTrack.value?.id === trackId) ?? false
}

function handlePlay(track: SubsonicSong) {
  if (currentQueue.value?.tracks.some(queueTrack => queueTrack.id === track.id)) {
    setCurrentlyPlayingTrackInQueue(track)
  }
  else if (routeTracks.value?.some(queueTrack => queueTrack.musicBrainzId === track.musicBrainzId)) {
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
  const index = props.tracks.findIndex(track => track.musicBrainzId === newTrack.musicBrainzId)
  currentRow.value = rowRefs.value[index]
  currentRow.value.scrollIntoView({ behavior: 'smooth', block: 'nearest' })
})

watch(playcount_updated_musicbrainz_track_id, (newTrack) => {
  routeTracks.value?.forEach((track) => {
    if (track.musicBrainzId === newTrack) {
      track.playcount = (track.playcount ?? 0) + 1
    }
  })
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
            Play Count
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
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(track, index) in tracks"
          :key="track.path"
          :ref="el => rowRefs[index] = el"
          class="group cursor-pointer transition-colors duration-200 ease-out hover:bg-zene-200/20"
          :class="{ 'bg-white/02': index % 2 === 0, 'bg-zene-200/40': isTrackPlaying(track.id) }"
          @click="handlePlay(track)"
        >
          <td
            class="relative h-full w-15 flex items-center justify-center"
          >
            <div class="relative translate-x-0 opacity-100 transition-all duration-300 group-hover:translate-x-[1rem] group-hover:opacity-0">
              <div v-if="!showAlbum">
                <div>
                  {{ track.discNumber }}
                </div>
                <span>{{ track.track }}</span>
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
                  :to="`/tracks/${track.id}`"
                  @click.stop
                >
                  {{ track.title }}
                </RouterLink>
                <RouterLink
                  class="text-sm text-white/80 no-underline hover:underline hover:underline-white"
                  :to="`/artists/${track.artistId}`"
                  @click.stop
                >
                  {{ track.artist }}
                </RouterLink>
              </div>
            </div>
          </td>

          <td v-if="showAlbum" class="relative w-15 flex items-center justify-center">
            <div>
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
                <img class="size-10 rounded-lg rounded-md object-cover" :src="getCoverArtUrl(track.albumId)" alt="Album Cover" @error="onImageError" />
              </RouterLink>
              <RouterLink
                class="text-white/80 no-underline hover:underline hover:underline-white"
                :to="`/albums/${track.albumId}`"
                @click.stop
              >
                {{ track.album }}
              </RouterLink>
            </div>
          </td>

          <td class="w-15 cursor-pointer text-center" @click="handlePlay(track)">
            {{ track.playcount ?? 0 }}
          </td>
          <td class="w-15 cursor-pointer text-center" @click="handlePlay(track)">
            {{ formatTime(track.duration) }}
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
