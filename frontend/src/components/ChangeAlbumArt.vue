<script setup lang="ts">
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import { postNewAlbumArt, useServerSentEventsForAlbumArt } from '~/logic/backendFetch'

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
  <Modal :show-modal="true" modal-title="Change Album Art" @close="$emit('close')">
    <template #content>
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
    </template>
  </Modal>
</template>
