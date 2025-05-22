<script setup lang="ts">
import type { AlbumMetadata, TrackMetadataWithImageUrl } from '../types'
import { useDark, useSessionStorage, useToggle } from '@vueuse/core'
import dayjs from 'dayjs'
import { backendFetchRequest } from '../composables/fetchFromBackend'

const inputText = useSessionStorage<string>('searchInput', '')
const isDark = useDark()
const toggleDark = useToggle(isDark)
const searchResults = ref<TrackMetadataWithImageUrl[]>([])

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
  <header>
    <div class="flex p-2">
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
    </div>
    <div v-if="inputText.length > 3" class="z-100 h-full w-full rounded-lg bg-zene-600/70 px-2">
      <div class="flex flex-col gap-2 p-4">
        <h3>
          Search results:
        </h3>
        <div>
          Albums: {{ searchResultsAlbums.length }}
        </div>
        <div class="flex flex-wrap gap-6">
          <div v-for="album in searchResultsAlbums" :key="album.album" class="w-30 flex flex-col gap-y-1 overflow-hidden">
            <Album :album="album" size="lg" />
          </div>
        </div>
      </div>
    </div>
  </header>
</template>
