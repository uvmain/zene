<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { RecycleScroller } from 'vue-virtual-scroller'
import { currentlyPlayingItem, currentQueuePosition } from '~/logic/playbackQueue'
import { routeTracks } from '~/logic/routeTracks'
import 'vue-virtual-scroller/dist/vue-virtual-scroller.css'

const props = defineProps({
  showAlbum: { type: Boolean, default: false },
  primaryArtist: { type: String, required: false },
  tracks: { type: Object as PropType<SubsonicSong[]>, required: true },
  autoScrolling: { type: Boolean, default: true },
})

const scroller = useTemplateRef('scroller')

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

function scrollToActiveTrack() {
  if (!props.autoScrolling || !scroller.value)
    return
  scroller.value.scrollToItem(currentQueuePosition.value - 1, { smooth: true, align: 'start' })
}

watch(() => props.tracks, (newtracks) => {
  routeTracks.value = newtracks
}, { immediate: true })

if (props.autoScrolling) {
  watch(() => currentlyPlayingItem.value, () => {
    scrollToActiveTrack()
  }, { deep: true })
}

onMounted(() => {
  if (props.tracks && props.tracks.length > 0) {
    routeTracks.value = props.tracks
  }
  scrollToActiveTrack()
})
</script>

<template>
  <div v-if="routeTracks && routeTracks.length > 0" class="corner-cut-large background-2">
    <div class="p-2 text-left flex flex-col h-full lg:p-4">
      <div
        class="text-lg text-muted mb-2 px-2 py-1 gap-4 grid items-center"
        :class="{
          'grid-cols-[60px_minmax(0,_1.2fr)_60px_minmax(0,_0.9fr)_minmax(0,_0.9fr)_60px_60px_60px_2px]': showAlbum,
          'grid-cols-[60px_minmax(0,_1fr)_60px_minmax(0,_1fr)_60px_60px_60px]': !showAlbum,
        }"
      >
        <div class="text-center cursor-pointer" @click="currentSortOption === 'trackNumberAsc' ? sorttracksBy('trackNumberDesc') : sorttracksBy('trackNumberAsc')">
          #
        </div>
        <div class="cursor-pointer" @click="currentSortOption === 'titleAsc' ? sorttracksBy('titleDesc') : sorttracksBy('titleAsc')">
          Title
        </div>
        <div class="flex cursor-pointer items-center justify-center" @click="currentSortOption === 'durationAsc' ? sorttracksBy('durationDesc') : sorttracksBy('durationAsc')">
          <icon-nrk-clock class="text-base" />
        </div>
        <div v-if="showAlbum" class="cursor-pointer" @click="currentSortOption === 'albumAsc' ? sorttracksBy('albumDesc') : sorttracksBy('albumAsc')">
          Album
        </div>
        <div class="cursor-pointer">
          Genres
        </div>
        <div class="text-center cursor-pointer">
          Year
        </div>
        <div class="flex cursor-pointer cursor-pointer items-center justify-center">
          <icon-nrk-star class="text-base" />
        </div>
        <div class="cursor-pointer" @click="sorttracksBy('playCount')">
          Plays
        </div>
      </div>
      <RecycleScroller
        v-slot="{ item, index }"
        ref="scroller"
        class="h-full"
        :items="routeTracks"
        :item-size="68"
        key-field="id"
        @visible="scrollToActiveTrack()"
      >
        <Track
          :track="item"
          :track-index="index"
          :primary-artist="primaryArtist"
          :show-album="showAlbum"
        />
      </RecycleScroller>
    </div>
  </div>
  <div v-else class="text-muted flex h-48 items-center justify-center">
    No tracks found
  </div>
</template>

<style lang="css" scoped>
.vue-recycle-scroller {
  scrollbar-color: var(--colors-background-600) var(--colors-background-100);
  .dark & {
    scrollbar-color: var(--colors-background-800) var(--colors-background-600);
  }
}
</style>
