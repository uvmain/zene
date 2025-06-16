<script setup lang="ts">
import { useBackendFetch } from '../composables/useBackendFetch'

const { backendFetchRequest } = useBackendFetch()

const recentlyAddedAlbums = ref()

async function getAlbums() {
  const response = await backendFetchRequest('albums?recent=true&limit=60')
  const json = await response.json()
  const albums = json.map((album: any) => ({
    album: album.album,
    artist: album.artist,
    album_artist: album.album_artist ?? album.artist,
    musicbrainz_album_id: album.musicbrainz_album_id,
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
    <RecentlyAddedAlbums :limit="100" />
  </div>
</template>
