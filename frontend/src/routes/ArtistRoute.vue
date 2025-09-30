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
const albums = ref<SubsonicAlbum[]>()

const musicbrainz_artist_id = computed(() => `${route.params.musicbrainz_artist_id}`)

const artistArtUrl = computed(() => {
  return getCoverArtUrl(musicbrainz_artist_id.value)
})

async function getData() {
  const promisesArray = [
    fetchArtist(musicbrainz_artist_id.value),
    fetchAlbumsForArtist(musicbrainz_artist_id.value),
    fetchArtistTopSongs(musicbrainz_artist_id.value),
  ]

  await Promise.all(promisesArray)
    .then(
      (results) => {
        artist.value = results[0] as SubsonicArtist
        albums.value = results[1] as SubsonicAlbum[]
        tracks.value = results[2] as SubsonicSong[]
        routeTracks.value = tracks.value
      },
    )
}

onBeforeMount(async () => {
  await getData()
})
</script>

<template>
  <section v-if="artist" class="h-80">
    <div
      class="corner-cut-large h-full w-full bg-cover bg-center"
      :style="{ backgroundImage: `url(${artistArtUrl})` }"
    >
      <div class="h-full w-full flex items-center justify-center gap-6 align-middle backdrop-blur-lg">
        <div class="w-full flex items-center justify-center gap-6 background-grad-2 p-4 align-middle">
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
