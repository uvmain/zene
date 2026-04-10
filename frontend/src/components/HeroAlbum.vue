<script setup lang="ts">
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import { fetchAlbums, postStarToggle } from '~/logic/backendFetch'
import { artSizes, cacheBustAlbumArt, getCoverArtUrl, onImageError, parseReleaseDate } from '~/logic/common'
import { albumsStore, heroAlbumsIndex, heroAlbumsStore } from '~/logic/store'

const props = defineProps({
  album: { type: Object as PropType<SubsonicAlbum>, required: false },
})

const router = useRouter()

const albumArray = ref<SubsonicAlbum[]>([])
const showChangeArtModal = ref(false)
const artUpdatedTime = ref<string | undefined>(undefined)

const currentAlbum = computed(() => {
  return albumArray.value[heroAlbumsIndex.value]
})

function nextIndex() {
  if (heroAlbumsIndex.value < albumArray.value.length - 1) {
    heroAlbumsIndex.value += 1
  }
  else {
    heroAlbumsIndex.value = 0
  }
}

async function getRandomAlbums() {
  const limit = 100
  if (heroAlbumsStore.value.length > 0) {
    albumArray.value = heroAlbumsStore.value
    return
  }
  if (albumsStore.value.length > 0) {
    heroAlbumsStore.value = albumsStore.value.toSorted(() => 0.5 - Math.random()).slice(0, limit)
    albumArray.value = heroAlbumsStore.value
    return
  }
  const response = await fetchAlbums({ type: 'random', size: limit })
  if (response) {
    albumArray.value = response
    heroAlbumsStore.value = response
    heroAlbumsIndex.value = 0
  }
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

function actOnUpdatedArt() {
  showChangeArtModal.value = false
  cacheBustAlbumArt(`${currentAlbum.value.id}`)
  artUpdatedTime.value = Date.now().toString()
}

function toggleStarred(album: SubsonicAlbum) {
  if (album.starred) {
    postStarToggle(album.id, false)
    album.starred = undefined
  }
  else {
    postStarToggle(album.id, true)
    album.starred = new Date().toDateString()
  }
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
      class="corner-cut-large h-full w-full shadow-background-500 shadow-md overflow-hidden bg-cover bg-center dark:shadow-background-950"
      :style="{ backgroundImage: `url(${coverArtUrl})` }"
    >
      <div class="corner-cut-large flex h-full w-full items-center justify-between overflow-hidden background-grad-2 backdrop-blur-md">
        <div class="p-8">
          <div class="flex flex-row gap-2 h-30 lg:gap-6 lg:h-52">
            <img
              :src="coverArtUrl"
              class="border-muted rounded-md h-30 aspect-square cursor-pointer shadow-background-500 shadow-md lg:h-52 dark:shadow-background-900"
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
              <div class="mt-2 flex flex-row gap-8">
                <PlayButton class="flex justify-start" :album="currentAlbum" />
                <div class="flex cursor-pointer items-center justify-center" @click="toggleStarred(currentAlbum)" @click.stop>
                  <icon-nrk-star-active v-if="currentAlbum.starred" class="text-primary-400" />
                  <icon-nrk-star v-else class="text-muted opacity-70 hover:opacity-100" />
                </div>
              </div>
            </div>
          </div>
          <div class="opacity-50 right-2 top-2 absolute hover:opacity-100">
            <!-- Change Album Art -->
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
                :album="albumArray[heroAlbumsIndex]"
                @close="showChangeArtModal = false"
                @art-updated="actOnUpdatedArt"
              />
            </div>
            <!-- next hero album button -->
            <ZButton v-else size12 class="right-0 top-0 absolute" @click="nextIndex()">
              <icon-nrk-media-next class="footer-icon" />
            </ZButton>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
