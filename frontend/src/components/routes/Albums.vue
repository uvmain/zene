<script setup lang="ts">
import { backendFetchRequest } from '../../composables/fetchFromBackend'

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
    <h2 class="mb-2 text-lg font-semibold">
      Recently Added Albums
    </h2>
    <div class="flex flex-wrap gap-6">
      <div v-for="album in recentlyAddedAlbums" :key="album.album" class="w-30 flex flex-col gap-y-1 overflow-hidden transition duration-200 hover:scale-110">
        <Album :album="album" size="lg" />
      </div>
    </div>
  </div>
</template>
