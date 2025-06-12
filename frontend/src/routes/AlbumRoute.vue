<script setup lang="ts">
import type { AlbumMetadata, TrackMetadataWithImageUrl } from '../types'
import { backendFetchRequest } from '../composables/fetchFromBackend'

const route = useRoute()
const album = ref<AlbumMetadata>()
const tracks = ref<TrackMetadataWithImageUrl[]>()
const musicbrainz_album_id = computed(() => `${route.params.musicbrainz_album_id}`)

async function getAlbum() {
  const response = await backendFetchRequest(`albums/${musicbrainz_album_id.value}`)
  const json = await response.json()
  album.value = {
    album: json.album,
    artist: json.artist,
    album_artist: json.album_artist ?? json.artist,
    musicbrainz_album_id: json.musicbrainz_album_id,
    musicbrainz_artist_id: json.musicbrainz_artist_id,
    release_date: json.release_date,
    genres: json.genres ? json.genres.split(';').filter((genre: string) => genre !== '') : [],
    image_url: `/api/albums/${json.musicbrainz_album_id}/art?size=lg`,
  }
}

async function getAlbumTracks() {
  const response = await backendFetchRequest(`albums/${musicbrainz_album_id.value}/tracks`)
  const json = await response.json()
  tracks.value = json
}

watch(() => route.params.musicbrainz_album_id, async () => {
  getAlbum()
  getAlbumTracks()
})

onBeforeMount(async () => {
  await getAlbum()
  await getAlbumTracks()
})
</script>

<template>
  <div v-if="album && tracks">
    <div class="flex flex-grow flex-col gap-6">
      <Album :album="album" size="xl" class="rounded-lg" />
      <Tracks :tracks="tracks" />
    </div>
  </div>
</template>
