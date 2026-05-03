<script setup lang="ts">
import type { SubsonicArtist } from '~/types/subsonicArtist'
import { postNewArtistArt, useServerSentEventsForArtistArt } from '~/logic/backendFetch'

const props = defineProps({
  artist: { type: Object as PropType<SubsonicArtist>, required: true },
})

const emits = defineEmits(['close', 'artUpdated'])

const loading = ref(true)
const deezerArtUrl = ref<string | null>(null)
const coverArtArchiveUrl = ref<string | null>(null)
const localFolderArtUrl = ref<string | null>(null)
const artistArt = ref<string | null>(null)

async function updateArt(source: 'deezer' | 'coverartarchive' | 'manual' | 'localfolder') {
  let artUrl: string | null = null
  switch (source) {
    case 'deezer':
      artUrl = deezerArtUrl.value
      break
    case 'coverartarchive':
      artUrl = coverArtArchiveUrl.value
      break
    case 'manual':
      artUrl = artistArt.value
      break
    case 'localfolder':
      artUrl = localFolderArtUrl.value
      break
  }
  if (artUrl) {
    const imageBlob = await (await fetch(artUrl)).blob()
    const response = await postNewArtistArt(props.artist.id, imageBlob)
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
  }
}

function onErrorReceived(error: any) {
  console.error('SSE Error Received:', error)
}

onMounted(() => {
  useServerSentEventsForArtistArt(props.artist.id, onMessageReceived, onErrorReceived)
})
</script>

<template>
  <Modal :show-modal="true" modal-title="Change Artist Art" @close="$emit('close')">
    <template #content>
      <Loading v-if="loading" class="h-56" />
      <div v-else class="flex flex-wrap gap-4 justify-center">
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
          v-if="artistArt"
          :image-url="artistArt"
          label="Custom"
          type="manual"
          @update-art="updateArt"
        />
      </div>
      <ImageSelector v-model="artistArt" />
    </template>
  </Modal>
</template>
