<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { artSizes, getCoverArtUrl, onImageError } from '~/logic/common'

const props = defineProps({
  genre: { type: String, required: true },
  tracks: { type: Array as PropType<SubsonicSong[]>, required: true },
})

const router = useRouter()

const current = ref<number>(0)
const nextIndex = ref<number>(0)
const intervalId = ref<NodeJS.Timeout | null>(null)

function updateRandomIndex() {
  current.value = nextIndex.value
  nextIndex.value = getRandomIndex()
  // prefetch the next image to ensure it's loaded when we switch to it
  const nextTrack = props.tracks[nextIndex.value]
  const img = new Image()
  img.src = getCoverArtUrl(nextTrack.coverArt, artSizes.size200)
}

function getRandomIndex(): number {
  return Math.floor(Math.random() * props.tracks.length)
}

const coverArtUrl = computed(() => {
  const track = props.tracks[current.value]
  return getCoverArtUrl(track.coverArt, artSizes.size200)
})

interface ArtistCount {
  name: string
  count: number
  musicBrainzId: string
}

const genreArtists = computed(() => {
  // get artists from tracks, filter out duplicates and sort by most frequent
  const artistCounts: ArtistCount[] = []
  for (const track of props.tracks) {
    const artist = track.displayArtist ?? track.artist ?? 'Unknown Artist'
    const existing = artistCounts.find(a => a.name === artist)
    if (existing) {
      existing.count += 1
    }
    else {
      artistCounts.push({ name: artist, count: 1, musicBrainzId: track.artistId })
    }
  }
  return artistCounts.sort((a, b) => b.count - a.count)
})

function navigateToArtist(artist: string) {
  router.push(`/artists/${artist}`)
}

onBeforeMount(() => {
  if (props.tracks.length <= 1) {
    return
  }
  intervalId.value = setInterval(updateRandomIndex, 4000)
  nextIndex.value = getRandomIndex()
  updateRandomIndex()
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
            class="fade border-muted rounded-md size-32 aspect-square cursor-pointer shadow-background-500 shadow-md bg-cover bg-center lg:size-52 dark:shadow-background-900"
            loading="eager"
            @error="onImageError"
          />
          <div class="my-auto text-left flex flex-col gap-1 lg:gap-4">
            <div class="text-xl font-bold line-clamp-1 lg:text-4xl">
              {{ genre }}
            </div>
            <div class="flex-wrap gap-2 hidden justify-start overflow-hidden md:flex">
              <ZButton v-for="artist in genreArtists.slice(0, 6)" :key="artist.name" @click="navigateToArtist(artist.musicBrainzId)">
                <span class="text-nowrap">{{ artist.name }}</span>
              </ZButton>
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
