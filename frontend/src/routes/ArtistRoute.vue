<script setup lang="ts">
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import type { SubsonicArtist } from '~/types/subsonicArtist'
import type { SubsonicSong } from '~/types/subsonicSong'
import { fetchAlbumsForArtist, fetchArtist, fetchArtistTopSongs } from '~/composables/backendFetch'
import { getCoverArtUrl, onImageError } from '~/composables/logic'
import { useRouteTracks } from '~/composables/useRouteTracks'

const route = useRoute()
const { routeTracks } = useRouteTracks()

const artist = ref<SubsonicArtist>()
const tracks = ref<SubsonicSong[]>()
const albumArtistAlbums = ref<SubsonicAlbum[]>([] as SubsonicAlbum[])
const artistAlbums = ref<SubsonicAlbum[]>([] as SubsonicAlbum[])
const canLoadMore = ref<boolean>(true)
const isLoading = ref<boolean>(false)
const limit = ref<number>(100)
const offset = ref<number>(0)

const musicbrainz_artist_id = computed(() => `${route.params.musicbrainz_artist_id}`)
const artistArtUrl = computed(() => getCoverArtUrl(musicbrainz_artist_id.value))

async function getData() {
  const promisesArray = [
    fetchArtist(musicbrainz_artist_id.value),
    fetchAlbumsForArtist(musicbrainz_artist_id.value),
    getTopSongs(),
  ]

  await Promise.all(promisesArray)
    .then(
      (results) => {
        artist.value = results[0] as SubsonicArtist
        const albums = results[1] as SubsonicAlbum[]
        albumArtistAlbums.value = albums.filter(album => album.albumArtists[0].name !== artist.value?.name) ?? []
        artistAlbums.value = albums.filter(album => album.albumArtists[0].name === artist.value?.name) ?? []
      },
    )
}

async function getTopSongs() {
  if (isLoading.value || !canLoadMore.value)
    return
  isLoading.value = true

  const newSongs = await fetchArtistTopSongs(musicbrainz_artist_id.value, limit.value, offset.value)

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
  offset.value = 0
  canLoadMore.value = true
}

watch(musicbrainz_artist_id, async () => {
  resetRefs()
  await getData()
})

onBeforeMount(async () => {
  await getData()
})
</script>

<template>
  <div class="flex flex-col gap-6">
    <div v-if="artist" class="h-80">
      <div
        class="corner-cut-large h-full w-full bg-cover bg-center"
        :style="{ backgroundImage: `url(${artistArtUrl})` }"
      >
        <div class="h-full w-full flex items-center gap-6 background-grad-2 align-middle backdrop-blur-lg">
          <div class="size-60">
            <img
              class="h-full w-full object-cover"
              :src="artistArtUrl"
              @error="onImageError"
            />
          </div>
          <div class="text-7xl text-primary font-bold">
            {{ artist.name }}
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
    <Tracks v-if="tracks" :tracks="tracks" :show-album="true" :observer-enabled="canLoadMore" @observer-visible="getTopSongs" />
  </div>
</template>
