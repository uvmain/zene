<script setup lang="ts">
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import type { SubsonicArtist } from '~/types/subsonicArtist'
import type { SubsonicSong } from '~/types/subsonicSong'
import { fetchAlbum, fetchArtist } from '~/composables/backendFetch'
import { getCoverArtUrl, onImageError } from '~/composables/logic'
import { useRouteTracks } from '~/composables/useRouteTracks'

const route = useRoute()
const { routeTracks } = useRouteTracks()

const artist = ref<SubsonicArtist>()
const tracks = ref<SubsonicSong[]>()
const albums = ref<SubsonicAlbum[]>()

const musicbrainz_artist_id = computed(() => `${route.params.musicbrainz_artist_id}`)

const artistArtUrl = computed(() => {
  return getCoverArtUrl(musicbrainz_artist_id.value)
})

async function getData() {
  artist.value = await fetchArtist(musicbrainz_artist_id.value)

  const albumTracks: SubsonicSong[] = []
  const artistAlbums: SubsonicAlbum[] = []
  artist.value.album.forEach(async (album) => {
    const albumWithSongs = await fetchAlbum(album.id)
    artistAlbums.push(albumWithSongs)
    for (const track of albumWithSongs.song) {
      albumTracks.push(track)
    }
  })

  albums.value = artistAlbums
  tracks.value = albumTracks
  routeTracks.value = albumTracks
}

onBeforeMount(async () => {
  await getData()
})
</script>

<template>
  <section v-if="artist" class="h-80 rounded-lg">
    <div
      class="h-full w-full bg-cover bg-center"
      :style="{ backgroundImage: `url(${artistArtUrl})` }"
    >
      <div class="h-full w-full flex items-center justify-center gap-6 align-middle backdrop-blur-md">
        <div class="mx-auto flex items-center justify-center gap-6 rounded-lg bg-black/40 p-4 align-middle">
          <div class="size-60">
            <img
              class="h-full w-full rounded-md object-cover"
              :src="artistArtUrl"
              @error="onImageError"
            />
          </div>
          <div class="text-7xl text-gray-300 font-bold">
            {{ artist.name }}
          </div>
        </div>
      </div>
    </div>
  </section>
  <div>
    <h2 class="text-lg font-semibold">
      Albums
    </h2>
    <div class="flex flex-wrap gap-6">
      <div v-for="album in albums" :key="album.id" class="flex flex-col gap-y-1 overflow-hidden transition duration-200 hover:scale-110">
        <Album :album="album" size="lg" />
      </div>
    </div>
  </div>
  <Tracks v-if="tracks" :tracks="tracks" :show-album="true" />
</template>
