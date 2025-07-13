<script setup lang="ts">
import dayjs from 'dayjs'
import { useBackendFetch } from '../composables/useBackendFetch'

const props = defineProps({
  limit: { type: Number, default: 30 },
})

const { backendFetchRequest } = useBackendFetch()

const recentlyAddedAlbums = ref()

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
  recentlyAddedAlbums.value = albums
}

onBeforeMount(async () => {
  await getAlbums()
})
</script>

<template>
  <div>
    <RefreshHeader title="Recently Added Albums" @refreshed="getAlbums()" />
    <div class="flex flex-wrap justify-center gap-6 md:justify-start">
      <div v-for="album in recentlyAddedAlbums" :key="album.album" class="flex flex-col gap-y-1 overflow-hidden transition duration-200 hover:scale-110">
        <Album :album="album" size="lg" />
      </div>
    </div>
  </div>
</template>
