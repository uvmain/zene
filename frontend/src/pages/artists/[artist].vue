<script setup lang="ts">
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import type { SubsonicArtist, SubsonicArtistInfo } from '~/types/subsonicArtist'
import type { SubsonicSong } from '~/types/subsonicSong'
import { fetchAlbumsForArtist, fetchArtist, fetchArtistInfo, fetchArtistTopSongs } from '~/composables/backendFetch'
import { getCoverArtUrl, onImageError } from '~/composables/logic'
import { useRouteTracks } from '~/composables/useRouteTracks'

const route = useRoute()
const { routeTracks } = useRouteTracks()

const artist = ref<SubsonicArtist>()
const tracks = ref<SubsonicSong[]>()
const albumArtistAlbums = ref<SubsonicAlbum[]>([] as SubsonicAlbum[])
const artistAlbums = ref<SubsonicAlbum[]>([] as SubsonicAlbum[])
const similarArtists = ref<SubsonicArtist[]>([])
const artistGenres = ref<string[]>([])
const canLoadMore = ref<boolean>(true)
const isLoading = ref<boolean>(false)
const limit = ref<number>(100)
const offset = ref<number>(0)

const musicbrainzArtistId = computed(() => `${route.params.artist}`)
const artistArtUrl = computed(() => getCoverArtUrl(musicbrainzArtistId.value, 240))

async function getData() {
  const promisesArray = [
    fetchArtist(musicbrainzArtistId.value),
    fetchAlbumsForArtist(musicbrainzArtistId.value),
    fetchArtistInfo(musicbrainzArtistId.value, 10),
    getTopSongs(),
  ]

  await Promise.all(promisesArray)
    .then(
      (results) => {
        artist.value = results[0] as SubsonicArtist
        const albums = results[1] as SubsonicAlbum[]
        const info = results[2] as SubsonicArtistInfo
        albumArtistAlbums.value = albums.filter(album => album.albumArtists[0].name !== artist.value?.name) ?? []
        artistAlbums.value = albums.filter(album => album.albumArtists[0].name === artist.value?.name) ?? []
        const albumGenres = albums.flatMap(album => album.genres).filter(genre => genre.name !== '').map(genre => genre.name) ?? []
        artistGenres.value = Array.from(new Set(albumGenres)).slice(0, 12)
        similarArtists.value = info.similarArtists.map((artist) => {
          return {
            ...artist,
            musicBrainzId: artist.id,
            album: [],
          } as SubsonicArtist
        }) ?? []
      },
    )
}

async function getTopSongs() {
  if (isLoading.value || !canLoadMore.value)
    return
  isLoading.value = true

  const newSongs = await fetchArtistTopSongs(musicbrainzArtistId.value, limit.value, offset.value)

  if (newSongs.length === 0) {
    canLoadMore.value = false
  }
  else {
    tracks.value = tracks.value?.concat(newSongs) ?? newSongs
    routeTracks.value = tracks.value
    offset.value += newSongs.length

    if (newSongs.length < limit.value) {
      canLoadMore.value = false
    }
  }

  isLoading.value = false
}

function resetRefs() {
  tracks.value = []
  routeTracks.value = []
  artistGenres.value = []
  isLoading.value = false
  offset.value = 0
  canLoadMore.value = true
}

watch(musicbrainzArtistId, async () => {
  resetRefs()
  await getData()
})

onBeforeMount(async () => {
  await getData()
})
</script>

<template>
  <div class="flex flex-col gap-6">
    <div v-if="artist" class="h-60">
      <div
        class="corner-cut-large h-full w-full bg-cover bg-center"
        :style="{ backgroundImage: `url(${artistArtUrl})` }"
      >
        <div class="h-full w-full flex items-center gap-6 background-grad-2 align-middle backdrop-blur-lg">
          <div class="flex flex-row">
            <img
              class="object-cover"
              :src="artistArtUrl"
              width="240"
              height="240"
              @error="onImageError"
            />
            <div id="verticalRule" class="w-2 bg-zshade-100 dark:bg-zshade-400" />
          </div>
          <div class="flex flex-col gap-6">
            <div class="text-7xl text-primary font-bold">
              {{ artist.name }}
            </div>
            <div v-if="artistGenres.length > 0" class="flex flex-wrap justify-center gap-2 md:justify-start">
              <GenreBottle v-for="genre in artistGenres" :key="genre" :genre="genre" />
            </div>
          </div>
        </div>
      </div>
    </div>
    <div v-if="artistAlbums.length > 0">
      <div class="text-lg font-semibold">
        Albums
      </div>
      <div class="flex flex-wrap gap-6">
        <div v-for="album in artistAlbums" :key="album.id" class="flex flex-col gap-y-1 overflow-hidden transition duration-200 hover:scale-110">
          <Album :album="album" size="sm" />
        </div>
      </div>
    </div>
    <div v-if="albumArtistAlbums.length > 0">
      <div class="text-lg font-semibold">
        Appears on albums
      </div>
      <div class="flex flex-wrap gap-6">
        <div v-for="album in albumArtistAlbums" :key="album.id" class="flex flex-col gap-y-1 overflow-hidden transition duration-200 hover:scale-110">
          <Album :album="album" size="sm" />
        </div>
      </div>
    </div>
    <div v-if="similarArtists.length > 0">
      <div class="text-lg font-semibold">
        Similar Artists
      </div>
      <div class="flex flex-wrap gap-6">
        <div v-for="similarArtist in similarArtists" :key="similarArtist.musicBrainzId" class="flex flex-col gap-y-1 overflow-hidden transition duration-200 hover:scale-110">
          <ArtistThumb :artist="similarArtist" />
        </div>
      </div>
    </div>
    <Tracks v-if="tracks" :tracks="tracks" :show-album="true" :observer-enabled="canLoadMore" @observer-visible="getTopSongs" />
  </div>
</template>
