<script setup lang="ts">
import type { SubsonicAlbum } from '../types/subsonicAlbum'
import { postNewAlbumArt, useServerSentEventsForAlbumArt } from '~/composables/backendFetch'

const props = defineProps({
  album: { type: Object as PropType<SubsonicAlbum>, required: true },
})

const emits = defineEmits(['close', 'artUpdated'])

const loading = ref(true)
const deezerArtUrl = ref<string | null>(null)
const coverArtArchiveUrl = ref<string | null>(null)
const localFolderArtUrl = ref<string | null>(null)
const localEmbeddedArtUrl = ref<string | null>(null)
const albumArt = ref<string | null>(null)

// async function getAlbumArtUrls() {
//   const options = await fetchAlbumArtOptions(props.album.albumArtists[0].name, props.album.name)
//   deezerArtUrl.value = options.deezer
//   coverArtArchiveUrl.value = options.cover_art_archive
//   localFolderArtUrl.value = options.local_folder_art
//   localEmbeddedArtUrl.value = options.local_embedded_art
// }

async function updateArt(source: 'deezer' | 'coverartarchive' | 'manual' | 'localfolder' | 'localembedded') {
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
    case 'localfolder':
      artUrl = localFolderArtUrl.value
      break
    case 'localembedded':
      artUrl = localEmbeddedArtUrl.value
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

function onMessageReceived(data: any) {
  loading.value = false
  if (data.source === 'Deezer') {
    deezerArtUrl.value = data.data
  }
  else if (data.source === 'CoverArtArchive') {
    coverArtArchiveUrl.value = data.data
  }
  else if (data.source === 'LocalArt') {
    localFolderArtUrl.value = data.data.folderArt
    localEmbeddedArtUrl.value = data.data.embeddedArt
  }
}

function onErrorReceived(error: any) {
  console.error('SSE Error Received:', error)
}

onMounted(() => {
  useServerSentEventsForAlbumArt(props.album.albumArtists[0].name, props.album.name, onMessageReceived, onErrorReceived)
})
</script>

<template>
  <teleport to="body">
    <div class="fixed inset-0 z-50 flex items-center justify-center backdrop-blur-lg">
      <div class="relative w-full flex flex-col gap-4 border-1 border-zshade-500 border-solid background-3 p-4 lg:w-80dvw">
        <div class="flex flex-row items-center justify-center gap-4">
          <ZButton aria-label="Close" @click="$emit('close')">
            X
          </ZButton>
          <p class="text-lg text-primary font-bold">
            Change Album Art
          </p>
          <div />
        </div>
        <Loading v-if="loading" class="h-56" />
        <div v-else class="flex flex-wrap justify-center gap-4">
          <ImageSelectorImage
            v-if="deezerArtUrl"
            :image-url="deezerArtUrl"
            label="Deezer"
            type="deezer"
            @update-art="updateArt"
          />
          <ImageSelectorImage
            v-if="coverArtArchiveUrl"
            :image-url="coverArtArchiveUrl"
            label="Cover Art Archive"
            type="coverartarchive"
            @update-art="updateArt"
          />
          <ImageSelectorImage
            v-if="localFolderArtUrl"
            :image-url="localFolderArtUrl"
            label="Album folder"
            type="localfolder"
            @update-art="updateArt"
          />
          <ImageSelectorImage
            v-if="localEmbeddedArtUrl"
            :image-url="localEmbeddedArtUrl"
            label="Embedded"
            type="localembedded"
            @update-art="updateArt"
          />
          <ImageSelectorImage
            v-if="albumArt"
            :image-url="albumArt"
            label="Custom"
            type="manual"
            @update-art="updateArt"
          />
        </div>
        <ImageSelector v-model="albumArt" />
      </div>
    </div>
  </teleport>
</template>
