<script setup lang="ts">
import type { ArtSseMessage, PostArtOptions } from '~/logic/backendFetch'
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
const showError = ref<string | null>(null)

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
    const options: PostArtOptions = {
      musicbrainz_id: props.artist.id,
    }
    if (artUrl.startsWith('http')) {
      options.url = artUrl
    }
    else {
      options.image = await (await fetch(artUrl)).blob()
    }
    const response = await postNewArtistArt(options)
    if (response.status === 'ok') {
      emits('artUpdated', artUrl)
    }
  }
}

function onMessageReceived(data: ArtSseMessage) {
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

function onCompleteReceived() {
  loading.value = false
  if (`${deezerArtUrl.value}${coverArtArchiveUrl.value}${localFolderArtUrl.value}` === '') {
    showError.value = 'No artist art options found.'
  }
}

function onErrorReceived(error: any) {
  console.error('SSE Error Received:', error)
  showError.value = 'An error occurred while fetching artist art options.'
  loading.value = false
}

onMounted(() => {
  useServerSentEventsForArtistArt(props.artist.id, onMessageReceived, onErrorReceived, onCompleteReceived)
})
</script>

<template>
  <Modal :show-modal="true" modal-title="Change Artist Art" @close="$emit('close')">
    <template #content>
      <Loading v-if="loading" class="h-56" />
      <div v-if="showError" class="text-primary-400 text-center">
        {{ showError }}
      </div>
      <div class="flex flex-wrap gap-4 justify-center">
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
