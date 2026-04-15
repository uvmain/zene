<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { artSizes, getCoverArtUrl, onImageError } from '~/logic/common'

const props = defineProps({
  genre: { type: String, required: true },
  tracks: { type: Array as PropType<SubsonicSong[]>, required: true },
})

const current = ref<number>(0)
const nextIndex = ref<number>(0)
const intervalId = ref<NodeJS.Timeout | null>(null)

function updateRandomIndex() {
  current.value = nextIndex.value
  nextIndex.value = Math.floor(Math.random() * props.tracks.length)
  // prefetch the next image to ensure it's loaded when we switch to it
  const nextTrack = props.tracks[nextIndex.value]
  const img = new Image()
  img.src = getCoverArtUrl(nextTrack.coverArt, artSizes.size200)
}

const coverArtUrl = computed(() => {
  const track = props.tracks[current.value]
  return getCoverArtUrl(track.coverArt, artSizes.size200)
})

onBeforeMount(() => {
  intervalId.value = setInterval(updateRandomIndex, 8000)
})

onUnmounted(() => {
  if (intervalId.value) {
    clearInterval(intervalId.value)
  }
})
</script>

<template>
  <div
    class="fade corner-cut-large w-full shadow-background-500 shadow-md bg-cover bg-center dark:shadow-background-950"
    :style="{ backgroundImage: `url(${coverArtUrl})` }"
  >
    <div class="corner-cut background-grad-2 backdrop-blur-md lg:(corner-cut-large)">
      <div class="p-4 lg:p-8">
        <div class="flex flex-row gap-4 items-center">
          <div
            :style="{ backgroundImage: `url(${coverArtUrl})` }"
            class="fade border-muted rounded-md h-32 aspect-square cursor-pointer shadow-background-500 shadow-md lg:h-52 dark:shadow-background-900"
            loading="eager"
            @error="onImageError"
          />
          <div class="my-auto text-left flex flex-col gap-1 lg:gap-4">
            <div class="text-xl font-bold line-clamp-1 lg:text-4xl">
              {{ genre }}
            </div>
            <div class="mt-2 flex flex-row gap-8">
              <PlayButton class="flex justify-start" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="css" scoped>
.fade {
  @apply transition-all duration-2000;
}
</style>
