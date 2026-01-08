<script setup lang="ts">
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import type { SubsonicSong } from '~/types/subsonicSong'
import { fetchAlbum } from '~/composables/backendFetch'
import { useRouteTracks } from '~/composables/useRouteTracks'

const route = useRoute<'/albums/[album]'>()
const { routeTracks, clearRouteTracks } = useRouteTracks()

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
      <Album :album="album" size="md" :show-change-art-button="true" />
      <Tracks :tracks="tracks" />
    </div>
  </div>
</template>
