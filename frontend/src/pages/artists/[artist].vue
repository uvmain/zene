<script setup lang="ts">
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import type { SubsonicArtist, SubsonicArtistInfo } from '~/types/subsonicArtist'
import type { SubsonicSong } from '~/types/subsonicSong'
import { fetchArtist, fetchArtistInfo, fetchArtistTopSongs } from '~/logic/backendFetch'
import { routeTracks } from '~/logic/routeTracks'

const route = useRoute('/artists/[artist]')

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

async function getData() {
  const promisesArray = [
    fetchArtist(musicbrainzArtistId.value),
    fetchArtistInfo(musicbrainzArtistId.value, 10),
    getTopSongs(),
  ]

  await Promise.all(promisesArray)
    .then(
      (results) => {
        artist.value = results[0] as SubsonicArtist
        const info = results[1] as SubsonicArtistInfo
        const albums = artist.value.album ?? []
        artistAlbums.value = albums.filter(album => album.albumArtists[0].name === artist.value?.name) ?? []
        albumArtistAlbums.value = albums.filter(album => album.albumArtists[0].name !== artist.value?.name) ?? []
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
    <HeroArtist v-if="artist" :artist="artist" :genres="artistGenres" />
    <div v-if="artistAlbums.length > 0">
      <div class="mb-2 text-lg font-semibold">
        Albums
      </div>
      <div class="flex flex-wrap gap-6">
        <div v-for="album in artistAlbums" :key="album.id" class="flex flex-col gap-y-1 overflow-hidden transition duration-200 hover:scale-110">
          <Album :album="album" :show-artist="false" :show-date="false" />
        </div>
      </div>
    </div>
    <div v-if="albumArtistAlbums.length > 0">
      <div class="mb-2 text-lg font-semibold">
        Appears on albums
      </div>
      <div class="flex flex-wrap gap-6">
        <div v-for="album in albumArtistAlbums" :key="album.id" class="flex flex-col gap-y-1 overflow-hidden transition duration-200 hover:scale-110">
          <Album :album="album" />
        </div>
      </div>
    </div>
    <div v-if="similarArtists.length > 0">
      <div class="mb-2 text-lg font-semibold">
        Similar Artists
      </div>
      <div class="flex flex-wrap gap-6">
        <div v-for="similarArtist in similarArtists" :key="similarArtist.musicBrainzId" class="flex flex-col gap-y-1 overflow-hidden transition duration-200 hover:scale-110">
          <ArtistThumb :artist="similarArtist" />
        </div>
      </div>
    </div>
    <Tracks
      v-if="tracks"
      :auto-scrolling="false"
      :tracks="tracks"
      :show-album="true"
      :primary-artist="artist?.name"
      :observer-enabled="canLoadMore"
      @observer-visible="getTopSongs"
    />
  </div>
</template>
