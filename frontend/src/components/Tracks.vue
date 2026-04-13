<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { RecycleScroller } from 'vue-virtual-scroller'
import { currentlyPlayingItem, currentQueuePosition } from '~/logic/playbackQueue'
import { routeTracks } from '~/logic/routeTracks'
import 'vue-virtual-scroller/dist/vue-virtual-scroller.css'

const props = defineProps({
  primaryArtist: { type: String, required: false },
  tracks: { type: Object as PropType<SubsonicSong[]>, required: true },
  autoScrolling: { type: Boolean, default: true },
})

const scroller = useTemplateRef('scroller')

function scrollToActiveTrack() {
  if (!props.autoScrolling || !scroller.value)
    return

  const index = currentQueuePosition.value - 1
  const scrollerEl = scroller.value.el as HTMLElement
  const itemTop = index * 68
  const itemBottom = itemTop + 68
  const scrollTop = scrollerEl.scrollTop
  const scrollBottom = scrollTop + scrollerEl.clientHeight

  if (itemTop < scrollTop || itemBottom > scrollBottom) {
    scroller.value.scrollToItem(index, { smooth: true, align: 'start' })
  }
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
  <div v-if="routeTracks && routeTracks.length > 0" class="corner-cut background-2 lg:corner-cut-large">
    <div class="p-2 text-left flex flex-col h-full lg:p-4">
      <TracksHeader class="hidden md:grid" />
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
