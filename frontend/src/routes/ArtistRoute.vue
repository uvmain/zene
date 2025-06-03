<script setup lang="ts">
import type { ArtistMetadata, TrackMetadataWithImageUrl } from '../types'
import { backendFetchRequest, getArtistTracks } from '../composables/fetchFromBackend'

const route = useRoute()
const artist = ref<ArtistMetadata>()
const tracks = ref<TrackMetadataWithImageUrl[]>()

const musicbrainz_artist_id = computed(() => `${route.params.musicbrainz_artist_id}`)

async function getArtist() {
  const response = await backendFetchRequest(`artists/${musicbrainz_artist_id.value}`)
  const json = await response.json() as ArtistMetadata
  artist.value = json
}

async function getTracks() {
  tracks.value = await getArtistTracks(musicbrainz_artist_id.value)
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
  <Tracks v-if="tracks" :tracks="tracks" :show-album="true" />
</template>
