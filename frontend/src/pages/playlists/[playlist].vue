<script setup lang="ts">
import type { SubsonicPlaylistResponse } from '~/types/subsonic'
import type { SubsonicPlaylist } from '~/types/subsonicPlaylists'
import { openSubsonicFetchRequest } from '~/composables/backendFetch'
import { getCoverArtUrl, onImageError } from '~/composables/logic'

const route = useRoute()
const playlistId = computed(() => `${route.params.playlist}`)

const playlist = ref<SubsonicPlaylist>()

async function getPlaylist() {
  const formData = new FormData()
  formData.append('id', playlistId.value)
  const response = await openSubsonicFetchRequest<SubsonicPlaylistResponse>('getPlaylist', {
    body: formData,
  })
  playlist.value = response?.playlist
  playlist.value.coverArt = getCoverArtUrl(playlist.value.coverArt, 200)
}

onBeforeMount(getPlaylist)
</script>

<template>
  <div class="p-4 space-y-4">
    <div v-if="playlist" class="flex flex-wrap justify-center gap-6 md:justify-start">
      <div>
        <img
          :src="playlist.coverArt"
          alt="Playlist Cover"
          class="size-40 object-cover"
          loading="lazy"
          width="200"
          height="200"
          @error="onImageError"
        />
        <div class="my-auto flex flex-col gap-4">
          <div class="text-2xl font-bold">
            {{ playlist.name }}
          </div>
        </div>
      </div>
      <Tracks :show-album="true" :tracks="playlist.entry" />
    </div>
  </div>
</template>
