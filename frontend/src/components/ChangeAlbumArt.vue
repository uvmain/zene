<script setup lang="ts">
import type { SubsonicAlbum } from '../types/subsonicAlbum'
import { fetchDeezerArtUrl, postNewAlbumArt } from '~/composables/backendFetch'

const props = defineProps({
  album: { type: Object as PropType<SubsonicAlbum>, required: true },
})

const emits = defineEmits(['close', 'artUpdated'])

const loading = ref(true)
const deezerArtUrl = ref<string | null>(null)

async function getDeezerUrl() {
  deezerArtUrl.value = await fetchDeezerArtUrl(props.album.albumArtists[0].name, props.album.name)
}

async function updateArt() {
  if (deezerArtUrl.value) {
    const imageBlob = await (await fetch(deezerArtUrl.value)).blob()
    const response = await postNewAlbumArt(props.album.id, imageBlob)
    if (response.status === 'ok') {
      emits('artUpdated', deezerArtUrl.value)
    }
  }
}

onMounted(async () => {
  await getDeezerUrl()
  loading.value = false
})
</script>

<template>
  <teleport to="body">
    <div class="fixed inset-0 z-50 flex items-center justify-center backdrop-blur-lg">
      <div class="relative w-full flex flex-col gap-4 border-1 border-zshade-400 border-solid background-3 p-4 lg:w-60dvw">
        <div class="flex flex-row items-center gap-4">
          <button class="z-button" aria-label="Close" @click="$emit('close')">
            X
          </button>
          <p class="text-lg text-primary font-bold">
            Change Album Art
          </p>
          <div />
        </div>
        <svg v-if="loading" class="text-primary" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
          <path fill="currentColor" d="M12 2A10 10 0 1 0 22 12A10 10 0 0 0 12 2Zm0 18a8 8 0 1 1 8-8A8 8 0 0 1 12 20Z" opacity="0.5" /><path fill="currentColor" d="M20 12h2A10 10 0 0 0 12 2V4A8 8 0 0 1 20 12Z"><animateTransform attributeName="transform" dur="1s" from="0 12 12" repeatCount="indefinite" to="360 12 12" type="rotate" /></path>
        </svg>
        <div v-if="deezerArtUrl" class="relative size-56">
          <img
            class="size-56"
            :src="deezerArtUrl"
            alt="Deezer Album Art"
          />
          <button
            class="z-button absolute bottom-2 right-2"
            aria-label="Choose art"
            @click="updateArt()"
          >
            Use This Art
          </button>
        </div>
      </div>
    </div>
  </teleport>
</template>
