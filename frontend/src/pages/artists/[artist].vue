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
  const newSongs = await fetchArtistTopSongs(musicbrainzArtistId.value)
  tracks.value = newSongs
  routeTracks.value = tracks.value
}

function resetRefs() {
  tracks.value = []
  routeTracks.value = []
  artistGenres.value = []
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
  <div class="flex flex-col gap-4 lg:gap-6">
    <HeroArtist v-if="artist" :artist="artist" :genres="artistGenres" />
    <Albums v-if="artistAlbums.length > 0" :albums="artistAlbums" :order-disabled="true" />
    <Albums v-if="albumArtistAlbums.length > 0" :albums="albumArtistAlbums" title="Appears on albums" :order-disabled="true" />
    <Artists v-if="similarArtists.length > 0" :artists="similarArtists" title="Similar Artists" :order-disabled="true" />
    <Tracks
      v-if="tracks"
      :auto-scrolling="false"
      :tracks="tracks"
      :show-album="true"
      :primary-artist="artist?.name"
    />
  </div>
</template>
