<script setup lang="ts">
import dayjs from 'dayjs'
import { useBackendFetch } from '../composables/useBackendFetch'

const props = defineProps({
  limit: { type: Number, default: 30 },
})

const { backendFetchRequest } = useBackendFetch()

const recentlyAddedAlbums = ref()
const refreshed = ref(false)

async function getAlbums() {
  const response = await backendFetchRequest(`albums?recent=true&limit=${props.limit}`)
  const json = await response.json()
  const albums = json.map((album: any) => ({
    album: album.album,
    artist: album.artist,
    album_artist: album.album_artist ?? album.artist,
    musicbrainz_album_id: album.musicbrainz_album_id,
    release_date: dayjs(album.release_date).format('YYYY'),
    image_url: `/api/albums/${album.musicbrainz_album_id}/art?size=lg`,
  }))
  refreshed.value = true
  setTimeout(() => {
    refreshed.value = false
  }, 1000)
  recentlyAddedAlbums.value = albums
}

onBeforeMount(async () => {
  await getAlbums()
})
</script>

<template>
  <div>
    <div class="flex flex-row items-center gap-x-2 py-2">
      <h2 class="text-lg font-semibold">
        Recently Added Albums
      </h2>
      <icon-tabler-refresh class="cursor-pointer text-sm" :class="{ spin: refreshed }" @click="getAlbums" />
    </div>
    <div class="flex flex-wrap gap-6">
      <div v-for="album in recentlyAddedAlbums" :key="album.album" class="flex flex-col gap-y-1 overflow-hidden transition duration-200 hover:scale-110">
        <Album :album="album" size="lg" />
      </div>
    </div>
  </div>
</template>

<style scoped>
@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(-360deg);
  }
}

.spin {
  animation: spin 0.5s linear;
}
</style>
