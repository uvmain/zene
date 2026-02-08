<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { useElementVisibility } from '@vueuse/core'
import { useRouteTracks } from '~/composables/useRouteTracks'

const props = defineProps({
  showAlbum: { type: Boolean, default: false },
  tracks: { type: Object as PropType<SubsonicSong[]>, required: true },
  observerEnabled: { type: Boolean, default: false },
  autoScrolling: { type: Boolean, default: true },
})

const emits = defineEmits(['observerVisible'])

const { routeTracks } = useRouteTracks()

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
      <div class="flex flex-row justify-between gap-2 text-lg text-muted">
        <div class="w-6 cursor-pointer text-center md:w-15" @click="currentSortOption === 'trackNumberAsc' ? sorttracksBy('trackNumberDesc') : sorttracksBy('trackNumberAsc')">
          #
        </div>
        <div class="cursor-pointer px-2" @click="currentSortOption === 'titleAsc' ? sorttracksBy('titleDesc') : sorttracksBy('titleAsc')">
          Title
        </div>
        <div v-if="showAlbum" class="w-6 cursor-pointer text-center md:w-16" @click="currentSortOption === 'trackNumberAsc' ? sorttracksBy('trackNumberDesc') : sorttracksBy('trackNumberAsc')">
          track
        </div>
        <div v-if="showAlbum" class="cursor-pointer px-2" @click="currentSortOption === 'albumAsc' ? sorttracksBy('albumDesc') : sorttracksBy('albumAsc')">
          Album
        </div>
        <div class="w-6 cursor-pointer text-center text-sm md:w-16" @click="sorttracksBy('playCount')">
          Play Count
        </div>
        <div class="w-6 cursor-pointer text-center md:w-16" @click="currentSortOption === 'durationAsc' ? sorttracksBy('durationDesc') : sorttracksBy('durationAsc')">
          <icon-nrk-clock class="inline" />
        </div>
      </div>
      <Track
        v-for="(track, index) in routeTracks"
        :key="track.id"
        :track="track"
        :track-index="index"
        :show-album="showAlbum"
        :auto-scrolling="autoScrolling"
      />
    </div>
  </div>
  <Loading v-if="observerEnabled" ref="observer" class="mb-6 text-center text-muted" />
</template>
