<script setup lang="ts">
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import { fetchAlbums } from '~/composables/backendFetch'
import { albumArtSizes, cacheBustAlbumArt, getCoverArtUrl, onImageError, parseReleaseDate } from '~/composables/logic'

const props = defineProps({
  album: { type: Object as PropType<SubsonicAlbum>, required: false },
})

const router = useRouter()

const METADATA_COUNT = 20

const isShaking = ref(false)
const albumArray = ref<SubsonicAlbum[]>([])
const index = ref(0)
const showChangeArtModal = ref(false)

const currentAlbum = computed(() => {
  return albumArray.value[index.value]
})

function nextIndex() {
  if (index.value < albumArray.value.length - 1) {
    index.value += 1
  }
  else {
    index.value = 0
  }
}

function prevIndex() {
  if (index.value > 0) {
    index.value -= 1
  }
  else {
    index.value = albumArray.value.length - 1
  }
}

async function getRandomAlbums(limit: number) {
  albumArray.value = await fetchAlbums('random', limit)
  index.value = 0
}

const coverArtUrl = computed(() => {
  return getCoverArtUrl(currentAlbum.value.coverArt, albumArtSizes.size200)
})

const artist = computed(() => {
  return currentAlbum.value.displayAlbumArtist ?? currentAlbum.value.artist ?? currentAlbum.value.displayArtist ?? 'Unknown Artist'
})

const artistAndDate = computed(() => {
  if (currentAlbum.value.releaseDate) {
    return `${artist.value} • ${parseReleaseDate(currentAlbum.value.releaseDate)}`
  }
  else if (currentAlbum.value.year) {
    return `${artist.value} • ${currentAlbum.value.year}`
  }
  else {
    return artist.value
  }
})

function navigateAlbum() {
  router.push(`/albums/${currentAlbum.value.id}`)
}

function navigateArtist() {
  router.push(`/artists/${currentAlbum.value.artistId}`)
}

function handleDiceClick() {
  isShaking.value = true
  setTimeout(() => {
    isShaking.value = false
  }, 200)
  getRandomAlbums(METADATA_COUNT)
}

function actOnUpdatedArt() {
  showChangeArtModal.value = false
  cacheBustAlbumArt(`${currentAlbum.value.id}`)
}

onBeforeMount(async () => {
  if (!props.album) {
    await getRandomAlbums(METADATA_COUNT)
    index.value = 0
  }
  else {
    albumArray.value = [props.album]
  }
})
</script>

<template>
  <section v-if="albumArray.length" class="corner-cut-large overflow-hidden">
    <div
      class="h-full w-full bg-cover bg-center"
      :style="{ backgroundImage: `url(${coverArtUrl})` }"
    >
      <div class="h-full w-full flex items-center justify-between background-grad-2 backdrop-blur-md">
        <div class="p-8">
          <div class="h-30 flex flex-row gap-2 lg:h-52 lg:gap-6">
            <img
              :src="coverArtUrl"
              class="aspect-square h-30 cursor-pointer border-muted rounded-md shadow-md shadow-zshade-500 lg:h-52 dark:shadow-zshade-900"
              loading="lazy"
              @error="onImageError"
              @click="navigateAlbum"
            >
            <div class="my-auto flex flex-col gap-1 text-left lg:gap-4">
              <div class="line-clamp-1 cursor-pointer text-xl font-bold lg:text-4xl" @click="navigateAlbum">
                {{ currentAlbum.name }}
              </div>
              <div class="cursor-pointer text-lg lg:text-xl" @click="navigateArtist()">
                {{ artistAndDate }}
              </div>
              <div v-if="currentAlbum.genres?.length > 0" class="hidden lg:(block flex flex-nowrap justify-start gap-2 overflow-hidden)">
                <GenreBottle v-for="genre in currentAlbum.genres.filter(g => g.name !== '').slice(0, 8)" :key="genre.name" :genre="genre.name" />
              </div>
              <PlayButton class="flex justify-start" :album="currentAlbum" />
            </div>
          </div>
          <!-- Change Album Art section -->
          <div v-if="props.album" class="absolute right-2 top-2 opacity-70 hover:opacity-100">
            <ZButton
              @click="showChangeArtModal = true"
            >
              <div>
                Change Album Art
              </div>
            </ZButton>
            <ChangeAlbumArt
              v-if="showChangeArtModal"
              :album="albumArray[index]"
              @close="showChangeArtModal = false"
              @art-updated="actOnUpdatedArt"
            />
          </div>
          <!-- Dice and navigation buttons -->
          <div
            v-else
            id="album-dice"
            class="corner-cut absolute right-0 top-0 flex gap-2 background-2 p-3 lg:p-2"
          >
            <icon-nrk-chevron-left
              class="cursor-pointer text-2xl opacity-80 hover:scale-105 lg:text-3xl hover:text-primary2 active:opacity-100"
              @click="prevIndex"
            />
            <icon-nrk-dice-3
              class="cursor-pointer text-2xl opacity-80 hover:scale-105 lg:text-3xl hover:text-primary2 active:opacity-100"
              :class="{ shake: isShaking }"
              @click="handleDiceClick()"
            />
            <icon-nrk-chevron-right
              class="cursor-pointer text-2xl opacity-80 hover:scale-105 lg:text-3xl hover:text-primary2 active:opacity-100"
              @click="nextIndex"
            />
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
@keyframes shake {
  0% {
    transform: rotate(0deg);
  }
  25% {
    transform: rotate(-15deg);
  }
  50% {
    transform: rotate(15deg);
  }
  75% {
    transform: rotate(-15deg);
  }
  100% {
    transform: rotate(0deg);
  }
}

.shake {
  animation: shake 0.2s ease-in-out;
}
</style>
