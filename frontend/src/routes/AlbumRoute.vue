<script setup lang="ts">
import type { AlbumMetadata, TrackMetadata } from '../types'
import { backendFetchRequest } from '../composables/fetchFromBackend'
import { formatTime } from '../composables/logic'

const route = useRoute()
const album = ref<AlbumMetadata>()
const tracks = ref<TrackMetadata[]>()
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
      <div class="rounded-lg bg-black/20 p-4">
        <div class="flex flex-row gap-2">
          <div class="w-15 flex justify-center">
            #
          </div>
          <div class="flex flex-grow">
            Title
          </div>
          <icon-tabler-clock-hour-3 class="w-15 flex justify-center" />
        </div>
        <hr class="my-4 border-white/10" />
        <div class="w-full flex flex-col gap-2">
          <div
            v-for="track in tracks"
            :key="track.title"
            class="group w-full flex flex-none flex-row overflow-hidden rounded p-1 duration-200 transition-ease-out hover:bg-zene-200/20"
          >
            <div class="w-15 flex items-center justify-center">
              <span class="group-hover:hidden">
                {{ track.track_number }}
              </span>
              <icon-tabler-player-play-filled class="hidden text-xl group-hover:block" />
            </div>
            <div class="flex flex-grow flex-col gap-1">
              <div class="font-semibold group-hover:ml-1">
                {{ track.title }}
              </div>
              <div class="text-sm">
                {{ track.artist }}
              </div>
            </div>
            <div class="w-15 flex items-center justify-center">
              {{ formatTime(Number.parseInt(track.duration)) }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
