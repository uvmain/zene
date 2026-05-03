<script setup lang="ts">
import type { PostArtOptions } from '~/logic/backendFetch'
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
const showError = ref<string | null>(null)

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
    const options: PostArtOptions = {
      musicbrainz_id: props.album.id,
    }
    if (artUrl.startsWith('http')) {
      options.url = artUrl
    }
    else {
      options.image = await (await fetch(artUrl)).blob()
    }
    const response = await postNewAlbumArt(options)
    if (response.status === 'ok') {
      emits('artUpdated', artUrl)
    }
  }
}

function onMessageReceived(data: any) {
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
  showError.value = 'An error occurred while fetching album art options.'
  loading.value = false
}

function onCompleteReceived() {
  loading.value = false
  if (`${deezerArtUrl.value}${coverArtArchiveUrl.value}${localFolderArtUrl.value}${localEmbeddedArtUrl.value}` === '') {
    showError.value = 'No album art options found.'
  }
}

onMounted(() => {
  useServerSentEventsForAlbumArt(props.album.id, onMessageReceived, onErrorReceived, onCompleteReceived)
})
</script>

<template>
  <Modal :show-modal="true" modal-title="Change Album Art" @close="$emit('close')">
    <template #content>
      <Loading v-if="loading" class="h-10" />
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
