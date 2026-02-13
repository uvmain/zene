<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { useElementVisibility } from '@vueuse/core'
import { routeTracks } from '~/logic/store'

const props = defineProps({
  showAlbum: { type: Boolean, default: false },
  primaryArtist: { type: String, required: false },
  tracks: { type: Object as PropType<SubsonicSong[]>, required: true },
  observerEnabled: { type: Boolean, default: false },
  autoScrolling: { type: Boolean, default: true },
})

const emits = defineEmits(['observerVisible'])

const observer = useTemplateRef('observer')
const observerIsVisible = useElementVisibility(observer)

watch(observerIsVisible, (newValue) => {
  if (newValue && props.observerEnabled) {
    emits('observerVisible')
  }
})

type SortOptions = 'titleAsc' | 'titleDesc' | 'artistAsc' | 'artistDesc' | 'albumAsc' | 'albumDesc' | 'playCount' | 'durationAsc' | 'durationDesc' | 'trackNumberAsc' | 'trackNumberDesc'
const currentSortOption = ref<SortOptions>('trackNumberAsc')

function sorttracksBy(sortOption: SortOptions) {
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

watch(() => props.tracks, (newtracks) => {
  routeTracks.value = newtracks
}, { immediate: true })
</script>

<template>
  <div class="corner-cut-large background-2">
    <div class="h-full flex flex-col p-2 text-left lg:p-4">
      <div
        class="grid mb-2 items-center gap-4 p-2 text-lg text-muted"
        :class="{
          'grid-cols-[60px_minmax(0,_1.2fr)_60px_minmax(0,_0.9fr)_minmax(0,_0.9fr)_60px_60px_60px]': showAlbum,
          'grid-cols-[60px_minmax(0,_1fr)_60px_minmax(0,_1fr)_60px_60px_60px]': !showAlbum,
        }"
      >
        <div class="cursor-pointer text-center" @click="currentSortOption === 'trackNumberAsc' ? sorttracksBy('trackNumberDesc') : sorttracksBy('trackNumberAsc')">
          #
        </div>
        <div class="cursor-pointer" @click="currentSortOption === 'titleAsc' ? sorttracksBy('titleDesc') : sorttracksBy('titleAsc')">
          Title
        </div>
        <div class="mx-auto flex cursor-pointer items-center" @click="currentSortOption === 'durationAsc' ? sorttracksBy('durationDesc') : sorttracksBy('durationAsc')">
          <icon-nrk-clock class="text-base" />
        </div>
        <div v-if="showAlbum" class="cursor-pointer" @click="currentSortOption === 'albumAsc' ? sorttracksBy('albumDesc') : sorttracksBy('albumAsc')">
          Album
        </div>
        <div class="cursor-pointer">
          Genres
        </div>
        <div class="cursor-pointer text-center">
          Year
        </div>
        <div class="mx-auto flex cursor-pointer items-center">
          <icon-nrk-star class="text-base" />
        </div>
        <div class="cursor-pointer text-center" @click="sorttracksBy('playCount')">
          Plays
        </div>
      </div>
      <Track
        v-for="(track, index) in routeTracks"
        :key="track.id"
        :track="track"
        :track-index="index"
        :primary-artist="primaryArtist"
        :show-album="showAlbum"
        :auto-scrolling="autoScrolling"
      />
    </div>
  </div>
  <Loading v-if="observerEnabled" ref="observer" class="mb-6 text-center text-muted" />
</template>
