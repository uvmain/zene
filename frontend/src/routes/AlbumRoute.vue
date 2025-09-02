<script setup lang="ts">
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import type { SubsonicSong } from '~/types/subsonicSong'
import { fetchAlbum } from '~/composables/backendFetch'
import { useRouteTracks } from '~/composables/useRouteTracks'

const route = useRoute()
const { routeTracks, clearRouteTracks } = useRouteTracks()

const album = ref<SubsonicAlbum>()
const tracks = ref<SubsonicSong[]>()
const musicbrainz_album_id = computed(() => `${route.params.musicbrainz_album_id}`)

async function getAlbum() {
  const response = await fetchAlbum(musicbrainz_album_id.value)
  album.value = response
  tracks.value = response.song
  routeTracks.value = response.song
}

watch(() => route.params.musicbrainz_album_id, async () => {
  getAlbum()
})

onBeforeMount(async () => {
  await getAlbum()
})

onUnmounted(() => clearRouteTracks())
</script>

<template>
  <div v-if="album && tracks">
    <div class="flex flex-grow flex-col gap-6">
      <Album :album="album" size="xl" class="rounded-lg" />
      <Tracks :tracks="tracks" />
    </div>
  </div>
</template>
