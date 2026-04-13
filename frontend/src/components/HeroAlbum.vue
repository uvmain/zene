<script setup lang="ts">
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import { fetchAlbums } from '~/logic/backendFetch'
import { artSizes, cacheBustAlbumArt, getCoverArtUrl, onImageError, parseReleaseDate } from '~/logic/common'
import { albumsStore } from '~/logic/store'

const props = defineProps({
  album: { type: Object as PropType<SubsonicAlbum>, required: false },
})

const router = useRouter()

const albumArray = ref<SubsonicAlbum[]>([])
const showChangeArtModal = ref(false)
const artUpdatedTime = ref<string | undefined>(undefined)
const index = ref(0)

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

async function getRandomAlbums() {
  const limit = 100
  if (albumsStore.value.length > 0) {
    albumArray.value = albumsStore.value.toSorted(() => 0.5 - Math.random()).slice(0, limit)
    index.value = 0
    return
  }
  const response = await fetchAlbums({ type: 'random', size: limit })
  if (response) {
    albumArray.value = response
    albumsStore.value = response
    index.value = 0
  }
}

const coverArtUrl = computed(() => {
  return getCoverArtUrl(currentAlbum.value.coverArt, artSizes.size200, artUpdatedTime.value)
})

const artist = computed(() => {
  return currentAlbum.value.displayAlbumArtist ?? currentAlbum.value.artist ?? currentAlbum.value.displayArtist ?? 'Unknown Artist'
})

const date = computed(() => {
  if (currentAlbum.value.releaseDate) {
    return parseReleaseDate(currentAlbum.value.releaseDate)
  }
  else if (currentAlbum.value.year) {
    return currentAlbum.value.year.toString()
  }
  else {
    return ''
  }
})

function navigateAlbum() {
  router.push(`/albums/${currentAlbum.value.id}`)
}

function navigateArtist() {
  router.push(`/artists/${currentAlbum.value.artistId}`)
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
    getRandomAlbums()
  }
}, { immediate: true })

onBeforeMount(async () => {
  if (!props.album) {
    await getRandomAlbums()
  }
  else {
    albumArray.value = [props.album]
  }
})
</script>

<template>
  <section v-if="albumArray.length > 0">
    <div
      class="corner-cut shadow-background-500 shadow-md bg-cover bg-center lg:(corner-cut-large) dark:shadow-background-950"
      :style="{ backgroundImage: `url(${coverArtUrl})` }"
    >
      <div class="corner-cut background-grad-2 backdrop-blur-md lg:(corner-cut-large)">
        <div class="p-4 lg:p-8">
          <div class="flex flex-row gap-4 items-center">
            <img
              :src="coverArtUrl"
              class="border-muted rounded-md h-32 aspect-square cursor-pointer shadow-background-500 shadow-md lg:h-52 dark:shadow-background-900"
              loading="lazy"
              @error="onImageError"
              @click="navigateAlbum"
            >
            <div class="text-left flex flex-col gap-1 justify-center lg:gap-4">
              <div class="text-2xl font-bold link line-clamp-1 lg:text-4xl" @click="navigateAlbum()">
                {{ currentAlbum.name }}
              </div>
              <span class="text-xl link hidden lg:block" @click="navigateArtist()">
                {{ artist }} | {{ date }}
              </span>
              <div class="text-lg link line-clamp-1 lg:hidden" @click="navigateArtist()">
                {{ artist }}
              </div>
              <div v-if="currentAlbum.genres?.length > 0" class="flex-wrap gap-2 hidden justify-start overflow-hidden md:flex">
                <GenreBottle v-for="genre in currentAlbum.genres.slice(0, 3)" :key="genre.name" :genre="genre.name" class="hidden md:block lg:hidden" />
                <GenreBottle v-for="genre in currentAlbum.genres.slice(0, 6)" :key="genre.name" :genre="genre.name" class="hidden lg:block xl:hidden" />
                <GenreBottle v-for="genre in currentAlbum.genres.slice(0, 9)" :key="genre.name" :genre="genre.name" class="hidden xl:block 2xl:hidden" />
                <GenreBottle v-for="genre in currentAlbum.genres" :key="genre.name" :genre="genre.name" class="hidden 2xl:block" />
              </div>
              <div class="flex flex-row gap-4 lg:gap-8">
                <PlayButton class="flex justify-start" :album="currentAlbum" />
                <Fave v-model="currentAlbum.starred" :musicbrainz-id="currentAlbum.id" />
                <Rating v-model="currentAlbum.userRating" :musicbrainz-id="currentAlbum.id" />
              </div>
            </div>
          </div>
          <div class="opacity-50 right-1 top-1 absolute hover:opacity-100 lg:(right-2 top-2)">
            <!-- Change Album Art -->
            <div v-if="props.album">
              <ZButton
                class=""
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
            <!-- next hero album button -->
            <icon-nrk-media-next v-else class="footer-icon right-0 top-16 absolute" @click="nextIndex()">
            </icon-nrk-media-next>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
