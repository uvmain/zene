<script setup lang="ts">
import type { AlbumMetadata, ArtistMetadata, TrackMetadataWithImageUrl } from '../types'
import dayjs from 'dayjs'
import { backendFetchRequest, getArtistAlbums, getArtistTracks } from '../composables/fetchFromBackend'

const route = useRoute()
const artist = ref<ArtistMetadata>()
const tracks = ref<TrackMetadataWithImageUrl[]>()
const albums = ref<AlbumMetadata[]>()

const musicbrainz_artist_id = computed(() => `${route.params.musicbrainz_artist_id}`)

async function getArtist() {
  const response = await backendFetchRequest(`artists/${musicbrainz_artist_id.value}`)
  const json = await response.json() as ArtistMetadata
  artist.value = json
}

async function getTracks() {
  tracks.value = await getArtistTracks(musicbrainz_artist_id.value, 30)
}

async function getAlbums() {
  const json = await getArtistAlbums(musicbrainz_artist_id.value)
  const albumMetadata: AlbumMetadata[] = []
  json.forEach((metadata: any) => {
    const metadataInstance = {
      artist: metadata.artist,
      album: metadata.album,
      album_artist: metadata.album_artist,
      musicbrainz_album_id: metadata.musicbrainz_album_id as string,
      musicbrainz_artist_id: metadata.musicbrainz_artist_id,
      genres: metadata.genres.split(';').filter((genre: string) => genre !== ''),
      release_date: dayjs(metadata.release_date).format('YYYY'),
      image_url: `/api/albums/${metadata.musicbrainz_album_id}/art?size=xl`,
    }
    albumMetadata.push(metadataInstance)
  })
  albums.value = albumMetadata
}

function onImageError(event: Event) {
  const target = event.target as HTMLImageElement
  target.onerror = null
  target.src = '/default-square.png'
}

onBeforeMount(async () => {
  await getArtist()
})

onMounted(async () => {
  await getTracks()
  await getAlbums()
})
</script>

<template>
  <section v-if="artist" class="h-80 rounded-lg">
    <div
      class="h-full w-full bg-cover bg-center"
      :style="{ backgroundImage: `url(${artist.image_url})` }"
    >
      <div class="h-full w-full flex items-center justify-center gap-6 align-middle backdrop-blur-md">
        <div class="mx-auto flex items-center justify-center gap-6 rounded-lg bg-black/40 p-4 align-middle">
          <div class="size-60">
            <img
              class="h-full w-full rounded-md object-cover"
              :src="artist.image_url"
              @error="onImageError"
            />
          </div>
          <div class="text-7xl text-gray-300 font-bold">
            {{ artist.artist }}
          </div>
        </div>
      </div>
    </div>
  </section>
  <div>
    <h2 class="text-lg font-semibold">
      Albums
    </h2>
    <div class="flex flex-wrap gap-6">
      <div v-for="album in albums" :key="album.album" class="flex flex-col gap-y-1 overflow-hidden transition duration-200 hover:scale-110">
        <Album :album="album" size="lg" />
      </div>
    </div>
  </div>
  <Tracks v-if="tracks" :tracks="tracks" :show-album="true" />
</template>
