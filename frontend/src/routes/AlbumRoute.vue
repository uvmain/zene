<script setup lang="ts">
import type { AlbumMetadata, TrackMetadataWithImageUrl } from '../types'
import { useBackendFetch } from '../composables/useBackendFetch'
import { usePlaycounts } from '../composables/usePlaycounts'
import { useRouteTracks } from '../composables/useRouteTracks'

const route = useRoute()
const { routeTracks, clearRouteTracks } = useRouteTracks()
const { backendFetchRequest, getAlbumTracks } = useBackendFetch()
const { last_updated_musicbrainz_track_id } = usePlaycounts()

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

async function getAlbumTracksAndRouteTracks() {
  const response = await getAlbumTracks(musicbrainz_album_id.value)
  tracks.value = response
  routeTracks.value = response
}

watch(() => route.params.musicbrainz_album_id, async () => {
  getAlbum()
  getAlbumTracksAndRouteTracks()
})

watch(last_updated_musicbrainz_track_id, (newTrack) => {
  tracks.value?.forEach((track) => {
    if (track.musicbrainz_track_id === newTrack) {
      track.user_play_count = track.user_play_count + 1
      track.global_play_count = track.global_play_count + 1
    }
  })
})

onBeforeMount(async () => {
  await getAlbum()
  await getAlbumTracksAndRouteTracks()
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
