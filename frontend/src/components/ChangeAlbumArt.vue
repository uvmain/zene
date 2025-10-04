<script setup lang="ts">
import type { SubsonicAlbum } from '../types/subsonicAlbum'
import { fetchAlbumArtOptions, postNewAlbumArt } from '~/composables/backendFetch'

const props = defineProps({
  album: { type: Object as PropType<SubsonicAlbum>, required: true },
})

const emits = defineEmits(['close', 'artUpdated'])

const loading = ref(true)
const deezerArtUrl = ref<string | null>(null)
const coverArtArchiveUrl = ref<string | null>(null)
const albumArt = ref<string | null>(null)

async function getAlbumArtUrls() {
  const options = await fetchAlbumArtOptions(props.album.albumArtists[0].name, props.album.name)
  deezerArtUrl.value = options.deezer
  coverArtArchiveUrl.value = options.cover_art_archive
}

async function updateArt(source: 'deezer' | 'coverartarchive' | 'manual') {
  let artUrl: string | null = null
  switch (source) {
    case 'deezer':
      artUrl = deezerArtUrl.value
      break
    case 'coverartarchive':
      artUrl = coverArtArchiveUrl.value
      break
    case 'manual':
      artUrl = albumArt.value
      break
  }
  if (artUrl) {
    const imageBlob = await (await fetch(artUrl)).blob()
    const response = await postNewAlbumArt(props.album.id, imageBlob)
    if (response.status === 'ok') {
      emits('artUpdated', artUrl)
    }
  }
}

onMounted(async () => {
  await getAlbumArtUrls()
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
        <div v-else class="flex flex-wrap justify-center gap-4">
          <div v-if="deezerArtUrl" class="relative size-56">
            <img
              class="size-56"
              :src="deezerArtUrl"
              alt="Deezer Album Art"
            />
            <button
              class="z-button absolute bottom-2 right-2"
              aria-label="Choose art"
              @click="updateArt('deezer')"
            >
              Use This Art
            </button>
          </div>
          <div v-if="coverArtArchiveUrl" class="relative size-56">
            <img
              class="size-56"
              :src="coverArtArchiveUrl"
              alt="Cover Art Archive"
            />
            <button
              class="z-button absolute bottom-2 right-2"
              aria-label="Choose art"
              @click="updateArt('coverartarchive')"
            >
              Use This Art
            </button>
          </div>
          <div v-if="albumArt" class="relative size-56">
            <img
              class="size-56"
              :src="albumArt"
              alt="Album Art"
            />
            <button
              class="z-button absolute bottom-2 right-2"
              aria-label="Choose art"
              @click="updateArt('manual')"
            >
              Use This Art
            </button>
          </div>
        </div>
        <ImageSelector v-model="albumArt" />
      </div>
    </div>
  </teleport>
</template>
