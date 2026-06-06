<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { RecycleScroller } from 'vue-virtual-scroller'
import { debugLog } from '~/logic/logger'
import { currentlyPlayingItem } from '~/logic/playbackQueue'
import 'vue-virtual-scroller/index.css'

const props = defineProps({
  primaryArtist: { type: String, required: false },
  tracks: { type: Object as PropType<SubsonicSong[]>, required: true },
  autoScrolling: { type: Boolean, default: true },
})

const scroller = useTemplateRef('scroller')

function scrollToActiveTrack() {
  if (!props.autoScrolling || !scroller.value) {
    debugLog('Auto-scrolling is disabled or scroller not available, skipping scrollToActiveTrack')
    return
  }

  const index = props.tracks.findIndex(track => track.id === currentlyPlayingItem.value.track?.id) - 1
  if (index === -1) {
    debugLog(`Currently playing track: ${currentlyPlayingItem.value.track?.id} not found in the track list, skipping scrollToActiveTrack`)
    return
  }
  const scrollerEl = scroller.value.el as HTMLElement
  const itemTop = index * 68
  const itemBottom = itemTop + 68
  const scrollTop = scrollerEl.scrollTop
  const scrollBottom = scrollTop + scrollerEl.clientHeight

  if (itemTop < scrollTop || itemBottom > scrollBottom) {
    scroller.value.scrollToItem(index, { smooth: true, align: 'start' })
  }
}

if (props.autoScrolling) {
  watch(() => currentlyPlayingItem.value, () => {
    scrollToActiveTrack()
  }, { deep: true })
}

onMounted(async () => {
  await nextTick()
  scrollToActiveTrack()
})
</script>

<template>
  <div v-if="tracks && tracks.length > 0" class="corner-cut background-2 lg:corner-cut-large">
    <div class="p-2 text-left flex flex-col h-full lg:p-4">
      <TracksHeader class="hidden md:grid" />
      <RecycleScroller
        v-slot="{ item, index }"
        ref="scroller"
        class="h-full"
        :items="tracks"
        :item-size="68"
        :flow-mode="true"
        key-field="id"
        @visible="scrollToActiveTrack()"
      >
        <Track
          :track="item"
          :track-index="index"
          :primary-artist="primaryArtist"
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
  scrollbar-color: hsl(from var(--main-colour) h 40% 65%) var(--colors-background-200);
  .dark & {
    scrollbar-color: hsl(from var(--main-colour) h 40% 35%) var(--colors-background-800);
  }
}
</style>
