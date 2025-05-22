<script setup lang="ts">
import type { AlbumMetadata, TrackMetadataWithImageUrl } from '../types'
import { useDark, useSessionStorage, useToggle } from '@vueuse/core'
import dayjs from 'dayjs'
import { backendFetchRequest } from '../composables/fetchFromBackend'

const inputText = useSessionStorage<string>('searchInput', '')
const isDark = useDark()
const toggleDark = useToggle(isDark)
const searchResults = ref<TrackMetadataWithImageUrl[]>([])
const searchDialog = ref<HTMLDialogElement | null>(null)

async function search() {
  if (!inputText.value || inputText.value.length < 3) {
    searchResults.value = []
    return
  }
  const response = await backendFetchRequest(`search?search=${inputText.value}`)
  const json = await response.json() as TrackMetadataWithImageUrl[]
  const albumMetadata: TrackMetadataWithImageUrl[] = []
  json.forEach((metadata: any) => {
    const metadataInstance = {
      id: metadata.id,
      file_id: metadata.file_id,
      filename: metadata.filename,
      format: metadata.format,
      duration: metadata.duration,
      size: metadata.size,
      bitrate: metadata.bitrate,
      title: metadata.title,
      artist: metadata.artist,
      album: metadata.album,
      album_artist: metadata.album_artist,
      track_number: metadata.track_number,
      total_tracks: metadata.total_tracks,
      disc_number: metadata.disc_number,
      total_discs: metadata.total_discs,
      musicbrainz_artist_id: metadata.musicbrainz_artist_id,
      musicbrainz_album_id: metadata.musicbrainz_album_id,
      musicbrainz_track_id: metadata.musicbrainz_track_id,
      label: metadata.label,
      genres: metadata.genre.split(';').filter((genre: string) => genre !== ''),
      release_date: dayjs(metadata.release_date).format('YYYY'),
      image_url: `/api/art/albums/${metadata.musicbrainz_album_id}?size=lg`,
    }
    albumMetadata.push(metadataInstance)
  })
  searchResults.value = albumMetadata
  if (searchResults.value.length > 0) {
    searchDialog.value?.showModal()
  }
  else {
    searchDialog.value?.close()
  }
}

const searchResultsAlbums = computed(() => {
  const uniqueAlbums = new Map<string, AlbumMetadata>()
  searchResults.value.forEach((album: TrackMetadataWithImageUrl) => {
    if (!uniqueAlbums.has(album.musicbrainz_album_id)) {
      uniqueAlbums.set(album.musicbrainz_album_id, {
        artist: album.artist,
        album: album.album,
        musicbrainz_track_id: album.musicbrainz_track_id,
        musicbrainz_album_id: album.musicbrainz_album_id,
        musicbrainz_artist_id: album.musicbrainz_artist_id,
        genres: album.genres,
        release_date: album.release_date,
        image_url: album.image_url,
      })
    }
  })
  return Array.from(uniqueAlbums.values())
})
</script>

<template>
  <header class="flex p-2">
    <div class="flex flex-grow justify-center">
      <input
        id="search-input"
        v-model="inputText"
        placeholder="Type here to search"
        type="text"
        class="block w-1/2 border border-zene-400 rounded rounded-lg bg-gray-800 px-4 py-2 text-white focus-border-1 focus:border-zene-200 focus:border-solid focus:shadow-zene-400 hover:shadow-lg focus:outline-none"
        @change="search()"
        @input="search()"
      >
    </div>
    <div id="user-and-settings" class="flex gap-4">
      <div class="hover:cursor-pointer" @click="toggleDark()">
        <icon-tabler-sun v-if="isDark" class="text-2xl" />
        <icon-tabler-moon-stars v-else class="text-2xl" />
      </div>
      <div class="hover:cursor-pointer">
        <icon-tabler-user class="text-2xl" />
      </div>
    </div>
    <dialog ref="searchDialog" class="absolute px-2 backdrop-blur-md" data-cy="modal-dialog">
      <div class="rounded-lg bg-white pb-4 md:pb-8">
        <div class="flex gap-6">
          <div v-for="album in searchResultsAlbums" :key="album.album" class="w-30 flex flex-col gap-y-1 overflow-hidden">
            <Album :album="album" size="lg" />
          </div>
        </div>
        <button class="bg-university-400 mx-auto mt-4 block rounded px-6 py-3 text-white" data-cy="close" @click="searchDialog?.close()">
          Close
        </button>
      </div>
    </dialog>
  </header>
</template>
