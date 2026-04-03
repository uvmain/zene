<script setup lang="ts">
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import { fetchAlbums } from '~/logic/backendFetch'
import { artSizes, cacheBustAlbumArt, getCoverArtUrl, onImageError, parseReleaseDate } from '~/logic/common'
import { albumsStore } from '~/logic/store'

const props = defineProps({
  album: { type: Object as PropType<SubsonicAlbum>, required: false },
})

const router = useRouter()

const METADATA_COUNT = 20

const isShaking = ref(false)
const albumArray = ref<SubsonicAlbum[]>([])
const index = ref(0)
const showChangeArtModal = ref(false)
const artUpdatedTime = ref<string | undefined>(undefined)

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
  if (albumsStore.value.length > 0) {
    const randomAlbums = albumsStore.value.sort(() => 0.5 - Math.random()).slice(0, limit)
    albumArray.value = randomAlbums
    index.value = 0
    return
  }
  albumArray.value = await fetchAlbums({ type: 'random', size: limit })
  index.value = 0
}

const coverArtUrl = computed(() => {
  return getCoverArtUrl(currentAlbum.value.coverArt, artSizes.size200, artUpdatedTime.value)
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
  artUpdatedTime.value = Date.now().toString()
}

watch(() => props.album, (newAlbum) => {
  if (newAlbum) {
    albumArray.value = [newAlbum]
  }
  else {
    getRandomAlbums(METADATA_COUNT)
  }
}, { immediate: true })

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
  <section v-if="albumArray.length">
    <div
      class="corner-cut-large h-full w-full shadow-md shadow-zshade-500 overflow-hidden bg-cover bg-center dark:shadow-zshade-950"
      :style="{ backgroundImage: `url(${coverArtUrl})` }"
    >
      <div class="corner-cut-large flex h-full w-full items-center justify-between overflow-hidden background-grad-2 backdrop-blur-md">
        <div class="p-8">
          <div class="flex flex-row gap-2 h-30 lg:gap-6 lg:h-52">
            <img
              :src="coverArtUrl"
              class="border-muted rounded-md h-30 aspect-square cursor-pointer shadow-md shadow-zshade-500 lg:h-52 dark:shadow-zshade-900"
              loading="lazy"
              @error="onImageError"
              @click="navigateAlbum"
            >
            <div class="my-auto text-left flex flex-col gap-1 lg:gap-4">
              <div class="text-xl font-bold cursor-pointer line-clamp-1 lg:text-4xl" @click="navigateAlbum">
                {{ currentAlbum.name }}
              </div>
              <div class="text-lg cursor-pointer lg:text-xl" @click="navigateArtist()">
                {{ artistAndDate }}
              </div>
              <div v-if="currentAlbum.genres?.length > 0" class="hidden lg:(flex flex-nowrap gap-2 justify-start overflow-hidden)">
                <GenreBottle v-for="genre in currentAlbum.genres.filter(g => g.name !== '').slice(0, 8)" :key="genre.name" :genre="genre.name" />
              </div>
              <PlayButton class="flex justify-start" :album="currentAlbum" />
            </div>
          </div>
          <div class="opacity-50 right-2 top-2 absolute hover:opacity-100">
            <!-- Change Album Art section -->
            <div v-if="props.album">
              <ZButton
                @click="showChangeArtModal = true"
              >
                <div>
                  Change Art
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
            <div v-else class="p-3 corner-cut background-2 flex gap-2 right-0 top-0 absolute lg:p-2">
              <icon-nrk-chevron-left
                class="text-2xl opacity-80 cursor-pointer lg:text-3xl hover:text-primary2 active:opacity-100 hover:scale-105"
                @click="prevIndex"
              />
              <icon-nrk-dice-3
                class="text-2xl opacity-80 cursor-pointer lg:text-3xl hover:text-primary2 active:opacity-100 hover:scale-105"
                :class="{ shake: isShaking }"
                @click="handleDiceClick()"
              />
              <icon-nrk-chevron-right
                class="text-2xl opacity-80 cursor-pointer lg:text-3xl hover:text-primary2 active:opacity-100 hover:scale-105"
                @click="nextIndex"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped lang="css">
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
