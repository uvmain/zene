<script setup lang="ts">
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import type { SubsonicSong } from '~/types/subsonicSong'
import { fetchAlbum } from '~/logic/backendFetch'
import { clearRouteTracks } from '~/logic/routeTracks'
import { routeTracks } from '~/logic/store'

const route = useRoute('/albums/[album]')

const album = ref<SubsonicAlbum>()
const tracks = ref<SubsonicSong[]>()
const musicbrainzAlbumId = computed(() => `${route.params.album}`)

async function getAlbum() {
  const response = await fetchAlbum(musicbrainzAlbumId.value)
  album.value = response
  tracks.value = response.song
  routeTracks.value = response.song
}

watch(() => route.params.album, async () => {
  getAlbum()
})

onBeforeMount(async () => {
  await getAlbum()
})

onUnmounted(() => clearRouteTracks())
</script>

<template>
  <div v-if="album && tracks">
    <div class="flex flex-grow flex-col gap-4 lg:gap-6">
      <HeroAlbum :album="album" />
      <Tracks :tracks="tracks" :primary-artist="album.artist" />
    </div>
  </div>
</template>
